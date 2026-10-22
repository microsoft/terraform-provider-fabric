## Review Summary

- **Round**: 4
- **Stage**: 4
- **Attempt**: 1
- **Theme**: Detailed correctness
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Artifact**: reviews/task-terraform-folder-fixtures-attempt-1-stage-attempt-1-review-4-grok-4.7.md
- **Issues Found**: 0
- **Verdict**: CLEAN

## Evidence Checklist

- [x] Diff inspected against the public checkout: workflow lock, developer notes, direct module requires, MLExperiment/MLModel acceptance step builders, and the three new unit files.
- [x] Folder data flow traced from `internal/testhelp/testacc.go` `FolderResource` (162-174) and `internal/testhelp/utils.go` `RefByFQN` (190-192) through `experimentCRUDSteps` (`internal/services/mlexperiment/resource_ml_experiment_test.go` 220-265) and `modelCRUDSteps` (`internal/services/mlmodel/resource_ml_model_test.go` 246-291). `JoinConfigs` keeps both folder blocks in both steps; `folder_id` is `${fabric_folder.<name>.id}`, which remains an HCL reference after quoted compilation.
- [x] Move versus replace checked in `internal/pkg/fabricitem/resource_schema.go` 204-209 (`folder_id` has no `RequiresReplace`), `resource_item.go` 222-225 (`MoveItem` on folder change), and `models_resource_item.go` 132-136. Experiment rename is allowed (`internal/services/mlexperiment/resource_ml_experiment.go` 16); the model step keeps `display_name` unchanged despite `NameRenameAllowed: false` (`internal/services/mlmodel/resource_ml_model.go` 16, model steps 264-282).
- [x] Fake lifetime checked: lifecycle tests replace only `FoldersServer` after a second `NewFakeServer` (`internal/services/mlexperiment/folder_fixtures_test.go` 21-33 and the model twin at 21-33), then wrap `Do` with a pointer mutex (35-46). Item move stays on the original item handler (`internal/testhelp/fakes/fabric_item.go` 156-172). `WellKnown` is not cached (`internal/testhelp/wellknown.go` 40-58) and the new step builders do not read `FolderRS1`/`FolderRS2`.
- [x] HCL indexing checked in `internal/services/mlexperiment/folder_fixtures_test.go` 48-110 (model twin 48-110): two runs, exactly two steps, traversal length 3, root `fabric_folder`, attribute `id`, both folder labels retained, and the second target differs. Alpha resource names from `RandomName` make `TraverseAttr` the actual shape, not an index traversal.
- [x] Workflow types checked: `test.yml` 23-27 and `cleanup-script.yml` 12-14 share the constant group `fabric-testacc-environment` with boolean `cancel-in-progress: false`. The only `concurrency:` key in `test.yml` is workflow-level; job `test` starts at line 310 and has no job concurrency key. `internal/testhelp/folder_fixtures_test.go` 16-52 unmarshals those fields and rejects an expression group or a job-level lock on `test`.
- [x] Docs checked against code: `DEVELOPER.md` 318-321 matches env-over-file precedence in `wellknown.go` and `vars.FABRIC_TESTACC_WELLKNOWN` at `test.yml` 31 and `cleanup-script.yml` 17. Setup writes the local fixture file and does not publish the Actions variable. Acceptance text at `DEVELOPER.md` 353-366 matches owned-folder steps and the non-FIFO pending-run limit.
- [x] `go.mod` 13 and 28 promote `github.com/hashicorp/hcl/v2 v2.24.0` and `gopkg.in/yaml.v3 v3.0.1` to direct requires. Readonly module resolution succeeded for the unit runs below.
- [x] Actually ran, with `TF_ACC=0` and `task testunit` only: `FolderFixtures` in `internal/testhelp` (1 pass), `internal/services/mlexperiment` (2 pass), and `internal/services/mlmodel` (2 pass); plus `MLExperimentResource` (3 pass) and `MLModelResource` (3 pass). No acceptance, cloud, or cleanup run.
- [ ] Full ML package suites, actionlint, and live acceptance were not re-run. Reported package statement coverage is not product-wide coverage. See limitations.

