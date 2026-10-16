# Terraform folder fixtures: architecture and patterns

## Review Summary

- **Tier**: High
- **Stage / Round**: STAGE2 / 2 of 6
- **Theme**: Architecture & patterns
- **Mode**: resolved-plan; independent, artifact-capable REVIEW ONLY
- **Model**: gpt-6.1-sol
- **Reasoning**: max
- **Review Attempt**: 1
- **Stage Attempt**: ATTEMPT1 / 1
- **Artifact**: `reviews\task-terraform-folder-fixtures-attempt-1-stage-attempt-1-review-2-gpt-6.1-sol.md`
- **Issues Found**: 0
- **Verdict**: CLEAN

### Scope and effective policy

Reviewed the complete supplied Stage 2 lens and frozen diff, then independently
inspected the corresponding public source and repository instructions. No prior
review artifacts or private investigation evidence were consulted.

The supplied human policy selects two slots, `gpt-6.1-sol/max` and
`grok-4.7/xhigh`, with the normal two-attempt High gates for stages 2–6 and no
third-slot exception. This artifact is the GPT slot only. It neither reports the
other slot's result nor advances the stage. The supplied plan does not establish
the coordinating agent's exact model identity.

The architectural scope is ownership and lifetime of the two folders used by
each positive ML CRUD test, reuse of those CRUD steps in offline regressions,
test-only fake composition, and the shared-environment cleanup lock. Provider
API, retry, error-classification and authentication behavior are not changed by
the reviewed diff.

## Evidence Checklist

- [x] Read the applicable conventions in `.github\copilot-instructions.md`,
  `.github\instructions\coding-standards.instructions.md`,
  `.github\instructions\testing-patterns.instructions.md` and
  `CODING_GUIDELINES.md`; inspected the fake-handler guidance as supporting
  architecture documentation. No checkout-local `AGENTS.md` was found.
- [x] Traced both acceptance entrypoints into their shared step builders:
  `internal\services\mlmodel\resource_ml_model_test.go:240-302` and
  `internal\services\mlexperiment\resource_ml_experiment_test.go:214-277`.
  Each resolves only `WorkspaceRS` from well-known data and generates two owned
  folder resources once per CRUD sequence.
- [x] Verified helper reuse and graph construction against
  `internal\testhelp\testacc.go:162-180` and
  `internal\testhelp\utils.go:189-192`. Both folders remain in both step
  configurations; the item references the first folder's `id` for creation and
  the second folder's `id` for the move.
- [x] Compared the composition with the existing owned-parent pattern in
  `internal\services\folder\resource_folder_test.go:310-399`. Inspected the
  unchanged ML constructors and `internal\pkg\fabricitem\resource_item.go:186-280`
  to confirm that model rename restrictions, experiment rename support,
  description updates, move/read behavior and deletion retain their existing
  production implementation.
- [x] Inspected both new service regression files,
  `internal\services\mlmodel\folder_fixtures_test.go:21-111` and
  `internal\services\mlexperiment\folder_fixtures_test.go:21-111`, together with
  `internal\testhelp\testunit.go` and
  `internal\testhelp\fakes\fake_server.go`,
  `internal\testhelp\fakes\fake_typedhandler.go` and
  `internal\testhelp\fakes\fabric_folder.go`. Confirmed that the separate folder
  store and serialized transport are local test accommodations, not runtime
  provider changes or a replacement fake API.
- [x] Compared `.github\workflows\test.yml:23-27` with
  `.github\workflows\cleanup-script.yml:12-14`, inspected the matrix and searched
  the other workflow definitions for acceptance/cleanup consumers. The two
  environment-consuming workflows use the same constant workflow-level group
  and disable cancellation of the active run; no per-matrix-leg lock was added.
- [x] Reviewed the workflow regression in
  `internal\testhelp\folder_fixtures_test.go:16-53` and the operational contract in
  `DEVELOPER.md:318-376`. The documentation explicitly describes pending-run
  replacement, non-FIFO scheduling, and the coordination required for local,
  fork and manual cleanup paths.
- [x] Independently ran all five new regression tests through `task testunit`,
  workflow lint, offline `go mod tidy -diff`, and the scoped whitespace check.
  All passed. Before/after SHA-256 comparisons matched for all eight reviewed
  changed/new source, workflow, documentation and module files.

## Clean Assessment

1. **Fixture lifetime is modeled in Terraform rather than in external setup.**
   `FolderResource` remains the authoritative generator. Reusing the same two
   generated blocks across both steps preserves their resource addresses and
   prevents the create-to-move transition from implicitly deleting the first
   folder. `RefByFQN` and paired state checks express the actual dependency and
   expected destination instead of comparing against stale fixture IDs.

