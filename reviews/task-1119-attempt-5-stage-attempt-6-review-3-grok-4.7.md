---
producer: "review-code"
schema_version: 1
tier: "Medium"
contract: "direct"
execution_backend: "copilot-cli"
attempt: 5
stage_attempt: 6
round: 3
timestamp: "2026-10-31T08:50:09Z"
outcome: "issues_found"
reviewer_model: "grok-4.7"
input_sha256: "efec2e9f92c08c50c3cd0a5b5ddaf0704405d765dbf2fd33d3f861c5e85525eb"
---

## Review Summary
- **Round**: 3
- **Theme**: Edge cases & robustness
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Artifact**: reviews/task-1119-attempt-5-stage-attempt-6-review-3-grok-4.7.md
- **Issues Found**: 3
- **Verdict**: ISSUES_FOUND

## Evidence Checklist
- [x] Compared the delta with `DEVELOPER.md` lines 348-390 against `.github/workflows/test.yml` and `.github/workflows/cleanup-script.yml` (shared group `fabric-testacc-environment`, `cancel-in-progress: false`, `0 2 * * *`, `environment: development`, timeouts 30 and 100).
- [x] Checked `context-only-public-docs/verified-github-workflow-facts.json`: one running run, one replaceable pending run, incoming `cancel-in-progress: true` cancels a running run, job environment approval is not concurrency approval.
- [x] Confirmed standalone Cleanup runs `tools/scripts/Remove-WorkspaceRSItems.ps1 -Force` (`.github/workflows/cleanup-script.yml` line 45); in-workflow pre-clean adds `-RecoverableItemsOnly` (`.github/workflows/test.yml` line 83). Script header hard-deletes items and non-preserved folders.
- [x] Confirmed `check-test-matrix` is the documented PR status check and fails when `needs.test.result` is `failure` or `cancelled` (`.github/workflows/test.yml` lines 596-607). Test environment jobs are staged: pre-clean, then `test-auth-spn` and `test`.
- [x] Confirmed unchanged bounds and advisories remain at `DEVELOPER.md` lines 388-394, and the merge-ref / workflow-vs-job approval caveat remains at lines 368-375. Supplied diff is `DEVELOPER.md` only.
- [ ] Markdown CLI, actionlint, and live Actions were not re-run; this pass is read-only source inspection.

## Issues

### Issue 1: Approval is listed as a neutral way to clear an environment wait
- **Severity**: Medium
- **File**: DEVELOPER.md
- **Line(s)**: 381-385
- **Description**: The new recovery text treats approving the job, rejecting it, or cancelling the run as equivalent ways to resolve an environment-approval wait. They are not. While the job is still waiting, approval starts it. Standalone Cleanup then runs `Remove-WorkspaceRSItems.ps1 -Force`, which hard-deletes active WorkspaceRS items and non-preserved folders. Reject or cancel does not start that script. On the Test workflow, approving `cleanup-workspace` also does not finish the wait: `test-auth-spn` and the test matrix each use `environment: development` and can wait again after pre-clean. Other jobs (`checkbuild`, `checkdocs`, `test-summary`, `check-test-matrix`) have no environment and can keep the workflow-level slot until the whole run exits.
- **Risk**: An owner trying only to free the concurrency slot approves the waiting Cleanup job and hard-deletes shared fixtures before assessment is applied. Or they approve one Test job, believe the wait is over, and dispatch Cleanup while later matrix jobs still hold the group. With `cancel-in-progress: false`, that dispatch does not stop the in-progress run and can cancel a different pending required check.
- **Suggested Fix**: Keep all three actions, but state the side effects. Approve only if that waiting job is the work intended after assessment; say standalone Cleanup then runs `-Force`. To free the slot without starting the job, reject the deployment or cancel the workflow run. Say one approval does not clear later Test environment jobs, and that the group stays held until the whole run exits. Leave the existing statements that job timeout and Cleanup redispatch do not release the wait.

### Issue 2: Required-check re-run is only tied to dispatch supersession, and "idle" is undefined
- **Severity**: Medium
- **File**: DEVELOPER.md
- **Line(s)**: 364-366, 385-387
- **Description**: The new procedure tells the owner to cancel the run or reject the waiting job, but the re-run rule applies only "if a dispatch supersedes a required check." The run holding the slot is often the Test workflow. Cancelling that run, or failing its test jobs, does not satisfy the required check: `check-test-matrix` exits 1 when the test result is `failure` or `cancelled`, and a cancelled workflow run is not success. "After the group is idle" is also never defined. The only precise gate, "no in-progress or pending run," is attached to Cleanup redispatch, not to the re-run. A run waiting on environment approval holds the running slot even when no test job is executing.
- **Risk**: Slot release leaves the PR required check cancelled or failed, and the runbook never says to re-run it. If the owner treats "idle" as "no test jobs executing" and re-runs while the approval wait or Cleanup still occupies the group, the re-run sits in the single pending slot and the next newer run, including the 02:00 Test or Cleanup schedule, cancels it again.
- **Suggested Fix**: Define idle once as no in-progress group run and no pending group run, and state that an environment-approval wait is in-progress. After that state, re-run the required Test check if it was superseded or if the owner cancelled or rejected that Test run to free the slot. Do not re-run in the same idle window as a Cleanup dispatch; finish one run, confirm the group is idle again, then start the other.

### Issue 3: New redispatch rule sits outside the executed-file caveat
- **Severity**: Low
- **File**: DEVELOPER.md
- **Line(s)**: 364-366, 370, 384-386
- **Description**: Line 370 limits "the guarantees above" when an incoming workflow renames or drops the group or sets `cancel-in-progress: true`. The new redispatch rule is below that sentence. It says a retained false setting will not cancel the in-progress run, but it does not say the setting comes from the ref selected for `workflow_dispatch`. Lines 364-366 also say to coordinate with the environment owner rather than stating, in that warning, that the dispatch cancels the pending Test run and must not be sent while any group run is in-progress or pending.
- **Risk**: Dispatching a PR ref whose executed `cleanup-script.yml` sets `cancel-in-progress: true` cancels the in-progress run, including an approval wait or running tests, which can skip destroy. A ref that drops the group does not take the lock and can run `-Force` beside tests. Owner coordination alone does not stop either outcome.
- **Suggested Fix**: At the redispatch step, require the executed Cleanup ref to still use group `fabric-testacc-environment` and `cancel-in-progress: false`. State that a true setting on that incoming run can cancel the in-progress run, and that dropping the group runs outside the lock. In the line 364 warning, say the replaced Test run is cancelled and that the no-in-progress-or-pending gate is mandatory, not replaced by owner coordination.

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
