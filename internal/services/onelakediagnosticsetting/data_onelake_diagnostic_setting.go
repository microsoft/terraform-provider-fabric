// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakediagnosticsetting

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceOneLakeDiagnosticSetting)(nil)

type dataSourceOneLakeDiagnosticSetting struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.OneLakeSettingsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceOneLakeDiagnosticSetting() datasource.DataSource {
	return &dataSourceOneLakeDiagnosticSetting{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceOneLakeDiagnosticSetting) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

func (d *dataSourceOneLakeDiagnosticSetting) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = itemSchema().GetDataSource(ctx)
}

func (d *dataSourceOneLakeDiagnosticSetting) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dataSourceOneLakeDiagnosticSetting) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourceOneLakeDiagnosticSettingModel
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
}

func (d *dataSourceOneLakeDiagnosticSetting) get(ctx context.Context, model *dataSourceOneLakeDiagnosticSettingModel) diag.Diagnostics {
	respGet, err := d.client.GetSettings(ctx, model.WorkspaceID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	if respGet.Diagnostics == nil {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic(common.ErrorReadHeader, "The OneLake settings response omitted diagnostics."),
		}
	}

	return model.set(ctx, model.WorkspaceID.ValueString(), *respGet.Diagnostics)
}
