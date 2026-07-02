// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/terraform-provider-fabric/tools/internal/toolutil"
)

func Test_pascalToSnake(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{"Workspace", "workspace"},
		{"WorkspaceRoleAssignment", "workspace_role_assignment"},
		{"GitOutboundPolicy", "git_outbound_policy"},
		{"NetworkCommunicationPolicy", "network_communication_policy"},
		{"ManagedPrivateEndpoint", "managed_private_endpoint"},
		{"ExternalDataShare", "external_data_share"},
		{"Lakehouse", "lakehouse"},
		// Acronym runs.
		{"KQLDatabase", "kql_database"},
		{"KQLDashboard", "kql_dashboard"},
		{"KQLQueryset", "kql_queryset"},
		{"SQLDatabase", "sql_database"},
		{"SQLEndpoint", "sql_endpoint"},
		{"MLExperiment", "ml_experiment"},
		{"MLModel", "ml_model"},
		// Trailing acronym (entire run to end of string).
		{"GraphQLAPI", "graph_qlapi"},
		// Mixed case.
		{"OutboundCloudConnectionRules", "outbound_cloud_connection_rules"},
		{"InboundAzureResourceRules", "inbound_azure_resource_rules"},
		{"CosmosDBDatabase", "cosmos_db_database"},
		{"OneLakeDataAccessSecurity", "one_lake_data_access_security"},
		// Simple.
		{"Folder", "folder"},
		{"Domain", "domain"},
		{"Gateway", "gateway"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			got := pascalToSnake(tt.input)
			if got != tt.want {
				t.Errorf("pascalToSnake(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func Test_splitVerbNoun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		method   string
		wantVerb string
		wantNoun string
	}{
		{"CreateWorkspace", "Create", "Workspace"},
		{"DeleteWorkspace", "Delete", "Workspace"},
		{"GetGitOutboundPolicy", "Get", "GitOutboundPolicy"},
		{"SetInboundAzureResourceRules", "Set", "InboundAzureResourceRules"},
		{"AddWorkspaceRoleAssignment", "Add", "WorkspaceRoleAssignment"},
		{"ListWorkspaces", "List", "Workspaces"},
		{"UpdateWorkspace", "Update", "Workspace"},
		{"CreateOrUpdateDataAccessRoles", "CreateOrUpdate", "DataAccessRoles"},
		// Action verbs.
		{"BeginCreateWorkspace", "Begin", "CreateWorkspace"},
		{"TestConnection", "Test", "Connection"},
		{"NewListItemsPager", "New", "ListItemsPager"},
		// Formerly action verbs: no longer recognized, so they parse to no verb
		// and surface as gaps to be accounted for via exclusions.yaml.
		{"AssignToCapacity", "", ""},
		{"ApplyTags", "", ""},
		{"DeployStageContent", "", ""},
		// Unrecognized.
		{"", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			t.Parallel()

			gotVerb, gotNoun := splitVerbNoun(tt.method)
			if gotVerb != tt.wantVerb || gotNoun != tt.wantNoun {
				t.Errorf("splitVerbNoun(%q) = (%q, %q), want (%q, %q)",
					tt.method, gotVerb, gotNoun, tt.wantVerb, tt.wantNoun)
			}
		})
	}
}

func Test_classifyNounKind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		verbs map[string]struct{}
		want  string
	}{
		{"full CRUD", map[string]struct{}{"Create": {}, "Delete": {}, "Get": {}, "Update": {}}, "resource"},
		{"add/delete", map[string]struct{}{"Add": {}, "Delete": {}, "Get": {}}, "resource"},
		{"get/set singleton", map[string]struct{}{"Get": {}, "Set": {}}, "singleton"},
		{"createOrUpdate + delete", map[string]struct{}{"CreateOrUpdate": {}, "Delete": {}, "Get": {}}, "resource"},
		{"get only", map[string]struct{}{"Get": {}}, ""},
		{"list only", map[string]struct{}{"List": {}}, ""},
		{"update only", map[string]struct{}{"Update": {}}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := classifyNounKind(tt.verbs)
			if got != tt.want {
				t.Errorf("classifyNounKind(%v) = %q, want %q", tt.verbs, got, tt.want)
			}
		})
	}
}

