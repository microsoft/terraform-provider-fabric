---
producer: "review-code"
schema_version: 1
tier: "Medium"
contract: "direct"
execution_backend: "copilot-cli"
attempt: 5
stage_attempt: 8
round: 3
timestamp: "2026-10-31T09:37:52Z"
outcome: "issues_found"
reviewer_model: "grok-4.7"
input_sha256: "1dac92489d63852a9f043a3894791e5199b14f97983b7114e81c089449bee49a"
---

## Review Summary
- **Round**: 3
- **Theme**: Edge cases & robustness
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Artifact**: reviews/task-1119-attempt-5-stage-attempt-8-review-3-grok-4.7.md
- **Issues Found**: 3
- **Verdict**: ISSUES_FOUND

## Evidence Checklist
- [x] Recovery text in `DEVELOPER.md` lines 360-416 checked against `context-only-public-docs/verified-github-workflow-facts.json` (one running run, one replaceable pending run, job-level environment approval, selected-ref dispatch).
- [x] Standalone Cleanup shape checked in `.github/workflows/cleanup-script.yml` lines 5-42: shared group, `cancel-in-progress: false`, single `development` job, 30-minute timeout, `Remove-WorkspaceRSItems.ps1 -Force`.
- [x] Test workflow checked in `.github/workflows/test.yml`: same group at lines 20-24; pre-clean `-Force -RecoverableItemsOnly` at line 83 with `continue-on-error`; matrix timeout 100; non-environment jobs (`changes`, `checkbuild`, `checkdocs`, `test-summary`) still hold the run.
- [x] `-Force` skips confirmation and hard-deletes with no undo in `tools/scripts/Remove-WorkspaceRSItems.ps1` lines 17 and 498-545; preserved names/IDs are filters, not a rollback.
- [x] Unchanged qualifications remain at `DEVELOPER.md` lines 360-376 and 411-416: merge-ref caveat, workflow approval vs job environment approval, idle definition, 30/100-minute bounds, shared `FolderRS1`/`FolderRS2`. No tests were executed.

## Issues

### Issue 1: Approval can start -Force without the selected-ref review
- **Severity**: Medium
- **File**: DEVELOPER.md
- **Line(s)**: 392-403
- **Description**: Step 2 says to approve intended work after assessment, and that standalone Cleanup then runs `-Force`. The checklist that actually bounds that work (target/preservation, `development`, `-Force`, same group, `cancel-in-progress: false`, changed bytes not approved) is only in step 3, and only "Before `workflow_dispatch`". A scheduled or already dispatched approval-waiting Cleanup run is not a new dispatch. Its executing SHA is fixed at queue time and can differ from a branch tip reviewed earlier. Following the numbered steps can approve that deployment before those files are reviewed.
- **Risk**: Approval is what starts hard deletion. An owner can release `-Force` against the shared workspace using an unreviewed or stale script/workflow, including one that no longer preserves fixtures or no longer shares the lock.
- **Suggested Fix**: Make the step 3 checklist a precondition of approving any Cleanup deployment, using that run's head SHA, not only of a new dispatch. If the SHA's workflow or script was not reviewed, cancel or reject; do not approve.

### Issue 2: Pending work can start during assessment or after rejection
- **Severity**: Medium
- **File**: DEVELOPER.md
- **Line(s)**: 383-397
- **Description**: The intro correctly says cancellation or rejection is not an atomic group clear and that a pending run can start as soon as the active run ends. Step 1 still starts with workspace assessment, and only then says to cancel pending runs "before ending" the active run. A run can also end on its own (completion, 30- or 100-minute job timeout, runner loss) while that assessment is in progress. Separately, step 2 says rejecting a deployment does not stop the whole workflow. That is true for multi-job Test runs, but `.github/workflows/cleanup-script.yml` has only `cleanup-workspace`. Rejecting that deployment fails the only job, the run can exit, and the single pending run can be promoted without the queue-first cancel that step 1 attaches only to cancelling the Cleanup run.
- **Risk**: Unwanted Cleanup becomes the in-progress run and can run `-Force`, or it occupies the slot, while the operator is still assessing or believes rejection is holding the group. Deletions are not undone.
- **Suggested Fix**: Say that any exit of the active run, including natural completion, timeout, runner loss, cancellation, or rejection of Cleanup's only job, can promote the pending run immediately and is not an approval prompt. Cancel unwanted pending runs before slow workspace assessment and again immediately before reject, cancel, or expected exit; recheck until that run has exited. Do not call this sequence atomic.

### Issue 3: Idle after Cleanup is treated as success
- **Severity**: Low
- **File**: DEVELOPER.md
- **Line(s)**: 404-408
- **Description**: Step 4 says an owner-approved Cleanup should "finish", then the group should be confirmed idle, then required Test checks re-run in a separate window. Cleanup's job timeout is 30 minutes (`.github/workflows/cleanup-script.yml` line 25). Timeout, step failure, or runner loss also makes the run exit and the group look idle, after a partial hard-delete. Step 1 tells the operator to assess resource changes only on the cancel path. The success path does not require checking the run conclusion or re-assessing the workspace.
- **Risk**: Tests are re-run against a partially deleted `WorkspaceRS`. The in-workflow pre-clean only removes recoverable items, so active leftovers from a failed `-Force` pass are not reclaimed.
- **Suggested Fix**: After the Cleanup run exits, require a successful conclusion and a fresh workspace assessment before any Test re-run. Treat timeout, failure, and runner loss like cancellation: partial deletions stand.

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
