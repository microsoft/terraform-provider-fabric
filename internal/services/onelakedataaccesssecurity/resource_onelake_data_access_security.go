// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	_ resource.ResourceWithConfigure   = (*resourceOneLakeDataAccessSecurity)(nil)
	_ resource.ResourceWithImportState = (*resourceOneLakeDataAccessSecurity)(nil)
)

type resourceOneLakeDataAccessSecurity struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.OneLakeDataAccessSecurityClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewResourceOneLakeDataAccessSecurity() resource.Resource {
	return &resourceOneLakeDataAccessSecurity{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceOneLakeDataAccessSecurity) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(true)
}

func (r *resourceOneLakeDataAccessSecurity) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema().GetResource(ctx)
}

func (r *resourceOneLakeDataAccessSecurity) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewOneLakeDataAccessSecurityClient()
}

func (r *resourceOneLakeDataAccessSecurity) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceOneLakeDataAccessSecurityModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	r.get(ctx, &state)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceOneLakeDataAccessSecurity) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan, state resourceOneLakeDataAccessSecurityModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	var reqCreate requestCreateOrUpdateOneLakeDataAccessSecurity

	diags := reqCreate.set(ctx, plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	respCreate, err := r.client.CreateOrUpdateDataAccessRoles(ctx, plan.WorkspaceID.ValueString(), plan.ItemID.ValueString(), reqCreate.CreateOrUpdateDataAccessRolesRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	state = plan
	state.setEtag(respCreate.Etag)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceOneLakeDataAccessSecurity) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	var plan resourceOneLakeDataAccessSecurityModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	var reqUpdate requestCreateOrUpdateOneLakeDataAccessSecurity

	reqUpdate.set(ctx, plan)

	respUpdate, err := r.client.CreateOrUpdateDataAccessRoles(ctx, plan.WorkspaceID.ValueString(), plan.ItemID.ValueString(), reqUpdate.CreateOrUpdateDataAccessRolesRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	// Set the etag from the response
	plan.setEtag(respUpdate.Etag)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceOneLakeDataAccessSecurity) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	resp.Diagnostics.AddWarning(
		"delete operation not supported",
		fmt.Sprintf(
			"Resource %s does not support deletion. It will be removed from Terraform state, but no action will be taken in the Fabric. All current settings will remain.",
			r.TypeInfo.Names,
		),
	)

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceOneLakeDataAccessSecurity) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "IMPORT", map[string]any{
		"id": req.ID,
	})

	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			common.ErrorImportIdentifierHeader,
			fmt.Sprintf(common.ErrorImportIdentifierDetails, "WorkspaceID/ItemID"),
		)

		return
	}

	workspaceID, itemID := parts[0], parts[1]

	uuidWorkspaceID, diags := customtypes.NewUUIDValueMust(workspaceID)
	resp.Diagnostics.Append(diags...)

	uuidID, diags := customtypes.NewUUIDValueMust(itemID)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	state := resourceOneLakeDataAccessSecurityModel{
		WorkspaceID: uuidWorkspaceID,
		ItemID:      uuidID,
	}

	if resp.Diagnostics.Append(r.get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceOneLakeDataAccessSecurity) get(ctx context.Context, model *resourceOneLakeDataAccessSecurityModel) diag.Diagnostics {
	respList, err := r.client.ListDataAccessRoles(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		diags.AddError(
			common.ErrorReadHeader,
			"Unable to find an item with 'item_id' "+model.ItemID.ValueString()+" and 'workspace_id': "+model.WorkspaceID.ValueString()+".",
		)

		return diags
	}

	return model.set(ctx, respList)
}
