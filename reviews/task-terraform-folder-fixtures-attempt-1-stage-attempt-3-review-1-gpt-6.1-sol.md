# Terraform Folder Fixtures — Independent Broad Review

## Review Summary

- **Round**: 1
- **Stage / Stage Attempt**: 1 / 3
- **Top-level Review Attempt**: 1
- **Theme**: Broad sweep
- **Mode**: resolved-plan
- **Effort**: High; explicitly authorized single replacement review for this stage attempt
- **Model**: gpt-6.1-sol
- **Reasoning**: max
- **Base**: `86c83beedf75c9a894c36a5bd516e0ccde7037cb`
- **Artifact**: `reviews\task-terraform-folder-fixtures-attempt-1-stage-attempt-3-review-1-gpt-6.1-sol.md`
- **Issues Found**: 0
- **Verdict**: CLEAN

This is a substantive review of the current candidate, not recovery of an earlier
review. No prior review artifacts or private investigation evidence were consulted,
and no other reviewers were invoked. The verdict covers this broad-sweep review
only, not completion of the remaining High stages or live acceptance validation.

## Evidence Checklist

- [x] **Scope and baseline verified independently.** `git rev-parse HEAD` matched
  the stated base. Inspected the current tracked diff in
  `.github\workflows\test.yml`, `DEVELOPER.md`, `go.mod`, and the two ML resource
  test files, plus all three new `folder_fixtures_test.go` files. Read applicable
  `.github\copilot-instructions.md`, testing/coding instructions and
  `CODING_GUIDELINES.md`; used no prior-review conclusions.
- [x] **Each positive CRUD case owns two folders in its requested workspace.**
  `internal\services\mlmodel\resource_ml_model_test.go:240-250` and
  `internal\services\mlexperiment\resource_ml_experiment_test.go:214-224` obtain
  only `WorkspaceRS` from well-known data and call the existing
  `testhelp.FolderResource` twice. That helper generates fresh resource and
  display names and emits managed `fabric_folder` resources with the supplied
  workspace (`internal\testhelp\testacc.go:162-177`).
- [x] **Both folder definitions survive both steps; references establish the graph.**
  The model's configurations at
  `internal\services\mlmodel\resource_ml_model_test.go:259-288` and the
  experiment's at
  `internal\services\mlexperiment\resource_ml_experiment_test.go:234-263`
  join the same two folder HCL blocks in each step. The first item configuration
  references the first folder's `id`; the second references the other folder's
  `id`. `TestCheckResourceAttrPair` checks the resulting state against those
  resources rather than a copied fixture ID. `RefByFQN` emits a Terraform
  interpolation (`internal\testhelp\utils.go:189-192`), not a literal identifier.
- [x] **Existing create/read/update/move behavior is retained.** The model keeps
  its original display name while changing its description and folder
  (`internal\services\mlmodel\resource_ml_model_test.go:252-288`); the experiment
  still changes display name, description and folder
  (`internal\services\mlexperiment\resource_ml_experiment_test.go:226-263`).
  The unchanged generic update path separately performs metadata update,
  `MoveItem`, and readback
  (`internal\pkg\fabricitem\resource_item.go:186-249`). The unchanged schema
  permits folder updates rather than requiring replacement.
- [x] **Teardown remains Terraform-managed.** Both acceptance entrypoints still
  call `resource.Test` with `NewTestAccCase`, without disabling destruction
  (`internal\services\mlmodel\resource_ml_model_test.go:240-244`,
  `internal\services\mlexperiment\resource_ml_experiment_test.go:214-218`,
  `internal\testhelp\testacc.go:58-82`). Both folders remain managed resources in
  the final configuration. The final item-to-destination reference gives
  Terraform the reverse dependency needed to delete the item before its occupied
  destination folder; the retained source folder is empty after a successful
  move. The underlying item/folder delete methods remain unchanged
  (`internal\pkg\fabricitem\resource_item.go:252-280`,
  `internal\services\folder\resource_folder.go:235-263`).
- [x] **Static regressions reject a stale-folder dependency and premature removal.**
  `internal\services\mlmodel\folder_fixtures_test.go:48-110` and
  `internal\services\mlexperiment\folder_fixtures_test.go:48-110` parse the actual
  returned step configurations. They require two steps, two owned folders in
  each step, the requested workspace, a three-component folder-ID traversal,
  retention of the same folder resources, a changed move target, and distinct
  folder resource names across repeated builders. Synthetic stale well-known
  folder values are present throughout these tests.
- [x] **Lifecycle tests exercise the shared steps through isolated SDK fakes.**
  `internal\services\mlmodel\folder_fixtures_test.go:21-45` and
  `internal\services\mlexperiment\folder_fixtures_test.go:21-45` pass the same CRUD
  builders to `NewTestUnitCase` and `resource.Test`. They use fresh fake servers,
  replace only the folder fake with a separately backed store, and serialize
  transport `Do` calls with a shared mutex for each test. The split avoids the
  generic fake's item-only conversion assumptions
  (`internal\testhelp\fakes\fake_typedhandler.go:203-234,277-292`); the provider
  factory uses a fake credential and the injected transport
  (`internal\testhelp\testunit.go:113-126`). These adaptations are test-only.
