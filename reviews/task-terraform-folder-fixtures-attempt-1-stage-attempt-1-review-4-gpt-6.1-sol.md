## Review Summary

- **Round**: 4
- **Stage**: 4
- **Review Attempt**: 1
- **Stage Attempt**: 1
- **Theme**: Detailed correctness
- **Mode**: resolved-plan
- **Model**: gpt-6.1-sol
- **Reasoning**: max
- **Artifact**: `reviews\task-terraform-folder-fixtures-attempt-1-stage-attempt-1-review-4-gpt-6.1-sol.md`
- **Issues Found**: 0
- **Verdict**: CLEAN

This is this reviewer's independent stage result, not an aggregate approval. The supplied High plan retains its human-selected two-model reviewer set and maximum of two attempts per stage; the driver owns aggregation and advancement.

## Evidence Checklist

- [x] Independently inspected the five tracked changed files and all three new regression-test files against the supplied scope. No earlier review artifacts or private evidence were used.
- [x] Traced model configuration from `internal\services\mlmodel\resource_ml_model_test.go:240-291`: the acceptance entrypoint reads only the workspace fixture; `modelCRUDSteps` creates both folder configurations once, includes both in both steps, and changes the item reference from the first folder's `id` to the second folder's `id`.
- [x] Traced the corresponding experiment configuration and assertions in `internal\services\mlexperiment\resource_ml_experiment_test.go:214-266`. Both folder resources remain present while the display name, description and destination folder change in the second step.
- [x] Checked the helper return order and interpolation construction in `internal\testhelp\testacc.go:162-179` and `internal\testhelp\utils.go:190-192`. The returned configuration/FQN pairs are consumed in the correct order, and folder IDs remain Terraform references rather than copied static identifiers.
- [x] Followed the unchanged provider dataflow through `internal\pkg\fabricitem\resource_item.go:80-139,186-249,252-281` and `internal\pkg\fabricitem\models_resource_item.go:42-44,108-110,132-139`: creation receives the resolved folder ID, an actual folder-ID difference triggers movement, and the updated item is read back before state is saved.
- [x] Checked replacement semantics against `internal\pkg\fabricitem\resource_schema.go:162-216` and `internal\services\mlmodel\resource_ml_model.go:12-21` / `internal\services\mlexperiment\resource_ml_experiment.go:12-21`. The model retains its display name; experiment renaming is allowed; changing `folder_id` is not itself a replacement trigger.
- [x] Reviewed both `internal\services\mlmodel\folder_fixtures_test.go:48-111` and `internal\services\mlexperiment\folder_fixtures_test.go:48-111` line by line: two steps are required, workspace values are checked, traversal length is guarded before indexing, both folders must be retained, the second target must differ, and a second builder invocation must use different folder resource names.
- [x] Reviewed fake isolation and transport lifetime in both new service test files at lines `21-45`, together with `internal\testhelp\testunit.go:67-88,113-139` and `internal\testhelp\fakes\fake_typedhandler.go:176-231`. Each lifecycle uses a fresh item store, a separate folder store, fake credentials and the same pointer-backed serialized transport for provider clients.
- [x] Compared the effective workflow-level lock in `.github\workflows\test.yml:23-27` with `.github\workflows\cleanup-script.yml:12-14`, and checked `internal\testhelp\folder_fixtures_test.go:16-53`. The constant groups match, active runs are not cancelled, and the test matrix has no per-job lock.
- [x] Checked the documentation and dependency claims against source: `DEVELOPER.md:318-322,353-370`, `internal\testhelp\wellknown.go:40-61`, `.github\workflows\test.yml:29-32`, and `go.mod:5-30`. The environment/file precedence explanation matches the populated-variable case, operational lock limitations are explicit, and the newly direct dependencies match the new test imports without version changes.
- [x] Independently executed all five new targeted regression tests through `task testunit` with `TF_ACC=0`; all passed. Also independently ran workflow lint, dependency consistency, formatting-diff and tracked-diff whitespace checks described below.
- [ ] Live acceptance, service-side deletion verification, GitHub scheduling/cancellation execution, race instrumentation and the full supported Terraform/OpenTofu matrix were not performed.

## Detailed Correctness Assessment

1. **Folder ownership and step lifetime are consistent.** Each builder obtains two independent helper results before constructing its steps. Reusing those exact configuration strings prevents the move step from introducing replacement folder names or dropping the previous folder. Both configurations receive the requested workspace. Each item targets an owned folder through a reference, and each state check compares against that folder resource's actual ID. No static folder-fixture lookup remains in either changed acceptance entrypoint.

