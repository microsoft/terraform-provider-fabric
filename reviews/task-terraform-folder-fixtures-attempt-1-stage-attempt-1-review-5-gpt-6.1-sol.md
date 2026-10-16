# Terraform folder fixtures — testing and coverage review

## Review Summary
- **Round**: 5
- **Stage**: 5
- **Theme**: Testing & coverage
- **Mode**: resolved-plan; review only
- **Review Attempt**: 1
- **Stage Attempt**: 1
- **Model**: gpt-6.1-sol
- **Reasoning**: max, as assigned by the caller
- **Artifact**: `reviews\task-terraform-folder-fixtures-attempt-1-stage-attempt-1-review-5-gpt-6.1-sol.md`
- **Issues Found**: 0
- **Verdict**: CLEAN

This verdict covers this reviewer and this stage only. The supplied High plan's two-slot override and normal maximum of two stage attempts remain unchanged; this review neither advances the aggregate gate nor requests an exception.

## Evidence Checklist
- [x] Read the supplied current-stage scope and diff, then independently inspected all eight changed/new source and documentation files. No earlier review artifacts or private investigation evidence were consulted.
- [x] Traced the acceptance entrypoints and shared builders in `internal\services\mlmodel\resource_ml_model_test.go:240` and `internal\services\mlexperiment\resource_ml_experiment_test.go:214`. The corresponding new lifecycle tests execute those same builders, rather than separate approximations.
- [x] Inspected both new `folder_fixtures_test.go` files, including their HCL traversal assertions, repeated-builder isolation checks, actual Terraform test execution, and serialized fake transports.
- [x] Inspected `internal\testhelp\testunit.go`, `testacc.go`, `utils.go`, and the item/folder fake implementations. Verified the shared helpers and fake implementation files are unchanged against the checkout's current base.
- [x] Traced the real provider create/read/update/move/delete path through `internal\pkg\fabricitem\resource_item.go` and folder CRUD through `internal\services\folder\resource_folder.go`; the tests substitute the SDK transport, not the Terraform CLI or provider resource implementation.
- [x] Inspected `internal\testhelp\folder_fixtures_test.go:16`, both workflow concurrency declarations, the preserved test matrix, and the related operational caveats in `DEVELOPER.md`.
- [x] Independently ran the correctly bounded unit commands below with `TF_ACC=0`: seven MLModel tests, seven MLExperiment tests, and the workflow-isolation test all passed. This includes all five newly added tests.
- [x] Ran a scoped tracked-diff whitespace check successfully. This is not presented as formatting validation of untracked files or as a substitute for the unit results.
- [ ] Did not rerun the reported red-before-fix demonstrations, actionlint, or dependency-tidiness check. Those are caller-reported results, not independently reproduced results in this review.
- [ ] Did not execute live acceptance, cloud calls, cleanup workflows, hosted workflow scheduling, race detection, or mutation testing. No source changes, commits, pushes, publication, or nested agents were used.

## Justified Clean Assessment

### 1. The new runtime tests exercise the acceptance configuration, not merely similar HCL

`modelCRUDSteps` and `experimentCRUDSteps` are called by both their positive acceptance tests and the new `TestUnit_FolderFixtures_*Lifecycle` tests. Each builder calls the existing `FolderResource` helper twice, retains the resulting HCL strings for both steps, and uses `RefByFQN` to connect the ML item to the first folder and subsequently the second.

The lifecycle tests call `resource.Test` with `NewTestUnitCase`, whose `IsUnitTest` setting permits actual Terraform execution without selecting acceptance tests. `GetTestUnitProtoV6ProviderFactories` supplies the real provider with a fake credential and SDK fake transport. This is stronger than merely parsing the generated configuration: Terraform resolves the graph and runs provider operations against the in-memory SDK implementation.

The state assertions check display name, description, and the relationship between the ML item's `folder_id` and the managed folder's actual `id`. The second step retains the MLModel name, updates its description, and changes its folder; the MLExperiment step also exercises its supported rename. The fake item implementation updates `FolderID` in its move operation, and the provider reads the item back after updating. These checks therefore observe the create/read and update/move/read outcomes, rather than only accepting a successful command exit.

The Terraform testing library performs a final destroy and propagates destroy errors. The existing test helper does not independently query the fake stores for remote absence afterward, so this review does not claim such an assertion exists.

### 2. The HCL assertions directly protect the changed fixture-lifetime contract

Both `TestUnit_FolderFixtures_*CRUD` tests require exactly two steps and parse every generated configuration with the HCL parser. Their checks establish:

- Exactly two owned folder resource names are present in each step.
- Each folder uses the requested workspace.
- The ML item's folder expression contains a traversal to an owned `fabric_folder` resource's `id`, rejecting a literal stale identifier.
- Both original folder resource names remain in the second step.
- The second step targets the other owned folder.
- A second invocation of the builder does not reuse the first invocation's folder resource names.

Both structural and lifecycle tests install deliberately stale synthetic well-known folder entries. Together with the inspected acceptance entrypoints, this gives relevant regression protection for removing the old static-folder dependency. The builders no longer read those folder entries.

The structural tests and runtime tests are complementary: HCL inspection establishes the intended ownership and dependency shape, while real Terraform verifies that this generated configuration can be applied and moved through the real provider.

### 3. The mock adaptations are appropriate but have explicit boundaries