- [x] **Destructive cleanup and acceptance runs share the same whole-workflow lock.**
  `.github\workflows\test.yml:23-27` and
  `.github\workflows\cleanup-script.yml:12-14` use the identical constant
  `fabric-testacc-environment` with `cancel-in-progress: false`. Cleanup's
  destructive invocation is inside its locked workflow
  (`.github\workflows\cleanup-script.yml:24-45`). The test matrix remains
  internally concurrent with `fail-fast: false`
  (`.github\workflows\test.yml:309-342`); the change does not introduce a
  per-matrix-job lock. The new guard explicitly checks matching constant groups,
  non-cancellation and absence of a job-level test lock
  (`internal\testhelp\folder_fixtures_test.go:16-53`).
- [x] **Queue/cancel semantics and operational exclusions are explicit.**
  `DEVELOPER.md:353-369` explains deliberate serialization, one pending run,
  replacement rather than FIFO ordering, rerunning superseded checks, and
  operator coordination for local, fork, other-repository and manual cleanup.
  This is not represented as a durable queue or a cross-repository mutex.
- [x] **Documentation, dependencies and security blast radius checked.**
  `DEVELOPER.md:318-322` agrees with the populated-environment-variable lookup
  preceding file loading in `internal\testhelp\wellknown.go:41-64`.
  `go.mod:13,28` promotes existing HCL/YAML versions to direct requirements for
  the new tests; no version upgrade or `go.sum` change is present. The diff
  changes no production Fabric API, retry, classification, authentication or
  authorization code, and introduces no credential material or new workflow
  permission. New Go files retain the repository's license headers and
  black-box test packages.
- [x] **Independent read-only validation passed.** Ran whitespace validation
  against the base, actionlint on both relevant workflows, and `gofmt -l` on all
  five reviewed Go test files; all exited successfully with no diagnostics.
- [ ] **Unit and live acceptance execution were not rerun by this reviewer.**
  The distinction from caller-reported results and the reasons are recorded below.

## Clean Assessment

No actionable correctness, security, requirement-conformance or blast-radius
defect was identified in the frozen candidate.

The folder lifetime change addresses the stated dependency rather than masking
an API error: ownership moves from pre-existing folder IDs into the same Terraform
configuration as the item. Keeping both blocks in both steps avoids deleting the
source folder while the item is still moving. Reusing the builders in the
fake-backed lifecycle tests also prevents a separate, simplified test
configuration from silently diverging from the acceptance sequence.

The cleanup protection is at the correct granularity for the stated shared
environment. Its constant group excludes cleanup for the duration of an active
test workflow, including the matrix's ordinary teardown, while leaving matrix
legs free to run concurrently. Serializing documentation/build jobs along with
acceptance work is a deliberate cost of the specified workflow-level lock, not
an unnoticed per-job serialization change.

The scope remains narrow: two acceptance configurations, their regressions,
workflow coordination, documentation, and direct dependency bookkeeping.
Production item behavior and shared fake implementations are unchanged.

## Tests Actually Run and Caller-Reported Evidence

### Executed independently

| Check | Result |
| --- | --- |
| `git rev-parse HEAD` | Matched `86c83beedf75c9a894c36a5bd516e0ccde7037cb`. |
| `git diff --check 86c83beedf75c9a894c36a5bd516e0ccde7037cb` | Exit 0; no whitespace diagnostics. |
| `actionlint -no-color .github\workflows\test.yml .github\workflows\cleanup-script.yml` | Exit 0; no diagnostics. |
| `gofmt -l` on the two modified ML resource test files and three new fixture test files | Exit 0; no files listed. |
| Base comparison for `internal\pkg`, `internal\provider`, shared test-case helpers, `internal\testhelp\fakes`, `go.sum`, and `.github\workflows\cleanup-script.yml` | No differences. |

These are static checks, not successful executions of the Go regression tests.
Existing isolated tools were used; no tool installation or cloud access occurred.

### Caller-reported, not independently rerun

The supplied review context reports:

- Both complete ML unit packages passed, seven tests per package.
- All five new tests passed, including the workflow-isolation guard.
- HCL and workflow regressions were exercised in their failing state before the fix.
- Actionlint and `go mod tidy -diff` were clean.
- No live acceptance/cloud validation was executed.

Actionlint was independently repeated as listed above. The unit-test counts,
red-phase execution and tidy result remain caller-reported evidence; no earlier
reports or review artifacts were opened to substantiate them.

### Validation boundaries and residual risks

- Unit tests were not rerun because the required `task testunit` entrypoint
  recreates `.coverdata` and writes `testresults.json` plus coverage outputs
  (`Taskfile.yml:274-290,362-391`). Those writes exceed this review's
  single-artifact change allowance. Raw Go test execution was not substituted.
- Live acceptance, `TF_ACC`-enabled execution, cloud resource creation and cleanup
  were deliberately not invoked. Actual API permissions, eventual consistency,
  service-side move behavior and absence of leaked cloud folders remain unverified.
- Source inspection confirms that the lifecycle tests use real Terraform test
  steps against fakes. The separate fake stores do not reproduce service-side
  folder membership constraints, so those tests are not proof of real
  nonempty-folder deletion enforcement or end-to-end cloud recovery.
- GitHub's one-pending replacement rule can supersede a queued test or cleanup
  run; this is the documented, explicitly selected policy. It does not guarantee
  fairness, FIFO execution or eventual cleanup under sustained submissions.
- `cancel-in-progress: false` prevents cancellation by a newer competing run,
  not manual cancellation, job timeout or runner failure. Such interruptions
  can still leave resources requiring coordinated cleanup. Local and
  cross-repository invocations remain outside the lock.

Only this requested review artifact was created. Product source, existing review
artifacts and generated test outputs were left untouched; no commit, push or
external write was performed.
