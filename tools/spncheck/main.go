// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

// Command spncheck audits every Terraform resource/data source service package
// and verifies that its declared service-principal support (the IsSPNSupported
// field in base.go's ItemTypeInfo) matches the actual service-principal support
// of the Microsoft Fabric Go SDK APIs that the package calls.
//
// An item supports service principal (IsSPNSupported = true) only when ALL of
// the fabric-sdk-go client functions it invokes support service principal. Each
// SDK client function documents its supported identities in a table under the
// following heading in the doc comment directly above the function:
//
//	MICROSOFT ENTRA SUPPORTED IDENTITIES
//	| Identity | Support | |-|-| | User | Yes | | Service principal ... | Yes |
//
// The "Service principal" row's Support cell is one of:
//
//	Yes  -> the API supports service principal
//	No   -> the API does NOT support service principal
//	<conditional sentence> -> supported only under certain conditions
//
// spncheck flags an item when ANY called SDK API has a hard "No" in its Service
// principal cell: that item cannot fully support service principal. Because the
// SDK's identity tables are not always present, findings are split by confidence:
//
//   - OVER-MARKED (high confidence): declared IsSPNSupported = true but a called
//     API is documented as NOT supporting service principal. Should be false.
//   - REVIEW (low confidence): declared IsSPNSupported = false but every called
//     API supports service principal. Possibly promotable to true, but requires
//     manual confirmation because SDK annotations are incomplete.
//   - UNDETERMINED: no fabric-sdk-go calls found in the package (e.g. generic
//     fabricitem resources whose CRUD runs through the shared abstraction).
//
// Usage:
//
//	go run ./tools/spncheck                  # report findings, exit 1 if any
//	go run ./tools/spncheck -dir DIR         # scan a different services directory
//	go run ./tools/spncheck -exclusions PATH # use a specific exclusions file
//
// An item whose called SDK API carries a stale or overly conservative identity
// annotation can be suppressed by listing its service package in exclusions.yaml
// (see that file).
package main

import (
	"flag"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/microsoft/terraform-provider-fabric/tools/internal/toolutil"
)

// spnSupport is the parsed service-principal support state of a single SDK API,
// derived from the "Service principal" row of its MICROSOFT ENTRA SUPPORTED
// IDENTITIES table.
type spnSupport int

const (
	spnUnknown     spnSupport = iota // no identity table (or unparsable)
	spnYes                           // "| Yes |"
	spnNo                            // "| No |"
	spnConditional                   // supported only under stated conditions
)

// sdkPackageOverrides maps a FabricItemType item name to its dedicated
// fabric-sdk-go package directory when the default lowercase derivation does not
// match. Add entries here for any item whose SDK package cannot be resolved
// automatically.
var sdkPackageOverrides = map[string]string{ //nolint:gochecknoglobals
	"Map": "maps",
}

// exit codes.
const (
	exitOK       = 0
	exitMismatch = 1
	exitError    = 2
)

// CLI flags.
var (
	dirFlag        = flag.String("dir", toolutil.DefaultServicesDir, "services directory to scan (relative to the module root)") //nolint:gochecknoglobals
	exclusionsFlag = flag.String("exclusions", "", "path to exclusions YAML file (default: auto-detected)")                      //nolint:gochecknoglobals
)

func main() {
	flag.Parse()

	os.Exit(run())
}

func run() int {
	root, err := toolutil.ModuleRoot()
	if err != nil {
		toolutil.Errf("error: %v\n", err)

		return exitError
	}

	pkgs, err := toolutil.LoadServicePackages(root, *dirFlag)
	if err != nil {
		toolutil.Errf("error loading packages: %v\n", err)

		return exitError
	}

	exclusionsPath := *exclusionsFlag
	if exclusionsPath == "" {
		exclusionsPath = filepath.Join(root, "tools", "spncheck", "exclusions.yaml")
	}

	exclusions, err := toolutil.LoadExclusions(exclusionsPath)
	if err != nil {
		toolutil.Errf("error loading exclusions: %v\n", err)

		return exitError
	}

	dc := toolutil.NewDocCache()

	results := make([]result, 0, len(pkgs))

	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			for _, e := range pkg.Errors {
				toolutil.Errf("package %s: %v\n", pkg.PkgPath, e)
			}

			return exitError
		}

		res, ok := analyzePackage(pkg, dc)
		if !ok {
			continue
		}

		results = append(results, res)
	}

	slices.SortFunc(results, func(a, b result) int { return strings.Compare(a.service, b.service) })

	return report(results, exclusions)
}

