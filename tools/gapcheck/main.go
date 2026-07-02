// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

// Command gapcheck detects control-plane Fabric entities exposed by
// fabric-sdk-go that have no corresponding Terraform resource in this provider.
//
// It classifies SDK client methods into noun-grouped operation sets and
// identifies control-plane resources (Create+Delete, Add+Delete, or Get+Set
// pairs) that are not yet covered by a TF resource type.
//
// Usage:
//
//	go run ./tools/gapcheck            # human-readable report
//	go run ./tools/gapcheck -json      # machine-readable JSON output
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"golang.org/x/tools/go/packages"
	"gopkg.in/yaml.v3"

	"github.com/microsoft/terraform-provider-fabric/tools/internal/toolutil"
)

// exit codes.
const (
	exitOK   = 0
	exitGaps = 1
	exitErr  = 2
)

// CLI flags.
var (
	jsonFlag       = flag.Bool("json", false, "output machine-readable JSON")                               //nolint:gochecknoglobals
	dirFlag        = flag.String("dir", toolutil.DefaultServicesDir, "services directory to scan")          //nolint:gochecknoglobals
	exclusionsFlag = flag.String("exclusions", "", "path to exclusions YAML file (default: auto-detected)") //nolint:gochecknoglobals
)

// nounTypeOverrides maps SDK noun (PascalCase) to TF type when the automatic
// derivation does not match. Used for known naming mismatches.
var nounTypeOverrides = map[string]string{ //nolint:gochecknoglobals
	// SDK uses "CosmosDBDatabase" but TF type is "cosmos_db"
	"CosmosDBDatabase":  "cosmos_db",
	"CosmosDBDatabases": "cosmos_db",
	// SDK uses "GraphQLAPI"/"GraphQLApis" but TF type is "graphql_api"
	"GraphQLAPI":           "graphql_api",
	"GraphQLApis":          "graphql_api",
	"GraphQLAPIDefinition": "graphql_api",
	// OneLake data access: SDK nouns map to "onelake_data_access_security"
	"DataAccessRoles":      "onelake_data_access_security",
	"DataAccessRole":       "onelake_data_access_security",
	"SingleDataAccessRole": "onelake_data_access_security",
	// Reflex is the old name for Activator
	"Reflex":   "activator",
	"Reflexes": "activator",
	// Spark custom pools: SDK splits workspace/capacity but TF has one resource
	"WorkspaceCustomPool":    "spark_custom_pool",
	"WorkspaceCustomPools":   "spark_custom_pool",
	"CapacityCustomPoolBeta": "spark_custom_pool",
}

// intentionallySkipped is no longer hard-coded — see exclusions.yaml.

// exclusionEntry represents a single exclusion in the YAML file.
type exclusionEntry struct {
	Package string `json:"package" yaml:"package"`
	Client  string `json:"client"  yaml:"client"`
	Method  string `json:"method"  yaml:"method"`
	Reason  string `json:"reason"  yaml:"reason"`
}

// exclusionKey returns the unique key for matching against noun groups.
func (e exclusionEntry) exclusionKey() string {
	return e.Package + "/" + e.Client + "/" + e.Method
}

// exclusionsFile is the top-level YAML structure.
type exclusionsFile struct {
	Exclusions []exclusionEntry `yaml:"exclusions"`
}

// loadExclusions reads the exclusions YAML file and returns a map keyed by
// package/client/noun.
func loadExclusions(path string) (map[string]exclusionEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading exclusions file %s: %w", path, err)
	}

	var ef exclusionsFile
	if err := yaml.Unmarshal(data, &ef); err != nil {
		return nil, fmt.Errorf("parsing exclusions file %s: %w", path, err)
	}

	result := make(map[string]exclusionEntry, len(ef.Exclusions))

	for _, e := range ef.Exclusions {
		if e.Package == "" || e.Client == "" || e.Reason == "" {
			return nil, fmt.Errorf("exclusion entry missing required field (package, client, reason): %+v", e)
		}

		// If method is empty or "*", it excludes all methods on the client.
		if e.Method == "" {
			e.Method = "*"
		}

		result[e.exclusionKey()] = e
	}

	return result, nil
}