func Test_itemNounCandidates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		method string
		want   []string
	}{
		// Compound package dirs whose lowercase name can't snake-convert, but
		// whose method noun can.
		{"ListPaginatedReports", []string{"paginated_report", "paginated_reports"}},
		{"ListMirroredWarehouses", []string{"mirrored_warehouse", "mirrored_warehouses"}},
		// Singular noun yields a single candidate.
		{"GetPaginatedReport", []string{"paginated_report"}},
		// Override takes precedence over derivation.
		{"ListReflexes", []string{"activator"}},
		// Unrecognized method shape yields no candidates.
		{"Frobnicate", nil},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			t.Parallel()

			got := itemNounCandidates(tt.method)
			if !slices.Equal(got, tt.want) {
				t.Errorf("itemNounCandidates(%q) = %v, want %v", tt.method, got, tt.want)
			}
		})
	}
}

func Test_findGaps_actionVerbsNotReported(t *testing.T) {
	t.Parallel()

	// Action-style methods on a client with no CRUD lifecycle must never be
	// reported as gaps, even when neither covered nor excluded.
	allMethods := []sdkMethod{
		{Package: "core", Client: "ItemJobSchedulesClient", Method: "RunOnDemandItemJob"},
		{Package: "core", Client: "ItemJobSchedulesClient", Method: "CancelItemJobInstance"},
		// A genuine uncovered resource method should still be reported.
		{Package: "core", Client: "WidgetsClient", Method: "CreateWidget"},
	}

	gaps := findGaps(nil, allMethods, map[string]struct{}{}, map[string]exclusionEntry{}, map[string]struct{}{})

	reported := map[string]struct{}{}
	for _, g := range gaps {
		reported[g.SDKMethod.Method] = struct{}{}
	}

	for _, action := range []string{"RunOnDemandItemJob", "CancelItemJobInstance"} {
		if _, ok := reported[action]; ok {
			t.Errorf("action method %q was reported as a gap, want it skipped", action)
		}
	}

	if _, ok := reported["CreateWidget"]; !ok {
		t.Error("control-plane method CreateWidget was not reported as a gap")
	}
}

// findGapFor returns the gap reported for a given method name, if any.
func findGapFor(gaps []gap, method string) (gap, bool) {
	for _, g := range gaps {
		if g.SDKMethod.Method == method {
			return g, true
		}
	}

	return gap{}, false
}

func Test_findGaps_nounCoverage(t *testing.T) {
	t.Parallel()

	// A multi-resource client (WorkspacesClient) covers individual nouns. A
	// method whose noun matches the inventory (directly or via singular/plural
	// variant) must not be reported; an unrelated noun still is.
	nounGroups := []nounGroup{
		{Package: "core", Client: "WorkspacesClient", Noun: "Workspace"},
	}
	inventory := map[string]struct{}{"workspace": {}}
	allMethods := []sdkMethod{
		{Package: "core", Client: "WorkspacesClient", Method: "CreateWorkspace"},
		{Package: "core", Client: "WorkspacesClient", Method: "ListWorkspaces"},
		{Package: "core", Client: "WorkspacesClient", Method: "CreateGateway"},
	}

	gaps := findGaps(nounGroups, allMethods, inventory, map[string]exclusionEntry{}, map[string]struct{}{})

	for _, covered := range []string{"CreateWorkspace", "ListWorkspaces"} {
		if _, ok := findGapFor(gaps, covered); ok {
			t.Errorf("method %q noun is covered, want it absent from gaps", covered)
		}
	}

	if _, ok := findGapFor(gaps, "CreateGateway"); !ok {
		t.Error("CreateGateway noun is not covered, want it reported as a gap")
	}
}

