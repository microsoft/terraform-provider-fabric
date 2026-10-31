---
producer: "review-code"
schema_version: 1
tier: "Medium"
contract: "direct"
execution_backend: "copilot-cli"
attempt: 5
stage_attempt: 4
round: 6
timestamp: "2026-10-31T07:29:37Z"
outcome: "issues_found"
reviewer_model: "grok-4.7"
input_sha256: "9a3a6b620db1aa53f7a5a836a249868549207118d4bdaa3727d86e84e74ad6ea"
---

## Review Summary

- **Round**: 6
- **Theme**: Polish & hardening
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Artifact**: reviews/task-1119-attempt-5-stage-attempt-4-review-6-grok-4.7.md
- **Issues Found**: 3
- **Verdict**: ISSUES_FOUND

## Evidence Checklist

- [x] Read `.github/workflows/test.yml` concurrency (lines 23–27), `src` filter (lines 51–55), recoverable pre-clean (`continue-on-error` line 64, `-RecoverableItemsOnly` line 83), matrix (lines 354–366), `checkdocs` with no `if` (line 295), and timeouts (lines 68, 347).
- [x] Read `.github/workflows/cleanup-script.yml` schedule (line 7), comment placement (lines 8–12), shared group (lines 12–14), and 30-minute timeout (line 29).
- [x] Compared `DEVELOPER.md` lines 353–370 with workflow behavior, `internal/testhelp/wellknown.go` precedence, and `internal/testhelp/folder_fixtures_test.go` structural checks.
- [x] Confirmed yaml.v3 string decode uses the raw scalar (`context-only-public-dependencies/yaml.v3@v3.0.1/decode.go`, `reflect.String`), so the guard can read `with` values that are booleans in YAML.
- [x] Confirmed ML CRUD steps no longer read `FolderRS1`/`FolderRS2`, while other acceptance tests still do (31 `*_test.go` files; example `internal/services/activator/resource_activator_test.go` line 248).
- [ ] Tests and GitHub Actions were not executed in this pass. Supplied validation was treated as prior evidence only.

## Issues

### Issue 1: Fork pull-request concurrency policy is documented as guaranteed and non-cancelling

- **Severity**: High
- **File**: DEVELOPER.md
- **Line(s)**: 360–367
- **Description**: The new text says active runs are not cancelled and that upstream fork Test runs join the group and may wait for environment approval. `pull_request` uses the merge-commit workflow file. A fork joins `fabric-testacc-environment` only if that file still declares it. If that file sets `cancel-in-progress: true`, the incoming run can cancel the active shared-group run. Workflow concurrency is applied when the run is queued. `environment: development` is only on `cleanup-workspace`, `test-auth-spn`, and `test`. `checkdocs` (`.github/workflows/test.yml` line 295) has no environment gate, so approval is not a concurrency boundary.
- **Risk**: The change exists to stop cancelled or overlapping cleanup. A fork-controlled flag can cancel in-progress tests or full cleanup and leave WorkspaceRS mid-mutation. The current wording gives the opposite assurance.
- **Suggested Fix**: Document that fork participation and `cancel-in-progress` come from the merge-commit file, and that "not cancelled" applies only to runs that keep `cancel-in-progress: false`. Note that environment approval does not hold the slot closed. Keep the shared group. Add an advisory to require outside-collaborator workflow approval and to re-dispatch Cleanup after a cancelled run. Do not change the group or the trusted `cancel-in-progress: false` setting.

### Issue 2: Deferred timeout and shared-identity limits are not visible, and lock scope is easy to misread

- **Severity**: Medium
- **File**: DEVELOPER.md
- **Line(s)**: 357–370
- **Description**: The lock paragraph does not keep the deferred advisories visible. It also narrows replacement to "a newer commit." The group is the static repository-wide name `fabric-testacc-environment`, not the old per-ref expression. Every Test run holds it until the workflow exits, including docs-only runs, because `checkdocs` has no path `if` and concurrency is workflow-level. Test and Cleanup both schedule `0 2 * * *`. GitHub keeps one pending run, so nightly cleanup sits behind nightly Test and is replaced, not queued, if another run arrives. A superseded Cleanup run is only a cancelled Actions run; the winning run does not log it. Job timeouts stay 30 minutes for both cleanup jobs and 100 minutes for the test matrix. A timeout can stop mutation and then release the slot. Only MLModel and MLExperiment CRUD stopped using `FolderRS1`/`FolderRS2`. Other acceptance tests still share those fixtures, and the four matrix legs still run in parallel inside one lock.
- **Risk**: Operators can treat the lock as complete isolation, miss a dropped nightly cleanup, or assume folder identity races are gone. Timeout and stable-identity work was supposed to stay advisory, not disappear.
- **Suggested Fix**: Add a short advisory only. State that the slot is repository-wide and survives path-filter skips; that identical crons make pending Cleanup easy to replace; that supersession is visible only as a cancelled Cleanup run; that 30/100-minute timeouts are unchanged and can strand items; and that WorkspaceRS plus non-ML `FolderRS1`/`FolderRS2` use remains shared. Do not change timeouts, matrix, or fixture scope.

### Issue 3: Cleanup and fixture test names do not match the behavior they describe

- **Severity**: Low
- **File**: DEVELOPER.md
- **Line(s)**: 359
- **Description**: "Teardown" is not a job. The in-workflow `cleanup-workspace` job is a best-effort recoverable pre-clean (`continue-on-error: true`, `-RecoverableItemsOnly`). Standalone cleanup deletes active items and non-preserved folders. Terraform destroy happens inside the test job. In `.github/workflows/cleanup-script.yml` lines 8–11, the concurrency warning is indented under `workflow_dispatch` instead of under `concurrency:` as in `test.yml` lines 23–25. `TestUnit_FolderFixtures_ExperimentCRUD` and `TestUnit_FolderFixtures_ModelCRUD` only parse HCL, while `*Lifecycle` runs `resource.Test`. That inverts `.github/instructions/testing-patterns.instructions.md` line 12, where `_CRUD` means full lifecycle. The `FABRIC_TESTACC_WELLKNOWN` `t.Setenv` calls are not read by those tests, so they look like a stale-fixture assertion but cannot fail the test.
- **Risk**: Operators look for a post-test teardown job or treat pre-clean failure as a hard gate. Reviewers can mistake the static HCL test for lifecycle coverage.
- **Suggested Fix**: Name the three phases separately in `DEVELOPER.md`. Move the cleanup-script comments under `concurrency:`. Rename the static tests to a config/HCL suffix, or rename `*Lifecycle` to `_CRUD`. Drop or actually assert the unused well-known env tripwire.

## Resolution Log

_Updated by the driving agent as findings are addressed._

### Issue 1

- **Status**: Open
- **What changed**: pending
- **Why**: pending
- **How verified**: pending

### Issue 2

- **Status**: Open
- **What changed**: pending
- **Why**: pending
- **How verified**: pending

### Issue 3

- **Status**: Open
- **What changed**: pending
- **Why**: pending
- **How verified**: pending

## Publication Note

The machine-local Artifact summary path was normalized and generated review Markdown formatting was applied. Original response and finalized receipt remain immutable outside this source tree;findings,verdict and input provenance are unchanged.
