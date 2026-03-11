// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package deploymentpipeline

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure   = (*resourceDeploymentPipeline)(nil)
	_ resource.ResourceWithImportState = (*resourceDeploymentPipeline)(nil)
)

type resourceDeploymentPipeline struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.DeploymentPipelinesClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewResourceDeploymentPipeline() resource.Resource {
	return &resourceDeploymentPipeline{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceDeploymentPipeline) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceDeploymentPipeline) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema(false).GetResource(ctx)
}

func (r *resourceDeploymentPipeline) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(r.TypeInfo.Name, r.TypeInfo.IsPreview, r.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}

	r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewDeploymentPipelinesClient()
}

func (r *resourceDeploymentPipeline) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan, state resourceDeploymentPipelineModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	state.Timeouts = plan.Timeouts

	var reqCreate requestCreateDeploymentPipeline

	diags = reqCreate.set(ctx, plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	respCreate, err := r.client.CreateDeploymentPipeline(ctx, reqCreate.CreateDeploymentPipelineRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	diags = state.set(ctx, respCreate.DeploymentPipelineExtendedInfo)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	planStages, diags := plan.Stages.Get(ctx)

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	for i, stage := range planStages {
		if stage.WorkspaceID.ValueString() == "" {
			continue
		}

		stage.AssignWorkspace(ctx, r.client, &state, &resp.Diagnostics, i)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceDeploymentPipeline) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceDeploymentPipelineModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Read(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	diags = r.getByID(ctx, state.ID.ValueString(), &state)
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

func (r *resourceDeploymentPipeline) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	var plan, state resourceDeploymentPipelineModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateDeploymentPipeline

	reqUpdate.set(plan)

	respUpdate, err := r.client.UpdateDeploymentPipeline(ctx, plan.ID.ValueString(), reqUpdate.UpdateDeploymentPipelineRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	planStages, diags := plan.Stages.Get(ctx)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	diags = r.getByID(ctx, plan.ID.ValueString(), &state)
	if utils.IsErrNotFound(state.ID.ValueString(), &diags, fabcore.ErrCommon.EntityNotFound) {
		resp.State.RemoveResource(ctx)

		resp.Diagnostics.Append(diags...)

		return
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	stateStages, diags := state.Stages.Get(ctx)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	for i, stage := range planStages {
		if stage.WorkspaceID.ValueString() == "" && stateStages[i].WorkspaceID.ValueString() != "" {
			stage.UnassignWorkspace(ctx, r.client, &state, &resp.Diagnostics, i)
		} else if stage.WorkspaceID.ValueString() != "" && stateStages[i].WorkspaceID.ValueString() == "" {
			stage.AssignWorkspace(ctx, r.client, &state, &resp.Diagnostics, i)
		}

		if stage.DisplayName.ValueString() != stateStages[i].DisplayName.ValueString() ||
			stage.IsPublic.ValueBool() != stateStages[i].IsPublic.ValueBool() ||
			stage.Description.ValueString() != stateStages[i].Description.ValueString() {
			stage.UpdateStage(ctx, r.client, &state, &resp.Diagnostics, i)
		}
	}

	if resp.Diagnostics.Append(plan.set(ctx, respUpdate.DeploymentPipelineExtendedInfo)...); resp.Diagnostics.HasError() {
		return
	}

	stages, diags := state.Stages.Get(ctx)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.setStages(ctx, stages)...); resp.Diagnostics.HasError() {
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

func (r *resourceDeploymentPipeline) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	var state resourceDeploymentPipelineModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.DeleteDeploymentPipeline(ctx, state.ID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceDeploymentPipeline) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "IMPORT", map[string]any{
		"id": req.ID,
	})

	_, diags := customtypes.NewUUIDValueMust(req.ID)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceDeploymentPipeline) getByID(ctx context.Context, deploymentPipelineID string, model *resourceDeploymentPipelineModel) diag.Diagnostics {
	respGet, err := r.client.GetDeploymentPipeline(ctx, deploymentPipelineID, nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		return diags
	}

	if diags := model.set(ctx, respGet.DeploymentPipelineExtendedInfo); diags.HasError() {
		return diags
	}

	return nil
}

func (stage *baseDeploymentPipelineStageModel) AssignWorkspace(
	ctx context.Context,
	client *fabcore.DeploymentPipelinesClient,
	state *resourceDeploymentPipelineModel,
	respDiags *diag.Diagnostics,
	order int,
) {
	var req requestAssignStageToWorkspace
	req.set(*stage)

	stateStages, diags := state.Stages.Get(ctx)

	*respDiags = append(*respDiags, diags...)
	if respDiags.HasError() {
		return
	}

	tflog.Debug(ctx, "ASSIGN WORKSPACE", map[string]any{
		"action": "start",
		"id":     stateStages[order].ID.ValueString(),
	})

	_, err := client.AssignWorkspaceToStage(
		ctx,
		state.ID.ValueString(),
		stateStages[order].ID.ValueString(),
		req.DeploymentPipelineAssignWorkspaceRequest,
		nil,
	)

	*respDiags = append(*respDiags, utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...)
	if respDiags.HasError() {
		return
	}

	stateStages[order].WorkspaceID = stage.WorkspaceID

	tflog.Debug(ctx, "ASSIGN WORKSPACE", map[string]any{
		"action": "end",
		"id":     stateStages[order].ID.ValueString(),
	})

	*respDiags = append(*respDiags, state.setStages(ctx, stateStages)...)
	if respDiags.HasError() {
		return
	}
}

func (stage *baseDeploymentPipelineStageModel) UnassignWorkspace(
	ctx context.Context,
	client *fabcore.DeploymentPipelinesClient,
	state *resourceDeploymentPipelineModel,
	respDiags *diag.Diagnostics,
	order int,
) {
	stateStages, diags := state.Stages.Get(ctx)

	*respDiags = append(*respDiags, diags...)
	if respDiags.HasError() {
		return
	}

	tflog.Debug(ctx, "UNASSIGN WORKSPACE", map[string]any{
		"action": "start",
		"id":     stateStages[order].ID.ValueString(),
	})

	_, err := client.UnassignWorkspaceFromStage(
		ctx,
		state.ID.ValueString(),
		stateStages[order].ID.ValueString(),
		nil,
	)

	*respDiags = append(*respDiags, utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...)
	if respDiags.HasError() {
		return
	}

	stateStages[order].WorkspaceID = supertypes.NewStringNull().StringValue

	tflog.Debug(ctx, "UNASSIGN WORKSPACE", map[string]any{
		"action": "end",
		"id":     stateStages[order].ID.ValueString(),
	})

	*respDiags = append(*respDiags, state.setStages(ctx, stateStages)...)
	if respDiags.HasError() {
		return
	}
}

func (stage *baseDeploymentPipelineStageModel) UpdateStage(
	ctx context.Context,
	client *fabcore.DeploymentPipelinesClient,
	state *resourceDeploymentPipelineModel,
	respDiags *diag.Diagnostics,
	order int,
) {
	var req requestUpdateDeploymentPipelineStage
	req.set(*stage)

	stateStages, diags := state.Stages.Get(ctx)

	*respDiags = append(*respDiags, diags...)
	if respDiags.HasError() {
		return
	}

	initialStage := stateStages[order]

	respUpdate, err := client.UpdateDeploymentPipelineStage(
		ctx,
		state.ID.ValueString(),
		stage.ID.ValueString(),
		req.DeploymentPipelineStageRequest,
		nil,
	)

	*respDiags = append(*respDiags, utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...)
	if respDiags.HasError() {
		return
	}

	stage.set(respUpdate.DeploymentPipelineStage)
	stage.WorkspaceID = initialStage.WorkspaceID
	stateStages[order] = stage

	*respDiags = append(*respDiags, state.setStages(ctx, stateStages)...)
	if respDiags.HasError() {
		return
	}
}