// lookupExclusion checks if a method is excluded, either by exact match or
// by client-level wildcard (method="*").
func lookupExclusion(exclusions map[string]exclusionEntry, m sdkMethod) (exclusionEntry, bool) {
	// Try exact match first.
	if entry, ok := exclusions[m.sdkMethodExclKey()]; ok {
		return entry, true
	}

	// Try client-level wildcard.
	wildcardKey := m.Package + "/" + m.Client + "/*"
	if entry, ok := exclusions[wildcardKey]; ok {
		return entry, true
	}

	return exclusionEntry{}, false
}

// actionVerbs are verbs that never contribute to control-plane classification.
var actionVerbs = map[string]struct{}{ //nolint:gochecknoglobals
	"Begin":     {},
	"Test":      {},
	"Run":       {},
	"Execute":   {},
	"Query":     {},
	"Submit":    {},
	"Cancel":    {},
	"Move":      {},
	"Sync":      {},
	"Revoke":    {},
	"Accept":    {},
	"Reset":     {},
	"Associate": {},
	"New":       {}, // pager constructors
}

// sortedActionVerbs is actionVerbs as a sorted slice for deterministic iteration.
var sortedActionVerbs = func() []string { //nolint:gochecknoglobals
	keys := make([]string, 0, len(actionVerbs))
	for k := range actionVerbs {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	return keys
}()

func main() {
	flag.Parse()
	os.Exit(run())
}

func run() int {
	root, err := toolutil.ModuleRoot()
	if err != nil {
		toolutil.Errf("error: %v\n", err)

		return exitErr
	}

	sdkDir, err := sdkFabricDir(root)
	if err != nil {
		toolutil.Errf("error resolving SDK: %v\n", err)

		return exitErr
	}

	// Resolve exclusions file path.
	exclusionsPath := *exclusionsFlag
	if exclusionsPath == "" {
		// Default: tools/gapcheck/exclusions.yaml relative to module root.
		exclusionsPath = filepath.Join(root, "tools", "gapcheck", "exclusions.yaml")
	}

	exclusions, err := loadExclusions(exclusionsPath)
	if err != nil {
		toolutil.Errf("error loading exclusions: %v\n", err)

		return exitErr
	}

	// 1. Extract noun groups from SDK clients.
	nounGroups, err := extractNounGroups(sdkDir)
	if err != nil {
		toolutil.Errf("error extracting SDK methods: %v\n", err)

		return exitErr
	}

	// 2. Build TF resource inventory.
	inventory, err := buildInventory(root, *dirFlag)
	if err != nil {
		toolutil.Errf("error building inventory: %v\n", err)

		return exitErr
	}

	// 3. Detect which SDK methods are actually called in service implementations.
	calledMethods, err := extractCoveredMethods(root, *dirFlag)
	if err != nil {
		toolutil.Errf("error detecting covered methods: %v\n", err)

		return exitErr
	}

	// 4. Extract ALL SDK methods for completeness checking.
	allMethods, err := extractAllMethods(sdkDir)
	if err != nil {
		toolutil.Errf("error extracting SDK methods: %v\n", err)

		return exitErr
	}

	// 4. Classify and diff.
	gaps := findGaps(nounGroups, allMethods, inventory, exclusions, calledMethods)

	// 5. Report.
	return report(gaps)
}

// nounGroup represents a classified SDK operation group.
type nounGroup struct {
	Package string   `json:"package"`
	Client  string   `json:"client"`
	Noun    string   `json:"noun"`
	Verbs   []string `json:"verbs"`
	Kind    string   `json:"kind"` // "resource", "singleton"
}

// gap represents an SDK method that is not covered by a TF resource.
type gap struct {
	SDKMethod sdkMethod `json:"sdk_method"`
	Verb      string    `json:"verb,omitempty"`
	Skipped   string    `json:"skipped,omitempty"`
	Stale     bool      `json:"stale,omitempty"`
	Covered   bool      `json:"covered,omitempty"`
}

// sdkMethod represents a single SDK client method.
type sdkMethod struct {
	Package string `json:"package"`
	Client  string `json:"client"`
	Method  string `json:"method"`
}

// sdkMethodExclKey returns the exclusion key for an sdkMethod.
func (m sdkMethod) sdkMethodExclKey() string {
	return m.Package + "/" + m.Client + "/" + m.Method
}

// extractCoveredMethods type-checks the service packages under servicesDir and
// returns the set of package/client/method keys for every fabric-sdk-go client
// method they invoke.
//
// It uses go/packages type resolution (not regex) so that calls made through
// local variables, function parameters, cross-file client fields, import
// aliases, or embedded types are all resolved to the exact SDK method they
// reference. Keys are formatted "sdkPackage/ClientType/Method" to match
// sdkMethod.sdkMethodExclKey.
func extractCoveredMethods(root, servicesDir string) (map[string]struct{}, error) {
	pkgs, err := toolutil.LoadServicePackages(root, servicesDir)
	if err != nil {
		return nil, err
	}

	covered := map[string]struct{}{}

	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			for _, e := range pkg.Errors {
				toolutil.Errf("package %s: %v\n", pkg.PkgPath, e)
			}

			return nil, fmt.Errorf("type-checking failed for package %s", pkg.PkgPath)
		}

		for _, file := range pkg.Syntax {
			collectCoveredCalls(pkg, file, covered)
		}
	}

	return covered, nil
}

