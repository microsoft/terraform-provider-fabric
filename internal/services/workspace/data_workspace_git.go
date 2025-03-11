// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
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

var _ datasource.DataSourceWithConfigure = (*dataSourceWorkspaceGit)(nil)

type dataSourceWorkspaceGit struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GitClient
}

func NewDataSourceWorkspaceGit() datasource.DataSource {
	return &dataSourceWorkspaceGit{}
}

func (d *dataSourceWorkspaceGit) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + WorkspaceGitTFName
}

func (d *dataSourceWorkspaceGit) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get a Fabric " + WorkspaceGitName + ".",
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The Workspace ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"git_connection_state": schema.StringAttribute{
				MarkdownDescription: "The git connection state. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleGitConnectionStateValues(), true, true),
				Computed:            true,
			},
			"git_sync_details": schema.SingleNestedAttribute{
				MarkdownDescription: "The git sync details.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[gitSyncDetailsModel](ctx),
				Attributes: map[string]schema.Attribute{
					"head": schema.StringAttribute{
						MarkdownDescription: "The git head.",
						Computed:            true,
					},
					"last_sync_time": schema.StringAttribute{
						MarkdownDescription: "The last sync time.",
						Computed:            true,
						CustomType:          timetypes.RFC3339Type{},
					},
				},
			},
			"git_provider_details": schema.SingleNestedAttribute{
				MarkdownDescription: "The Git provider details.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[gitProviderDetailsModel](ctx),
				Attributes: map[string]schema.Attribute{
					"git_provider_type": schema.StringAttribute{
						MarkdownDescription: "The Git provider type. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleGitProviderTypeValues(), true, true),
						Computed:            true,
					},
					"organization_name": schema.StringAttribute{
						MarkdownDescription: "The Azure DevOps organization name.",
						Computed:            true,
					},
					"project_name": schema.StringAttribute{
						MarkdownDescription: "The Azure DevOps project name.",
						Computed:            true,
					},
					"owner_name": schema.StringAttribute{
						MarkdownDescription: "The GitHub owner name.",
						Computed:            true,
					},
					"repository_name": schema.StringAttribute{
						MarkdownDescription: "The repository name.",
						Computed:            true,
					},
					"branch_name": schema.StringAttribute{
						MarkdownDescription: "The branch name.",
						Computed:            true,
					},
					"directory_name": schema.StringAttribute{
						MarkdownDescription: "The directory name.",
						Computed:            true,
					},
				},
			},
			"git_credentials": schema.SingleNestedAttribute{
				MarkdownDescription: "The Git credentials details.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[gitCredentialsModel](ctx),
				Attributes: map[string]schema.Attribute{
					"source": schema.StringAttribute{
						MarkdownDescription: "The Git credentials source. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleGitCredentialsSourceValues(), true, true),
						Computed:            true,
					},
					"connection_id": schema.StringAttribute{
						MarkdownDescription: "The object ID of the connection.",
						Computed:            true,
						CustomType:          customtypes.UUIDType{},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceWorkspaceGit) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewGitClient()
}

func (d *dataSourceWorkspaceGit) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceWorkspaceGitModel

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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dataSourceWorkspaceGit) get(ctx context.Context, model *dataSourceWorkspaceGitModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s for Workspace ID: %s", WorkspaceGitName, model.WorkspaceID.ValueString()))

	respGet, err := d.client.GetConnection(ctx, model.WorkspaceID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	if diags := model.set(ctx, respGet.GitConnection); diags.HasError() {
		return diags
	}

	respGetCredentials, err := d.client.GetMyGitCredentials(ctx, model.WorkspaceID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	if diags := model.setCredentials(ctx, respGetCredentials.GitCredentialsConfigurationResponseClassification); diags.HasError() {
		return diags
	}

	return nil
}
