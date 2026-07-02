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
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
	"gopkg.in/yaml.v3"

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

// previewExclusion suppresses the preview check for a single service package,
// typically a GA item whose called SDK API still carries a stale preview marker.
type previewExclusion struct {
	Service string `json:"service" yaml:"service"`
	Reason  string `json:"reason"  yaml:"reason"`
}

// exclusionsFile is the top-level YAML structure of the exclusions file.
type exclusionsFile struct {
	Exclusions []previewExclusion `yaml:"exclusions"`
}

// loadExclusions reads the exclusions YAML file and returns a map of service
// package name to reason.
func loadExclusions(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading exclusions file %s: %w", path, err)
	}

	var ef exclusionsFile

	err = yaml.Unmarshal(data, &ef)
	if err != nil {
		return nil, fmt.Errorf("parsing exclusions file %s: %w", path, err)
	}

	result := make(map[string]string, len(ef.Exclusions))

	for _, e := range ef.Exclusions {
		if e.Service == "" || e.Reason == "" {
			return nil, fmt.Errorf("exclusion entry missing required field (service, reason): %+v", e)
		}

		result[e.Service] = e.Reason
	}

	return result, nil
}

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

	exclusions, err := loadExclusions(exclusionsPath)
	if err != nil {
		toolutil.Errf("error loading exclusions: %v\n", err)

		return exitError
	}

	dc := newDocCache()

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
func analyzePackage(pkg *packages.Package, dc *docCache) (result, bool) {
	res := result{service: pkgName(pkg)}

	findIsPreview(pkg, &res)

	previewSet := map[string]struct{}{}
	seen := map[string]struct{}{}

	for _, file := range pkg.Syntax {
		collectSDKCalls(pkg, file, dc, &res, previewSet, seen)
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
func determineViaItemType(pkg *packages.Package, dc *docCache, res *result, previewSet map[string]struct{}) {
	pkgDir, item, ok := sdkPackageDir(pkg)
	if !ok {
		return
	}

	crud, methods := dc.scanCRUD(pkgDir, item)
	if crud == 0 {
		return
	}

	res.sdkCalls += crud

	for _, m := range methods {
		previewSet[strings.ToLower(item)+"."+m] = struct{}{}
	}
}

// sdkPackageDir resolves a package's FabricItemType constant to its dedicated
// fabric-sdk-go package directory, returning the directory and the item name.
func sdkPackageDir(pkg *packages.Package) (string, string, bool) { //nolint:revive // dir + item name + ok
	constName, anyFile := findFabricItemType(pkg)
	if constName == "" || anyFile == "" {
		return "", "", false
	}

	// ItemType constants are named ItemType<Name>; the SDK package is fabric/<name>.
	item := strings.TrimPrefix(constName, "ItemType")
	if item == constName || item == "" {
		return "", "", false
	}

	fabricDir := fabricRoot(anyFile)
	if fabricDir == "" {
		return "", "", false
	}

	pkgName := strings.ToLower(item)
	if override, ok := sdkPackageOverrides[item]; ok {
		pkgName = override
	}

	return filepath.Join(fabricDir, pkgName), item, true
}

// findFabricItemType locates the FabricItemType const (= fabcore.ItemType<Name>)
// and returns the SDK constant name plus a file in the fabcore package (used to
// derive the SDK fabric/ directory).
func findFabricItemType(pkg *packages.Package) (string, string) { //nolint:revive // two strings: SDK const name and fabcore file path
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.CONST {
				continue
			}

			for _, spec := range gen.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok || !hasName(vs.Names, "FabricItemType") {
					continue
				}

				for _, val := range vs.Values {
					sel, ok := val.(*ast.SelectorExpr)
					if !ok {
						continue
					}

					obj := pkg.TypesInfo.ObjectOf(sel.Sel)
					if obj == nil || obj.Pkg() == nil {
						continue
					}

					return obj.Name(), pkg.Fset.Position(obj.Pos()).Filename
				}
			}
		}
	}

	return "", ""
}

// fabricRoot returns the fabric-sdk-go "fabric" directory given a file inside
// the fabric/core package (e.g. .../fabric-sdk-go@v/fabric/core/x.go -> .../fabric).
func fabricRoot(coreFile string) string {
	dir := filepath.Dir(coreFile)
	if filepath.Base(dir) != "core" {
		return ""
	}

	return filepath.Dir(dir)
}

// findIsPreview locates the IsPreview field inside the ItemTypeInfo composite
// literal and records its value.
func findIsPreview(pkg *packages.Package, res *result) {
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.VAR {
				continue
			}

			for _, spec := range gen.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok || !hasName(vs.Names, "ItemTypeInfo") {
					continue
				}

				for _, val := range vs.Values {
					lit, ok := val.(*ast.CompositeLit)
					if !ok {
						continue
					}

					if extractIsPreview(lit, res) {
						return
					}
				}
			}
		}
	}
}

func extractIsPreview(lit *ast.CompositeLit, res *result) bool {
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != "IsPreview" {
			continue
		}

		ident, ok := kv.Value.(*ast.Ident)
		if !ok || (ident.Name != "true" && ident.Name != "false") {
			continue
		}

		res.hasDeclaration = true
		res.declared = ident.Name == "true"

		return true
	}

	return false
}

