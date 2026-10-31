## Review Summary

- **Round**: 6
- **Stage**: 6 — Polish & hardening
- **Review Attempt**: 1
- **Stage Attempt**: 1
- **Theme**: Polish & hardening
- **Mode**: resolved-plan
- **Model**: gpt-6.1-sol
- **Reasoning**: max
- **Artifact**: `reviews\task-terraform-folder-fixtures-attempt-1-stage-attempt-1-review-6-gpt-6.1-sol.md`
- **Issues Found**: 2
- **Verdict**: ISSUES_FOUND

Independent final-stage review of the current public-source change. Both findings
are narrowly scoped Go lint blockers in the newly added workflow regression test,
not defects in the owned-folder lifecycle. The assigned gate is High, with the
human-selected two-model review and at most two attempts per stage; this artifact
does not aggregate another reviewer's result or advance the gate.

## Evidence Checklist

- [x] Compared the supplied scope/diff with the current checkout: the workflow,
  developer documentation, module declarations, two existing CRUD test files and
  three new fixture regression files comprise the reviewed change. No previous
  review artifacts or private evidence were consulted.
- [x] Traced folder construction and lifecycle through
  `internal\testhelp\testacc.go:162-180`,
  `internal\services\mlmodel\resource_ml_model_test.go:240-299` and
  `internal\services\mlexperiment\resource_ml_experiment_test.go:214-274`.
  Each helper creates its two folder configurations once, retains both in both
  steps and switches the item between Terraform folder-ID references.
- [x] Inspected both `folder_fixtures_test.go` files under the ML service packages,
  lines 21-111. The lifecycle tests reuse the actual CRUD steps; the separate
  folder store and serialized transport are test-local. The HCL checks cover
  workspace selection, owned references, folder retention and fresh run names.
- [x] Cross-checked workflow locking at `.github\workflows\test.yml:23-27` and
  `.github\workflows\cleanup-script.yml:10-14`, including the retained test matrix
  at `.github\workflows\test.yml:313-344`. The constant workflow-level group
  matches, cancellation is disabled and the matrix is not given a competing
  per-job lock.
- [x] Checked `DEVELOPER.md:318-322,353-367` against
  `internal\testhelp\wellknown.go:20-60`,
  `tools\scripts\Set-WellKnown.ps1:2309-2312` and both workflows. The added text
  correctly distinguishes local setup from the Actions variable, documents
  pending-run replacement rather than FIFO, and excludes local/fork/manual
  operations from the GitHub lock.
- [x] Checked the Go lint contract in `.golangci.yml:13-97` and its CI invocation
  at `.github\workflows\test.yml:225-233`. Independently compared the new
  `runtime.Caller` assignment with the explicitly suppressed equivalent at
  `internal\testhelp\wellknown.go:34-37`; also checked the public dogsled
  documentation and revive `nested-structs` implementation.
- [x] Independently ran both complete ML unit packages through `task testunit`
  with `TF_ACC=0`: **14 tests passed, seven per package**. Then independently ran
  the workflow-isolation regression through the same task: **one test passed**.
- [x] Ran read-only `gofmt -l` on all five reviewed Go test files, read-only
  `actionlint` on both workflows, `go mod tidy -diff` and `git diff --check` on
  the reviewed source. All completed successfully without source changes.
- [ ] Did not run the complete configured `golangci-lint` suite: its executable
  was unavailable in the existing tools. The findings below are static,
  rule-backed findings, not fabricated local linter output. `gofmt`, unit tests
  and `actionlint` do not establish Go lint cleanliness.
- [ ] Did not run live acceptance tests, cloud operations, cleanup workflows,
  the Terraform/OpenTofu CI version matrix, live queue/cancellation checks or
  performance benchmarks; these are outside the authorized offline review.

## Issues

### Issue 1: New workflow regression violates the enabled dogsled rule

- **Severity**: Medium
- **File**: `internal\testhelp\folder_fixtures_test.go`
- **Line(s)**: 17
- **Description**: `_, filename, _, _ := runtime.Caller(0)` discards three return
  values. Dogsled is enabled by `.golangci.yml`'s `default: all`; it is neither
  disabled nor covered by the test-file exclusions, and its default allows at
  most two blank identifiers. The existing equivalent in
  `internal\testhelp\wellknown.go:35` already carries `//nolint:dogsled`, whereas
  the new assignment does not. Passing `task testunit` cannot detect this rule
  violation.
- **Risk**: The new test introduces a diagnostic in the configured CI Go lint
  check even though its behavioral assertions pass, preventing the change from
  satisfying the build gate.
- **Suggested Fix**: Retain and check the `ok` return from `runtime.Caller`,
  leaving only two blank identifiers, or use the existing narrowly scoped
  dogsled suppression with an appropriate explanation. Do not disable the
  linter globally. Re-run the targeted workflow unit test and the configured
  Go linter with fixing explicitly disabled.

### Issue 2: Inline concurrency struct violates revive's enabled nested-structs rule

- **Severity**: Medium
- **File**: `internal\testhelp\folder_fixtures_test.go`
- **Line(s)**: 19-24
- **Description**: The local `workflow` type contains an anonymous struct type
  in its `Concurrency` field. `.golangci.yml:38-64` enables all revive rules and
  does not disable `nested-structs`; the exclusions at lines 80-97 do not exempt
  this rule for tests or `internal\testhelp`. The rule visits struct fields and
  rejects precisely this nested struct declaration. This is separate from YAML
  decoding correctness, which the passing unit regression exercises.
