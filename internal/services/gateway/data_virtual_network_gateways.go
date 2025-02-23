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

var _ datasource.DataSourceWithConfigure = (*dataSourceVirtualNetworkGateways)(nil)

type dataSourceVirtualNetworkGateways struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GatewaysClient
}

func NewDataSourceVirtualNetworkGateways() datasource.DataSource {
	return &dataSourceVirtualNetworkGateways{}
}

func (d *dataSourceVirtualNetworkGateways) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + VirtualNetworkItemsTFType
}

func (d *dataSourceVirtualNetworkGateways) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all Fabric Virtual Network Gateways.",
		Attributes: map[string]schema.Attribute{
			"values": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "A list of Virtual Network Gateways.",
				CustomType:          supertypes.NewListNestedObjectTypeOf[virtualNetworkGatewayModelBase](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The gateway ID.",
							Computed:            true,
							CustomType:          customtypes.UUIDType{},
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The display name of the gateway.",
							Computed:            true,
						},
						"capacity_id": schema.StringAttribute{
							MarkdownDescription: "The Fabric license capacity ID.",
							Computed:            true,
							CustomType:          customtypes.UUIDType{},
						},
						"inactivity_minutes_before_sleep": schema.Int64Attribute{
							MarkdownDescription: "Minutes of inactivity before auto-sleep.",
							Computed:            true,
						},
						"number_of_member_gateways": schema.Int64Attribute{
							MarkdownDescription: "The number of member gateways.",
							Computed:            true,
						},
						"virtual_network_azure_resource": schema.SingleNestedAttribute{
							MarkdownDescription: "The Azure resource details for this gateway's Virtual Network.",
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
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceVirtualNetworkGateways) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dataSourceVirtualNetworkGateways) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ-ALL-Virtual-Network-Gateways", map[string]any{"action": "start"})

	var data dataSourceVirtualNetworkGatewaysModel
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
	tflog.Debug(ctx, "READ-ALL-Virtual-Network-Gateways", map[string]any{"action": "end"})
}

func (d *dataSourceVirtualNetworkGateways) list(ctx context.Context, model *dataSourceVirtualNetworkGatewaysModel) diag.Diagnostics {
	tflog.Trace(ctx, "Listing all virtual network gateways")

	gatewaysResp, err := d.client.ListGateways(ctx, nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	return model.setValues(ctx, gatewaysResp)
}
