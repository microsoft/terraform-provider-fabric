// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakediagnosticsetting

import (
	"context"
	"fmt"

	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts" //revive:disable-line:import-alias-naming
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
	_ resource.ResourceWithConfigure   = (*resourceOneLakeDiagnosticSetting)(nil)
	_ resource.ResourceWithImportState = (*resourceOneLakeDiagnosticSetting)(nil)
)

type resourceOneLakeDiagnosticSetting struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.OneLakeSettingsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewResourceOneLakeDiagnosticSetting() resource.Resource {
	return &resourceOneLakeDiagnosticSetting{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceOneLakeDiagnosticSetting) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceOneLakeDiagnosticSetting) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema().GetResource(ctx)
}

func (r *resourceOneLakeDiagnosticSetting) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceOneLakeDiagnosticSetting) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{"action": "start"})

	var plan resourceOneLakeDiagnosticSettingModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var request requestModifyOneLakeDiagnosticSetting
	if resp.Diagnostics.Append(request.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(r.modifyDiagnostics(ctx, plan.WorkspaceID.ValueString(), request.OneLakeDiagnosticSettings, utils.OperationCreate)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(r.get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	tflog.Debug(ctx, "CREATE", map[string]any{"action": "end"})
}

func (r *resourceOneLakeDiagnosticSetting) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{"action": "start"})

	var state resourceOneLakeDiagnosticSettingModel
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

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	tflog.Debug(ctx, "READ", map[string]any{"action": "end"})
}

func (r *resourceOneLakeDiagnosticSetting) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{"action": "start"})

	var plan resourceOneLakeDiagnosticSettingModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var request requestModifyOneLakeDiagnosticSetting
	if resp.Diagnostics.Append(request.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(r.modifyDiagnostics(ctx, plan.WorkspaceID.ValueString(), request.OneLakeDiagnosticSettings, utils.OperationUpdate)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(r.get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	tflog.Debug(ctx, "UPDATE", map[string]any{"action": "end"})
}

func (r *resourceOneLakeDiagnosticSetting) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{"action": "start"})

	var state resourceOneLakeDiagnosticSettingModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp.Diagnostics.AddWarning(
		"Disabling OneLake diagnostics",
		"OneLake diagnostic settings cannot be deleted. Destroying this resource disables diagnostics and clears the destination.",
	)

	settings := fabcore.OneLakeDiagnosticSettings{
		Status: new(diagnosticStatusDisabled),
		Destination: &fabcore.LakehouseOneLakeDiagnosticSettingsDestination{
			Type: new(destinationTypeLakehouse),
		},
	}

	if resp.Diagnostics.Append(r.modifyDiagnostics(ctx, state.WorkspaceID.ValueString(), settings, utils.OperationDelete)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{"action": "end"})
}

func (r *resourceOneLakeDiagnosticSetting) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "IMPORT", map[string]any{"action": "start"})

	workspaceID, diags := customtypes.NewUUIDValueMust(req.ID)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var timeout timeoutsR.Value
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeout)...); resp.Diagnostics.HasError() {
		return
	}

	state := resourceOneLakeDiagnosticSettingModel{
		WorkspaceID: workspaceID,
		Timeouts:    timeout,
	}

	if resp.Diagnostics.Append(r.get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, "IMPORT", map[string]any{"action": "end"})
}

func (r *resourceOneLakeDiagnosticSetting) modifyDiagnostics(
	ctx context.Context,
	workspaceID string,
	settings fabcore.OneLakeDiagnosticSettings,
	operation utils.Operation,
) diag.Diagnostics {
	_, err := r.client.ModifyDiagnostics(
		ctx,
		workspaceID,
		&fabcore.OneLakeSettingsClientBeginModifyDiagnosticsOptions{
			ModifyOneLakeDiagnosticSettingRequest: &settings,
		},
	)

	return utils.GetDiagsFromError(ctx, err, operation, nil)
}

func (r *resourceOneLakeDiagnosticSetting) get(ctx context.Context, model *resourceOneLakeDiagnosticSettingModel) diag.Diagnostics {
	respGet, err := r.client.GetSettings(ctx, model.WorkspaceID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		return diags
	}

	if respGet.Diagnostics == nil {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic(common.ErrorReadHeader, "The OneLake settings response omitted diagnostics."),
		}
	}

	return model.set(ctx, model.WorkspaceID.ValueString(), *respGet.Diagnostics)
}
