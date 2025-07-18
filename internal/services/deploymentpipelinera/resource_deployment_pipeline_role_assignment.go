// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package deploymentpipelinera

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.ResourceWithConfigure   = (*resourceDeploymentPipelineRoleAssignment)(nil)
	_ resource.ResourceWithImportState = (*resourceDeploymentPipelineRoleAssignment)(nil)
)

type resourceDeploymentPipelineRoleAssignment struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.DeploymentPipelinesClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewResourceDeploymentPipelineRoleAssignment() resource.Resource {
	return &resourceDeploymentPipelineRoleAssignment{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceDeploymentPipelineRoleAssignment) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceDeploymentPipelineRoleAssignment) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema(false).GetResource(ctx)
}

func (r *resourceDeploymentPipelineRoleAssignment) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceDeploymentPipelineRoleAssignment) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan resourceDeploymentPipelineRoleAssignmentModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqCreate requestCreateDeploymentPipelineRoleAssignment

	reqCreate.set(ctx, plan)

	respCreate, err := r.client.AddDeploymentPipelineRoleAssignment(ctx, plan.DeploymentPipelineID.ValueString(), reqCreate.AddDeploymentPipelineRoleAssignmentRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.set(ctx, plan.DeploymentPipelineID.ValueString(), respCreate.DeploymentPipelineRoleAssignment)...); resp.Diagnostics.HasError() {
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

func (r *resourceDeploymentPipelineRoleAssignment) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceDeploymentPipelineRoleAssignmentModel

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

func (r *resourceDeploymentPipelineRoleAssignment) Update(ctx context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	// in real world, this should not reach here
	resp.Diagnostics.AddError(
		common.ErrorUpdateHeader,
		"Update is not supported. Requires delete and recreate.",
	)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceDeploymentPipelineRoleAssignment) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	var state resourceDeploymentPipelineRoleAssignmentModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.DeleteDeploymentPipelineRoleAssignment(ctx, state.DeploymentPipelineID.ValueString(), state.ID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceDeploymentPipelineRoleAssignment) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "IMPORT", map[string]any{
		"id": req.ID,
	})

	deploymentPipelineID, deploymentPipelineRoleAssignmentID, found := strings.Cut(req.ID, "/")
	if !found {
		resp.Diagnostics.AddError(
			common.ErrorImportIdentifierHeader,
			fmt.Sprintf(common.ErrorImportIdentifierDetails, "DeploymentPipelineID/DeploymentPipelineRoleAssignmentID"),
		)

		return
	}

	uuidDeploymentPipelineID, diags := customtypes.NewUUIDValueMust(deploymentPipelineID)
	resp.Diagnostics.Append(diags...)

	uuidDeploymentPipelineRoleAssignmentID, diags := customtypes.NewUUIDValueMust(deploymentPipelineRoleAssignmentID)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var timeout timeouts.Value
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeout)...); resp.Diagnostics.HasError() {
		return
	}

	state := resourceDeploymentPipelineRoleAssignmentModel{
		baseDeploymentPipelineRoleAssignmentModel: baseDeploymentPipelineRoleAssignmentModel{
			ID:                   uuidDeploymentPipelineRoleAssignmentID,
			DeploymentPipelineID: uuidDeploymentPipelineID,
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

func (r *resourceDeploymentPipelineRoleAssignment) get(ctx context.Context, model *resourceDeploymentPipelineRoleAssignmentModel) diag.Diagnostics {
	respGetInfo, err := r.client.ListDeploymentPipelineRoleAssignments(ctx, model.DeploymentPipelineID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		return diags
	}

	for _, entity := range respGetInfo {
		if *entity.ID == model.ID.ValueString() {
			if diags := model.set(ctx, model.DeploymentPipelineID.ValueString(), entity); diags.HasError() {
				return diags
			}

			break
		}
	}

	return nil
}
