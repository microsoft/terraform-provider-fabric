## Review Summary

- **Round**: 6
- **Stage**: 6
- **Attempt**: 1
- **Theme**: Polish & hardening
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Reasoning effort**: xhigh
- **Artifact**: reviews/task-terraform-folder-fixtures-attempt-1-stage-attempt-1-review-6-grok-4.7.md
- **Issues Found**: 0
- **Verdict**: CLEAN

## Evidence Checklist

- [x] Inspected current source independently. No prior review artifacts were used as evidence.
- [x] Workflow lock is a shared constant with cancellation disabled, and the acceptance matrix job has no per-job lock: `.github/workflows/test.yml:23-27`, `.github/workflows/test.yml:310-329`, `.github/workflows/cleanup-script.yml:12-14`.
- [x] Developer docs match the fixture and lock contract: `DEVELOPER.md:318-321`, `DEVELOPER.md:353-367`.
- [x] Populated well-known environment data takes precedence over the setup file, and setup writes only the local fixture file: `internal/testhelp/wellknown.go:40-52`, `tools/scripts/Set-WellKnown.ps1:2312`, `tools/scripts/Remove-WorkspaceRSItems.ps1:246-259`.
- [x] CRUD steps own two folders, keep both across create and move, and reference folder IDs instead of static fixture IDs: `internal/testhelp/testacc.go:162-176`, `internal/services/mlexperiment/resource_ml_experiment_test.go:214-268`, `internal/services/mlmodel/resource_ml_model_test.go:240-296`.
- [x] Unit lifecycle and HCL regressions reject the stale fixture IDs and serialize fake HTTP: `internal/services/mlexperiment/folder_fixtures_test.go:21-111`, `internal/services/mlmodel/folder_fixtures_test.go:21-111`.
- [x] Workflow regression reads both workflow files and rejects expression groups, cancellation, and a per-matrix-job lock: `internal/testhelp/folder_fixtures_test.go:16-53`.
- [x] Test-only modules are direct requirements and `go mod tidy -diff` is clean: `go.mod:13`, `go.mod:28`.
- [x] `gofmt -d` reported no diff for the changed Go files. `actionlint` exited 0 on `.github/workflows/test.yml` and `.github/workflows/cleanup-script.yml`.
- [x] Ran `task testunit` with `TF_ACC=0` for both ML packages and the workflow regression. Both ML packages reported 7 passing unit tests. The workflow package reported 1 passing test. No raw `go test` was used.
- [ ] Full `./...` unit suite was not run. Affected packages were run instead; see limits.
- [ ] `golangci-lint` and Markdown autofix linters were not run. They are not in the local tool cache, and the repository lint tasks rewrite source.
- [ ] Live acceptance, cloud calls, and `TF_ACC=1` were not run.

## Clean Assessment

Stage 6 checks performance, observability, documentation accuracy, and naming on the current fixture and workflow changes. No product runtime, retry, or classification code is changed. The implementation matches the resolved contract, and the new names and comments are specific enough to maintain.

### Documentation

`DEVELOPER.md:318-321` is accurate. Setup writes `internal/testhelp/fixtures/.wellknown.json` and does not publish the Actions variable. A populated `FABRIC_TESTACC_WELLKNOWN` value is preferred by both `WellKnown()` (`internal/testhelp/wellknown.go:47`) and the cleanup script (`tools/scripts/Remove-WorkspaceRSItems.ps1:246-259`). Both Test and Cleanup workflows inject that Actions variable at workflow scope.

`DEVELOPER.md:353-355` matches the CRUD helpers. Each MLModel and MLExperiment acceptance case reads `WorkspaceRS`, creates two owned folders with `FolderResource`, and passes Terraform references. Both folder blocks are joined into the create step and the move step, so the previous folder is not dropped from configuration. Terraform destroys the owned resources after the last step. These cases no longer read `FolderRS1` or `FolderRS2`. Data-source acceptance cases still use their own well-known entities; the new sentence does not claim otherwise.

`DEVELOPER.md:357-367` matches the workflow files. The group is the constant `fabric-testacc-environment`, cancellation is false, and the lock is workflow-level. That holds the whole run, including matrix teardown, before the other workflow can enter the group. The one-pending replacement behavior is stated as not FIFO, and local, fork, other-repository, and manual cleanup are explicitly outside the lock.

### Naming

`fabric-testacc-environment` names the shared acceptance environment rather than a ref-specific queue. `experimentCRUDSteps` and `modelCRUDSteps` identify the reused create/move steps. `TestUnit_FolderFixtures_*` is discoverable by the unit-test filter. `folderFixtureTransport` states that the mutex is a fixture concern, not a provider transport. Comments at the separate folder-server assignment and the mutex explain why those test-only mechanisms exist.

### Performance

