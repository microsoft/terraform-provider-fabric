---
producer: "review-code"
schema_version: 1
tier: "Low"
contract: "direct"
execution_backend: "copilot-cli"
attempt: 1
stage_attempt: null
round: 0
timestamp: "2026-11-24T10:41:53Z"
outcome: "issues_found"
reviewer_model: "claude-haiku-5.5"
input_sha256: "a5a9650b9e2d0c7c41516e16ad9a4b91c884c065437fcfe1be7c1ba8975fb3ff"
---

## Provenance Summary

- **Tier**: Low
- **Round**: 0 (harness, not a lens)
- **Outcome**: issues_found

## Provenance Checklist

- [x] Input snapshot recorded: SHA-256 `a5a9650b9e2d0c7c41516e16ad9a4b91c884c065437fcfe1be7c1ba8975fb3ff`.
- [x] Supplied harness response or failure diagnostic preserved below.

## Harness Response

## Review Summary

- **Round**: 0
- **Theme**: Low harness review
- **Mode**: resolved-plan
- **Model**: claude-haiku-5.5
- **Artifact**: C:\Users\orbauberg\.copilot\repos\copilot-worktrees\terraform-provider-fabric\orbaubergmicrosoft-curly-sniffle\Features\1169\reviews\task-1169-attempt-1-review-0-harness.md
- **Issues Found**: 4
- **Verdict**: ISSUES_FOUND

## Evidence Checklist

- [x] Re-read the new resource (`internal/services/onelakedefaulttier/resource_onelake_default_tier.go` 1–235), data source `data_onelake_default_tier.go` 90–106, `models.go` 1–51, `schema.go` 1–92, `fake_test.go` 1–92, and `resource_onelake_default_tier_test.go` 1–200.
- [x] Verified the not-found contract end to end in `internal/pkg/utils/errors.go` (58–98, 163–165). Resource Read at lines 127–133 removes state on 404 correctly. Issue 3 is a test gap, not a behavior bug.
- [x] Checked repo rules: `.github/instructions/schema-model-patterns.instructions.md` 108–110 (`set()` on every model); `non-item-patterns.instructions.md` 69 (Read not-found uses `IsErrNotFound` then `RemoveResource`); `testing-patterns.instructions.md` 13 (`_ImportState` group) and 185–198 (404 fakes).
- [x] Confirmed example layout via glob: `examples/resources/fabric_workspace_onelake_default_tier/` has only `resource.tf` and `providers.tf`; `examples/resources/fabric_workspace_encryption/import.sh` exists.
- [x] Reviewed supplied diff hunks for provider registration, changelog entries #1169/#1170, and the `Set-WellKnown.ps1` WellKnown blocks.
- [ ] `Invoke-FabricRest(...).Response` contract is unverified: the helper comes from an external module not in this repo. Not flagged.
- [ ] I executed no commands. Supplied validation (9 unit tests, 82.7% coverage, golangci-lint, changie dry-run, Set-WellKnown parse, `git diff --check`) is taken as provided, not re-run. Acceptance tests and fixture provisioning were not run.

## Issues

### Issue 1: Model mappers do not follow the `set()` convention

- **Severity**: Low
- **File**: internal/services/onelakedefaulttier/models.go
- **Line(s)**: 21 (`setFromSettings`), 34 (`setFromTier`); call sites in resource_onelake_default_tier.go 104, 175, 231 and data_onelake_default_tier.go 103
- **Description**: The repo requires every model to have `set()` mapping the SDK DTO to the TF model. This package defines two differently named mappers. Sibling non-item models (`workspacencp`, `workspaceencryption`) use `set`.
- **Risk**: Convention drift. Behavior is correct today (nil tier maps to null).
- **Suggested Fix**: Replace both with `func (to *baseOneLakeDefaultTierModel) set(workspaceID string, tier *fabcore.OneLakeAccessTier)`, which sets ID, WorkspaceID, and `types.StringPointerValue((*string)(tier))`. The data source passes `respGet.Lifecycle.DefaultTier` after a nil check on `Lifecycle`, so nil still maps to null. The resource `get` keeps its existing nil checks.

