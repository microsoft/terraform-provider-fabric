---
producer: "review-code"
schema_version: 1
tier: "Low"
contract: "direct"
execution_backend: "copilot-cli"
attempt: 1
stage_attempt: null
round: 0
timestamp: "2026-11-24T10:56:38Z"
outcome: "issues_found"
reviewer_model: "claude-haiku-5.5"
input_sha256: "8e1d2ca071f8a9e1f95253618e83574e435e06bd2db70565834e3ced2369c60e"
---

## Provenance Summary
- **Tier**: Low
- **Round**: 0 (harness, not a lens)
- **Outcome**: issues_found

## Provenance Checklist
- [x] Input snapshot recorded: SHA-256 `8e1d2ca071f8a9e1f95253618e83574e435e06bd2db70565834e3ced2369c60e`.
- [x] Supplied harness response or failure diagnostic preserved below.

## Harness Response

## Review Summary
- **Round**: 0
- **Theme**: Low harness review
- **Mode**: resolved-plan
- **Model**: claude-haiku-5.5
- **Artifact**: C:\Users\orbauberg\.copilot\repos\copilot-worktrees\terraform-provider-fabric\orbaubergmicrosoft-curly-happiness\Features\1167\reviews\task-1167-attempt-1-review-0-harness.md
- **Issues Found**: 7
- **Verdict**: ISSUES_FOUND

## Evidence Checklist
- [x] Resource source read: `internal/services/workspaceonelakeimmutability/resource_workspace_onelake_immutability_policy.go` (ModifyPlan 75-108, Create 110-138, Read 140-179, Delete 211-227, ImportState 229-264, modify 266-314, validateRetentionIncrease 327-333) and `helper.go` (getImmutabilityPolicy 15-47).
- [x] Schema, models, base, and data source read (`schema.go`, `models.go`, `base.go`, `data_workspace_onelake_immutability_policy.go`). No defects found. Descriptions and ForceNew markers match the two docs pages.
- [x] Tests and fake read: resource test (RetentionCannotDecrease 119-153, ReadRemovesMissingPolicy 213-242, acceptance 301-305), data source test (acceptance 150-172), `fake_test.go` (GET 26-56, modify 58-88).
- [x] Precedent compared: `internal/services/workspacencp/resource_workspace_network_communication_policy.go` (Read 117-158, IsErrNotFound 137-143, get 244-253) and `internal/pkg/utils/errors.go` (processFabricError 149-165).
- [x] Repository rules checked: `.github\instructions\non-item-patterns.instructions.md` line 69 and `.github\instructions\testing-patterns.instructions.md` lines 19-20.
- [x] `tools/scripts/Set-WellKnown.ps1` read: Write-Log 1-34, Invoke-FabricRest 102-194, Set-OneLakeDiagnostic 626-671, Set-OneLakeImmutabilityPolicy 673-703, new workspace/Lakehouse block 1491-1517 (its `Set-FabricItem -Type` call matches the existing loop at 1526), policy step 2460-2468, fixture save 2470-2473.
- [x] Supplied validation evidence inspected, not re-executed: build, unit (14 tests, 81.6%), golangci-lint (0 issues), `task docs`, terraform fmt, changie dry-run. The `tfproviderlintx` crypto error is tooling, not a finding.
- [x] `provider.go` registration hunks and both changelog YAMLs (Issues 1167 and 1168) are consistent with conventions in the diff.
- [ ] `configure_test.go` and `base_test.go` not reviewed: not opened on disk, and the supplied diff truncates before them.

## Issues

### Issue 1: Policy step runs before the well-known fixture is saved
- **Severity**: High
- **File**: tools/scripts/Set-WellKnown.ps1
- **Line(s)**: 2460-2468 (policy step), 2470-2473 (fixture save); consumers at internal/services/workspaceonelakeimmutability/resource_workspace_onelake_immutability_policy_test.go 302-305 and data_workspace_onelake_immutability_policy_test.go 151-154
- **Description**: `Set-OneLakeImmutabilityPolicy` runs before `.wellknown.json` is written. The script's own comment (line 2461) says `WorkspaceNotSupported` can persist for up to an hour. On that run, `Write-Log ... ERROR` exits, so the fixture is never written, including all pre-existing entries. The acceptance tests read the new key through unchecked type assertions.
- **Risk**: One first-run failure discards the well-known fixtures for every service, not just this feature. The acceptance tests then panic on the missing key instead of failing with a clear message.
- **Suggested Fix**: Write the fixture before the immutability step, or move that block after the save and save again once it succeeds. Use checked type assertions (`, ok`) in both acceptance tests.

