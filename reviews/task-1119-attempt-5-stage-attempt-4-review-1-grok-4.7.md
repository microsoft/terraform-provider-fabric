---
producer: "review-code"
schema_version: 1
tier: "Medium"
contract: "direct"
execution_backend: "copilot-cli"
attempt: 5
stage_attempt: 4
round: 1
timestamp: "2026-10-31T07:42:20Z"
outcome: "issues_found"
reviewer_model: "grok-4.7"
input_sha256: "0ef5c80b66be43214a65eb40d5e492da3a5862593bd9bfb1d8940429feca69f6"
---

## Review Summary

- **Round**: 1
- **Theme**: Broad sweep
- **Mode**: resolved-plan
- **Model**: grok-4.7
- **Artifact**: reviews/task-1119-attempt-5-stage-attempt-4-review-1-grok-4.7.md
- **Issues Found**: 4
- **Verdict**: ISSUES_FOUND

## Evidence Checklist

- [x] Inspected `.github/workflows/test.yml` workflow-level group `fabric-testacc-environment`, `cancel-in-progress: false`, src filter including both workflow files, development job `if`/`needs`, unchanged matrix, and recoverable-only pre-clean (`timeout-minutes` absent on `changes`, `test-auth-spn`, `checkdocs`).
- [x] Inspected `.github/workflows/cleanup-script.yml` for the same workflow-level group, no job-level `concurrency`/`if`/`needs`, and full cleanup via `Remove-WorkspaceRSItems.ps1 -Force`.
- [x] Inspected `internal/testhelp/folder_fixtures_test.go` against `context-only-public-dependencies/yaml.v3@v3.0.1/decode.go` (`scalar` string coercion, `Node.Decode`, Kind 0 null). Guard checks wiring, not Picomatch.
- [x] Inspected ML CRUD steps, `testhelp.FolderResource`, `fabricitem` move path (`folder_id` is not `RequiresReplace`), and folder `IsPreview: true` with Taskfile `FABRIC_PREVIEW`.
- [x] Compared `DEVELOPER.md` acceptance-test text with requirements 1–5. No other workflow uses `fabric-testacc-environment`.
- [ ] Did not execute tests, actionlint, or a hosted concurrency run. Supplied validation was not re-run.

## Issues

### Issue 1: Shared group can still be cancelled by an incoming run

- **Severity**: Medium
- **File**: DEVELOPER.md
- **Line(s)**: 360-361
- **Description**: Published Test and Cleanup files set `cancel-in-progress: false` (`.github/workflows/test.yml` 23-27, `.github/workflows/cleanup-script.yml` 10-14). GitHub applies cancellation from the incoming run in a repository-scoped group. A pull request whose workflow sets this group with `cancel-in-progress: true`, or any other workflow that does, cancels the active run. The previous per-ref group could not cancel Cleanup or another ref. Docs say a newer commit does not cancel an active run, which is true only for incoming runs that also set false.
- **Risk**: Cancelling Cleanup mid hard-delete can leave WorkspaceRS partially deleted. Cancelling Test interrupts acceptance tests and drops the lock without a completed teardown. Fork or same-repo PRs execute their workflow file on `pull_request` before environment approval.
- **Suggested Fix**: Document that false only protects against other false runs, and that a PR changing this group or adding another workflow on it can cancel an active Test or Cleanup run. Optionally scan every `.github/workflows/*.yml` and fail if this group is paired with `cancel-in-progress: true`. Do not drop the shared group or fork-PR participation.

### Issue 2: Isolation guard treats a missing cancel-in-progress key as safe

- **Severity**: Low
- **File**: internal/testhelp/folder_fixtures_test.go
- **Line(s)**: 25-28, 74-76
- **Description**: `CancelInProgress` is a `bool`. yaml.v3 leaves a missing key false, and GitHub also defaults to false, so deleting the explicit key still passes the guard. The persistent regression is structural wiring, not default equivalence.
- **Risk**: A later edit can remove the explicit setting without a failing guard, so the required wiring is no longer locked.
- **Suggested Fix**: Decode workflow `concurrency` as `yaml.Node` and require a boolean `cancel-in-progress` key whose value is false on both files.

