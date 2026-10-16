# Stage 6 review

Independent current-source review. Prior review artifacts were not read and no earlier clean verdict was carried forward.

## Review Summary
- **Stage**: 6
- **Attempt**: 2
- **Round**: 6
- **Theme**: Polish & hardening
- **Mode**: resolved-plan
- **Model**: grok-4.7 (xhigh)
- **Artifact**: reviews/task-terraform-folder-fixtures-attempt-1-stage-attempt-2-review-6-grok-4.7.md
- **Issues Found**: 0
- **Verdict**: CLEAN

## Evidence Checklist
- [x] Full final scope, not only the three lint-touched tests: `.github/workflows/test.yml`, `DEVELOPER.md`, `go.mod`, `internal/services/mlmodel/resource_ml_model_test.go`, `internal/services/mlexperiment/resource_ml_experiment_test.go`, plus `internal/services/mlmodel/folder_fixtures_test.go`, `internal/services/mlexperiment/folder_fixtures_test.go`, and `internal/testhelp/folder_fixtures_test.go`. No product service files are in the diff.
- [x] Shared lock: `test.yml` lines 23-27 set group `fabric-testacc-environment` and `cancel-in-progress: false`, with a comment that a per-job lock would drop matrix legs. `cleanup-script.yml` lines 12-14 use the same group and `cancel-in-progress: false`. The `test` job (around lines 309-316) still uses environment `development` and has no job-level `concurrency` key. Other workflows do not use that environment name (`release` uses `release`; pages uses `github-pages`).
- [x] Docs match helpers: `DEVELOPER.md` lines 318-322 and 353-367. Setup writes `internal/testhelp/fixtures/.wellknown.json` (`tools/scripts/Set-WellKnown.ps1` line 2312) and does not publish the Actions variable. `WellKnown()` (`internal/testhelp/wellknown.go` lines 40-52) returns a set `FABRIC_TESTACC_WELLKNOWN` value before the file. Cleanup treats whitespace as unset and logs which source it used (`tools/scripts/Remove-WorkspaceRSItems.ps1` lines 246-259).
- [x] CRUD helpers: `FolderResource` (`internal/testhelp/testacc.go` line 162) and `RefByFQN` (`internal/testhelp/utils.go` line 190). Both acceptance step builders keep both folder configs on create and move and check the folder id with `TestCheckResourceAttrPair`.
- [x] Naming and polish in the three new tests: `runtime.Caller` ok check (`internal/testhelp/folder_fixtures_test.go` lines 17-19); named `concurrencyConfig` and a one-field `tagliatelle` suppression with an external-name reason (lines 23-26); named interface parameter `request` on `folderFixtureTransport` in both service tests. No `TODO`, `FIXME`, or `HACK` in the changed files. LF only, no trailing whitespace.
- [x] Actually run, pinned golangci-lint 2.14.0 `--fix=false` on `internal/testhelp`, `internal/services/mlmodel`, and `internal/services/mlexperiment`: version 2.14.0, `0 issues.`, exit 0. Config deprecation warnings only; they are not findings in this diff.
- [x] Actually run: actionlint on `test.yml` and `cleanup-script.yml` exit 0; `go mod tidy -diff` exit 0; `git diff --check` exit 0.
- [x] Actually run unit tests with `TF_ACC=0` through the repository task `testunit` (the defined unit task; it sets `TF_ACC=0`). `task test:unit` is not defined. Result: mlmodel 7 passed, mlexperiment 7 passed, testhelp 1 passed (15 total), matching the 15 `TestUnit_` functions in those packages. No raw `go test`. No live acceptance.
- [ ] Live acceptance and cloud cleanup were not run. This review is prohibited from live `TF_ACC` or cloud cleanup. Acceptance behavior is covered here only by fake lifecycle tests and HCL regressions.

## Clean reasoning

This pass reviewed the whole candidate, including the workflow lock, developer docs, module promotion, acceptance step builders, and the three new tests. The lint-only edits (caller ok check, named concurrency type, narrow YAML tag suppression, named interface parameter, and formatter spacing) are present and are not the only thing checked.

