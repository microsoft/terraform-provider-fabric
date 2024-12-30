// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"
	"fmt"
	"strings"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithModifyPlan  = (*ResourceFabricItemDefinitionProperties[struct{}, struct{}])(nil)
	_ resource.ResourceWithConfigure   = (*ResourceFabricItemDefinitionProperties[struct{}, struct{}])(nil)
	_ resource.ResourceWithImportState = (*ResourceFabricItemDefinitionProperties[struct{}, struct{}])(nil)
)

type ResourceFabricItemDefinitionProperties[Ttfprop, Titemprop any] struct {
	ResourceFabricItemDefinition
	PropertiesAttributes map[string]schema.Attribute
	PropertiesSetter     func(ctx context.Context, from *Titemprop, to *ResourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]) diag.Diagnostics
	ItemGetter           func(ctx context.Context, fabricClient fabric.Client, model ResourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop], fabricItem *FabricItemProperties[Titemprop]) error
}

func NewResourceFabricItemDefinitionProperties[Ttfprop, Titemprop any](config ResourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) resource.Resource {
	return &config
}

func (r *ResourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) { //revive:disable-line:confusing-naming
	resp.TypeName = req.ProviderTypeName + "_" + r.TFName
}

func (r *ResourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) { //revive:disable-line:confusing-naming
	tflog.Debug(ctx, "MODIFY PLAN", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "MODIFY PLAN", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
		"state":  req.State,
	})

	if !req.State.Raw.IsNull() && !req.Plan.Raw.IsNull() {
		var plan, state ResourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]

		resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

		if resp.Diagnostics.HasError() {
			return
		}

		var reqUpdateDefinition requestUpdateFabricItemDefinition

		doUpdateDefinition, diags := fabricItemCheckUpdateDefinition(ctx, plan.Definition, state.Definition, plan.Format, plan.DefinitionUpdateEnabled, r.DefinitionEmpty, r.DefinitionPathKeys, &reqUpdateDefinition)
		if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
			return
		}

		if doUpdateDefinition {
			resp.Diagnostics.AddWarning(
				common.WarningItemDefinitionUpdateHeader,
				fmt.Sprintf(common.WarningItemDefinitionUpdateDetails, r.Name),
			)
		}
	}

	tflog.Debug(ctx, "MODIFY PLAN", map[string]any{
		"action": "end",
	})
}

func (r *ResourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) { //revive:disable-line:confusing-naming
	resp.Schema = getResourceFabricItemDefinitionPropertiesSchema(ctx, *r)
}

func (r *ResourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) { //revive:disable-line:confusing-naming
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
	r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewItemsClient()
}

func (r *ResourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { //revive:disable-line:confusing-naming
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "CREATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
	})

	var plan ResourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqCreate requestCreateFabricItem

	reqCreate.setDisplayName(plan.DisplayName)
	reqCreate.setDescription(plan.Description)
	reqCreate.setType(r.Type)

	if resp.Diagnostics.Append(reqCreate.setDefinition(ctx, plan.Definition, plan.Format, plan.DefinitionUpdateEnabled)...); resp.Diagnostics.HasError() {
		return
	}

	respCreate, err := r.client.CreateItem(ctx, plan.WorkspaceID.ValueString(), reqCreate.CreateItemRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	plan.ID = customtypes.NewUUIDPointerValue(respCreate.ID)
	plan.WorkspaceID = customtypes.NewUUIDPointerValue(respCreate.WorkspaceID)

	if resp.Diagnostics.Append(r.get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) { //revive:disable-line:confusing-naming
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"state": req.State,
	})

	var state ResourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]

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
	if utils.IsErrNotFound(state.ID.ValueString(), &diags, fabcore.ErrCommon.EntityNotFound) {
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

func (r *ResourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { //revive:disable-line:confusing-naming
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "UPDATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
		"state":  req.State,
	})

	var plan, state ResourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdatePlan requestUpdateFabricItem

	if fabricItemCheckUpdate(plan.DisplayName, plan.Description, state.DisplayName, state.Description, &reqUpdatePlan) {
		tflog.Trace(ctx, fmt.Sprintf("updating %s (WorkspaceID: %s ItemID: %s)", r.Name, plan.WorkspaceID.ValueString(), plan.ID.ValueString()))

		_, err := r.client.UpdateItem(ctx, plan.WorkspaceID.ValueString(), plan.ID.ValueString(), reqUpdatePlan.UpdateItemRequest, nil)
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
			return
		}

		if resp.Diagnostics.Append(r.get(ctx, &plan)...); resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}

	var reqUpdateDefinition requestUpdateFabricItemDefinition

	doUpdateDefinition, diags := fabricItemCheckUpdateDefinition(ctx, plan.Definition, state.Definition, plan.Format, plan.DefinitionUpdateEnabled, r.DefinitionEmpty, r.DefinitionPathKeys, &reqUpdateDefinition)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	if doUpdateDefinition {
		tflog.Trace(ctx, fmt.Sprintf("updating %s definition", r.Name))

		_, err := r.client.UpdateItemDefinition(ctx, plan.WorkspaceID.ValueString(), plan.ID.ValueString(), reqUpdateDefinition.UpdateItemDefinitionRequest, nil)
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) { //revive:disable-line:confusing-naming
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "DELETE", map[string]any{
		"state": req.State,
	})

	var state ResourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.DeleteItem(ctx, state.WorkspaceID.ValueString(), state.ID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *ResourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) { //revive:disable-line:confusing-naming
	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "IMPORT", map[string]any{
		"id": req.ID,
	})

	workspaceID, fabricItemID, found := strings.Cut(req.ID, "/")
	if !found {
		resp.Diagnostics.AddError(
			common.ErrorImportIdentifierHeader,
			fmt.Sprintf(
				common.ErrorImportIdentifierDetails,
				fmt.Sprintf("WorkspaceID/%sID", string(r.Type)),
			),
		)

		return
	}

	uuidWorkspaceID, diags := customtypes.NewUUIDValueMust(workspaceID)
	resp.Diagnostics.Append(diags...)

	uuidFabricItemID, diags := customtypes.NewUUIDValueMust(fabricItemID)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var timeout timeouts.Value
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeout)...); resp.Diagnostics.HasError() {
		return
	}

	var definitionUpdateEnabled types.Bool
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("definition_update_enabled"), &definitionUpdateEnabled)...); resp.Diagnostics.HasError() {
		return
	}

	var definition supertypes.MapNestedObjectValueOf[resourceFabricItemDefinitionPartModel]
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("definition"), &definition)...); resp.Diagnostics.HasError() {
		return
	}

	state := ResourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]{
		FabricItemPropertiesModel: FabricItemPropertiesModel[Ttfprop, Titemprop]{
			ID:          uuidFabricItemID,
			WorkspaceID: uuidWorkspaceID,
		},
		DefinitionUpdateEnabled: definitionUpdateEnabled,
		Definition:              definition,
		Timeouts:                timeout,
	}

	if resp.Diagnostics.Append(r.get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) get(ctx context.Context, model *ResourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]) diag.Diagnostics { //revive:disable-line:confusing-naming
	tflog.Trace(ctx, fmt.Sprintf("getting %s by ID: %s", r.Name, model.ID.ValueString()))

	var fabricItem FabricItemProperties[Titemprop]

	err := r.ItemGetter(ctx, *r.pConfigData.FabricClient, *model, &fabricItem)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	model.set(fabricItem)

	return r.PropertiesSetter(ctx, fabricItem.Properties, model)
}
