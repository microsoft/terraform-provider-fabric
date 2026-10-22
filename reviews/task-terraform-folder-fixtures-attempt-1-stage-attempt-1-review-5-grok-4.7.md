## Review Summary

- **Round**: 5
- **Stage**: 5
- **Attempt**: 1
- **Theme**: Testing & coverage
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Reasoning**: xhigh
- **Artifact**: reviews/task-terraform-folder-fixtures-attempt-1-stage-attempt-1-review-5-grok-4.7.md
- **Issues Found**: 0
- **Verdict**: CLEAN

## Evidence Checklist

- [x] Acceptance and unit lifecycle share one builder. `TestAcc_MLExperimentResource_CRUD` (`internal/services/mlexperiment/resource_ml_experiment_test.go:214`) and `TestUnit_FolderFixtures_ExperimentLifecycle` (`internal/services/mlexperiment/folder_fixtures_test.go:21`) both call `experimentCRUDSteps` (`internal/services/mlexperiment/resource_ml_experiment_test.go:220`). The model pair is the same shape: `internal/services/mlmodel/resource_ml_model_test.go:240`, `internal/services/mlmodel/folder_fixtures_test.go:21`, and `modelCRUDSteps` at `internal/services/mlmodel/resource_ml_model_test.go:246`.
- [x] The lifecycle path is real Terraform with SDK fakes, not a hand-rolled apply. `newTestUnitCase` sets `IsUnitTest` and a protocol-6 factory (`internal/testhelp/testunit.go:66-90`). The fixture test then swaps in a mutex-wrapped `ServerFactoryTransport` over that same factory (`internal/services/mlexperiment/folder_fixtures_test.go:24-32` and the model copy). Passing `resource.Test` executes create, refresh/read, move, refresh/read, and destroy.
- [x] HCL assertions walk the expression graph of that same builder. They require two `fabric_folder` blocks, the requested workspace literal, a length-3 `fabric_folder.<name>.id` traversal for `folder_id`, both folders on the move step, a changed target, and no name reuse across two calls (`internal/services/mlexperiment/folder_fixtures_test.go:48-111` and `internal/services/mlmodel/folder_fixtures_test.go:48-111`).
- [x] The split folder store is a fake limitation, not a skipped product path. Item listing converts non-matching elements when the entity has item fields (`internal/testhelp/fakes/fake_typedhandler.go:152-168`), and conversion dereferences `Type` (`internal/testhelp/fakes/fake_typedhandler.go:278-279`). The lifecycle test therefore installs a second server's `FoldersServer` (`folder_fixtures_test.go:23-24`) while item calls stay on the original factory. `MoveItem` still stores `TargetFolderID` (`internal/testhelp/fakes/fabric_item.go:171-172`). `TestCheckResourceAttrPair` binds item `folder_id` to the Terraform-managed folder `id`.
- [x] The workflow guard matches current YAML. `TestUnit_FolderFixtures_WorkflowIsolation` (`internal/testhelp/folder_fixtures_test.go:16-53`) requires one non-expression group, `cancel-in-progress` false on both workflows, a `test` job, and no job-level `concurrency` on that job. Source has that group and `cancel-in-progress: false` at `.github/workflows/test.yml:26-27` and `.github/workflows/cleanup-script.yml:13-14`. The matrix job is `test` (`.github/workflows/test.yml:310`). Neither of those two workflow files has a job-level `concurrency` key.
- [x] Re-ran the unit entrypoint with `TF_ACC=0`. Experiment package: 7 passed in 16.559s. Model package: 7 passed in 17.835s. Workflow guard: 1 passed in 3.172s. Did not run acceptance, cloud, cleanup, actionlint, or module tidy.
- [ ] GitHub concurrency scheduling, including one-pending replacement, was not executed. The guard is a YAML contract test.
- [ ] No mutation or branch-coverage score. The runner's package statement percentages are not evidence that every fake or workflow branch was exercised.

## Clean Assessment

No current source-backed regression is left untested. Extra tests that would only guard hypothetical future edits are not required.

Same builder. Acceptance no longer reads static folder fixtures. `experimentCRUDSteps` and `modelCRUDSteps` call `FolderResource` (`internal/testhelp/testacc.go:162`) twice, pass `RefByFQN(..., "id")`, and keep both folder blocks on the move step. The unit lifecycle calls those functions and runs them through Terraform. Step checks require the created and moved item `folder_id` to equal the owned folder `id`, so a stale literal cannot satisfy the pair check. Destroy is executed by `resource.Test`; these runs passed, so teardown of that configuration did not error.

HCL graph. The CRUD tests parse builder output and walk `Variables()`. That rejects a non-reference `folder_id`, a missing folder block, and a move that drops the previous folder, without depending on the fake. It is not `terraform graph` output, and it does not need to be: the lifecycle run is the dynamic dependency execution. The config helper quotes string values, including references. A quoted interpolation still yields a three-part traversal, and the passing lifecycle run shows Terraform accepted that form.

Fake limitations. The second folder server and the serializing transport are test-double constraints. They do not skip folder create, item create, or move in the provider. The item fake does not check that the target folder exists in the folder server. That does not hide the fixture bug this change guards, because the pair check still requires the provider to persist the interpolated folder id. An explicit post-destroy API read-back, or a test that the mutex is required, would be additional coverage rather than a fix for a failing contract. `CheckDestroy` in `newTestUnitCase` inspects pre-destroy state for the item address; the harness destroy that follows is what these passing tests executed. Both folders remain in the configuration, so a folder destroy error fails the test.

Workflow guard. The new test locks the property that changed: one constant group shared with cleanup, cancellation disabled, and no per-job lock on the matrix job. It does not pin a coordinated rename of that shared constant, and it does not scan every job for a future concurrency override. Current cleanup has no job-level concurrency, and the matrix job is the leg called out by the workflow comment. Those extra scans would not catch a regression in the present diff. The test also does not prove the scheduler keeps only one pending run. That limit is documented and is not a unit-testable defect in this change.

## Tests Run Versus Reported

Reported by the review input and not re-executed here:

- Historical red HCL and workflow failures before the fix.
- Workflow lint on both workflow files.
- Module tidy diff clean.
- Live acceptance, cloud, and cleanup. Those stay out of scope.

Actually run, with `TF_ACC=0`, through the Taskfile unit task only. The acceptance task was not used, because it forces `TF_ACC=1`. Raw `go test` was not used.

- `task testunit -- . ./internal/services/mlexperiment` — 7 passed, 16.559s; runner reported 100.0% statement coverage of that package.
- `task testunit -- . ./internal/services/mlmodel` — 7 passed, 17.835s; runner reported 100.0% statement coverage of that package.
- `task testunit -- FolderFixtures ./internal/testhelp` — 1 passed, 3.172s; runner reported 1.1% statement coverage of `testhelp`.

The pattern `^TestUnit_.` covers every `TestUnit_` function in the two ML packages, including both new fixture tests in each package. The workflow command covers `TestUnit_FolderFixtures_WorkflowIsolation` only. Together those runs are the five new tests plus the existing ML unit tests (7 per package).

## Limitations

- Statement-coverage percentages are the unit runner's package figures, not branch or mutation scores. `testhelp` stays low because the guard only reads two workflow files.
- Fake item move does not consult the separate folder store. The pair check proves persistence of the Terraform-interpolated id, not folder-store foreign-key validation.
- Workflow YAML parsing does not execute the hosted concurrency scheduler or its one-pending replacement rule.
- This review did not use prior review artifacts. The artifact path is repository-relative so the public result does not embed a local absolute path.
