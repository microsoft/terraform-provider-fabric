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

// dataSourceOnPremisesGatewayPersonals is the "plural" version of data_on_premises_gateway_personal.go.
type dataSourceOnPremisesGatewayPersonals struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GatewaysClient
}

func NewDataSourceOnPremisesGatewayPersonals() datasource.DataSource {
	return &dataSourceOnPremisesGatewayPersonals{}
}

func (d *dataSourceOnPremisesGatewayPersonals) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// e.g.: fabric_on_premises_gateway_personals
	resp.TypeName = req.ProviderTypeName + "_on_premises_gateway_personals"
}

func (d *dataSourceOnPremisesGatewayPersonals) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all on-premises personal gateways.",
		Attributes: map[string]schema.Attribute{
			"values": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "A list of on-premises personal gateways.",
				CustomType:          supertypes.NewListNestedObjectTypeOf[onPremisesGatewayPersonalModelBase](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The on-premises personal gateway ID.",
							Computed:            true,
							CustomType:          customtypes.UUIDType{},
						},
						"public_key": schema.SingleNestedAttribute{
							MarkdownDescription: "The public key settings.",
							Computed:            true,
							CustomType:          supertypes.NewSingleNestedObjectTypeOf[publicKeyModel](ctx),
							Attributes: map[string]schema.Attribute{
								"exponent": schema.StringAttribute{
									MarkdownDescription: "RSA exponent.",
									Computed:            true,
								},
								"modulus": schema.StringAttribute{
									MarkdownDescription: "RSA modulus.",
									Computed:            true,
								},
							},
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The on-premises personal gateway type.",
							Computed:            true,
						},
						"version": schema.StringAttribute{
							MarkdownDescription: "The personal gateway version.",
							Computed:            true,
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceOnPremisesGatewayPersonals) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dataSourceOnPremisesGatewayPersonals) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ-ALL-On-Premises-Gateway-Personals", map[string]any{"action": "start"})

	var data dataSourceOnPremisesGatewayPersonalsModel
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
	tflog.Debug(ctx, "READ-ALL-On-Premises-Gateway-Personals", map[string]any{"action": "end"})
}

func (d *dataSourceOnPremisesGatewayPersonals) list(ctx context.Context, model *dataSourceOnPremisesGatewayPersonalsModel) diag.Diagnostics {
	tflog.Trace(ctx, "Listing all on-premises personal gateways")

	allItems, err := d.client.ListGateways(ctx, nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	return model.setValues(ctx, allItems)
}
