---
name: pre-pr-review
description: >
  Self-review your terraform-provider-fabric PR before requesting human / Copilot
  review. Runs three independent parallel code reviews across different models,
  cross-verifies findings against the Fabric REST API docs and fabric-sdk-go, removes
  false positives, and presents consolidated findings for you to fix BEFORE you push.
  USE FOR: catching the issues the Copilot reviewer and maintainers usually flag on
  provider PRs (changelog format, example/doc schema-validity, diagnostics-before-state,
  fake handler bugs, preview-mode parity). Triggers: "review my PR", "self-review",
  "pre-review", "deep review my PR", "review before submitting", "what will reviewers flag".
---

# Pre-PR Review — Self-Review Your Provider PR Before Submitting

Run a rigorous, cross-verified review on YOUR OWN PR before a human or the Copilot
reviewer does. Scoped to the recurring findings that drive reviewer flags in
`microsoft/terraform-provider-fabric`.

## When to Use

- After `task lint` and the relevant `task testunit` pass locally
- Before pushing the commit you want a maintainer to review
- After addressing a review round, before re-pushing (catches issues introduced by the fix)

## When NOT to Use

- Before local lint/tests pass — you'll waste a pass on issues `golangci-lint` already catches
- For pure dependency-bump PRs (`build(deps)` / `ci(deps)`) — little to review
- For draft PRs still mid-edit — review the complete diff or you'll get noise

## Requirements

Designed for **Copilot CLI** (or a compatible multi-agent assistant such as Claude Code)
that exposes the `task` tool for spawning sub-agents in parallel. If your environment
only supports a single agent, run the workflow sequentially — replace each "launch in
parallel" step with three serial passes, swapping the model between runs.

## Why Cross-Verification

Single-pass reviews over-report on Fabric API / SDK specifics (fabricated endpoints,
wrong enum casing, wrong SDK method names). Cross-verification has each independent pass
validate the others' findings against the Fabric REST API docs and the `fabric-sdk-go`
source. Findings confirmed by 2+ reviewers with a citation are near-certain; lone
unverifiable claims are dropped.

---

## Workflow Overview

```
Phase 1: Fetch PR diff, launch 3 independent parallel reviews
    |
Phase 2: Each reviewer cross-verifies the other two's findings against docs/SDK
    |
Phase 3: Consolidate — remove false positives, revise partials, deduplicate
    |
Phase 4: Present findings to you for fix-up
    |
Phase 5 (optional): If your PR is open, post findings as PR comments AFTER your approval
```

Default behavior surfaces findings to you so you can fix them on your branch before
pushing. Posting to the PR is opt-in.

---

## Phase 1 — Independent Parallel Review

### Step 1.1: Fetch the PR diff

If you have a PR number:
```bash
gh pr view <PR_NUMBER> --json number,title,body,headRefOid,files
gh api repos/microsoft/terraform-provider-fabric/pulls/<PR_NUMBER> \
  -H "Accept: application/vnd.github.v3.diff" > pr-diff.txt
```

If you have a local branch but no PR yet:
```bash
git diff origin/main...HEAD > pr-diff.txt
```

### Step 1.2: Launch 3 reviews in parallel

Use the `task` tool with `agent_type: "code-review"` and `mode: "background"`, three
times, with three different models. Give each reviewer the full diff + the checklist
below. If your CLI does not expose `task`, run the three reviews sequentially.

**Recommended model combination:**

| Reviewer | Model | Strength |
|---|---|---|
| 1 | `claude-opus-4.8` | Deep reasoning, Go correctness |
| 2 | `gpt-5.5` | Broad knowledge, doc-grounded evidence |
| 3 | `gemini-3.1-pro-preview` | Fast, strong at structural / pattern issues |

If a model is unavailable, fall back to the closest peer (e.g. `claude-opus-4.7-high`,
`gpt-5.4`).

### Step 1.3: Reviewer prompt (use for all 3)

