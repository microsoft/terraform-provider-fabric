---
producer: "review-code"
schema_version: 1
tier: "Medium"
contract: "direct"
execution_backend: "copilot-cli"
attempt: 5
stage_attempt: 9
round: 3
timestamp: "2026-10-31T10:00:31Z"
outcome: "issues_found"
reviewer_model: "grok-4.7"
input_sha256: "194801dc887fdf6225f70048e77933f2a7edb3c4b3d4f3f2034d7b2c847957c1"
---

## Review Summary
- **Round**: 3
- **Theme**: Edge cases & robustness
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Artifact**: reviews/task-1119-attempt-5-stage-attempt-9-review-3-grok-4.7.md
- **Issues Found**: 4
- **Verdict**: ISSUES_FOUND

## Evidence Checklist
- [x] Recovery delta is `DEVELOPER.md` lines 383–414 only. Group `fabric-testacc-environment`, `cancel-in-progress: false`, merge-ref, and workflow-vs-environment approval text at lines 357–376 is unchanged and matches `.github/workflows/test.yml` lines 26–27, `.github/workflows/cleanup-script.yml` lines 10–14, and `context-only-public-docs/verified-github-workflow-facts.json` lines 9, 17, 32, and 39.
- [x] Standalone Cleanup is one workflow-level job, `development`, 30-minute timeout, checkout with no `ref`, and `Remove-WorkspaceRSItems.ps1 -Force` (`.github/workflows/cleanup-script.yml` lines 6–8, 23–29, 42–45). Test pre-clean is a different job: `continue-on-error: true` and `-Force -RecoverableItemsOnly` (`.github/workflows/test.yml` lines 60–64 and 83). Facts JSON lines 48–52 match those two modes.
- [x] Script success is not delete success. `exit 1` happens only for `ERROR` with default `Stop` (`tools/scripts/Remove-WorkspaceRSItems.ps1` lines 74–75). Item, recoverable, and folder delete failures pass `-Stop $false` (lines 566, 578, 611, 623, 651, 656). Zero deletes and partial deletes also use `-Stop $false` (lines 684–691), then main logs completion (line 769). Folder fetch failure (line 338) and standalone recoverable fetch failure (`-Stop $RecoverableItemsOnly`, line 367) continue, and a zero failed-count still logs full success (line 679).
- [x] Slot hold is until workflow exit (`DEVELOPER.md` line 359). Test siblings do not all need Cleanup: `checkbuild` needs only `changes` (`.github/workflows/test.yml` lines 215–219); `checkdocs` has no `needs` (line 295). Cleanup’s only job is `cleanup-workspace`.
- [ ] No live Actions run, test, or command was executed. Conclusions are from checkout source and `context-only-public-docs/verified-github-workflow-facts.json` only.

## Issues

### Issue 1: A green Cleanup conclusion can hide partial or skipped hard-deletes
- **Severity**: High
- **File**: DEVELOPER.md
- **Line(s)**: 408–412
- **Description**: Step 4 tells the operator to check the run conclusion and lists only failure, timeout, runner loss, and cancellation as partial-change causes. Standalone Cleanup’s only step is the script (`.github/workflows/cleanup-script.yml` lines 42–45) and does not inspect the deletion summary. Per-item delete errors do not `exit 1`. A total delete failure is logged as `ERROR` with `-Stop $false` and still falls through to “cleanup process completed”. Worse, a failed folder fetch or, for `-Force`, a failed recoverable-item fetch is swallowed, those sets become empty, and if the attempted deletes then report zero failures the script logs “All items, recoverable items, and folders deleted successfully!” (line 679) even though folders or recoverable items were never enumerated. The new “did not succeed” gate, keyed off conclusion, will not trip.
- **Risk**: The operator treats a green run as successful Cleanup, skips a real inventory, and starts Test on `WorkspaceRS` after some active items were hard-deleted and others, folders, or recoverable items remain. Idle-is-not-success does not correct that, because the unsafe signal here is success.
- **Suggested Fix**: State that conclusion `success` and the “deleted successfully” summary are inconclusive. After every standalone Cleanup exit, require the log to show folder and recoverable fetches succeeded and all three failed counts are zero, plus a fresh workspace inventory against the pre-run assessment. Any fetch warning, missing summary, non-zero failed count, or inventory mismatch is unsuccessful Cleanup: stop normal recovery, do not retry automatically, and do not start Test.