2. **The intended create/read/update/move/read path remains intact.** The model's unchanged display name avoids its existing rename-replacement rule; its description and folder change exercise update and movement. The experiment retains its permitted rename plus description update and folder movement. The unchanged provider sends the resolved folder ID on create, compares planned/current folder IDs on update, performs the move and reads state back. Terraform retains ownership of both folders for ordinary destroy; the final item-to-destination reference supplies the dependency needed to destroy the item before its occupied destination folder.

3. **The new regression checks do not contain a verified indexing, state-propagation or conversion defect.** The integer-range loops are supported by the module's Go version. Each inner iteration allocates a fresh folder map; assigning it to `previous` does not subsequently mutate an earlier step's map. The combined two-folder count, prior-folder membership checks and changed-target check rule out silently dropping one owned folder. The traversal-length condition short-circuits before indexing. The remaining HCL body/traversal assertions operate on this controlled generator's output, not arbitrary external HCL; both generated configurations also passed real local Terraform execution.

4. **The fake-only workaround remains bounded.** The second fake server's folder callbacks retain their own backing store, avoiding generic-item conversion of folder entities. The transport is installed by pointer, so provider clients share its mutex rather than independent lock copies. Unlocking is deferred, and the underlying response/error is returned unchanged. These changes are test-local and do not alter product API, retry or classification behavior.

5. **The workflow change implements the stated exclusion boundary.** The shared constant is at workflow scope, so all matrix legs remain inside the active run's lock without competing for one pending job slot. `cancel-in-progress: false` protects an active run from automatic replacement. This does not promise FIFO scheduling or preservation of every pending run; the documentation explicitly describes pending replacement and operator coordination for commands outside that repository's workflow lock.

No actionable correctness defect was verified within the reviewed diff.

## Checks Actually Executed

All test execution used the repository's `task testunit` entrypoint with `TF_ACC=0`, the provided local toolchain and local caches. Module/network lookup was disabled; no live acceptance or cloud command was invoked.

| Check | Independently observed result |
| --- | --- |
| `task testunit -- FolderFixtures github.com/microsoft/terraform-provider-fabric/internal/services/mlmodel` | Passed: 2 tests, including the generated-HCL regression and fake-backed Terraform lifecycle. |
| `task testunit -- FolderFixtures github.com/microsoft/terraform-provider-fabric/internal/services/mlexperiment` | Passed: 2 tests, including the generated-HCL regression and fake-backed Terraform lifecycle. |
| `task testunit -- FolderFixtures github.com/microsoft/terraform-provider-fabric/internal/testhelp` | Passed: 1 workflow-isolation regression test. |
| `actionlint -shellcheck="" -pyflakes="" .github\workflows\test.yml .github\workflows\cleanup-script.yml` | Exit 0; embedded shell/Python lint integrations were disabled. |
| `go mod tidy -diff` | Exit 0, no diff. |
| `gofmt -d` on the two modified service test files and three new regression-test files | Exit 0, no formatting diff. |
| `git diff --check` | Exit 0 for the tracked diff. |

An initial `task testunit` invocation used named selector-variable assignments that the task's own variable definitions did not honor. It expanded to a wider unit run and was stopped after out-of-scope failures, including tests requiring unavailable local setup. It was not a completed baseline or a green full-suite run, and those failures are not attributed to this change. The corrected positional selectors above completed independently.

## Reported Checks and Limitations

- The supplied context reported earlier red regressions and complete seven-test runs for each ML package. Those historical results were not used as independent evidence; this reviewer completed only the five targeted tests listed above.
- The successful lifecycle tests use an actual local Terraform CLI with SDK fake transports. They verify executable generated configurations and the normal fake-backed lifecycle, not real Fabric authorization, eventual consistency, folder referential constraints or service-side cleanup success.
- The two ML package runs displayed 100% package statement coverage, but those packages contain thin resource constructors and call into substantial imported implementation code. That number is not full provider, branch, teardown-order or production-behavior coverage.
- The fake stores are deliberately separate and requests deliberately serialized. Success is not a proof of behavior under real service concurrency, and no race-detector run was performed.
- Workflow lint and YAML assertions do not execute GitHub's concurrency scheduler. The pending-run replacement and outside-lock coordination caveats remain operational requirements.
- No source edits, live acceptance, cloud cleanup, external writes, commits, pushes, forks or pull requests were performed. Only this review artifact was authored.
