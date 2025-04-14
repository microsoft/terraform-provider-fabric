// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure datasource implements expected interfaces.
var _ datasource.DataSourceWithConfigure = (*dataSourceConnection)(nil)

type dataSourceConnection struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.ConnectionsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceConnection() datasource.DataSource {
	return &dataSourceConnection{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceConnection) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

func (d *dataSourceConnection) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = itemSchema(false).GetDataSource(ctx)
}

func (d *dataSourceConnection) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	// Check if this is a preview resource
	if d.TypeInfo.IsPreview && !d.pConfigData.Preview {
		resp.Diagnostics.AddError(
			"Preview Resource Type",
			fmt.Sprintf("The %s data source is currently in preview and must be opted in using the provider's preview setting", d.TypeInfo.Name),
		)

		return
	}

	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewConnectionsClient()
}

func (d *dataSourceConnection) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourceConnectionModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, d.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if data.ID.ValueString() != "" {
		diags = d.getByID(ctx, &data)
	} else {
		diags = d.getByDisplayName(ctx, &data)
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dataSourceConnection) getByID(ctx context.Context, model *dataSourceConnectionModel) diag.Diagnostics {
	respGet, err := d.client.GetConnection(ctx, model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	if diags := model.set(ctx, respGet.Connection); diags.HasError() {
		return diags
	}

	return nil
}

func (d *dataSourceConnection) getByDisplayName(ctx context.Context, model *dataSourceConnectionModel) diag.Diagnostics {
	var diags diag.Diagnostics

	pager := d.client.NewListConnectionsPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Value {
			var entityDisplayName string

			if entityDisplayName == model.DisplayName.ValueString() {
				model.ID = customtypes.NewUUIDPointerValue(entity.ID)

				return d.getByID(ctx, model)
			}
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		"Unable to find Connection with 'display_name': "+model.DisplayName.ValueString(),
	)

	return diags
}
