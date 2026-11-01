---
applyTo: "examples/**/*.tf,examples/**/import.sh"
---

# Example (.tf) Review Checklist — terraform-provider-fabric

## Purpose & Scope

Guidance for Copilot code review of Terraform HCL examples under `examples/`.
These examples must be valid, appliable configurations — they are also the
source `task docs` renders into `docs/`, so a defect here surfaces in the docs
too. Flag the recurring issues below.

## Schema-Valid Values

Example values must pass the provider's schema validators. Flag:

- `tokens_delimiter` not one of `{{}}`, `<<>>`, `@{}@`, `____`. A common bad
  value is `"##"`.

```hcl
// Avoid — rejected by schema
tokens_delimiter = "##"
// Prefer
tokens_delimiter = "<<>>"
```

- `processing_mode` not exactly `GoTemplate`, `Parameters`, or `None`
  (case-sensitive — `"parameters"` fails).
- A `definition` (resource) or `output_definition = true` (data source) set
  without the required `format`. Flag the missing `format`, and flag a `format`
  value not allowed by that item's schema. Most items accept only `Default`,
  but some accept only specific formats — e.g. `fabric_report`
  (`PBIR` / `PBIR-Legacy`), `fabric_semantic_model` (`TMSL` / `TMDL`),
  `fabric_notebook` (`ipynb` / `py`). Verify against the item's schema, not the
  SDK enum names.

## Fixture References

- `definition` sources reference fixture files (e.g.
  `"${local.path}/Foo.json"`). Flag a `source` file name that does not match a
  key/path in the resource's `definition` block, and references that clearly
  point outside the example's fixtures directory
  (`internal/testhelp/fixtures/<item>/`).
- Flag use of `${local.path}` when the example module has no `locals { path = ... }`
  block (defined in `providers.tf`). Either add the block or use `${path.module}`.

## Naming & Consistency

- The example folder name must match the actual resource/data source type
  (e.g. `examples/resources/fabric_workspace_outbound_cloud_connection_rules/`,
  not a shortened `..._outbound_connection_rules`).
- `import.sh` comment and command must reference the correct type and ID
  placeholders (e.g. `fabric_cosmos_db` with `<CosmosDBID>`, not a copy/pasted
  `fabric_map` / `<MapID>`).
- Flag copy/paste artifacts: duplicated blocks, stray comment fragments, and
  output names with doubled prefixes (e.g. `example_example_custom_delimiter`).

## HCL Correctness

- Object attributes use `=`, not `:` (e.g. `connection_type = "..."`, not
  `connection_type : "..."`).
- Placeholder URLs/IDs should be realistic and generic (e.g. a
  `00000000-...` UUID), not a misleading value presented as required
  (e.g. a non-tenant-specific `https://microsoft.sharepoint.com`).