// collectSDKCalls walks call expressions in a file, resolves each callee, and
// records distinct fabric-sdk-go functions (flagging preview ones). The seen set
// is shared across the package's files so each callee is counted at most once.
func collectSDKCalls(pkg *packages.Package, file *ast.File, dc *docCache, res *result, previewSet, seen map[string]struct{}) {
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Only x.Y(...) style calls; SDK funcs are always reached via a selector.
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		// Resolve the selector to its real declared object via the type checker
		// (relies on NeedTypes/NeedTypesInfo/NeedDeps when loading packages).
		obj := pkg.TypesInfo.ObjectOf(sel.Sel)

		fn, ok := obj.(*types.Func)
		if !ok || fn.Pkg() == nil {
			return true
		}

		if !strings.HasPrefix(fn.Pkg().Path(), toolutil.SDKModulePath) {
			return true
		}

		// Identify the callee as "pkgName.FuncName" and skip duplicates.
		id := fn.Pkg().Name() + "." + fn.Name()
		if _, dup := seen[id]; dup {
			return true
		}

		seen[id] = struct{}{}
		res.sdkCalls++

		// fn.Pos() points at the SDK declaration; check its doc comment for a
		// preview marker and record the func if it is preview.
		if dc.isPreview(pkg.Fset.Position(fn.Pos())) {
			previewSet[id] = struct{}{}
		}

		return true
	})
}

// docCache reads SDK source files lazily and caches the preview status of the
// doc comment directly above a declaration position.
type docCache struct {
	files   map[string][]string
	results map[string]bool
	methods map[string][]string
	crud    map[string]int
}

func newDocCache() *docCache {
	return &docCache{
		files:   map[string][]string{},
		results: map[string]bool{},
		methods: map[string][]string{},
		crud:    map[string]int{},
	}
}

// scanCRUD scans the *_client.go files in an SDK package directory and returns
// the count of the item's exported CRUD methods plus the names of those flagged
// preview. A non-zero count means the item status is determinable. Cached per
// directory and item.
func (d *docCache) scanCRUD(pkgDir, item string) (int, []string) {
	cacheKey := pkgDir + "\x00" + item
	if v, ok := d.methods[cacheKey]; ok {
		return d.crud[cacheKey], v
	}

	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		d.crud[cacheKey] = 0
		d.methods[cacheKey] = nil

		return 0, nil
	}

	var (
		crud    int
		preview []string
	)

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), "_client.go") {
			continue
		}

		lines, rErr := d.readFile(filepath.Join(pkgDir, e.Name()))
		if rErr != nil {
			continue
		}

		c, p := scanItemMethods(lines, item)
		crud += c

		preview = append(preview, p...)
	}

	slices.Sort(preview)

	d.crud[cacheKey] = crud
	d.methods[cacheKey] = preview

	return crud, preview
}

// isPreview reports whether the doc comment immediately above the declaration at
// pos contains a preview marker.
func (d *docCache) isPreview(pos token.Position) bool {
	if pos.Filename == "" || pos.Line <= 0 {
		return false
	}

	key := fmt.Sprintf("%s:%d", pos.Filename, pos.Line)
	if v, ok := d.results[key]; ok {
		return v
	}

	lines, err := d.readFile(pos.Filename)
	if err != nil {
		d.results[key] = false

		return false
	}

	doc := docCommentAbove(lines, pos.Line)
	preview := containsMarker(doc)
	d.results[key] = preview

	return preview
}

func (d *docCache) readFile(path string) ([]string, error) {
	if lines, ok := d.files[path]; ok {
		return lines, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	d.files[path] = lines

	return lines, nil
}

// docCommentAbove returns the contiguous //-style comment block that sits
// directly above the 1-indexed declLine.
func docCommentAbove(lines []string, declLine int) string {
	// declLine is 1-indexed; the line above it is at slice index declLine-2.
	var doc []string

	for i := declLine - 2; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "//") {
			doc = append(doc, trimmed)

			continue
		}

		break
	}

	return strings.Join(doc, "\n")
}

// scanItemMethods returns the count of the item's CRUD client methods in the
// file and the names of those whose preceding doc comment is marked preview.
func scanItemMethods(lines []string, item string) (int, []string) {
	var (
		crud    int
		preview []string
	)

	for i, line := range lines {
		_, name := toolutil.ParseClientMethod(line)
		if name == "" || !isItemCRUD(name, item) {
			continue
		}

		crud++

		// Line numbers passed to docCommentAbove are 1-indexed.
		if containsMarker(docCommentAbove(lines, i+1)) {
			preview = append(preview, name)
		}
	}

	return crud, preview
}

// isItemCRUD reports whether method is a CRUD operation on the item itself,
// e.g. GetNotebook, BeginCreateNotebook, ListNotebooks, UpdateNotebook.
func isItemCRUD(method, item string) bool {
	if !strings.Contains(strings.ToLower(method), strings.ToLower(item)) {
		return false
	}

	for _, verb := range []string{"Get", "List", "Create", "Update", "Delete"} {
		if strings.Contains(method, verb) {
			return true
		}
	}

	return false
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

// loadServicePackages loads all packages under the services directory with full
// type information.
func pkgName(pkg *packages.Package) string {
	if pkg.Name != "" {
		return pkg.Name
	}

	return filepath.Base(pkg.PkgPath)
}

func hasName(names []*ast.Ident, target string) bool {
	for _, n := range names {
		if n.Name == target {
			return true
		}
	}

	return false
}