// collectCoveredCalls walks call expressions in file, resolves each callee via
// the type checker, and records every fabric-sdk-go client method it finds in
// covered as "sdkPackage/ClientType/Method".
func collectCoveredCalls(pkg *packages.Package, file *ast.File, covered map[string]struct{}) {
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// SDK methods are always reached through a selector: x.Method(...).
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		fn, ok := pkg.TypesInfo.ObjectOf(sel.Sel).(*types.Func)
		if !ok || fn.Pkg() == nil {
			return true
		}

		if !strings.HasPrefix(fn.Pkg().Path(), toolutil.SDKModulePath) {
			return true
		}

		sig, ok := fn.Type().(*types.Signature)
		if !ok || sig.Recv() == nil {
			return true // package-level SDK func, not a client method
		}

		clientType := receiverTypeName(sig.Recv().Type())
		if !strings.HasSuffix(clientType, "Client") {
			return true
		}

		// The SDK directory name (used by extractAllMethods) is the last element
		// of the import path, which is also the package name in fabric-sdk-go.
		sdkPkg := fn.Pkg().Path()
		if i := strings.LastIndex(sdkPkg, "/"); i >= 0 {
			sdkPkg = sdkPkg[i+1:]
		}

		covered[sdkPkg+"/"+clientType+"/"+fn.Name()] = struct{}{}

		return true
	})
}

// receiverTypeName returns the bare type name of a method receiver, unwrapping
// a pointer receiver (e.g. *fabcore.WorkspacesClient -> "WorkspacesClient").
func receiverTypeName(t types.Type) string {
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	if named, ok := t.(*types.Named); ok {
		return named.Obj().Name()
	}

	return ""
}

// extractAllMethods scans SDK client files and returns every public method
// (excluding internal helpers, pager constructors, and Begin* LRO twins).
func extractAllMethods(sdkDir string) ([]sdkMethod, error) {
	var results []sdkMethod

	entries, err := os.ReadDir(sdkDir)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if !e.IsDir() || e.Name() == "fake" {
			continue
		}

		pkgDir := filepath.Join(sdkDir, e.Name())

		pkgEntries, readErr := os.ReadDir(pkgDir)
		if readErr != nil {
			return nil, fmt.Errorf("reading %s: %w", e.Name(), readErr)
		}

		for _, f := range pkgEntries {
			if f.IsDir() || !strings.HasSuffix(f.Name(), "_client.go") || strings.Contains(f.Name(), "example") {
				continue
			}

			data, fileErr := os.ReadFile(filepath.Join(pkgDir, f.Name()))
			if fileErr != nil {
				return nil, fmt.Errorf("reading %s/%s: %w", e.Name(), f.Name(), fileErr)
			}

			for line := range strings.SplitSeq(string(data), "\n") {
				client, method := parseClientMethod(line)
				if client == "" {
					continue
				}

				// Skip Begin* LRO twins (the sync version is tracked).
				if strings.HasPrefix(method, "Begin") {
					continue
				}

				// Skip pager constructors.
				if strings.HasPrefix(method, "New") && strings.Contains(method, "Pager") {
					continue
				}

				// Skip core.ItemsClient (generic fan-out).
				if e.Name() == "core" && client == "ItemsClient" {
					continue
				}

				results = append(results, sdkMethod{
					Package: e.Name(),
					Client:  client,
					Method:  method,
				})
			}
		}
	}

	return results, nil
}

