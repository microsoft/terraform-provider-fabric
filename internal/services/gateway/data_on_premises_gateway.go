// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var (
	_ datasource.DataSourceWithConfigValidators = (*dataSourceOnPremisesGateway)(nil)
	_ datasource.DataSourceWithConfigure        = (*dataSourceOnPremisesGateway)(nil)
)

type dataSourceOnPremisesGateway struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GatewaysClient
}

func NewDataSourceOnPremisesGateway() datasource.DataSource {
	return &dataSourceOnPremisesGateway{}
}

func (d *dataSourceOnPremisesGateway) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + OnPremisesItemTFType
}

func (d *dataSourceOnPremisesGateway) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get a Fabric " + ItemName + ".\n\n" +
			"Use this data source to fetch [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
			ItemDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("The %s ID.", ItemName),
				Optional:            true,
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("The %s display name.", ItemName),
				Optional:            true,
				Computed:            true,
			},
			"allow_cloud_connection_refresh": schema.BoolAttribute{
				MarkdownDescription: "Defines if cloud connection refresh is allowed.",
				Computed:            true,
			},
			"allow_custom_connectors": schema.BoolAttribute{
				MarkdownDescription: "Defines if custom connectors are allowed.",
				Computed:            true,
			},
			"load_balancing_setting": schema.StringAttribute{
				MarkdownDescription: "Gateway load balancing setting.",
				Computed:            true,
			},
			"number_of_member_gateways": schema.Int32Attribute{
				MarkdownDescription: "The number of member gateways.",
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
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceOnPremisesGateway) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("id"),
			path.MatchRoot("display_name"),
		),
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("display_name"),
		),
	}
}

// Configure adds the provider configured client to the data source.
func (d *dataSourceOnPremisesGateway) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the latest data.
func (d *dataSourceOnPremisesGateway) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data datasourceOnPremisesGatewayModel

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

func (d *dataSourceOnPremisesGateway) getByID(ctx context.Context, model *datasourceOnPremisesGatewayModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY ID", map[string]any{
		"id": model.ID.ValueString(),
	})

	respGet, err := d.client.GetGateway(ctx, model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	if gw, ok := respGet.GatewayClassification.(*fabcore.OnPremisesGateway); ok {
		model.set(ctx, *gw)
		return nil
	} else {
		var diags diag.Diagnostics
		diags.AddError(common.ErrorReadHeader, "expected gateway to be an on-premises gateway")
		return diags
	}
}

func (d *dataSourceOnPremisesGateway) getByDisplayName(ctx context.Context, model *datasourceOnPremisesGatewayModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s by 'display_name'", ItemName))

	gateways, err := d.client.ListGateways(ctx, nil)

	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	for _, gw := range gateways {
		if OnPremisesGateway, ok := gw.(*fabcore.OnPremisesGateway); ok {
			if *OnPremisesGateway.DisplayName == model.DisplayName.ValueString() {
				model.set(ctx, *OnPremisesGateway)
				return nil
			}
		}
	}

	var diags diag.Diagnostics
	diags.AddError(common.ErrorReadHeader, "no on-premises gateway with display name found")
	return diags
}