The constant group serializes runs that share the acceptance workspace, including newer commits, instead of cancelling an in-progress run. That is the intended cost of keeping teardown intact. Within one run, the matrix stays parallel: `.github/workflows/test.yml` has no job-level `concurrency` on the `test` job, and the regression test fails if one is added. A job-level lock would be the slower and incorrect alternative because a pending matrix leg can be replaced.

The mutex wraps each fake `Do` call only in the two unit lifecycle tests. Acceptance tests and provider runtime do not use it. The fake transport finishes store updates before `Do` returns, so the lock covers the in-memory critical section without a product hot path. Building a second fake server only isolates the folder store from the generic item conversion map. That allocation is test-only and is explained in source.

Promoting `github.com/hashicorp/hcl/v2` and `gopkg.in/yaml.v3` to direct requirements does not add new runtime modules. They were already in the module graph, and `go mod tidy -diff` produced no diff.

### Observability

Failure text in the new tests names the broken contract: stale static folder IDs, a missing owned folder, a destroyed previous folder, a cancelled in-progress run, or a per-matrix-job lock. Provider client construction already logs initialization failure. GitHub shows the constant concurrency group in the Actions UI, and `DEVELOPER.md` tells operators to re-run a superseded required check. No production log point is missing because no production control flow changed.

### CI lint and formatting

`actionlint` accepted both workflows. `gofmt -d` was clean on the changed Go files. New test files use the repository copyright header and SPDX line. The new developer prose is wrapped, has no new lists or untyped fences, and does not depend on the disabled line-length rule. Duplicate fixture tests are covered by the existing test-file duplication exclusion; the two copies differ only by ML resource name and are not drifting.

## Considered And Not Filed

- Collateral waiting of non-acceptance jobs in the Test workflow is a consequence of the required workflow-level constant group. A per-job or expression-based group would either cancel matrix legs or fail the new regression. `DEVELOPER.md:357-361` already describes workflow-level serialization. This is not a separate defect.
- The pre-existing Cleanup comment at `.github/workflows/cleanup-script.yml:10-11` says a queued run waits and is indented above the concurrency key. It does not contradict `cancel-in-progress: false`. The new developer text is the accurate operator contract, including one-pending replacement. Rewriting the untouched comment is not required for this change.
- Empty or whitespace-only `FABRIC_TESTACC_WELLKNOWN` handling differs slightly between Go lookup and the cleanup script. The new documentation claims precedence only for a populated value, which both paths implement. The empty-value edge is pre-existing and outside this claim.
- No changelog fragment is present. The changelog workflow fails a pull request unless a fragment exists or the skip-changelog label is applied. A fragment requires a real public issue number, which this review must not invent. That is a merge-time process follow-up, not a source defect in the fixture implementation.

## Tests Run Versus Reported

Reported claim checked here: the five new tests pass, and each full ML unit package passes 7 tests. No live acceptance was claimed as executed by this review.

Actually run, with `TF_ACC=0`, through `task testunit` only:

- `task testunit -- . ./internal/services/mlexperiment`: passed, 7 tests, about 34s. This package has exactly 7 `TestUnit_` functions, including `TestUnit_FolderFixtures_ExperimentLifecycle` and `TestUnit_FolderFixtures_ExperimentCRUD`.
- `task testunit -- . ./internal/services/mlmodel`: passed, 7 tests, about 33s. This package has exactly 7 `TestUnit_` functions, including the two model folder-fixture tests.
- `task testunit -- FolderFixtures ./internal/testhelp`: passed, 1 test, about 13s. That is `TestUnit_FolderFixtures_WorkflowIsolation`.

The first argument `.` makes the task run pattern `^TestUnit_.`, so the ML commands execute the full unit set in each package rather than only the new tests. Acceptance tests were not selected. The reported 7-and-7 unit result matches this run. The five new tests are the four ML fixture tests included in those 7-test counts plus the one workflow test.

Also run as checks, not as the unit task:

- `actionlint` on `.github/workflows/test.yml` and `.github/workflows/cleanup-script.yml`: exit 0.
- `gofmt -d` on the five changed Go files: no diff.
- `go mod tidy -diff`: exit 0, no module diff.

## Limits

- `TF_ACC` remained 0. No live acceptance, cloud API, cleanup, or credential use was executed.
- The full unit suite under `./...` was not run. Only the two changed ML packages and the workflow regression package were executed.
- `golangci-lint`, Markdown lint, and other autofix linters were not run. They are absent from the local tool cache, and the repository tasks apply fixes to source. Formatter confidence is limited to `gofmt -d` plus inspection against the checked-in lint config.
- `gotestsum` was invoked in fail-only package format, so individual passing names were confirmed by the package test inventory and the passing counts, not by a per-test success log.
- This review did not modify source, commit, push, or write outside this artifact.
