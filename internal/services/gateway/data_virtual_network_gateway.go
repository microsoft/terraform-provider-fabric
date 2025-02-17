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
	_ datasource.DataSourceWithConfigValidators = (*dataSourceVirtualNetworkGateway)(nil)
	_ datasource.DataSourceWithConfigure        = (*dataSourceVirtualNetworkGateway)(nil)
)

type dataSourceVirtualNetworkGateway struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GatewaysClient
}

func NewDataSourceVirtualNetworkGateway() datasource.DataSource {
	return &dataSourceVirtualNetworkGateway{}
}

func (d *dataSourceVirtualNetworkGateway) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + VirtualNetworkItemTFType
}

func (d *dataSourceVirtualNetworkGateway) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"capacity_id": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("The %s capacity Id.", ItemName),
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"inactivity_minutes_before_sleep": schema.Int32Attribute{
				MarkdownDescription: "The number of minutes of inactivity before the gateway goes to sleep.",
				Computed:            true,
			},
			"number_of_member_gateways": schema.Int32Attribute{
				MarkdownDescription: "The number of member gateways.",
				Computed:            true,
			},
			"virtual_network_azure_resource": schema.SingleNestedAttribute{
				MarkdownDescription: "The Azure resource of the virtual network.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[virtualNetworkAzureResourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"subscription_id": schema.StringAttribute{
						MarkdownDescription: "The subscription ID.",
						Computed:            true,
						CustomType:          customtypes.UUIDType{},
					},
					"resource_group_name": schema.StringAttribute{
						MarkdownDescription: "The name of the resource group.",
						Computed:            true,
					},
					"virtual_network_name": schema.StringAttribute{
						MarkdownDescription: "The name of the virtual network.",
						Computed:            true,
					},
					"subnet_name": schema.StringAttribute{
						MarkdownDescription: "The name of the subnet.",
						Computed:            true,
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceVirtualNetworkGateway) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
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
func (d *dataSourceVirtualNetworkGateway) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the latest data.
func (d *dataSourceVirtualNetworkGateway) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data datasourceVirtualNetworkGatewayModel

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

func (d *dataSourceVirtualNetworkGateway) getByID(ctx context.Context, model *datasourceVirtualNetworkGatewayModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY ID", map[string]any{
		"id": model.ID.ValueString(),
	})

	respGet, err := d.client.GetGateway(ctx, model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	if gw, ok := respGet.GatewayClassification.(*fabcore.VirtualNetworkGateway); ok {
		model.set(ctx, *gw)
		return nil
	}

	return diag.Diagnostics{
		diag.NewErrorDiagnostic(
			common.ErrorReadHeader,
			"expected gateway to be a virtual network gateway",
		),
	}
}

func (d *dataSourceVirtualNetworkGateway) getByDisplayName(ctx context.Context, model *datasourceVirtualNetworkGatewayModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s by 'display_name'", ItemName))

	gateways, err := d.client.ListGateways(ctx, nil)

	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	for _, gw := range gateways {
		if virtualNetworkGateway, ok := gw.(*fabcore.VirtualNetworkGateway); ok {
			if *virtualNetworkGateway.DisplayName == model.DisplayName.ValueString() {
				model.set(ctx, *virtualNetworkGateway)
				return nil
			}
		}
	}

	return diag.Diagnostics{diag.NewErrorDiagnostic(common.ErrorReadHeader, "virtual network gateway not found")}
}