### Issue 3: Development-job lock guard is a fixed allow-list

- **Severity**: Low
- **File**: internal/testhelp/folder_fixtures_test.go
- **Line(s)**: 109-137
- **Description**: Only `cleanup-workspace`, `test-auth-spn`, and `test` are checked for a job-level concurrency override and the cleanup dependency. The job struct does not read `environment`. Commented `test-auth-msi` (`.github/workflows/test.yml` 151-158) targets `development`, needs only `changes`, and would not be seen if uncommented. A new development job with its own `concurrency` block would also pass.
- **Risk**: That job can run beside recoverable cleanup, or a per-job lock can drop matrix legs, without failing the regression test.
- **Suggested Fix**: Fail if any job with `environment.name: development` has job-level concurrency, a non-`changes.src` `if`, or (except `cleanup-workspace`) is missing `needs: cleanup-workspace`.

### Issue 4: Deferred timeout and workspace-lifecycle advisories are not stated

- **Severity**: Low
- **File**: DEVELOPER.md
- **Line(s)**: 353-370
- **Description**: The new text covers lock lifetime, approval waits, fork participation, uncovered local/manual paths, and pending Cleanup replacement. It does not say matrix legs still run in parallel in one WorkspaceRS, other tests still share FolderRS1/FolderRS2, or that recoverable-only pre-clean (`.github/workflows/test.yml` 81-83) does not delete active leaked test folders. Several jobs have no `timeout-minutes`, so a hung job can hold the lock until the platform default or a manual cancel. These remain scope-deferred, not resolved.
- **Risk**: Operators can treat the lock as serializing all fixture use and as self-expiring, and can assume the next Test run cleans leaked ML folders.
- **Suggested Fix**: Add those residual advisories to DEVELOPER.md only. Do not add mandatory timeouts or per-run workspace identity.

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

The machine-local Artifact summary path was normalized and generated review Markdown formatting was applied. Original response and finalized receipt remain immutable outside this source tree;findings,verdict and input provenance are unchanged.

## Driver Decision and Aggregate Outcome

- **Plan/source:** Human-selected Medium, all six concurrent `grok-4.7` invocations at `xhigh`, public Task1119, attempt5, stage-attempt4. Current primary runtime identity is not independently exposed; exact reviewer selection is explicit. The complete candidate retains moderate YAML-wiring and workflow/fixture interaction complexity.
- **Unchanged source:** Four-file diff `082e8e972f14d3143900099fcca91a563b68606b31cecc154a9f4ef51e6628d9`; frozen complete review diff `170e459f45e87bdfb17216b354a98b8ff7e187e078b9dac8d3bb97b4dd239896`. Published head remains `6a27d3e1ad0c5b2a3aa73175917d60042e96366e`.
- **Corrected context:** Declared public yaml.v3 v3.0.1 and the needed testify v1.12.1 require/assert implementation sources and licenses were supplied only under each checkout's context-only directory. Existing Go module verification, exact go.sum h1 values, cache ZIP hashes, independent ZIP h1 and go.mod h1 checks passed. Copied files have per-file hashes and preserved LICENSE/NOTICE metadata. No dependency/version/install/global-path change occurred, and these copies are not product-diff or commit content.
- **Actual execution:** All six skill-owned Execute invocations finalized with runtime `grok-4.7`. Round4 CLEAN; rounds1/2/3/5/6 ISSUES_FOUND, with 4/4/5/3/3 findings respectively. Aggregate: **ISSUES_FOUND**, 19 original findings retained; not clean or merge clearance.
- **Independence/read audit:** All six physical source trees were separate and feedback-free. Actual tool reads stayed inside the allowed snapshots, including the supplied public dependency paths. No outside read, feedback exposure, identity substitution or missing-context failure occurred. The earlier six stopped receipts remain failures with zero signoffs; none was reused.
- **Finding disposition:** Incoming-workflow concurrency/documentation qualification remains an open owner review question, including the original High finding. Future structural guard strictness/coverage, fixture-helper/refactor and diagnostic suggestions remain visible; no current product-source failure was established merely by a proposed future workflow mutation. Optional job bounds and workspace/lifecycle enhancements remain scope-deferred, not newly attributed to a fresh human decision.
- **Validation:** Prior executed structural red/green, configured Go lint, YAML, actionlint, Markdown and whitespace evidence applies to these unchanged source bytes. No code or product test was rerun for this context-only correction. Actual Picomatch evidence remains distinct from the persistent structural Go/YAML test.
- **Publication:** No further source fix, commit, push or PR write occurred. Coordinator/human must reconcile the retained findings and inspect the exact final diff/evidence before publication. Context-only module sources must never be staged with product changes.

