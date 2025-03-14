// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package spark

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure        = (*resourceSparkEnvironmentSettings)(nil)
	_ resource.ResourceWithConfigValidators = (*resourceSparkEnvironmentSettings)(nil)
)

type resourceSparkEnvironmentSettings struct {
	pConfigData     *pconfig.ProviderData
	client          *fabenvironment.SparkComputeClient
	clientLibraries *fabenvironment.SparkLibrariesClient
}

func NewResourceSparkEnvironmentSettings() resource.Resource {
	return &resourceSparkEnvironmentSettings{}
}

func (r *resourceSparkEnvironmentSettings) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + SparkEnvironmentSettingsTFName
}

func (r *resourceSparkEnvironmentSettings) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage a Fabric " + SparkEnvironmentSettingsName + ".\n\n" +
			"See [" + SparkEnvironmentSettingsName + "](" + SparkEnvironmentSettingsDocsURL + ") for more information.\n\n" +
			SparkEnvironmentSettingsDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The Workspace ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "The Environment ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"publication_status": schema.StringAttribute{
				MarkdownDescription: "Publication Status. Accepted values: " + utils.ConvertStringSlicesToString(SparkEnvironmentPublicationStatusValues, true, true) + ".",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(SparkEnvironmentPublicationStatusValues...),
				},
			},
			"driver_cores": schema.Int32Attribute{
				MarkdownDescription: "Spark driver core. Accepted values: " + utils.ConvertStringSlicesToString(SparkEnvironmentDriverCoresValues, true, false) + ".",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int32{
					int32validator.OneOf(SparkEnvironmentDriverCoresValues...),
				},
			},
			"driver_memory": schema.StringAttribute{
				MarkdownDescription: "Spark driver memory. Accepted values: " + utils.ConvertStringSlicesToString(SparkEnvironmentDriverMemoryValues, true, false) + ".",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dynamic_executor_allocation": schema.SingleNestedAttribute{
				MarkdownDescription: "Dynamic executor allocation.",
				Optional:            true,
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[dynamicExecutorAllocationPropertiesModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "The status of the dynamic executor allocation. Accepted values: `false` - Disabled, `true` - Enabled.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"min_executors": schema.Int32Attribute{
						MarkdownDescription: "The minimum executor number.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Int32{
							int32planmodifier.UseStateForUnknown(),
						},
						Validators: []validator.Int32{
							int32validator.AtLeast(1),
							int32validator.AlsoRequires(path.MatchRelative().AtParent().AtName("max_executors")),
						},
					},
					"max_executors": schema.Int32Attribute{
						MarkdownDescription: "The maximum executor number.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Int32{
							int32planmodifier.UseStateForUnknown(),
						},
						Validators: []validator.Int32{
							int32validator.AtLeast(1),
							int32validator.AlsoRequires(path.MatchRelative().AtParent().AtName("min_executors")),
						},
					},
				},
			},
			"executor_cores": schema.Int32Attribute{
				MarkdownDescription: "Spark executor core. Accepted values: " + utils.ConvertStringSlicesToString(SparkEnvironmentExecutorCoresValues, true, false) + ".",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int32{
					int32validator.OneOf(SparkEnvironmentExecutorCoresValues...),
				},
			},
			"executor_memory": schema.StringAttribute{
				MarkdownDescription: "Spark executor memory. Accepted values: " + utils.ConvertStringSlicesToString(SparkEnvironmentExecutorMemoryValues, true, false) + ".",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(SparkEnvironmentExecutorMemoryValues...),
				},
			},
			"pool": schema.SingleNestedAttribute{
				MarkdownDescription: "Environment pool.",
				Optional:            true,
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[instancePoolPropertiesModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "The Pool ID.",
						Computed:            true,
						CustomType:          customtypes.UUIDType{},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The Pool name. `Starter Pool` means use the starting pool.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("type")),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "The Pool type. Accepted values: " + utils.ConvertStringSlicesToString(utils.ConvertEnumsToStringSlices(fabenvironment.PossibleCustomPoolTypeValues(), false), true, true) + ".",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("name")),
							stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabenvironment.PossibleCustomPoolTypeValues(), false)...),
						},
					},
				},
			},
			"runtime_version": schema.StringAttribute{
				MarkdownDescription: "[Runtime](https://review.learn.microsoft.com/fabric/data-engineering/runtime) version. Possible values: " + utils.ConvertStringSlicesToString(SparkRuntimeVersionValues, true, true) + ".",
				Description:         "Runtime version. Accepted values: " + utils.ConvertStringSlicesToString(SparkRuntimeVersionValues, true, true) + ".",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(SparkRuntimeVersionValues...),
				},
			},
			"spark_properties": schema.MapAttribute{
				MarkdownDescription: "A map of key/value pairs of Spark properties.",
				Optional:            true,
				Computed:            true,
				CustomType:          supertypes.NewMapTypeOf[string](ctx),
				Validators: []validator.Map{
					mapvalidator.SizeAtLeast(1),
					mapvalidator.KeysAre(stringvalidator.RegexMatches(
						regexp.MustCompile(`^spark\.[a-zA-Z0-9]+([\.]?[a-zA-Z0-9]+)*$`),
						"Spark properties:\n"+
							"- must starts with 'spark.'\n"+
							"- cannot contains any white spaces\n"+
							"- dot '.' is allowed but not at the start or end of the property key",
					)),
				},
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *resourceSparkEnvironmentSettings) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	pConfigData, ok := req.ProviderData.(*pconfig.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			common.ErrorResourceConfigType,
			fmt.Sprintf(common.ErrorFabricClientType, req.ProviderData),
		)

		return
	}

	r.pConfigData = pConfigData
	r.client = fabenvironment.NewClientFactoryWithClient(*pConfigData.FabricClient).NewSparkComputeClient()
	r.clientLibraries = fabenvironment.NewClientFactoryWithClient(*pConfigData.FabricClient).NewSparkLibrariesClient()
}