Each lifecycle test creates a fresh server. Its folder endpoints are taken from another fresh server, avoiding the generic item store's conversion of folder objects through item-only fields. The reason is supported by `fake_typedhandler.go`, including `asFabricItem`, rather than being an unexplained bypass.

The added `folderFixtureTransport.Do` serializes access to the in-memory stores with a mutex and releases it with `defer`. The same wrapped transport is supplied to the provider factories for that test. The tests are not marked parallel, and their environment changes use `t.Setenv`.

Important limitations remain:

- Item creation and movement in `fabric_item.go` do not check that a referenced folder exists in the separate folder store.
- Generic deletion in `fake_handlers.go` checks entity presence but does not model rejection of deleting a nonempty folder.
- Serialized SDK fake requests do not establish production concurrency correctness.
- The existing `CheckDestroy` callback checks the supplied pre-destroy state, not a post-destroy service inventory.

These boundaries prevent treating the fake lifecycle pass alone as proof of live Fabric folder validation or cleanup-race behavior. They do not identify an actual regression in this change: the inspected builders retain both folders and express the required references, and the separate workflow guard addresses the supported shared-environment exclusion. Requiring unrelated service fault injection or complete production concurrency simulation would exceed the changed behavior.

### 4. The workflow regression checks the actual exclusion mechanism

`TestUnit_FolderFixtures_WorkflowIsolation` reads both real workflow YAML files. It checks for a nonempty, non-expression, identical root-level concurrency group; rejects active-run cancellation; requires the test matrix job; and rejects a per-test-job concurrency declaration.

The inspected workflow uses that shared constant at workflow scope with cancellation disabled. The test matrix and its `fail-fast: false` behavior remain in place. Consequently the regression is aimed at the actual mechanism protecting the whole run, including teardown, rather than at a misleading per-matrix-job lock.

The documentation correctly limits this protection to the coordinating workflows and discloses pending-run replacement and the absence of a FIFO guarantee. It does not claim that local commands, forks, other repositories, or manual cleanup are covered.

This local YAML test is not a hosted GitHub scheduler integration test. That limitation is explicit and is not, by itself, a source-backed defect.

### 5. Existing coverage is preserved without overstating the coverage report

The two package runs include their existing resource attribute, import, duplicate-create/CRUD, and data-source tests, as well as the two new tests per package. These tests were not replaced by the fixture-specific checks.

The two ML packages report 100.0% statement coverage, but their product files largely configure shared implementations. That number does not establish branch coverage of `internal\pkg\fabricitem`, coverage of the new `_test.go` builders themselves, or complete mock fidelity. The helper-only run reports 1.1%, as expected for selecting one workflow test rather than the entire helper suite.

No changed product error-handling, API, retry, or classification path was found that would justify a mandatory additional negative-path regression test for this patch.

## Test Commands Actually Run

All actual test invocations used `TF_ACC=0`, the existing isolated toolchain and cached Terraform executable, and caches/work directories confined to the approved sibling tooling area. Module downloads were disabled. No dependencies were installed.

For public-safe presentation, `$module` below denotes the module declaration in `go.mod`; the actual invocations passed the resulting module-qualified package strings literally. These are Go import arguments, not filesystem paths.

| Executed runner and selector | Repository package | Observed result |
| --- | --- | --- |
| `task testunit -- '(FolderFixtures\|MLModel)' "$module/internal/services/mlmodel"` | `internal\services\mlmodel` | PASS — 7 tests; 17.712 s |
| `task testunit -- '(FolderFixtures\|MLExperiment)' "$module/internal/services/mlexperiment"` | `internal\services\mlexperiment` | PASS — 7 tests; 15.720 s |
| `task testunit -- 'FolderFixtures' "$module/internal/testhelp"` | `internal\testhelp` | PASS — 1 test; 4.271 s |

Each successful invocation also completed the runner's coverage merge and report-generation steps. Thus all five new regressions and all existing unit tests in the two ML packages were independently verified.

### Additional attempted invocation and dry-run disclosure

An initial `task testunit` invocation attempted to set its selector and package through `TEST_NAME`/`TEST_PATH` command-line variables. The task's local definitions ignored those overrides and expanded to all unit tests. That invocation was explicitly stopped after the unintended expansion and out-of-scope failures became apparent, including a missing local well-known fixture file. It is **not** a successful full-suite result and is not used to assess unrelated packages.

A subsequent `task --dry testunit` inspection exposed the multi-package quoting behavior. The three successful runs therefore used the supported positional selector and a single package per invocation.

The initial broad attempt, the dry-run, and the final targeted runs are disclosed separately; only the final targeted runs establish the passing results above.

## Reported-Only Evidence and Remaining Limits

The supplied review context reports earlier failing-then-passing HCL/workflow regressions, actionlint success, and a clean dependency-tidiness check. This reviewer did not reproduce those claims, inspect earlier review artifacts, or modify source to perform mutations.

No live acceptance result, deployed behavior, actual shared-workspace recovery, complete-repository unit pass, hosted queue/cancellation execution, race-detector result, or mutation score is claimed. Real-service folder existence/emptiness behavior and operator coordination outside the workflow lock remain outside this local review.

Within the supplied change and this testing lens, the independently inspected source and passing targeted regressions support a **CLEAN** assessment with no actionable source-backed regression finding.