### Issue 2: Step 3 denies the slot release that promotes pending work
- **Severity**: Medium
- **File**: DEVELOPER.md
- **Line(s)**: 383–384, 388–390, 401–405
- **Description**: The intro says cancellation or rejection does not atomically clear the group and that queued runs can be promoted. Step 3 then says rejecting standalone Cleanup’s only job can end the run and promote pending work, and that cancellation stops only the selected run, but also that neither action “releases the group by itself”. Those claims conflict. Concurrency is workflow-level and the slot lasts until that run exits (`DEVELOPER.md` line 359; `.github/workflows/cleanup-script.yml` lines 10–14 and 24). Ending the in-progress run, including by cancelling it or by rejecting Cleanup’s only job, frees the single running slot. The facts file (line 9) allows one running run and one replaceable pending run, so that pending run is promoted. “Without another approval prompt” (line 390) also overstates the handoff: facts line 32 say environment approval is what permits the job to proceed. Promotion is not itself an approval, but the promoted job can still wait on its own `development` gate. That gate is not a concurrency lock (lines 372–373) and must not be described as absent.
- **Risk**: An operator cancels the active run, or rejects Cleanup’s only waiting job, while leaving a pending `-Force` Cleanup queued, believing the group stays held. The pending run starts immediately. They may also ignore a real approval wait on the promoted run, or assume hard-delete has already begun.
- **Suggested Fix**: Keep the no-atomic-clear wording from lines 383–384. Replace “neither action releases the group” with: cancelling the in-progress run, and rejecting standalone Cleanup’s only job, ends that run and promotes the pending run; rejecting one Test job does not free the slot while siblings such as `checkbuild` or `checkdocs` still run. Say the handoff is not an approval prompt, and that any later environment prompt is not a barrier to rely on. Require unwanted pending runs to be cancelled before either action.

### Issue 3: Expected exit can still promote a queued Test onto a partial Cleanup
- **Severity**: Medium
- **File**: DEVELOPER.md
- **Line(s)**: 386–392, 406–414
- **Description**: Unwanted pending work is cancelled “before slow workspace assessment”, and the recheck is attached to rejecting or cancelling. Natural completion, the 30-minute Cleanup timeout (`.github/workflows/cleanup-script.yml` line 29), and runner loss are named as promoters, but they are not given the same cancel-pending action when the operator intends the run to finish. Step 4 says not to start Test after a failed intended Cleanup, and it monitors arrivals only for an operator-initiated Cleanup dispatch. A one-time idle observation then authorizes the Test re-run. Retained line 366 already says owner coordination does not replace the idle check, and facts line 9 plus the same 02:00 UTC cron mean a new run can arrive or replace the single pending slot between observation and start.
- **Risk**: Intended Cleanup times out or finishes with partial hard-deletes while a Test or second Cleanup is pending. That run is promoted before the post-exit assessment. “Do not start Test” is too late. A later idle glance is also not a hold: a scheduled or dispatched Cleanup can replace the pending Test re-run or occupy the group before the re-run is actually in progress.
- **Suggested Fix**: Require the queue to be drained before any expected exit (completion, approaching timeout, or runner loss), and keep watching it until that run has exited and the post-exit inventory is done. Authorize a Test re-run only when idle is observed again immediately before the re-run, and keep monitoring until that re-run is the in-progress run in a window with no Cleanup. State that this manual sequence is not an atomic lock.

### Issue 4: “Any Cleanup” mixes best-effort pre-clean with standalone hard-delete
- **Severity**: Low
- **File**: DEVELOPER.md
- **Line(s)**: 378–380, 395–399, 406–412
- **Description**: Step 4 applies the post-exit stop-Test rule to “any” Cleanup run. The Test workflow’s `cleanup-workspace` job is the best-effort pre-clean: `continue-on-error: true` and `-Force -RecoverableItemsOnly` (`.github/workflows/test.yml` lines 60–64 and 83; facts JSON line 50). Standalone Cleanup is bare `-Force` (`.github/workflows/cleanup-script.yml` line 45; facts JSON line 52). Step 2’s “known `-Force` invocation” matches both command lines. The retained paragraph says pre-clean is best-effort and tests continue; step 4 says if Cleanup did not succeed, do not start Test.
- **Risk**: In an incident the operator either halts a Test matrix the workflow is built to continue after pre-clean failure, or treats a standalone hard-delete like best-effort pre-clean and continues after a partial `-Force` run.
- **Suggested Fix**: Name the workflow. Apply the partial-state stop, no automatic retry, and no Test start only to standalone `cleanup-script.yml`. Say Test `cleanup-workspace` remains best-effort recoverable-item cleanup, while its residue still needs an owner inventory before the next pre-clean is trusted. Require the executing command to be distinguished as `-Force` versus `-Force -RecoverableItemsOnly`, not merely “contains `-Force`”.

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
## Publication Note

Only the machine-local artifact path was normalized;original response,findings,verdict and finalized input/response receipt remain unchanged outside the source tree.
