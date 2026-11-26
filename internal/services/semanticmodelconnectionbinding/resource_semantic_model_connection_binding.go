// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package semanticmodelconnectionbinding

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabsemanticmodel "github.com/microsoft/fabric-sdk-go/fabric/semanticmodel"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var (
	_ resource.ResourceWithConfigure   = (*resourceSemanticModelConnectionBinding)(nil)
	_ resource.ResourceWithImportState = (*resourceSemanticModelConnectionBinding)(nil)
)

type resourceSemanticModelConnectionBinding struct {
	pConfigData *pconfig.ProviderData
	client      *fabsemanticmodel.ItemsClient
	coreClient  *fabcore.ItemsClient
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
	resp.Schema = itemSchema(ctx, false).GetResource(ctx)
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

	fabricClient := *pConfigData.FabricClient
	r.client = fabsemanticmodel.NewClientFactoryWithClient(fabricClient).NewItemsClient()
	r.coreClient = fabcore.NewClientFactoryWithClient(fabricClient).NewItemsClient()
}

func (r *resourceSemanticModelConnectionBinding) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{"action": "start"})

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

	resp.Diagnostics.Append(r.bind(ctx, plan, fabsemanticmodel.ConnectivityType(plan.ConnectivityType.ValueString()), utils.OperationCreate)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.setPostWriteState(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, "CREATE", map[string]any{"action": "end"})
}

func (r *resourceSemanticModelConnectionBinding) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{"action": "start"})

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

	bindings, err := listItemConnections(ctx, r.coreClient, state.WorkspaceID.ValueString(), state.SemanticModelID.ValueString())
	if err != nil {
		errorDiags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound)
		if utils.IsErrNotFound(state.SemanticModelID.ValueString(), &errorDiags, fabcore.ErrCommon.EntityNotFound) {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.Append(errorDiags...)

		return
	}

	details, diags := state.ConnectionDetails.Get(ctx)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	binding, found, err := findConnectionBinding(bindings, details.Path.ValueString(), details.Type.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Semantic Model connection binding response", err.Error())

		return
	}

	if !found || binding.ConnectivityType == nil || *binding.ConnectivityType == fabcore.ConnectivityTypeNone {
		resp.State.RemoveResource(ctx)

		return
	}

	if resp.Diagnostics.Append(state.set(ctx, state.WorkspaceID.ValueString(), state.SemanticModelID.ValueString(), binding)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, "READ", map[string]any{"action": "end"})
}

func (r *resourceSemanticModelConnectionBinding) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{"action": "start"})

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

	resp.Diagnostics.Append(r.bind(ctx, plan, fabsemanticmodel.ConnectivityType(plan.ConnectivityType.ValueString()), utils.OperationUpdate)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.setPostWriteState(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{"action": "end"})
}

func (r *resourceSemanticModelConnectionBinding) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{"action": "start"})

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

	resp.Diagnostics.Append(r.bind(ctx, state, fabsemanticmodel.ConnectivityTypeNone, utils.OperationDelete)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{"action": "end"})
}