## Final Driver Reconciliation

This later driver decision supersedes the open disposition above, not the original reviewer response, findings, verdict or finalized input/response provenance.

- **Effective plan:** Human-selected Medium, six independent `grok-4.7` lenses at `xhigh` through the skill-owned Execute harness. The original High complexity rationale and unavailable-model blocker remain preserved. Primary runtime identity remains independently unverified; the explicit human reviewer override supplies the selection.
- **Final source:** Follow-up four-file diff SHA256 `97d98ae354f5a34468c1fabb4f1747e1223a1b6e7e64dc03cd6156d3348531a4`; complete frozen review diff SHA256 `101f6008c56a6baf905a3b6da17aa8b9389a03145d04c5da7504ede395ad55b4`. The published head remains `6a27d3e1ad0c5b2a3aa73175917d60042e96366e`.
- **Issue 1 and matching High concern:** Resolved by qualifying assurances to executed workflows retaining the same group and `cancel-in-progress: false`, distinguishing workflow approval from environment-job approval, and documenting incoming workflow changes and PR merge-ref execution. Official GitHub and pinned source evidence is retained. Qualification lenses1/6 at stage-attempt5 finalized CLEAN with actual `grok-4.7/xhigh` provenance.
- **Issue 4 and concrete recovery notes:** Documentation now covers whole-workflow slots, approval waits, pending replacement, promotion, owner-coordinated redispatch, unchanged bounds, remaining shared fixtures, pre-clean versus Terraform destroy versus standalone cleanup, exact executing revision, and independent cleanup-state evidence. Eighteen subsequent concrete edge-documentation findings were corrected; final edge lens3 at stage-attempt11 finalized CLEAN. No scheduling, timeout, provider, fixture or cleanup-script behavior was changed.
- **Issues 2/3 and remaining advisory notes:** Nonblocking structural-strictness/future-job suggestions remain visible. Current checked jobs are complete; missing cancellation-key decoding matches the platform's false default. Original-PR helper/refactor, assertion/naming diagnostics and timeout/workspace-lifecycle suggestions retain source/base evidence and scope dispositions. No speculative generic helper was added, and no suggestion is labeled human-accepted.
- **Preserved results:** All19 original and18 later edge notes remain retained, including dirty verdicts and failed/stopped attempts. Unchanged executable/Go scope retains stage-attempt4 lenses2/4/5 with their original outcomes and driver dispositions. The result is reconciled in-scope source, not a claim that every historical raw verdict or every lens on the latest doc bytes is CLEAN.
- **Validation:** Structural original-source RED and corrected-source GREEN, final focused Task guard1/1, configured Go lint on unchanged guard bytes, YAML10, actionlint, actual Markdown CLI and whitespace passed. Persistent coverage is Go/YAML structural wiring, not full Picomatch or hosted GitHub execution; the separate one-off matcher evidence remains distinct.
- **Evidence:** Current receipts and per-note source/base dispositions are in `..\..\..\driver-reconciled-final\final-review-disposition.json`, `driver-dispositions.json`, `source-base-evidence.json` and `approval-bundle.json`. Public dependency copies remain context-only and excluded from product/commit content.
- **Publication disposition:** Ready for the coordinator to show the exact revised diff and evidence to the human and obtain actual publication approval. No commit, push or PR write occurred. Hosted CI/configuration, maintainer acceptance, current-target validation, live acceptance and recovery remain separate gates.
