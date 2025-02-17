// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

type dataSourceOnPremisesGatewayPersonal struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GatewaysClient
}

func NewDataSourceOnPremisesGatewayPersonal() datasource.DataSource {
	return &dataSourceOnPremisesGatewayPersonal{}
}

func (d *dataSourceOnPremisesGatewayPersonal) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + OnPremisesPersonalItemType
}

func (d *dataSourceOnPremisesGatewayPersonal) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve an on-premises gateway in its 'personal' form (ID, public key, type, version).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The gateway ID.",
				CustomType:          customtypes.UUIDType{},
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "The gateway version.",
				Computed:            true,
			},
			"public_key": schema.SingleNestedAttribute{
				MarkdownDescription: "The public key settings of the gateway.",
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
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceOnPremisesGatewayPersonal) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewGatewaysClient()
}

func (d *dataSourceOnPremisesGatewayPersonal) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data datasourceOnPremisesGatewayPersonalModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, d.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	gatewayResp, err := d.client.GetGateway(ctx, data.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError("GetGateway failed", err.Error())
		return
	}

	realGw, ok := gatewayResp.GatewayClassification.(*fabcore.OnPremisesGatewayPersonal)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Gateway Type", "Result is not an OnPremisesGatewayPersonal.")
		return
	}

	data.set(ctx, *realGw)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if diags := resp.State.Set(ctx, data); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}
