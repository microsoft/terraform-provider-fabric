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

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceGatewayRoleAssignment)(nil)

type dataSourceGatewayRoleAssignment struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GatewaysClient
	Name        string
	IsPreview   bool
}

func NewDataSourceGatewayRoleAssignment() datasource.DataSource {
	return &dataSourceGatewayRoleAssignment{
		Name:      GatewayRoleAssignmentName,
		IsPreview: ItemPreview,
	}
}

func (d *dataSourceGatewayRoleAssignment) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + GatewayRoleAssignmentTFName
}

func (d *dataSourceGatewayRoleAssignment) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attributes := getDataSourceGatewayRoleAssignmentAttributes(ctx, false)
	attributes["timeouts"] = timeouts.Attributes(ctx)

	attributes["gateway_id"] = schema.StringAttribute{
		MarkdownDescription: "The Gateway ID.",
		Required:            true,
		CustomType:          customtypes.UUIDType{},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: fabricitem.GetDataSourcePreviewNote("Get a Fabric "+GatewayRoleAssignmentName+".\n\n"+
			"Use this data source to get ["+GatewayRoleAssignmentName+"].\n\n"+
			ItemDocsSPNSupport, d.IsPreview),
		Attributes: attributes,
	}
}

func (d *dataSourceGatewayRoleAssignment) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dataSourceGatewayRoleAssignment) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceGatewayRoleAssignmentModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, d.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(d.get(ctx, &data)...); resp.Diagnostics.HasError() {
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

func (d *dataSourceGatewayRoleAssignment) get(ctx context.Context, model *dataSourceGatewayRoleAssignmentModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET "+GatewayRoleAssignmentName, map[string]any{
		"id": model.ID.ValueString(),
	})

	respGet, err := d.client.GetGatewayRoleAssignment(ctx, model.GatewayID.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	if diags := model.set(ctx, respGet.GatewayRoleAssignment); diags.HasError() {
		return diags
	}

	return nil
}
