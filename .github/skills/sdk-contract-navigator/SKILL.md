---
name: sdk-contract-navigator
description: Navigate the fabric-sdk-go SDK source to find the correct client factory, DTOs, item type constants, and pager methods for a given Fabric Item or non-item resource. USE FOR: identifying SDK packages, client factories, DTO structs, item type constants, and pager methods for any Fabric resource.
---

# Skill: SDK Contract Navigator

Navigate the `fabric-sdk-go` SDK source to find the correct client factory, DTOs, item type constants, and pager methods for a given Fabric Item or non-item resource.

## Prerequisites

- The user has identified a Fabric resource name (e.g. "Lakehouse", "Connection", "Shortcut")

## Step 1 — Check SDK Version and Ensure Local Cache

Read `go.mod` in the provider repo root and find the `github.com/microsoft/fabric-sdk-go` dependency line to determine the current version:

```
github.com/microsoft/fabric-sdk-go <version>
```

Record the exact version string. All SDK browsing must target the version found in `go.mod`.

> **Important:** Always read `go.mod` dynamically — never assume or hardcode the SDK version. It changes with dependency upgrades.

Then ensure the SDK is in the local Go module cache:

```powershell
go mod download github.com/microsoft/fabric-sdk-go@<version>
```

The SDK source is now available at a **deterministic path**:

```
$(go env GOPATH)/pkg/mod/github.com/microsoft/fabric-sdk-go@<version>/
```

> **Primary source:** Always use the local Go module cache for SDK browsing. It is faster and more reliable than GitHub MCP tools (`get_file_contents` returns large JSON blobs that are hard to search; `search_code` has indexing lag). Use `read_file` and `Select-String` (PowerShell) to browse and search SDK files directly.
>
> **Fallback only:** Use GitHub MCP tools or `pkg.go.dev` only if the local cache is unavailable.

### Useful search commands

```powershell
# Find a struct definition
Select-String -Path "<sdk-path>\models.go" -Pattern "^type <StructName> struct"

# Find a field across all structs (with context)
Select-String -Path "<sdk-path>\models.go" -Pattern "<FieldName>" -Context 5

# List all structs matching a pattern
Select-String -Path "<sdk-path>\models.go" -Pattern "^type .*Connection.*struct" | Select-Object -ExpandProperty Line

# List directory contents
Get-ChildItem "<sdk-path>\fabric\" -Directory | Select-Object -ExpandProperty Name
```

## Step 2 — Determine Resource Category

This provider has **two** resource categories with different SDK patterns:

### Category A: Fabric Items (~60% of resources)

These use the generic `fabricitem` abstraction with per-item SDK packages.

**Identifying trait:** The SDK has a dedicated package under `fabric/<itempackage>/` with its own client factory, and the provider service uses `fabricitem.NewResource*` constructors.

**Examples:** Lakehouse, Eventhouse, Data Pipeline, SQL Database, KQL Database, Notebook, Semantic Model, Spark Job Definition, ML Experiment, ML Model, Environment, Warehouse, etc.

### Category B: Non-item resources (~40% of resources)

These use bespoke CRUD implementations with clients from the `fabric/core/` package or other shared packages.

**Identifying trait:** The SDK client is accessed via `fabcore.*Client` (e.g. `fabcore.ConnectionsClient`, `fabcore.ShortcutsClient`, `fabcore.GatewaysClient`). The provider service implements `resource.Resource` directly without `fabricitem` generics.

**Examples:** Connection, Shortcut, Gateway, Workspace, Domain, Workspace Role Assignment, Connection Role Assignment, Gateway Role Assignment, Deployment Pipeline Role Assignment, Workspace Git, etc.

Proceed with **Step 3A** for Fabric Items or **Step 3B** for non-item resources.

---

## Fabric Items Path (Category A)

### Step 3A — Locate the SDK Package

The SDK organizes item packages under `fabric/<packagename>/`.

**Using local module cache** — List directories to find the item package:

```powershell
Get-ChildItem "<sdk-path>\fabric\" -Directory | Select-Object -ExpandProperty Name
```

Look for a directory matching the item name (lowercased, no spaces). Common patterns:

| Item Name            | SDK Package Path             |
| -------------------- | ---------------------------- |
| Lakehouse            | `fabric/lakehouse/`          |
| Eventhouse           | `fabric/eventhouse/`         |
| SQL Database         | `fabric/sqldatabase/`        |
| Data Pipeline        | `fabric/datapipeline/`       |
| KQL Database         | `fabric/kqldatabase/`        |
| Spark Job Definition | `fabric/sparkjobdefinition/` |
| Mirrored Database    | `fabric/mirroreddatabase/`   |
| Copy Job             | `fabric/copyjob/`            |

### Step 4A — Identify the Client Factory

Read the SDK package's client factory file (usually `client_factory.go`) from the local cache. The pattern is:

```go
fab<package>.NewClientFactoryWithClient(fabricClient).NewItemsClient()
```

For example:

- `fablakehouse.NewClientFactoryWithClient(fabricClient).NewItemsClient()`
- `fabeventhouse.NewClientFactoryWithClient(fabricClient).NewItemsClient()`

### Step 5A — Identify CRUD Methods

Read the items client file (usually `items_client.go`) for these methods:

| Method Type    | Pattern                                                           | Example                                               |
| -------------- | ----------------------------------------------------------------- | ----------------------------------------------------- |
| **Get**        | `client.Get<ItemName>(ctx, workspaceID, itemID, nil)`             | `client.GetLakehouse(ctx, wID, iID, nil)`             |
| **List pager** | `client.NewList<ItemNames>Pager(workspaceID, nil)`                | `client.NewListLakehousesPager(wID, nil)`             |
| **Create**     | `client.Create<ItemName>(ctx, workspaceID, payload, nil)`         | `client.CreateLakehouse(ctx, wID, payload, nil)`      |
| **Update**     | `client.Update<ItemName>(ctx, workspaceID, itemID, payload, nil)` | `client.UpdateLakehouse(ctx, wID, iID, payload, nil)` |
| **Delete**     | `client.Delete<ItemName>(ctx, workspaceID, itemID, nil)`          | `client.DeleteLakehouse(ctx, wID, iID, nil)`          |

Also check for definition methods (indicates `definition` or `*-definition-*` archetype):

- `client.Get<ItemName>Definition(ctx, workspaceID, itemID, nil)`
- `client.Update<ItemName>Definition(ctx, workspaceID, itemID, payload, nil)`

### Step 6A — Identify the ItemType Constant

The item type constant lives in the `core` package (`fabric/core/`), typically in `constants.go` or `models.go`:

```go
fabcore.ItemType<ItemName>
```

Examples: `fabcore.ItemTypeLakehouse`, `fabcore.ItemTypeEventhouse`, `fabcore.ItemTypeSQLDatabase`

### Step 7A — Identify DTO Structs

Read the SDK package's `models.go` file for these key structs:

| DTO                  | Purpose                                | How to Identify                                                             |
| -------------------- | -------------------------------------- | --------------------------------------------------------------------------- |
| **Main item struct** | Top-level item (e.g. `Lakehouse`)      | Has fields: `ID`, `DisplayName`, `Description`, `WorkspaceID`, `Properties` |
| **Properties**       | Read-only properties from Get response | Struct named `Properties` with server-computed fields                       |
| **CreationPayload**  | Create-time configuration              | Struct named `CreationPayload` with user-settable fields                    |

Also look for:

- **Enum types** — Any `type <Name> string` with `Possible<Name>Values()` func
- **Nested structs** — Structs referenced as pointer fields within `Properties`

### Step 8A — Determine the Item Archetype

Use the **"Item Archetypes"** table in `.github/instructions/fabric-item-patterns.instructions.md` to match SDK capabilities to the correct archetype.

**How to check:**

- **Has Properties** → The Get response struct has a `Properties` field of a named struct type
- **Has CreationPayload** → A `CreationPayload` struct exists in the package
- **Has Definition** → The items client has `GetDefinition` / `UpdateDefinition` methods

---

## Non-item resources path (Category B)

### Step 3B — Locate the SDK Client

Non-item resources typically use clients from the `fabric/core/` package. For the SDK client mapping, see the **"SDK Client Mapping"** table in `.github/instructions/non-item-patterns.instructions.md`.

Client access pattern:

```go
client := fabcore.NewClientFactoryWithClient(fabricClient).NewConnectionsClient()
```

### Step 4B — Identify CRUD Methods

Non-item clients have bespoke method names. Read the client source to find:

- Create/Get/Update/Delete methods
- List pager methods
- Any sub-resource methods (e.g. role assignments)

### Step 5B — Identify DTO Structs

Read `fabric/core/models.go` or the specific client models for:

- Request/Response structs (e.g. `CreateConnectionRequest`, `Connection`)
- Nested configuration structs
- Enum types

---

## Output Format

### For Fabric Items (Category A)

| Field                   | Value                                                                    |
| ----------------------- | ------------------------------------------------------------------------ |
| **Category**            | Fabric Item                                                              |
| **SDK Version**         | `<version from go.mod>`                                                  |
| **SDK Package**         | `github.com/microsoft/fabric-sdk-go/fabric/<package>`                    |
| **Import Alias**        | `fab<package>`                                                           |
| **Client Factory**      | `fab<package>.NewClientFactoryWithClient(fabricClient).NewItemsClient()` |
| **Get Method**          | `client.Get<ItemName>(ctx, workspaceID, itemID, nil)`                    |
| **List Pager**          | `client.NewList<ItemNames>Pager(workspaceID, nil)`                       |
| **ItemType Constant**   | `fabcore.ItemType<ItemName>`                                             |
| **Archetype**           | `basic` / `definition` / `properties` / etc.                             |
| **Properties DTO**      | `fab<package>.Properties` — list fields                                  |
| **CreationPayload DTO** | `fab<package>.CreationPayload` — list fields (or "N/A")                  |
| **Enum Types**          | List any enum types with their possible values                           |

### For non-item resources (Category B)

| Field             | Value                                                       |
| ----------------- | ----------------------------------------------------------- |
| **Category**      | Non-item (bespoke CRUD)                                     |
| **SDK Version**   | `<version from go.mod>`                                     |
| **SDK Package**   | `github.com/microsoft/fabric-sdk-go/fabric/core` (or other) |
| **Import Alias**  | `fabcore` (or other)                                        |
| **Client**        | `fabcore.<Resource>Client`                                  |
| **CRUD Methods**  | List all relevant methods                                   |
| **Request DTOs**  | List create/update request structs with fields              |
| **Response DTOs** | List response structs with fields                           |
| **Enum Types**    | List any enum types with their possible values              |

### Go Struct Fields

For each DTO struct found, list all fields with their Go types. Example:

```go
type Properties struct {
    OneLakeFilesPath      *string                  `json:"oneLakeFilesPath,omitempty"`
    OneLakeTablesPath     *string                  `json:"oneLakeTablesPath,omitempty"`
    SQLEndpointProperties *SQLEndpointProperties   `json:"sqlEndpointProperties,omitempty"`
    DefaultSchema         *string                  `json:"defaultSchema,omitempty"`
}
```

The output from this skill feeds directly into `#skill:schema-model-generator` and `#skill:itemgen-command-builder`.
