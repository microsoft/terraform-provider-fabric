// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package spark

import (
	"context"
	"fmt"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceSparkEnvironmentSettings)(nil)

type dataSourceSparkEnvironmentSettings struct {
	pConfigData *pconfig.ProviderData
	client      *fabenvironment.SparkComputeClient
}

func NewDataSourceSparkEnvironmentSettings() datasource.DataSource {
	return &dataSourceSparkEnvironmentSettings{}
}

func (d *dataSourceSparkEnvironmentSettings) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + SparkEnvironmentSettingsTFName
}

func (d *dataSourceSparkEnvironmentSettings) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get a Fabric " + SparkEnvironmentSettingsName + ".\n\n" +
			"See [" + SparkEnvironmentSettingsName + "](" + SparkEnvironmentSettingsDocsURL + ") for more information.\n\n" +
			SparkEnvironmentSettingsDocsSPNSupport,
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
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "The Environment ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"publication_status": schema.StringAttribute{
				MarkdownDescription: "Publication Status. Accepted values: " + utils.ConvertStringSlicesToString(SparkEnvironmentPublicationStatusValues, true, true) + ".",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(SparkEnvironmentPublicationStatusValues...),
				},
			},
			"driver_cores": schema.Int32Attribute{
				MarkdownDescription: "Spark driver core. Possible values: " + utils.ConvertStringSlicesToString(SparkEnvironmentDriverCoresValues, true, false) + ".",
				Computed:            true,
			},
			"driver_memory": schema.StringAttribute{
				MarkdownDescription: "Spark driver memory. Possible values: " + utils.ConvertStringSlicesToString(SparkEnvironmentDriverMemoryValues, true, false) + ".",
				Computed:            true,
			},
			"dynamic_executor_allocation": schema.SingleNestedAttribute{
				MarkdownDescription: "Dynamic executor allocation.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[dynamicExecutorAllocationPropertiesModel](ctx),
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "The status of the dynamic executor allocation. Possible values: `false` - Disabled, `true` - Enabled.",
						Computed:            true,
					},
					"min_executors": schema.Int32Attribute{
						MarkdownDescription: "The minimum executor number.",
						Computed:            true,
					},
					"max_executors": schema.Int32Attribute{
						MarkdownDescription: "The maximum executor number.",
						Computed:            true,
					},
				},
			},
			"executor_cores": schema.Int32Attribute{
				MarkdownDescription: "Spark executor core. Possible values: " + utils.ConvertStringSlicesToString(SparkEnvironmentExecutorCoresValues, true, false) + ".",
				Computed:            true,
			},
			"executor_memory": schema.StringAttribute{
				MarkdownDescription: "Spark executor memory. Possible values: " + utils.ConvertStringSlicesToString(SparkEnvironmentExecutorMemoryValues, true, false) + ".",
				Computed:            true,
			},
			"pool": schema.SingleNestedAttribute{
				MarkdownDescription: "Environment pool.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[instancePoolPropertiesModel](ctx),
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
						MarkdownDescription: "The Pool type. Possible values: " + utils.ConvertStringSlicesToString(fabenvironment.PossibleCustomPoolTypeValues(), true, true) + ".",
						Computed:            true,
					},
				},
			},
			"runtime_version": schema.StringAttribute{
				MarkdownDescription: "[Runtime](https://review.learn.microsoft.com/fabric/data-engineering/runtime) version. Possible values: " + utils.ConvertStringSlicesToString(SparkRuntimeVersionValues, true, false) + ".",
				Description:         "Runtime version. Possible values: " + utils.ConvertStringSlicesToString(SparkRuntimeVersionValues, true, false) + ".",
				Computed:            true,
			},
			"spark_properties": schema.MapAttribute{
				MarkdownDescription: "A map of key/value pairs of Spark properties.",
				Computed:            true,
				CustomType:          supertypes.NewMapTypeOf[string](ctx),
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceSparkEnvironmentSettings) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fabenvironment.NewClientFactoryWithClient(*pConfigData.FabricClient).NewSparkComputeClient()
}

func (d *dataSourceSparkEnvironmentSettings) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceSparkEnvironmentSettingsModel

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

	data.ID = data.EnvironmentID

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dataSourceSparkEnvironmentSettings) get(ctx context.Context, model *dataSourceSparkEnvironmentSettingsModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s for Workspace ID: %s", SparkEnvironmentSettingsName, model.WorkspaceID.ValueString()))

	var respEntity fabenvironment.SparkCompute

	if model.PublicationStatus.ValueString() == SparkEnvironmentPublicationStatusPublished {
		respGet, err := d.client.GetPublishedSettings(ctx, model.WorkspaceID.ValueString(), model.EnvironmentID.ValueString(), nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
			return diags
		}

		respEntity = respGet.SparkCompute
	} else {
		respGet, err := d.client.GetStagingSettings(ctx, model.WorkspaceID.ValueString(), model.EnvironmentID.ValueString(), nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
			return diags
		}

		respEntity = respGet.SparkCompute
	}

	return model.set(ctx, respEntity)
}