// extractNounGroups scans all *_client.go files in the SDK fabric/ directory
// (core + per-item packages) and returns classified noun groups that look like
// control-plane resources.
func extractNounGroups(sdkDir string) ([]nounGroup, error) {
	var results []nounGroup

	// Scan core package.
	coreGroups, err := scanPackageClients(filepath.Join(sdkDir, "core"), "core")
	if err != nil {
		return nil, fmt.Errorf("scanning core: %w", err)
	}

	results = append(results, coreGroups...)

	// Scan per-item packages (non-ItemsClient clients that have control-plane signatures).
	entries, err := os.ReadDir(sdkDir)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if !e.IsDir() || e.Name() == "core" || e.Name() == "fake" {
			continue
		}

		pkgDir := filepath.Join(sdkDir, e.Name())

		groups, scanErr := scanPackageClients(pkgDir, e.Name())
		if scanErr != nil {
			return nil, fmt.Errorf("scanning package %s: %w", e.Name(), scanErr)
		}

		results = append(results, groups...)
	}

	return results, nil
}

// scanPackageClients reads all *_client.go files in a directory, extracts method
// names per client, and classifies noun groups.
func scanPackageClients(dir, pkgName string) ([]nounGroup, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Collect methods per client.
	clientMethods := map[string][]string{}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), "_client.go") || strings.Contains(e.Name(), "example") {
			continue
		}

		data, readErr := os.ReadFile(filepath.Join(dir, e.Name()))
		if readErr != nil {
			return nil, fmt.Errorf("reading %s: %w", e.Name(), readErr)
		}

		for line := range strings.SplitSeq(string(data), "\n") {
			client, method := parseClientMethod(line)
			if client == "" {
				continue
			}

			clientMethods[client] = append(clientMethods[client], method)
		}
	}

	var results []nounGroup

	for client, methods := range clientMethods {
		// Skip the generic core.ItemsClient (per-item packages cover individual items).
		if pkgName == "core" && client == "ItemsClient" {
			continue
		}

		groups := classifyMethods(pkgName, client, methods)
		results = append(results, groups...)
	}

	return results, nil
}

// parseClientMethod extracts the client type and method name from an SDK client
// method declaration line, restricted to *…Client receivers and excluding the
// generated request/response helpers.
func parseClientMethod(line string) (string, string) {
	recvType, method := toolutil.ParseClientMethod(line)
	if recvType == "" || !strings.HasSuffix(recvType, "Client") {
		return "", ""
	}

	// Skip internal helpers.
	if strings.HasSuffix(method, "CreateRequest") || strings.HasSuffix(method, "HandleResponse") {
		return "", ""
	}

	return recvType, method
}

