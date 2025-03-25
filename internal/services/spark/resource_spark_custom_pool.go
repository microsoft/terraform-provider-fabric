// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package spark

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabspark "github.com/microsoft/fabric-sdk-go/fabric/spark"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	superint32validator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/int32validator"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure   = (*resourceSparkCustomPool)(nil)
	_ resource.ResourceWithImportState = (*resourceSparkCustomPool)(nil)
)

type resourceSparkCustomPool struct {
	pConfigData *pconfig.ProviderData
	client      *fabspark.CustomPoolsClient
	Name        string
	IsPreview   bool
}

func NewResourceSparkCustomPool() resource.Resource {
	return &resourceSparkCustomPool{
		Name:      SparkCustomPoolName,
		IsPreview: SparkCustomPoolPreview,
	}
}

func (r *resourceSparkCustomPool) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + SparkCustomPoolTFName
}

func (r *resourceSparkCustomPool) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage a Fabric " + SparkCustomPoolName + ".\n\n" +
			"See [" + SparkCustomPoolName + "](" + SparkCustomPoolDocsURL + ") for more information.\n\n" +
			SparkCustomPoolDocsSPNSupport,
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
				MarkdownDescription: "The " + SparkCustomPoolName + " ID.",
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The " + SparkCustomPoolName + " name.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(64),
					stringvalidator.NoneOfCaseInsensitive("Starter Pool"),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9-_ ]+$`), "The name must contain only letters, numbers, dashes, underscores and spaces."),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The " + SparkCustomPoolName + " type. Accepted values: " + utils.ConvertStringSlicesToString(
					utils.RemoveSliceByValue(fabspark.PossibleCustomPoolTypeValues(), fabspark.CustomPoolTypeCapacity),
					true,
					true,
				) + ".",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(utils.RemoveSliceByValue(fabspark.PossibleCustomPoolTypeValues(), fabspark.CustomPoolTypeCapacity), false)...),
				},
			},
			"node_family": schema.StringAttribute{
				MarkdownDescription: "The Node family. Accepted values: " + utils.ConvertStringSlicesToString(fabspark.PossibleNodeFamilyValues(), true, true) + ".",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabspark.PossibleNodeFamilyValues(), false)...),
				},
			},
			"node_size": schema.StringAttribute{
				MarkdownDescription: "The Node size. Accepted values: " + utils.ConvertStringSlicesToString(fabspark.PossibleNodeSizeValues(), true, true) + ".",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabspark.PossibleNodeSizeValues(), false)...),
				},
			},
			"auto_scale": schema.SingleNestedAttribute{
				MarkdownDescription: "Auto-scale properties.",
				Required:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[sparkCustomPoolAutoScaleModel](ctx),
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "The status of the auto scale. Accepted values: `false` - Disabled, `true` - Enabled.",
						Required:            true,
					},
					"min_node_count": schema.Int32Attribute{
						MarkdownDescription: "The minimum node count.",
						Required:            true,
					},
					"max_node_count": schema.Int32Attribute{
						MarkdownDescription: "The maximum node count.",
						Required:            true,
					},
				},
			},
			"dynamic_executor_allocation": schema.SingleNestedAttribute{
				MarkdownDescription: "Dynamic Executor Allocation properties.",
				Required:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[sparkCustomPoolDynamicExecutorAllocationModel](ctx),
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "The status of the dynamic executor allocation. Accepted values: `false` - Disabled, `true` - Enabled.",
						Required:            true,
					},
					"min_executors": schema.Int32Attribute{
						MarkdownDescription: "The minimum executors.",
						Computed:            true,
						Optional:            true,
						Validators: []validator.Int32{
							superint32validator.NullIfAttributeIsOneOf(
								path.MatchRoot("dynamic_executor_allocation").AtName("enabled"),
								[]attr.Value{types.BoolValue(false)},
							),
							superint32validator.RequireIfAttributeIsOneOf(
								path.MatchRoot("dynamic_executor_allocation").AtName("enabled"),
								[]attr.Value{types.BoolValue(true)},
							),
						},
					},
					"max_executors": schema.Int32Attribute{
						MarkdownDescription: "The maximum executors.",
						Computed:            true,
						Optional:            true,
						Validators: []validator.Int32{
							superint32validator.NullIfAttributeIsOneOf(
								path.MatchRoot("dynamic_executor_allocation").AtName("enabled"),
								[]attr.Value{types.BoolValue(false)},
							),
							superint32validator.RequireIfAttributeIsOneOf(
								path.MatchRoot("dynamic_executor_allocation").AtName("enabled"),
								[]attr.Value{types.BoolValue(true)},
							),
						},
					},
				},
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *resourceSparkCustomPool) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = fabspark.NewClientFactoryWithClient(*pConfigData.FabricClient).NewCustomPoolsClient()

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(r.Name, r.IsPreview, r.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSparkCustomPool) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan resourceSparkCustomPoolModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqCreate requestCreateSparkCustomPool

	if resp.Diagnostics.Append(reqCreate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	respCreate, err := r.client.CreateWorkspaceCustomPool(ctx, plan.WorkspaceID.ValueString(), reqCreate.CreateCustomPoolRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.set(ctx, respCreate.CustomPool)...); resp.Diagnostics.HasError() {
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

func (r *resourceSparkCustomPool) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceSparkCustomPoolModel

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
	if utils.IsErrNotFound(state.ID.ValueString(), &diags, fabcore.ErrSpark.SparkSettingsManagementUserError) {
		resp.State.RemoveResource(ctx)

		resp.Diagnostics.Append(diags...)

		return
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSparkCustomPool) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	var plan resourceSparkCustomPoolModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateSparkCustomPool

	if resp.Diagnostics.Append(reqUpdate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	respUpdate, err := r.client.UpdateWorkspaceCustomPool(ctx, plan.WorkspaceID.ValueString(), plan.ID.ValueString(), reqUpdate.UpdateCustomPoolRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.set(ctx, respUpdate.CustomPool)...); resp.Diagnostics.HasError() {
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

func (r *resourceSparkCustomPool) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	var state resourceSparkCustomPoolModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.DeleteWorkspaceCustomPool(ctx, state.WorkspaceID.ValueString(), state.ID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceSparkCustomPool) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "start",
	})

	workspaceID, poolID, found := strings.Cut(req.ID, "/")
	if !found {
		resp.Diagnostics.AddError(
			common.ErrorImportIdentifierHeader,
			fmt.Sprintf(common.ErrorImportIdentifierDetails, "WorkspaceID/PoolID"),
		)

		return
	}

	uuidWorkspaceID, diags := customtypes.NewUUIDValueMust(workspaceID)
	resp.Diagnostics.Append(diags...)

	uuidID, diags := customtypes.NewUUIDValueMust(poolID)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var timeout timeouts.Value
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeout)...); resp.Diagnostics.HasError() {
		return
	}

	state := resourceSparkCustomPoolModel{
		baseSparkCustomPoolModel: baseSparkCustomPoolModel{
			ID:          uuidID,
			WorkspaceID: uuidWorkspaceID,
		},
		Timeouts: timeout,
	}

	if resp.Diagnostics.Append(r.get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSparkCustomPool) get(ctx context.Context, model *resourceSparkCustomPoolModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting "+SparkCustomPoolName)

	respGet, err := r.client.GetWorkspaceCustomPool(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrSpark.SparkSettingsManagementUserError); diags.HasError() {
		return diags
	}

	return model.set(ctx, respGet.CustomPool)
}
