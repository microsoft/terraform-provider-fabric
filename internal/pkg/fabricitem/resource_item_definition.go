// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithModifyPlan  = (*ResourceFabricItemDefinition)(nil)
	_ resource.ResourceWithConfigure   = (*ResourceFabricItemDefinition)(nil)
	_ resource.ResourceWithImportState = (*ResourceFabricItemDefinition)(nil)
)

type ResourceFabricItemDefinition struct {
	pConfigData                 *pconfig.ProviderData
	client                      *fabcore.ItemsClient
	Type                        fabcore.ItemType
	Name                        string
	NameRenameAllowed           bool
	TFName                      string
	MarkdownDescription         string
	DisplayNameMaxLength        int
	DescriptionMaxLength        int
	FormatTypeDefault           string
	FormatTypes                 []string
	DefinitionPathDocsURL       string
	DefinitionPathKeys          []string
	DefinitionPathKeysValidator []validator.Map
	DefinitionRequired          bool
	DefinitionEmpty             string
}

func NewResourceFabricItemDefinition(config ResourceFabricItemDefinition) resource.Resource { //revive:disable-line:argument-limit
	return &config
}

func (r *ResourceFabricItemDefinition) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.TFName
}

func (r *ResourceFabricItemDefinition) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	tflog.Debug(ctx, "MODIFY PLAN", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "MODIFY PLAN", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
		"state":  req.State,
	})

	if !req.State.Raw.IsNull() && !req.Plan.Raw.IsNull() {
		var plan, state resourceFabricItemDefinitionModel

		resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

		if resp.Diagnostics.HasError() {
			return
		}

		var reqUpdate requestUpdateFabricItemDefinitionDefinition

		doUpdateDefinition, diags := r.checkUpdateDefinition(ctx, plan, state, &reqUpdate)
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

func (r *ResourceFabricItemDefinition) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = GetResourceFabricItemDefinitionSchema(ctx, *r)
}

func (r *ResourceFabricItemDefinition) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ResourceFabricItemDefinition) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "CREATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
	})

	var plan resourceFabricItemDefinitionModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqCreate requestCreateFabricItemDefinition

	if resp.Diagnostics.Append(reqCreate.set(ctx, plan, r.Type)...); resp.Diagnostics.HasError() {
		return
	}

	respCreate, err := r.client.CreateItem(ctx, plan.WorkspaceID.ValueString(), reqCreate.CreateItemRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	plan.set(respCreate.Item)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourceFabricItemDefinition) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"state": req.State,
	})

	var state resourceFabricItemDefinitionModel

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

func (r *ResourceFabricItemDefinition) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "UPDATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
		"state":  req.State,
	})

	var plan, state resourceFabricItemDefinitionModel

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

	var reqUpdatePlan requestUpdateFabricItemDefinition

	if r.checkUpdateItem(plan, state, &reqUpdatePlan) {
		tflog.Trace(ctx, fmt.Sprintf("updating %s (WorkspaceID: %s ItemID: %s)", r.Name, plan.WorkspaceID.ValueString(), plan.ID.ValueString()))

		respUpdate, err := r.client.UpdateItem(ctx, plan.WorkspaceID.ValueString(), plan.ID.ValueString(), reqUpdatePlan.UpdateItemRequest, nil)
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
			return
		}

		plan.set(respUpdate.Item)

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}

	var reqUpdateDefinition requestUpdateFabricItemDefinitionDefinition

	doUpdateDefinition, diags := r.checkUpdateDefinition(ctx, plan, state, &reqUpdateDefinition)
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

func (r *ResourceFabricItemDefinition) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "DELETE", map[string]any{
		"state": req.State,
	})

	var state resourceFabricItemDefinitionModel

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

func (r *ResourceFabricItemDefinition) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

	var definition supertypes.MapNestedObjectValueOf[ResourceFabricItemDefinitionPartModel]
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("definition"), &definition)...); resp.Diagnostics.HasError() {
		return
	}

	state := resourceFabricItemDefinitionModel{
		baseFabricItemModel: baseFabricItemModel{
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

func (r *ResourceFabricItemDefinition) get(ctx context.Context, model *resourceFabricItemDefinitionModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s by ID: %s", r.Name, model.ID.ValueString()))

	respGet, err := r.client.GetItem(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		return diags
	}

	model.set(respGet.Item)

	return nil
}

func (r *ResourceFabricItemDefinition) checkUpdateItem(plan, state resourceFabricItemDefinitionModel, reqUpdatePlan *requestUpdateFabricItemDefinition) bool {
	var reqUpdateState requestUpdateFabricItemDefinition

	reqUpdatePlan.set(plan)
	reqUpdateState.set(state)

	return !reflect.DeepEqual(reqUpdatePlan.UpdateItemRequest, reqUpdateState.UpdateItemRequest)
}

func (r *ResourceFabricItemDefinition) checkUpdateDefinition(ctx context.Context, plan, state resourceFabricItemDefinitionModel, reqUpdate *requestUpdateFabricItemDefinitionDefinition) (bool, diag.Diagnostics) {
	if !plan.Definition.Equal(state.Definition) && plan.DefinitionUpdateEnabled.ValueBool() {
		if diags := reqUpdate.set(ctx, plan, r.DefinitionEmpty, r.DefinitionPathKeys); diags.HasError() {
			return false, diags
		}

		if len(reqUpdate.Definition.Parts) > 0 && !plan.Definition.Equal(state.Definition) {
			return true, nil
		}
	}

	return false, nil
}