### Issue 2: Read does not remove state when the settings GET returns 404
- **Severity**: Medium
- **File**: internal/services/workspaceonelakeimmutability/helper.go; internal/services/workspaceonelakeimmutability/resource_workspace_onelake_immutability_policy.go
- **Line(s)**: helper.go 17; resource Read 140-179 (error path 159-168)
- **Description**: `getImmutabilityPolicy` passes `nil` as `errIs` to `utils.GetDiagsFromError`. In `processFabricError` (internal/pkg/utils/errors.go 149-165), the 404 mapping applies only when `errIs != nil` (line 163). A 404 from `GET workspaces/{id}/onelake/settings`, for example after the workspace is deleted outside Terraform, becomes a generic read error. The precedent `workspacencp` passes `fabcore.ErrCommon.EntityNotFound` (line 248) and calls `utils.IsErrNotFound` then `RemoveResource` (lines 137-143). `non-item-patterns.instructions.md` line 69 requires that pattern.
- **Risk**: Refresh fails permanently, so users must run `terraform state rm` by hand. This also misses the stated requirement to remove state when the policy disappears.
- **Suggested Fix**: Pass `fabcore.ErrCommon.EntityNotFound` in `getImmutabilityPolicy`. In Read, handle not-found with `utils.IsErrNotFound` and `resp.State.RemoveResource(ctx)`, as `workspacencp` does. Add a unit test where the GET returns 404 and assert the state is removed.

### Issue 3: Apply-time retention-reduction guard is untested
- **Severity**: Low
- **File**: internal/services/workspaceonelakeimmutability/resource_workspace_onelake_immutability_policy.go; internal/services/workspaceonelakeimmutability/resource_workspace_onelake_immutability_policy_test.go; internal/services/workspaceonelakeimmutability/fake_test.go
- **Line(s)**: resource modify 283-287 and validateRetentionIncrease 327-333; test RetentionCannotDecrease 119-153; fake modify 72-77
- **Description**: `ModifyPlan` (105-106) rejects reductions at plan time, so the apply-time check in `modify` and the fake's 400 branch are unreachable from the current suite. Nothing covers remote retention higher than config, which is the drift case the apply-time check guards.
- **Risk**: A regression in the apply-time guard would go unnoticed. A reduction could reach the API and surface only as an opaque service error.
- **Suggested Fix**: Add a unit test that calls `Update` directly, with the fake GET returning 30 and the plan requesting 7. Assert the error diagnostic and that no modify call is recorded.

### Issue 4: ReadRemovesMissingPolicy depends on an exact GET ordinal
- **Severity**: Low
- **File**: internal/services/workspaceonelakeimmutability/resource_workspace_onelake_immutability_policy_test.go; internal/services/workspaceonelakeimmutability/fake_test.go
- **Line(s)**: test 213-242 (ordinal at 214); fake GET 26-56 (comparison at 42)
- **Description**: `policyDisappearsAfterGet: 4` is compared with the fake's GET counter. It assumes how many refreshes the framework and plugin-testing harness issue, which the test does not control.
- **Risk**: The test can pass or fail for the wrong reason when refresh behavior changes. A failure may not show whether the Read removal path works.
- **Suggested Fix**: Replace the ordinal with an explicit trigger. Use a `PreConfig` hook that clears `store.policy` before the step that expects removal, then assert `ExpectNonEmptyPlan: true`.

### Issue 5: Acceptance test name omits `_CRUD`
- **Severity**: Low
- **File**: internal/services/workspaceonelakeimmutability/resource_workspace_onelake_immutability_policy_test.go
- **Line(s)**: 301
- **Description**: The acceptance test is named `TestAcc_WorkspaceOneLakeImmutabilityPolicyResource`. `testing-patterns.instructions.md` lines 19-20 require `TestAcc_<TypeName>Resource_CRUD`. The unit test name already follows the rule.
- **Risk**: Selectors that target the `_CRUD` suffix would miss this test, and the naming diverges from the documented convention.
- **Suggested Fix**: Rename to `TestAcc_WorkspaceOneLakeImmutabilityPolicyResource_CRUD`.

