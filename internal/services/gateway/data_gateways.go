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

var _ datasource.DataSourceWithConfigure = (*dataSourceGateways)(nil)

type dataSourceGateways struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GatewaysClient
}

func NewDataSourceGateways() datasource.DataSource {
	return &dataSourceGateways{}
}

func (d *dataSourceGateways) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ItemsTFName
}

func (d *dataSourceGateways) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List a Fabric " + ItemsName + ".\n\n" +
			"Use this data source to list [" + ItemsName + "](" + ItemDocsURL + ").\n\n" +
			ItemDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"values": schema.ListNestedAttribute{
				MarkdownDescription: "The list of " + ItemsName + ".",
				Computed:            true,
				CustomType:          supertypes.NewListNestedObjectTypeOf[baseDataSourceGatewayModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The " + ItemName + " ID.",
							Optional:            true,
							Computed:            true,
							CustomType:          customtypes.UUIDType{},
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The " + ItemName + " display name.",
							Optional:            true,
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The " + ItemName + " type. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleGatewayTypeValues(), true, true),
							Computed:            true,
						},
						"capacity_id": schema.StringAttribute{
							MarkdownDescription: "The " + ItemName + " capacity ID.",
							Computed:            true,
							CustomType:          customtypes.UUIDType{},
						},
						"inactivity_minutes_before_sleep": schema.Int32Attribute{
							MarkdownDescription: "The " + ItemName + " inactivity minutes before sleep. Possible values: " + utils.ConvertStringSlicesToString(PossibleInactivityMinutesBeforeSleepValues, true, true),
							Computed:            true,
						},
						"number_of_member_gateways": schema.Int32Attribute{
							MarkdownDescription: "The number of member gateways. Possible values: " + fmt.Sprint(MinNumberOfMemberGatewaysValues) + " to " + fmt.Sprint(MaxNumberOfMemberGatewaysValues) + ".",
							Computed:            true,
						},
						"virtual_network_azure_resource": schema.SingleNestedAttribute{
							MarkdownDescription: "The Azure virtual network resource.",
							Computed:            true,
							CustomType:          supertypes.NewSingleNestedObjectTypeOf[virtualNetworkAzureResourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"resource_group_name": schema.StringAttribute{
									MarkdownDescription: "The resource group name.",
									Computed:            true,
								},
								"subnet_name": schema.StringAttribute{
									MarkdownDescription: "The subnet name.",
									Computed:            true,
								},
								"subscription_id": schema.StringAttribute{
									MarkdownDescription: "The subscription ID.",
									Computed:            true,
									CustomType:          customtypes.UUIDType{},
								},
								"virtual_network_name": schema.StringAttribute{
									MarkdownDescription: "The virtual network name.",
									Computed:            true,
								},
							},
						},
						"allow_cloud_connection_refresh": schema.BoolAttribute{
							MarkdownDescription: "Allow cloud connection refresh.",
							Computed:            true,
						},
						"allow_custom_connectors": schema.BoolAttribute{
							MarkdownDescription: "Allow custom connectors.",
							Computed:            true,
						},
						"load_balancing_setting": schema.StringAttribute{
							MarkdownDescription: "The load balancing setting. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleLoadBalancingSettingValues(), true, true),
							Computed:            true,
						},
						"public_key": schema.SingleNestedAttribute{
							MarkdownDescription: "The public key of the primary gateway member. Used to encrypt the credentials for creating and updating connections.",
							Computed:            true,
							CustomType:          supertypes.NewSingleNestedObjectTypeOf[publicKeyModel](ctx),
							Attributes: map[string]schema.Attribute{
								"exponent": schema.StringAttribute{
									MarkdownDescription: "The exponent.",
									Computed:            true,
								},
								"modulus": schema.StringAttribute{
									MarkdownDescription: "The modulus.",
									Computed:            true,
								},
							},
						},
						"version": schema.StringAttribute{
							MarkdownDescription: "The " + ItemName + " version.",
							Computed:            true,
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceGateways) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dataSourceGateways) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceGatewaysModel

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

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dataSourceGateways) list(ctx context.Context, model *dataSourceGatewaysModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting "+ItemsName)

	respList, err := d.client.ListGateways(ctx, nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	return model.setValues(ctx, respList)
}
