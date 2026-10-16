## Review Summary

- **Round**: 6
- **Stage**: 6 — final stage
- **Theme**: Polish & hardening
- **Mode**: resolved-plan
- **Tier**: High, as explicitly assigned
- **Review Attempt**: 1
- **Stage Attempt**: 2 — final ordinary attempt; no third attempt authorized
- **Model**: gpt-6.1-sol
- **Reasoning**: max, per assigned reviewer contract
- **Artifact**: `reviews\task-terraform-folder-fixtures-attempt-1-stage-attempt-2-review-6-gpt-6.1-sol.md`
- **Issues Found**: 0
- **Verdict**: CLEAN

This is an independent assessment of the complete current candidate, not just the
three recently adjusted test files. No previous review artifacts or private
investigation evidence were consulted. CLEAN applies to this reviewer's assigned
scope; it is not an aggregate verdict for other reviewers or stages.

## Evidence Checklist

- [x] Compared the current tracked diff against HEAD and read all three new
  `folder_fixtures_test.go` files directly: eight changed/new source, workflow,
  module and documentation files in total. Separately checked unchanged
  `go.sum`, shared helpers, fake handlers and cleanup workflow as context.
- [x] `internal\services\mlmodel\resource_ml_model_test.go:240-293`: inspected
  the acceptance entrypoint and both generated CRUD steps. The same two owned
  folder configurations persist across creation and movement; `folder_id`
  references switch between their Terraform resource IDs, with paired state
  assertions and the original display-name/description behavior preserved.
- [x] `internal\services\mlexperiment\resource_ml_experiment_test.go:214-268`:
  independently traced the corresponding experiment path, including its
  display-name update, description update and destination-folder checks.
- [x] `internal\testhelp\testacc.go:162-180` and
  `internal\testhelp\utils.go:190-192`: verified that the existing folder helper
  generates owned resources with randomized names and that `RefByFQN` produces
  a Terraform reference, rather than a stale literal folder ID.
- [x] Both service `folder_fixtures_test.go:50-126` files: inspected HCL parsing,
  requested workspace checks, exactly two retained folder declarations,
  destination references, movement to the other folder and name separation
  across two generated runs. The seeded stale well-known IDs are synthetic.
- [x] Both service `folder_fixtures_test.go:21-48` files and
  `internal\testhelp\testunit.go:66-139`: verified that lifecycle regressions
  reuse the actual acceptance-step builders with unit-mode provider factories
  and SDK fake transport. The separate folder store and per-test transport
  mutex are confined to tests. Both lifecycle regressions actually passed.
- [x] `.github\workflows\test.yml:23-27,310-350`,
  `.github\workflows\cleanup-script.yml:12-14` and
  `internal\testhelp\folder_fixtures_test.go:16-66`: checked the identical
  constant workflow-level group, non-cancelling behavior, retained matrix and
  absence of a conflicting lock on the matrix job. Ran the workflow regression.
- [x] `internal\testhelp\folder_fixtures_test.go:17-26`: checked the
  `runtime.Caller` success guard, named concurrency configuration and narrowly
  scoped, explained external-YAML-tag suppression. Both transport interfaces
  now name their request parameter. These final files were included in the
  actual configured linter run.
- [x] `DEVELOPER.md:318-322,353-367`: cross-checked setup output and fixture
  precedence against `tools\scripts\Set-WellKnown.ps1`,
  `internal\testhelp\wellknown.go` and both workflows. Reviewed the documented
  serialization, one-pending replacement and local/fork/manual coordination
  limitations.
- [x] `go.mod` and `go.sum`: confirmed that HCL and YAML retain their existing
  versions while becoming direct test dependencies; no checksum-file change.
  Independently ran `go mod tidy -diff`.
- [x] Independently ran 15 correctly scoped offline tests, configured
  golangci-lint with `--fix=false`, actionlint on both workflows and the source
  whitespace check. Results and the initial selector mistake are recorded below.
- [x] Recomputed candidate hashes after validation: all eight candidate files
  and unchanged `go.sum` retained the reviewed contents. No product files were
  edited.
- [ ] Live acceptance, cloud deletion and actual hosted workflow scheduling:
  deliberately not executed, as prohibited.

## Effective Review Decision

