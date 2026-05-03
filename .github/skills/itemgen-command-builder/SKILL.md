---
name: itemgen-command-builder
description: Given SDK analysis results, automatically determine the correct itemgen archetype and build the full go run tools/itemgen/main.go command. USE FOR: scaffolding new Fabric Item resources using the itemgen code generator. Only applies to Fabric Item resources (not bespoke resources like Connection, Gateway, Workspace).
---

# Skill: Itemgen Command Builder

Given SDK analysis results (from `#skill:sdk-contract-navigator`), automatically determine the correct `itemgen` archetype and build the full `go run tools/itemgen/main.go` command.

> **Important:** This skill applies ONLY to Fabric Item resources (Category A from `#skill:sdk-contract-navigator`). Non-item resources (Connection, Shortcut, Gateway, Workspace, etc.) do NOT use `itemgen` — they require manual bespoke implementation.

## Prerequisites

- SDK analysis has been completed (from `#skill:sdk-contract-navigator`)
- The resource is confirmed as a Fabric Item (not a non-item resource)

## Step 1 — Determine the Archetype

Use the SDK analysis to select the correct archetype. Refer to the **"Item Archetypes"** table in `.github/instructions/fabric-item-patterns.instructions.md` for the archetype capabilities matrix.

Also read `tools/itemgen/main.go` for the canonical list of valid item types from the `validItemTypes()` function.

### How to Check Each Capability

- **Has Properties** → The SDK Get response main struct has a `Properties` field pointing to a named struct type (e.g. `fablakehouse.Properties`)
- **Has CreationPayload** → A `CreationPayload` struct exists in the SDK package
- **Has Definition** → The items client has `Get<Item>Definition()` and/or `Update<Item>Definition()` methods

## Step 2 — Gather Flag Values

The `itemgen` tool accepts 9 command-line flags. Determine each value from the SDK analysis and Fabric API docs:

| Flag                 | Type   | How to Determine                                                                                                    | Default        |
| -------------------- | ------ | ------------------------------------------------------------------------------------------------------------------- | -------------- |
| `-item-name`         | string | Display name with spaces (e.g. `"Data Pipeline"`, `"Eventhouse"`)                                                   | **required**   |
| `-items-name`        | string | Plural form (e.g. `"Data Pipelines"`, `"Eventhouses"`)                                                              | **required**   |
| `-item-type`         | string | Archetype from Step 1                                                                                               | **required**   |
| `-definition-path`   | string | The definition file path from the issue's "Definition Paths" field (e.g. `"definition.json"`, `"eventstream.json"`) | `content.json` |
| `-rename-allowed`    | bool   | Check SDK for Update/Rename method on the items client                                                              | `true`         |
| `-is-preview`        | bool   | Check Fabric API docs for "preview" badge or header                                                                 | `false`        |
| `-is-spn-supported`  | bool   | Check API docs for service principal authentication support                                                         | `false`        |
| `-generate-fakes`    | bool   | set to `true` unless item is of archetype `basic` or `definition` — generates fake test handlers                    | `true`         |
| `-generate-examples` | bool   | Always set to `true` — generates TF example files                                                                   | `true`         |

### Flag Value Details

**`-item-name`**: The human-readable display name. Use the form from Microsoft docs (e.g. "Data Pipeline" not "DataPipeline"). The tool derives:

- `Package` = lowercased, no spaces (e.g. `datapipeline`)
- `Type` = lowercased, spaces→underscores (e.g. `data_pipeline`)
- `TypeInfo` = no spaces (e.g. `DataPipeline`)

**`-items-name`**: The plural form. Usually just append "s" but check API docs for irregular plurals (e.g. "KQL Databases", "Warehouses", "Variable Libraries").

**`-definition-path`**: The definition file path as listed in the issue's "Definition Paths" section (populated by `#skill:issue-composer` from the Fabric definition article). Use the first/primary definition path (e.g. `"eventstream.json"`, `"definition.json"`, `"notebook-content.ipynb"`). This determines the definition key used in Terraform HCL blocks and template source references.

> **Note:** This flag is only relevant for item types that have a definition (`definition`, `definition-properties`, `config-definition-properties`). For archetypes without a definition (`basic`, `properties`, `config-properties`), omit this flag — it will be ignored. If the item archetype includes a definition but the "Definition Paths" field is missing from the issue, **prompt the user** to provide the definition file path before proceeding.

**`-rename-allowed`**: Most items support rename. Set to `false` if the SDK items client lacks an `Update<ItemName>` method.

**`-is-preview`**: Fetch the Create item API docs page and check for the preview note. The URL pattern is:

```
https://learn.microsoft.com/rest/api/fabric/<itemlowercase>/items/create-<item-kebab-case>
```

To check, run:

```bash
curl -sL -H 'User-Agent: Mozilla/5.0' \
  'https://learn.microsoft.com/en-us/rest/api/fabric/lakehouse/items/create-lakehouse?tabs=HTTP' \
  | grep -ci 'currently in preview'
```

If the output is `0`, the item is NOT in preview (`-is-preview=false`). If greater than `0`, the item IS in preview (`-is-preview=true`).

The preview note appears in the HTML as:

```html
<p>{ItemType} item is currently in Preview (<a href="...">learn more</a>).</p>
```

**`-is-spn-supported`**: Check if the API documentation mentions service principal support. Also check if the existing `base.go` similar items use `IsSPNSupported: true`.

## Step 3 — Build the Command

Construct the full command:

```bash
go run tools/itemgen/main.go \
  -item-name "<Display Name>" \
  -items-name "<Plural Display Name>" \
  -item-type "<archetype>" \
  -definition-path "<definition-file-path>" \
  -rename-allowed=<true|false> \
  -is-preview=<true|false> \
  -is-spn-supported=<true|false> \
  -generate-fakes=true \
  -generate-examples=true
```

## Reference

- Itemgen source: `tools/itemgen/main.go`
- Template directory: `tools/itemgen/templates/`
- Canonical example output: `internal/services/lakehouse/`
