// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceConnections)(nil)

type dataSourceConnections struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.ConnectionsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceConnections() datasource.DataSource {
	return &dataSourceConnections{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceConnections) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(true)
}

func (d *dataSourceConnections) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	s := itemSchema(true).GetDataSource(ctx)

	resp.Schema = schema.Schema{
		MarkdownDescription: s.GetMarkdownDescription(),
		Attributes: map[string]schema.Attribute{
			"values": schema.SetNestedAttribute{
				MarkdownDescription: "The set of " + d.TypeInfo.Names + ".",
				Computed:            true,
				CustomType:          supertypes.NewSetNestedObjectTypeOf[baseConnectionModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: s.Attributes,
				},
			},
		},
	}
}

func (d *dataSourceConnections) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dataSourceConnections) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceConnectionsModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	// Use a default timeout since the model doesn't have a Timeouts field
	timeout := d.pConfigData.Timeout

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(d.list(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})
}

func (d *dataSourceConnections) list(ctx context.Context, model *dataSourceConnectionsModel) diag.Diagnostics {
	respList, err := d.client.ListConnections(ctx, nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	return model.setValues(ctx, respList)
}
