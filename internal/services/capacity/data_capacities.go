// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package capacity

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

var _ datasource.DataSourceWithConfigure = (*dataSourceCapacities)(nil)

type dataSourceCapacities struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.CapacitiesClient
}

func NewDataSourceCapacities() datasource.DataSource {
	return &dataSourceCapacities{}
}

func (d *dataSourceCapacities) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ItemsTFName
}

func (d *dataSourceCapacities) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List a Fabric " + ItemsName + ".\n\n" +
			"Use this data source to list [" + ItemsName + "](" + ItemDocsURL + ") for more information.\n\n" +
			ItemDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"values": schema.ListNestedAttribute{
				MarkdownDescription: fmt.Sprintf("The list of %s.", ItemsName),
				Computed:            true,
				CustomType:          supertypes.NewListNestedObjectTypeOf[baseCapacityModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: fmt.Sprintf("The %s ID.", ItemName),
							Computed:            true,
							CustomType:          customtypes.UUIDType{},
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: fmt.Sprintf("The %s display name.", ItemName),
							Computed:            true,
						},
						"region": schema.StringAttribute{
							MarkdownDescription: "The Azure region where the " + ItemName + " has been provisioned. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleCapacityRegionValues(), true, true),
							Computed:            true,
						},
						"sku": schema.StringAttribute{
							MarkdownDescription: fmt.Sprintf("The %s SKU.", ItemName),
							Computed:            true,
						},
						"state": schema.StringAttribute{
							MarkdownDescription: "The " + ItemName + " state. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleCapacityStateValues(), true, true),
							Computed:            true,
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceCapacities) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewCapacitiesClient()
}

func (d *dataSourceCapacities) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceCapacitiesModel

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

func (d *dataSourceCapacities) list(ctx context.Context, model *dataSourceCapacitiesModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting "+ItemsName)

	respList, err := d.client.ListCapacities(ctx, nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	return model.setValues(ctx, respList)
}
