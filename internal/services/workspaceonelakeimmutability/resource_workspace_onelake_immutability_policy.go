// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceonelakeimmutability

import (
	"context"
	"fmt"
	"strings"

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

var (
	_ resource.ResourceWithConfigure   = (*resourceWorkspaceOneLakeImmutabilityPolicy)(nil)
	_ resource.ResourceWithImportState = (*resourceWorkspaceOneLakeImmutabilityPolicy)(nil)
	_ resource.ResourceWithModifyPlan  = (*resourceWorkspaceOneLakeImmutabilityPolicy)(nil)
)

type resourceWorkspaceOneLakeImmutabilityPolicy struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.OneLakeSettingsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewResourceWorkspaceOneLakeImmutabilityPolicy() resource.Resource {
	return &resourceWorkspaceOneLakeImmutabilityPolicy{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceWorkspaceOneLakeImmutabilityPolicy) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceWorkspaceOneLakeImmutabilityPolicy) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema().GetResource(ctx)
}

func (r *resourceWorkspaceOneLakeImmutabilityPolicy) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewOneLakeSettingsClient()
}

func (r *resourceWorkspaceOneLakeImmutabilityPolicy) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var plan resourceWorkspaceOneLakeImmutabilityPolicyModel
	var state resourceWorkspaceOneLakeImmutabilityPolicyModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	if plan.WorkspaceID.IsNull() || plan.WorkspaceID.IsUnknown() ||
		state.WorkspaceID.IsNull() || state.WorkspaceID.IsUnknown() ||
		plan.Scope.IsNull() || plan.Scope.IsUnknown() ||
		state.Scope.IsNull() || state.Scope.IsUnknown() {
		return
	}

	if plan.WorkspaceID.ValueString() != state.WorkspaceID.ValueString() ||
		plan.Scope.ValueString() != state.Scope.ValueString() ||
		plan.RetentionDays.IsNull() || plan.RetentionDays.IsUnknown() ||
		state.RetentionDays.IsNull() || state.RetentionDays.IsUnknown() {
		return
	}

	if plan.RetentionDays.ValueInt64() < state.RetentionDays.ValueInt64() {
		resp.Diagnostics.Append(retentionDecreaseDiagnostics(plan.RetentionDays.ValueInt64(), state.RetentionDays.ValueInt64())...)
	}
}

func (r *resourceWorkspaceOneLakeImmutabilityPolicy) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan resourceWorkspaceOneLakeImmutabilityPolicyModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(r.modify(ctx, &plan, utils.OperationCreate)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})
}

func (r *resourceWorkspaceOneLakeImmutabilityPolicy) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceWorkspaceOneLakeImmutabilityPolicyModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Read(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	policy, found, diags := getImmutabilityPolicy(ctx, r.client, state.WorkspaceID.ValueString(), state.Scope.ValueString())
	if utils.IsErrNotFound(state.WorkspaceID.ValueString(), &diags, fabcore.ErrCommon.EntityNotFound) {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.Append(diags...)

		return
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	if !found {
		resp.State.RemoveResource(ctx)

		return
	}

	if resp.Diagnostics.Append(state.set(ctx, state.WorkspaceID.ValueString(), policy)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})
}

func (r *resourceWorkspaceOneLakeImmutabilityPolicy) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	var plan resourceWorkspaceOneLakeImmutabilityPolicyModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(r.modify(ctx, &plan, utils.OperationUpdate)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})
}

func (r *resourceWorkspaceOneLakeImmutabilityPolicy) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceWorkspaceOneLakeImmutabilityPolicyModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddWarning(
		"OneLake immutability policy remains configured",
		fmt.Sprintf(
			"The Fabric API does not support deleting or resetting the OneLake immutability policy for workspace %q and scope %q. Terraform will remove the resource from state, but the remote policy remains configured with a %d-day immutability period.",
			state.WorkspaceID.ValueString(),
			state.Scope.ValueString(),
			state.RetentionDays.ValueInt64(),
		),
	)
}

func (r *resourceWorkspaceOneLakeImmutabilityPolicy) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "start",
	})

	segments := strings.Split(req.ID, "/")
	if len(segments) != 2 {
		resp.Diagnostics.AddError(
			common.ErrorImportHeader,
			"Expected import ID in the format <workspace_id>/<scope>.",
		)

		return
	}

	workspaceID, diags := customtypes.NewUUIDValueMust(segments[0])
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	scope := segments[1]
	if !containsImmutabilityScope(scope) {
		resp.Diagnostics.AddError(
			common.ErrorImportHeader,
			fmt.Sprintf("Unsupported OneLake immutability policy scope %q.", scope),
		)

		return
	}

	if resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_id"), workspaceID)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scope"), scope)...)
}

func (r *resourceWorkspaceOneLakeImmutabilityPolicy) modify(ctx context.Context, plan *resourceWorkspaceOneLakeImmutabilityPolicyModel, operation utils.Operation) diag.Diagnostics {
	if plan.Scope.IsNull() || plan.Scope.IsUnknown() ||
		plan.RetentionDays.IsNull() || plan.RetentionDays.IsUnknown() {
		var diags diag.Diagnostics
		diags.AddError(
			"Invalid OneLake immutability policy plan",
			"The policy scope and retention period must be known before applying the policy.",
		)

		return diags
	}

	currentPolicy, found, diags := getImmutabilityPolicy(ctx, r.client, plan.WorkspaceID.ValueString(), plan.Scope.ValueString())
	if diags.HasError() {
		return diags
	}

	if found {
		if diags := validateRetentionNotReduced(plan.RetentionDays.ValueInt64(), int64(*currentPolicy.RetentionDays)); diags.HasError() {
			return diags
		}
	}

	var request requestModifyWorkspaceOneLakeImmutabilityPolicy
	if diags := request.set(*plan); diags.HasError() {
		return diags
	}

	_, err := r.client.ModifyImmutabilityPolicy(ctx, plan.WorkspaceID.ValueString(), request.ImmutabilityPolicyRequest, nil)
	if diags := utils.GetDiagsFromError(ctx, err, operation, nil); diags.HasError() {
		return diags
	}

	updatedPolicy, found, diags := getImmutabilityPolicy(ctx, r.client, plan.WorkspaceID.ValueString(), plan.Scope.ValueString())
	if diags.HasError() {
		return diags
	}

	if !found {
		diags.AddError(
			"OneLake immutability policy was not returned after modification",
			fmt.Sprintf("The API did not return a policy for workspace %q and scope %q after the modify request succeeded.", plan.WorkspaceID.ValueString(), plan.Scope.ValueString()),
		)

		return diags
	}

	return plan.set(ctx, plan.WorkspaceID.ValueString(), updatedPolicy)
}

func retentionDecreaseDiagnostics(requested, current int64) diag.Diagnostics {
	var diags diag.Diagnostics
	diags.AddAttributeError(
		path.Root("retention_days"),
		"Retention period cannot be reduced",
		fmt.Sprintf("The existing OneLake immutability policy is set to %d days, but the requested value is %d days. The Fabric API does not allow reducing the retention period.", current, requested),
	)

	return diags
}

func validateRetentionNotReduced(requested, current int64) diag.Diagnostics {
	if requested < current {
		return retentionDecreaseDiagnostics(requested, current)
	}

	return nil
}

func containsImmutabilityScope(scope string) bool {
	for _, possibleScope := range fabcore.PossibleImmutabilityScopeValues() {
		if string(possibleScope) == scope {
			return true
		}
	}

	return false
}