### Issue 6: Diagnostics setup has no read-back, and Invoke-FabricRest is silent on some success codes
- **Severity**: Low
- **File**: tools/scripts/Set-WellKnown.ps1
- **Line(s)**: 626-671 (POST at 670, no re-read); Invoke-FabricRest 173-189
- **Description**: After `modifyDiagnostics`, `Set-OneLakeDiagnostic` does not re-read settings to confirm the Lakehouse destination is active. `Invoke-FabricRest` returns nothing for statuses other than 200/201, or 202 with an operation id, so a 204 is treated as success without a check. `Set-OneLakeImmutabilityPolicy` does re-read (690-700).
- **Risk**: Diagnostics may not be routed while the script reports success. The failure then appears later as a `WorkspaceNotSupported` policy error, which is harder to diagnose.
- **Suggested Fix**: After the POST, re-GET `workspaces/{id}/onelake/settings` and call `Write-Log -Level 'ERROR'` if the status or destination Lakehouse ID does not match, mirroring the policy helper.

### Issue 7: validateRetentionIncrease name does not match its behavior
- **Severity**: Low
- **File**: internal/services/workspaceonelakeimmutability/resource_workspace_onelake_immutability_policy.go
- **Line(s)**: 327-333
- **Description**: The function rejects a retention value below the current remote value and allows equal values. The name suggests it validates an increase, which reads as rejecting unchanged values.
- **Risk**: A maintainer could "fix" it to reject equal values, breaking reruns that keep the same retention.
- **Suggested Fix**: Rename it to something like `validateRetentionNotReduced`, keeping the behavior.

## Resolution Log
_Updated by the driving agent as findings are addressed._

### Issue 1
- **Status**: Fixed
- **What changed**: Persisted the current well-known fixture before the immutability request and added a checked acceptance-test helper that skips when the policy key is not available.
- **Why**: OneLake diagnostic events can take time to propagate, so a `WorkspaceNotSupported` response must not prevent unrelated well-known resources from being saved or cause acceptance tests to panic.
- **How verified**: PowerShell parser passed; targeted Go tests passed. The Fabric setup script was not run because it can apply a non-reversible retention policy.

### Issue 2
- **Status**: Fixed
- **What changed**: `GetSettings` read diagnostics now recognize Fabric `EntityNotFound`, and resource `Read` removes state using the provider's standard `IsErrNotFound` pattern.
- **Why**: A workspace deleted outside Terraform should be dropped from state instead of leaving refresh stuck on a generic 404 error.
- **How verified**: Added a fake 404 refresh test; the focused unit suite passed.

### Issue 3
- **Status**: Fixed
- **What changed**: Added an apply-time drift test that raises the remote retention to 30 days after planning an increase to 8 days.
- **Why**: The apply-time guard must stop an update if another actor raises retention after Terraform's plan-time refresh.
- **How verified**: The update test expects the clear reduction diagnostic; all 16 focused unit tests passed.

### Issue 4
- **Status**: Fixed
- **What changed**: Removed the fake's GET ordinal behavior. The state-removal test now explicitly clears the fake policy in `PreConfig`.
- **Why**: The test controls the external deletion event instead of depending on Terraform's internal GET count.
- **How verified**: The focused unit test and suite passed.

### Issue 5
- **Status**: Fixed
- **What changed**: Renamed the resource acceptance test to `TestAcc_WorkspaceOneLakeImmutabilityPolicyResource_CRUD`.
- **Why**: Match repository test selectors and naming conventions.
- **How verified**: Package tests compiled and passed.

### Issue 6
- **Status**: Fixed
- **What changed**: `Set-OneLakeDiagnostic` now reads settings after the modify call and errors unless diagnostics are enabled with the exact dedicated Lakehouse destination.
- **Why**: The setup should catch a failed or ineffective diagnostics update at its source instead of surfacing later as an immutability precondition failure.
- **How verified**: PowerShell parser passed. The API documents 200/202 responses; the existing wrapper handles these, while the read-back verifies the applied state. No live Fabric setup was run.

### Issue 7
- **Status**: Fixed
- **What changed**: Renamed `validateRetentionIncrease` to `validateRetentionNotReduced`.
- **Why**: The helper permits an unchanged value and rejects only reductions; the new name describes that invariant.
- **How verified**: Focused unit tests and Go lint passed.

