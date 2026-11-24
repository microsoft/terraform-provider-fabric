---
producer: "review-code"
schema_version: 1
tier: "Low"
contract: "direct"
execution_backend: "copilot-cli"
attempt: 1
stage_attempt: null
round: 0
timestamp: "2026-11-24T10:43:57Z"
outcome: "issues_found"
reviewer_model: "claude-haiku-5.5"
input_sha256: "a829797e392fbadd4f746ead58b80e1d8b4dc4ce4993a73088f2b5d086a9c588"
---

## Provenance Summary

- **Tier**: Low
- **Round**: 0 (harness, not a lens)
- **Outcome**: issues_found

## Provenance Checklist

- [x] Input snapshot recorded: SHA-256 `a829797e392fbadd4f746ead58b80e1d8b4dc4ce4993a73088f2b5d086a9c588`.
- [x] Supplied harness response or failure diagnostic preserved below.

## Harness Response

## Review Summary

- **Round**: 0
- **Theme**: Low harness review
- **Mode**: resolved-plan
- **Model**: claude-haiku-5.5
- **Artifact**: C:\Users\orbauberg\.copilot\repos\copilot-worktrees\terraform-provider-fabric\orbaubergmicrosoft-automatic-guacamole\Features\1172\reviews\task-1172-attempt-1-review-0-harness.md
- **Issues Found**: 4
- **Verdict**: ISSUES_FOUND

## Evidence Checklist

- [x] Read resource code: `internal/services/onelakediagnosticsetting/resource_onelake_diagnostic_setting.go` (Read 107-137, Delete 172-205, ImportState 207-233, modify/get 235-265).
- [x] Read mapping code: `models.go` (state `set` 38-130, request `set` 165-234) and `schema.go` (23-130, singleton `workspace_id` with `RequiresReplace`).
- [x] Read tests: `resource_onelake_diagnostic_setting_test.go` (CRUD 88-159, step 2 at 121-134, ImportState 161-179, CreateError 229-246), `data_onelake_diagnostic_setting_test.go` (20-224), `fake_test.go` (18-104).
- [x] Verified helper semantics in `internal/pkg/utils/errors.go` lines 58-79 (not-found becomes a warning), 163-165 (404 maps to EntityNotFound).
- [x] Verified fixture: `tools/scripts/Set-WellKnown.ps1` 1376-1408 mirrors WorkspaceDS; line 1435 creates the well-known Lakehouse in WorkspaceDS, matching the acceptance destination.
- [x] Verified precedents: `internal/services/workspacencp/fake_test.go` 17-54 (inline singleton fake); 55 `examples/resources/*/import.sh` files; no `resources` template in `templates/`.
- [x] Supplied validation is consistent with source: 7 `TestUnit_` functions (4 resource, 1 data source, 2 Configure). Not re-executed; acceptance tests are excluded from that run.
- [ ] Not verified: live acceptance behavior and SDK `ModifyDiagnostics` polling internals (module source is outside the checkout and was not searched).

## Issues

### Issue 1: Disabled-state destination mapping rests on an unverified API reset contract

- **Severity**: Medium
- **File**: internal/services/onelakediagnosticsetting/models.go (also resource_onelake_diagnostic_setting.go)
- **Line(s)**: models.go 50-58, 60-77, 182-193; resource 99, 164, 200
- **Description**: `set()` writes any returned Lakehouse reference into state, with no status guard for Disabled (lines 60-77). It only suppresses the destination when the API omits it or the lakehouse is nil. The reset payload `{type: Lakehouse}` with no `lakehouse` (lines 189-191, 195-197) is proven only by the fake, which echoes whatever was stored (`fake_test.go` line 84). CRUD step 2 (test lines 121-134) passes only because of that echo.
- **Risk**: If the service retains the last Lakehouse reference after disable, or rejects the partial reset payload, Create/Update fail with "inconsistent result after apply" and refresh creates a perpetual diff. No unit test can detect this, and the supplied evidence contains no live acceptance run.
- **Suggested Fix**: Run `task testacc -- OneLakeDiagnosticSettingResource_CRUD ./internal/services/onelakediagnosticsetting/` against a tenant and record the result. Add a unit case where GET returns Disabled with a retained reference, then decide the mapping explicitly (for example, keep destination null when status is Disabled and config has none).

### Issue 2: Read not-found and successful-import paths are unit-untested

- **Severity**: Low
- **File**: internal/services/onelakediagnosticsetting/resource_onelake_diagnostic_setting_test.go
- **Line(s)**: 88-159 (CRUD), 161-179 (ImportState); source 123-129
- **Description**: No step removes `state.settings[workspaceID]`, so the `IsErrNotFound` → `RemoveResource` branch never runs. The `_ImportState` test covers only an invalid UUID; a successful import is implicit at CRUD step 3 (lines 135-141). Delete and Update error branches are also unexercised. I confirmed the not-found mapping by reading `utils/errors.go`, but nothing guards it.
- **Risk**: A regression in drift removal or import would pass CI. Coverage is 80.5%, just above the target.
- **Suggested Fix**: Add a RefreshState step whose PreConfig deletes `state.settings[workspaceID]`, expecting removal, plus a positive `_ImportState` case.

### Issue 3: Missing import example and Import docs section

- **Severity**: Low
- **File**: examples/resources/fabric_onelake_diagnostic_setting/ (and docs/resources/onelake_diagnostic_setting.md)
- **Line(s)**: resource 27 (implements `ResourceWithImportState`), 207-233
- **Description**: The resource supports import but has no `import.sh`. 55 resource examples have one, including the sibling `fabric_onelake_data_access_security/import.sh`. The generated doc has no `## Import` section.
- **Risk**: The import ID format (workspace ID alone) is undocumented, and generated docs diverge from convention.
- **Suggested Fix**: Add `import.sh` containing `terraform import fabric_onelake_diagnostic_setting.example "00000000-0000-0000-0000-000000000000"`, then run `task docs`.

