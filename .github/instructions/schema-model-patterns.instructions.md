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

---

## Design Decisions — When to Apply Patterns

### Polymorphic / Union Types

When the SDK has an interface or struct with a type discriminator and variant-specific fields:

- Model each variant as a **separate optional nested block** in the schema
- Add `ExactlyOneOf` or `AtLeastOneOf` resource-level validator across the variant blocks
- Each variant gets its own model struct with its own `set()` method
- The parent model holds all variant fields; non-applicable variants are set to `Null`

References: `internal/services/shortcut/` (target variants), `internal/services/connection/` (credential types)

### Immutability — When to Use `RequiresReplace`

Apply `RequiresReplace()` plan modifier when:

- The field exists only in the Create request DTO, not in the Update request DTO
- The field is used as a path/identity parameter (e.g., `workspace_id`, `item_id`, `path`, `name` for sub-resources)
- The API returns an error if you attempt to change the value after creation

Do **not** apply `RequiresReplace` to fields that are simply omitted from Update for convenience — only when the API fundamentally cannot change them.

### Enum Filtering

When exposing SDK enum values via `OneOf` validators:

- **Remove non-creatable variants** — enum values that only appear in responses or represent system-managed states (e.g., `OnPremises`/`OnPremisesPersonal` gateway types cannot be created via API)
- **Remove deprecated/internal values** — enum values marked for removal or internal-only use
- Use `utils.RemoveSlicesByValues()` to filter `Possible*Values()` slices

Reference: `internal/services/gateway/schema.go`

### Conditional Nulling by Type Discriminator

When a resource has type-specific fields that only apply to certain variants:

- Use `superstringvalidator.NullIfAttributeIsOneOf(path, values...)` to enforce null when type doesn't match
- Use `superstringvalidator.RequireIfAttributeIsOneOf(path, values...)` to enforce presence when type matches
- In `set()` methods, explicitly set non-applicable fields to `types.<Type>Null()` based on the type discriminator

Reference: `internal/services/gateway/models.go`

### Computed-Only Fields

Mark as `Computed: true` (without `Optional`) when:

- The field is only in the response DTO, never in create/update requests
- The field represents server-generated state (e.g., `provisioning_status`, `capacity_region`, endpoint URLs)
- Add `UseStateForUnknown()` if the value never changes after creation

### Provider-Only Attributes

Add attributes not in the SDK only when **operational control** is needed:

- Validation bypass flags (`skip_capacity_state_validation`) — controls provider-side logic (e.g., whether to call the capacity list API before proceeding). Always `Optional` + `BoolDefault(false)`. Never sent to the Fabric API.

Do **not** invent provider-only attributes speculatively. Only add them when a concrete operational need exists (e.g., the provider performs a pre-flight check that some callers cannot execute due to permissions).

### Resource vs Data Source Field Subsetting

- **Data sources** expose only read fields (response DTO) — no create/update-only fields, no write-only attributes
- **Resources** expose the full schema including request-only fields
- Use separate model types when the field set differs significantly (e.g., `dsConnectionDetailsModel` vs `rsConnectionDetailsModel`)
