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
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceGatewayRoleAssignments)(nil)

type dataSourceGatewayRoleAssignments struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GatewaysClient
	Names       string
	Name        string
	IsPreview   bool
}

func NewDataSourceGatewayRoleAssignments() datasource.DataSource {
	return &dataSourceGatewayRoleAssignments{
		Names:     GatewayRoleAssignmentsName,
		Name:      GatewayRoleAssignmentName,
		IsPreview: ItemPreview,
	}
}

func (d *dataSourceGatewayRoleAssignments) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + GatewayRoleAssignmentsTFName
}

func (d *dataSourceGatewayRoleAssignments) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attributes := getDataSourceGatewayRoleAssignmentAttributes(ctx, true)

	resp.Schema = schema.Schema{
		MarkdownDescription: fabricitem.GetDataSourcePreviewNote("List Fabric "+GatewayRoleAssignmentsName+".\n\n"+
			"Use this data source to list ["+GatewayRoleAssignmentsName+"].\n\n"+
			ItemDocsSPNSupport, d.IsPreview),
		Attributes: map[string]schema.Attribute{
			"gateway_id": schema.StringAttribute{
				MarkdownDescription: "The Gateway ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"values": schema.ListNestedAttribute{
				MarkdownDescription: "The list of " + GatewayRoleAssignmentsName + ".",
				Computed:            true,
				CustomType:          supertypes.NewListNestedObjectTypeOf[baseGatewayRoleAssignmentModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: attributes,
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceGatewayRoleAssignments) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(d.Name, d.IsPreview, d.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}

	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewGatewaysClient()
}

func (d *dataSourceGatewayRoleAssignments) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	var data dataSourceGatewayRoleAssignmentsModel

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

func (d *dataSourceGatewayRoleAssignments) list(ctx context.Context, model *dataSourceGatewayRoleAssignmentsModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting "+GatewayRoleAssignmentsName)

	respList, err := d.client.ListGatewayRoleAssignments(ctx, model.GatewayID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	return model.setValues(ctx, respList)
}
