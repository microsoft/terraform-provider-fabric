---
producer: "review-code"
schema_version: 1
tier: "Medium"
contract: "direct"
execution_backend: "copilot-cli"
attempt: 5
stage_attempt: 10
round: 3
timestamp: "2026-10-31T10:20:47Z"
outcome: "issues_found"
reviewer_model: "grok-4.7"
input_sha256: "78ab5a5ca8ae14578003f557f783099c6b5d52d6b549b6703267e54ccc8a1481"
---

## Review Summary
- **Round**: 3
- **Theme**: Edge cases & robustness
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Artifact**: reviews/task-1119-attempt-5-stage-attempt-10-review-3-grok-4.7.md
- **Issues Found**: 3
- **Verdict**: ISSUES_FOUND

## Evidence Checklist
- [x] Inspected `DEVELOPER.md` lines 388-424 against `tools/scripts/Remove-WorkspaceRSItems.ps1` fetch, summary, and exit paths (lines 74-75, 167-171, 338, 367, 430-441, 536-539, 666-691, 719-723, 769). No tests or workflows were executed.
- [x] Confirmed mode split in `.github/workflows/cleanup-script.yml` lines 26-27 and 42-45 (`environment: development`, bare `-Force`) versus `.github/workflows/test.yml` lines 64-65 and 83 (`continue-on-error: true`, `-Force -RecoverableItemsOnly`). The new text does not claim a Test pre-clean hard gate.
- [x] Compared promotion wording in `DEVELOPER.md` lines 388-390 and 405-408 with the single Cleanup job, shared `fabric-testacc-environment` / `cancel-in-progress: false` (cleanup-script.yml lines 13-14, test.yml lines 26-27), and `context-only-public-docs/verified-github-workflow-facts.json` (job-level environment approval is not concurrency approval; one running run and one replaceable pending run).
- [x] Confirmed step 4 denies an atomic post-run inventory lock, and prior SHA, selected-ref, idle, timeout, and shared-fixture sentences remain outside this delta.

## Issues

### Issue 1: Standalone success gate is unsatisfiable and mislabels swallowed fetches
- **Severity**: High
- **File**: DEVELOPER.md
- **Line(s)**: 419-424
- **Description**: The new rule correctly treats a green conclusion and "deleted successfully" as inconclusive, then requires fetch evidence, complete deletion summaries with every failed count zero, and a matching inventory. It also says a warning or a missing summary means unsuccessful Cleanup. Under `cleanup-script.yml`'s bare `-Force`, that conjunction cannot be met. Any run that deletes at least one object logs a non-fatal WARN, `Force mode enabled - proceeding with deletion...` (`Remove-WorkspaceRSItems.ps1` line 539), before it can print `=== DELETION SUMMARY ===` and `All items, recoverable items, and folders deleted successfully!` (lines 666-679). Rate-limit retries (lines 167-171) and tenant or subscription mismatch (lines 719-723) are also non-fatal WARNs on paths that can still exit 0. The zero-delete path returns at lines 430-441 with `only preserved items found` and no summary, then line 769 still logs `WorkspaceRS cleanup process completed`. `Write-Log` exits 1 only when level is ERROR and `-Stop` is true (lines 74-75), and these WARNs pass `-Stop $false`. The swallowed fetches the paragraph is meant to catch do not match "a warning" as a discriminator: folder fetch failure is WARN and continues with an empty folder list (line 338), and standalone recoverable fetch failure is ERROR with `-Stop $RecoverableItemsOnly` false (line 367), clears the list, and does not fail the process. Either can still reach failed counts of zero and "deleted successfully", or the false `only preserved items found` line. Those WARN/ERROR lines are ordinary `Write-Host` text, not GitHub `::warning::` annotations.
- **Risk**: Literal use blocks every normal standalone recovery, so operators will waive the rule. Waiving "any warning" or "missing summary" drops the only stop-list tokens that catch a green run whose folder or recoverable enumeration was swallowed. Hard-deleted or leftover objects can then be treated as a clean workspace, and Test can start on a partial state. `only preserved items found` is especially easy to mistake for proof that nothing deletable remained.
- **Suggested Fix**: Replace the blanket warning and missing-summary rule with exact log evidence. Require the Found/No-items, Found/No-folders, and Found/No-recoverable result lines, and treat any `Failed to fetch` line as unsuccessful at either WARN or ERROR even if "deleted successfully" follows. Do not treat the expected `-Force` banner, rate-limit retry WARNs, or `WorkspaceRS cleanup process completed` as success or failure by themselves. Require zero failed counts only when a deletion summary was produced. Accept the zero-delete early return only when all three fetch result lines exist, no fetch-failure line exists, and an independent inventory of active items, folders, and recoverable items matches the pre-run assessment. State that script totals and `only preserved items found` are not that inventory.

### Issue 2: Step 1 still says promotion happens without an approval prompt
- **Severity**: Medium
- **File**: DEVELOPER.md
- **Line(s)**: 388-390, 405-408
- **Description**: Step 1 says natural completion, timeout, runner loss, or rejection promotes pending work immediately "without another approval prompt." Step 3 says promotion is not approval and a later environment gate may or may not wait. Those statements conflict. `cleanup-script.yml` has one job, and that job is `environment: development` (lines 26-27); the `-Force` script runs only after the job starts (lines 42-45). Checked public facts say environment approval permits a job to proceed and is not approval of the concurrency group. Ending the running slot can promote the single pending run, but promotion does not skip that job's environment rule. If reviewers or a wait timer exist, the promoted Cleanup job waits before any step. If they do not, hard-delete starts with no prompt. "Keep observing until that run exits" also stops at the moment the slot is freed.
- **Risk**: An operator who trusts the absolute will not watch for a new rejectable `development` wait, or will treat a missing prompt as already approved or as idle. A wait timer or another reviewer can then start standalone hard-delete. Cancellation and rejection still do not make the group atomically idle.
- **Suggested Fix**: Delete "without another approval prompt." Say freeing the running slot can promote pending work immediately, that promotion is neither approval nor idle, and the promoted run's `development` jobs may wait for a new review or may start at once. Continue watching that run until the group is idle: no in-progress run, including an approval wait, and no pending run. Cancel or reject unwanted Cleanup there.

### Issue 3: Approval checklist still names bare `-Force` as the known invocation
- **Severity**: Low
- **File**: DEVELOPER.md
- **Line(s)**: 396-401
- **Description**: Step 2 says to check "the known `-Force` invocation" before approving any waiting Cleanup job, then says to distinguish standalone bare `-Force` from Test's `-Force -RecoverableItemsOnly`. Both jobs wait on `development` (`cleanup-script.yml` line 27, `test.yml` lines 65 and 83). A presence check for `-Force` matches the safe Test pre-clean and a mutated Test command that dropped `-RecoverableItemsOnly`. The distinguish sentence is not written as a fail-closed exact match, while the checklist item names the hard-delete token as the known invocation.
- **Risk**: Approving Test `cleanup-workspace` after `-RecoverableItemsOnly` is removed starts active-item hard-delete during acceptance tests. `continue-on-error: true` (`test.yml` line 64) does not roll those deletes back. The later "no new pre-clean hard gate" sentence does not prevent that approval mistake.
- **Suggested Fix**: Require an exact invocation match in the approval checklist. Standalone `cleanup-script.yml` must be `-Force` and must not pass `-RecoverableItemsOnly`. Test pre-clean must be `-Force -RecoverableItemsOnly`. Any other invocation is changed execution and is not approved.

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

Only the machine-local artifact path was normalized;original response,findings,verdict and finalized input/response receipt remain unchanged outside the source tree.
