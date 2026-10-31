---
producer: "review-code"
schema_version: 1
tier: "Medium"
contract: "direct"
execution_backend: "copilot-cli"
attempt: 5
stage_attempt: 4
round: 2
timestamp: "2026-10-31T07:32:02Z"
outcome: "issues_found"
reviewer_model: "grok-4.7"
input_sha256: "1b165472424a4c3a3a299c2ebdf5eec391bffd1df32086ff7128b3a0e8c72c0a"
---

## Review Summary

- **Round**: 2
- **Theme**: Architecture & patterns
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Artifact**: reviews/task-1119-attempt-5-stage-attempt-4-review-2-grok-4.7.md
- **Issues Found**: 4
- **Verdict**: ISSUES_FOUND

## Evidence Checklist

- [x] Compared `internal/services/mlexperiment/folder_fixtures_test.go` and `internal/services/mlmodel/folder_fixtures_test.go` with `internal/testhelp/testacc.go` (`FolderResource`, `RefByFQN`) and `.github/instructions/fake-handler-patterns.instructions.md`.
- [x] Read `.github/workflows/test.yml` and `.github/workflows/cleanup-script.yml` for the shared group, `cancel-in-progress`, matrix, path filter, `environment: development`, and `timeout-minutes`.
- [x] Read `DEVELOPER.md` acceptance-test section against the required lock, approval, fork, and pending-cleanup description.
- [x] Read `internal/testhelp/fakes/fake_typedhandler.go` (`Elements`, `asFabricItem`) and `context-only-public-dependencies/yaml.v3@v3.0.1/yaml.go` (`Kind` has no null kind; omitted bools stay false).
- [x] Confirmed `folder_id` has no `RequiresReplace` in `internal/pkg/fabricitem/resource_schema.go` and moves go through `fabricItemCheckMove`. Matrix include list is unchanged.
- [ ] Tests were not executed. Supplied validation was inspected only as evidence, not re-run.

## Issues

### Issue 1: Folder-fixture harness is copied instead of shared

- **Severity**: Medium
- **File**: internal/services/mlexperiment/folder_fixtures_test.go
- **Line(s)**: 21-126
- **Description**: The ML experiment and ML model fixture files are the same harness twice. Each owns an HCL walker, a mutex transport, a second-server splice, and a lifecycle `resource.Test`. Shared helpers already live in `internal/testhelp` (`FolderResource` at `internal/testhelp/testacc.go:162`, `RefByFQN`). Acceptance steps in `experimentCRUDSteps` (`internal/services/mlexperiment/resource_ml_experiment_test.go:220`) and `modelCRUDSteps` (`internal/services/mlmodel/resource_ml_model_test.go:246`) correctly use those helpers; the regression checks do not. `t.Setenv` also forces `resource.Test` instead of `resource.ParallelTest`, with no comment, and the stale `FABRIC_TESTACC_WELLKNOWN` value is never read.
- **Risk**: The next owned-folder test will copy this again. A fix to traversal checks, transport locking, or fake setup must be applied twice and can drift. That breaks the repo rule that test utilities stay in `internal/testhelp`.
- **Suggested Fix**: Put `AssertOwnedFolderCRUD` and the lifecycle runner next to `FolderResource`. Parameterize the item type. Keep the service files as thin callers of `experimentCRUDSteps` / `modelCRUDSteps`. Document why `resource.Test` is required. Do not spread this copy to other resources.

### Issue 2: Lifecycle tests bypass the fake-handler layer

- **Severity**: Medium
- **File**: internal/services/mlexperiment/folder_fixtures_test.go
- **Line(s)**: 24-32
- **Description**: Both lifecycle tests build a full `fakes.NewFakeServer()`, then replace `ServerFactory.Core.FoldersServer` with the folder server from a second `NewFakeServer()`. They then discard the factories from `NewTestUnitCase` (`internal/testhelp/testunit.go:85`) and install a private transport. The comment is right: item `Elements` (`internal/testhelp/fakes/fake_typedhandler.go:152-163`) sends every non-item through `asFabricItem`, which dereferences `Type` (`fake_typedhandler.go:279`). Fake patterns say that wiring belongs in `internal/testhelp/fakes`, not in service tests. The same block is repeated in `internal/services/mlmodel/folder_fixtures_test.go:24-32`.
- **Risk**: Callers depend on an unexported storage collision and on overwriting `ProtoV6ProviderFactories`. A later change to `typedHandler` can panic again, or the two copies can diverge. The mutex policy is also local, so the next multi-resource unit test rediscovers the unsynchronized `elements` slice.
- **Suggested Fix**: Add a fakes constructor that installs folder handlers on a store item list/delete will not convert. Serialize at that boundary if the in-memory store is shared. Extend the unit-test helper to accept a transport so tests do not replace factories after construction.

### Issue 3: Workflow guard model is narrower than the lock invariant

- **Severity**: Low
- **File**: internal/testhelp/folder_fixtures_test.go
- **Line(s)**: 25-28, 69-71, 112-120
- **Description**: Workflow `cancel-in-progress` is a bool. In `gopkg.in/yaml.v3`, an omitted field stays false, so deleting the key still passes the check at lines 69-71. Job concurrency correctly uses `yaml.Node`, but `workflowJob` does not decode `environment`. Development jobs are a fixed list (`cleanup-workspace`, `test-auth-spn`, `test`) rather than `environment.name == development`. A new development job with its own concurrency group would not fail this guard. Standalone cleanup is checked; that part is sound.
- **Risk**: The structural test is the regression lock for this change. It misses an omitted `cancel-in-progress` and any future development job that overrides the workflow group. GitHub's default is false, so omission is not a behavior change today; the guard still does not lock the explicit wiring.
- **Suggested Fix**: Decode workflow concurrency as `yaml.Node` and require `cancel-in-progress` to be present and false. Decode `environment` and assert every `development` job has no job-level concurrency, while still requiring the current three jobs.

### Issue 4: Deferred lock limits are not stated

- **Severity**: Low
- **File**: DEVELOPER.md
- **Line(s)**: 357-370
- **Description**: The new section correctly describes whole-workflow lifetime, approval waits, upstream fork PRs, uncovered fork/local/manual paths, and pending Cleanup replacement. It does not say that job-timeout bounds and stable workspace identity are out of scope. `test-auth-spn` (`.github/workflows/test.yml:85`), `changes`, `checkdocs`, `test-summary`, and `check-test-matrix` set no `timeout-minutes`. Matrix legs still share `WorkspaceRS`. Other CRUD tests still use `FolderRS1` / `FolderRS2`. The stuck-run sentence covers manual cancel only.
- **Risk**: Operators can read the new lock as complete isolation. A hung development or docs job holds `fabric-testacc-environment` until the platform default or a manual cancel, and non-ML folder fixtures remain shared.
- **Suggested Fix**: Add a short advisory to `DEVELOPER.md` only. Do not add timeouts, a workflow timeout, or a new workspace. State that those bounds are deferred, name the jobs without `timeout-minutes`, and note that matrix legs still share `WorkspaceRS` and non-ML tests still use `FolderRS1` / `FolderRS2`.

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

### Issue 4

- **Status**: Open
- **What changed**: pending
- **Why**: pending
- **How verified**: pending

## Publication Note

The machine-local Artifact summary path was normalized and generated review Markdown formatting was applied. Original response and finalized receipt remain immutable outside this source tree;findings,verdict and input provenance are unchanged.
