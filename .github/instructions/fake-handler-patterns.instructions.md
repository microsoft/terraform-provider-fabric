---
applyTo: "internal/testhelp/fakes/**"
---

# Fake Handler Patterns

Fake handlers provide mock SDK servers for unit tests. They live in `internal/testhelp/fakes/`.

## Operations Struct

Implements typed handler interfaces. Key methods:

| Method                                    | Purpose                                  | Required for                   |
| ----------------------------------------- | ---------------------------------------- | ------------------------------ |
| `CreateWithParentID(parentID, data)`      | Create entity under a parent             | `parentIDOperations`           |
| `Create(data)`                            | Create entity (no parent)                | `simpleIDOperations`           |
| `Filter(entities, parentID)`              | Filter by parent workspace ID            | `parentIDOperations`           |
| `GetID(entity)`                           | Generate composite ID                    | All                            |
| `TransformCreate/Get/List/Update(entity)` | Wrap in SDK response types               | All                            |
| `Update(base, data)`                      | Apply update request to existing entity  | All (except NoUpdate variants) |
| `Validate(newEntity, existing)`           | Check duplicates, return conflict errors | All                            |
| `CreateDefinition`/`UpdateDefinition`     | Definition CRUD                          | `definitionOperations`         |
| `TransformDefinition`                     | Wrap definition in response type         | `definitionOperations`         |

### Interface Hierarchy

- **`simpleIDOperations`** = `operationsBase` + `creator` — simple ID resources (workspace, gateway)
- **`parentIDOperations`** = `operationsBase` + `creatorWithParentID` + `parentFilter` — workspace-scoped (Fabric Items, folders)
- **`definitionOperations`** = `definitionCreator` + `definitionUpdater` + `definitionTransformer`

## Configure Function

### Handler Construction

- **`newTypedHandler(server, ops)`** — uses reflection-based `defaultConverter` for item-to-entity conversion. Use for entities without `Properties` or non-item resources.
- **`newTypedHandlerWithConverter(server, ops, converter)`** — uses a custom `itemConverter`. **Required** when your entity has `Properties` that must survive cross-type conversion (e.g. lakehouse). Implement `ConvertItemToEntity(fabcore.Item)` on your operations struct.

### Configure Variants

Choose the right wiring function based on your resource's API shape:

| Function                                   | Create Type | Has Update | Has Parent ID | Use When                                                |
| ------------------------------------------ | ----------- | ---------- | ------------- | ------------------------------------------------------- |
| `configureEntityPagerWithSimpleID`         | Sync        | ✓          | ✗             | Simple-ID resources (workspace, gateway)                |
| `configureEntityWithParentID`              | LRO         | ✓          | ✓             | Standard Fabric Items (lakehouse, notebook)             |
| `configureEntityWithParentIDNoLRO`         | Sync        | ✓          | ✓             | Parent-ID resources with sync create (folder)           |
| `configureEntityWithParentIDNoLRONoUpdate` | Sync        | ✗          | ✓             | Parent-ID resources with no update (create+delete only) |
| `configureDefinitions`                     | LRO         | —          | —             | Items with definitions (call **after** main configure)  |
| `configureDefinitionsNonLROCreation`       | Sync        | —          | —             | Items with non-LRO definition creation                  |

## Sub-Operations

Some resources have SDK calls beyond standard CRUD (e.g. `AssignToDomain`, `MoveFolder`).

**Key insight:** If the SDK fake server has no handler wired for a method and a test calls it, the fake returns `nonRetriableError("fake for method X not implemented")`, causing the test to fail.

### Decision Tree

```
Does your resource code call a sub-operation SDK method?
│
├── NO → Standard CRUD fakes handle it
│
└── YES → Does any unit test exercise that code path?
    │
    ├── NO → No fake needed (only acceptance tests hit it)
    │
    └── YES → You MUST wire a fake handler
            │
            ├── Updates entity state? → Use handler.Contains/Get/Upsert
            │   Example: AssignToDomain updates DomainID on workspace
            │
            ├── Fire-and-forget? → Return empty OK response
            │   Example: ProvisionIdentity returns 200 with no body
            │
            └── Returns static/reference data? → Return hardcoded response
                Example: ListSupportedConnectionTypes returns fixed metadata
```

### Override Pattern

For sub-operations that need custom list/filter logic or non-standard signatures, define **exported functions** that take the `typedHandler` and return the handler func, then assign them after the standard configure call:

```go
// In configure function, after configureEntityWithParentIDNoLRO(...)
server.ServerFactory.Core.FoldersServer.NewListFoldersPager = FakeListFolders(handler)
server.ServerFactory.Core.FoldersServer.MoveFolder = FakeMoveFolder(handler)
```

Reference: `fabric_folder.go` — `FakeMoveFolder`, `FakeListFolders`

## `typedHandler` API Reference

The `typedHandler[TEntity]` struct (from `fake_typedhandler.go`) embeds `fakeServer` and provides entity store access:

| Method           | Signature          | Purpose                                  |
| ---------------- | ------------------ | ---------------------------------------- |
| `Contains(id)`   | `(string) bool`    | Check if entity exists in store          |
| `Get(id)`        | `(string) TEntity` | Retrieve entity by ID                    |
| `Upsert(entity)` | `(TEntity)`        | Insert or update entity in store         |
| `Delete(id)`     | `(string)`         | Remove entity from store                 |
| `Elements()`     | `() []TEntity`     | Get all entities of this type            |
| `ServerFactory`  | field              | Access to `fabfake.ServerFactory` routes |

`fakeServer.Upsert(element any)` can also be used directly to pre-seed entities in tests (type must be registered via `handleEntity`).

## Polymorphic / Multi-Subtype Resources

Resources with `Classification` interfaces (gateway, connection) need **multiple type registrations** so `fakeServer.Upsert()` accepts all subtypes. Each subtype gets its own `configure*` wrapper that delegates to a shared configure function:

```go
func configureShareableCloudConnection(server *fakeServer) fabcore.ShareableCloudConnection {
    configureConnection(server) // shared wiring
    return fabcore.ShareableCloudConnection{}
}
func configureVirtualNetworkGatewayConnection(server *fakeServer) fabcore.VirtualNetworkGatewayConnection {
    configureConnection(server)
    return fabcore.VirtualNetworkGatewayConnection{}
}
```

Use type switches in `Create`/`Update` to handle subtype-specific fields. Reference: `fabric_gateway.go`, `fabric_connection.go`.

## Inline Fakes

**When to use:** The resource's API signature doesn't fit the `parentID/childID` pattern (e.g., shortcut uses `workspaceID, itemID, path, name` — 4 path params).

Inline fakes live in the service package's test directory (`fake_test.go`) and use:

- A **package-level map** as the store (e.g., `var fakeShortcutStore = map[string]fabcore.Shortcut{}`)
- **Direct handler functions** assigned to `ServerFactory` fields in test setup — no `typedHandler` or generic configure functions
- A **custom composite ID** function (e.g., `GenerateShortcutID(workspaceID, itemID, path, name)`)

Reference: `internal/services/shortcut/fake_test.go`

## Canonical References

- Fabric Item (with definitions): `internal/testhelp/fakes/fabric_lakehouse.go`
- Parent-ID resource (with sub-operation overrides): `internal/testhelp/fakes/fabric_folder.go`
- Non-item (simple ID): `internal/testhelp/fakes/fabric_workspace.go`
- Non-item (polymorphic): `internal/testhelp/fakes/fabric_gateway.go`
- Inline fakes: `internal/services/shortcut/fake_test.go`