func (r *resourceSemanticModelConnectionBinding) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "IMPORT", map[string]any{"action": "start"})
	tflog.Trace(ctx, "IMPORT", map[string]any{"id": req.ID})

	identifier, err := parseConnectionBindingIdentifier(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			common.ErrorImportIdentifierHeader,
			fmt.Sprintf("%s: %v", fmt.Sprintf(common.ErrorImportIdentifierDetails, "WorkspaceID/SemanticModelID/Base64URLPath/Base64URLType"), err),
		)

		return
	}

	workspaceUUID, diags := customtypes.NewUUIDValueMust(identifier.WorkspaceID)
	resp.Diagnostics.Append(diags...)
	semanticModelUUID, diags := customtypes.NewUUIDValueMust(identifier.SemanticModelID)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var timeout timeouts.Value
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeout)...); resp.Diagnostics.HasError() {
		return
	}

	state := resourceSemanticModelConnectionBindingModel{
		WorkspaceID:     workspaceUUID,
		SemanticModelID: semanticModelUUID,
		Timeouts:        timeout,
	}

	readTimeout, diags := timeout.Read(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	bindings, err := listItemConnections(ctx, r.coreClient, identifier.WorkspaceID, identifier.SemanticModelID)
	if err != nil {
		resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound)...)

		return
	}

	binding, found, err := findConnectionBinding(bindings, identifier.ConnectionPath, identifier.ConnectionType)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Semantic Model connection binding response", err.Error())

		return
	}

	if !found || binding.ConnectivityType == nil || *binding.ConnectivityType == fabcore.ConnectivityTypeNone {
		resp.Diagnostics.AddError(
			"Semantic Model connection binding not found",
			fmt.Sprintf(
				"No binding exists or is currently unbound for connection details path %q and type %q on semantic model %q.",
				identifier.ConnectionPath,
				identifier.ConnectionType,
				identifier.SemanticModelID,
			),
		)

		return
	}

	if resp.Diagnostics.Append(state.set(ctx, identifier.WorkspaceID, identifier.SemanticModelID, binding)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, "IMPORT", map[string]any{"action": "end"})
}

func (r *resourceSemanticModelConnectionBinding) bind(
	ctx context.Context,
	plan resourceSemanticModelConnectionBindingModel,
	connectivityType fabsemanticmodel.ConnectivityType,
	operation utils.Operation,
) diag.Diagnostics {
	var request requestBindSemanticModelConnection
	if diags := request.set(ctx, plan, connectivityType); diags.HasError() {
		return diags
	}

	_, err := r.client.BindSemanticModelConnection(
		ctx,
		plan.WorkspaceID.ValueString(),
		plan.SemanticModelID.ValueString(),
		request.BindSemanticModelConnectionRequest,
		nil,
	)

	return utils.GetDiagsFromError(ctx, err, operation, nil)
}

func (r *resourceSemanticModelConnectionBinding) setPostWriteState(ctx context.Context, model *resourceSemanticModelConnectionBindingModel) diag.Diagnostics {
	details, diags := model.ConnectionDetails.Get(ctx)
	if diags.HasError() {
		return diags
	}

	model.ID = types.StringValue(connectionBindingIdentifier(
		model.WorkspaceID.ValueString(),
		model.SemanticModelID.ValueString(),
		details.Path.ValueString(),
		details.Type.ValueString(),
	))
	model.DisplayName = types.StringNull()
	model.GatewayID = customtypes.NewUUIDPointerValue(nil)

	return nil
}

func listItemConnections(ctx context.Context, client *fabcore.ItemsClient, workspaceID, semanticModelID string) ([]fabcore.ItemConnection, error) {
	return client.ListItemConnections(ctx, workspaceID, semanticModelID, nil)
}

func findConnectionBinding(bindings []fabcore.ItemConnection, connectionPath, connectionType string) (fabcore.ItemConnection, bool, error) {
	var foundBinding fabcore.ItemConnection
	found := false

	for _, binding := range bindings {
		if binding.ConnectionDetails == nil || binding.ConnectionDetails.Path == nil || binding.ConnectionDetails.Type == nil {
			return fabcore.ItemConnection{}, false, errors.New("item connection response is missing required connection details")
		}

		if binding.ConnectivityType == nil {
			return fabcore.ItemConnection{}, false, errors.New("item connection response is missing required connectivity type")
		}

		if *binding.ConnectionDetails.Path != connectionPath || *binding.ConnectionDetails.Type != connectionType {
			continue
		}

		if found {
			return fabcore.ItemConnection{}, false, fmt.Errorf("multiple bindings match connection details path %q and type %q", connectionPath, connectionType)
		}

		foundBinding = binding
		found = true
	}

	return foundBinding, found, nil
}