```
You are reviewing a pull request for the microsoft/terraform-provider-fabric repo (a Go
Terraform provider built on the HashiCorp Terraform Plugin Framework and fabric-sdk-go).
Be precise. Cite file:line. Never flag style or formatting handled by golangci-lint.

MUST-CHECK (highest priority):
1. Correctness — SDK method names, enum constants, and REST endpoints actually exist in
   fabric-sdk-go / the Fabric REST API. No fabricated methods or flags.
2. Diagnostics & state — operation diagnostics are appended and checked BEFORE
   resp.State.Set(...); a failed sub-operation (e.g. SyncTags) must not persist state.
   Delete calls resp.State.RemoveResource on success. set()/setter diagnostics are not
   ignored.
3. Nil / bounds safety — SDK pointers (e.g. Etag) are nil-checked before deref; slices
   (Value[0]) are length-checked. Applies to fakes too.
4. Fake handlers (internal/testhelp/fakes/**) — new(...) is never called with a VALUE or
   enum constant (does not compile); use azto.Ptr / to.Ptr. Fake responses populate every
   field the model's set() dereferences. *.json(.tmpl) fixtures render to valid JSON.
5. Schema/model type match — schema customtypes (e.g. customtypes.UUIDType{}) match the
   model field type (customtypes.UUID, not types.String).
6. Validators — write-only *_wo attributes pair with AlsoRequires to their *_wo_version;
   definition path-key validators use pattern matching when the format has wildcard parts
   (e.g. EntityTypes/*); definition SizeAtMost is not smaller than the format's part count.
7. Preview parity — if the resource Configure calls fabricitem.IsPreviewMode, the matching
   data source does too; no hard-coded false for a preview/GA flag that carries IsPreview.
8. Security — no hardcoded secrets/tokens/connection strings in code or fixtures.

SHOULD-CHECK:
9. Changelog — .changes/unreleased/*.yaml: custom.Issue is present and a valid
   issue/PR number (changie generates it quoted, which is correct — do NOT flag
   the quoting); the body names the ACTUAL resource/data source added.
10. Examples & docs — examples/**/*.tf and docs/**/*.md use values that pass schema
    validation: tokens_delimiter in {{}} / <<>> / @{}@ / ____; processing_mode is exact
    case (GoTemplate / Parameters / None); format is set when definition/output_definition
    is set and uses a valid value; referenced fixture files exist; ${local.path} has a
    matching locals block; example folder / import.sh name matches the actual resource.
11. Descriptions — MarkdownDescription (never Description); no copy/paste wrong subject
    (e.g. inbound text on an outbound attribute); attribute names in prose are snake_case;
    learn.microsoft.com URLs have no en-us (or any) locale; enum lists have a leading space.
12. Tests — new behavior (new definition format, job type, attribute) is asserted, not just
    generic attributes; black-box <pkg>_test package; resource.ParallelTest is NOT used when
    the test mutates the shared fakes.FakeServer; test names follow TestUnit_/TestAcc_<Type>_*.
13. Registration — a new resource/data source is registered in internal/provider/provider.go
    and ships an examples/ entry with regenerated docs.

NICE-TO-HAVE (nit, non-blocking):
14. Error messages name the right operation/entity (no "delete" text in Create; no
    "Workspace" in a Tag data source).
15. No dead/unused code or commented-out validator invocations (golangci-lint fails these).

For each finding, return:
- Severity: MUST-CHECK / SHOULD-CHECK / NICE-TO-HAVE
- Dimension (number from the list above)
- File and line
- What is wrong, and what the fix should be
- Evidence (Fabric REST API or fabric-sdk-go URL, or your reasoning if unverifiable)

PR diff:
<paste pr-diff.txt content here>
```

Wait for all 3 background agents to complete. In **Copilot CLI** the runtime surfaces a
system notification automatically; retrieve results with `read_agent`. Collect findings.

---

## Phase 2 — Cross-Verify

Re-launch the same 3 agents (background mode) with this prompt, giving each reviewer the
other two's findings:

```
You previously reviewed this PR. Now verify the OTHER TWO reviewers' findings against
official sources. For EACH finding, decide:

- CONFIRMED: cite the doc / source URL that supports it
- REFUTED: cite the doc / source URL that disproves it
- UNVERIFIABLE: cannot confirm or deny from available sources (explain why)

Key sources:
- Fabric REST API: https://learn.microsoft.com/rest/api/fabric/
- Fabric docs: https://learn.microsoft.com/fabric/
- fabric-sdk-go: https://github.com/microsoft/fabric-sdk-go
- Terraform Plugin Framework: https://developer.hashicorp.com/terraform/plugin/framework
- Repo conventions: .github/instructions/*.instructions.md, .changie.yaml, Taskfile.yml

Also: do a fresh correctness scan — flag any new issue you missed in Round 1.

Other reviewers' findings to verify:
<paste R1 findings from the other two>

PR diff:
<paste pr-diff.txt content>
```

---

## Phase 3 — Consolidate

Build the consensus matrix:

| # | Finding | R1 | R2 | R3 | Verdict |
|---|---|---|---|---|---|
| F1 | SDK method wrong | CONFIRMED | CONFIRMED | CONFIRMED | **CONFIRMED (3/3)** |
| F2 | Default value | CONFIRMED | REFUTED | -- | **DISPUTED** |
| F3 | Enum casing | -- | REFUTED | REFUTED | **LIKELY FALSE** |

("--" means this reviewer originally raised the finding, so they did not verify it.)

| Consensus | Action |
|---|---|
| 2/2 or 3/3 CONFIRMED | Include in final findings |
| Mixed, 0 REFUTED | Include |
| Any REFUTED with documented evidence | Investigate; remove if evidence is strong |
| Majority REFUTED | Remove — false positive |

Then:
1. **Remove false positives** — anything REFUTED by 2+ with evidence, or by 1 with a definitive citation.
2. **Revise partials** — if the core insight is valid but the specific claim is wrong, update to verified facts only.
3. **Deduplicate** — merge same-issue findings; keep the most precise description, combine evidence, use highest severity.
4. **Add new findings** — include issues newly surfaced in Round 2 by 2+ reviewers with evidence.

---

## Phase 4 — Present to You

Output consolidated findings as a numbered list:

```markdown
## Self-Review Findings — PR #<NUMBER>

| # | Tier | Dimension | File:Line | Summary |
|---|---|---|---|---|
| 1 | MUST-CHECK | Diagnostics & state | internal/services/foo/resource_foo.go:120 | State written before checking SyncTags diags |
| 2 | SHOULD-CHECK | Changelog | .changes/unreleased/added-...yaml:2 | body names `fabric_foo_bind` but PR adds `fabric_foo_binding` |
| 3 | SHOULD-CHECK | Examples & docs | examples/resources/fabric_foo/resource.tf:12 | tokens_delimiter "##" fails schema validation |

### Finding 1 — MUST-CHECK: Diagnostics & state
**File:** `internal/services/foo/resource_foo.go:120`
**Issue:** `SyncTags` diagnostics are appended after `resp.State.Set(...)`, so a failed tag
sync still persists new state on a failed apply.
**Evidence:** Terraform Plugin Framework — handle diagnostics before writing state.
**Suggested fix:** Append `tagDiags` and `return` on error before calling `State.Set`.

### Finding 2 — ...
```

Then ask: **"Fix these locally first? Or post them as PR comments now?"**

Default recommendation: **fix locally first.** A clean PR avoids review-cycle ping-pong.

---

## Phase 5 — Post as PR Comments (Opt-In Only)

Only execute when you explicitly request it AND a PR exists.

```python
import json
comments = []
for finding in approved_findings:
    comments.append({
        "path": finding["file"],
        "line": finding["line"],
        "body": f"**{finding['tier']}** — {finding['dimension']}\n\n"
                f"{finding['description']}\n\n"
                f"**Evidence:** {finding['evidence']}"
    })
review = {"body": " ", "event": "COMMENT", "commit_id": "<HEAD_SHA>", "comments": comments}
with open("review.json", "w", encoding="utf-8") as f:
    json.dump(review, f)
```

