# spncheck

`spncheck` verifies that every Terraform resource/data source is correctly
marked as **supporting service principal (SPN) authentication vs. not**. It
cross-checks each service package's declared support (`IsSPNSupported` in
`base.go`'s `ItemTypeInfo`) against the actual service-principal support of the
[`fabric-sdk-go`](https://github.com/microsoft/fabric-sdk-go) APIs the package
calls.

It shares its scanning engine with [`previewcheck`](../previewcheck/): both tools
resolve SDK call sites, read SDK doc comments, and compare the result against a
declared `ItemTypeInfo` flag via the common
[`tools/internal/toolutil`](../internal/toolutil/) package. Only the marker they
look for (and the declared field they compare) differ.

## How it works

A package should be `IsSPNSupported = true` only when **every** SDK client
function it invokes supports service principal. Each SDK client function
documents its supported identities in a table under this heading in the doc
comment directly above the function:

```text
MICROSOFT ENTRA SUPPORTED IDENTITIES
| Identity | Support | |-|-| | User | Yes | | Service principal ... | Yes |
```

The **Service principal** row's Support cell is one of:

- `Yes` — the API supports service principal
- `No` — the API does **not** support service principal
- a conditional sentence — supported only under stated conditions

`spncheck` resolves each SDK call site to the exact SDK function (using the Go
type checker, so import aliases and receivers resolve correctly), reads its
identity table, and flags the package when **any** called API has a hard `No`.

## Run it

```sh
go run ./tools/spncheck                  # report findings, exit 1 if any
go run ./tools/spncheck -dir DIR         # scan a different services directory
go run ./tools/spncheck -exclusions PATH # use a specific exclusions file
```

Exit codes: `0` = all consistent, `1` = mismatch (or stale exclusion) found,
`2` = error.

## Reading the output

Findings are split by confidence, because the SDK's identity tables are not
always present (a missing table does **not** prove an API supports SPN):

- **OVER-MARKED** (high confidence) — declared `IsSPNSupported = true`, but a
  called SDK API is documented as **not** supporting service principal. The item
  **should** be `false`. Fix these.
- **REVIEW** (low confidence) — declared `IsSPNSupported = false`, but every
  called API supports service principal. *Possibly* promotable to `true`, but
  confirm manually.
- **EXCLUDED** — a failing item suppressed via `exclusions.yaml` (see below).
- **STALE EXCLUSIONS** — an excluded item that is no longer a mismatch; remove
  the entry. Stale entries also fail the run (exit `1`).
- **UNDETERMINED** — no `fabric-sdk-go` calls found (e.g. generic `fabricitem`
  resources whose CRUD runs through the shared abstraction). Informational.

## Acting on a failure

- **OVER-MARKED** → set `IsSPNSupported = false` on the item's `ItemTypeInfo` in
  its `base.go`, **or** — if the item genuinely supports SPN and the SDK
  annotation is stale/conservative — exclude it (see below).
- **REVIEW** → manually verify against Fabric docs; flip `IsSPNSupported` to
  `true` only if every operation truly supports service principal.

If an item calls an SDK package whose directory can't be derived from the item
name automatically, add a mapping to `sdkPackageOverrides` in
[`main.go`](./main.go).

## Excluding items with stale SDK annotations

The SDK sometimes annotates an API more conservatively than reality. That makes
an item show up as **OVER-MARKED** even though keeping it `true` is correct. List
such items in [`exclusions.yaml`](./exclusions.yaml) with a reason:

```yaml
exclusions:
  - service: someitem
    reason: SPN-supported; someclient.DoThing carries a stale SDK identity annotation
```

`service` is the service package name (the value shown in the report). An
excluded item moves to the **EXCLUDED** section instead of failing the run.
Exclusion is deliberately narrow — a stale entry (the item is no longer a
mismatch, e.g. the SDK corrected the annotation) is reported as **STALE** so it
can be removed. Note that a service-level exclusion also suppresses *future*
genuine non-SPN APIs that item may start calling, so keep the list minimal and
reasoned.

The tool's unit tests run via `task test:gaptools`.
