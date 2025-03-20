// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package spark

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
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
	fabspark "github.com/microsoft/fabric-sdk-go/fabric/spark"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure        = (*resourceSparkWorkspaceSettings)(nil)
	_ resource.ResourceWithConfigValidators = (*resourceSparkWorkspaceSettings)(nil)
)

type resourceSparkWorkspaceSettings struct {
	pConfigData *pconfig.ProviderData
	client      *fabspark.WorkspaceSettingsClient
}

func NewResourceSparkWorkspaceSettings() resource.Resource {
	return &resourceSparkWorkspaceSettings{}
}

func (r *resourceSparkWorkspaceSettings) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + SparkWorkspaceSettingsTFName
}

func (r *resourceSparkWorkspaceSettings) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage a Fabric " + SparkWorkspaceSettingsName + ".\n\n" +
			"See [" + SparkWorkspaceSettingsName + "](" + SparkWorkspaceSettingsDocsURL + ") for more information.\n\n" +
			SparkWorkspaceSettingsDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The Workspace ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"automatic_log": schema.SingleNestedAttribute{
				MarkdownDescription: "Automatic Log properties.",
				Optional:            true,
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[automaticLogPropertiesModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "The status of the automatic log. Possible values: `false` - Disabled, `true` - Enabled.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"environment": schema.SingleNestedAttribute{
				MarkdownDescription: "Environment properties.",
				Optional:            true,
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[environmentPropertiesModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the default environment. Empty indicated there is no workspace default environment.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"runtime_version": schema.StringAttribute{
						MarkdownDescription: "[Runtime](https://review.learn.microsoft.com/fabric/data-engineering/runtime) version. Accepted values: " + utils.ConvertStringSlicesToString(SparkRuntimeVersionValues, true, false) + ".",
						Description:         "Runtime version. Accepted values: " + utils.ConvertStringSlicesToString(SparkRuntimeVersionValues, true, false) + ".",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.OneOf(SparkRuntimeVersionValues...),
						},
					},
				},
			},
			"high_concurrency": schema.SingleNestedAttribute{
				MarkdownDescription: "High Concurrency properties.",
				Optional:            true,
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[highConcurrencyPropertiesModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"notebook_interactive_run_enabled": schema.BoolAttribute{
						MarkdownDescription: "The status of the high concurrency for notebook interactive run. `false` - Disabled, `true` - Enabled.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"notebook_pipeline_run_enabled": schema.BoolAttribute{
						MarkdownDescription: "The status of the high concurrency for notebook pipeline run. `false` - Disabled, `true` - Enabled.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"job": schema.SingleNestedAttribute{
				MarkdownDescription: "Jobs properties.",
				Optional:            true,
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[jobPropertiesModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"conservative_job_admission_enabled": schema.BoolAttribute{
						MarkdownDescription: "Reserve maximum cores for active Spark jobs. When this setting is enabled, your Fabric capacity reserves the maximum number of cores needed for active Spark jobs, ensuring job reliability by making sure that cores are available if a job scales up. When this setting is disabled, jobs are started based on the minimum number of cores needed, letting more jobs run at the same time. `false` - Disabled, `true` - Enabled.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"session_timeout_in_minutes": schema.Int32Attribute{
						MarkdownDescription: "Time to terminate inactive Spark sessions. The maximum is 14 days (20160 minutes).",
						Optional:            true,
						Computed:            true,
						Validators: []validator.Int32{
							int32validator.AtMost(20160),
						},
						PlanModifiers: []planmodifier.Int32{
							int32planmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"pool": schema.SingleNestedAttribute{
				MarkdownDescription: "Pool properties.",
				Optional:            true,
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[poolPropertiesModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"customize_compute_enabled": schema.BoolAttribute{
						MarkdownDescription: "Customize compute configurations for items. `false` - Disabled, `true` - Enabled.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"default_pool": schema.SingleNestedAttribute{
						MarkdownDescription: "Default pool for workspace.",
						Optional:            true,
						Computed:            true,
						CustomType:          supertypes.NewSingleNestedObjectTypeOf[defaultPoolPropertiesModel](ctx),
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								MarkdownDescription: "The Pool ID. `00000000-0000-0000-0000-000000000000` means use the starter pool.",
								Optional:            true,
								Computed:            true,
								CustomType:          customtypes.UUIDType{},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
								Validators: []validator.String{
									stringvalidator.ConflictsWith(
										path.MatchRelative().AtParent().AtName("name"),
										path.MatchRelative().AtParent().AtName("type"),
									),
								},
							},
							"name": schema.StringAttribute{
								MarkdownDescription: "The Pool name. It should be a valid custom pool name. `Starter Pool` means use the starter pool.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
								Validators: []validator.String{
									stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("type")),
									stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("id")),
								},
							},
							"type": schema.StringAttribute{
								MarkdownDescription: "The Pool type. Accepted values: " + utils.ConvertStringSlicesToString(fabspark.PossibleCustomPoolTypeValues(), true, true) + ".",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
								Validators: []validator.String{
									stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("name")),
									stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("id")),
									stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabspark.PossibleCustomPoolTypeValues(), false)...),
								},
							},
						},
					},
					"starter_pool": schema.SingleNestedAttribute{
						MarkdownDescription: "Starter pool for workspace. For more information about configuring starter pool, see [configuring starter pool](https://review.learn.microsoft.com/fabric/data-engineering/configure-starter-pools).",
						Description:         "Starter pool for workspace.",
						Optional:            true,
						Computed:            true,
						CustomType:          supertypes.NewSingleNestedObjectTypeOf[starterPoolPropertiesModel](ctx),
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"max_node_count": schema.Int32Attribute{
								MarkdownDescription: "The maximum node count.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Int32{
									int32planmodifier.UseStateForUnknown(),
								},
							},
							"max_executors": schema.Int32Attribute{
								MarkdownDescription: "The maximum executors count.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Int32{
									int32planmodifier.UseStateForUnknown(),
								},
							},
						},
					},
				},
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *resourceSparkWorkspaceSettings) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = fabspark.NewClientFactoryWithClient(*pConfigData.FabricClient).NewWorkspaceSettingsClient()
}

