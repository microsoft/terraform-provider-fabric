# previewcheck

`previewcheck` verifies that every Terraform resource/data source is correctly
marked as **preview vs. GA**. It cross-checks each service package's declared
preview status (`IsPreview` in `base.go`'s `ItemTypeInfo`) against the actual
preview status of the [`fabric-sdk-go`](https://github.com/microsoft/fabric-sdk-go)
APIs the package calls.

## How it works

A package should be `IsPreview = true` when **any** SDK client function it
invokes is documented as preview. The SDK flags a preview API with one of these
phrases in the doc comment directly above the function:

- `is currently in preview`
- `is part of a Preview release`

`previewcheck` resolves each SDK call site to the exact SDK function (using the
Go type checker, so import aliases and receivers resolve correctly), reads its
doc comment, and compares the result against the declared `IsPreview` value.

## Run it

```sh
go run ./tools/previewcheck                  # report findings, exit 1 if any
go run ./tools/previewcheck -dir DIR         # scan a different services directory
go run ./tools/previewcheck -exclusions PATH # use a specific exclusions file
```

Exit codes: `0` = all consistent, `1` = mismatch (or stale exclusion) found,
`2` = error.

## Reading the output

Findings are split by confidence, because the SDK's preview annotations are
incomplete (a missing marker does **not** guarantee an API is GA):

- **UNDER-MARKED** (high confidence) — declared GA, but a called SDK API is
  flagged preview. The item **should** be marked preview. Fix these.
- **REVIEW** (low confidence) — declared preview, but no SDK preview marker was
  found on any called API. *Possibly* demotable to GA, but confirm manually.
- **EXCLUDED** — a failing item suppressed via `exclusions.yaml` (see below).
- **STALE EXCLUSIONS** — an excluded item that is no longer a mismatch; remove
  the entry. Stale entries also fail the run (exit `1`).
- **UNDETERMINED** — no `fabric-sdk-go` calls found (e.g. generic `fabricitem`
  resources whose CRUD runs through the shared abstraction). Informational.

## Acting on a failure

- **UNDER-MARKED** → set `IsPreview = true` on the item's `ItemTypeInfo` in its
  `base.go`, **or** — if the item is genuinely GA and the SDK marker is stale —
  exclude it (see below).
- **REVIEW** → manually verify against Fabric docs; flip `IsPreview` to `false`
  only if the API is truly GA.

If an item calls an SDK package whose directory can't be derived from the item
name automatically, add a mapping to `sdkPackageOverrides` in
[`main.go`](./main.go).

## Excluding GA items with stale SDK markers

The SDK sometimes leaves a preview marker on an API that is actually GA. That
makes an item show up as **UNDER-MARKED** even though keeping it GA is correct.
List such items in [`exclusions.yaml`](./exclusions.yaml) with a reason:

```yaml
exclusions:
  - service: tenantsetting
    reason: GA; admin.UpdateTenantSetting still carries a stale SDK preview marker
```

`service` is the service package name (the value shown in the report). An
excluded item moves to the **EXCLUDED** section instead of failing the run.
Exclusion is deliberately narrow — a stale entry (the item is no longer a
mismatch, e.g. the SDK removed the marker) is reported as **STALE** so it can be
removed. Note that a service-level exclusion also suppresses *future* genuine
preview APIs that item may start calling, so keep the list minimal and reasoned.

## CI

The `test-gaptools` job in `.github/workflows/test.yml` runs this tool's unit
tests (via `task test:gaptools`) whenever `tools/**` changes.
