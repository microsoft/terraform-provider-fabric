// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package spark

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabspark "github.com/microsoft/fabric-sdk-go/fabric/spark"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var (
	_ datasource.DataSourceWithConfigValidators = (*dataSourceSparkCustomPool)(nil)
	_ datasource.DataSourceWithConfigure        = (*dataSourceSparkCustomPool)(nil)
)

type dataSourceSparkCustomPool struct {
	pConfigData *pconfig.ProviderData
	client      *fabspark.CustomPoolsClient
}

func NewDataSourceSparkCustomPool() datasource.DataSource {
	return &dataSourceSparkCustomPool{}
}

func (d *dataSourceSparkCustomPool) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + SparkCustomPoolTFName
}

func (d *dataSourceSparkCustomPool) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get a Fabric " + SparkCustomPoolName + ".\n\n" +
			"See [" + SparkCustomPoolName + "](" + SparkCustomPoolDocsURL + ") for more information.\n\n" +
			SparkCustomPoolDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The Workspace ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The " + SparkCustomPoolName + " ID.",
				Optional:            true,
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The " + SparkCustomPoolName + " name.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(64),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9-_ ]+$`), "The name must contain only letters, numbers, dashes, underscores and spaces."),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The " + SparkCustomPoolName + " type. Possible values: " + utils.ConvertStringSlicesToString(fabspark.PossibleCustomPoolTypeValues(), true, true) + ".",
				Computed:            true,
			},
			"node_family": schema.StringAttribute{
				MarkdownDescription: "The Node family. Possible values: " + utils.ConvertStringSlicesToString(fabspark.PossibleNodeFamilyValues(), true, true) + ".",
				Computed:            true,
			},
			"node_size": schema.StringAttribute{
				MarkdownDescription: "The Node size. Possible values: " + utils.ConvertStringSlicesToString(fabspark.PossibleNodeSizeValues(), true, true) + ".",
				Computed:            true,
			},
			"auto_scale": schema.SingleNestedAttribute{
				MarkdownDescription: "Auto-scale properties.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[sparkCustomPoolAutoScaleModel](ctx),
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "The status of the auto scale. Possible values: `false` - Disabled, `true` - Enabled.",
						Computed:            true,
					},
					"min_node_count": schema.Int32Attribute{
						MarkdownDescription: "The minimum node count.",
						Computed:            true,
					},
					"max_node_count": schema.Int32Attribute{
						MarkdownDescription: "The maximum node count.",
						Computed:            true,
					},
				},
			},
			"dynamic_executor_allocation": schema.SingleNestedAttribute{
				MarkdownDescription: "Dynamic Executor Allocation properties.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[sparkCustomPoolDynamicExecutorAllocationModel](ctx),
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "The status of the dynamic executor allocation. Possible values: `false` - Disabled, `true` - Enabled.",
						Computed:            true,
					},
					"min_executors": schema.Int32Attribute{
						MarkdownDescription: "The minimum executors.",
						Computed:            true,
					},
					"max_executors": schema.Int32Attribute{
						MarkdownDescription: "The maximum executors.",
						Computed:            true,
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceSparkCustomPool) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

func (d *dataSourceSparkCustomPool) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fabspark.NewClientFactoryWithClient(*pConfigData.FabricClient).NewCustomPoolsClient()
}

func (d *dataSourceSparkCustomPool) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceSparkCustomPoolModel

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

func (d *dataSourceSparkCustomPool) getByID(ctx context.Context, model *dataSourceSparkCustomPoolModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s by 'id'", SparkCustomPoolName))

	respGet, err := d.client.GetWorkspaceCustomPool(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	return model.set(ctx, respGet.CustomPool)
}

func (d *dataSourceSparkCustomPool) getByDisplayName(ctx context.Context, model *dataSourceSparkCustomPoolModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s by 'display_name'", SparkCustomPoolName))

	var diags diag.Diagnostics

	pager := d.client.NewListWorkspaceCustomPoolsPager(model.WorkspaceID.ValueString(), nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Value {
			if *entity.Name == model.Name.ValueString() {
				return model.set(ctx, entity)
			}
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		fmt.Sprintf("Unable to find "+SparkCustomPoolName+" with name: '%s' in the Workspace ID: %s ", model.Name.ValueString(), model.WorkspaceID.ValueString()),
	)

	return diags
}