func Test_findGaps_itemsClientCoverage(t *testing.T) {
	t.Parallel()

	// For a per-item ItemsClient, a single inventory match covers the entire
	// client (one resource == one ItemsClient), so all its methods are covered.
	nounGroups := []nounGroup{
		{Package: "lakehouse", Client: "ItemsClient", Noun: "Lakehouse"},
	}
	inventory := map[string]struct{}{"lakehouse": {}}
	allMethods := []sdkMethod{
		{Package: "lakehouse", Client: "ItemsClient", Method: "GetLakehouse"},
		{Package: "lakehouse", Client: "ItemsClient", Method: "GetLakehouseDefinition"},
		{Package: "lakehouse", Client: "ItemsClient", Method: "ListLakehouses"},
	}

	gaps := findGaps(nounGroups, allMethods, inventory, map[string]exclusionEntry{}, map[string]struct{}{})

	if len(gaps) != 0 {
		t.Errorf("all ItemsClient methods should be covered, got %d gaps: %+v", len(gaps), gaps)
	}
}

func Test_findGaps_calledMethodCoverage(t *testing.T) {
	t.Parallel()

	// A method directly resolved as called in service code (calledMethods) is
	// covered even without any noun-group/inventory match.
	allMethods := []sdkMethod{
		{Package: "core", Client: "DeploymentPipelinesClient", Method: "UpdateDeploymentPipelineStage"},
		{Package: "core", Client: "DeploymentPipelinesClient", Method: "CreateDeploymentPipeline"},
	}
	calledMethods := map[string]struct{}{
		"core/DeploymentPipelinesClient/UpdateDeploymentPipelineStage": {},
	}

	gaps := findGaps(nil, allMethods, map[string]struct{}{}, map[string]exclusionEntry{}, calledMethods)

	if _, ok := findGapFor(gaps, "UpdateDeploymentPipelineStage"); ok {
		t.Error("UpdateDeploymentPipelineStage is called in service code, want it absent from gaps")
	}

	if _, ok := findGapFor(gaps, "CreateDeploymentPipeline"); !ok {
		t.Error("CreateDeploymentPipeline is uncovered, want it reported as a gap")
	}
}

func Test_findGaps_exclusions(t *testing.T) {
	t.Parallel()

	allMethods := []sdkMethod{
		{Package: "core", Client: "CatalogClient", Method: "GetCatalog"},
		{Package: "core", Client: "CatalogClient", Method: "ListCatalogs"},
		{Package: "core", Client: "OneLakeSettingsClient", Method: "CreateOneLakeSetting"},
	}
	exclusions := map[string]exclusionEntry{
		// Client-level wildcard suppresses every method on the client.
		"core/CatalogClient/*": {Package: "core", Client: "CatalogClient", Method: "*", Reason: "read-only catalog"},
		// Exact-method exclusion.
		"core/OneLakeSettingsClient/CreateOneLakeSetting": {
			Package: "core", Client: "OneLakeSettingsClient", Method: "CreateOneLakeSetting", Reason: "action-style",
		},
	}

	gaps := findGaps(nil, allMethods, map[string]struct{}{}, exclusions, map[string]struct{}{})

	for _, m := range []string{"GetCatalog", "ListCatalogs", "CreateOneLakeSetting"} {
		g, ok := findGapFor(gaps, m)
		if !ok {
			t.Errorf("excluded method %q should still appear (as skipped), but is absent", m)

			continue
		}

		if g.Skipped == "" {
			t.Errorf("method %q should carry a skip reason, got empty", m)
		}

		if g.Stale {
			t.Errorf("excluded-but-uncovered method %q must not be marked stale", m)
		}
	}
}

func Test_findGaps_staleExclusion(t *testing.T) {
	t.Parallel()

	// An exclusion for a method that is actually covered (here via a resolved
	// service call) is stale and must be flagged so it can be removed.
	allMethods := []sdkMethod{
		{Package: "core", Client: "WorkspacesClient", Method: "CreateWorkspace"},
	}
	calledMethods := map[string]struct{}{
		"core/WorkspacesClient/CreateWorkspace": {},
	}
	exclusions := map[string]exclusionEntry{
		"core/WorkspacesClient/CreateWorkspace": {
			Package: "core", Client: "WorkspacesClient", Method: "CreateWorkspace", Reason: "was skipped, now implemented",
		},
	}

	gaps := findGaps(nil, allMethods, map[string]struct{}{}, exclusions, calledMethods)

	g, ok := findGapFor(gaps, "CreateWorkspace")
	if !ok {
		t.Fatal("stale exclusion should be reported")
	}

	if !g.Stale {
		t.Errorf("covered + excluded method must be marked Stale, got Stale=%v", g.Stale)
	}
}

