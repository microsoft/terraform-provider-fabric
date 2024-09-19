// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package environment

import (
	"context"
	"fmt"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigValidators = (*dataSourceEnvironment)(nil)
	_ datasource.DataSourceWithConfigure        = (*dataSourceEnvironment)(nil)
)

type dataSourceEnvironment struct {
	pConfigData *pconfig.ProviderData
	client      *fabenvironment.ItemsClient
}

func NewDataSourceEnvironment() datasource.DataSource {
	return &dataSourceEnvironment{}
}

func (d *dataSourceEnvironment) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ItemTFName
}

func (d *dataSourceEnvironment) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	markdownDescription := "Get a Fabric " + ItemName + ".\n\n" +
		"Use this data source to fetch an [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
		ItemDocsSPNSupport

	publishStatePossibleValuesMarkdown := "Publish state. Possible values: " + utils.ConvertStringSlicesToString(fabenvironment.PossibleEnvironmentPublishStateValues(), true, true) + "."

	properties := schema.SingleNestedAttribute{
		MarkdownDescription: "The " + ItemName + " properties.",
		Computed:            true,
		CustomType:          supertypes.NewSingleNestedObjectTypeOf[environmentPropertiesModel](ctx),
		Attributes: map[string]schema.Attribute{
			"publish_details": schema.SingleNestedAttribute{
				MarkdownDescription: "Environment publish operation details.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[environmentPublishDetailsModel](ctx),
				Attributes: map[string]schema.Attribute{
					"state": schema.StringAttribute{
						MarkdownDescription: publishStatePossibleValuesMarkdown,
						Computed:            true,
					},
					"target_version": schema.StringAttribute{
						MarkdownDescription: "Target version to be published.",
						Computed:            true,
						CustomType:          customtypes.UUIDType{},
					},
					"start_time": schema.StringAttribute{
						MarkdownDescription: "Start time of publish operation.",
						Computed:            true,
						CustomType:          timetypes.RFC3339Type{},
					},
					"end_time": schema.StringAttribute{
						MarkdownDescription: "End time of publish operation.",
						Computed:            true,
						CustomType:          timetypes.RFC3339Type{},
					},
					"component_publish_info": schema.SingleNestedAttribute{
						MarkdownDescription: "Environment component publish information.",
						Computed:            true,
						CustomType:          supertypes.NewSingleNestedObjectTypeOf[environmentComponentPublishInfoModel](ctx),
						Attributes: map[string]schema.Attribute{
							"spark_libraries": schema.SingleNestedAttribute{
								MarkdownDescription: "Spark libraries publish information.",
								Computed:            true,
								CustomType:          supertypes.NewSingleNestedObjectTypeOf[environmentSparkLibrariesModel](ctx),
								Attributes: map[string]schema.Attribute{
									"state": schema.StringAttribute{
										MarkdownDescription: publishStatePossibleValuesMarkdown,
										Computed:            true,
									},
								},
							},
							"spark_settings": schema.SingleNestedAttribute{
								MarkdownDescription: "Spark settings publish information.",
								Computed:            true,
								CustomType:          supertypes.NewSingleNestedObjectTypeOf[environmentSparkSettingsModel](ctx),
								Attributes: map[string]schema.Attribute{
									"state": schema.StringAttribute{
										MarkdownDescription: publishStatePossibleValuesMarkdown,
										Computed:            true,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	resp.Schema = fabricitem.GetDataSourceFabricItemPropertiesSchema(ctx, ItemName, markdownDescription, true, properties)
}

func (d *dataSourceEnvironment) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
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

func (d *dataSourceEnvironment) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fabenvironment.NewClientFactoryWithClient(*pConfigData.FabricClient).NewItemsClient()
}

func (d *dataSourceEnvironment) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceEnvironmentModel

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

func (d *dataSourceEnvironment) getByID(ctx context.Context, model *dataSourceEnvironmentModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting Environment by 'id'")

	respGet, err := d.client.GetEnvironment(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	model.set(respGet.Environment)
	model.setProperties(ctx, respGet.Environment)

	return nil
}

func (d *dataSourceEnvironment) getByDisplayName(ctx context.Context, model *dataSourceEnvironmentModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting Environment by 'display_name'")

	var diags diag.Diagnostics

	pager := d.client.NewListEnvironmentsPager(model.WorkspaceID.ValueString(), nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Value {
			if *entity.DisplayName == model.DisplayName.ValueString() {
				model.set(entity)
				model.setProperties(ctx, entity)

				return nil
			}
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		fmt.Sprintf("Unable to find Environment with 'display_name': %s in the Workspace ID: %s ", model.DisplayName.ValueString(), model.WorkspaceID.ValueString()),
	)

	return diags
}
