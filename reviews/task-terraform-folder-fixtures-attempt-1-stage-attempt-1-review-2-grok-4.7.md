## Review Summary
- **Round**: 2
- **Stage**: 2
- **Attempt**: 1
- **Theme**: Architecture & patterns
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Reasoning**: xhigh
- **Artifact**: reviews/task-terraform-folder-fixtures-attempt-1-stage-attempt-1-review-2-grok-4.7.md
- **Issues Found**: 0
- **Verdict**: CLEAN

## Evidence Checklist
- [x] Frozen diff reviewed for `.github/workflows/test.yml`, `DEVELOPER.md`, `go.mod`, both ML resource tests, both `folder_fixtures_test.go` files, and `internal/testhelp/folder_fixtures_test.go`. Checkout status matches that file set; no product runtime files are modified.
- [x] Existing helper seam inspected: `FolderResource` in `internal/testhelp/testacc.go` (line 162) and `RefByFQN` in `internal/testhelp/utils.go` (line 190). Both ML CRUD builders call those helpers plus `at.JoinConfigs`. No new folder-construction API was added.
- [x] Move-versus-replace contract checked in `internal/pkg/fabricitem/resource_schema.go` (`folder_id` has no `RequiresReplace`) and `internal/pkg/fabricitem/resource_item.go` (`Update` calls `MoveItem` only when `fabricItemCheckMove` is true). The two-step config change follows that existing item abstraction.
- [x] Fake-store seam checked in `internal/testhelp/fakes/fake_typedhandler.go` (`Elements` / `asFabricItem`) and `internal/testhelp/fakes/fabric_folder.go` (`configureFolder`). The lifecycle tests swap in a second server's `FoldersServer` and wrap transport; they do not edit fake or provider code.
- [x] Workflow boundary checked in `.github/workflows/test.yml` (concurrency at lines 23-27; `test` job matrix remains, no job-level `concurrency`) and `.github/workflows/cleanup-script.yml` (same static group, `cancel-in-progress: false`). In-repo `environment: development` users are only those two workflows. `ghpages.yml` and `release.yml` use different environments and groups.
- [x] Convention sources read: `.github/copilot-instructions.md`, `.github/instructions/testing-patterns.instructions.md`, `.github/instructions/coding-standards.instructions.md`, and `Taskfile.yml` test filters. New tests stay in `_test` packages, use the `TestUnit_FolderFixtures_` designator, and keep `TestAcc_*Resource_CRUD` names.
- [x] `DEVELOPER.md` acceptance text compared with `internal/testhelp/wellknown.go`: environment variable precedence, no `FolderRS1`/`FolderRS2` reads in the CRUD builders, and the non-FIFO pending-run limitation are consistent with the code.
- [x] Actually ran local static and unit checks listed under Tests. Did not read prior review artifacts.

## Clean Assessment

No architectural change is required. The diff keeps one configuration source of truth per resource, reuses the existing test-helper layer, and places the shared-environment lock at the only workflow scope that preserves matrix parallelism.

### 1. Fixture ownership and dependency direction

`experimentCRUDSteps` and `modelCRUDSteps` are the single builders for acceptance, HCL contract checks, and fake lifecycle execution. Each builder creates two folders with `testhelp.FolderResource`, keeps both HCL blocks in the create step and the move step, and passes `testhelp.RefByFQN(..., "id")` instead of a well-known folder id. `TestCheckResourceAttrPair` compares the item `folder_id` to the owned folder `id`.

That dependency direction is correct. Service tests depend on `internal/testhelp`; `testhelp` does not depend on ML packages. The cross-workflow assertion lives in `internal/testhelp/folder_fixtures_test.go`, which is the right package for a contract that spans two workflow files. Direct `github.com/hashicorp/hcl/v2` and `gopkg.in/yaml.v3` requirements match the new test imports. Their versions are the previous indirect versions, and `go.sum` already contains those modules.

Stale `FABRIC_TESTACC_WELLKNOWN` values in the unit tests are a tripwire, not a second fixture source. The builders never call `WellKnown()` for folders. Acceptance still reads only `WorkspaceRS`, which matches the documented scope. Data-source tests still use their own well-known entities and were not retargeted.

### 2. Abstraction quality of the shared steps

Both folders remain in every step, so a folder-id change is an in-place move rather than a Terraform destroy of the previous folder. That matches `folder_id` being optional and not force-new, and matches `ResourceFabricItem.Update` splitting display/description updates from `MoveItem`. Experiment still changes display name on the second step; model still changes description and keeps display name. The extraction did not alter those existing update shapes.

Create, read, move, read, and teardown stay on the existing test-case path. There is no extra destroy step and no new provider operation. Teardown remains the harness destroy of the whole configuration, including owned folders. The shared `CheckDestroy` helper is unchanged and still only names the item address; that is pre-existing and is not the destroy mechanism for the new folders.

The HCL tests parse the same builder output and reject a static folder id, a missing owned-folder reference, a dropped folder on the move step, and a repeated folder name across runs. The lifecycle tests execute that same output. Those two layers cannot drift inside a package.

### 3. Fake-store workaround stays test-local

`typedHandler.Elements` converts non-matching elements when the entity type looks like a Fabric item. Sharing one `elements` slice between generic items and folders is the limitation the lifecycle tests avoid. Replacing `ServerFactory.Core.FoldersServer` with the folder server from a second `NewFakeServer()` gives folder CRUD its own store while item calls stay on the first factory. The HTTP transport is then built from the swapped factory and serialized with a test-only mutex.

The factory replacement calls `GetTestUnitProtoV6ProviderFactories`, so the echo provider and client construction stay on the helper path. `NewTestUnitCase` has no transport argument; overriding `ProtoV6ProviderFactories` after the swap is the local extension point. Both copies do the swap before building the transport. Product, SDK, retry, and classification code are untouched.