// classifyMethods groups methods by noun, classifies each noun group, and
// returns those that qualify as control-plane resources.
func classifyMethods(pkg, client string, methods []string) []nounGroup {
	// Build noun → verb set.
	nounVerbs := map[string]map[string]struct{}{}

	for _, method := range methods {
		verb, noun := splitVerbNoun(method)
		if verb == "" || noun == "" {
			continue
		}

		// For Begin* LRO wrappers, re-parse the inner method to extract
		// the real verb+noun. This handles future LRO-only resources that
		// may not have a synchronous twin.
		if verb == "Begin" {
			verb, noun = splitVerbNoun(noun)
			if verb == "" || noun == "" {
				continue
			}
		}

		// Skip action verbs entirely.
		if isActionVerb(verb) {
			continue
		}

		if _, ok := nounVerbs[noun]; !ok {
			nounVerbs[noun] = map[string]struct{}{}
		}

		nounVerbs[noun][verb] = struct{}{}
	}

	var results []nounGroup

	for noun, verbs := range nounVerbs {
		kind := classifyNounKind(verbs)
		if kind == "" {
			continue
		}

		verbList := toolutil.SortedKeys(verbs)

		results = append(results, nounGroup{
			Package: pkg,
			Client:  client,
			Noun:    noun,
			Verbs:   verbList,
			Kind:    kind,
		})
	}

	return results
}

// classifyNounKind determines if a noun group qualifies as a control-plane
// resource based on its verb signature.
func classifyNounKind(verbs map[string]struct{}) string {
	_, hasCreate := verbs["Create"]
	_, hasCreateOrUpdate := verbs["CreateOrUpdate"]
	_, hasDelete := verbs["Delete"]
	_, hasAdd := verbs["Add"]
	_, hasGet := verbs["Get"]
	_, hasSet := verbs["Set"]

	switch {
	case (hasCreate || hasCreateOrUpdate) && hasDelete:
		return "resource"
	case hasAdd && hasDelete:
		return "resource"
	case hasGet && hasSet:
		return "singleton"
	default:
		return ""
	}
}

// splitVerbNoun splits a PascalCase method name into its leading verb and the
// remaining noun. Returns empty strings for unrecognized shapes.
func splitVerbNoun(method string) (string, string) {
	// Try multi-word verbs first.
	if after, ok := strings.CutPrefix(method, "CreateOrUpdate"); ok {
		noun := after
		if noun != "" {
			return "CreateOrUpdate", noun
		}
	}

	// Try known control verbs.
	for _, v := range []string{"Create", "Delete", "Get", "Set", "Add", "Update", "List"} {
		if after, ok := strings.CutPrefix(method, v); ok {
			noun := after
			if noun != "" {
				return v, noun
			}
		}
	}

	// Try action verbs (iterated in sorted order for determinism).
	for _, v := range sortedActionVerbs {
		if after, ok := strings.CutPrefix(method, v); ok {
			return v, after
		}
	}

	return "", ""
}

func isActionVerb(verb string) bool {
	_, ok := actionVerbs[verb]

	return ok
}

