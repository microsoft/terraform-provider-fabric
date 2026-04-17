// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package warehousesqlauditsetting

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabwarehouse "github.com/microsoft/fabric-sdk-go/fabric/warehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceWarehouseSQLAuditSettings)(nil)

type dataSourceWarehouseSQLAuditSettings struct {
	pConfigData *pconfig.ProviderData
	client      *fabwarehouse.SQLAuditSettingsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceWarehouseSQLAuditSettings() datasource.DataSource {
	return &dataSourceWarehouseSQLAuditSettings{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceWarehouseSQLAuditSettings) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

func (d *dataSourceWarehouseSQLAuditSettings) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = itemSchema().GetDataSource(ctx)
}

func (d *dataSourceWarehouseSQLAuditSettings) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	pConfigData, ok := req.ProviderData.(*pconfig.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			common.ErrorDataSourceConfigType,
			fmt.Sprintf(common.ErrorFabricClientType, req.ProviderData),
		)

		return
	}

	d.pConfigData = pConfigData

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(d.TypeInfo.Name, d.TypeInfo.IsPreview, d.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}

	d.client = fabwarehouse.NewClientFactoryWithClient(*pConfigData.FabricClient).NewSQLAuditSettingsClient()
}

func (d *dataSourceWarehouseSQLAuditSettings) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceWarehouseSQLAuditSettingsModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, d.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(d.get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})
}

func (d *dataSourceWarehouseSQLAuditSettings) get(ctx context.Context, model *dataSourceWarehouseSQLAuditSettingsModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s for Warehouse ID: %s in Workspace ID: %s", d.TypeInfo.Name, model.ItemID.ValueString(), model.WorkspaceID.ValueString()))

	respGet, err := d.client.GetSQLAuditSettings(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	return model.set(ctx, respGet.SQLAuditSettings)
}
