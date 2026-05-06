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

## Generic Type Constructors

**Resources:** `fabricitem.NewResourceFabricItem[...]` — suffix matches archetype: `(config)`, `Definition(config)`, `Properties[T,S](config)`, `DefinitionProperties[T,S](config)`, `ConfigProperties[T,S,C,CS](config)`, `ConfigDefinitionProperties[T,S,C,CS](config)`

**Data Sources (singular):** `fabricitem.NewDataSourceFabricItem(config)`, `...Properties[T,S](config)`, `...DefinitionProperties[T,S](config)`

**Data Sources (plural):** `fabricitem.NewDataSourceFabricItems(config)`, `...Properties[T,S](config)`

## Closure Patterns

Resources with properties/config wire closures in the config struct. Which closures depend on the archetype:

```go
// propertiesSetter + itemGetter (all property archetypes)
PropertiesSetter: func(ctx context.Context, from *fab<pkg>.<Item>, model *<item>PropertiesModel) diag.Diagnostics {
    if from.Properties == nil { return nil }
    return model.set(ctx, *from.Properties)
},
ItemGetter: func(ctx context.Context, fabricClient fabric.Client, model fabricitem.FabricItemProperties[propertiesModel], fabricItem *fabricitem.FabricItemGetterData) error {
    client := fab<pkg>.NewClientFactoryWithClient(fabricClient).New<Items>Client()
    respGet, err := client.Get<Item>(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
    if err != nil { return err }
    fabricItem.Set(&respGet.<Item>)
    return nil
},

// creationPayloadSetter (config archetypes only)
CreationPayloadSetter: func(ctx context.Context, model <item>ConfigurationModel) (*fab<pkg>.CreationPayload, diag.Diagnostics) {
    return &fab<pkg>.CreationPayload{EnableSchemas: model.EnableSchemas.ValueBoolPointer()}, nil
},

// itemListGetter (plural data sources)
ItemListGetter: func(ctx context.Context, fabricClient fabric.Client, model fabricitem.FabricItemsProperties[propertiesModel], fabricItems *fabricitem.FabricItemsGetterData) error {
    client := fab<pkg>.NewClientFactoryWithClient(fabricClient).New<Items>Client()
    pager := client.NewList<Items>Pager(model.WorkspaceID.ValueString(), nil)
    for pager.More() {
        page, err := pager.NextPage(ctx)
        if err != nil { return err }
        fabricItems.Append(page.Value)
    }
    return nil
},
```

Reference: `internal/services/lakehouse/resource_lakehouse.go`

## Post-Itemgen Fix Guide

After scaffolding, fix `<TODO>` / `// TODO` placeholders:

|  #  | File(s)                  | Fix                                                                                                                                                                                                                                                    | Applies to          |
| :-: | ------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------- |
|  0  | `provider.go`            | Import service pkg (alphabetical); register resource in `Resources()` and data sources in `DataSources()`. If constructor accepts `ctx`, use wrapper: `func() resource.Resource { return pkg.NewResourceType(ctx) }`                                   | All                 |
|  1  | `base.go`                | Fill `FabricItemType`, `ItemFormatTypeDefault`, `ItemDefinitionEmpty`, `ItemDefinitionPathDocsURL`, `ItemTypeInfo` (`DocsURL`, `IsPreview`, `IsSPNSupported`), `itemDefinitionFormats` — values from issue. Ref: `internal/services/lakehouse/base.go` | All                 |
|  2  | `models.go`              | Replace stub fields/`set()` body (use `#skill:schema-model-generator`)                                                                                                                                                                                 | Properties+         |
|  3  | `schema_*.go`            | Replace `"TODO"` attributes (use `#skill:schema-model-generator`)                                                                                                                                                                                      | Properties+         |
|  4  | `resource_<type>.go`     | Set `DefinitionRequired`/`ConfigRequired` bools; wire `creationPayloadSetter`                                                                                                                                                                          | Definition+/Config+ |
|  5  | `data_*.go`              | Align `set()` calls with Fix 2 signature changes                                                                                                                                                                                                       | Properties+         |
|  6  | `fakes/fabric_<type>.go` | Populate `Properties` in `NewRandom*`                                                                                                                                                                                                                  | Properties+         |
|  7  | `*_test.go`              | Add `resource.TestCheckResourceAttrSet` for each property                                                                                                                                                                                              | Properties+         |

**`set()` signature:** Templates generate pointer `from *fab<pkg>.Properties`; canonical uses value `from fab<pkg>.Properties`. If changed to value, dereference at call sites.

**`ctx` parameter:** If schema uses `supertypes.NewSingleNestedObjectTypeOf[<model>](ctx)`, add `ctx` to constructor and wrap in `provider.go`.
