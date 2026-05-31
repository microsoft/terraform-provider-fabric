---
applyTo: "internal/services/**/schema_*.go,internal/services/**/models.go"
---

# Schema & Model Patterns

Ref: <https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas>

## Schema Approaches

### 1. Fabric Items — Separate Schema Functions

Define attributes in `schema_resource_<item>.go` / `schema_data_<item>.go`:

- `getResource<Item>PropertiesAttributes(ctx)`, `getResource<Item>ConfigurationAttributes()`
- `getDataSource<Item>PropertiesAttributes(ctx)`

Reference: `internal/services/lakehouse/schema_resource_lakehouse.go`

### 2. Non-item resources — Superschema (Combined)

Single `schema.go` using [`superschema`](https://github.com/orange-cloudavenue/terraform-plugin-framework-superschema) with `Common`, `Resource`, `DataSource` variants per attribute.

Import aliases: `schemaD` (datasource), `schemaR` (resource), `superschema`, `supertypes`.

| Scenario                         | Attribute Type                                    |
| -------------------------------- | ------------------------------------------------- |
| Same base attribute for both     | `superschema.StringAttribute`                     |
| Different `CustomType` or models | `superschema.SuperStringAttribute`                |
| Nested, same model               | `superschema.SingleNestedAttribute`               |
| Nested, different models         | `superschema.SuperSingleNestedAttribute`          |
| List nested, different models    | `superschema.SuperListNestedAttributeOf[<model>]` |
| Set nested, different models     | `superschema.SuperSetNestedAttributeOf[<model>]`  |

Schema consumption: `itemSchema(ctx, false).GetResource(ctx)` / `.GetDataSource(ctx)`.

Reference: `internal/services/connection/schema.go`

## Attribute Behaviors

| Behavior                               | Flags                    | Use Case                       |
| -------------------------------------- | ------------------------ | ------------------------------ |
| `Required: true`                       | Must be set by user      | `workspace_id`, `display_name` |
| `Optional: true`                       | May be set               | `description`                  |
| `Computed: true`                       | Server-computed          | `id` after create              |
| `Optional: true, Computed: true`       | Server fills default     | `capacity_id`                  |
| `Required: true` + `RequiresReplace()` | Immutable after creation | `enable_schemas`               |

**Always use `MarkdownDescription`** (never `Description`).

### Null vs Unknown — `Value<Type>Pointer()` Gotcha

- **Null** → `Value<Type>Pointer()` returns `nil`
- **Unknown** → `Value<Type>Pointer()` returns a pointer to the Go zero value (`""`, `false`, `0`)

`Computed: true` makes unset fields **unknown** during planning — not null. Avoid `Computed: true` on fields where the API or validators distinguish `nil` from the zero value. For example, `customtypes.UUID` fields managed via separate assign/unassign endpoints should use `Optional: true` only, because `""` fails UUID validation.

### Plan Modifiers

- `<type>planmodifier.RequiresReplace()` — destroy and recreate on change (string/bool/int32/float64 variants)
- `stringplanmodifier.UseStateForUnknown()` — preserve prior state during plan

### Validators

- `stringvalidator.OneOf(...)` — enum values (use `utils.ConvertEnumsToStringSlices()` for SDK enums)
- `stringvalidator.ConflictsWith(path)` / `AlsoRequires(path)` — mutual exclusion / co-requirement
- `stringvalidator.PreferWriteOnlyAttribute(path)` — deprecation path
- `float64validator.Between(min, max)` / `OneOf(...)` — numeric constraints
- `superstringvalidator.NullIfAttributeIsOneOf(path, values...)` — conditional null

## SDK Type → Schema Mapping

| SDK Type                | Schema Type                                                                                      |
| ----------------------- | ------------------------------------------------------------------------------------------------ |
| `*string`               | `schema.StringAttribute`                                                                         |
| `*bool`                 | `schema.BoolAttribute`                                                                           |
| `*int32` / `*int64`     | `schema.Int32Attribute` / `schema.Int64Attribute`                                                |
| UUID (`*string`)        | `schema.StringAttribute{CustomType: customtypes.UUIDType{}}`                                     |
| Nested struct           | `schema.SingleNestedAttribute{CustomType: supertypes.NewSingleNestedObjectTypeOf[<model>](ctx)}` |
| Slice of structs (list) | `schema.ListNestedAttribute{CustomType: supertypes.NewListNestedObjectTypeOf[<model>](ctx)}`     |
| Slice of structs (set)  | `schema.SetNestedAttribute{CustomType: supertypes.NewSetNestedObjectTypeOf[<model>](ctx)}`       |
| Enum types              | `schema.StringAttribute` (string representation)                                                 |

---

## Model Structs

**Naming:** `<item>PropertiesModel`, `<item>ConfigurationModel`, `<nested>Model`

Every field needs a `tfsdk:"snake_case"` tag.

### SDK Type → Model Type Mapping

| SDK Type                | Model Type                                            | Setter Pattern                                    |
| ----------------------- | ----------------------------------------------------- | ------------------------------------------------- |
| `*string`               | `types.String`                                        | `types.StringPointerValue(from.Field)`            |
| `*bool`                 | `types.Bool`                                          | `types.BoolPointerValue(from.Field)`              |
| `*int32` / `*int64`     | `types.Int64`                                         | `types.Int64Value(int64(*from.Field))`            |
| `*float32` / `*float64` | `types.Float64`                                       | `types.Float64Value(float64(*from.Field))`        |
| UUID (`*string`)        | `customtypes.UUID`                                    | `customtypes.NewUUIDPointerValue(from.ID)`        |
| Nested struct           | `supertypes.SingleNestedObjectValueOf[<nestedModel>]` | Create sub-model, call `.set()`, wrap             |
| Slice of structs (list) | `supertypes.ListNestedObjectValueOf[<nestedModel>]`   | Iterate, set each, collect                        |
| Slice of structs (set)  | `supertypes.SetNestedObjectValueOf[<nestedModel>]`    | Iterate, set each, collect                        |
| Enum (`*EnumType`)      | `types.String`                                        | `types.StringPointerValue((*string)(from.Field))` |

> **Note:** Prefer **Set** over List for collection attributes unless the order of elements is meaningful. Sets provide better drift detection since Terraform compares elements regardless of order.

### The `set()` Method

Every model must have `set()` mapping SDK DTO → TF model. Models with nested objects take `(ctx, from) diag.Diagnostics`; leaf models can omit `ctx` and return nothing. For nested sub-models: create `Null` value, check nil, create sub-model, call `sub.set()`, wrap with `.Set(ctx, sub)`.

Reference: `internal/services/lakehouse/models.go`
