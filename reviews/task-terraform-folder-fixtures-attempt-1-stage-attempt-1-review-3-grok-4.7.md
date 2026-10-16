## Review Summary
- **Round**: 3
- **Theme**: Edge cases & robustness
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Stage**: 3
- **Attempt**: 1
- **Artifact**: reviews/task-terraform-folder-fixtures-attempt-1-stage-attempt-1-review-3-grok-4.7.md
- **Issues Found**: 0
- **Verdict**: CLEAN

## Evidence Checklist
- [x] Workflow lock is a shared constant with cancellation disabled: `.github/workflows/test.yml:21-27` and `.github/workflows/cleanup-script.yml:12-14` both use `fabric-testacc-environment` and `cancel-in-progress: false`.
- [x] The acceptance matrix stays inside that workflow lock and is not given its own cancelling lock: `.github/workflows/test.yml:311-326` (`test` job, `environment: development`, `fail-fast: false`, no job `concurrency`). Commented MSI job is inactive.
- [x] Other repository workflows do not join that group. `development` is used only by the test auth job (`.github/workflows/test.yml:63`), the test matrix (`.github/workflows/test.yml:315`), and cleanup (`.github/workflows/cleanup-script.yml:26`). Pages and release workflows use different environment names.
- [x] Isolation regression reads both workflow files and rejects an empty group, an expression group, a group mismatch, `cancel-in-progress: true`, a missing `test` job, and a per-job lock: `internal/testhelp/folder_fixtures_test.go:16-52`.
- [x] Owned-folder configs keep both folders in create and move steps and bind `folder_id` with `RefByFQN`: `internal/services/mlexperiment/resource_ml_experiment_test.go:220-269` and `internal/services/mlmodel/resource_ml_model_test.go:246-294`. Helper emits a fresh address and display name: `internal/testhelp/testacc.go:162-179`. Reference helper: `internal/testhelp/utils.go:190-192`.
- [x] HCL regression requires a single `fabric_folder` id traversal, both folder addresses in every step, no cross-run address reuse, and a real move: `internal/services/mlexperiment/folder_fixtures_test.go:48-110` and the matching model test.
- [x] Quoted interpolation from `CompileConfig` (`%q` of `${...}`) is still an HCL traversal, so the parser check above is a real dependency check, not a string-literal check.
- [x] Fake lifecycle serializes concurrent Terraform graph walks (`Taskfile.yml:26-29` sets destroy/apply parallelism to 3) with one mutex around `Do`: `internal/services/mlexperiment/folder_fixtures_test.go:34-46` and the model copy. Folder entities are stored on a separate folder server (`folder_fixtures_test.go:23-31`) because the generic item map cannot hold them.
- [x] Item move updates `FolderID` on the item record (`internal/testhelp/fakes/fabric_item.go` `moveItemOperations.Update`). Provider update calls move only when folder ids differ and returns before writing state on error: `internal/pkg/fabricitem/resource_item.go` `Update` and `internal/pkg/fabricitem/models_resource_item.go` `fabricItemCheckMove`.
- [x] Description omitted from the create config is defaulted to empty string, so the create check for `""` matches the schema: `internal/pkg/fabricitem/resource_schema.go` `description` `stringdefault.StaticString("")`.
- [x] Cleanup already deletes non-preserved folders after items, and preserves only `FolderRS1`/`FolderRS2`: `tools/scripts/Remove-WorkspaceRSItems.ps1` (preserve around line 372; folder delete at line 583). Owned test folders are therefore in the cleanup set if Terraform teardown does not run.
- [x] Operator text matches the lock: `DEVELOPER.md:353-368` states workflow-level sharing, no automatic cancel, one pending replacement, and that local/fork/manual cleanup is outside the lock.
- [x] Promoted direct modules `github.com/hashicorp/hcl/v2 v2.24.0` and `gopkg.in/yaml.v3 v3.0.1` already have sums in `go.sum`. No new version was introduced.
- [ ] Live acceptance, cloud cleanup, and `task testunit` were not re-executed in this review. See Validation.

## Assessed failure modes

No introduced defect met the bar. Checked paths:

