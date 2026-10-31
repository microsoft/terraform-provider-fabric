---
producer: "review-code"
schema_version: 1
tier: "Medium"
contract: "direct"
execution_backend: "copilot-cli"
attempt: 5
stage_attempt: 4
round: 5
timestamp: "2026-10-31T07:38:06Z"
outcome: "issues_found"
reviewer_model: "grok-4.7"
input_sha256: "cc69c4b660e6341b319b114a8057bdc5ad7f1a4732289a5af857339ce515cb5f"
---

## Review Summary

- **Round**: 5
- **Theme**: Testing & coverage
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Artifact**: reviews/task-1119-attempt-5-stage-attempt-4-review-5-grok-4.7.md
- **Issues Found**: 3
- **Verdict**: ISSUES_FOUND

## Evidence Checklist

- [x] Compared `internal/testhelp/folder_fixtures_test.go` with `.github/workflows/test.yml` and `.github/workflows/cleanup-script.yml` for group, `cancel-in-progress`, paths-filter id, `src` output, selector set, and job `if`/`needs`.
- [x] Checked yaml.v3 `scalar()` (`context-only-public-dependencies/yaml.v3@v3.0.1/decode.go`) so bool/int `with` values unmarshal into `map[string]string` and do not make the workflow guard fail closed.
- [x] Reviewed owned-folder tests in `internal/services/mlexperiment/folder_fixtures_test.go`, `internal/services/mlmodel/folder_fixtures_test.go`, and the shared `*CRUDSteps` checks against `testhelp.FolderResource`, fake folder/item handlers, and `MoveItem`.
- [x] Confirmed existing `TestUnit_MLExperimentResource_CRUD` and `TestUnit_MLModelResource_CRUD` keep the same `folder_id` on update, so the new lifecycle tests are the only move execution.
- [ ] Tests were not executed; this review is source inspection only.

## Issues

### Issue 1: Workflow guard accepts extra dependencies and ignores approval/redispatch wiring

- **Severity**: Medium
- **File**: internal/testhelp/folder_fixtures_test.go
- **Line(s)**: 113-143
- **Description**: The dependency check only requires `needs` to contain `changes` and, except for `cleanup-workspace`, `cleanup-workspace`. Extra needs still pass. The job struct also has no `continue-on-error` or `environment` field, and the workflow struct does not parse `on`. Removing `continue-on-error: true` from `.github/workflows/test.yml` line 64, removing `environment: development` from the development jobs or standalone cleanup, or removing `workflow_dispatch` from `.github/workflows/cleanup-script.yml` line 8 would leave the guard green. Development jobs are a hardcoded name list, so a new `environment: development` job, or a job-level `concurrency` block on `checkbuild`, `checkdocs`, `test-summary`, or `check-test-matrix`, is not rejected.
- **Risk**: The stated regression is condition and dependency wiring, plus documented approval waits and cleanup redispatch. Those contracts can regress while this guard still passes, allowing cleanup and tests to overlap or a superseded cleanup run to have no tested recovery trigger.
- **Suggested Fix**: Require exact `needs` sets (`cleanup-workspace`: `changes`; `test-auth-spn` and `test`: `changes` and `cleanup-workspace`). Assert `continue-on-error: true` on in-workflow cleanup, `environment.name: development` on every development job and standalone cleanup, and `workflow_dispatch` on the cleanup workflow. Discover development jobs by that environment marker, reject job-level `concurrency` on every job, and include actual versus expected values in failures.

### Issue 2: Move and stale-fixture isolation are not actually asserted

- **Severity**: Medium
- **File**: internal/services/mlexperiment/folder_fixtures_test.go
- **Line(s)**: 20-32, 49-125
- **Description**: `TestUnit_FolderFixtures_ExperimentLifecycle` and the matching model test only inherit `display_name`, `description`, and `TestCheckResourceAttrPair` for `folder_id`. They never record the item `id` on the create step or require it to be unchanged after the folder change. The same tests set `FABRIC_TESTACC_WELLKNOWN` to `stale-first` / `stale-second`, but `experimentCRUDSteps` / `modelCRUDSteps` and `FolderResource` never call `WellKnown()`, and neither test scans the config for `FolderRS1`, `FolderRS2`, or those stale IDs. The HCL walk also type-asserts `refs[0][1].(hcl.TraverseAttr)` and `refs[0][2]` without an `ok` check, so an unexpected traversal panics before later assertions. The same gaps are copied in `internal/services/mlmodel/folder_fixtures_test.go`.
- **Risk**: A `folder_id` change that replaces the item still satisfies `TestCheckResourceAttrPair`. The stale environment variable looks like an isolation tripwire but cannot fail the test. The only existing unit CRUD paths (`resource_ml_experiment_test.go` lines 183-207 and `resource_ml_model_test.go` lines 186-232) never change `folder_id`, so this hole is the move coverage.
- **Suggested Fix**: In the lifecycle steps, capture the item `id` and both folder addresses after create, then assert the same item `id` and both folders after the move. Assert the rendered config does not contain `stale-first`, `stale-second`, `FolderRS1`, or `FolderRS2`. Replace the traversal type assertions with comma-ok failures. Apply the same checks to both packages.

### Issue 3: Guard failures omit the actual wiring

- **Severity**: Low
- **File**: internal/testhelp/folder_fixtures_test.go
- **Line(s)**: 64-72, 113-136, 138-143
- **Description**: Group mismatch, `cancel-in-progress`, `if` mismatch, needs mismatch, and the standalone-cleanup check use `t.Fatal` without printing the decoded group, condition, needs, or node kind. Only the missing selector pattern is included in the message.
- **Risk**: A regression in the shared lock or job wiring fails with a generic sentence, which makes the intended diagnostic weaker than the missing-selector failure this guard is meant to preserve.
- **Suggested Fix**: Include actual and expected values in each fatal, including both concurrency groups, the job `if`, the decoded `needs`, and the concurrency node kind.

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