func (r *resourceSparkEnvironmentSettings) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("driver_cores"),
			path.MatchRoot("driver_memory"),
			path.MatchRoot("dynamic_executor_allocation"),
			path.MatchRoot("executor_cores"),
			path.MatchRoot("executor_memory"),
			path.MatchRoot("pool"),
			path.MatchRoot("runtime_version"),
			path.MatchRoot("spark_properties"),
		),
	}
}

func (r *resourceSparkEnvironmentSettings) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan resourceSparkEnvironmentSettingsModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqCreate requestUpdateSparkEnvironmentSettings

	if resp.Diagnostics.Append(reqCreate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	respCreate, err := r.client.UpdateStagingSettings(ctx, plan.WorkspaceID.ValueString(), plan.EnvironmentID.ValueString(), reqCreate.UpdateEnvironmentSparkComputeRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.set(ctx, respCreate.SparkCompute)...); resp.Diagnostics.HasError() {
		return
	}

	plan.ID = plan.EnvironmentID

	if resp.Diagnostics.Append(r.publish(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSparkEnvironmentSettings) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceSparkEnvironmentSettingsModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Read(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	diags = r.get(ctx, &state)
	if utils.IsErrNotFound(state.WorkspaceID.ValueString(), &diags, fabcore.ErrCommon.EntityNotFound) {
		resp.State.RemoveResource(ctx)

		resp.Diagnostics.Append(diags...)

		return
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	state.ID = state.EnvironmentID

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSparkEnvironmentSettings) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	var plan resourceSparkEnvironmentSettingsModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateSparkEnvironmentSettings

	if resp.Diagnostics.Append(reqUpdate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	respUpdate, err := r.client.UpdateStagingSettings(ctx, plan.WorkspaceID.ValueString(), plan.EnvironmentID.ValueString(), reqUpdate.UpdateEnvironmentSparkComputeRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.set(ctx, respUpdate.SparkCompute)...); resp.Diagnostics.HasError() {
		return
	}

	plan.ID = plan.EnvironmentID

	if resp.Diagnostics.Append(r.publish(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSparkEnvironmentSettings) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	resp.Diagnostics.AddWarning(
		"delete operation not supported",
		fmt.Sprintf("Resource %s does not support deletion. It will be removed from Terraform state, but no action will be taken in the Fabric. All current settings will remain.", SparkEnvironmentSettingsName),
	)

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceSparkEnvironmentSettings) get(ctx context.Context, model *resourceSparkEnvironmentSettingsModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s for Workspace ID: %s", SparkEnvironmentSettingsName, model.WorkspaceID.ValueString()))

	var respEntity fabenvironment.SparkCompute

	if model.PublicationStatus.ValueString() == SparkEnvironmentPublicationStatusPublished {
		respGet, err := r.client.GetPublishedSettings(ctx, model.WorkspaceID.ValueString(), model.EnvironmentID.ValueString(), nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
			return diags
		}

		respEntity = respGet.SparkCompute
	} else {
		respGet, err := r.client.GetStagingSettings(ctx, model.WorkspaceID.ValueString(), model.EnvironmentID.ValueString(), nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
			return diags
		}

		respEntity = respGet.SparkCompute
	}

	return model.set(ctx, respEntity)
}

func (r *resourceSparkEnvironmentSettings) publish(ctx context.Context, model resourceSparkEnvironmentSettingsModel) diag.Diagnostics {
	if model.PublicationStatus.ValueString() == SparkEnvironmentPublicationStatusPublished {
		for {
			respPublish, err := r.clientLibraries.PublishEnvironment(ctx, model.WorkspaceID.ValueString(), model.EnvironmentID.ValueString(), nil)
			if diags := utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil); diags.HasError() {
				return diags
			}

			if respPublish.PublishDetails == nil || respPublish.PublishDetails.State == nil {
				tflog.Info(ctx, "Environment publishing not done, waiting 30 seconds before retrying")
				time.Sleep(30 * time.Second) // lintignore:R018

				continue
			}

			switch strings.ToLower((string)(*respPublish.PublishDetails.State)) {
			case strings.ToLower((string)(fabenvironment.PublishStateFailed)):
				var diags diag.Diagnostics

				diags.AddError(
					"publishing failed",
					"Environment publishing failed")

				return diags

			case strings.ToLower((string)(fabenvironment.PublishStateCancelled)):
				var diags diag.Diagnostics

				diags.AddError(
					"publishing cancelled",
					"Environment publishing cancelled")

				return diags

			case strings.ToLower((string)(fabenvironment.PublishStateSuccess)):
				return nil

			default:
				tflog.Info(ctx, "Environment provisioning in progress, waiting 30 seconds before retrying")
				time.Sleep(30 * time.Second) // lintignore:R018
			}
		}
	}

	return nil
}