2. **The new abstraction has a bounded, useful responsibility.**
   `modelCRUDSteps` and `experimentCRUDSteps` return test steps; they do not
   authenticate, provision fixtures out of band, or hide acceptance execution.
   Each original acceptance entrypoint still owns the well-known workspace
   lookup and `NewTestAccCase`. The offline lifecycle test executes the very
   same steps through `NewTestUnitCase`, avoiding a second implementation of
   the behavior being tested. The existing create/read/update/move assertions
   remain service-specific.

3. **Fake limitations are contained without widening the production change.**
   The existing fake server has an unsynchronized shared element slice, and its
   generic item conversion assumes item-shaped entities. Isolating the folder
   handlers onto another existing fake server avoids feeding folder entities
   into that conversion. A mutex around the local fake transport protects
   concurrent Terraform requests to those stores. This does not serialize the
   real provider or the CI matrix. The small duplicated test scaffolding is
   acceptable here; extracting a new cross-service framework or refactoring the
   entire fake store is not necessary to make this scoped fix maintainable.

4. **Cleanup exclusion is placed at the correct ownership boundary.**
   A workflow-level lock spans the matrix's resource creation through teardown.
   Reusing the cleanup workflow's constant group, rather than locking each
   matrix job, preserves matrix concurrency and avoids the single-pending-slot
   problem between matrix legs. Broader workflow serialization and possible
   replacement of a pending required check are explicit accepted trade-offs,
   not an undocumented FIFO guarantee. Repository-local concurrency does not
   claim to protect commands issued outside that repository's workflows.

5. **Dependency and convention changes remain proportionate.**
   New tests stay in black-box test packages and use the existing SDK fake and
   provider-factory mechanisms. HCL and YAML imports promote already-pinned
   modules to direct dependencies without changing versions or adding a new
   runtime subsystem. The offline tidy check confirms that module
   classification is consistent. No architectural defect requiring a source
   change was identified in this stage.

## Validation and Limitations

### Actually run by this reviewer

Validation used the existing local tools, Go 1.27.1 and Task 3.54.0. Module
network resolution was disabled. `TF_ACC=0` was retained through the unit-test
runner; an existing Terraform executable was used only with SDK fakes.

| Invocation | Result |
| --- | --- |
| `task testunit -- FolderFixtures github.com/microsoft/terraform-provider-fabric/internal/services/mlmodel` | PASS: `TestUnit_FolderFixtures_ModelCRUD` and `TestUnit_FolderFixtures_ModelLifecycle`; 2 tests. |
| `task testunit -- FolderFixtures github.com/microsoft/terraform-provider-fabric/internal/services/mlexperiment` | PASS: `TestUnit_FolderFixtures_ExperimentCRUD` and `TestUnit_FolderFixtures_ExperimentLifecycle`; 2 tests. |
| `task testunit -- FolderFixtures github.com/microsoft/terraform-provider-fabric/internal/testhelp` | PASS: `TestUnit_FolderFixtures_WorkflowIsolation`; 1 test. |
| `actionlint -shellcheck= -pyflakes= .github\workflows\test.yml .github\workflows\cleanup-script.yml` | PASS with actionlint 1.7.12; ShellCheck and Pyflakes integrations disabled. |
| `go mod tidy -diff` | PASS: no proposed module changes; offline cache only. |
| `git diff --check` scoped to the five changed tracked files | PASS: no whitespace errors. |
| Before/after SHA-256 comparison of the eight frozen review files | PASS: source inputs unchanged by the review and validation. |

The module identifiers in the Task commands are public Go package identifiers,
not local filesystem paths. The Task runner also completed its normal coverage
conversion steps. These results are targeted regression evidence, not a claim of
comprehensive provider or cloud coverage.

### Reported by the caller, not independently rerun here

The supplied context reports full unit runs of both ML packages, seven tests
each, and earlier failing HCL/workflow regression runs before the fix. This slot
did not reproduce those full-package or pre-fix runs or inspect their logs.
The five new passing regressions, actionlint result and tidy result above were
independently verified rather than accepted solely from that report.

### Boundaries and remaining verification

- No live acceptance tests, Fabric/cloud operations, cleanup execution,
  authentication actions or external writes were performed. Fake-backed
  lifecycle success does not establish live folder permissions, eventual
  consistency, or cloud recovery.
- GitHub scheduling, pending-run replacement and cancellation behavior were
  inspected statically, not exercised by dispatching workflows. Manual
  cancellation, runner termination and cross-repository/local interference
  remain operational risks outside the stated lock boundary.
- A full Go lint run, the complete provider unit suite and the Terraform/OpenTofu
  CI version matrix were not executed in this architecture-review slot.
- No source edits, commits, pushes, forks or pull requests were made. Only this
  review artifact was added; the Task runner refreshed its local generated
  coverage/test outputs.

**Conclusion:** CLEAN for Stage 2 architecture and patterns, review attempt 1,
stage attempt 1. This is not a verdict on the other review slot, completion of
the remaining stages, or approval to publish or deploy.
