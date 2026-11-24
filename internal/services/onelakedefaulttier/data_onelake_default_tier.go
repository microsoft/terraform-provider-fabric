// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedefaulttier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceOneLakeDefaultTier)(nil)

type dataSourceOneLakeDefaultTier struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.OneLakeSettingsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceOneLakeDefaultTier() datasource.DataSource {
	return &dataSourceOneLakeDefaultTier{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceOneLakeDefaultTier) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

func (d *dataSourceOneLakeDefaultTier) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = itemSchema().GetDataSource(ctx)
}

func (d *dataSourceOneLakeDefaultTier) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewOneLakeSettingsClient()
}

func (d *dataSourceOneLakeDefaultTier) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{"action": "start"})

	var data dataSourceOneLakeDefaultTierModel

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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Debug(ctx, "READ", map[string]any{"action": "end"})
}

func (d *dataSourceOneLakeDefaultTier) get(ctx context.Context, model *dataSourceOneLakeDefaultTierModel) diag.Diagnostics {
	workspaceID := model.WorkspaceID.ValueString()
	tflog.Trace(ctx, "getting OneLake default tier for Workspace ID: "+workspaceID)

	respGet, err := d.client.GetSettings(ctx, workspaceID, nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	var tier *fabcore.OneLakeAccessTier
	if respGet.Lifecycle != nil {
		tier = respGet.Lifecycle.DefaultTier
	}

	model.set(workspaceID, tier)

	return nil
}
