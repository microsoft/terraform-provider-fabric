// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceencryption

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
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

	wsDetail, diags := r.get(ctx, &state.baseWorkspaceEncryptionModel)
	if utils.IsErrNotFound(state.WorkspaceID.ValueString(), &diags, fabcore.ErrCommon.EntityNotFound) {
		resp.State.RemoveResource(ctx)

		resp.Diagnostics.Append(diags...)

		return
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	detail, diagsGet := state.EncryptionDetails.Get(ctx)
	if resp.Diagnostics.Append(diagsGet...); resp.Diagnostics.HasError() {
		return
	}

	// A workspace without a customer-managed key reports Disabled, which is the absence of this resource.
	if detail == nil || detail.EncryptionStatus.ValueString() == string(fabcore.WorkspaceEncryptionStatusDisabled) {
		resp.State.RemoveResource(ctx)

		return
	}

	// Encryption can fail outside of Terraform, for example when the key is revoked. Without this warning the
	// plan would be empty and the broken workspace would go unnoticed.
	if detail.EncryptionStatus.ValueString() == string(fabcore.WorkspaceEncryptionStatusFailed) {
		var failedItemsMsg string

		if wsDetail != nil {
			if failedItems := formatFailedItems(*wsDetail); failedItems != "" {
				failedItemsMsg = "\n\nFailed items:\n" + failedItems
			}
		}

		resp.Diagnostics.AddWarning(
			r.TypeInfo.Name+" failed",
			fmt.Sprintf(
				"%s is in the Failed state for Workspace ID: %s. Verify that the key exists, is enabled, and that the 'Fabric Platform CMK' application can wrap and unwrap it, then re-apply to retry.%s",
				r.TypeInfo.Name,
				state.WorkspaceID.ValueString(),
				failedItemsMsg,
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

	var rawResp *http.Response
	ctxCapture := policy.WithCaptureResponse(ctx, &rawResp)

	// Resetting removes the customer-managed key, after which the workspace falls back to Microsoft-managed keys.
	_, err := r.client.ResetWorkspaceEncryption(ctxCapture, state.WorkspaceID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(r.waitForStatus(ctx, state.WorkspaceID.ValueString(), fabcore.WorkspaceEncryptionStatusDisabled, nil, getRetryAfter(rawResp))...); resp.Diagnostics.HasError() {
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

	_, diags = r.get(ctx, &state.baseWorkspaceEncryptionModel)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	detail, diagsGet := state.EncryptionDetails.Get(ctx)
	if resp.Diagnostics.Append(diagsGet...); resp.Diagnostics.HasError() {
		return
	}

	if detail == nil || detail.EncryptionStatus.ValueString() == string(fabcore.WorkspaceEncryptionStatusDisabled) {
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

	if diags := reqAssign.set(ctx, *model); diags.HasError() {
		return diags
	}

	var rawResp *http.Response
	ctxCapture := policy.WithCaptureResponse(ctx, &rawResp)

	_, err := r.client.AssignWorkspaceEncryption(ctxCapture, model.WorkspaceID.ValueString(), reqAssign.AssignWorkspaceEncryptionRequest, nil)
	if diags := utils.GetDiagsFromError(ctx, err, operation, nil); diags.HasError() {
		return diags
	}

	return r.waitForStatus(ctx, model.WorkspaceID.ValueString(), fabcore.WorkspaceEncryptionStatusActive, &model.baseWorkspaceEncryptionModel, getRetryAfter(rawResp))
}

func (r *resourceWorkspaceEncryption) get(ctx context.Context, model *baseWorkspaceEncryptionModel) (*fabcore.WorkspaceEncryptionDetail, diag.Diagnostics) {
	tflog.Trace(ctx, fmt.Sprintf("getting %s for Workspace ID: %s", r.TypeInfo.Name, model.WorkspaceID.ValueString()))

	respGet, err := r.client.GetWorkspaceEncryption(ctx, model.WorkspaceID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		return nil, diags
	}

	return &respGet.WorkspaceEncryptionDetail, model.set(ctx, model.WorkspaceID.ValueString(), respGet.WorkspaceEncryptionDetail)
}

// waitForStatus polls the encryption status until it settles on want, because assign and reset are asynchronous.
func (r *resourceWorkspaceEncryption) waitForStatus(
	ctx context.Context,
	workspaceID string,
	want fabcore.WorkspaceEncryptionStatus,
	model *baseWorkspaceEncryptionModel,
	initialPollInterval time.Duration,
) diag.Diagnostics {
	var diags diag.Diagnostics

	pollInterval := initialPollInterval
	if pollInterval <= 0 {
		pollInterval = encryptionPollInterval
	}

	for {
		tflog.Trace(ctx, fmt.Sprintf("waiting for %s of Workspace ID: %s to become %s (polling in %s)", r.TypeInfo.Name, workspaceID, want, pollInterval))

		select {
		case <-ctx.Done():
			diags.AddError(
				common.ErrorGenericUnexpected,
				fmt.Sprintf("Timeout waiting for %s of Workspace ID: %s to become %s", r.TypeInfo.Name, workspaceID, want),
			)

			return diags
		case <-time.After(pollInterval):
		}

		var rawResp *http.Response
		ctxCapture := policy.WithCaptureResponse(ctx, &rawResp)

		respGet, err := r.client.GetWorkspaceEncryption(ctxCapture, workspaceID, nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
			return diags
		}

		status := encryptionStatus(respGet.WorkspaceEncryptionDetail)

		if status == fabcore.WorkspaceEncryptionStatusFailed {
			var failedItemsMsg string
			if failedItems := formatFailedItems(respGet.WorkspaceEncryptionDetail); failedItems != "" {
				failedItemsMsg = "\n\nFailed items:\n" + failedItems
			}

			diags.AddError(
				common.ErrorGenericUnexpected,
				fmt.Sprintf(
					"%s failed for Workspace ID: %s. Verify that the key exists, is enabled, and that the 'Fabric Platform CMK' application can wrap and unwrap it.%s",
					r.TypeInfo.Name,
					workspaceID,
					failedItemsMsg,
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

		pollInterval = getRetryAfter(rawResp)
	}
}

func getRetryAfter(resp *http.Response) time.Duration {
	if resp == nil || resp.Header == nil {
		return encryptionPollInterval
	}

	val := resp.Header.Get("Retry-After")
	if val == "" {
		return encryptionPollInterval
	}

	seconds, err := strconv.Atoi(val)
	if err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}

	return encryptionPollInterval
}

func formatFailedItems(detail fabcore.WorkspaceEncryptionDetail) string {
	var failedItems []string

	for _, itemsDetail := range detail.WorkspaceEncryptionItemsDetails {
		if itemsDetail.EncryptionStatus == nil || *itemsDetail.EncryptionStatus != fabcore.WorkspaceEncryptionStatusFailed {
			continue
		}

		for _, item := range itemsDetail.Items {
			if itemStr := formatFailedItem(item); itemStr != "" {
				failedItems = append(failedItems, itemStr)
			}
		}
	}

	return strings.Join(failedItems, "\n")
}

func formatFailedItem(item fabcore.WorkspaceEncryptionItem) string {
	var label string
	if item.DisplayName != nil && *item.DisplayName != "" {
		label = *item.DisplayName
	}

	var details []string
	if item.Type != nil && *item.Type != "" {
		details = append(details, "Type: "+*item.Type)
	}

	if item.ID != nil && *item.ID != "" {
		details = append(details, "ID: "+*item.ID)
	}

	if len(details) > 0 {
		if label != "" {
			return fmt.Sprintf("- %s (%s)", label, strings.Join(details, ", "))
		}

		return "- " + strings.Join(details, ", ")
	}

	if label != "" {
		return "- " + label
	}

	return ""
}