func Test_findGaps_readOnlyMethodsAreReported(t *testing.T) {
	t.Parallel()

	// The gap report is intentionally exhaustive: uncovered read-only (Get/List)
	// and update-only methods MUST be reported so every SDK operation is
	// accounted for (and pruned via exclusions.yaml, not by the classifier).
	// This guards against narrowing findGaps to classified control-plane groups.
	allMethods := []sdkMethod{
		{Package: "core", Client: "DomainsClient", Method: "GetDomain"},
		{Package: "core", Client: "DomainsClient", Method: "ListDomains"},
		{Package: "core", Client: "TagsClient", Method: "UpdateTag"},
	}

	gaps := findGaps(nil, allMethods, map[string]struct{}{}, map[string]exclusionEntry{}, map[string]struct{}{})

	for _, m := range []string{"GetDomain", "ListDomains", "UpdateTag"} {
		g, ok := findGapFor(gaps, m)
		if !ok {
			t.Errorf("read/update-only method %q must be reported as a gap (exhaustive report)", m)

			continue
		}

		if g.Skipped != "" || g.Stale {
			t.Errorf("uncovered method %q should be a plain gap, got skipped=%q stale=%v", m, g.Skipped, g.Stale)
		}
	}
}

func Test_receiverTypeName(t *testing.T) {
	t.Parallel()

	pkg := types.NewPackage("github.com/microsoft/fabric-sdk-go/fabric/core", "core")
	client := types.NewNamed(
		types.NewTypeName(token.NoPos, pkg, "WorkspacesClient", nil),
		types.NewStruct(nil, nil),
		nil,
	)

	tests := []struct {
		name string
		typ  types.Type
		want string
	}{
		{"pointer receiver", types.NewPointer(client), "WorkspacesClient"},
		{"value receiver", client, "WorkspacesClient"},
		{"non-named type", types.Typ[types.Int], ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := receiverTypeName(tt.typ); got != tt.want {
				t.Errorf("receiverTypeName(%s) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}

func Test_buildCandidates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ng   nounGroup
		want []string
	}{
		{
			"explicit override wins",
			nounGroup{Client: "CosmosDBDatabasesClient", Noun: "CosmosDBDatabase"},
			[]string{"cosmos_db"},
		},
		{
			"simple item noun (singular == plural stem)",
			nounGroup{Client: "ItemsClient", Noun: "Lakehouse"},
			[]string{"lakehouse", "item_lakehouse"},
		},
		{
			"plural list noun yields singular + plural + prefixed",
			nounGroup{Client: "WorkspacesClient", Noun: "InboundAzureResourceRules"},
			[]string{
				"inbound_azure_resource_rule",
				"inbound_azure_resource_rules",
				"workspace_inbound_azure_resource_rule",
				"workspace_inbound_azure_resource_rules",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := buildCandidates(tt.ng)
			if !slices.Equal(got, tt.want) {
				t.Errorf("buildCandidates(%+v) = %v, want %v", tt.ng, got, tt.want)
			}
		})
	}
}

// Test_extractCoveredMethods_integration exercises the full go/packages
// type-aware coverage detection against the real service packages. It is the
// regression guard for the migration away from regex line-scanning: it proves
// that a call made through a bare local variable (not r.client.X) is resolved.
// It type-checks the whole dependency graph, so it is skipped under -short.
// RequireSDK turns a missing module-root / SDK condition into a hard failure
// when GAPCHECK_REQUIRE_SDK is set (as CI must set it), so the SDK-backed guard
// tests cannot silently pass by skipping. Locally, without the env var, the
// test still skips when the SDK is not resolvable. If err is non-nil this never
// returns (it calls t.Fatalf or t.Skipf, both of which stop the test).
func requireSDK(t *testing.T, err error, what string) {
	t.Helper()

	if err == nil {
		return
	}

	if os.Getenv("GAPCHECK_REQUIRE_SDK") != "" {
		t.Fatalf("GAPCHECK_REQUIRE_SDK is set but %s: %v", what, err)
	}

	t.Skipf("%s: %v", what, err)
}

func Test_extractCoveredMethods_integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping type-checking integration test in -short mode")
	}

	t.Parallel()

	root, err := toolutil.ModuleRoot()
	requireSDK(t, err, "find module root")

	covered, err := extractCoveredMethods(root, toolutil.DefaultServicesDir)
	if err != nil {
		t.Fatalf("extractCoveredMethods: %v", err)
	}

	if len(covered) == 0 {
		t.Fatal("extractCoveredMethods returned no covered methods")
	}

	// Every key must be "sdkPackage/ClientType/Method" so it matches
	// sdkMethod.sdkMethodExclKey used by findGaps.
	for key := range covered {
		if strings.Count(key, "/") != 2 {
			t.Errorf("covered key %q is not in pkg/Client/Method form", key)
		}
	}

	// Methods that must resolve. UpdateDeploymentPipelineStage is called through
	// a bare `client` variable — the exact shape the old regex scanner missed
	// and the reason for the go/packages migration.
	mustCover := []string{
		"core/WorkspacesClient/CreateWorkspace",
		"core/DeploymentPipelinesClient/UpdateDeploymentPipelineStage",
	}
	for _, key := range mustCover {
		if _, ok := covered[key]; !ok {
			t.Errorf("expected %q to be detected as covered, but it was not", key)
		}
	}
}

func Test_loadExclusions(t *testing.T) {
	t.Parallel()

	write := func(t *testing.T, content string) string {
		t.Helper()

		path := filepath.Join(t.TempDir(), "exclusions.yaml")

		err := os.WriteFile(path, []byte(content), 0o600)
		if err != nil {
			t.Fatalf("writing temp exclusions file: %v", err)
		}

		return path
	}

	t.Run("valid entry with explicit method", func(t *testing.T) {
		t.Parallel()

		path := write(t, "exclusions:\n  - package: core\n    client: CatalogClient\n    method: GetCatalog\n    reason: read-only\n")

		got, err := loadExclusions(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, ok := got["core/CatalogClient/GetCatalog"]; !ok {
			t.Errorf("expected key core/CatalogClient/GetCatalog, got %v", got)
		}
	})

	t.Run("empty method normalized to wildcard", func(t *testing.T) {
		t.Parallel()

		path := write(t, "exclusions:\n  - package: core\n    client: CatalogClient\n    reason: whole client excluded\n")

		got, err := loadExclusions(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		entry, ok := got["core/CatalogClient/*"]
		if !ok {
			t.Fatalf("expected wildcard key core/CatalogClient/*, got %v", got)
		}

		if entry.Method != "*" {
			t.Errorf("Method = %q, want %q", entry.Method, "*")
		}
	})

	invalid := []struct {
		name    string
		content string
	}{
		{"missing reason", "exclusions:\n  - package: core\n    client: CatalogClient\n    method: GetCatalog\n"},
		{"missing client", "exclusions:\n  - package: core\n    method: GetCatalog\n    reason: x\n"},
		{"missing package", "exclusions:\n  - client: CatalogClient\n    method: GetCatalog\n    reason: x\n"},
		{"malformed yaml (scalar where list expected)", "exclusions: not-a-list\n"},
	}

	for _, tt := range invalid {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := write(t, tt.content)

			_, err := loadExclusions(path)
			if err == nil {
				t.Errorf("expected error for %q, got nil", tt.name)
			}
		})
	}

	t.Run("nonexistent path", func(t *testing.T) {
		t.Parallel()

		_, err := loadExclusions(filepath.Join(t.TempDir(), "does-not-exist.yaml"))
		if err == nil {
			t.Error("expected error for nonexistent path, got nil")
		}
	})
}

// Test_nounTypeOverrides_notStale verifies that every key in nounTypeOverrides
// corresponds to an actual noun extracted from the SDK. If the SDK renames or
// removes a method, this test flags the stale override entry.
func Test_nounTypeOverrides_notStale(t *testing.T) {
	t.Parallel()

	root, err := toolutil.ModuleRoot()
	requireSDK(t, err, "find module root")

	sdkDir, err := sdkFabricDir(root)
	requireSDK(t, err, "resolve SDK")

	// Collect all nouns from the SDK.
	allNouns := collectAllSDKNouns(t, sdkDir)

	for key := range nounTypeOverrides {
		if _, ok := allNouns[key]; !ok {
			t.Errorf("nounTypeOverrides key %q does not match any SDK noun — stale entry", key)
		}
	}
}

// Test_exclusionsFile_notStale verifies that every entry in exclusions.yaml
// corresponds to an actual noun group in the SDK. Stale entries indicate the
// SDK removed the operation.
func Test_exclusionsFile_notStale(t *testing.T) {
	t.Parallel()

	root, err := toolutil.ModuleRoot()
	requireSDK(t, err, "find module root")

	sdkDir, err := sdkFabricDir(root)
	requireSDK(t, err, "resolve SDK")

	exclusionsPath := filepath.Join(root, "tools", "gapcheck", "exclusions.yaml")

	exclusions, err := loadExclusions(exclusionsPath)
	if err != nil {
		t.Fatalf("loadExclusions: %v", err)
	}

	allMethods, err := extractAllMethods(sdkDir)
	if err != nil {
		t.Fatalf("extractAllMethods: %v", err)
	}

	// Build set of all method keys and client keys from the SDK.
	allMethodKeys := map[string]struct{}{}
	allClientKeys := map[string]struct{}{}

	for _, m := range allMethods {
		allMethodKeys[m.sdkMethodExclKey()] = struct{}{}
		allClientKeys[m.Package+"/"+m.Client+"/*"] = struct{}{}
	}

	for key, entry := range exclusions {
		if entry.Method == "*" {
			// Wildcard: check that the client exists.
			if _, ok := allClientKeys[key]; !ok {
				t.Errorf("exclusions.yaml entry %q does not match any SDK client — stale entry", key)
			}
		} else {
			// Exact method: check that it exists.
			if _, ok := allMethodKeys[key]; !ok {
				t.Errorf("exclusions.yaml entry %q does not match any SDK method — stale entry", key)
			}
		}
	}
}

// collectAllSDKNouns extracts all noun names from SDK client methods across all packages.
func collectAllSDKNouns(t *testing.T, sdkDir string) map[string]struct{} {
	t.Helper()

	nouns := map[string]struct{}{}

	packages := []string{"core", "admin"}

	// Add per-item packages.
	entries, err := os.ReadDir(sdkDir)
	if err != nil {
		t.Fatalf("reading SDK dir: %v", err)
	}

	for _, e := range entries {
		if e.IsDir() && e.Name() != "core" && e.Name() != "admin" && e.Name() != "fake" {
			packages = append(packages, e.Name())
		}
	}

	for _, pkg := range packages {
		pkgDir := filepath.Join(sdkDir, pkg)

		files, readErr := os.ReadDir(pkgDir)
		if readErr != nil {
			continue
		}

		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), "_client.go") || strings.Contains(f.Name(), "example") {
				continue
			}

			data, fileErr := os.ReadFile(filepath.Join(pkgDir, f.Name()))
			if fileErr != nil {
				continue
			}

			for line := range strings.SplitSeq(string(data), "\n") {
				_, method := parseClientMethod(line)
				if method == "" {
					continue
				}

				verb, noun := splitVerbNoun(method)
				if verb == "Begin" {
					// Strip Begin and re-parse the inner method.
					_, noun = splitVerbNoun(noun)
				}

				if noun != "" {
					nouns[noun] = struct{}{}
				}
			}
		}
	}

	return nouns
}