// findGaps determines which classified noun groups are not covered by the TF inventory.
// findGaps determines which SDK methods are not accounted for — either covered
// by a TF resource (via noun-group matching) or excluded in exclusions.yaml.
func findGaps(nounGroups []nounGroup, allMethods []sdkMethod, inventory map[string]struct{}, exclusions map[string]exclusionEntry, calledMethods map[string]struct{}) []gap {
	// Build coverage sets.
	// coveredClients: for per-item ItemsClients, if any noun matches TF inventory
	// the entire client is covered (one resource = one ItemsClient).
	// coveredNouns: for multi-resource clients (e.g. core.WorkspacesClient),
	// track individual noun coverage.
	coveredClients := map[string]struct{}{} // key: pkg/client
	coveredNouns := map[string]struct{}{}   // key: pkg/client/noun

	for _, ng := range nounGroups {
		candidates := buildCandidates(ng)

		for _, c := range candidates {
			if _, ok := inventory[c]; ok {
				if ng.Client == "ItemsClient" {
					// Per-item package: entire client is one resource.
					coveredClients[ng.Package+"/"+ng.Client] = struct{}{}
				}

				coveredNouns[ng.Package+"/"+ng.Client+"/"+ng.Noun] = struct{}{}

				break
			}
		}
	}

	// Also mark per-item packages as covered if their package name matches
	// any TF inventory type (handles List-only data sources like dashboard).
	for _, m := range allMethods {
		if m.Client != "ItemsClient" {
			continue
		}

		clientKey := m.Package + "/" + m.Client
		if _, ok := coveredClients[clientKey]; ok {
			continue // already covered
		}

		// Check if the package name (or its snake_case) matches an inventory type.
		pkgSnake := pascalToSnake(m.Package)
		if _, ok := inventory[pkgSnake]; ok {
			coveredClients[clientKey] = struct{}{}

			continue
		}

		// Fall back to the SDK method noun. Compound package directories are
		// lowercase concatenations (e.g. "paginatedreport"), so pascalToSnake
		// leaves them unchanged and misses inventory types like "paginated_report".
		// The method noun (e.g. "PaginatedReports") snake-converts correctly.
		for _, c := range itemNounCandidates(m.Method) {
			if _, ok := inventory[c]; ok {
				coveredClients[clientKey] = struct{}{}

				break
			}
		}
	}

	// For each SDK method, check if it is covered or excluded.
	//
	// Gap reporting is intentionally exhaustive and method-level: every SDK
	// method that is not covered, not an action verb, and not excluded is
	// reported — including read-only (Get/List) and update-only methods. The
	// verb-signature classifier (classifyNounKind) governs coverage matching
	// only; it deliberately does NOT filter this report. The goal is that every
	// SDK operation is explicitly accounted for (implemented, excluded with a
	// reason, or a pending gap), so the noise is pruned via exclusions.yaml
	// rather than by narrowing what is enumerated here. See README.md.
	var gaps []gap

	for _, m := range allMethods {
		// Determine which noun group this method belongs to (if any).
		verb, noun := splitVerbNoun(m.Method)

		// Check if the entire client is covered (per-item packages).
		clientKey := m.Package + "/" + m.Client
		if _, ok := coveredClients[clientKey]; ok {
			// Method is covered. Check for stale exclusion.
			if entry, ok := lookupExclusion(exclusions, m); ok {
				gaps = append(gaps, gap{
					SDKMethod: m,
					Verb:      verb,
					Skipped:   entry.Reason,
					Stale:     true,
				})
			}

			continue
		}

		// Check if this specific method is called in a service implementation.
		methodKey := m.sdkMethodExclKey()
		if _, ok := calledMethods[methodKey]; ok {
			// Method is called in service code. Check for stale exclusion.
			if entry, ok := lookupExclusion(exclusions, m); ok {
				gaps = append(gaps, gap{
					SDKMethod: m,
					Verb:      verb,
					Skipped:   entry.Reason,
					Stale:     true,
				})
			}

			continue
		}

		// Check noun-level coverage (for multi-resource clients like WorkspacesClient).
		if noun != "" {
			// Try the noun as-is and its singular/plural variants.
			nounVariants := []string{noun, strings.TrimSuffix(noun, "s")}
			if !strings.HasSuffix(noun, "s") {
				nounVariants = append(nounVariants, noun+"s")
			}

			covered := false

			for _, nv := range nounVariants {
				nounKey := m.Package + "/" + m.Client + "/" + nv
				if _, ok := coveredNouns[nounKey]; ok {
					covered = true

					break
				}
			}

			if covered {
				// Method is covered. Check for stale exclusion.
				if entry, ok := lookupExclusion(exclusions, m); ok {
					gaps = append(gaps, gap{
						SDKMethod: m,
						Verb:      verb,
						Skipped:   entry.Reason,
						Stale:     true,
					})
				}

				continue
			}
		}

		// Action-style methods (Run, Cancel, Sync, etc.) are never
		// control-plane resources, so they can never be gaps. classifyMethods
		// already excludes them from noun-group classification; mirror that here
		// so an uncovered action method is not misreported as a gap.
		if isActionVerb(verb) {
			continue
		}

		// Check if excluded — by exact method or by client wildcard.
		if entry, ok := lookupExclusion(exclusions, m); ok {
			gaps = append(gaps, gap{
				SDKMethod: m,
				Verb:      verb,
				Skipped:   entry.Reason,
			})

			continue
		}

		// Not covered and not excluded — this is a real gap.
		gaps = append(gaps, gap{
			SDKMethod: m,
			Verb:      verb,
		})
	}

	slices.SortFunc(gaps, func(a, b gap) int {
		if c := strings.Compare(a.SDKMethod.Package, b.SDKMethod.Package); c != 0 {
			return c
		}

		if c := strings.Compare(a.SDKMethod.Client, b.SDKMethod.Client); c != 0 {
			return c
		}

		return strings.Compare(a.SDKMethod.Method, b.SDKMethod.Method)
	})

	return gaps
}

