---
producer: "review-code"
schema_version: 1
tier: "Medium"
contract: "direct"
execution_backend: "copilot-cli"
attempt: 5
stage_attempt: 5
round: 3
timestamp: "2026-10-31T08:27:34Z"
outcome: "issues_found"
reviewer_model: "grok-4.7"
input_sha256: "43ff42cddecf0d583dc529d2321f15ed8fe9cc7537f5b4df272cdea91f56b1b1"
---

## Review Summary
- **Round**: 3
- **Theme**: Edge cases & robustness
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Artifact**: reviews/task-1119-attempt-5-stage-attempt-5-review-3-grok-4.7.md
- **Issues Found**: 2
- **Verdict**: ISSUES_FOUND

## Evidence Checklist
- [x] Compared `DEVELOPER.md` lines 357-388 with `.github/workflows/test.yml` (concurrency lines 23-27, triggers lines 9-21, cron line 19) and `.github/workflows/cleanup-script.yml` (concurrency lines 10-14, cron line 7, `workflow_dispatch` line 8). Group, `cancel-in-progress: false`, and 02:00 UTC schedules match. Required qualifications for the executed group, merge-ref edits, workflow vs job approval, and follow-up isolation are present and match `context-only-public-docs/verified-github-workflow-facts.json`.
- [x] Checked slot duration against jobs that keep a run in progress after or without test execution: `checkdocs` (line 295, no `if`, no `timeout-minutes`), `test-summary` (line 516, `if: always()`, no timeout), `check-test-matrix` (line 597, `if: always()`), and environment-gated jobs at `test.yml` lines 65, 91, and 344.
- [x] Checked pre-clean vs destroy vs Cleanup against `test.yml` lines 60-68 (`continue-on-error`, `timeout-minutes: 30`, `-RecoverableItemsOnly`), matrix `timeout-minutes: 100` at lines 347, 429, and 446, `cleanup-script.yml` line 29 and `-Force`, and `tools/scripts/Remove-WorkspaceRSItems.ps1` lines 14 and 373-380. Cleanup-job 30-minute and matrix 100-minute bounds match; `test-auth-spn` (line 85) has no job timeout.
- [x] Checked shared-identity claims against `test.yml` matrix `strategy` (line 354, no `max-parallel`) and ML folder creation via `testhelp.FolderResource` in `internal/services/mlmodel/resource_ml_model_test.go` and `internal/services/mlexperiment/resource_ml_experiment_test.go`, versus non-ML `FolderRS1`/`FolderRS2` use. YAML comments at `test.yml` lines 24-25 and `cleanup-script.yml` lines 11-12 match the executed group and do not change settings.
- [ ] Live Actions scheduling, approval, or recovery was not run. This review used source and the in-checkout verified facts only.

## Issues

### Issue 1: Cleanup redispatch "when no tests are active" can cancel a pending required run
- **Severity**: Medium
- **File**: DEVELOPER.md
- **Line(s)**: 361-364, 379-380
- **Description**: The recovery text says to re-dispatch standalone Cleanup with `workflow_dispatch` when no tests are active, and calls that safe. Cleanup joins `fabric-testacc-environment` with `cancel-in-progress: false` (`.github/workflows/cleanup-script.yml` lines 8-14). That group allows one in-progress run and one replaceable pending run. "No tests are active" is not "the group is idle." The workflow slot continues through docs/build work (`checkdocs`, line 295), `test-summary` (line 516), `check-test-matrix` (line 597), and environment-approval waits. A later Test can already occupy the pending slot while no test job is executing. A Cleanup dispatch is a newer group run and replaces that pending run.
- **Risk**: An owner following the documented recovery can cancel a queued required Test, including `check-test-matrix`, which treats `cancelled` as failure (`.github/workflows/test.yml` lines 597-608). That repeats the supersession the previous paragraph says to repair, and it can delay both cleanup and the required check.
- **Suggested Fix**: State that `workflow_dispatch` Cleanup is itself a newer group run and replaces the single pending slot. Dispatch only when the group has no in-progress run and no pending run, not merely when no test job is executing. If a dispatch supersedes a required check, re-run that check after the group is idle. Do not change the group, trigger, or `cancel-in-progress`.

### Issue 2: Named recoveries do not release a stuck environment-approval slot
- **Severity**: Medium
- **File**: DEVELOPER.md
- **Line(s)**: 357-367, 379-380, 383-384
- **Description**: The text correctly says the workflow-level slot includes environment-approval waits and that the 30-minute cleanup and 100-minute matrix job limits do not bound those waits. In-progress runs that still have `cancel-in-progress: false` are not cancelled by a later push, schedule, or dispatch; only an incoming executed file that sets `cancel-in-progress: true` is described as able to cancel an active run. The named recoveries are to re-run a superseded check and to re-dispatch Cleanup. Neither ends an in-progress run waiting on the `development` environment (`test.yml` lines 65, 91, and 344; `cleanup-script.yml` lines 26-29). A Cleanup dispatch while that wait holds the running slot only queues, or replaces the pending slot.
- **Risk**: One shared running slot can stay occupied with no job timeout, blocking Test and Cleanup. Operators who wait for the documented job limits, or who apply the documented redispatch, will not release it. Redispatch during the wait can also cancel a pending required run (Issue 1).
- **Suggested Fix**: Add that, while the executed workflows keep `cancel-in-progress: false`, a later run does not clear an in-progress environment-approval wait. The environment owner must approve, reject, or cancel that run to release the slot. Do not treat job timeout or Cleanup redispatch as that release. Do not change timeouts, environments, or concurrency settings.

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
## Publication Note

Machine-local artifact path normalized only;original response,verdict,findings and input provenance remain retained unchanged in the finalized run.