The duplicated transport and HCL walker are not a layering defect. The two files differ only by package and resource type, which is required because `mlexperiment_test` and `mlmodel_test` cannot share unexported step builders. A `testhelp` helper cannot see those builders without a new callback API. Keeping the workaround next to the only two callers matches the per-service external-test pattern and the instruction not to change fake runtime behavior for this task.

### 4. Workflow lock matches the concurrency constraint

GitHub treats a matrix job concurrency group as one group per matrix leg and keeps only one pending leg. A job-level lock would therefore drop legs or serialize the CLI matrix. The change keeps concurrency on the workflow, uses one static group shared with cleanup, and sets `cancel-in-progress: false`. The `test` job still has `fail-fast: false` and four CLI entries, and it has no job-level `concurrency` key. `test-auth-spn` is covered by the same workflow group.

That group is coarser than the jobs that actually use `environment: development`. `changes`, `checkbuild`, and `checkdocs` also occupy it, and a docs-only run can hold or replace the single pending run. That coarseness is the mechanism that keeps the matrix inside one concurrency unit. `DEVELOPER.md` states the non-FIFO pending-run behavior, the decision not to cancel an active run, and the paths the lock does not cover: local commands, forks, other repositories, and manual cleanup. No third in-repo workflow uses the development environment. Splitting lint and docs into another workflow would preserve a narrower lock, but it is not required by the current contract and would change required-check topology beyond this fix.

The workflow test enforces the architectural invariant that can be checked statically: both files share a non-empty group with no expression, neither cancels in progress, the `test` job exists, and that job has no concurrency key. It does not hard-code the group string, so a coordinated rename of both files remains valid.

### 5. Naming and runner conventions

Copyright headers, MPL-2.0 lines, and `package *_test` match coding standards. `TestAcc_MLExperimentResource_CRUD` and `TestAcc_MLModelResource_CRUD` keep the acceptance names. New unit tests use `TestUnit_FolderFixtures_`, which is a valid Task designator: `Taskfile.yml` runs `^TestUnit_{{.TEST_NAME}}`. `task testunit -- FolderFixtures` selects the new contract tests without replacing `TestUnit_MLExperimentResource_CRUD` or `TestUnit_MLModelResource_CRUD`.

`*CRUD` in the new names means the shared CRUD step contract, not a second resource lifecycle. `*Lifecycle` is the fake Terraform execution of those steps. The pair is internally consistent. The testing-patterns suffix table still describes the existing `*Resource_CRUD` tests, which remain present. Non-parallel `resource.Test` plus `t.Setenv` is the correct choice: `Setenv` cannot be used after `t.Parallel()`, and these tests must not race the process environment with other tests that call `WellKnown()`.

### Residuals considered and not filed

- Workflow lock includes non-environment jobs. Accepted, as recorded above, because a per-job lock drops matrix legs.
- Duplicate test adapters. Accepted, because the external test packages cannot share the unexported builders without a new helper API.
- Lifecycle tests do not themselves set `FABRIC_PREVIEW`. Accepted. `fabric_folder` is preview-gated, and `Taskfile.yml` sets `FABRIC_PREVIEW` for every task, which is the documented runner contract. A direct `go test` without that variable fails closed during plan; it does not skip the preview check.
- Shared `CheckDestroy` still names only the item. Pre-existing. Owned folders are destroyed with the configuration by the test harness, not by a new destroy callback.
- One-pending replacement can drop a queued required check. Documented operational limit, not an implementation mismatch.

## Tests

### Actually run
- Local `actionlint` v1.7.12 on `.github/workflows/test.yml` and `.github/workflows/cleanup-script.yml`: exit 0, no findings.
- `go test -count=1 -run TestUnit_FolderFixtures_ ./internal/services/mlexperiment ./internal/services/mlmodel ./internal/testhelp` with `TF_ACC=0` and no preview flag: `internal/testhelp` passed; both lifecycle tests failed closed in plan because `fabric_folder` requires preview mode. HCL contract tests in those packages were not the failing tests.
- Same command with `FABRIC_PREVIEW=true` and `TF_ACC=0`, using the existing local Go toolchain and an existing local Terraform CLI: all three packages passed. This covers the five new tests (two lifecycle, two HCL contract, one workflow isolation).
- `go test -count=1 -run TestUnit_ ./internal/services/mlexperiment ./internal/services/mlmodel` with `FABRIC_PREVIEW=true` and `TF_ACC=0`: both packages passed. `go test -list TestUnit_` shows 7 unit tests in each package, including the two new tests.
- `go mod tidy -diff` with modules readonly and proxy disabled: exit 0, no diff.

### Reported and not re-executed here
- Earlier full ML package runs and the five new tests were reported passing before this review. The unit portion of that claim was independently re-run as listed above.
- Actionlint and tidy were reported clean. Both were independently re-run and were clean.
- Live acceptance and any cloud or GitHub-hosted workflow execution were not run. The workflow concurrency behavior is verified by source inspection, the YAML contract test, and actionlint, not by a queued GitHub run.
- The full provider suite, OpenTofu matrix, and auth-method acceptance jobs were not run.

### Limitations
- No live acceptance, cloud calls, or external writes. `TF_ACC` stayed unset for every command.
- GitHub's one-pending replacement and cross-run queueing were not executed on a hosted runner. The static contract matches the documented GitHub concurrency rule, but hosted scheduling was not observed.
- Direct `go test` is not the documented entrypoint. The passing rerun reproduced the Taskfile preview variable without invoking `task testacc`.
- No global tool installation. Commands used the already present local toolchain.
