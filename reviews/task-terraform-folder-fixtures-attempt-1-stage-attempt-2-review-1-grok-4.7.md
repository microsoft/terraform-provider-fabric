## Review Summary
- **Round**: 1
- **Theme**: Broad sweep
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Artifact**: reviews/task-terraform-folder-fixtures-attempt-1-stage-attempt-2-review-1-grok-4.7.md
- **Issues Found**: 0
- **Verdict**: CLEAN

## Evidence Checklist
- [x] Working tree inspected against the public fixture-fix repository. Changed source is limited to `.github/workflows/test.yml`, `DEVELOPER.md`, `go.mod`, `internal/services/mlexperiment/resource_ml_experiment_test.go`, `internal/services/mlmodel/resource_ml_model_test.go`, and the three new `folder_fixtures_test.go` files. No Fabric API, retry, or classification product files are in the diff.
- [x] CRUD steps in `internal/services/mlexperiment/resource_ml_experiment_test.go` and `internal/services/mlmodel/resource_ml_model_test.go` call existing `testhelp.FolderResource` twice, join both folder configs into both steps, and set `folder_id` with `testhelp.RefByFQN`. They no longer read `FolderRS1` or `FolderRS2`. Acceptance entry points still take `WorkspaceRS` from `testhelp.WellKnown`.
- [x] `internal/testhelp/testacc.go` `FolderResource` emits a `fabric_folder` with `display_name` and `workspace_id`. `internal/testhelp/utils.go` `RefByFQN` builds `${<fqn>.id}`. Dependency `CompileConfig` quotes string values, so the emitted `folder_id` is a Terraform interpolation and a graph dependency, not a static fixture ID. `JoinConfigs` only concatenates the three configs.
- [x] HCL regressions `TestUnit_FolderFixtures_ExperimentCRUD` and `TestUnit_FolderFixtures_ModelCRUD` parse each step, require exactly two `fabric_folder` blocks, require the item `folder_id` traversal to be `fabric_folder.<name>.id`, require both folders to remain on the move step, require the target to change, and require a second run to allocate new folder names. Stale `FABRIC_TESTACC_WELLKNOWN` folder IDs are set and are not embedded.
- [x] Lifecycle tests run those same steps through `resource.Test`. `folder_id` in `internal/pkg/fabricitem/resource_schema.go` has no `RequiresReplace`. `ResourceFabricItem.Update` calls `MoveItem` when the planned folder differs, so the second step is an in-place move plus attribute update, followed by the harness destroy of the final config. Both folders remain in that final config, so teardown includes the owned folders and the item.
- [x] Separate folder store is justified by `internal/testhelp/fakes/fake_typedhandler.go`: item handlers convert non-item elements through `asFabricItem`, which folders cannot satisfy. Replacing `ServerFactory.Core.FoldersServer` with a second fake server's folder server happens before transport construction. `ProtoV6ProviderFactories` is replaced after `NewTestUnitCase`, so the mutex-wrapped transport is the one the provider uses. The mutex serializes concurrent folder creates on the unsynchronized in-memory fake.
- [x] Workflow lock: `.github/workflows/test.yml` and `.github/workflows/cleanup-script.yml` both use workflow-level group `fabric-testacc-environment` and `cancel-in-progress: false`. `test.yml` has one `concurrency` key. Job `test` keeps `strategy.fail-fast: false` and has no job-level `concurrency`. Other workflows use different groups. `internal/testhelp/folder_fixtures_test.go` asserts the shared constant group, explicit non-cancellation, the matrix job, and absence of a per-job lock.
- [x] `DEVELOPER.md` matches `internal/testhelp/wellknown.go` and `tools/scripts/Remove-WorkspaceRSItems.ps1`: a populated `FABRIC_TESTACC_WELLKNOWN` value overrides the fixture file, GitHub workflows inject the Actions variable, and a local setup does not publish that variable. The same section documents non-FIFO pending replacement and that local, fork, other-repository, and manual cleanup are outside the GitHub lock.
- [x] `go.mod` promotes already-used `github.com/hashicorp/hcl/v2 v2.24.0` and `gopkg.in/yaml.v3 v3.0.1` from indirect to direct. `go.sum` already contains both module and hash lines and is not modified. Imports in the new tests match those modules. No product dependency version changed.
- [x] Caller-reported validation accepted as consistent with the source: five new tests pass, each ML unit package passes seven tests, actionlint passed both workflows, and `go mod tidy -diff` was clean. Count of unit tests in each ML package matches that report.
- [ ] Live acceptance and cloud execution not performed in this slot; explicitly out of scope.
- [ ] Unit tests and actionlint were not re-executed in this slot because no Go toolchain was available. Source inspection was used instead of a fresh run.

## Driver stage gate

The paired GPT artifact is incomplete. This complete clean slot alone does not
complete High stage 1; the driver has not advanced to later stages.
