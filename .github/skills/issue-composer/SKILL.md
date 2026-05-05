---
name: issue-composer
description: Compose a GitHub issue for a new resource, data source, ephemeral resource, or enhancement to an existing resource, pre-filling the appropriate template from Fabric API documentation and SDK analysis. USE FOR: creating well-structured GitHub issues for new Terraform resources or enhancements. INVOKES: GitHub MCP tools for issue creation, sdk-contract-navigator skill for SDK analysis.
---

# Skill: Issue Composer

Compose a GitHub issue for a new resource, data source, ephemeral resource, or enhancement to an existing resource, pre-filling the appropriate template from Fabric API documentation and SDK analysis.

## Prerequisites

- The user has described what they want to add or change
- Ideally, `#skill:sdk-contract-navigator` has already been run to identify SDK availability
- If the caller (issue-creator agent or user) has not provided a milestone name, this skill should prompt for one, then resolve it to a numeric milestone ID via the GitHub API

## Step 1 — Resolve Milestone Name to Number

If a milestone name was provided (e.g. "2026-04"), resolve it to its numeric ID:

1. Query the GitHub REST API:

   ```
   GET /repos/microsoft/terraform-provider-fabric/milestones?state=all&per_page=100
   ```

2. Parse the response and find the milestone object where `title` matches the user's input (case-insensitive)

3. Extract the `number` field from that milestone object

4. If no match is found, report an error with the list of available milestones and ask the user to retry

If no milestone was specified, set `milestone: null` and proceed.

## Step 2 — Determine Issue Type

Choose the correct issue template based on what is being requested.

| Scenario                                   | Title Prefix | Label            | Template File                               |
| ------------------------------------------ | ------------ | ---------------- | ------------------------------------------- |
| **New resource** (Fabric Item or non-item) | `[RS]`       | `tf/resource`    | `tfprovider_resource_request.yml`           |
| **New data source**                        | `[DS]`       | `tf/data-source` | `tfprovider_data_source_request.yml`        |
| **New ephemeral resource**                 | `[EPH]`      | `tf/ephemeral`   | `tfprovider_ephemeral_resource_request.yml` |
| **Enhancement to existing resource**       | `[FEAT]`     | `feature`        | `feature_request.yml`                       |

### When to Use Each

- **`[RS]` / `[DS]`** — A completely new Terraform resource or data source that doesn't exist yet. Applies to both Fabric Items (Lakehouse, Eventhouse, etc.) and non-item resources (Connection, Shortcut, Gateway, Workspace Role Assignment, etc.)
- **`[EPH]`** — A new ephemeral resource (short-lived, not stored in state)
- **`[FEAT]`** — Adding new attributes to an existing resource, changing behavior, adding support for a new API feature on an existing resource, or any other enhancement that modifies existing code

If both a resource and data source are needed for the same item, create **two** separate issues.

## Step 3 — Extract Resource Details

From the user's description and SDK analysis, gather:

| Detail                | How to Determine                                                                                                             |
| --------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| **Resource name**     | `fabric_<snake_case>` — e.g. `fabric_lakehouse`, `fabric_connection`, `fabric_shortcut`                                      |
| **Display name**      | PascalCase with spaces — e.g. "Lakehouse", "Connection", "Shortcut"                                                          |
| **API doc links**     | Browse `learn.microsoft.com/rest/api/fabric/` for API pages                                                                  |
| **SDK availability**  | Check if the SDK package/client exists (from `#skill:sdk-contract-navigator`)                                                |
| **Resource category** | Fabric Item (~80%) or Non-Item (~20%) — see below                                                                            |
| **Item archetype**    | For Fabric Items only: basic, definition, properties, definition-properties, config-properties, config-definition-properties |
| **Complexity**        | `easy` (basic/definition), `moderate` (properties/non-item), `hard` (config-\*/complex non-item)                             |
| **Definition paths**  | For Fabric Items with definitions: fetch the definition article to list supported definition parts — see below               |
| **Related resources** | Existing resources that interact with this one                                                                               |
| **Preview status**    | Check if the Fabric API is marked as preview                                                                                 |
| **SPN support**       | Check if the API supports service principal authentication                                                                   |

### Resource Category Identification

**Fabric Items** — Standard items managed in workspaces (Lakehouse, Eventhouse, SQL Database, Data Pipeline, Notebook, etc.). These use the `fabricitem` generic abstraction and have `ItemType` constants in the core SDK.

**Non-Item Resources** — Specialized resources with bespoke CRUD logic. Each belongs to an **implementation pattern (A–H)** that determines canonical reference, lifecycle semantics, and test structure:

