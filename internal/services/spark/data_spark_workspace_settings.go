// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package spark

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabspark "github.com/microsoft/fabric-sdk-go/fabric/spark"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceSparkWorkspaceSettings)(nil)

type dataSourceSparkWorkspaceSettings struct {
	pConfigData *pconfig.ProviderData
	client      *fabspark.WorkspaceSettingsClient
}

func NewDataSourceSparkWorkspaceSettings() datasource.DataSource {
	return &dataSourceSparkWorkspaceSettings{}
}

func (d *dataSourceSparkWorkspaceSettings) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + SparkWorkspaceSettingsTFName
}

func (d *dataSourceSparkWorkspaceSettings) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get a Fabric " + SparkWorkspaceSettingsName + ".\n\n" +
			"See [" + SparkWorkspaceSettingsName + "](" + SparkWorkspaceSettingsDocsURL + ") for more information.\n\n" +
			SparkWorkspaceSettingsDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: customtypes.UUIDType{},
			},
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The Workspace ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"automatic_log": schema.SingleNestedAttribute{
				MarkdownDescription: "Automatic Log properties.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[automaticLogPropertiesModel](ctx),
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "The status of the automatic log. Possible values: `false` - Disabled, `true` - Enabled.",
						Computed:            true,
					},
				},
			},
			"environment": schema.SingleNestedAttribute{
				MarkdownDescription: "Environment properties.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[environmentPropertiesModel](ctx),
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the default environment. Empty indicated there is no workspace default environment.",
						Computed:            true,
					},
					"runtime_version": schema.StringAttribute{
						MarkdownDescription: "[Runtime](https://review.learn.microsoft.com/fabric/data-engineering/runtime) version. Possible values: " + utils.ConvertStringSlicesToString(SparkRuntimeVersionValues, true, false) + ".",
						Description:         "Runtime version. Possible values: " + utils.ConvertStringSlicesToString(SparkRuntimeVersionValues, true, false) + ".",
						Computed:            true,
					},
				},
			},
			"high_concurrency": schema.SingleNestedAttribute{
				MarkdownDescription: "High Concurrency properties.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[highConcurrencyPropertiesModel](ctx),
				Attributes: map[string]schema.Attribute{
					"notebook_interactive_run_enabled": schema.BoolAttribute{
						MarkdownDescription: "The status of the high concurrency for notebook interactive run. `false` - Disabled, `true` - Enabled.",
						Computed:            true,
					},
				},
			},
			"pool": schema.SingleNestedAttribute{
				MarkdownDescription: "Pool properties.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[poolPropertiesModel](ctx),
				Attributes: map[string]schema.Attribute{
					"customize_compute_enabled": schema.BoolAttribute{
						MarkdownDescription: "Customize compute configurations for items. `false` - Disabled, `true` - Enabled.",
						Computed:            true,
					},
					"default_pool": schema.SingleNestedAttribute{
						MarkdownDescription: "Default pool for workspace.",
						Computed:            true,
						CustomType:          supertypes.NewSingleNestedObjectTypeOf[defaultPoolPropertiesModel](ctx),
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								MarkdownDescription: "The Pool ID.",
								Computed:            true,
								CustomType:          customtypes.UUIDType{},
							},
							"name": schema.StringAttribute{
								MarkdownDescription: "The Pool name. `Starter Pool` means using the starting pool.",
								Computed:            true,
							},
							"type": schema.StringAttribute{
								MarkdownDescription: "The Pool type. Possible values: " + utils.ConvertStringSlicesToString(fabspark.PossibleCustomPoolTypeValues(), true, true) + ".",
								Computed:            true,
							},
						},
					},
					"starter_pool": schema.SingleNestedAttribute{
						MarkdownDescription: "Starter pool for workspace.",
						Computed:            true,
						CustomType:          supertypes.NewSingleNestedObjectTypeOf[starterPoolPropertiesModel](ctx),
						Attributes: map[string]schema.Attribute{
							"max_node_count": schema.Int32Attribute{
								MarkdownDescription: "The maximum node count.",
								Computed:            true,
							},
							"max_executors": schema.Int32Attribute{
								MarkdownDescription: "The maximum executors count.",
								Computed:            true,
							},
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceSparkWorkspaceSettings) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fabspark.NewClientFactoryWithClient(*pConfigData.FabricClient).NewWorkspaceSettingsClient()
}

func (d *dataSourceSparkWorkspaceSettings) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceSparkWorkspaceSettingsModel

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

	data.ID = data.WorkspaceID

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dataSourceSparkWorkspaceSettings) get(ctx context.Context, model *dataSourceSparkWorkspaceSettingsModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s for Workspace ID: %s", SparkWorkspaceSettingsName, model.WorkspaceID.ValueString()))

	respGet, err := d.client.GetSparkSettings(ctx, model.WorkspaceID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	return model.set(ctx, respGet.WorkspaceSparkSettings)
}