Performance. The constant workflow group and `cancel-in-progress: false` deliberately serialize Test and Cleanup so an active matrix, including teardown, is not cancelled. That is the required protection, and the docs say so: the group is workflow-level, replacement is not FIFO, and local, fork, other-repo, and manual cleanup stay outside the lock. The matrix `strategy` remains inside the workflow, and the `test` job has no second concurrency key. Jobs that do not declare the development environment still occupy the group for the life of a Test run. That is inherent to a workflow-level group and is already stated as workflow-level, so it is not a separate defect. The fake HTTP mutex and the extra in-memory folder server are test-only. `NewFakeServer` registers handlers in memory and does not start a listener, so the discarded second server is not a connection leak. No product request path changed.

Observability. The workflow comment names the shared group and the per-job hazard. `DEVELOPER.md` tells an operator why a newer run waits, why a superseded check must be re-run, and which paths the lock does not cover. Cleanup already logs whether it used the environment variable or the file. New tests fail with specific messages (static id, missing folder, cancelled run, per-job lock) rather than a bare success check. No new product log or metric is required because product runtime is unchanged. There is no unresolved `TODO`/`FIXME`/`HACK` in the changed files.

Documentation. The new well-known paragraph matches the setup writer and the populated-variable precedence. It does not claim an empty variable falls back; that empty-string `LookupEnv` edge is pre-existing and outside this diff. The CRUD paragraph matches the step builders: owned folders in `WorkspaceRS`, Terraform references, both folders retained across the move, and destroy of resources still present in the final config. It does not claim `WorkspaceRS` itself is optional. The concurrency paragraph matches the two workflow files, including one pending run rather than a FIFO queue.

Naming. `fabric-testacc-environment` is the same literal in both workflows and is not an expression. `concurrencyConfig`, `experimentCRUDSteps`, `modelCRUDSteps`, and `folderFixtureTransport` say what they are. The YAML suppression is one external field, with the reason on that line. The interface parameter is `request`, not a blank name. `runtime.Caller` failure is a fatal with a clear message. Direct `hcl/v2` and `yaml.v3` requirements are the same versions previously indirect, which matches the new test imports. `go.sum` did not need a content change; tidy diff was empty.

Considered and not filed. A docs-only Test run still holds the shared group until the workflow finishes. Splitting that off would fight the required workflow-level lock, and the docs already call the unit a workflow. The duplicated transport type stays in the two external test packages so the fake client can be injected there; the name and comment are clear. CI checkbuild still asks for a latest linter, which is pre-existing and was not changed here. This review's gate is the pinned 2.14.0 run, which reported 0 issues.

## Validation

Reported: pinned golangci-lint 2.14.0 has 0 issues; Task 15 tests, actionlint, tidy diff, and whitespace pass.

Actually run this attempt:
- golangci-lint 2.14.0 `--fix=false` on `internal/testhelp`, `internal/services/mlmodel`, and `internal/services/mlexperiment`: confirmed version 2.14.0, output `0 issues.`, exit 0. Two pre-existing config warnings (deprecated formatter extra-rules key, deprecated gomodguard name) did not produce issues.
- actionlint on `.github/workflows/test.yml` and `.github/workflows/cleanup-script.yml`: exit 0, no findings.
- `go mod tidy -diff`: exit 0, no diff. `git diff --check`: exit 0. Changed and new files are LF with no trailing whitespace.
- `task test:unit` is not a task in this repository. The unit entrypoint is `task testunit`, which sets `TF_ACC=0`; `TF_ACC=0` was also set in the environment. First call failed before any test ran because `bash` resolved to a stub with no distribution (0 tests parsed). Rerun with a working bash first on PATH: mlmodel 7 passed, mlexperiment 7 passed, testhelp 1 passed. Those 15 names are the `TestUnit_` functions in the three packages, including both lifecycle tests, both HCL regressions, and `TestUnit_FolderFixtures_WorkflowIsolation`.
- The repository `golangci-lint` task was not used because it passes `--fix` and would edit source.

Residuals: no live acceptance or cloud cleanup was executed, by constraint. CI may still run a newer linter than 2.14.0 because checkbuild uses `version: latest`; that setting is unchanged. Generated coverage files from the unit task are ignored and are not part of this diff. No source file was edited by this review.

Clean review round: zero issues found.