1. **Shared-environment overlap.** Matching workflow-level groups are repository-scoped, so a running Test workflow and Cleanup workflow cannot overlap. `cancel-in-progress: false` does not cancel the running workflow. Internal matrix legs stay concurrent, which is the stated design; a per-job group would be the cancelling behavior the new test rejects.
2. **Move deleting the previous folder.** Both folder blocks are passed to both steps. The move changes only the item reference. The previous folder stays in config, so apply does not destroy it. The HCL test fails if either address disappears or the target does not change.
3. **Destroy order.** `folder_id = "${fabric_folder.<name>.id}"` is a dependency. Terraform destroys the item before the folder it still references. After a successful move, the old folder has no dependent and may be destroyed in parallel with the item under parallelism 3. That is safe only because the item no longer references it. A failed move leaves state pointing at the old folder, so that folder remains a dependency.
4. **Fake-store races.** SDK `ServerFactoryTransport` locks only lazy transport init, then calls the handler outside that lock. The new wrapper holds one mutex for the whole `Do`, covering concurrent folder creates and the later item create/move against the unsynchronized in-memory slices. Handlers upsert before returning, so the mutation is inside the locked call. This wrapper is test-only.
5. **Split folder server.** Item and folder fakes do not share a slice. That avoids the item converter treating a folder as an item. It also means the unit lifecycle does not prove the folder id exists in the folder store. Terraform still will not apply if the folder resource fails, and the HCL test rejects a static id. This is a stated test-double limit, not a product behavior change.
6. **Empty, stale, and hostile inputs.** Acceptance still requires `WorkspaceRS` and panics on a missing map the same way as before. Unit tests install a stale folder-only well-known document with `t.Setenv` and do not call `t.Parallel`, so the override does not leak into the next test. `FolderResource` does not sanitize the workspace id; `CompileConfig` quotes it, and the folder schema is a UUID, so a bad id fails closed instead of becoming a static folder dependency. Names are 20 letters, within the 255-character folder limit and the 123/256 item limits.
7. **Partial provider update.** A successful rename followed by a failed move returns before `State.Set`. That drift is the existing update function. The previous acceptance step already changed name and folder together. This change does not alter retry or classification. `MoveItem` still has no retry; create retry remains limited to the existing display-name-not-ready case.
8. **Dependency promotion.** Moving `hcl/v2` and `yaml.v3` to direct requirements does not change versions or runtime provider behavior.

## Documented residual risks (not findings)

These are real operational limits. They are already stated in `DEVELOPER.md:357-368` or are the same class of interrupt, and this diff does not claim to close them.

- **Pending replacement is not a queue of every dispatched run.** Current GitHub concurrency defaults to one pending run. A newer run cancels the pending run even when `cancel-in-progress` is false. It does not cancel the run that is already in progress. The new text matches that behavior. GitHub's note that runs which remain waiting are ordered by wait start does not make the default single-pending slot a FIFO of all requested runs.
- **Manual cancel and job timeout can still kill Terraform before destroy.** The test job and its test step are limited to 100 minutes (`.github/workflows/test.yml:318` and `:400`). `fail-fast: false` does not protect a timed-out or manually cancelled leg. The workflow lock cannot finish teardown that the runner already killed. Cleanup is the backstop, and it is blocked until that workflow run ends.
- **The lock is not a workspace mutex for every client.** Local commands, forks, other repositories, and manual cleanup do not take the GitHub group. In-workflow jobs that use `development` still run together. That includes the auth matrix and the four CLI legs.
- **Workflow-level scope also covers non-mutating jobs** such as docs checks. That is the cost of holding one group for the whole run so matrix legs are not cancelled. It can delay cleanup; it is not a second cancellation bug.
- **Environment deployment policy is outside the repository.** This review confirmed only the workflow files. It cannot see whether the `development` environment has a separate deployment concurrency rule.

## Validation

Reported validation, not re-executed here: red HCL and workflow regressions, fake CRUD lifecycle, five new tests passing, both ML unit packages passing, tidy diff clean, and actionlint on both workflows. No live acceptance was claimed.

This review did not re-run `task testunit`. The isolated Go, Task, and gotestsum binaries are present, but Terraform is not on PATH, and the Task entrypoint invokes bash. A lifecycle re-run would not have been a valid confirmation, and raw `go test` is prohibited. Source was not modified. No live acceptance, cloud call, or cleanup was run.

Independent confirmation is therefore static: the workflow YAML, HCL construction, fake transport, provider move/delete order, cleanup scope, and module sums cited above. Equality at the documented thresholds is not relevant to this fixture change. Unverified items are the reported test execution, rendering of the workflows by GitHub, and any environment-level deployment concurrency not stored in the repository.
