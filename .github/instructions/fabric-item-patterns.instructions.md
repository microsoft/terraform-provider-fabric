---
applyTo: "internal/services/**/*.go"
---

# Fabric Item Patterns

## Item Archetypes

~80% of resources use the generic `fabricitem` abstraction. Choose archetype based on SDK capabilities:

| Archetype                      | Definition | Properties | Config | Example              |
| ------------------------------ | :--------: | :--------: | :----: | -------------------- |
| `basic`                        |     ✗      |     ✗      |   ✗    | KQL Dashboard        |
| `definition`                   |     ✓      |     ✗      |   ✗    | Data Pipeline        |
| `properties`                   |     ✗      |     ✓      |   ✗    | Environment          |
| `definition-properties`        |     ✓      |     ✓      |   ✗    | Spark Job Definition |
| `config-properties`            |     ✗      |     ✓      |   ✓    | Warehouse            |
| `config-definition-properties` |     ✓      |     ✓      |   ✓    | Lakehouse            |

## `base.go` Pattern

Constants: `FabricItemType`, `ItemFormatTypeDefault`, `ItemDefinitionEmpty`, `ItemDefinitionPathDocsURL`, `ItemTypeInfo`, `itemDefinitionFormats`. See `internal/services/lakehouse/base.go`.

## Generic Type Constructors

**Resources:** `fabricitem.NewResourceFabricItem[...]` — suffix matches archetype: `(config)`, `Definition(config)`, `Properties[T,S](config)`, `DefinitionProperties[T,S](config)`, `ConfigProperties[T,S,C,CS](config)`, `ConfigDefinitionProperties[T,S,C,CS](config)`

**Data Sources (singular):** `fabricitem.NewDataSourceFabricItem(config)`, `...Properties[T,S](config)`, `...DefinitionProperties[T,S](config)`

**Data Sources (plural):** `fabricitem.NewDataSourceFabricItems(config)`, `...Properties[T,S](config)`

## Closure Patterns

Resources with properties/config require closures:

- `creationPayloadSetter` — maps TF config model → SDK creation payload
- `propertiesSetter` — maps SDK properties → TF model
- `itemGetter` — fetches single item via SDK typed client
- `itemListGetter` (data sources) — fetches via pager, matches by DisplayName

Reference: `internal/services/lakehouse/resource_lakehouse.go`

## Code Generator

```bash
go run tools/itemgen/main.go -item-name "<Name>" -items-name "<Names>" -item-type <archetype> -rename-allowed=<bool> -is-preview=<bool> -is-spn-supported=<bool>
```

Additional flags: `-generate-fakes`, `-generate-examples`

## Post-Itemgen Fix Guide

After scaffolding, fix `<TODO>` / `// TODO` placeholders:

|  #  | File(s)                  | Fix                                                                                                                                                                                                                  | Applies to          |
| :-: | ------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------- |
|  0  | `provider.go`            | Import service pkg (alphabetical); register resource in `Resources()` and data sources in `DataSources()`. If constructor accepts `ctx`, use wrapper: `func() resource.Resource { return pkg.NewResourceType(ctx) }` | All                 |
|  1  | `base.go`                | Fill `DocsURL`, `IsPreview`, `IsSPNSupported`, definition URLs and paths (values from issue)                                                                                                                         | All                 |
|  2  | `models.go`              | Replace stub fields/`set()` body (use `#skill:schema-model-generator`)                                                                                                                                               | Properties+         |
|  3  | `schema_*.go`            | Replace `"TODO"` attributes (use `#skill:schema-model-generator`)                                                                                                                                                    | Properties+         |
|  4  | `resource_<type>.go`     | Set `DefinitionRequired`/`ConfigRequired` bools; wire `creationPayloadSetter`                                                                                                                                        | Definition+/Config+ |
|  5  | `data_*.go`              | Align `set()` calls with Fix 2 signature changes                                                                                                                                                                     | Properties+         |
|  6  | `fakes/fabric_<type>.go` | Populate `Properties` in `NewRandom*`                                                                                                                                                                                | Properties+         |
|  7  | `*_test.go`              | Add `resource.TestCheckResourceAttrSet` for each property                                                                                                                                                            | Properties+         |

**`set()` signature:** Templates generate pointer `from *fab<pkg>.Properties`; canonical uses value `from fab<pkg>.Properties`. If changed to value, dereference at call sites.

**`ctx` parameter:** If schema uses `supertypes.NewSingleNestedObjectTypeOf[<model>](ctx)`, add `ctx` to constructor and wrap in `provider.go`.
