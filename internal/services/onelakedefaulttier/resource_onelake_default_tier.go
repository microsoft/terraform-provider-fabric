// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedefaulttier

import (
	"context"
	"fmt"

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
	_ resource.ResourceWithConfigure   = (*resourceOneLakeDefaultTier)(nil)
	_ resource.ResourceWithImportState = (*resourceOneLakeDefaultTier)(nil)
)

type resourceOneLakeDefaultTier struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.OneLakeSettingsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewResourceOneLakeDefaultTier() resource.Resource {
	return &resourceOneLakeDefaultTier{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceOneLakeDefaultTier) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceOneLakeDefaultTier) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema().GetResource(ctx)
}

func (r *resourceOneLakeDefaultTier) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *resourceOneLakeDefaultTier) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{"action": "start"})

	var plan resourceOneLakeDefaultTierModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	respCreate, err := r.client.ModifyDefaultTier(ctx, plan.WorkspaceID.ValueString(), fabcore.OneLakeAccessTier(plan.DefaultTier.ValueString()), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if respCreate.DefaultTier == nil {
		resp.Diagnostics.AddError(
			"Missing OneLake Default Tier",
			fmt.Sprintf("The Modify Default Tier API did not return the resulting tier for workspace %s.", plan.WorkspaceID.ValueString()),
		)

		return
	}

	plan.set(plan.WorkspaceID.ValueString(), respCreate.DefaultTier)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, "CREATE", map[string]any{"action": "end"})
}

func (r *resourceOneLakeDefaultTier) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{"action": "start"})

	var state resourceOneLakeDefaultTierModel

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

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, "READ", map[string]any{"action": "end"})
}

func (r *resourceOneLakeDefaultTier) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{"action": "start"})

	var plan resourceOneLakeDefaultTierModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	respUpdate, err := r.client.ModifyDefaultTier(ctx, plan.WorkspaceID.ValueString(), fabcore.OneLakeAccessTier(plan.DefaultTier.ValueString()), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if respUpdate.DefaultTier == nil {
		resp.Diagnostics.AddError(
			"Missing OneLake Default Tier",
			fmt.Sprintf("The Modify Default Tier API did not return the resulting tier for workspace %s.", plan.WorkspaceID.ValueString()),
		)

		return
	}

	plan.set(plan.WorkspaceID.ValueString(), respUpdate.DefaultTier)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{"action": "end"})
}

func (r *resourceOneLakeDefaultTier) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"OneLake default tier remains configured",
		fmt.Sprintf(
			"Resource %s does not support deletion. The current OneLake default tier remains configured in Fabric; change it explicitly before removing the resource if a different tier is desired.",
			r.TypeInfo.Name,
		),
	)
	resp.State.RemoveResource(ctx)
}

func (r *resourceOneLakeDefaultTier) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "IMPORT", map[string]any{"action": "start"})
	tflog.Trace(ctx, "IMPORT", map[string]any{"id": req.ID})

	workspaceID, diags := customtypes.NewUUIDValueMust(req.ID)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_id"), workspaceID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), workspaceID)...)

	tflog.Debug(ctx, "IMPORT", map[string]any{"action": "end"})
}

func (r *resourceOneLakeDefaultTier) get(ctx context.Context, model *resourceOneLakeDefaultTierModel) diag.Diagnostics {
	workspaceID := model.WorkspaceID.ValueString()
	tflog.Trace(ctx, "getting OneLake default tier for Workspace ID: "+workspaceID)

	respGet, err := r.client.GetSettings(ctx, workspaceID, nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		return diags
	}

	if respGet.Lifecycle == nil || respGet.Lifecycle.DefaultTier == nil {
		var diags diag.Diagnostics
		diags.AddError(
			"OneLake Default Tier Unavailable",
			fmt.Sprintf("Get Settings did not return lifecycle.defaultTier for workspace %s.", workspaceID),
		)

		return diags
	}

	model.set(workspaceID, respGet.Lifecycle.DefaultTier)

	return nil
}
