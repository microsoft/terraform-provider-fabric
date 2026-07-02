# gapcheck

`gapcheck` finds control-plane Fabric entities exposed by
[`fabric-sdk-go`](https://github.com/microsoft/fabric-sdk-go) that have **no
corresponding Terraform resource** in this provider — so we can spot SDK
capabilities the provider does not yet cover.

## How it works

1. Scans the SDK's `*_client.go` files and collects every exported client method
   (e.g. `core/WorkspacesClient.CreateWorkspace`).
2. Scans this provider's service packages to learn which SDK methods are already
   called (i.e. already covered by a resource).
3. Reports every SDK method that is **not** covered, **not** an action-style
   verb (e.g. `Cancel`, `Assign`), and **not** listed in `exclusions.yaml`.

The report is intentionally **exhaustive**: it lists read-only/list-only/
update-only methods too, not just missing Create/Delete lifecycles. Everything
that isn't genuinely a resource gap is pruned explicitly via `exclusions.yaml`,
which keeps the list honest and reviewable.

## Run it

```sh
go run ./tools/gapcheck            # human-readable report
go run ./tools/gapcheck -json      # machine-readable JSON
go run ./tools/gapcheck -dir DIR   # scan a different services directory
go run ./tools/gapcheck -exclusions PATH   # use a specific exclusions file
```

Exit codes: `0` = no gaps, `1` = gaps or stale exclusions found, `2` = error.

## Reading the output

- **GAPS** — SDK methods with no TF coverage and no exclusion. Each is a
  candidate to either implement or exclude.
- **EXCLUDED** — count of methods suppressed by `exclusions.yaml`.
- **STALE EXCLUSIONS** — methods that are now implemented but still listed in
  `exclusions.yaml`; remove them. Stale entries also fail the run (exit `1`).

## Acting on a failure

For each reported gap, do **one** of:

- **Implement it** as a Terraform resource/data source (the tool will then
  auto-detect the SDK calls and stop reporting it), **or**
- **Exclude it** in [`exclusions.yaml`](./exclusions.yaml) with a reason.

Exclusion format — `package/client/method → reason`. Omit `method` (or set it to
`"*"`) to exclude every method on a client:

```yaml
exclusions:
  - package: core
    client: CatalogClient
    reason: read-only catalog search, no CRUD lifecycle
  - package: core
    client: WorkspacesClient
    method: ProvisionIdentity
    reason: action-style operation, not a resource
```

You do **not** need to exclude methods that are already used by an existing
resource — those are detected automatically from the provider source.

Run the audit locally with the commands above before adding new SDK-backed
resources. The tool's unit tests run via `task test:gaptools`.