| Pattern | Characteristic                                       | SDK Client                         | Canonical Reference                    |
| :-----: | ---------------------------------------------------- | ---------------------------------- | -------------------------------------- |
|  **A**  | Workspace policy singleton (no ID, delete=reset)     | `fabcore.WorkspacesClient`         | `internal/services/workspacencp/`      |
|  **B**  | Workspace settings (dedicated client, validators)    | `fabspark.*`                       | `internal/services/sparkwssettings/`   |
|  **C**  | Role assignment (parent+principal+role, ImportState) | Various `*Client`                  | `internal/services/workspacera/`       |
|  **D**  | Batch assignment (immutable set, no update)          | `fabadmin.DomainsClient`           | `internal/services/domainra/`          |
|  **E**  | Standalone entity, standard CRUD                     | Various dedicated clients          | `internal/services/workspace/`         |
|  **E**  | Standalone entity, polymorphic types                 | `fabcore.GatewaysClient`           | `internal/services/gateway/`           |
|  **E**  | Standalone entity, tenant-scoped                     | `fabadmin.DomainsClient`           | `internal/services/domain/`            |
|  **E**  | Standalone entity, connect/disconnect lifecycle      | `fabcore.GitClient`                | `internal/services/workspacegit/`      |
|  **F**  | Item-scoped (workspace_id+item_id, 3+ path params)   | `fabcore.OneLakeShortcutsClient`   | `internal/services/shortcut/`          |
|  **F**  | Item-scoped (ModifyPlan, conditional validation)     | `fabcore.JobSchedulerClient`       | `internal/services/itemjobscheduler/`  |
|  **F**  | Item-scoped (simple CRUD, no Update)                 | `fabcore.ExternalDataSharesClient` | `internal/services/externaldatashare/` |
|  **G**  | Tenant-level, custom identity/delete semantics       | `fabadmin.TenantsClient`           | `internal/services/tenantsetting/`     |
|  **H**  | Complex: dual clients, write-only secrets, KV refs   | `fabcore.ConnectionsClient`        | `internal/services/connection/`        |

**Pattern classification decision tree:**

```
Is it a workspace policy/settings with no real entity ID?
├── YES → Uses WorkspacesClient sub-endpoint, delete=reset? → Pattern A
│         Uses dedicated Spark/Environment client, ConfigValidators? → Pattern B
│
├── NO → Is it an assignment of principals/items to a parent?
│         ├── Single-item assignment with updatable role, ImportState? → Pattern C
│         └── Batch set assignment, fully immutable, no Import? → Pattern D
│
├── NO → Is it scoped to a specific item (workspace_id + item_id)?
│         └── YES → Pattern F
│
├── NO → Is it a tenant-level admin resource with non-UUID identity or custom delete?
│         └── YES → Pattern G
│
├── NO → Does it have write-only secrets, dual clients, KV references?
│         └── YES → Pattern H
│
└── NO → Standalone entity with dedicated client → Pattern E
```

When composing the issue, include the pattern letter in the "Details / References" section so downstream agents can immediately route to the correct implementation reference.

### Item Definition Paths

For Fabric Items that support definitions (archetypes: `definition`, `definition-properties`, `config-definition-properties`), fetch the item definition article to discover the supported definition parts/paths:

```
https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/<item-kebab-case>-definition
```

For example:

- Notebook → `https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/notebook-definition`
- Report → `https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/report-definition`
- Spark Job Definition → `https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/spark-job-definition-definition`

Use the `fetch_webpage` tool to read the article and extract the list of definition paths (e.g. `notebook-content.py`, `definition.pbir/report.json`). Each definition part has a path and format. Record all supported definition paths for inclusion in the issue.

If the article does not exist or returns a 404, the item likely does not support definitions — adjust the archetype accordingly.

### API Documentation URL Pattern

Fabric REST API docs follow this pattern (never include `en-us` locale):

```
https://learn.microsoft.com/rest/api/fabric/<service>/items
```

For non-item resources, the API path varies:

```
https://learn.microsoft.com/rest/api/fabric/core/connections
https://learn.microsoft.com/rest/api/fabric/core/shortcuts
https://learn.microsoft.com/rest/api/fabric/core/gateways
https://learn.microsoft.com/rest/api/fabric/core/workspaces
```

## Step 4 — Compose Issue Title

| Type            | Format                                      |
| --------------- | ------------------------------------------- |
| New resource    | `[RS] fabric_<snake_case_name>`             |
| New data source | `[DS] fabric_<snake_case_name>`             |
| New ephemeral   | `[EPH] fabric_<snake_case_name>`            |
| Enhancement     | `[FEAT] <short description of enhancement>` |

