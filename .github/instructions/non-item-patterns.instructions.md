---
applyTo: "internal/services/**/*.go"
---

# Non-Item Resource Patterns

Non-Item resources (~20%) do **not** use the `fabricitem` generic abstraction. They include Connection, Shortcut, Gateway, Workspace, role assignments, and similar bespoke CRUD resources.

> For Fabric Item patterns (Category A), see `fabric-item-patterns.instructions.md`.

## File Structure

**Key difference from Fabric Items:** Non-Item resources use a single `schema.go` with `superschema` instead of separate `schema_resource_*.go` / `schema_data_*.go` files. They may split models across `models_resource_<type>.go` and `models_data_<type>.go`.

| File                                 | Purpose                                                                          |
| ------------------------------------ | -------------------------------------------------------------------------------- |
| `base.go`                            | `ItemTypeInfo` only (no `FabricItemType` or `itemDefinitionFormats`)             |
| `schema.go`                          | Shared `superschema` for both resource and data source                           |
| `resource_<type>.go`                 | Resource struct with `Create`, `Read`, `Update`, `Delete`                        |
| `data_<type>.go` / `data_<types>.go` | Singular / plural data sources                                                   |
| `models.go`                          | Shared models; optionally split into `models_resource_*.go` / `models_data_*.go` |

## `base.go` Pattern

Only defines `ItemTypeInfo` (including `IsPreview` and `IsSPNSupported` — values sourced from the GitHub issue) — no `FabricItemType`, `ItemFormatTypeDefault`, or `itemDefinitionFormats`. Optionally add resource-specific constants (e.g. `PossibleInactivityMinutesBeforeSleepValues` in gateway).

Reference: `internal/services/connection/base.go`, `internal/services/gateway/base.go`

## Superschema

Non-Item resources use `superschema` to define a single schema function `itemSchema(ctx, isList)` in `schema.go`. For superschema details (imports, attribute types, consumption pattern), see `schema-model-patterns.instructions.md` § "Non-Item Resources — Superschema".

Reference: `internal/services/connection/schema.go`, `internal/services/gateway/schema.go`

## Resource Implementation

Non-Item resources implement `resource.Resource` directly — **no closures** like Fabric Items.

**Struct:** Fields are `pConfigData *pconfig.ProviderData`, `client *fabcore.<Resource>Client`, `TypeInfo tftypeinfo.TFTypeInfo`.

**Configure:** Extract `pConfigData` from `req.ProviderData`, create SDK client via `fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).New<Resource>Client()`, call `fabricitem.IsPreviewMode(...)`.

**CRUD pattern** — each method follows these steps:

1. Read plan/state into model struct
2. Get timeout from `Timeouts`
3. Build SDK request DTO via request builder's `set()` method
4. Call SDK client (e.g. `r.client.Create<Resource>(ctx, ...)`)
5. Map response → model via `model.set(ctx, response)`
6. Save model to state

**Read not-found:** Use `utils.IsErrNotFound(...)` → `resp.State.RemoveResource(ctx)`.

Reference: `internal/services/connection/resource_connection.go`, `internal/services/gateway/resource_gateway.go`

## Data Source Implementation

Same struct pattern as resource but implements `datasource.DataSource`. Constructor: `NewDataSource<Type>()`.

Reference: `internal/services/connection/data_connection.go`, `internal/services/gateway/data_gateway.go`

## Model Pattern

Non-Item models may use **generic type parameters** to share a base model between resource and data source with different nested types (e.g. `baseConnectionModel[ConnectionDetails, CredentialDetails]`). The `set()` method uses type switches to handle both variants.

Reference: `internal/services/connection/models.go`

## SDK Client Mapping

| Resource Type    | SDK Client                       |
| ---------------- | -------------------------------- |
| Connection       | `fabcore.ConnectionsClient`      |
| Shortcut         | `fabcore.ShortcutsClient`        |
| Gateway          | `fabcore.GatewaysClient`         |
| Workspace        | `fabcore.WorkspacesClient`       |
| Domain           | `fabcore.DomainsClient`          |
| Role Assignments | Various `*Client` from `fabcore` |

## Canonical References

| Resource Type    | Reference                         |
| ---------------- | --------------------------------- |
| Connection       | `internal/services/connection/`   |
| Shortcut         | `internal/services/shortcut/`     |
| Gateway          | `internal/services/gateway/`      |
| Workspace        | `internal/services/workspace/`    |
| Role Assignments | `internal/services/*ra/` packages |
