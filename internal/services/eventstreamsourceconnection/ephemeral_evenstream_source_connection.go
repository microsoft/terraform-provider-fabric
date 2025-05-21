// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventstreamsourceconnection

import (
	"context"
	"fmt"

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

// _ ephemeral.EphemeralResourceWithConfigValidators = (*ephemeralEventstreamSourceConnection)(nil)
var (
	_ ephemeral.EphemeralResourceWithConfigure = (*ephemeralEventstreamSourceConnection)(nil)
)

type ephemeralEventstreamSourceConnection struct {
	pConfigData *pconfig.ProviderData
	client      *fabeventstream.TopologyClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewEphemeralResourceEventstreamSourceConnection() ephemeral.EphemeralResource {
	return &ephemeralEventstreamSourceConnection{
		TypeInfo: ItemTypeInfo,
	}
}

func (e *ephemeralEventstreamSourceConnection) Metadata(_ context.Context, _ ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = e.TypeInfo.FullTypeName(false)
}

func (e *ephemeralEventstreamSourceConnection) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	tflog.Debug(ctx, "OPEN", map[string]any{
		"action": "start",
	})

	var data ephemeralEventstreamSourceConnectionModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(e.get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)

	tflog.Debug(ctx, "OPEN", map[string]any{
		"action": "end",
	})
}

func (e *ephemeralEventstreamSourceConnection) Schema(ctx context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fabricitem.NewEphemeralResourceMarkdownDescription(e.TypeInfo, false),
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
					},
					"secondary_key": schema.StringAttribute{
						MarkdownDescription: "The secondary key for the event hub.",
						Computed:            true,
					},
					"secondary_connection_string": schema.StringAttribute{
						MarkdownDescription: "The secondary connection string for the event hub.",
						Computed:            true,
					},
					"primary_connection_string": schema.StringAttribute{
						MarkdownDescription: "The primary connection string for the event hub.",
						Computed:            true,
					},
				},
			},
		},
	}
}

// func (e *ephemeralEventstreamSourceConnection) ConfigValidators(_ context.Context) []ephemeral.ConfigValidator {
// 	return []ephemeral.ConfigValidator{
// 		ephemeralvalidator.Conflicting(
// 			path.MatchRoot("id"),
// 			path.MatchRoot("display_name"),
// 		),
// 		ephemeralvalidator.ExactlyOneOf(
// 			path.MatchRoot("id"),
// 			path.MatchRoot("display_name"),
// 		),
// 	}
// }

func (e *ephemeralEventstreamSourceConnection) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (e *ephemeralEventstreamSourceConnection) get(ctx context.Context, model *ephemeralEventstreamSourceConnectionModel) diag.Diagnostics {
	respGet, err := e.client.GetEventstreamSourceConnection(ctx, model.WorkspaceID.ValueString(), model.EventstreamID.ValueString(), model.SourceID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationOpen, nil); diags.HasError() {
		return diags
	}

	return model.set(ctx, model.WorkspaceID.ValueString(), model.EventstreamID.ValueString(), model.SourceID.ValueString(), respGet.SourceConnectionResponse)
}