- **Risk**: This adds another CI Go lint diagnostic in the same new test file.
  Fixing the blank identifiers alone would therefore leave the lint gate dirty.
- **Suggested Fix**: Introduce a small named concurrency configuration type and
  use that type for `workflow.Concurrency`, preserving the existing YAML tags
  and assertions. Keep the change local to the regression test, then re-run its
  unit test and the configured Go linter without autofix.

## Tests Actually Run and Reported Evidence

All test execution used the repository task runner, existing local tools,
project-local scratch/cache directories, offline Go module resolution and
`TF_ACC=0`. No raw test command, live acceptance entrypoint or tool installation
was used.

| Independent command/check | Actual result |
| --- | --- |
| `task testunit -- '.*' 'github.com/microsoft/terraform-provider-fabric/internal/services/ml...'` | Passed: 14 tests, seven in each ML package, including the four new ML fixture tests. |
| `task testunit -- FolderFixtures_WorkflowIsolation 'github.com/microsoft/terraform-provider-fabric/internal/testhelp'` | Passed: one workflow regression test. |
| `gofmt -l` on the two changed CRUD test files and three new fixture test files | Exit 0; no files reported. |
| `actionlint -no-color` on `.github\workflows\test.yml` and `.github\workflows\cleanup-script.yml` | Exit 0; no diagnostics. |
| `go mod tidy -diff` | Exit 0; no differences. Existing HCL/YAML versions are promoted to direct requirements, not upgraded. |
| `git diff --check` on reviewed source, excluding review artifacts | Exit 0; no whitespace errors. |

The initial task invocation selected a WSL launcher with no installed
distribution and executed **zero** tests. Prepending the already installed Git
Bash to the process PATH resolved that local harness issue; the successful runs
above completed through coverage generation.

The supplied scope reports earlier red regressions and passing validation.
The final passing unit results, module tidiness and workflow lint were
independently reproduced here. The earlier red state was not reconstructed
because this assignment prohibits source edits. No claim is made that the full
Go lint suite passed.

## Final-Stage Assessment and Limits

Beyond the two lint findings, no additional actionable performance,
observability, documentation or naming defect was found in the scoped change.
The extra folder operations are bounded and required for fixture ownership;
there is no changed production request path, retry policy or logging behavior.
Serializing the in-memory fake transport is confined to the new tests. The
broader workflow serialization is intentional, retains intra-run matrix
parallelism and has its throughput/queue limitations documented.

Both CRUD helpers preserve create/read, move/read and automatic Terraform
teardown while removing the stale static-folder dependency. The offline fake
lifecycle is not proof of live service acceptance or complete real-workspace
cleanup. The driver should fix the two local lint violations and obtain actual
configured Go lint evidence before treating this final-stage attempt as clean.
No source edits, commits, pushes, publications or external writes were made.

## Resolution Log

_Updated by the driving agent as findings are addressed._

### Issue 1

- **Status**: Open
- **What changed**: Nothing; this invocation is review-only.
- **Why**: Preserve the reviewed source and leave remediation to the driver.
- **How verified**: Static rule/configuration comparison and the existing
  `runtime.Caller` suppression pattern; the targeted behavioral unit test
  passed, but no lint fix was applied or verified.

### Issue 2

- **Status**: Open
- **What changed**: Nothing; this invocation is review-only.
- **Why**: A local named type can address the policy violation without changing
  workflow behavior or widening the source change.
- **How verified**: Static comparison of the anonymous struct with the enabled
  revive rule; YAML workflow behavior passed its unit test, but no lint fix was
  applied or verified.

## Driver Resolution Before Stage 6 Attempt 2

- **Issue 1 resolved:** the test now retains and checks `runtime.Caller`'s `ok`
  result instead of discarding three return values. Failure to locate the test
  file reports a clear test failure.
- **Issue 2 resolved:** a local named `concurrencyConfig` type replaces the
  anonymous nested struct without changing the YAML schema or assertions.
- The configured Go linter was initially unavailable. After that failed
  validation invocation, an official checksum-verified `golangci-lint` 2.14.0
  release was installed in isolated local tooling, not globally.
- Running the actual configured linter reproduced both findings and identified
  additional new-test-only style diagnostics: unnamed interface parameter,
  return/statement spacing, inline error assignment and the required external
  YAML field spelling. The three new test files were corrected. The external
  `cancel-in-progress` field keeps its exact spelling with a narrow, explained
  `tagliatelle` suppression; no global rule was weakened.
- Only the three new test files changed in this correction. Production code,
  acceptance builders, workflow behavior and existing tracked helpers were
  unchanged. Automatic spacing fixes were restricted to the reviewed packages,
  and before/after checks confirmed no other tracked source edits.
- **Validation:** configured `golangci-lint run --fix=false` on the three
  affected packages reports **0 issues**. Rebuilt through `task testunit`:
  both complete ML unit packages pass **14 tests**, and the workflow guard
  passes **1 test**. Both workflow actionlint checks, `go mod tidy -diff` and
  `git diff --check` also pass.
- These are source changes after attempt 1, not an infrastructure-only retry.
  Both slots must independently review the current candidate at the ordinary
  stage-6 attempt 2. No previous clean slot is carried across this correction.
