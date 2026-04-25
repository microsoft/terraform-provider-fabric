// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package warehousesqlauditsetting

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabwarehouse "github.com/microsoft/fabric-sdk-go/fabric/warehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var _ resource.ResourceWithConfigure = (*resourceWarehouseSQLAuditSettings)(nil)

type resourceWarehouseSQLAuditSettings struct {
	pConfigData *pconfig.ProviderData
	client      *fabwarehouse.SQLAuditSettingsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewResourceWarehouseSQLAuditSettings() resource.Resource {
	return &resourceWarehouseSQLAuditSettings{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceWarehouseSQLAuditSettings) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceWarehouseSQLAuditSettings) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema().GetResource(ctx)
}

func (r *resourceWarehouseSQLAuditSettings) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = fabwarehouse.NewClientFactoryWithClient(*pConfigData.FabricClient).NewSQLAuditSettingsClient()
}

func (r *resourceWarehouseSQLAuditSettings) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan resourceWarehouseSQLAuditSettingsModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateWarehouseSQLAuditSettings

	reqUpdate.set(plan)

	_, err := r.client.UpdateSQLAuditSettings(
		ctx,
		plan.WorkspaceID.ValueString(),
		plan.ItemID.ValueString(),
		reqUpdate.SQLAuditSettingsUpdate,
		nil,
	)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	// Set audit actions and groups
	var reqSetAuditActions requestSetAuditActionsAndGroups

	if resp.Diagnostics.Append(reqSetAuditActions.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	if reqSetAuditActions.AuditActionsAndGroups != nil {
		_, err = r.client.SetAuditActionsAndGroups(
			ctx,
			plan.WorkspaceID.ValueString(),
			plan.ItemID.ValueString(),
			reqSetAuditActions.AuditActionsAndGroups,
			nil,
		)
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
			return
		}
	}

	// Read back the full state
	if resp.Diagnostics.Append(r.get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})
}

func (r *resourceWarehouseSQLAuditSettings) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceWarehouseSQLAuditSettingsModel

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
	if utils.IsErrNotFound(state.ItemID.ValueString(), &diags, fabcore.ErrCommon.EntityNotFound) {
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
}

func (r *resourceWarehouseSQLAuditSettings) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	var plan resourceWarehouseSQLAuditSettingsModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Update settings
	var reqUpdate requestUpdateWarehouseSQLAuditSettings

	reqUpdate.set(plan)

	_, err := r.client.UpdateSQLAuditSettings(
		ctx,
		plan.WorkspaceID.ValueString(),
		plan.ItemID.ValueString(),
		reqUpdate.SQLAuditSettingsUpdate,
		nil,
	)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	// Set audit actions and groups
	var reqSetAuditActions requestSetAuditActionsAndGroups

	if resp.Diagnostics.Append(reqSetAuditActions.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	if reqSetAuditActions.AuditActionsAndGroups != nil {
		_, err = r.client.SetAuditActionsAndGroups(
			ctx,
			plan.WorkspaceID.ValueString(),
			plan.ItemID.ValueString(),
			reqSetAuditActions.AuditActionsAndGroups,
			nil,
		)
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
			return
		}
	}

	if resp.Diagnostics.Append(r.get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})
}

func (r *resourceWarehouseSQLAuditSettings) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	resp.Diagnostics.AddWarning(
		"delete operation not supported",
		fmt.Sprintf(
			"Resource %s does not support deletion. It will be removed from Terraform state, but no action will be taken in the Fabric. All current settings will remain.",
			r.TypeInfo.Name,
		),
	)

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceWarehouseSQLAuditSettings) get(ctx context.Context, model *resourceWarehouseSQLAuditSettingsModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s for Warehouse ID: %s in Workspace ID: %s", r.TypeInfo.Name, model.ItemID.ValueString(), model.WorkspaceID.ValueString()))

	respGet, err := r.client.GetSQLAuditSettings(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		return diags
	}

	return model.set(ctx, respGet.SQLAuditSettings)
}
