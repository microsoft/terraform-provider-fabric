## Review Summary

- **Round**: 3
- **Stage**: 3 — Edge cases & robustness
- **Theme**: Edge cases & robustness
- **Focus**: Error handling, boundary conditions, concurrency, failure modes
- **Mode**: resolved-plan; review only
- **Model**: gpt-6.1-sol
- **Reasoning**: max
- **Review Attempt**: 1
- **Stage Attempt**: 1
- **Tier / Gate**: High; ordinary maximum of two attempts per stage, including the initial attempt. No stage-3 exception. Aggregation and advancement remain caller-owned.
- **Artifact**: `reviews/task-terraform-folder-fixtures-attempt-1-stage-attempt-1-review-3-gpt-6.1-sol.md`
- **Issues Found**: 0
- **Verdict**: CLEAN

This verdict covers this reviewer slot and the current public diff, not other reviewers or stages. Source was inspected independently; no previous review or private evidence was consulted.

## Evidence Checklist

- [x] Traced the acceptance entrypoints into the shared CRUD builders: [MLModel, lines 240–292](../internal/services/mlmodel/resource_ml_model_test.go#L240-L292) and [MLExperiment, lines 214–267](../internal/services/mlexperiment/resource_ml_experiment_test.go#L214-L267). Each builder allocates two folders once, reuses both HCL blocks in both steps, and changes the item's reference from the first folder's ID to the second folder's ID.
- [x] Checked fixture identity and dependency construction against the existing [folder helper, lines 162–178](../internal/testhelp/testacc.go#L162-L178), [random-name helper, lines 23–32](../internal/testhelp/utils.go#L23-L32), and [reference helper, lines 190–192](../internal/testhelp/utils.go#L190-L192). The new configs do not obtain folder IDs from well-known resources.
- [x] Inspected the stale-fixture and cross-invocation regression assertions in [MLModel, lines 47–111](../internal/services/mlmodel/folder_fixtures_test.go#L47-L111) and [MLExperiment, lines 47–111](../internal/services/mlexperiment/folder_fixtures_test.go#L47-L111). They check workspace scope, two retained folder resources, an ID traversal to an owned folder, different move targets, and distinct folder addresses between builder invocations.
- [x] Examined lifecycle-test isolation and synchronization in [MLModel, lines 20–45](../internal/services/mlmodel/folder_fixtures_test.go#L20-L45) and [MLExperiment, lines 20–45](../internal/services/mlexperiment/folder_fixtures_test.go#L20-L45), together with the [unit-case factory, lines 66–86 and 113–140](../internal/testhelp/testunit.go#L66-L140). Each lifecycle test uses fresh fake stores and one mutex-protected transport shared by its provider factory.
- [x] Verified the identical workflow-level group and disabled automatic cancellation in [Test, lines 23–27](../.github/workflows/test.yml#L23-L27) and [Cleanup, lines 10–14](../.github/workflows/cleanup-script.yml#L10-L14). Inspected the retained [test matrix](../.github/workflows/test.yml#L310-L349) and the [workflow regression, lines 16–53](../internal/testhelp/folder_fixtures_test.go#L16-L53).
- [x] Compared cleanup expectations with the existing [folder Delete implementation, lines 235–265](../internal/services/folder/resource_folder.go#L235-L265) and the documented [environment-lock boundaries, lines 353–367](../DEVELOPER.md#L353-L367). Neither the destructive cleanup implementation nor product API/retry behavior is changed by this diff.
- [x] Independently executed 15 passing unit tests through the required [Task entrypoint, lines 274–291](../Taskfile.yml#L274-L291), including all five new tests, and successfully linted both workflows. SHA-256 comparisons confirmed all eight reviewed change files plus `go.sum` and `Taskfile.yml` remained unchanged; `git diff --check` passed.
- [ ] Live acceptance, cloud operations, destructive cleanup, and actual GitHub queue/cancellation behavior were not exercised; these were outside the authorized review scope.

## Clean Assessment

1. **Stale or missing shared folders no longer affect these positive CRUD configurations.** Both builders call the existing owned-folder helper rather than reading `FolderRS1` or `FolderRS2`. Folder HCL and addresses are captured once per invocation, not regenerated per step. The target-folder references therefore establish creation dependencies without binding either case to another run's folders.

2. **The move does not destroy the source folder mid-test.** Both folder resources remain in the second configuration. Only the item's folder reference changes, and the attribute-pair checks compare against the appropriate managed folder. The final state retains the target-folder dependency for normal Terraform destruction. The now-empty source folder can be destroyed independently without requiring a new manual teardown path. Normal harness teardown remains in place; successful unit execution is not a guarantee of recovery after process termination or a real service failure.

3. **The new fake-backed lifecycle coverage does not introduce a shared-store race or an evident lock cycle.** Each case owns its server instances rather than mutating the package-global fake server. Separating folder storage avoids passing non-item folders through the generic item conversion code. The same transport mutex covers concurrent requests from that case, including folder creation and teardown. The inspected folder/item fake callbacks operate on in-memory handlers and do not recursively call this transport. This is a test-only serialization measure, not evidence that production concurrency or every failure path is race-free.

4. **Workflow exclusion is scoped consistently with the documented operational contract.** The root-level lock covers the entire workflow, so matrix legs retain their internal concurrency while another participating Test or Cleanup workflow cannot overlap the active run. The one-pending-run replacement behavior is expressly documented; this is not a FIFO scheduler. Manual cancellation, timeouts, runner loss, and cleanup initiated locally or outside the participating repository workflows remain limitations, not defects newly established by this diff. The lock does not make those paths transactional or guarantee cleanup following interruption.

No actionable introduced robustness defect was found.

## Validation — Executed Versus Reported

### Independently executed

Used the existing isolated toolchain, Go **1.27.1**, Terraform **1.16.6**, and Git Bash on Windows. Every test ran through `task testunit` with **`TF_ACC=0`**, Task's preview configuration, and `GOFLAGS="-mod=readonly -count=1"`. Compiler/module caches and temporary work stayed within the isolated toolchain area. No global installation or raw test invocation was used.

| Command | Observed result |
| --- | --- |
| `task testunit -- ".*" github.com/microsoft/terraform-provider-fabric/internal/services/mlmodel` | **PASS — 7 tests**, including both new fixture tests |
| `task testunit -- ".*" github.com/microsoft/terraform-provider-fabric/internal/services/mlexperiment` | **PASS — 7 tests**, including both new fixture tests |
| `task testunit -- FolderFixtures github.com/microsoft/terraform-provider-fabric/internal/testhelp` | **PASS — 1 test**, the new workflow-isolation regression |
| `actionlint .github\workflows\test.yml .github\workflows\cleanup-script.yml` | **PASS — exit 0** |
| `git diff --check` | **PASS — exit 0** |
| Before/after SHA-256 comparison of the eight reviewed files, `go.sum`, and `Taskfile.yml` | **Unchanged — 10/10** |

The initial Task invocation executed **zero tests** because the default Bash launcher selected an unconfigured WSL installation. Selecting the existing Git Bash in the process-local executable search path resolved that runner issue. All three subsequent scoped Task invocations completed successfully against unchanged source.

### Reported, not independently reproduced

The supplied context reports pre-fix failing HCL/workflow regressions and a clean dependency-tidying comparison. This review did not recreate the pre-fix tree or rerun dependency tidying. The green unit results and workflow lint above were independently repeated rather than inferred from those claims.

### Limits and follow-up

These checks used one local Terraform version, not the complete Terraform/OpenTofu workflow matrix. No race-detector run, prolonged concurrency stress, injected transport failures, or live eventual-consistency/deletion behavior was tested. Actual service permissions, interruption recovery, and end-to-end deployment-environment behavior require separately authorized acceptance validation. No source changes, commits, pushes, publication, or external-state changes were performed.
