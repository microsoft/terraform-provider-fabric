---
name: schema-model-generator
description: Given a Go SDK contract (DTOs), generate the corresponding Terraform schema attributes and model structs with high fidelity. USE FOR: creating schema_*.go and models.go files from SDK DTOs, mapping Go types to Terraform attribute types and custom types.
---

# Skill: Schema & Model Generator

Given a Go SDK contract (DTOs), generate the corresponding Terraform schema attributes and model structs with high fidelity.

## Prerequisites

- SDK DTO struct fields have been identified (from `#skill:sdk-contract-navigator`)
- You know the item archetype (basic, definition, properties, etc.)
- If the issue contains a **🌳 DTO Nesting Depth Map**, use it to determine:
  - One model struct per tree node that introduces `[]Type` or `*Type`
  - Each such struct needs its own `set()` method
  - Which fields use `supertypes.ListNestedObjectValueOf` (nested objects) vs `supertypes.ListValueOf` (scalars) vs `supertypes.SingleNestedObjectValueOf` (optional nested)

## Step 1 — Classify Each SDK Field

For each field in the SDK DTO struct, determine its category using the **"Attribute Behaviors"** table in `.github/instructions/schema-model-patterns.instructions.md`.

Classification criteria:

| Category                  | Criterion                                          |
| ------------------------- | -------------------------------------------------- |
| **Read-only**             | Only appears in Get response `Properties` struct   |
| **Create-time-only**      | Appears in `CreationPayload` and cannot be updated |
| **Updatable required**    | User must provide, can be changed after create     |
| **Updatable optional**    | User may provide, can be changed after create      |
| **Optional with default** | Server sets a default if not provided              |

## Step 2 — Generate Model Struct Fields

Map each SDK field type to the corresponding Terraform model type using the **"SDK Type → Model Type Mapping"** table in `.github/instructions/schema-model-patterns.instructions.md`.

Every field must have a `tfsdk:"<snake_case_name>"` tag.

### Naming Conventions

- Model struct: `<item><purpose>Model` — e.g. `lakehousePropertiesModel`, `lakehouseConfigurationModel`
- Nested model: `<item><nested>Model` — e.g. `lakehouseSQLEndpointPropertiesModel`
- Field names: PascalCase in Go, with `tfsdk:"snake_case"` tag
- SDK PascalCase → Terraform snake_case: `OneLakeFilesPath` → `onelake_files_path`

### Example Properties Model

```go
type lakehousePropertiesModel struct {
    OneLakeFilesPath      types.String                                                              `tfsdk:"onelake_files_path"`
    OneLakeTablesPath     types.String                                                              `tfsdk:"onelake_tables_path"`
    SQLEndpointProperties supertypes.SingleNestedObjectValueOf[lakehouseSQLEndpointPropertiesModel] `tfsdk:"sql_endpoint_properties"`
    DefaultSchema         types.String                                                              `tfsdk:"default_schema"`
}
```

### Example Configuration Model (CreationPayload)

```go
type lakehouseConfigurationModel struct {
    EnableSchemas types.Bool `tfsdk:"enable_schemas"`
}
```

### Example Nested Model

```go
type lakehouseSQLEndpointPropertiesModel struct {
    ID                 customtypes.UUID `tfsdk:"id"`
    ConnectionString   types.String     `tfsdk:"connection_string"`
    ProvisioningStatus types.String     `tfsdk:"provisioning_status"`
}
```

## Step 3 — Generate Model Methods

Generate both directions of mapping: **response `set()`** (SDK → TF) and **request builders** (TF → SDK).

### 3a. Response `set()` — SDK → TF (both Fabric Items and Non-Items)

Every model struct needs a `set()` method that maps SDK response DTO → TF model.

**Top-level `set()` (with nested objects):**

Signature includes `context.Context` and returns `diag.Diagnostics`:

```go
func (to *<item>PropertiesModel) set(ctx context.Context, from fab<package>.<DTO>) diag.Diagnostics {
    to.SimpleField = types.StringPointerValue(from.SimpleField)
    // ... other simple fields

    // Handle nested struct
    nestedValue := supertypes.NewSingleNestedObjectValueOfNull[<nestedModel>](ctx)
    if from.NestedField != nil {
        nestedModel := &<nestedModel>{}
        nestedModel.set(*from.NestedField)  // or with ctx if nested has its own nested
        if diags := nestedValue.Set(ctx, nestedModel); diags.HasError() {
            return diags
        }
    }
    to.NestedField = nestedValue

    return nil
}
```

**Leaf `set()` (no nested objects):**

Simpler signature without `context.Context` or `diag.Diagnostics`:

```go
func (to *<nestedModel>) set(from fab<package>.<NestedDTO>) {
    to.ID = customtypes.NewUUIDPointerValue(from.ID)
    to.StringField = types.StringPointerValue(from.StringField)
    to.EnumField = types.StringPointerValue((*string)(from.EnumField))
}
```