## Step 5 — Compose Issue Body

### For New Resources (`[RS]`) and Data Sources (`[DS]`)

#### 📝 Description

Use Job Story format:

```
When managing Microsoft Fabric infrastructure as code,
I want to create/manage <ResourceName> resources via Terraform,
so I can automate provisioning and maintain consistent <ResourceName> configurations across environments.
```

#### 🔬 Details / References

**For Fabric Items:**

> **Important:** Do NOT include SDK CRUD method signatures (e.g. `Get`, `Create`, `Update`, `Delete`, `List`, `GetDefinition`, `UpdateDefinition`) in the issue body. Fabric Items follow standardized method patterns determined by the archetype — listing them adds noise without value. Only include the SDK package, archetype, DTO fields (Properties/CreationPayload), enum types, and definition paths.

```markdown
- Resource Name: `fabric_<snake_case_name>`
- API documentation:
  - https://learn.microsoft.com/rest/api/fabric/<itemtype>/items/create-<item>
  - https://learn.microsoft.com/rest/api/fabric/<itemtype>/items/get-<item>
  - https://learn.microsoft.com/rest/api/fabric/<itemtype>/items/list-<items>
- Definition article: https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/<item-kebab-case>-definition
- SDK Package: `github.com/microsoft/fabric-sdk-go/fabric/<package>`
- Item Archetype: `<archetype>`
- Definition Paths (if applicable):
  - `<path>` (format: `<format>`)
  - `<path>` (format: `<format>`)
- Estimated complexity/effort: <easy|moderate|hard>
- Preview: <yes|no>
- SPN Supported: <yes|no>
- Related resources/data-sources:
  - `fabric_workspace` (required parent)
```

**For Non-Item Resources:**

```markdown
- Resource Name: `fabric_<snake_case_name>`
- API documentation:
  - https://learn.microsoft.com/rest/api/fabric/core/<service>/create-<resource>
  - https://learn.microsoft.com/rest/api/fabric/core/<service>/get-<resource>
  - https://learn.microsoft.com/rest/api/fabric/core/<service>/list-<resources>
- SDK Client: `fabcore.<Resource>Client`
- Resource Category: Non-Item (bespoke CRUD)
- Implementation Pattern: <A|B|C|D|E|F|G|H> — <pattern description>
- Estimated complexity/effort: <easy|moderate|hard>
- Preview: <yes|no>
- SPN Supported: <yes|no>
- Related resources/data-sources:
  - <related resources>
```

> **Pattern key:** A=Workspace policy singleton, B=Workspace settings (Spark), C=Role assignment, D=Batch assignment, E=Standalone entity, F=Item-scoped, G=Tenant-level custom, H=Connection (complex credentials)

#### 🌳 DTO Nesting Depth Map (complex resources only)

**Include this section only for `[RS]` and `[DS]` issues where the SDK DTOs have 3+ nesting levels** (slices containing structs with nested slices/pointers to other structs). Skip for flat resources (Workspace, Domain, role assignments, basic Fabric Items without properties).

Render a tree showing the SDK DTO hierarchy with type annotations. This helps implementors plan model structs, `set()` methods, and `SetNull(ctx)` placement.

Format — use indented tree with type annotations at leaf/branch nodes:

```markdown
<RootDTO>
├── Field1 (string)
├── Field2 (enum)
├── NestedSlice []ChildDTO                    ← Level 1
│   ├── ScalarField (string)
│   ├── DeepSlice []GrandchildDTO             ← Level 2
│   │   ├── LeafSlice []string               ← Level 3
│   │   └── LeafField (enum)
│   └── OptionalNested *AnotherDTO            ← Level 2
│       └── Items []ItemDTO                   ← Level 3
└── OptionalTop *TopDTO                       ← Level 1
    └── Children []ChildDTO                   ← Level 2
```

Rules:

- Annotate slices as `[]Type`, optional nested as `*Type`, scalars as `(type)`
- Mark nesting levels with `← Level N` comments on lines introducing a new struct depth
- Only show fields that map to Terraform schema attributes (skip internal/wire-only fields)
- Maximum depth shown: 5 levels (truncate deeper with `...`)

#### 🚧 Potential Terraform Configuration

Generate sample HCL based on SDK properties discovered:

```terraform
resource "fabric_<snake_case_name>" "example" {
  display_name = "example"
  description  = "Example resource"
  workspace_id = fabric_workspace.example.id

  # Item-specific attributes
}
```

#### ☑️ Acceptance Criteria

```markdown
- [ ] Can create a new <ResourceName> with required attributes
- [ ] Can read <ResourceName> by ID
- [ ] Can update <ResourceName> mutable attributes
- [ ] Can delete <ResourceName>
- [ ] Can import existing <ResourceName>
- [ ] Properties are correctly mapped from SDK response
```

