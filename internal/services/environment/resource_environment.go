// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package environment

import (
	"context"
	"fmt"
	"strings"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure   = (*resourceEnvironment)(nil)
	_ resource.ResourceWithImportState = (*resourceEnvironment)(nil)
)

type resourceEnvironment struct {
	pConfigData *pconfig.ProviderData
	client      *fabenvironment.ItemsClient
}

func NewResourceEnvironment() resource.Resource {
	return &resourceEnvironment{}
}

func (r *resourceEnvironment) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ItemTFName
}

func (r *resourceEnvironment) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	markdownDescription := "This resource manages a Fabric " + ItemName + ".\n\n" +
		"See [" + ItemName + "](" + ItemDocsURL + ") for more information.\n\n" +
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

	resp.Schema = fabricitem.GetResourceFabricItemPropertiesSchema(ctx, ItemName, markdownDescription, 123, 256, true, properties)
}

func (r *resourceEnvironment) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = fabenvironment.NewClientFactoryWithClient(*pConfigData.FabricClient).NewItemsClient()
}

func (r *resourceEnvironment) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "CREATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
	})

	var plan resourceEnvironmentModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqCreate requestCreateEnvironment

	reqCreate.set(plan)

	respCreate, err := r.client.CreateEnvironment(ctx, plan.WorkspaceID.ValueString(), reqCreate.CreateEnvironmentRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	plan.set(respCreate.Environment)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	if resp.Diagnostics.Append(r.get(ctx, &plan)...); resp.Diagnostics.HasError() {
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

func (r *resourceEnvironment) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"state": req.State,
	})

	var state resourceEnvironmentModel

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
	if utils.IsErrNotFound(state.ID.ValueString(), &diags, fabcore.ErrCommon.EntityNotFound) {
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

func (r *resourceEnvironment) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "UPDATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
		"state":  req.State,
	})

	var plan resourceEnvironmentModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateEnvironment

	reqUpdate.set(plan)

	respUpdate, err := r.client.UpdateEnvironment(ctx, plan.WorkspaceID.ValueString(), plan.ID.ValueString(), reqUpdate.UpdateEnvironmentRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	plan.set(respUpdate.Environment)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	if resp.Diagnostics.Append(r.get(ctx, &plan)...); resp.Diagnostics.HasError() {
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

func (r *resourceEnvironment) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "DELETE", map[string]any{
		"state": req.State,
	})

	var state resourceEnvironmentModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.DeleteEnvironment(ctx, state.WorkspaceID.ValueString(), state.ID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceEnvironment) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "IMPORT", map[string]any{
		"id": req.ID,
	})

	workspaceID, environmentID, found := strings.Cut(req.ID, "/")
	if !found {
		resp.Diagnostics.AddError(
			common.ErrorImportIdentifierHeader,
			fmt.Sprintf(common.ErrorImportIdentifierDetails, "WorkspaceID/EnvironmentID"),
		)

		return
	}

	uuidWorkspaceID, diags := customtypes.NewUUIDValueMust(workspaceID)
	resp.Diagnostics.Append(diags...)

	uuidID, diags := customtypes.NewUUIDValueMust(environmentID)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var timeout timeouts.Value
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeout)...); resp.Diagnostics.HasError() {
		return
	}

	state := resourceEnvironmentModel{}
	state.ID = uuidID
	state.WorkspaceID = uuidWorkspaceID
	state.Timeouts = timeout

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

func (r *resourceEnvironment) get(ctx context.Context, model *resourceEnvironmentModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET", map[string]any{
		"workspace_id": model.WorkspaceID.ValueString(),
		"id":           model.ID.ValueString(),
	})

	respGet, err := r.client.GetEnvironment(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		return diags
	}

	model.set(respGet.Environment)

	return model.setProperties(ctx, respGet.Environment)
}