// result captures the analysis outcome for a single service package.
type result struct {
	service        string   // package name (e.g. "lakehouse")
	declared       bool     // current IsSPNSupported value
	hasDeclaration bool     // ItemTypeInfo.IsSPNSupported was found
	sdkCalls       int      // number of distinct SDK funcs called
	nonSPNAPIs     []string // SDK funcs documented as NOT supporting SPN (pkg.Func)
}

// expected returns the IsSPNSupported value implied by the SDK usage: an item
// supports service principal only when none of its called APIs is a hard "No".
func (r result) expected() bool { return len(r.nonSPNAPIs) == 0 }

// determinable is true when at least one SDK API call was found, so the
// expected status can be trusted.
func (r result) determinable() bool { return r.sdkCalls > 0 }

func report(results []result, exclusions map[string]string) int {
	b := categorize(results, exclusions)

	for _, r := range b.undeclared {
		toolutil.Outf("%-40s no ItemTypeInfo.IsSPNSupported field found\n", r.service)
	}

	printOverMarked(b.overMarked)

	printReview(b.review)

	printExcluded(b.excluded)

	printStaleExclusions(b.staleExcl)

	toolutil.Outf("\nScanned %d services: %d over-marked, %d to review, %d undeclared, %d excluded, %d stale, %d undetermined\n",
		len(results), len(b.overMarked), len(b.review), len(b.undeclared), len(b.excluded), len(b.staleExcl), b.undetermined)

	if len(b.overMarked) > 0 || len(b.review) > 0 || len(b.undeclared) > 0 || len(b.staleExcl) > 0 {
		return exitMismatch
	}

	return exitOK
}

// category is the service-principal-support classification of a single service.
type category int

const (
	catConsistent   category = iota // declared status matches SDK usage
	catUndetermined                 // no SDK calls, status not determinable
	catOverMarked                   // declared SPN-supported but a called API is not
	catReview                       // declared unsupported but every called API supports SPN
	catUndeclared                   // no ItemTypeInfo.IsSPNSupported field found
)

// classify returns the service-principal-support category for a single result.
func classify(r result) category {
	if !r.hasDeclaration {
		return catUndeclared
	}

	if !r.determinable() {
		return catUndetermined
	}

	switch {
	case r.declared && !r.expected():
		return catOverMarked
	case !r.declared && r.expected():
		return catReview
	default:
		return catConsistent
	}
}

// excludedResult is a failing service suppressed by the exclusions file.
type excludedResult struct {
	res    result
	reason string
}

// buckets groups categorized services by confidence level.
type buckets struct {
	overMarked   []result         // declared SPN-supported but a called API is not
	review       []result         // declared unsupported but every called API supports SPN
	undeclared   []result         // no ItemTypeInfo.IsSPNSupported field found
	excluded     []excludedResult // failing but suppressed via exclusions.yaml
	staleExcl    []string         // excluded services that no longer mismatch
	undetermined int              // no SDK calls, status not determinable
}

// categorize splits services into over-marked, review, undeclared, excluded,
// and undetermined groups, routing failing services listed in exclusions into
// the excluded bucket and flagging exclusions that no longer apply as stale.
func categorize(results []result, exclusions map[string]string) buckets {
	var b buckets

	used := make(map[string]struct{})

	for _, r := range results {
		switch cat := classify(r); cat {
		case catUndetermined:
			b.undetermined++
		case catConsistent:
			// declared status is correct; nothing to report.
		default:
			b.addFailing(r, cat, exclusions, used)
		}
	}

	for svc := range exclusions {
		if _, ok := used[svc]; !ok {
			b.staleExcl = append(b.staleExcl, svc)
		}
	}

	slices.Sort(b.staleExcl)

	return b
}

// addFailing routes a failing result into the excluded bucket (if its service
// is listed in exclusions) or its failing-category bucket.
func (b *buckets) addFailing(r result, cat category, exclusions map[string]string, used map[string]struct{}) {
	if reason, ok := exclusions[r.service]; ok {
		b.excluded = append(b.excluded, excludedResult{res: r, reason: reason})
		used[r.service] = struct{}{}

		return
	}

	switch cat {
	case catOverMarked:
		b.overMarked = append(b.overMarked, r)
	case catReview:
		b.review = append(b.review, r)
	case catUndeclared:
		b.undeclared = append(b.undeclared, r)
	default:
	}
}

func printOverMarked(overMarked []result) {
	if len(overMarked) == 0 {
		return
	}

	toolutil.Outf("\nOVER-MARKED — declared SPN-supported but the SDK marks a called API as NOT supporting service principal (should be IsSPNSupported = false):\n")

	for _, r := range overMarked {
		toolutil.Outf("  ✗ %-38s\n", r.service)
		printNonSPNAPIs(r)
	}
}