#### ✅ Definition of Done

For Resource:

```markdown
- [ ] Data Transfer Objects (DTOs)
- [ ] Resource Implementation
- [ ] Resource Added to Provider
- [ ] Unit Tests for Happy path
- [ ] Unit Tests for Error path
- [ ] Acceptance Tests
- [ ] Example in the ./examples folder
- [ ] Schema documentation in code
- [ ] Updated auto-generated provider docs with `task docs`
```

For Data Source:

```markdown
- [ ] Data Transfer Objects (DTOs)
- [ ] Data-Source Implementation
- [ ] Data-Source Added to Provider
- [ ] Unit Tests for Happy path
- [ ] Unit Tests for Error path
- [ ] Acceptance Tests
- [ ] Example in the ./examples folder
- [ ] Schema documentation in code
- [ ] Updated auto-generated provider docs with `task docs`
```

### For Enhancements (`[FEAT]`)

Use the `feature_request.yml` template structure.

#### 🚀 Feature description

Job Story format describing the enhancement:

```
When using the existing `fabric_<resource>` resource,
I want to <what is missing or needs to change>,
so I can <expected outcome>.
```

#### 🔈 Motivation

Explain why the enhancement is needed — e.g. new API capability, missing attribute, user request.

#### 🛰 Alternatives

Describe workarounds or alternative approaches considered.

#### 🔬 SDK Diff (from `#skill:sdk-contract-navigator`)

Compare the current SDK DTOs against the existing implementation to identify what's new. Include this section in the issue body so the implementor agent has a clear change list:

**For Fabric Item enhancements:**

```markdown
### SDK Diff

Resource: `fabric_<snake_case_name>`
Category: Fabric Item (`<archetype>`)
Service package: `internal/services/<package>/`

| Change    | SDK Field                     | Go Type                  | Current Status                    |
| --------- | ----------------------------- | ------------------------ | --------------------------------- |
| + New     | `Properties.<FieldName>`      | `*string`                | Not in `<item>PropertiesModel`    |
| + New     | `Properties.<NestedField>`    | `*<NestedDTO>`           | New sub-model needed              |
| + New     | `CreationPayload.<FieldName>` | `*bool`                  | Not in `<item>ConfigurationModel` |
| ~ Changed | `Properties.<FieldName>`      | `*int32` (was `*string`) | Type change in model              |
```

**For Non-Item enhancements:**

```markdown
### SDK Diff

Resource: `fabric_<snake_case_name>`
Category: Non-Item (bespoke CRUD)
Service package: `internal/services/<package>/`
SDK Client: `fabcore.<Resource>Client`

| Change       | SDK Field                   | Go Type   | Current Status                          |
| ------------ | --------------------------- | --------- | --------------------------------------- |
| + New        | `<ResponseDTO>.<FieldName>` | `*string` | Not in `base<Resource>Model`            |
| + New        | `<RequestDTO>.<FieldName>`  | `*bool`   | Not in `request<Action><Resource>Model` |
| + New method | `client.<NewMethod>(...)`   | —         | No CRUD handler                         |
```

#### 🚧 Potential Configuration / Desired Solution

Show HCL demonstrating the desired new behavior:

```terraform
resource "fabric_<existing_resource>" "example" {
  # existing attributes...

  # NEW: proposed enhancement
  new_attribute = "value"
}
```

#### ☑️ Acceptance Criteria

```markdown
- [ ] New attribute `<name>` is added to the schema
- [ ] SDK mapping is correct for new fields
- [ ] `set()` methods updated for new fields
- [ ] Fakes updated with new field test data
- [ ] Existing tests still pass
- [ ] New test assertions for added attributes
- [ ] Documentation updated
```

## Step 6 — Create the Issue

Use the GitHub MCP server `create_issue` tool:

```
owner: microsoft
repo: terraform-provider-fabric
title: <composed title from Step 4>
body: <composed body from Step 5>
labels: [<appropriate label>]
milestone: <resolved milestone number from Step 1, or null if not specified>
```

Label mapping:

- `[RS]` → `["tf/resource"]`
- `[DS]` → `["tf/data-source"]`
- `[EPH]` → `["tf/ephemeral"]`
- `[FEAT]` → `["feature"]`

## Step 7 — Report Back

After creating the issue, report:

- Issue number and URL
- Summary of what was filed
- Resource category (Fabric Item vs Non-Item) and archetype (if applicable)
- Any gaps or unknowns that need follow-up (e.g. API not yet public, SDK package missing)
