---
producer: "review-code"
schema_version: 1
tier: "Medium"
contract: "direct"
execution_backend: "copilot-cli"
attempt: 5
stage_attempt: 4
round: 3
timestamp: "2026-10-31T07:35:34Z"
outcome: "issues_found"
reviewer_model: "grok-4.7"
input_sha256: "002669bc4a58e3bdae6298fa6e4d7c15b5a5dc3c69799e93bc6374a8913cbcd7"
---

## Review Summary

- **Round**: 3
- **Theme**: Edge cases & robustness
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Artifact**: reviews/task-1119-attempt-5-stage-attempt-4-review-3-grok-4.7.md
- **Issues Found**: 5
- **Verdict**: ISSUES_FOUND

## Evidence Checklist

- [x] Compared `.github/workflows/test.yml` (concurrency at lines 21–27, `changes` filters at 48–57, recoverable cleanup at 60–83, matrix at 338–368) with `.github/workflows/cleanup-script.yml` (group at 12–14, full `-Force` cleanup at 42–45).
- [x] Checked `internal/testhelp/folder_fixtures_test.go` wiring against `context-only-public-dependencies/yaml.v3@v3.0.1/decode.go` (`Node.Decode`, zero-node `null`, string decode via raw `n.Value`). Current `src` output, selector list, and test/auth/cleanup `if`/`needs` match the workflows.
- [x] Traced test-owned folders from `internal/testhelp/testacc.go` `FolderResource` through both `folder_fixtures_test.go` files, fake item/folder handlers, and `tools/scripts/Remove-WorkspaceRSItems.ps1` lines 19 and 324–388.
- [x] Compared `DEVELOPER.md` lines 353–370 with jobs that have no `timeout-minutes` (`changes` 40, `checkdocs` 295, `test-summary` 516, `check-test-matrix` 597) and with `on.pull_request` at `.github/workflows/test.yml` lines 8–13.
- [ ] Tests, actionlint, and Picomatch were not executed; this was a read-only source review. No live scheduling or recovery was established.

## Issues

### Issue 1: Cancelled or timed-out runs leave active test-owned folders

- **Severity**: Medium
- **File**: `.github/workflows/test.yml`
- **Line(s)**: 64–83, 338–347
- **Description**: ML CRUD steps now create active WorkspaceRS folders and destroy them only if Terraform finishes. A job timeout (`timeout-minutes: 100`), manual cancel, or runner kill skips that destroy. The in-workflow cleanup is `-RecoverableItemsOnly` with `continue-on-error: true`, and `Remove-WorkspaceRSItems.ps1` leaves active items and folders untouched in that mode. The only full delete is scheduled Cleanup (`-Force`), which shares `0 2 * * *` with Test and can be dropped when a newer run replaces the single pending slot. `DEVELOPER.md` says to cancel a stuck run and to re-dispatch Cleanup only if pending Cleanup was superseded. It does not say cancellation itself requires a full cleanup, and it does not keep the deferred workspace/stable-identity gap visible.
- **Risk**: Active test folders and ML items accumulate in the shared workspace. Later acceptance runs and a delayed full cleanup can then delete or collide with in-use resources. Quota pressure is unbounded if nightly Cleanup keeps losing the pending slot.
- **Suggested Fix**: Document only; do not change cleanup scope or add stable identities. State that cancel/timeout can skip destroy, that recoverable pre-cleanup will not reclaim those active objects, and that an owner must `workflow_dispatch` full Cleanup after a cancelled or timed-out Test run. Note that workspace/stable-identity lifecycle automation stays deferred.

### Issue 2: Shared lock lifetime is not bounded by job timeouts

- **Severity**: Medium
- **File**: `.github/workflows/test.yml`
- **Line(s)**: 21–27, 40–42, 295–297, 516–522, 597–601
- **Description**: Workflow-level `fabric-testacc-environment` is held until the whole run exits, including approval waits. Job `timeout-minutes` starts only after a job is running, so it does not bound queue or environment-approval time. `changes`, `checkdocs`, `test-summary`, and `check-test-matrix` set no timeout, so they use the hosted-runner default (360 minutes) while still holding the environment lock. `checkdocs` also has no `src` gate, so a docs-only run holds the lock. `DEVELOPER.md` mentions approval waits and stuck runs, but not this timeout gap. Adding bounds is scope-deferred; the advisory is not visible.
- **Risk**: A hung docs/summary job, or an unapproved fork run, blocks Cleanup and every other Test run for hours. The single pending slot can be replaced during that wait, so scheduled Cleanup never starts.
- **Suggested Fix**: Do not add timeouts in this change. In `DEVELOPER.md`, state that job timeouts do not cover approval or queue time, name the jobs with no `timeout-minutes`, and mark explicit lock bounds as deferred.