func printReview(review []result) {
	if len(review) == 0 {
		return
	}

	toolutil.Outf("\nREVIEW — declared NOT SPN-supported but every called API supports service principal (possibly IsSPNSupported = true, confirm manually):\n")

	for _, r := range review {
		toolutil.Outf("  ? %-38s\n", r.service)
	}
}

func printExcluded(excluded []excludedResult) {
	if len(excluded) == 0 {
		return
	}

	toolutil.Outf("\nEXCLUDED — suppressed via exclusions.yaml (confirmed intentional):\n")

	for _, e := range excluded {
		toolutil.Outf("  - %-38s %s\n", e.res.service, e.reason)
	}
}

func printStaleExclusions(stale []string) {
	if len(stale) == 0 {
		return
	}

	toolutil.Outf("\nSTALE EXCLUSIONS — no longer mismatched, remove from exclusions.yaml:\n")

	for _, svc := range stale {
		toolutil.Outf("  %s\n", svc)
	}
}

func printNonSPNAPIs(r result) {
	if len(r.nonSPNAPIs) == 0 {
		return
	}

	for _, api := range r.nonSPNAPIs {
		toolutil.Outf("    non-SPN API: %s\n", api)
	}
}

// analyzePackage inspects a service package: it locates ItemTypeInfo.IsSPNSupported
// and collects every fabric-sdk-go function it calls, flagging those that do not
// support service principal.
func analyzePackage(pkg *packages.Package, dc *toolutil.DocCache) (result, bool) {
	res := result{service: toolutil.PkgName(pkg)}

	if v, ok := toolutil.ExtractItemTypeInfoBool(pkg, "IsSPNSupported"); ok {
		res.hasDeclaration = true
		res.declared = v
	}

	nonSPNSet := map[string]struct{}{}

	for _, call := range toolutil.CollectSDKCalls(pkg) {
		res.sdkCalls++

		if parseSPNSupport(dc.DocCommentAt(call.Pos)) == spnNo {
			nonSPNSet[call.ID] = struct{}{}
		}
	}

	// Fallback: generic fabricitem resources make no direct SDK calls in-package.
	// Resolve their FabricItemType to the dedicated SDK package and scan its
	// exported client methods for service-principal support.
	if res.sdkCalls == 0 {
		determineViaItemType(pkg, dc, &res, nonSPNSet)
	}

	res.nonSPNAPIs = toolutil.SortedKeys(nonSPNSet)

	// Skip packages that have neither a declaration nor any relevance.
	if !res.hasDeclaration && res.sdkCalls == 0 {
		return res, false
	}

	return res, true
}

// determineViaItemType resolves the FabricItemType constant to its dedicated
// fabric-sdk-go package and scans that package's exported client methods for
// service-principal support, closing the gap for generic fabricitem resources.
func determineViaItemType(pkg *packages.Package, dc *toolutil.DocCache, res *result, nonSPNSet map[string]struct{}) {
	pkgDir, item, ok := toolutil.SDKPackageDir(pkg, sdkPackageOverrides)
	if !ok {
		return
	}

	methods := dc.ScanItemCRUD(pkgDir, item)
	if len(methods) == 0 {
		return
	}

	res.sdkCalls += len(methods)

	for _, m := range methods {
		if parseSPNSupport(m.Doc) == spnNo {
			nonSPNSet[strings.ToLower(item)+"."+m.Name] = struct{}{}
		}
	}
}

// parseSPNSupport reads the "Service principal" row of the MICROSOFT ENTRA
// SUPPORTED IDENTITIES table in doc and returns its support state. It normalizes
// the doc by stripping all whitespace (the table wraps across comment lines) and
// inspecting the support cell that immediately follows the "service principal"
// identity label.
func parseSPNSupport(doc string) spnSupport {
	// Collapse all whitespace so the wrapped, multi-line table becomes a single
	// contiguous string that is trivial to slice on the "|" cell delimiters.
	norm := strings.Join(strings.Fields(strings.ToLower(doc)), "")

	if !strings.Contains(norm, "microsoftentrasupportedidentities") {
		return spnUnknown
	}

	// Use LastIndex: "service principal" can appear in free-text (e.g.
	// PERMISSIONS lines) before the table. The last occurrence is always the
	// one inside the actual identity table row.
	idx := strings.LastIndex(norm, "serviceprincipal")
	if idx < 0 {
		return spnUnknown
	}

	// The support cell is the first "| value |" after the (bracketed) identity
	// label. Identity labels and URL references never contain "|".
	rest := norm[idx:]

	_, after, ok := strings.Cut(rest, "|")
	if !ok {
		return spnUnknown
	}

	cell := after

	before0, _, ok0 := strings.Cut(cell, "|")
	if !ok0 {
		return spnUnknown
	}

	switch before0 {
	case "yes":
		return spnYes
	case "no":
		return spnNo
	default:
		return spnConditional
	}
}
