// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure it implements the interface.
var _ datasource.DataSourceWithConfigure = (*dataSourceOnPremisesGateways)(nil)

// dataSourceOnPremisesGateways is analogous to data_virtual_network_gateways.go, but for on-premises gateways (plural).
type dataSourceOnPremisesGateways struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GatewaysClient
}

func NewDataSourceOnPremisesGateways() datasource.DataSource {
	return &dataSourceOnPremisesGateways{}
}

func (d *dataSourceOnPremisesGateways) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + OnPremisesItemsTFType
}

func (d *dataSourceOnPremisesGateways) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all Fabric on-premises gateways.",
		Attributes: map[string]schema.Attribute{
			"values": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "A list of on-premises gateways.",
				CustomType:          supertypes.NewListNestedObjectTypeOf[onPremisesGatewayModelBase](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The on-premises gateway ID.",
							Computed:            true,
							CustomType:          customtypes.UUIDType{},
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The display name of the on-premises gateway.",
							Computed:            true,
						},
						"allow_custom_connectors": schema.BoolAttribute{
							MarkdownDescription: "Allow custom connectors.",
							Computed:            true,
						},
						"allow_cloud_connection_refresh": schema.BoolAttribute{
							MarkdownDescription: "Allow custom connectors refresh.",
							Computed:            true,
						},
						"number_of_member_gateways": schema.Int64Attribute{
							MarkdownDescription: "The number of member gateways.",
							Computed:            true,
						},
						"load_balancing_setting": schema.StringAttribute{
							MarkdownDescription: "The load balancing setting.",
							Computed:            true,
						},
						"public_key": schema.SingleNestedAttribute{
							MarkdownDescription: "The public key settings.",
							Computed:            true,
							CustomType:          supertypes.NewSingleNestedObjectTypeOf[publicKeyModel](ctx),
							Attributes: map[string]schema.Attribute{
								"exponent": schema.StringAttribute{
									MarkdownDescription: "The RSA exponent.",
									Computed:            true,
								},
								"modulus": schema.StringAttribute{
									MarkdownDescription: "The RSA modulus.",
									Computed:            true,
								},
							},
						},
						"version": schema.StringAttribute{
							MarkdownDescription: "The gateway version.",
							Computed:            true,
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceOnPremisesGateways) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = (*fabcore.GatewaysClient)(fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewGatewaysClient())
}

func (d *dataSourceOnPremisesGateways) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ-ALL-On-Premises-Gateways", map[string]any{"action": "start"})

	var data dataSourceOnPremisesGatewaysModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, d.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(d.list(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
	tflog.Debug(ctx, "READ-ALL-On-Premises-Gateways", map[string]any{"action": "end"})
}

// list retrieves all gateways from the Fabric SDK and filters only the on-premises gateway ones.
func (d *dataSourceOnPremisesGateways) list(ctx context.Context, model *dataSourceOnPremisesGatewaysModel) diag.Diagnostics {
	tflog.Trace(ctx, "Listing all on-premises gateways")

	gatewaysResp, err := d.client.ListGateways(ctx, nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	return model.setValues(ctx, gatewaysResp)
}
