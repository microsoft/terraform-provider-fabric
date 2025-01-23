// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// dataSourceOnPremisesGatewayPersonal fetches the simplified on-prem gateway: OnPremisesGatewayPersonal.
type dataSourceOnPremisesGatewayPersonal struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GatewaysClient
}

func NewDataSourceOnPremisesGatewayPersonal() datasource.DataSource {
	return &dataSourceOnPremisesGatewayPersonal{}
}

// Metadata sets the data source type name.
func (d *dataSourceOnPremisesGatewayPersonal) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_on_premises_gateway_personal"
}

// Schema defines the attributes for OnPremisesGatewayPersonal.
func (d *dataSourceOnPremisesGatewayPersonal) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve an on-premises gateway in its 'personal' form (ID, public key, type, version).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The gateway ID (UUID).",
				Optional:            true,
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the gateway.",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "The gateway version.",
				Computed:            true,
			},
			"public_key": schema.SingleNestedAttribute{
				MarkdownDescription: "The public key settings of the gateway.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[PublicKeyModel](ctx),
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
		},
	}
}

// Configure stores provider data and creates a new GatewaysClient.
func (d *dataSourceOnPremisesGatewayPersonal) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read fetches the OnPremisesGateway from the Fabric service, then maps to the personal model.
func (d *dataSourceOnPremisesGatewayPersonal) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OnPremisesGatewayPersonalModel

	// Parse config into data model
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	// If no ID was set, we can return early or handle lookups by name, etc. Here, handle if ID is blank
	if data.ID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing ID",
			"An ID is required to look up a personal on-premises gateway.",
		)
		return
	}

	// Do the actual GET call
	gatewayResp, errResp := d.client.GetGateway(ctx, data.ID.ValueString(), nil)
	if errResp != nil {
		resp.Diagnostics.AddError("GetGateway failed", errResp.Error())
		return
	}

	// Type-assert to OnPremisesGateway
	realGw, ok := gatewayResp.GatewayClassification.(*fabcore.OnPremisesGateway)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Gateway Type", "Result is not an OnPremisesGateway.")
		return
	}

	// Map the returned gateway to the personal model
	gateway := OnPremisesGatewayPersonalModel{}
	diags := gateway.set(ctx, *realGw)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save data back into state
	if diags := resp.State.Set(ctx, gateway); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}
