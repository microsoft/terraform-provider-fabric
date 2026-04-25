// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package eventstreamdestinationconnection

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/ephemeral/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
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

var _ ephemeral.EphemeralResourceWithConfigure = (*ephemeralEventstreamDestinationConnection)(nil)

type ephemeralEventstreamDestinationConnection struct {
	pConfigData *pconfig.ProviderData
	client      *fabeventstream.TopologyClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewEphemeralResourceEventstreamDestinationConnection() ephemeral.EphemeralResource {
	return &ephemeralEventstreamDestinationConnection{
		TypeInfo: ItemTypeInfo,
	}
}

func (e *ephemeralEventstreamDestinationConnection) Metadata(_ context.Context, _ ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = e.TypeInfo.FullTypeName(false)
}

func (e *ephemeralEventstreamDestinationConnection) Schema(ctx context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fabricitem.NewEphemeralResourceMarkdownDescription(e.TypeInfo, false),
		Attributes: map[string]schema.Attribute{
			"destination_id": schema.StringAttribute{
				MarkdownDescription: "The destination ID.",
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
			"consumer_group_name": schema.StringAttribute{
				MarkdownDescription: "The consumer group name.",
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

func (e *ephemeralEventstreamDestinationConnection) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	pConfigData, ok := req.ProviderData.(*pconfig.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			common.ErrorEphemeralResourceConfigType,
			fmt.Sprintf(common.ErrorFabricClientType, req.ProviderData),
		)

		return
	}

	e.pConfigData = pConfigData

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(e.TypeInfo.Name, e.TypeInfo.IsPreview, e.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}

	e.client = fabeventstream.NewClientFactoryWithClient(*pConfigData.FabricClient).NewTopologyClient()
}

func (e *ephemeralEventstreamDestinationConnection) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	tflog.Debug(ctx, "OPEN", map[string]any{
		"action": "start",
	})

	var data ephemeralEventstreamDestinationConnectionModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Open(ctx, e.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(e.get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)

	tflog.Debug(ctx, "OPEN", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (e *ephemeralEventstreamDestinationConnection) get(ctx context.Context, model *ephemeralEventstreamDestinationConnectionModel) diag.Diagnostics {
	respGet, err := e.client.GetEventstreamDestinationConnection(ctx, model.WorkspaceID.ValueString(), model.EventstreamID.ValueString(), model.DestinationID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationOpen, nil); diags.HasError() {
		return diags
	}

	return model.set(ctx, model.WorkspaceID.ValueString(), model.EventstreamID.ValueString(), model.DestinationID.ValueString(), respGet.DestinationConnectionResponse)
}