### Issue 3: pull_request runs can opt out of the shared lock

- **Severity**: Medium
- **File**: `.github/workflows/test.yml`
- **Line(s)**: 8–13, 21–27
- **Description**: `on.pull_request` executes the merge commit’s workflow file, so a fork or same-repo PR can remove or rename `concurrency`. That run does not join `fabric-testacc-environment`. Environment approval, if granted, still supplies development credentials. `TestUnit_FolderFixtures_WorkflowIsolation` fails the check only after the run has already started; it cannot put the in-flight run back into the group. `DEVELOPER.md` lines 366–368 say upstream fork-PR Test runs join the group. That is true only when the executed file still contains it. Fork-owned workflows are a separate, already documented gap.
- **Risk**: An approved PR whose workflow edit drops the group can mutate WorkspaceRS while scheduled Cleanup or another Test run holds, or believes it holds, the lock. Shared fixtures and new test-owned folders can be deleted mid-apply.
- **Suggested Fix**: Correct the doc sentence. State that group membership is whatever the executed workflow file declares, that `pull_request` uses the PR merge workflow, and that approving `development` for a run that fails the isolation check can still race. Do not switch the trigger to `pull_request_target`.

### Issue 4: Folder-fixture HCL checks panic instead of failing

- **Severity**: Low
- **File**: `internal/services/mlexperiment/folder_fixtures_test.go`
- **Line(s)**: 72–95
- **Description**: The same pattern is in `internal/services/mlmodel/folder_fixtures_test.go`. After a successful parse, the test indexes `block.Labels[0]` and `block.Labels[1]` without a length check, calls `.Expr` on `Attributes["workspace_id"]` and `Attributes["folder_id"]` without a nil check, and type-asserts `refs[0][1]` and `refs[0][2]` to `hcl.TraverseAttr`. A short label list, a missing attribute, or an index traversal panics before the intended `t.Fatal`.
- **Risk**: A generator regression becomes an opaque panic. Rerun tooling may treat it differently from an assertion failure, and the diagnostic that the folder reference regressed is lost.
- **Suggested Fix**: Require two labels, require the attributes to be non-nil, and use comma-ok traversal asserts with `t.Fatalf` that includes the block labels.

### Issue 5: Isolation guard will not see a new development job

- **Severity**: Low
- **File**: `internal/testhelp/folder_fixtures_test.go`
- **Line(s)**: 34–40, 112–120, 139–141
- **Description**: Job-level concurrency is rejected only for the hardcoded names `cleanup-workspace`, `test-auth-spn`, and `test`, plus standalone cleanup. `workflowJob` does not decode `environment`. Those three jobs are the only active `development` jobs today, so the current file is covered. A new job with `environment.name: development` and its own `concurrency` would still pass.
- **Risk**: Job-level concurrency overrides the workflow group for that job. A later development job could run beside Cleanup, which is the collision this guard is meant to keep from returning. GitHub keeps only one pending run per group, so a per-job group can also cancel queued matrix legs.
- **Suggested Fix**: Decode `environment.name` and reject job-level `concurrency` on every `development` job, while keeping the existing `if`/`needs` checks for the three named jobs and the standalone cleanup job.

## Resolution Log

_Updated by the driving agent as findings are addressed._

### Issue 1

- **Status**: Open
- **What changed**: pending
- **Why**: pending
- **How verified**: pending

### Issue 2

- **Status**: Open
- **What changed**: pending
- **Why**: pending
- **How verified**: pending

### Issue 3

- **Status**: Open
- **What changed**: pending
- **Why**: pending
- **How verified**: pending

### Issue 4

- **Status**: Open
- **What changed**: pending
- **Why**: pending
- **How verified**: pending

### Issue 5

- **Status**: Open
- **What changed**: pending
- **Why**: pending
- **How verified**: pending

## Publication Note

The machine-local Artifact summary path was normalized and generated review Markdown formatting was applied. Original response and finalized receipt remain immutable outside this source tree;findings,verdict and input provenance are unchanged.
