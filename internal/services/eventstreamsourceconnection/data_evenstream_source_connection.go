// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventstreamsourceconnection

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabeventstream "github.com/microsoft/fabric-sdk-go/fabric/eventstream"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceEventstreamSourceConnection)(nil)

type dataSourceEventstreamSourceConnection struct {
	pConfigData *pconfig.ProviderData
	client      *fabeventstream.TopologyClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceEventstreamSourceConnection() datasource.DataSource {
	return &dataSourceEventstreamSourceConnection{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceEventstreamSourceConnection) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

func (d *dataSourceEventstreamSourceConnection) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fabricitem.NewDataSourceMarkdownDescription(d.TypeInfo, false),
		Attributes: map[string]schema.Attribute{
			"source_id": schema.StringAttribute{
				MarkdownDescription: "The source ID.",
				CustomType:          customtypes.UUIDType{},
				Required:            true,
			},
			"eventstream_id": schema.StringAttribute{
				MarkdownDescription: "The eventstream ID.",
				CustomType:          customtypes.UUIDType{},
				Required:            true,
			},
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The workspace ID.",
				CustomType:          customtypes.UUIDType{},
				Required:            true,
			},
			"event_hub_name": schema.StringAttribute{
				MarkdownDescription: "The name of the event hub.",
				Computed:            true,
			},
			"fully_qualified_namespace": schema.StringAttribute{
				MarkdownDescription: "The fully qualified namespace of the event hub.",
				Computed:            true,
			},
			"access_keys": schema.SingleNestedAttribute{
				MarkdownDescription: "The access keys for the event hub.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[accessKeysModel](ctx),
				Attributes: map[string]schema.Attribute{
					"primary_key": schema.StringAttribute{
						MarkdownDescription: "The primary key for the event hub.",
						Computed:            true,
						Sensitive:           true,
					},
					"secondary_key": schema.StringAttribute{
						MarkdownDescription: "The secondary key for the event hub.",
						Computed:            true,
						Sensitive:           true,
					},
					"secondary_connection_string": schema.StringAttribute{
						MarkdownDescription: "The secondary connection string for the event hub.",
						Computed:            true,
						Sensitive:           true,
					},
					"primary_connection_string": schema.StringAttribute{
						MarkdownDescription: "The primary connection string for the event hub.",
						Computed:            true,
						Sensitive:           true,
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceEventstreamSourceConnection) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = fabeventstream.NewClientFactoryWithClient(*pConfigData.FabricClient).NewTopologyClient()
}

func (d *dataSourceEventstreamSourceConnection) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceEventstreamSourceConnectionModel

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

func (d *dataSourceEventstreamSourceConnection) get(ctx context.Context, model *dataSourceEventstreamSourceConnectionModel) diag.Diagnostics {
	respGet, err := d.client.GetEventstreamSourceConnection(ctx, model.WorkspaceID.ValueString(), model.EventstreamID.ValueString(), model.SourceID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	return model.set(ctx, model.WorkspaceID.ValueString(), model.EventstreamID.ValueString(), model.SourceID.ValueString(), respGet.SourceConnectionResponse)
}
