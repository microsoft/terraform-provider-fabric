// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

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

var _ datasource.DataSourceWithConfigure = (*dataSourceWorkspaceRoleAssignments)(nil)

type dataSourceWorkspaceRoleAssignments struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.WorkspacesClient
}

func NewDataSourceWorkspaceRoleAssignments() datasource.DataSource {
	return &dataSourceWorkspaceRoleAssignments{}
}

func (d *dataSourceWorkspaceRoleAssignments) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + WorkspaceRoleAssignmentsTFName
}

func (d *dataSourceWorkspaceRoleAssignments) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List a Fabric " + WorkspaceRoleAssignmentsName + ".\n\n" +
			"Use this data source to list [" + WorkspaceRoleAssignmentsName + "](" + WorkspaceRoleAssignmentDocsURL + ").\n\n" +
			ItemDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The Workspace ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"values": schema.ListNestedAttribute{
				MarkdownDescription: "The list of " + WorkspaceRoleAssignmentsName + ".",
				Computed:            true,
				CustomType:          supertypes.NewListNestedObjectTypeOf[baseWorkspaceRoleAssignmentModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The " + WorkspaceRoleAssignmentName + " ID.",
							Computed:            true,
							CustomType:          customtypes.UUIDType{},
						},
						"principal_id": schema.StringAttribute{
							MarkdownDescription: "The Principal ID.",
							Computed:            true,
							CustomType:          customtypes.UUIDType{},
						},
						"role": schema.StringAttribute{
							MarkdownDescription: "The workspace role of the principal. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleWorkspaceRoleValues(), true, true) + ".",
							Computed:            true,
						},
						"principal_display_name": schema.StringAttribute{
							MarkdownDescription: "The principal's display name.",
							Computed:            true,
						},
						"principal_type": schema.StringAttribute{
							MarkdownDescription: "The type of the principal. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossiblePrincipalTypeValues(), true, true) + ".",
							Computed:            true,
						},
						"principal_details": schema.SingleNestedAttribute{
							MarkdownDescription: "The principal details.",
							Computed:            true,
							CustomType:          supertypes.NewSingleNestedObjectTypeOf[principalDetailsModel](ctx),
							Attributes: map[string]schema.Attribute{
								"user_principal_name": schema.StringAttribute{
									MarkdownDescription: "The user principal name.",
									Computed:            true,
								},
								"group_type": schema.StringAttribute{
									MarkdownDescription: "The type of the group. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleGroupTypeValues(), true, true) + ".",
									Computed:            true,
								},
								"app_id": schema.StringAttribute{
									MarkdownDescription: "The service principal's Microsoft Entra App ID.",
									Computed:            true,
									CustomType:          customtypes.UUIDType{},
								},
								"parent_principal_id": schema.StringAttribute{
									MarkdownDescription: "The parent principal ID of Service Principal Profile.",
									Computed:            true,
									CustomType:          customtypes.UUIDType{},
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

func (d *dataSourceWorkspaceRoleAssignments) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewWorkspacesClient()
}

func (d *dataSourceWorkspaceRoleAssignments) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceWorkspaceRoleAssignmentsModel

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

func (d *dataSourceWorkspaceRoleAssignments) list(ctx context.Context, model *dataSourceWorkspaceRoleAssignmentsModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting "+WorkspaceRoleAssignmentsName)

	respList, err := d.client.ListWorkspaceRoleAssignments(ctx, model.WorkspaceID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	return model.setValues(ctx, respList)
}
