---
applyTo: "internal/services/**/*.go"
---

# Non-Item Resource Patterns

Non-Item resources (~20%) do **not** use the `fabricitem` generic abstraction. They include Connection, Shortcut, Gateway, Workspace, role assignments, and similar bespoke CRUD resources.

> For Fabric Item patterns (Category A), see `fabric-item-patterns.instructions.md`.

## File Structure

**Key difference from Fabric Items:** Non-Item resources use a single `schema.go` with `superschema` instead of separate `schema_resource_*.go` / `schema_data_*.go` files. They may split models across `models_resource_<type>.go` and `models_data_<type>.go`.

| File                                 | Purpose                                                                                                                                                      |
| ------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `base.go`                            | `ItemTypeInfo` only (`IsPreview`, `IsSPNSupported` from issue). No `FabricItemType`/`itemDefinitionFormats`. Optionally add resource-specific constants.     |
| `schema.go`                          | Shared `superschema` — single `itemSchema(ctx, isList)` function for both resource and data source. See `schema-model-patterns.instructions.md` for details. |
| `resource_<type>.go`                 | Resource struct with `Create`, `Read`, `Update`, `Delete`                                                                                                    |
| `data_<type>.go` / `data_<types>.go` | Singular / plural data sources                                                                                                                               |
| `models.go`                          | Shared models; optionally split into `models_resource_*.go` / `models_data_*.go`                                                                             |

## Resource Implementation

Non-Item resources implement `resource.Resource` directly — **no closures** like Fabric Items.

```go
type resource<Type> struct {
    pConfigData *pconfig.ProviderData
    client      *fabcore.<Type>Client
    TypeInfo    tftypeinfo.TFTypeInfo
}

func (r *resource<Type>) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    pConfigData := pconfig.Configure(req, resp)
    if resp.Diagnostics.HasError() { return }
    r.pConfigData = pConfigData
    r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).New<Type>Client()
    fabricitem.IsPreviewMode(r.TypeInfo, r.pConfigData, &resp.Diagnostics)
}
```

**CRUD pattern:**

```go
func (r *resource<Type>) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan <type>ResourceModel
    if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() { return }

    timeout, diags := plan.Timeouts.Create(ctx, r.TypeInfo.Timeouts.Create)
    if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() { return }
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()

    // Build request (request builder embeds SDK request type)
    var reqCreate requestCreate<Type>
    if resp.Diagnostics.Append(reqCreate.set(ctx, plan)...); resp.Diagnostics.HasError() { return }

    // Call SDK — pass the embedded SDK request field
    respCreate, err := r.client.Create<Type>(ctx, reqCreate.Create<Type>Request, nil)
    if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() { return }

    // Map response → model
    if resp.Diagnostics.Append(plan.set(ctx, respCreate.<Type>)...); resp.Diagnostics.HasError() { return }
    resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
```

**Read not-found:** Use `utils.IsErrNotFound(...)` → `resp.State.RemoveResource(ctx)`.

Reference: `internal/services/connection/resource_connection.go`, `internal/services/gateway/resource_gateway.go`

## Data Source Implementation

Same struct pattern as resource but implements `datasource.DataSource`. Constructor: `NewDataSource<Type>()`.

**Plural data source paging:**

```go
func (d *dataSource<Types>) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var state <types>Model
    if resp.Diagnostics.Append(req.Config.Get(ctx, &state)...); resp.Diagnostics.HasError() { return }

    pager := d.client.NewList<Types>Pager(nil)
    for pager.More() {
        page, err := pager.NextPage(ctx)
        if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationList, nil)...); resp.Diagnostics.HasError() { return }
        // Append page.Value items to state.Values
    }
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
```

Reference: `internal/services/connection/data_connection.go`, `internal/services/gateway/data_gateway.go`

## Model Pattern

Non-Item models may use **generic type parameters** to share a base model between resource and data source with different nested types (e.g. `baseConnectionModel[ConnectionDetails, CredentialDetails]`). The `set()` method uses type switches to handle both variants.

Reference: `internal/services/connection/models.go`

## Canonical References

| Resource Type    | SDK Client                       | Reference                         | Key Pattern                       |
| ---------------- | -------------------------------- | --------------------------------- | --------------------------------- |
| Connection       | `fabcore.ConnectionsClient`      | `internal/services/connection/`   | Generic type params, polymorphic  |
| Shortcut         | `fabcore.ShortcutsClient`        | `internal/services/shortcut/`     | Inline fakes, non-standard paths  |
| Gateway          | `fabcore.GatewaysClient`         | `internal/services/gateway/`      | Polymorphic, `simpleIDOperations` |
| Workspace        | `fabcore.WorkspacesClient`       | `internal/services/workspace/`    | Simple CRUD                       |
| Domain           | `fabcore.DomainsClient`          | `internal/services/domain/`       | Non-workspace-scoped              |
| Role Assignments | Various `*Client` from `fabcore` | `internal/services/*ra/` packages | Sub-resource, composite ID        |