func (r *resourceSparkWorkspaceSettings) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("automatic_log"),
			path.MatchRoot("environment"),
			path.MatchRoot("high_concurrency"),
			path.MatchRoot("job"),
			path.MatchRoot("pool"),
		),
	}
}

func (r *resourceSparkWorkspaceSettings) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan resourceSparkWorkspaceSettingsModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqCreate requestUpdateSparkWorkspaceSettings

	if resp.Diagnostics.Append(reqCreate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	respCreate, err := r.client.UpdateSparkSettings(ctx, plan.WorkspaceID.ValueString(), reqCreate.UpdateWorkspaceSparkSettingsRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.set(ctx, respCreate.WorkspaceSparkSettings)...); resp.Diagnostics.HasError() {
		return
	}

	plan.ID = plan.WorkspaceID

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSparkWorkspaceSettings) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceSparkWorkspaceSettingsModel

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

	state.ID = state.WorkspaceID

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSparkWorkspaceSettings) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	var plan resourceSparkWorkspaceSettingsModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateSparkWorkspaceSettings

	if resp.Diagnostics.Append(reqUpdate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	respUpdate, err := r.client.UpdateSparkSettings(ctx, plan.WorkspaceID.ValueString(), reqUpdate.UpdateWorkspaceSparkSettingsRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.set(ctx, respUpdate.WorkspaceSparkSettings)...); resp.Diagnostics.HasError() {
		return
	}

	plan.ID = plan.WorkspaceID

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSparkWorkspaceSettings) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	resp.Diagnostics.AddWarning(
		"delete operation not supported",
		fmt.Sprintf("Resource %s does not support deletion. It will be removed from Terraform state, but no action will be taken in the Fabric. All current settings will remain.", SparkWorkspaceSettingsName),
	)

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceSparkWorkspaceSettings) get(ctx context.Context, model *resourceSparkWorkspaceSettingsModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s for Workspace ID: %s", SparkWorkspaceSettingsName, model.WorkspaceID.ValueString()))

	respGet, err := r.client.GetSparkSettings(ctx, model.WorkspaceID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	return model.set(ctx, respGet.WorkspaceSparkSettings)
}
