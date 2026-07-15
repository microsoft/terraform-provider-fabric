// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

// Command previewcheck audits every Terraform resource/data source service
// package and verifies that its declared preview status (the IsPreview field in
// base.go's ItemTypeInfo) matches the actual preview status of the Microsoft
// Fabric Go SDK APIs that the package calls.
//
// An item must be marked as preview (IsPreview = true) when ANY of the
// fabric-sdk-go client functions it invokes is documented as preview. The SDK
// flags a preview API with one of the following phrases in the doc comment that
// sits directly above the relevant client function:
//
//	"is currently in preview"
//	"is part of a Preview release"
//
// If none of the called SDK functions are preview, the item is *likely* GA --
// but note that the SDK's preview annotations are incomplete: a missing marker
// does not guarantee the API is GA. Findings are therefore split by confidence:
//
//   - UNDER-MARKED (high confidence): declared GA but the SDK flags a called API
//     as preview. The item should be preview.
//   - REVIEW (low confidence): declared preview but no SDK preview marker was
//     found on any called API. Possibly demotable to GA, but requires manual
//     confirmation because SDK annotations are sparse.
//   - UNDETERMINED: no fabric-sdk-go calls found in the package (e.g. generic
//     fabricitem resources whose CRUD runs through the shared abstraction).
//
// Usage:
//
//	go run ./tools/previewcheck                  # report findings, exit 1 if any
//	go run ./tools/previewcheck -dir DIR         # scan a different services directory
//	go run ./tools/previewcheck -exclusions PATH # use a specific exclusions file
//
// A GA item whose called SDK API still carries a stale preview marker can be
// suppressed by listing its service package in exclusions.yaml (see that file).
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

// previewMarkers are the doc-comment phrases that flag an SDK API as preview
// (stored lowercase for case-insensitive matching).
var previewMarkers = []string{ //nolint:gochecknoglobals
	"is currently in preview",
	"is part of a preview release",
}

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
		exclusionsPath = filepath.Join(root, "tools", "previewcheck", "exclusions.yaml")
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
	declared       bool     // current IsPreview value
	hasDeclaration bool     // ItemTypeInfo.IsPreview was found
	sdkCalls       int      // number of distinct SDK funcs called
	previewAPIs    []string // SDK funcs flagged as preview (pkg.Func)
}

// expected returns the IsPreview value implied by the SDK usage.
func (r result) expected() bool { return len(r.previewAPIs) > 0 }

// determinable is true when at least one SDK API call was found, so the
// expected status can be trusted.
func (r result) determinable() bool { return r.sdkCalls > 0 }

func report(results []result, exclusions map[string]string) int {
	b := categorize(results, exclusions)

	for _, r := range b.undeclared {
		toolutil.Outf("%-40s no ItemTypeInfo.IsPreview field found\n", r.service)
	}

	printUnderMarked(b.underMarked)

	printReview(b.review)

	printExcluded(b.excluded)

	printStaleExclusions(b.staleExcl)

	toolutil.Outf("\nScanned %d services: %d under-marked, %d to review, %d undeclared, %d excluded, %d stale, %d undetermined\n",
		len(results), len(b.underMarked), len(b.review), len(b.undeclared), len(b.excluded), len(b.staleExcl), b.undetermined)

	if len(b.underMarked) > 0 || len(b.review) > 0 || len(b.undeclared) > 0 || len(b.staleExcl) > 0 {
		return exitMismatch
	}

	return exitOK
}

// category is the preview-status classification of a single service.
type category int

const (
	catConsistent   category = iota // declared status matches SDK usage
	catUndetermined                 // no SDK calls, status not determinable
	catUnderMarked                  // declared GA but a called API is preview
	catReview                       // declared preview but no SDK marker found
	catUndeclared                   // no ItemTypeInfo.IsPreview field found
)

// classify returns the preview-status category for a single result.
func classify(r result) category {
	if !r.hasDeclaration {
		return catUndeclared
	}

	if !r.determinable() {
		return catUndetermined
	}

	switch {
	case !r.declared && r.expected():
		return catUnderMarked
	case r.declared && !r.expected():
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
	underMarked  []result         // declared GA but a called API is preview
	review       []result         // declared preview but no SDK marker found
	undeclared   []result         // no ItemTypeInfo.IsPreview field found
	excluded     []excludedResult // failing but suppressed via exclusions.yaml
	staleExcl    []string         // excluded services that no longer mismatch
	undetermined int              // no SDK calls, status not determinable
}

// categorize splits services into under-marked, review, undeclared, excluded,
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
	case catUnderMarked:
		b.underMarked = append(b.underMarked, r)
	case catReview:
		b.review = append(b.review, r)
	case catUndeclared:
		b.undeclared = append(b.undeclared, r)
	default:
	}
}

func printUnderMarked(underMarked []result) {
	if len(underMarked) == 0 {
		return
	}

	toolutil.Outf("\nUNDER-MARKED — declared GA but the SDK marks a called API as preview (should be PREVIEW):\n")

	for _, r := range underMarked {
		toolutil.Outf("  ✗ %-38s\n", r.service)
		printPreviewAPIs(r)
	}
}

func printReview(review []result) {
	if len(review) == 0 {
		return
	}

	toolutil.Outf("\nREVIEW — declared PREVIEW but no SDK preview marker found (possibly GA, confirm manually):\n")

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

func printPreviewAPIs(r result) {
	if len(r.previewAPIs) == 0 {
		return
	}

	for _, api := range r.previewAPIs {
		toolutil.Outf("    preview API: %s\n", api)
	}
}

// analyzePackage inspects a service package: it locates ItemTypeInfo.IsPreview
// and collects every fabric-sdk-go function it calls, flagging preview ones.
func analyzePackage(pkg *packages.Package, dc *toolutil.DocCache) (result, bool) {
	res := result{service: toolutil.PkgName(pkg)}

	if v, ok := toolutil.ExtractItemTypeInfoBool(pkg, "IsPreview"); ok {
		res.hasDeclaration = true
		res.declared = v
	}

	previewSet := map[string]struct{}{}

	for _, call := range toolutil.CollectSDKCalls(pkg) {
		res.sdkCalls++

		if containsMarker(dc.DocCommentAt(call.Pos)) {
			previewSet[call.ID] = struct{}{}
		}
	}

	// Fallback: generic fabricitem resources make no direct SDK calls in-package.
	// Resolve their FabricItemType to the dedicated SDK package and scan its
	// exported client methods for preview markers.
	if res.sdkCalls == 0 {
		determineViaItemType(pkg, dc, &res, previewSet)
	}

	res.previewAPIs = toolutil.SortedKeys(previewSet)

	// Skip packages that have neither a declaration nor any relevance.
	if !res.hasDeclaration && res.sdkCalls == 0 {
		return res, false
	}

	return res, true
}

// determineViaItemType resolves the FabricItemType constant to its dedicated
// fabric-sdk-go package and scans that package's exported client methods for
// preview markers, closing the gap for generic fabricitem resources.
func determineViaItemType(pkg *packages.Package, dc *toolutil.DocCache, res *result, previewSet map[string]struct{}) {
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
		if containsMarker(m.Doc) {
			previewSet[strings.ToLower(item)+"."+m.Name] = struct{}{}
		}
	}
}

func containsMarker(doc string) bool {
	lower := strings.ToLower(doc)

	for _, marker := range previewMarkers {
		if strings.Contains(lower, marker) {
			return true
		}
	}

	return false
}
