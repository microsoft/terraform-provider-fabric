---
producer: "review-code"
schema_version: 1
tier: "Medium"
contract: "direct"
execution_backend: "copilot-cli"
attempt: 5
stage_attempt: 7
round: 3
timestamp: "2026-10-31T09:11:45Z"
outcome: "issues_found"
reviewer_model: "grok-4.7"
input_sha256: "ca9e1a397a1d2a210dd47f51d49ae1fbad92e4206743953de831f747ab1cd327"
---

## Review Summary
- **Round**: 3
- **Theme**: Edge cases & robustness
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Artifact**: reviews/task-1119-attempt-5-stage-attempt-7-review-3-grok-4.7.md
- **Issues Found**: 3
- **Verdict**: ISSUES_FOUND

## Evidence Checklist
- [x] Compared `DEVELOPER.md` lines 361–401 with the seven required recovery corrections and with workflow-level concurrency in `.github/workflows/test.yml` lines 23–27 and `.github/workflows/cleanup-script.yml` lines 10–14.
- [x] Checked the Test job graph: environment jobs at `test.yml` lines 65, 91, and 344; non-environment jobs `changes` (40), `checkbuild` (215), `checkdocs` (295), `test-summary` (516), and `check-test-matrix` (597); `continue-on-error: true` at line 64; timeouts 30 and 100 at lines 68, 221, and 347.
- [x] Checked `context-only-public-docs/verified-github-workflow-facts.json`: one running run plus one replaceable pending run; environment approval is job-level, not the concurrency declaration; selected ref configuration must be reviewed before `workflow_dispatch`.
- [x] Confirmed standalone Cleanup invokes `Remove-WorkspaceRSItems.ps1 -Force` (`cleanup-script.yml` line 45) while pre-clean adds `-RecoverableItemsOnly` (`test.yml` line 82); script lines 4 and 17–19 hard-delete active items unless that switch is set.
- [x] Confirmed unchanged qualifications remain: merge-ref and separate approval layers (`DEVELOPER.md` lines 369–376), 30/100 bounds (line 401), shared `WorkspaceRS`/`FolderRS1`/`FolderRS2` (lines 402–403), and owner coordination (line 382) that does not replace the idle check (lines 364–367).
- [ ] Live Actions runs and unit tests were not executed; this pass is source and checked-in facts only.

## Issues

### Issue 1: Reject or cancel is treated as freeing the group slot
- **Severity**: High
- **File**: DEVELOPER.md
- **Line(s)**: 384–390
- **Description**: Step 1 says to free the slot without starting the job by rejecting the deployment or cancelling the workflow run. Environment approval is job-level, not the repository concurrency lock. Rejecting one deployment does not stop sibling jobs that have no environment (`changes`, `checkdocs`, `checkbuild`, and later `test-summary` / `check-test-matrix`), so the workflow run remains the running group occupant. Cancelling that run completes only that run. The group still allows one pending run; that pending run is then promoted into the running slot. A pending Cleanup is the documented 02:00 case and will proceed to `-Force` once promoted. Step 2 tells the operator to confirm idle, but it has no branch for a group that is still occupied, and it cannot undo a promoted Cleanup that has already started.
- **Risk**: An operator following the written cancel/reject step can start standalone `-Force` deletion, or overlap remaining Test jobs, while believing the slot was released.
- **Suggested Fix**: State that rejection does not start that job and does not release the group. Before cancelling the running run, cancel any pending group run that must not start. After either action, wait until the acted-on run has exited and repeat until the group has no running run and no pending run. Do not call either action a slot release by itself.

### Issue 2: "No Cleanup action" leaves a queued Cleanup able to run
- **Severity**: High
- **File**: DEVELOPER.md
- **Line(s)**: 395–396
- **Description**: If the owner decides Cleanup is not needed, step 3 says no Cleanup action is required. That covers not dispatching a new run. It does not cancel a Cleanup run already pending behind the in-progress run. When the running run exits, that pending run becomes the running run. `cleanup-script.yml` only names the `development` environment; it does not declare required reviewers, and the checked-in facts do not say this environment always waits. The pending run can therefore enter `-Force` without a later approval prompt. The needed path says to let Cleanup finish, so the not-needed path must explicitly stop an already queued run.
- **Risk**: The owner’s decision not to delete active items and non-preserved folders is not enforced. Scheduled or previously queued Cleanup still hard-deletes once it acquires the single running slot.
- **Suggested Fix**: Keep “do not dispatch Cleanup” for the not-needed path, and add: cancel any pending or approval-waiting Cleanup run in the group before the current run exits. Recheck idle before re-running Test checks. Do not leave a not-needed Cleanup queued.

### Issue 3: Selected-ref gate checks only concurrency keys
- **Severity**: Medium
- **File**: DEVELOPER.md
- **Line(s)**: 391–394
- **Description**: Step 2 requires the selected Cleanup ref to keep the same group and `cancel-in-progress: false`. That matches the cancel/lock hazard, but the checked-in manual-run fact says to review the selected ref’s workflow configuration before invoking it. A ref can keep those two concurrency settings and still change the job command, drop preservation, or point `environment` at a different name. Step 3 then says to let that Cleanup finish.
- **Risk**: An operator who treats the concurrency pair as a sufficient pre-dispatch check can start a same-lock run whose executed job is not the intended `Remove-WorkspaceRSItems.ps1 -Force` cleanup, with no cancel of active work to reveal the mistake.
- **Suggested Fix**: Before `workflow_dispatch`, review the whole selected Cleanup workflow, not only the concurrency keys. Require the same group, `cancel-in-progress: false`, the `development` environment, and the same `-Force` script invocation. If any of those differ, do not dispatch and do not let that ref’s job finish as if it were the known cleanup.

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
## Publication Note

Only the machine-local Artifact path was normalized. Original response,findings,verdict and finalized input/response receipts are retained unchanged.
