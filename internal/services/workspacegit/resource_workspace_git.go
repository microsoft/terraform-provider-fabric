// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspacegit

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ resource.ResourceWithConfigure = (*resourceWorkspaceGit)(nil)

type resourceWorkspaceGit struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GitClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewResourceWorkspaceGit() resource.Resource {
	return &resourceWorkspaceGit{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceWorkspaceGit) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceWorkspaceGit) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema().GetResource(ctx)
}

func (r *resourceWorkspaceGit) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewGitClient()

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(r.TypeInfo.Name, r.TypeInfo.IsPreview, r.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceWorkspaceGit) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan resourceWorkspaceGitModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Connect.
	var reqGitConnect requestGitConnect

	if resp.Diagnostics.Append(reqGitConnect.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Connect(ctx, plan.WorkspaceID.ValueString(), reqGitConnect.GitConnectRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	// Initialize.
	var reqGitInitialize requestGitInitialize

	reqGitInitialize.set(plan)

	gitInitResp, err := r.client.InitializeConnection(ctx, plan.WorkspaceID.ValueString(), &fabcore.GitClientBeginInitializeConnectionOptions{
		GitInitializeConnectionRequest: &reqGitInitialize.InitializeGitConnectionRequest,
	})
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	// Git commit.
	switch *gitInitResp.RequiredAction {
	case fabcore.RequiredActionCommitToGit: // Commit to Git.
		var reqGitCommitTo requestGitCommitTo

		reqGitCommitTo.set(gitInitResp.WorkspaceHead)

		_, err = r.client.CommitToGit(ctx, plan.WorkspaceID.ValueString(), reqGitCommitTo.CommitToGitRequest, nil)

	case fabcore.RequiredActionUpdateFromGit: // Update from Git.
		var reqGitUpdateFrom requestGitUpdateFrom

		if resp.Diagnostics.Append(reqGitUpdateFrom.set(ctx, plan, gitInitResp.RemoteCommitHash, plan.InitializationStrategy.ValueStringPointer())...); resp.Diagnostics.HasError() {
			return
		}

		_, err = r.client.UpdateFromGit(ctx, plan.WorkspaceID.ValueString(), reqGitUpdateFrom.UpdateFromGitRequest, nil)
	case fabcore.RequiredActionNone:
		// Do nothing.
	default:
		resp.Diagnostics.AddError(
			common.ErrorCreateHeader,
			fmt.Sprintf("Unsupported required git action '%s'.", *gitInitResp.RequiredAction),
		)
	}

	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(r.get(ctx, &plan)...); resp.Diagnostics.HasError() {
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

func (r *resourceWorkspaceGit) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceWorkspaceGitModel

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
	if utils.IsErrNotFound(state.ID.ValueString(), &diags, fabcore.ErrGit.GitProviderResourceNotFound) {
		resp.State.RemoveResource(ctx)

		resp.Diagnostics.Append(diags...)

		return
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	if !state.GitConnectionState.IsNull() && !state.GitConnectionState.IsUnknown() && state.GitConnectionState.ValueString() != (string)(fabcore.GitConnectionStateConnectedAndInitialized) {
		resp.Diagnostics.AddWarning(
			"Unexpected Git connection state",
			fmt.Sprintf("Git connection state is '%s'.\nIt may have been deleted outside of Terraform. Removing object from state.", state.GitConnectionState.ValueString()),
		)

		resp.State.RemoveResource(ctx)

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

func (r *resourceWorkspaceGit) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	var plan resourceWorkspaceGitModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateGitCredentials

	if resp.Diagnostics.Append(reqUpdate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	respUpdate, err := r.client.UpdateMyGitCredentials(ctx, plan.WorkspaceID.ValueString(), reqUpdate.UpdateGitCredentialsRequestClassification, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.setCredentials(ctx, respUpdate.GitCredentialsConfigurationResponseClassification)...); resp.Diagnostics.HasError() {
		return
	}

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

func (r *resourceWorkspaceGit) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	var state resourceWorkspaceGitModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.Disconnect(ctx, state.WorkspaceID.ValueString(), nil)

	diags = utils.GetDiagsFromError(ctx, err, utils.OperationDelete, fabcore.ErrGit.WorkspaceNotConnectedToGit)
	if diags.HasError() && !utils.IsErr(diags, fabcore.ErrGit.WorkspaceNotConnectedToGit) {
		resp.Diagnostics.Append(diags...)

		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceWorkspaceGit) get(ctx context.Context, model *resourceWorkspaceGitModel) diag.Diagnostics {
	respGet, err := r.client.GetConnection(ctx, model.WorkspaceID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrGit.GitProviderResourceNotFound); diags.HasError() {
		return diags
	}

	if diags := model.set(ctx, respGet.GitConnection); diags.HasError() {
		return diags
	}

	respGetCredentials, err := r.client.GetMyGitCredentials(ctx, model.WorkspaceID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	if diags := model.setCredentials(ctx, respGetCredentials.GitCredentialsConfigurationResponseClassification); diags.HasError() {
		return diags
	}

	return nil
}
