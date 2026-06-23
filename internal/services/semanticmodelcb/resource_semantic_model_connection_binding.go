// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package semanticmodelcb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabsemanticmodel "github.com/microsoft/fabric-sdk-go/fabric/semanticmodel"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ resource.ResourceWithConfigure = (*resourceSemanticModelConnectionBinding)(nil)

type resourceSemanticModelConnectionBinding struct {
	pConfigData *pconfig.ProviderData
	client      *fabsemanticmodel.ItemsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewResourceSemanticModelConnectionBinding() resource.Resource {
	return &resourceSemanticModelConnectionBinding{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceSemanticModelConnectionBinding) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceSemanticModelConnectionBinding) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema().GetResource(ctx)
}

func (r *resourceSemanticModelConnectionBinding) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = fabsemanticmodel.NewClientFactoryWithClient(*pConfigData.FabricClient).NewItemsClient()
}

func (r *resourceSemanticModelConnectionBinding) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan resourceSemanticModelConnectionBindingModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(r.bind(ctx, &plan, utils.OperationCreate)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})
}

// Read is a state passthrough — the Fabric API exposes no read operation for connection bindings.
func (r *resourceSemanticModelConnectionBinding) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceSemanticModelConnectionBindingModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Read(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})
}

func (r *resourceSemanticModelConnectionBinding) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	var plan resourceSemanticModelConnectionBindingModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(r.bind(ctx, &plan, utils.OperationUpdate)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})
}

func (r *resourceSemanticModelConnectionBinding) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	var state resourceSemanticModelConnectionBindingModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if state.ConnectivityType.ValueString() == string(fabsemanticmodel.ConnectivityTypeNone) {
		tflog.Debug(ctx, "DELETE", map[string]any{
			"action": "skip-unbind",
			"reason": "connectivity_type is already None",
		})

		resp.State.RemoveResource(ctx)

		return
	}

	var reqUnbind requestBindSemanticModelConnection
	if resp.Diagnostics.Append(reqUnbind.setUnbind(ctx, state)...); resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.BindSemanticModelConnection(
		ctx,
		state.WorkspaceID.ValueString(),
		state.SemanticModelID.ValueString(),
		reqUnbind.BindSemanticModelConnectionRequest,
		nil,
	)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceSemanticModelConnectionBinding) bind(ctx context.Context, model *resourceSemanticModelConnectionBindingModel, op utils.Operation) diag.Diagnostics {
	var diags diag.Diagnostics

	var reqBind requestBindSemanticModelConnection
	if d := reqBind.set(ctx, *model); d.HasError() {
		diags.Append(d...)

		return diags
	}

	_, err := r.client.BindSemanticModelConnection(
		ctx,
		model.WorkspaceID.ValueString(),
		model.SemanticModelID.ValueString(),
		reqBind.BindSemanticModelConnectionRequest,
		nil,
	)

	diags.Append(utils.GetDiagsFromError(ctx, err, op, nil)...)

	return diags
}