// itemNounCandidates derives candidate TF type strings from a per-item
// ItemsClient method by extracting its noun. This lets compound package
// directories (e.g. "paginatedreport", "mirroredwarehouse") match inventory
// types ("paginated_report", "mirrored_warehouse") that the lowercase package
// name alone cannot produce via pascalToSnake.
func itemNounCandidates(method string) []string {
	_, noun := splitVerbNoun(method)
	if noun == "" {
		return nil
	}

	if override, ok := nounTypeOverrides[noun]; ok {
		return []string{override}
	}

	bareNoun := pascalToSnake(noun)
	bareSingular := strings.TrimSuffix(bareNoun, "s")

	if bareNoun == bareSingular {
		return []string{bareSingular}
	}

	return []string{bareSingular, bareNoun}
}

// buildCandidates generates the multi-candidate type strings for matching:
// 1. explicit override (if present)
// 2. bare noun → snake_case (singular)
// 3. bare noun → snake_case (plural, for "rules" style resources)
// 4. client-entity prefix + noun (singular)
// 5. client-entity prefix + noun (plural)
func buildCandidates(ng nounGroup) []string {
	// Check for explicit override first.
	if override, ok := nounTypeOverrides[ng.Noun]; ok {
		return []string{override}
	}

	bareNoun := pascalToSnake(ng.Noun)

	// Normalize plurals for list-derived nouns.
	bareSingular := strings.TrimSuffix(bareNoun, "s")

	// Derive prefix from client name: "WorkspacesClient" → "workspace"
	clientEntity := strings.TrimSuffix(ng.Client, "Client")
	clientEntity = strings.TrimSuffix(clientEntity, "s") // WorkspacesClient → Workspace
	prefix := pascalToSnake(clientEntity)

	prefixedSingular := prefix + "_" + bareSingular
	prefixedPlural := prefix + "_" + bareNoun

	seen := map[string]struct{}{}
	var candidates []string

	add := func(c string) {
		if _, ok := seen[c]; !ok {
			seen[c] = struct{}{}
			candidates = append(candidates, c)
		}
	}

	add(bareSingular)

	if bareNoun != bareSingular {
		add(bareNoun) // plural form
	}

	add(prefixedSingular)

	if prefixedPlural != prefixedSingular {
		add(prefixedPlural)
	}

	return candidates
}

// buildInventory collects all TF resource Type values from base.go files.
// Every service package is expected to have a base.go declaring its TFTypeInfo.
func buildInventory(root, servicesDir string) (map[string]struct{}, error) {
	servicesPath := filepath.Join(root, servicesDir)
	inventory := map[string]struct{}{}

	entries, err := os.ReadDir(servicesPath)
	if err != nil {
		return nil, err
	}

	typeRe := regexp.MustCompile(`Type:\s*"([a-z_]+)"`)

	var missing []string

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		baseFile := filepath.Join(servicesPath, e.Name(), "base.go")

		data, readErr := os.ReadFile(baseFile)
		if readErr != nil {
			missing = append(missing, e.Name())

			continue
		}

		matches := typeRe.FindAllSubmatch(data, -1)
		for _, m := range matches {
			inventory[string(m[1])] = struct{}{}
		}
	}

	if len(missing) > 0 {
		toolutil.Errf("warning: %d service packages have no base.go: %s\n", len(missing), strings.Join(missing, ", "))
	}

	if len(inventory) == 0 {
		return nil, fmt.Errorf("no TF types found in %s — check services directory", servicesPath)
	}

	return inventory, nil
}

func report(gaps []gap) int {
	if *jsonFlag {
		return reportJSON(gaps)
	}

	return reportHuman(gaps)
}