```bash
gh api repos/microsoft/terraform-provider-fabric/pulls/<PR_NUMBER>/reviews \
  --method POST --input review.json
```

Then cleanup: `rm pr-diff.txt review.json`.

---

## Must

- **Verify every API/SDK claim** against the Fabric REST API docs or `fabric-sdk-go` — correctness is #1
- **Use 3 independent reviewers across different models** — single-model passes over-report
- **Cross-verify before presenting** — never surface a finding that 2+ reviewers refuted with evidence
- **Present findings to you first** — never auto-post PR comments
- **No implementation details in posted comments** — never mention models, agents, AI, or the review methodology in anything that appears on the PR
- **All posted comments are inline** — placed on exact file:line, never as review body text
- **Fetch the full diff** — don't rely on file listings alone

## Prefer

- **Unanimous (3/3) findings** — highest signal
- **MUST-CHECK first** — the highest-priority correctness issues
- **Specific source URLs** as evidence, not general "the docs say"
- **Background mode** for the parallel reviewers — 3x throughput vs sync
- **Python** for the JSON payload — avoids shell escaping issues
- **Fix locally, then push** — one clean push beats 5 review-then-fix rounds

## Avoid

- **Posting comments without your approval** — always present findings first
- **Mentioning models, agents, or AI** in anything the PR author or reviewers will see
- **Flagging style / formatting** — golangci-lint, gofumpt, and markdownlint own those
- **Using the review body for findings** — use inline comments only
- **Assuming enum casing / API paths** — different Fabric APIs differ; verify `kind`, `type`, and enum values against fabric-sdk-go or actual responses
- **Recommending raw `go test`** — this repo requires `task testunit` / `task testacc` (they set env like `FABRIC_PREVIEW=true`)
- **Running this on a PR with hundreds of files** — split by changed-files focus or you'll exceed reviewer context

---

## Examples

### Example 1: Self-review a resource-onboarding PR

**User prompt:** "Deep review my PR #1014 before I request review."

**Workflow:**
1. Fetch diff for #1014.
2. Launch 3 background reviewers with the checklist prompt.
3. After all 3 complete, launch round 2 with cross-verification.
4. Consolidate: 9 raw → 5 confirmed, 3 refuted, 1 revised.
5. Present 6 findings (5 confirmed + 1 revised) to the user — e.g. a changelog
   body naming the wrong resource, `tokens_delimiter "##"` in the example, missing acceptance test.
6. User fixes locally, commits, re-pushes.

### Example 2: Pre-push self-review (no PR yet)

**User prompt:** "I'm about to push my branch. Run a deep review on `git diff origin/main...HEAD`."

**Workflow:**
1. Capture local diff.
2. Same 3-reviewer + cross-verify flow.
3. Findings surface BEFORE the PR exists — fix and push a clean first commit.

### Example 3: Handling a disputed finding

During Round 2, R1 claims an enum constant `AgentStateActive` doesn't exist. R2 and R3
refute, citing the `fabric-sdk-go` package where the constant is defined.

**Result:** Finding is removed before presenting (2/3 refuted with source evidence).
Cross-verification prevented a false positive that would have wasted a maintainer's time.

---

## Limitations

- **Model availability** — fall back to peer models while keeping the 3-reviewer minimum
- **Rate limits** — 3 concurrent agents per round (6 total across two rounds); batch if rate-limited
- **Documentation lag** — official docs may lag API reality; note when verification is inconclusive
- **Context window** — PRs > ~100KB diff may need to be split by changed files
- **Runtime correctness** — some bugs only surface by running `task testunit` / `task testacc`, not by docs review; run the tests too

---

## See Also

- `.github/instructions/code-review.instructions.md` — the Go review checklist this skill internalizes
- `.github/instructions/fakes-review.instructions.md` — fake-handler review rules
- `.github/instructions/examples-review.instructions.md` — example `.tf` review rules
- `Taskfile.yml` — `task lint`, `task testunit`, `task testacc`, `task docs`
- `.changie.yaml` — changelog entry format (`custom.Issue` is an int)
- `CONTRIBUTING.md` — contribution requirements