**Setter patterns by type:** Use the "Setter Pattern" column in the "SDK Type → Model Type Mapping" table in `schema-model-patterns.instructions.md`.

### 3b. Request Builders — TF → SDK (Create/Update)

Both Fabric Items and Non-Items need TF→SDK mapping for writable fields. The pattern differs by category:

- **Fabric Items:** Inline in `creationPayloadSetter` closure (simple — typically 1-3 fields from configuration model). See `fabric-item-patterns.instructions.md` § "Closure Examples".
- **Non-Items:** Dedicated request builder structs with `set()` method that builds the SDK request directly (complex — full request DTOs)

**Non-Item Request Builder Struct** — embeds the SDK request type, `set()` populates it:

```go
type requestCreate<Type> struct {
    fabcore.Create<Type>Request // embedded SDK request type
}

func (to *requestCreate<Type>) set(ctx context.Context, from <type>ResourceModel) diag.Diagnostics {
    to.DisplayName = from.DisplayName.ValueStringPointer()
    to.Description = from.Description.ValueStringPointer()
    // ... map each writable field into the embedded struct
    return nil
}
```

Usage: `r.client.Create<Type>(ctx, reqCreate.Create<Type>Request, nil)`

**Inverse mapping rules:** Use the inverse of the "SDK Type → Model Type Mapping" table in `schema-model-patterns.instructions.md`. For each TF type, call its `Value*Pointer()` method (e.g., `types.String` → `.ValueStringPointer()`, `types.Bool` → `.ValueBoolPointer()`).

**Non-obvious cases:**

| TF Model Type                             | SDK Type      | Pattern                                                            |
| ----------------------------------------- | ------------- | ------------------------------------------------------------------ |
| `types.Int64`                             | `*int32`      | `ptr.To(int32(from.Field.ValueInt64()))` — type narrowing required |
| `supertypes.SingleNestedObjectValueOf[M]` | `*NestedDTO`  | `.Get(ctx)` → construct nested DTO from sub-model                  |
| `supertypes.ListNestedObjectValueOf[M]`   | `[]NestedDTO` | `.Get(ctx)` → iterate slice, build each DTO                        |

Reference: `internal/services/connection/models_resource_connection.go`

## Step 4 — Generate Schema Attributes

For Fabric Item resources, schema attributes go in separate functions:

### Resource Properties Schema

```go
// schema_resource_<item>.go
func getResource<Item>PropertiesAttributes(ctx context.Context) map[string]schema.Attribute {
    return map[string]schema.Attribute{
        "<field_name>": schema.StringAttribute{
            MarkdownDescription: "<description>.",
            Computed:            true,
        },
        // ... more attributes
    }
}
```

### Resource Configuration Schema (for CreationPayload)

```go
func getResource<Item>ConfigurationAttributes() map[string]schema.Attribute {
    return map[string]schema.Attribute{
        "<field_name>": schema.BoolAttribute{
            MarkdownDescription: "<description>.",
            Required:            true,
            PlanModifiers: []planmodifier.Bool{
                boolplanmodifier.RequiresReplace(),
            },
        },
    }
}
```

### Schema Attribute Type Mapping

Use the **"SDK Type → Schema Mapping"** table in `.github/instructions/schema-model-patterns.instructions.md`.

### Rules for All Schema Attributes

For attribute behavior flags, plan modifiers, and validators, refer to the **"Attribute Behaviors"**, **"Plan Modifiers"**, and **"Validators"** sections in `.github/instructions/schema-model-patterns.instructions.md`.

Additional rules:

1. **Always** use `MarkdownDescription` (never `Description`)
2. **Nested objects**: Must include `CustomType: supertypes.NewSingleNestedObjectTypeOf[<model>](ctx)`
3. **UUID fields**: Must include `CustomType: customtypes.UUIDType{}`

### Example Nested Attribute

```go
"sql_endpoint_properties": schema.SingleNestedAttribute{
    MarkdownDescription: "An object containing the properties of the SQL endpoint.",
    Computed:            true,
    CustomType:          supertypes.NewSingleNestedObjectTypeOf[lakehouseSQLEndpointPropertiesModel](ctx),
    Attributes: map[string]schema.Attribute{
        "provisioning_status": schema.StringAttribute{
            MarkdownDescription: "The SQL endpoint provisioning status.",
            Computed:            true,
        },
        "connection_string": schema.StringAttribute{
            MarkdownDescription: "SQL endpoint connection string.",
            Computed:            true,
        },
        "id": schema.StringAttribute{
            MarkdownDescription: "SQL endpoint ID.",
            Computed:            true,
            CustomType:          customtypes.UUIDType{},
        },
    },
},
```

## Canonical References

- Model struct patterns: `internal/services/lakehouse/models.go`
- Resource schema patterns: `internal/services/lakehouse/schema_resource_lakehouse.go`
- Data source schema patterns: `internal/services/lakehouse/schema_data_lakehouse.go`
- Non-item schema (superschema): `internal/services/connection/schema.go`
- Request builder patterns: `internal/services/connection/models_resource_connection.go`