func reportJSON(gaps []gap) int {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	err := enc.Encode(gaps)
	if err != nil {
		toolutil.Errf("error encoding JSON: %v\n", err)

		return exitErr
	}

	realGaps := countRealGaps(gaps)
	if realGaps > 0 {
		return exitGaps
	}

	return exitOK
}

func reportHuman(gaps []gap) int {
	var real, excluded, stale []gap

	for _, g := range gaps {
		switch {
		case g.Stale:
			stale = append(stale, g)
		case g.Skipped != "":
			excluded = append(excluded, g)
		default:
			real = append(real, g)
		}
	}

	if len(real) > 0 {
		toolutil.Outf("GAPS — SDK methods with no TF coverage and no exclusion:\n\n")

		prevClient := ""

		for _, g := range real {
			currClient := g.SDKMethod.Package + "/" + g.SDKMethod.Client
			if prevClient != "" && currClient != prevClient {
				toolutil.Outf("\n")
			}

			toolutil.Outf("  %s/%s.%s\n", g.SDKMethod.Package, g.SDKMethod.Client, g.SDKMethod.Method)

			prevClient = currClient
		}
	}

	if len(excluded) > 0 {
		toolutil.Outf("\nEXCLUDED (per exclusions.yaml): %d methods\n", len(excluded))
	}

	if len(stale) > 0 {
		toolutil.Outf("\nSTALE EXCLUSIONS — now implemented, remove from exclusions.yaml:\n\n")

		for _, g := range stale {
			toolutil.Outf("  %s/%s.%s\n", g.SDKMethod.Package, g.SDKMethod.Client, g.SDKMethod.Method)
		}
	}

	toolutil.Outf("\nSummary: %d gaps, %d excluded, %d stale\n", len(real), len(excluded), len(stale))

	if len(real) > 0 || len(stale) > 0 {
		return exitGaps
	}

	return exitOK
}

func countRealGaps(gaps []gap) int {
	count := 0

	for _, g := range gaps {
		if g.Skipped == "" || g.Stale {
			count++
		}
	}

	return count
}

// sdkFabricDir resolves the fabric-sdk-go module's fabric/ directory.
func sdkFabricDir(root string) (string, error) {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", toolutil.SDKModulePath)
	cmd.Dir = root

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("go list -m: %w", err)
	}

	modDir := strings.TrimSpace(string(out))
	fabricDir := filepath.Join(modDir, "fabric")

	if _, statErr := os.Stat(fabricDir); statErr != nil {
		return "", fmt.Errorf("fabric dir not found at %s: %w", fabricDir, statErr)
	}

	return fabricDir, nil
}

// pascalToSnake converts PascalCase to snake_case.
// It handles runs of uppercase letters (acronyms) by keeping them together:
// "GraphQLApi" → "graph_ql_api", "KQLDatabase" → "kql_database".
func pascalToSnake(s string) string {
	var result strings.Builder

	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if !unicode.IsUpper(r) {
			result.WriteRune(r)

			continue
		}

		// We have an uppercase rune. Determine if it starts an acronym run.
		if i > 0 {
			result.WriteByte('_')
		}

		// Collect the run of consecutive uppercase letters.
		j := i

		for j < len(runes) && unicode.IsUpper(runes[j]) {
			j++
		}

		runLen := j - i

		if runLen == 1 {
			// Single uppercase: just a normal word start.
			result.WriteRune(unicode.ToLower(r))
		} else if j == len(runes) {
			// Uppercase run goes to end of string — entire run is the acronym.
			for k := i; k < j; k++ {
				result.WriteRune(unicode.ToLower(runes[k]))
			}

			i = j - 1
		} else {
			// Uppercase run followed by lowercase: last uppercase char starts next word.
			// e.g. "SQLDatabase" → run is "SQLD", but "D" starts "Database".
			for k := i; k < j-1; k++ {
				result.WriteRune(unicode.ToLower(runes[k]))
			}

			result.WriteByte('_')
			result.WriteRune(unicode.ToLower(runes[j-1]))
			i = j - 1
		}
	}

	return result.String()
}
