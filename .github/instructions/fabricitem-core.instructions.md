---
applyTo: "internal/pkg/fabricitem/**"
---

# Fabric Item Core Abstraction

`internal/pkg/fabricitem/` is the generic abstraction layer used by ~60% of provider resources. Service packages compose with it via configuration structs and closures — they never implement CRUD directly.

> **Do not modify this package for individual item implementations.** Service-specific logic belongs in `internal/services/<package>/`. Only modify this package to add capabilities that apply to ALL Fabric Items.

## Type Hierarchy — Embedding Matters

Resource structs compose via embedding. **Note:** definition archetypes embed `ResourceFabricItemDefinition` (which has its own `pConfigData`, `client`, and definition fields), NOT `ResourceFabricItem`.

| Struct                                         | Embeds                         | Type Params                                    | Closures Required                                         |
| ---------------------------------------------- | ------------------------------ | ---------------------------------------------- | --------------------------------------------------------- |
| `ResourceFabricItem`                           | —                              | —                                              | None                                                      |
| `ResourceFabricItemDefinition`                 | — (standalone, NOT Item)       | —                                              | None                                                      |
| `ResourceFabricItemProperties`                 | `ResourceFabricItem`           | `[Ttfprop, Titemprop]`                         | `PropertiesSetter`, `ItemGetter`                          |
| `ResourceFabricItemDefinitionProperties`       | `ResourceFabricItemDefinition` | `[Ttfprop, Titemprop]`                         | `PropertiesSetter`, `ItemGetter`                          |
| `ResourceFabricItemConfigProperties`           | `ResourceFabricItem`           | `[Ttfprop, Titemprop, Ttfconfig, Titemconfig]` | `PropertiesSetter`, `ItemGetter`, `CreationPayloadSetter` |
| `ResourceFabricItemConfigDefinitionProperties` | `ResourceFabricItemDefinition` | `[Ttfprop, Titemprop, Ttfconfig, Titemconfig]` | `PropertiesSetter`, `ItemGetter`, `CreationPayloadSetter` |

Data source types: `DataSourceFabricItem`, `DataSourceFabricItemProperties[T,S]`, `DataSourceFabricItemDefinitionProperties[T,S]`, `DataSourceFabricItems`, `DataSourceFabricItemsProperties[T,S]`.

## Closure Signatures — The Extension Points

Closures are how service packages inject item-specific logic. **The `to` parameter type changes per archetype** — this is a common source of type mismatch errors.

### `PropertiesSetter`

```go
// Signature varies — the model type matches the resource archetype:
func(ctx context.Context, from *Titemprop, to *ResourceFabricItemPropertiesModel[Ttfprop, Titemprop]) diag.Diagnostics
func(ctx context.Context, from *Titemprop, to *ResourceFabricItemConfigPropertiesModel[Ttfprop, Titemprop, Ttfconfig, Titemconfig]) diag.Diagnostics
func(ctx context.Context, from *Titemprop, to *ResourceFabricItemConfigDefinitionPropertiesModel[Ttfprop, Titemprop, Ttfconfig, Titemconfig]) diag.Diagnostics
```

### `ItemGetter`

Fetches a single item using the typed SDK client (not the generic `ItemsClient`):

```go
func(ctx context.Context, fabricClient fabric.Client, model <ModelType>, fabricItem *FabricItemProperties[Titemprop]) error
```

### `CreationPayloadSetter` (config archetypes only)

```go
func(ctx context.Context, from Ttfconfig) (*Titemconfig, diag.Diagnostics)
```

### `ItemListGetter` (plural data sources only)

```go
func(ctx context.Context, fabricClient fabric.Client, model <ModelType>, fabricItems []FabricItemProperties[Titemprop]) error
```

## `FabricItemProperties[T]` — Reflection-Based Mapper

`FabricItemProperties[Titemprop]` uses **reflection** (not interfaces) via `Set(from any)` to extract common fields (`ID`, `DisplayName`, `Description`, `WorkspaceID`, `FolderID`) and `Properties` from any typed SDK response struct.

This avoids requiring all SDK item types to implement a shared interface. The `getFieldStringValue` and `getFieldStructValue[T]` helpers do the reflection work.

**Key implication:** If a new SDK item struct names fields differently than the standard pattern, `Set()` will silently return `nil` for those fields. Always verify field names match when onboarding a new item type.

## Retry Logic

`CreateItem` and `UpdateItem` wrap SDK calls with `RetryOperationWithResult[T]`, which retries indefinitely on `ItemDisplayNameNotAvailableYet` errors (a transient Fabric API issue when display names haven't propagated). Default retry interval: 2 minutes.

## Schema Generation

Schema functions (`resource_schema.go`, `data_schema.go`) compose base attributes (`workspace_id`, `id`, `display_name`, `description`, `folder_id`, `timeouts`) with archetype-specific additions (`properties` nested object, `configuration` nested object with `RequiresReplace`, `definition` map + `format` + `definition_update_enabled`). Service packages pass their custom attributes via `PropertiesAttributes` and `ConfigAttributes` maps.

## When to Modify This Package

| ✅ DO modify for                      | ❌ DO NOT modify for                                  |
| ------------------------------------- | ----------------------------------------------------- |
| New base attribute ALL items share    | Item-specific properties (use closures)               |
| New archetype variant                 | Item-specific schema attributes (pass via maps)       |
| Bug in shared CRUD/retry/schema logic | Custom SDK response handling (implement `ItemGetter`) |
| Markdown description template changes | Item-specific validation (add in service schema)      |
| New shared plan modifier or validator | Anything only one item needs                          |
