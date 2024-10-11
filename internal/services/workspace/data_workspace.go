// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"
	"fmt"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigure        = (*dataSourceWorkspace)(nil)
	_ datasource.DataSourceWithConfigValidators = (*dataSourceWorkspace)(nil)
)

// DataSource is the data source for the Fabric Workspace.
type dataSourceWorkspace struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.WorkspacesClient
}

// NewDataSource creates a new data source for the Fabric Workspace.
func NewDataSourceWorkspace() datasource.DataSource {
	return &dataSourceWorkspace{}
}

// Metadata sets metadata for the Fabric Workspace data source.
func (d *dataSourceWorkspace) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ItemTFName
}

// Schema sets the schema for the Fabric Workspace data source.
func (d *dataSourceWorkspace) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get a Fabric Workspace.\n\n" +
			"Use this data source to fetch a [Workspace](https://learn.microsoft.com/fabric/get-started/workspaces).\n\n" +
			ItemDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Workspace ID.",
				Optional:            true,
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The Workspace display name.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The Workspace description.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The Workspace type.",
				Computed:            true,
			},
			"capacity_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the Capacity the Workspace is assigned to.",
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"capacity_region": schema.StringAttribute{
				MarkdownDescription: "The region of the capacity associated with this workspace. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleCapacityRegionValues(), true, true),
				Computed:            true,
			},
			"capacity_assignment_progress": schema.StringAttribute{
				MarkdownDescription: "A Workspace assignment to capacity progress status. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleCapacityAssignmentProgressValues(), true, true),
				Computed:            true,
			},
			"onelake_endpoints": schema.SingleNestedAttribute{
				MarkdownDescription: "The OneLake API endpoints associated with this workspace.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[oneLakeEndpointsModel](ctx),
				Attributes: map[string]schema.Attribute{
					"blob_endpoint": schema.StringAttribute{
						MarkdownDescription: "The OneLake API endpoint available for Blob API operations.",
						Computed:            true,
						CustomType:          customtypes.URLType{},
					},
					"dfs_endpoint": schema.StringAttribute{
						MarkdownDescription: "The OneLake API endpoint available for Distributed File System (DFS) or ADLSgen2 filesystem API operations.",
						Computed:            true,
						CustomType:          customtypes.URLType{},
					},
				},
			},
			"identity": schema.SingleNestedAttribute{
				MarkdownDescription: "A workspace identity object.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[workspaceIdentityModel](ctx),
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The workspace identity type. Possible values: " + utils.ConvertStringSlicesToString(workspaceIdentityTypes, true, true) + ".",
						Computed:            true,
					},
					"application_id": schema.StringAttribute{
						MarkdownDescription: "The application ID.",
						Computed:            true,
						CustomType:          customtypes.UUIDType{},
					},
					"service_principal_id": schema.StringAttribute{
						MarkdownDescription: "The service principal ID.",
						Computed:            true,
						CustomType:          customtypes.UUIDType{},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceWorkspace) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
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

func (d *dataSourceWorkspace) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dataSourceWorkspace) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceWorkspaceModel

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

func (d *dataSourceWorkspace) getByID(ctx context.Context, model *dataSourceWorkspaceModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY ID", map[string]any{
		"id": model.ID.ValueString(),
	})

	respGet, err := d.client.GetWorkspace(ctx, model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	if diags := checkWorkspaceType(respGet.WorkspaceInfo); diags.HasError() {
		return diags
	}

	model.set(ctx, respGet.WorkspaceInfo)

	return nil
}

func (d *dataSourceWorkspace) getByDisplayName(ctx context.Context, model *dataSourceWorkspaceModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY DISPLAY NAME", map[string]any{
		"display_name": model.DisplayName.ValueString(),
	})

	var diags diag.Diagnostics

	pager := d.client.NewListWorkspacesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Value {
			if *entity.DisplayName == model.DisplayName.ValueString() {
				model.ID = customtypes.NewUUIDPointerValue(entity.ID)

				return d.getByID(ctx, model)
			}
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		"Unable to find Workspace with 'display_name': "+model.DisplayName.ValueString(),
	)

	return diags
}