### Issue 2: No import example or Import docs for the resource

- **Severity**: Low
- **File**: examples/resources/fabric_workspace_onelake_default_tier/ and docs/resources/workspace_onelake_default_tier.md
- **Line(s)**: resource_onelake_default_tier.go 192–210 (`ImportState`); the docs file has no `## Import` section
- **Description**: The resource implements `ImportState` (import ID = workspace ID), but there is no `import.sh`. The sibling `fabric_workspace_encryption` ships one, and tfplugindocs renders the Import section only from that file.
- **Risk**: Users cannot discover the import ID format from the generated docs.
- **Suggested Fix**: Add `examples/resources/fabric_workspace_onelake_default_tier/import.sh` in the same format as `fabric_workspace_encryption/import.sh`, using a workspace ID as the import ID. Regenerate docs with `task docs`.

### Issue 3: Resource Read 404 (state removal) has no test

- **Severity**: Low
- **File**: internal/services/onelakedefaulttier/resource_onelake_default_tier_test.go
- **Line(s)**: Code under test at resource_onelake_default_tier.go 127–133. No resource test references `fakeGetOneLakeSettingsError()` (fake_test.go 83–92).
- **Description**: The drift-removal branch (`IsErrNotFound` → `RemoveResource`) is untested for the resource. It depends on exact sentinel string matching in `utils`, so a regression would go unnoticed.
- **Risk**: Refresh could start erroring on deleted workspaces, or remove state on other errors, with no failing test.
- **Suggested Fix**: After a successful create, add a `RefreshState: true` step with `PreConfig` that swaps `GetSettings` to `fakeGetOneLakeSettingsError()`, expecting a non-empty plan and no error.

### Issue 4: `_ImportState` invalid-ID and Update error paths are untested

- **Severity**: Low
- **File**: internal/services/onelakedefaulttier/resource_onelake_default_tier_test.go
- **Line(s)**: Code under test at resource_onelake_default_tier.go 196–199 (invalid import ID) and 161–173 (Update error and nil-tier). Tests at 93–98 (success import inside `_CRUD`), 102–119 and 121–138 (modify error and missing response, Create only).
- **Description**: `testing-patterns` expects an `_ImportState` group covering invalid formats. The only import exercise is a success step inside `_CRUD`. The modify error and missing-response tests run only on Create, so the Update branches are never reached.
- **Risk**: Regressions in import validation or Update error handling would pass the suite despite 82.7% coverage.
- **Suggested Fix**: Add `TestUnit_OneLakeDefaultTierResource_ImportState` with `ImportStateId: "invalid-uuid"` and an `ExpectError` matching the UUID error the implementation emits. Add an Update case that creates successfully, then swaps `ModifyDefaultTier` to `fakeModifyDefaultTierError()` or `fakeModifyDefaultTierMissingResponse()` while changing `default_tier`.

## Resolution Log

_Updated by the driving agent as findings are addressed._

### Issue 1
- **Status**: Fixed
- **What changed**: Replaced the two differently named mappers with `baseOneLakeDefaultTierModel.set` and updated resource/data-source mapping call sites.
- **Why**: Aligns the model with the repository's `set()` convention while preserving null handling for missing lifecycle tiers.
- **How verified**: `task testunit -- OneLakeDefaultTier ./internal/services/onelakedefaulttier/` passed (13 tests); `golangci-lint run ./internal/services/onelakedefaulttier/... ./internal/provider/...` reported 0 issues.

### Issue 2
- **Status**: Fixed
- **What changed**: Added the resource `import.sh` example and regenerated provider docs with an Import section showing a Workspace ID as the import ID.
- **Why**: Makes the implemented import format discoverable to users.
- **How verified**: `task docs` completed; `docs/resources/workspace_onelake_default_tier.md` now includes the `## Import` section and example.