The driving instruction supplies High effort, the Stage 6 polish/hardening lens,
the exact model/reasoning assignment and the second-attempt limit. The supplied
plan permits its two reviewer slots rather than requiring an additional model;
this reviewer did not delegate or dispatch any other reviewer. The primary
driver's exact model identity is unavailable.

The change is moderate in breadth: two acceptance-step builders, three regression
files, workflow concurrency, developer guidance and dependency classification.
The important interactions are resource lifetime across Terraform steps, fake
store access under concurrent graph execution, and workflow-wide exclusion
against destructive cleanup. Those observations justify reviewing the complete
candidate despite the latest adjustments being test-only style/hardening changes.
Earlier clean results were not carried forward.

## Clean Review Reasoning

### Fixture lifetime and scope

Each builder calls the established folder helper twice before constructing its
step slice. Reusing those two HCL fragments in both steps avoids accidental
replacement or removal of the original folder when the item moves. Terraform
references provide creation and destroy dependency edges; checks compare the
item's folder state with the corresponding owned resource's ID. The acceptance
entrypoints still obtain only their shared workspace from well-known data.
Model and experiment creation, subsequent reads, metadata updates, moves and
the test framework's automatic destroy phase remain intact.

No provider runtime, Fabric API handling, retries or error classification changed.
The existing non-fixture unit CRUD tests were also retained and passed with the
new regressions.

### Performance, isolation and observability

The HCL regressions use small, bounded loops. The lifecycle transport serializes
only its own SDK fake requests, not product traffic or other test cases. Reading
`internal\testhelp\fakes\fake_typedhandler.go` confirmed why combining folders
with the generic item store is unsuitable; the dedicated folder handler store
avoids that fake limitation without changing production behavior.

The workflow's broader serialization is a deliberate shared-environment safety
tradeoff, not an accidental per-matrix bottleneck. Matrix strategy and internal
parallelism are unchanged. The comments and developer documentation explicitly
describe the queue replacement cost rather than promising FIFO execution.

Test failures expose parsing, workspace, ownership, retained-folder, movement
and workflow-lock conditions through focused messages. The `runtime.Caller`
guard now fails explicitly if source location cannot be obtained. No new
production path needs additional logging or telemetry for this change.

### Naming, documentation and dependencies

The builder, lifecycle, HCL and transport names make their roles distinguishable.
The transport/store comments explain a real test constraint. The YAML suppression
is limited to the required GitHub field spelling and does not disable the
configured naming linter for the surrounding file.

Documentation correctly limits folder independence to these two CRUD tests;
it does not imply the shared workspace or other acceptance fixtures disappeared.
It also distinguishes local setup output from the hosted Actions variable.
Dependency changes match the new direct imports without upgrading libraries.
No actionable Stage 6 defect was found.

## Validation Actually Executed

Existing local tools were used without installation: Go 1.27.1, Task 3.54.0,
Terraform CLI 1.16.6, golangci-lint 2.14.0 and actionlint 1.7.12. Go build/module
and linter caches, subprocess working storage and temporary-file destinations
were confined to the supplied sibling tools area. Dependency network resolution
was disabled. Test invocations used only `task testunit`, explicitly with
`TF_ACC=0`; the new lifecycle cases used SDK fake provider transport.

| Check | Command | Observed result |
| --- | --- | --- |
| Full MLModel unit package | `task testunit -- '.*' github.com/microsoft/terraform-provider-fabric/internal/services/mlmodel` | PASS, 7 tests; package reported 100.0% statement coverage |
| Full MLExperiment unit package | `task testunit -- '.*' github.com/microsoft/terraform-provider-fabric/internal/services/mlexperiment` | PASS, 7 tests; package reported 100.0% statement coverage |
| Workflow isolation regression | `task testunit -- FolderFixtures_WorkflowIsolation github.com/microsoft/terraform-provider-fabric/internal/testhelp` | PASS, 1 test |
| Configured Go lint | `golangci-lint run --fix=false --timeout=5m .\internal\testhelp .\internal\services\mlmodel .\internal\services\mlexperiment` | PASS, exit 0, `0 issues.` |
| Both workflow files | `actionlint -color .\.github\workflows\test.yml .\.github\workflows\cleanup-script.yml` | PASS, exit 0, no diagnostics |
| Module consistency | `go mod tidy -diff` | PASS, exit 0, no diff |
| Source whitespace | `git diff --check -- . ':!reviews'` | PASS, exit 0, no diagnostics |
| Final candidate integrity | SHA-256 comparison before/after the targeted validation | PASS for eight candidate files and unchanged `go.sum` |

