// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceencryption

import (
	"context"
	"fmt"
	"time"

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

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure   = (*resourceWorkspaceEncryption)(nil)
	_ resource.ResourceWithImportState = (*resourceWorkspaceEncryption)(nil)
)

type resourceWorkspaceEncryption struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.WorkspacesClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewResourceWorkspaceEncryption() resource.Resource {
	return &resourceWorkspaceEncryption{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceWorkspaceEncryption) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceWorkspaceEncryption) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema().GetResource(ctx)
}

func (r *resourceWorkspaceEncryption) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewWorkspacesClient()
}

func (r *resourceWorkspaceEncryption) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan resourceWorkspaceEncryptionModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(r.assign(ctx, &plan, utils.OperationCreate)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})
}

func (r *resourceWorkspaceEncryption) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceWorkspaceEncryptionModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Read(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	diags = r.get(ctx, &state.baseWorkspaceEncryptionModel)
	if utils.IsErrNotFound(state.WorkspaceID.ValueString(), &diags, fabcore.ErrCommon.EntityNotFound) {
		resp.State.RemoveResource(ctx)

		resp.Diagnostics.Append(diags...)

		return
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// A workspace without a customer-managed key reports Disabled, which is the absence of this resource.
	if state.EncryptionStatus.ValueString() == string(fabcore.WorkspaceEncryptionStatusDisabled) {
		resp.State.RemoveResource(ctx)

		return
	}

	// Encryption can fail outside of Terraform, for example when the key is revoked. Without this warning the
	// plan would be empty and the broken workspace would go unnoticed.
	if state.EncryptionStatus.ValueString() == string(fabcore.WorkspaceEncryptionStatusFailed) {
		resp.Diagnostics.AddWarning(
			r.TypeInfo.Name+" failed",
			fmt.Sprintf(
				"%s is in the Failed state for Workspace ID: %s. Verify that the key exists, is enabled, and that the 'Fabric Platform CMK' application can wrap and unwrap it, then re-apply to retry.",
				r.TypeInfo.Name,
				state.WorkspaceID.ValueString(),
			),
		)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})
}

func (r *resourceWorkspaceEncryption) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	var plan resourceWorkspaceEncryptionModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(r.assign(ctx, &plan, utils.OperationUpdate)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})
}

func (r *resourceWorkspaceEncryption) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	var state resourceWorkspaceEncryptionModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Resetting removes the customer-managed key, after which the workspace falls back to Microsoft-managed keys.
	_, err := r.client.ResetWorkspaceEncryption(ctx, state.WorkspaceID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(r.waitForStatus(ctx, state.WorkspaceID.ValueString(), fabcore.WorkspaceEncryptionStatusDisabled, nil)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceWorkspaceEncryption) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "IMPORT", map[string]any{
		"id": req.ID,
	})

	uuidWorkspaceID, diags := customtypes.NewUUIDValueMust(req.ID)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	var timeout timeouts.Value
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeout)...); resp.Diagnostics.HasError() {
		return
	}

	state := resourceWorkspaceEncryptionModel{
		WorkspaceID: uuidWorkspaceID,
		Timeouts:    timeout,
	}

	if resp.Diagnostics.Append(r.get(ctx, &state.baseWorkspaceEncryptionModel)...); resp.Diagnostics.HasError() {
		return
	}

	if state.EncryptionStatus.ValueString() == string(fabcore.WorkspaceEncryptionStatusDisabled) {
		resp.Diagnostics.AddError(
			common.ErrorImportHeader,
			fmt.Sprintf("%s is not enabled for Workspace ID: %s", r.TypeInfo.Name, req.ID),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "end",
	})
}

func (r *resourceWorkspaceEncryption) assign(ctx context.Context, model *resourceWorkspaceEncryptionModel, operation utils.Operation) diag.Diagnostics {
	var reqAssign requestAssignWorkspaceEncryption

	reqAssign.set(*model)

	_, err := r.client.AssignWorkspaceEncryption(ctx, model.WorkspaceID.ValueString(), reqAssign.AssignWorkspaceEncryptionRequest, nil)
	if diags := utils.GetDiagsFromError(ctx, err, operation, nil); diags.HasError() {
		return diags
	}

	return r.waitForStatus(ctx, model.WorkspaceID.ValueString(), fabcore.WorkspaceEncryptionStatusActive, &model.baseWorkspaceEncryptionModel)
}

func (r *resourceWorkspaceEncryption) get(ctx context.Context, model *baseWorkspaceEncryptionModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s for Workspace ID: %s", r.TypeInfo.Name, model.WorkspaceID.ValueString()))

	respGet, err := r.client.GetWorkspaceEncryption(ctx, model.WorkspaceID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		return diags
	}

	return model.set(ctx, model.WorkspaceID.ValueString(), respGet.WorkspaceEncryptionDetail)
}

// waitForStatus polls the encryption status until it settles on want, because assign and reset are asynchronous.
func (r *resourceWorkspaceEncryption) waitForStatus(ctx context.Context, workspaceID string, want fabcore.WorkspaceEncryptionStatus, model *baseWorkspaceEncryptionModel) diag.Diagnostics {
	var diags diag.Diagnostics

	for {
		respGet, err := r.client.GetWorkspaceEncryption(ctx, workspaceID, nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
			return diags
		}

		status := encryptionStatus(respGet.WorkspaceEncryptionDetail)

		if status == fabcore.WorkspaceEncryptionStatusFailed {
			diags.AddError(
				common.ErrorGenericUnexpected,
				fmt.Sprintf(
					"%s failed for Workspace ID: %s. Verify that the key exists, is enabled, and that the 'Fabric Platform CMK' application can wrap and unwrap it.",
					r.TypeInfo.Name,
					workspaceID,
				),
			)

			return diags
		}

		if status == want {
			if model != nil {
				if diags := model.set(ctx, workspaceID, respGet.WorkspaceEncryptionDetail); diags.HasError() {
					return diags
				}
			}

			return diags
		}

		tflog.Trace(ctx, fmt.Sprintf("waiting for %s of Workspace ID: %s to become %s, current status: %s", r.TypeInfo.Name, workspaceID, want, status))

		select {
		case <-ctx.Done():
			diags.AddError(
				common.ErrorGenericUnexpected,
				fmt.Sprintf("Timeout waiting for %s of Workspace ID: %s to become %s, last known status: %s", r.TypeInfo.Name, workspaceID, want, status),
			)

			return diags
		case <-time.After(encryptionPollInterval):
		}
	}
}