### Issue 3
- **Status**: Fixed
- **What changed**: Added a resource refresh test that returns `EntityNotFound` after creation and expects a non-empty plan.
- **Why**: Verifies the documented 404 behavior removes the resource from state instead of failing refresh.
- **How verified**: `task testunit -- OneLakeDefaultTier ./internal/services/onelakedefaulttier/` passed (13 tests).

### Issue 4
- **Status**: Fixed
- **What changed**: Added invalid UUID import coverage plus update error and missing-response tests.
- **Why**: Covers the resource's import validation and both Update failure branches, not only Create.
- **How verified**: `task testunit -- OneLakeDefaultTier ./internal/services/onelakedefaulttier/` passed (13 tests, 88.5% statement coverage).

## Driver Decision

# Review decision

Tier: Low
Source: Explicit human invocation argument requests a low-tier review and names Claude Haiku 5.5; the review skill specifies one script-owned Low review.
Change size: Moderate — workspace-scoped resource and data source, tests and fake server, provider registration, test fixtures, examples, generated docs, changelog.
Cyclomatic complexity: Moderate — read/write lifecycle, import, unsupported remote reset, optional lifecycle response, stateful fake behavior, and workspace-wide file migration/billing consequences.
Rationale: Use exactly one read-only Copilot CLI review with the exact eligible model claude-haiku-5.5. The user mentioned GPT Luna, but the review skill prohibits every Luna model. Low tier was explicitly requested despite moderate change complexity. The low contract uses one harness reviewer, not multiple task-tool subagents.
Primary model: GPT-6 Luna (current parent runtime; excluded from review dispatch).
Reviewer: claude-haiku-5.5 (exact Copilot CLI model ID).
Reasoning: The live task/harness catalog advertises low, medium, high, xhigh, max for Claude Haiku 5.5; Copilot CLI accepts these reasoning values; Anthropic's effort docs list Haiku 5.5 supporting xhigh and max. Pass the full shared set so the runner selects max.
Override source: Human invocation arguments; exact Claude Haiku 5.5 preference; Luna suggestion rejected under the skill's absolute exclusion.
ZDR: Copilot CLI auth uses GH_TOKEN (COPILOT_GITHUB_TOKEN absent, GH_TOKEN present, GH_HOST=github.com); GitHub auth reports the current Microsoft-managed account. GitHub's model-hosting documentation covers generally available Anthropic Copilot features under ZDR except Claude Fable models; Haiku 5.5 is GA and outside that exception. No direct vendor API will be used.
Scope: Review the current uncommitted tracked and untracked workspace changes against HEAD, including the OneLake default-tier resource/data source, workspace test setup, examples, generated docs, registration, and changelog. Expected semantics: modify the workspace default tier using Hot/Cool/Cold; read lifecycle.defaultTier; report absent data-source tier as null; fail resource refresh if the tier becomes unavailable; delete removes only Terraform state and warns; separate empty workspaces isolate acceptance tests because tier changes can migrate files and incur charges.
User request: Do not post PR comments or GitHub artifacts. Review and identify actionable issues for local fixes only; the required review evidence remains local to this branch.
Validation: `task testunit -- OneLakeDefaultTier ./internal/services/onelakedefaulttier/` passed (9 tests, 82.7% coverage). `golangci-lint run ./internal/services/onelakedefaulttier/... ./internal/provider/...` passed (0 issues; deprecated-config warnings). `CHANGIE_GITHUB_REPOSITORY=microsoft/terraform-provider-fabric changie batch patch --dry-run --prerelease dev` passed with entries for #1169 and #1170. `Set-WellKnown.ps1` PowerShell parser and `git diff --check` passed. Live acceptance tests and fixture provisioning were not run. Earlier `tfproviderlintx` invocation failed with an internal missing-types error and was not rerun.
Artifact: `Features/1169/reviews/task-1169-attempt-1-review-0-harness.md`.
