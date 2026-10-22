## Review Summary

- **Round**: 1
- **Theme**: Broad sweep
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Artifact**: reviews/task-terraform-folder-fixtures-attempt-1-stage-attempt-1-review-1-grok-4.7.md
- **Issues Found**: 0
- **Verdict**: CLEAN

## Evidence Checklist

- [x] Inspected the isolated public diff and untracked tests: `.github/workflows/test.yml`, `DEVELOPER.md`, `internal/services/mlmodel/resource_ml_model_test.go`, `internal/services/mlexperiment/resource_ml_experiment_test.go`, `internal/services/mlmodel/folder_fixtures_test.go`, `internal/services/mlexperiment/folder_fixtures_test.go`, and `internal/testhelp/folder_fixtures_test.go`. No other tracked or untracked source changes.
- [x] Shared cleanup exclusion: both `.github/workflows/test.yml` and `.github/workflows/cleanup-script.yml` use the workflow-level group `fabric-testacc-environment` with `cancel-in-progress: false`. The `test` matrix job has no job-level `concurrency` key. Other workflows use different groups, and only the cleanup workflow invokes the workspace cleanup script. A workflow-level group holds the lock for every matrix leg and its Terraform destroy before another run in that group can start. A per-job lock would be the wrong shape because GitHub keeps only one pending job per group.
- [x] Queued-run semantics match the stated contract and `DEVELOPER.md`: `cancel-in-progress: false` leaves the active run, including teardown, running. GitHub keeps at most one pending run and may replace it. That can skip a queued cleanup or required check. It does not cancel an in-progress matrix. The constant group also serializes non-acceptance jobs in those two workflows. That blast radius is the documented workflow-level lock, not a separate per-matrix lock.
- [x] Fake transport and folder store: `folderFixtureTransport.Do` uses a pointer mutex around the fake transport. `Taskfile.yml` sets Terraform parallelism to 3, and the in-memory fake slices are not synchronized, so the wrapper is necessary. Replacing `ServerFactory.Core.FoldersServer` with the folder server from a second fake server keeps folder entities out of the item slice. Item `Get` calls `asFabricItem` on every element, so a shared slice would be unsafe. The wrapped transport is built after that swap. Item move still updates `FolderID` through the original item handler.
- [x] HCL fixture ordering: both CRUD steps join both `FolderResource` configs with the item config. `RefByFQN` emits `${fabric_folder.<name>.id}`, and `CompileConfig` quotes that string. Quoted interpolation still creates a Terraform dependency, so the referenced folder is created or known before the item create or move. Attribute order inside `CompileConfig` is map-random and does not affect the resource graph. The second step keeps the first folder address, so the move does not drop it. `folder_id` is not `RequiresReplace`. The update path calls the move API when the folder changes. The generated-HCL tests check both folder addresses, the workspace literal, a single `fabric_folder` ID traversal, and that the second step changes the target.
- [x] Stale static folder values are not read by the new CRUD steps. Acceptance tests still require `WorkspaceRS` from well-known data. The unit tests set a synthetic well-known value that contains only stale folder IDs and then build configs without reading it. `t.Setenv` runs in non-parallel tests, and other well-known reads happen before `t.Parallel`, so that override does not leak into those reads.
- [x] Teardown remains in the normal Terraform destroy after the steps. The item depends on its current folder, so that folder is destroyed after the item. The cleanup script still exempts only the two static folder names. The new random folders are therefore protected by the workflow lock while the matrix is alive, and they are eligible for later cleanup if destroy leaks them.
- [ ] Did not re-run unit tests, actionlint, or live acceptance. The review harness has no Go toolchain here, and live cloud execution was out of scope. The conclusion is from source inspection of the changed paths and their helpers.

## Driver persistence

The complete inline reviewer result was saved unchanged; no reviewer source edits or live execution occurred. Paired-stage status is recorded in the GPT stage-one artifact.
