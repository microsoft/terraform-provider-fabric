// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package semanticmodelconnectionbinding

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceSemanticModelConnectionBindings)(nil)

type dataSourceSemanticModelConnectionBindings struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.ItemsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceSemanticModelConnectionBindings() datasource.DataSource {
	return &dataSourceSemanticModelConnectionBindings{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceSemanticModelConnectionBindings) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(true)
}

func (d *dataSourceSemanticModelConnectionBindings) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dataSourceSchema(ctx, true)
}

func (d *dataSourceSemanticModelConnectionBindings) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewItemsClient()
}

func (d *dataSourceSemanticModelConnectionBindings) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{"action": "start"})

	var data dataSourceSemanticModelConnectionBindingsModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, d.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	bindings, err := listItemConnections(ctx, d.client, data.WorkspaceID.ValueString(), data.SemanticModelID.ValueString())
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationList, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(data.setValues(ctx, bindings)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, "READ", map[string]any{"action": "end"})
}