## Verified Assessment

1. Acceptance CRUD no longer reads static `FolderRS1`/`FolderRS2`. Each call to `FolderResource` allocates a new alpha resource name and display name, emits a literal `workspace_id`, and returns `fabric_folder.<name>`. Both generated blocks are concatenated into step 1 and step 2, so a folder-id change cannot drop the previous folder from configuration.
2. `RefByFQN` formats `${fabric_folder.<name>.id}`. Quoted compilation still yields one three-element traversal rooted at `fabric_folder` and ending at `id`. That is a real dependency, so Terraform orders folder create before item create and item destroy before folder destroy. `TestCheckResourceAttrPair` compares the item `folder_id` to the addressed folder `id` after apply, not to a baked-in fixture string.
3. The second experiment step changes display name, description, and folder. Rename is allowed, and `folder_id` is not force-new, so the provider update path calls item update and then `MoveItem`. The model step changes description and folder only, so its force-new display name is not triggered. Fake `MoveItem` assigns `FolderID` from `TargetFolderID` and does not delete folders.
4. Splitting `FoldersServer` onto a second in-memory server avoids item conversion of folder entities. The replacement happens before the transport is constructed, and the core dispatcher reads `ServerFactory.Core` on request, so it observes the replacement. The mutex is on a pointer receiver and is held for the full inner `Do`, which is the concurrent create/read path Terraform can exercise inside one apply.
5. The static HCL loops are not off-by-one: `range 2` is runs 0 and 1; step 0 records names; later steps require the same names, two folders, and a different target. An empty target fails `folders[target]`. A static folder id fails the traversal-length check.
6. Workflow concurrency is a shared root-level string group with boolean cancel false. That serializes test and cleanup workflow runs, including matrix teardown, without cancelling an active run. Job `test` has no nested concurrency key, so the workflow lock is not replaced by a per-matrix-job lock. The documented one-pending replacement is a platform limit, not an off-by-one in this diff.
7. No type or conversion defect was found in the changed steps, fake transport, YAML struct tags, or direct module promotion. Pre-existing helpers used by the new tests behave consistently with this data flow.

## Actually Run Versus Reported

Reported checks were not treated as coverage. This review ran only `task testunit` with `TF_ACC=0` and local module resolution:

- `FolderFixtures` / `internal/testhelp`: 1 passed.
- `FolderFixtures` / `internal/services/mlexperiment`: 2 passed.
- `FolderFixtures` / `internal/services/mlmodel`: 2 passed.
- `MLExperimentResource` / `internal/services/mlexperiment`: 3 passed.
- `MLModelResource` / `internal/services/mlmodel`: 3 passed.

Those runs compile the changed tests and execute the fake Terraform lifecycle for the new folder steps. They do not execute the acceptance tests, a live API, cleanup, or actionlint. The reported full-package count and actionlint result were not reproduced here.

## Limitations

- No live acceptance, cloud call, or cleanup was run. Acceptance ordering and destroy against a real workspace remain unverified.
- Statement-coverage percentages from these filters are package-local. They do not cover the shared item implementation, other services, or GitHub execution.
- Lifecycle checks do not assert that the item id is unchanged. Current schema and step attributes make the second step an in-place update plus move; the tests alone would not catch a future force-new regression.
- Plugin-testing calls `CheckDestroy` with the pre-destroy state. Destroy failure still fails the test; the helper does not itself prove the resource is absent afterward.
- The HCL test rejects a static `folder_id` expression. It does not search the whole configuration for the poisoned sample ids.
- The workflow test parses the two workflow files. It does not simulate GitHub's pending-run replacement. Matrix legs inside one run remain concurrent; that is pre-existing and outside this lock's job-serialization claim.
- Local, fork, and manual cleanup stay outside the GitHub group, as the developer note states.