## Follow-up Verification

- Read `base_test.go` and `configure_test.go`; the acceptance helper validates well-known values and the configure tests cover nil and invalid provider data.
- `task testunit -- WorkspaceOneLakeImmutabilityPolicy ./internal/services/workspaceonelakeimmutability/`: passed, 16 tests, 83.0% coverage.
- `task build`: passed.
- `golangci-lint run ./internal/services/workspaceonelakeimmutability/...`: passed, 0 issues.
- PowerShell parser check for `tools/scripts/Set-WellKnown.ps1`: passed.
- `changie batch patch --dry-run --prerelease dev` with `CHANGIE_GITHUB_REPOSITORY=microsoft/terraform-provider-fabric`: passed.
- Acceptance tests and live `Set-WellKnown.ps1` provisioning were not run because they require Fabric infrastructure and include a non-reversible policy.
- No PR comments or remote artifact posts were made.

## Driver Decision

{
  "tier": "Low",
  "source": "Explicit user tier request; exact reviewer override: claude-haiku-5.5",
  "changeSize": "Moderate: 24 status entries, 11 new Go files / 1162 Go lines plus provider registration, WellKnown setup, examples, docs, and changelogs.",
  "cyclomaticComplexity": "Medium-high: remote workspace-scoped setting with irreversible retention, create/update/read/import/no-op-delete lifecycle, API errors, and well-known diagnostic setup dependencies.",
  "rationale": "Honor the user-selected Low tier despite material state and API complexity. Run one review-code CLI harness review, not subagents; the exact Claude Haiku 5.5 model was explicitly named by the user. GPT-6 Luna is excluded from all direct review use. User requests no PR comments or remote artifacts; only the skill-owned local review evidence is generated.",
  "primaryModel": "gpt-6-luna (OpenAI; host context, excluded as reviewer)",
  "reviewers": [
    {
      "id": "claude-haiku-5.5",
      "family": "Anthropic",
      "reasoningLevels": [
        "low",
        "medium",
        "high",
        "xhigh",
        "max"
      ],
      "selectionSource": "Exact user-named model; current functions.task catalog and GitHub Copilot CLI docs.",
      "selectedReasoning": "max"
    }
  ],
  "capabilitySources": [
    "functions.task current model catalog for exact ID and reasoning levels",
    "https://docs.github.com/en/copilot/how-tos/copilot-cli/use-copilot-cli/choose-a-model",
    "https://docs.github.com/copilot/reference/ai-models/supported-models",
    "https://docs.github.com/copilot/reference/ai-models/model-hosting",
    "https://developers.openai.com/api/docs/guides/model-selection"
  ],
  "zdrEvidence": "Copilot CLI help confirms GH_TOKEN is a supported authentication source; COPILOT_GITHUB_TOKEN is unset and GH_TOKEN is set. gh auth status confirms the active github.com EMU identity and copilot scope. Claude Haiku 5.5 is GA; GitHub documents ZDR coverage for GA Anthropic features in Copilot. The review runner uses standard CLI tools view, rg, glob without experimental features or Messages API tool search.",
  "reviewArtifact": {
    "attempt": 1,
    "tier": "Low",
    "outputDir": "Features\\1167\\reviews",
    "taskId": "1167"
  },
  "validation": {
    "providerLint": "tfproviderlintx returned an internal type-loading error for package crypto; not treated as a code finding",
    "terraformFormat": "terraform fmt -check -recursive examples\\resources\\fabric_workspace_onelake_immutability_policy examples\\data-sources\\fabric_workspace_onelake_immutability_policy: passed",
    "changelog": "changie batch patch --dry-run --prerelease dev with CHANGIE_GITHUB_REPOSITORY=microsoft/terraform-provider-fabric: passed",
    "build": "task build: passed",
    "unit": "task testunit -- WorkspaceOneLakeImmutabilityPolicy ./internal/services/workspaceonelakeimmutability/: passed, 14 tests, 81.6% coverage",
    "docs": "task docs: passed",
    "goLint": "golangci-lint run ./internal/services/workspaceonelakeimmutability/...: passed, 0 issues"
  },
  "recovery": "One retry only for transient transport/service failure. Use at most one xAI/Grok fallback only when a demonstrated provider/model failure qualifies under review-code policy; do not retry findings or malformed evidence."
}