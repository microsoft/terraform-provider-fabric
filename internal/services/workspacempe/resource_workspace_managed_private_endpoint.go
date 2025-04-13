// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspacempe

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	_ resource.ResourceWithConfigure = (*resourceWorkspaceManagedPrivateEndpoint)(nil)
	// _ resource.ResourceWithImportState = (*resourceWorkspaceManagedPrivateEndpoint)(nil).
)

type resourceWorkspaceManagedPrivateEndpoint struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.ManagedPrivateEndpointsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewResourceWorkspaceManagedPrivateEndpoint() resource.Resource {
	return &resourceWorkspaceManagedPrivateEndpoint{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceWorkspaceManagedPrivateEndpoint) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceWorkspaceManagedPrivateEndpoint) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema(false).GetResource(ctx)
}

func (r *resourceWorkspaceManagedPrivateEndpoint) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewManagedPrivateEndpointsClient()
}

func (r *resourceWorkspaceManagedPrivateEndpoint) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan, state resourceWorkspaceManagedPrivateEndpointModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	state.Timeouts = plan.Timeouts

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqCreate requestCreateWorkspaceManagedPrivateEndpoint

	reqCreate.set(plan)

	respCreate, err := r.client.CreateWorkspaceManagedPrivateEndpoint(ctx, plan.WorkspaceID.ValueString(), reqCreate.CreateManagedPrivateEndpointRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.set(ctx, plan.WorkspaceID.ValueString(), respCreate.ManagedPrivateEndpoint)...); resp.Diagnostics.HasError() {
		return
	}

	state.ID = customtypes.NewUUIDPointerValue(respCreate.ID)
	state.WorkspaceID = plan.WorkspaceID
	state.RequestMessage = plan.RequestMessage

	if resp.Diagnostics.Append(r.get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceWorkspaceManagedPrivateEndpoint) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceWorkspaceManagedPrivateEndpointModel

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
	if utils.IsErrNotFound(state.ID.ValueString(), &diags, errors.New("PrivateEndpointNotFound")) {
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

func (r *resourceWorkspaceManagedPrivateEndpoint) Update(ctx context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
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

func (r *resourceWorkspaceManagedPrivateEndpoint) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	var state resourceWorkspaceManagedPrivateEndpointModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.DeleteWorkspaceManagedPrivateEndpoint(ctx, state.WorkspaceID.ValueString(), state.ID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceWorkspaceManagedPrivateEndpoint) get(ctx context.Context, model *resourceWorkspaceManagedPrivateEndpointModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET", map[string]any{
		"workspace_id": model.WorkspaceID.ValueString(),
		"id":           model.ID.ValueString(),
	})

	for {
		respGet, err := r.client.GetWorkspaceManagedPrivateEndpoint(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrManagedPrivateEndpoint.PrivateEndpointNotFound); diags.HasError() {
			return diags
		}

		provisioningStateStr := string(*respGet.ProvisioningState)

		switch *respGet.ProvisioningState {
		case fabcore.PrivateEndpointProvisioningStateFailed:
			var diags diag.Diagnostics

			diags.AddError(
				common.ErrorReadHeader,
				r.TypeInfo.Name+" provisioning state: "+provisioningStateStr,
			)

			return diags

		case fabcore.PrivateEndpointProvisioningStateSucceeded:
			return model.set(ctx, model.WorkspaceID.ValueString(), respGet.ManagedPrivateEndpoint)

		default:
			tflog.Info(ctx, r.TypeInfo.Name+" provisioning in progress, waiting 30 seconds before retrying", map[string]any{
				"provisioning_state": provisioningStateStr,
			})

			time.Sleep(30 * time.Second) // lintignore:R018
		}
	}
}