The service coverage numbers are the runner's package statement figures, not
branch coverage, fake-store coverage, or evidence of live service recovery.
Sequential Task invocations regenerate their normal local coverage/result
outputs; those outputs do not constitute a combined full-repository report.

The linter emitted configuration deprecation warnings for `gofumpt.extra-rules`
and `gomodguard`, but no issues. Independently checking `.golangci.yml` against
HEAD confirmed that configuration is unchanged; no configuration migration was
made as part of this review.

### Initial command-selection error

An initial invocation attempted to scope Task using a `TEST_PATH=...` variable
assignment. The task's local variable definition retained the default all-package
selector, which was visible in its emitted command. That over-broad unit run was
stopped, and the documented positional package selector was dry-run checked
before the three correctly scoped commands above were executed successfully.

The broad run reported failures outside the selected review packages, including
a missing local well-known fixture at
`internal\services\deploymentpipelinera\resource_deployment_pipeline_role_assignment_test.go:192`.
That call site and `internal\testhelp\wellknown.go` were independently confirmed
unchanged against HEAD. Other broader-run failures were not diagnosed or
baseline-compared. No clean full-suite result or blanket pre-existing-failure
claim is made.

### Reported versus independently observed

The supplied input reports earlier failing-then-passing regressions, rebuilds and
successful final checks. Those historical executions were not used as independent
validation. This attempt freshly inspected the full current source and reran the
15 targeted tests, linter and static checks listed above. It did not reproduce
the historical red phase or consume earlier review conclusions.

## Reviewed Candidate Fingerprints

SHA-256 values bind this assessment to the final current files, including the
three recently adjusted regression files.

| Repository-relative path | SHA-256 |
| --- | --- |
| `.github\workflows\test.yml` | `b29495131a66b7f125ddc62519626c59fd827ca48638d0e2cdef044048e0b4e7` |
| `DEVELOPER.md` | `cb6e2f0ef171fd15a9f83bace15d240713f9fc1ef39ec400419e622272288c75` |
| `go.mod` | `14b6e7c591f115541dc89e0551cccaf02b48643ed1a9cadf6f8e84ce38538b08` |
| `internal\services\mlmodel\resource_ml_model_test.go` | `0e3dd4f480daba5d6778c1ad9c1fad7c4469c759449f12fbbd09170817550913` |
| `internal\services\mlexperiment\resource_ml_experiment_test.go` | `361ddcb4e66297a5d78fbd0a44eedcb7fc99105ec705b0fc2d7ecb576c1f0611` |
| `internal\services\mlmodel\folder_fixtures_test.go` | `b5f4045feda97edfe2b7afa90b193af52309699ff4e0fe07e3c11ac27ba41314` |
| `internal\services\mlexperiment\folder_fixtures_test.go` | `3ca99bf713f0e5365c84103e6e8a6a066c42b092d807e4340cabbe3ac5b1c0c2` |
| `internal\testhelp\folder_fixtures_test.go` | `a32d52787c460c9eb30156cec182278c0791c1c13e6425ffa84749d4712dfd9e` |

## Residual Risks and Unverified Boundaries

- Live folder permissions, service-side move behavior and deletion were not
  exercised. Local Terraform execution against SDK fakes verifies the shared
  step configurations and lifecycle, not actual cloud behavior.
- Hosted GitHub scheduling and cleanup exclusion were inspected and regression
  checked, not executed. One-pending replacement, operator cancellation,
  timeouts and local/fork/manual cleanup remain operational boundaries; the
  workflow lock is not an external cleanup interlock.
- Validation used one Windows Terraform version. The hosted Terraform/OpenTofu
  matrix, race detector, dedicated performance benchmarks and a successful
  complete repository suite were not established by this review.
- No product edits, live acceptance, cloud cleanup, external writes, commits,
  pushes, forks, PR creation or nested agents were performed. This final ordinary
  Stage 6 attempt does not authorize another attempt or any publication action.