### Issue 4: Reset-on-destroy semantics absent from generated docs

- **Severity**: Low
- **File**: internal/services/onelakediagnosticsetting/schema.go (also resource_onelake_diagnostic_setting.go, docs/resources/onelake_diagnostic_setting.md)
- **Line(s)**: schema.go 26; resource 188-191
- **Description**: Destroy disables diagnostics and clears the destination, but this is disclosed only as a runtime warning at destroy time. The description uses the generic `fabricitem.NewResourceMarkdownDescription`, and the generated doc omits the behavior. A `workspace_id` change (ForceNew) also disables the old workspace's diagnostics.
- **Risk**: Users treat this as a normal create/delete object and learn the reset semantics only during destroy.
- **Suggested Fix**: Add a note to the resource `MarkdownDescription`, for example "Destroying this resource disables OneLake diagnostics and clears the destination; the API has no delete operation." Then regenerate docs.

## Resolution Log

_Updated by the driving agent as findings are addressed._

### Issue 1
- **Status**: Fixed
- **What changed**: Documented that a disabled setting retains its destination when configured, while omitting the destination clears it. Added resource and data-source unit coverage for both states.
- **Why**: The Microsoft API examples distinguish disabling diagnostics from explicitly removing the destination; the provider's type-only destination request implements the documented removal operation.
- **How verified**: `task testunit -- OneLakeDiagnosticSetting ./internal/services/onelakediagnosticsetting/` passed 8 tests at 82.5% coverage. Confirmed request shapes against the [Modify Diagnostics API examples](https://learn.microsoft.com/rest/api/fabric/core/onelake-settings/modify-diagnostics). Live acceptance was not run because neither `internal/testhelp/fixtures/.wellknown.json` nor `FABRIC_TESTACC_WELLKNOWN` is configured.

### Issue 2
- **Status**: Fixed
- **What changed**: Added a refresh-state test for a 404 removing the resource from state, and a standalone successful-import test.
- **Why**: Protects the resource's drift-removal and workspace-ID import paths independently of the CRUD sequence.
- **How verified**: `task testunit -- OneLakeDiagnosticSetting ./internal/services/onelakediagnosticsetting/` passed 8 tests at 82.5% coverage.

### Issue 3
- **Status**: Fixed
- **What changed**: Added `examples/resources/fabric_onelake_diagnostic_setting/import.sh`; regenerated the resource documentation with its workspace-ID import syntax.
- **Why**: Keeps the import format discoverable and consistent with other importable resources.
- **How verified**: `task docs` completed successfully; the generated `docs/resources/onelake_diagnostic_setting.md` contains an Import section showing a workspace UUID.

### Issue 4
- **Status**: Fixed
- **What changed**: Expanded the resource schema description to explain that destroy disables diagnostics and clears the destination, and that replacing `workspace_id` resets the old workspace.
- **Why**: Users need to understand the singleton reset behavior before planning a destroy or replacement.
- **How verified**: `task docs` completed successfully; the generated resource documentation includes the reset and replacement behavior.

## Driver Decision

{
  "tier": "Low",
  "source": "Explicit human instruction",
  "changeSize": "Moderate: new workspace-scoped resource and read-only data source, provider wiring, well-known fixture, fake-backed tests, examples, generated documentation, and changelogs; no dependency changes.",
  "cyclomaticComplexity": "Moderate: stateful CRUD/reset semantics, long-running operation polling, workspace identity/import, and polymorphic nested destination mapping.",
  "rationale": "The user explicitly selected the Low tier. Low uses one script-owned harness review. Claude Haiku 5.5 is the exact user-named eligible model override. GPT-6 Luna is forbidden for direct review by the skill.",
  "primaryModel": "gpt-6-luna (known primary model; excluded for direct review)",
  "reviewers": [
    {
      "id": "claude-haiku-5.5",
      "family": "Anthropic",
      "reasoning": "max",
      "supportedLevels": ["low", "medium", "high", "xhigh", "max"],
      "selectionSource": "Explicit user model preference; current task-tool catalog and CLI capability references"
    }
  ],
  "overrideSource": "User request: using low tier models such as gpt luna 6 and claude haiku 5.5 as subagents",
  "modelCapabilityEvidence": "GitHub lists Claude Haiku 5.5 as GA. Anthropic documents adaptive thinking and suitability for subagent tasks. The current task-tool model catalog supports max reasoning for claude-haiku-5.5; Copilot CLI accepts max. The review runner verifies runtime model identity.",
  "zdrEvidence": "The review runner invokes GitHub Copilot CLI. GH_TOKEN variable presence and active GitHub API identity were checked without reading token contents; identity matched the current Copilot session account. GitHub documents ZDR for generally available Anthropic Copilot features; Haiku 5.5 is GA, and the listed Fable/preview exceptions do not apply.",
  "validation": {
    "unitTests": "task testunit -- OneLakeDiagnosticSetting ./internal/services/onelakediagnosticsetting/ passed: 7 tests, 80.5% package coverage.",
    "lint": "golangci-lint run ./internal/services/onelakediagnosticsetting reported 0 issues.",
    "changelog": "With CHANGIE_GITHUB_REPOSITORY=microsoft/terraform-provider-fabric, changie batch patch --dry-run --prerelease dev passed and listed issues 1171 and 1172.",
    "scriptSyntax": "PowerShell parser and git diff --check passed for Set-WellKnown.ps1 and current diff."
  },
  "taskId": "1172",
  "artifact": "Features/1172/reviews/task-1172-attempt-1-review-0-harness.md"
}
