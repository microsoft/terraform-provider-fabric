// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithModifyPlan       = (*ResourceFabricItemConfigDefinitionProperties[struct{}, struct{}, struct{}, struct{}])(nil)
	_ resource.ResourceWithConfigValidators = (*ResourceFabricItemConfigDefinitionProperties[struct{}, struct{}, struct{}, struct{}])(nil)
	_ resource.ResourceWithConfigure        = (*ResourceFabricItemConfigDefinitionProperties[struct{}, struct{}, struct{}, struct{}])(nil)
	_ resource.ResourceWithImportState      = (*ResourceFabricItemConfigDefinitionProperties[struct{}, struct{}, struct{}, struct{}])(nil)
)

type ResourceFabricItemConfigDefinitionProperties[Ttfprop, Titemprop, Ttfconfig, Titemconfig any] struct {
	ResourceFabricItemDefinition

	ConfigRequired             bool
	ConfigOrDefinitionRequired bool
	ConfigAttributes           map[string]schema.Attribute
	PropertiesAttributes       map[string]schema.Attribute
	PropertiesSetter           func(ctx context.Context, from *Titemprop, to *ResourceFabricItemConfigDefinitionPropertiesModel[Ttfprop, Titemprop, Ttfconfig, Titemconfig]) diag.Diagnostics
	CreationPayloadSetter      func(ctx context.Context, from Ttfconfig) (*Titemconfig, diag.Diagnostics)
	ItemGetter                 func(ctx context.Context, fabricClient fabric.Client, model ResourceFabricItemConfigDefinitionPropertiesModel[Ttfprop, Titemprop, Ttfconfig, Titemconfig], fabricItem *FabricItemProperties[Titemprop]) error
}

func NewResourceFabricItemConfigDefinitionProperties[Ttfprop, Titemprop, Ttfconfig, Titemconfig any](
	config ResourceFabricItemConfigDefinitionProperties[Ttfprop, Titemprop, Ttfconfig, Titemconfig],
) resource.Resource {
	return &config
}

func (r *ResourceFabricItemConfigDefinitionProperties[Ttfprop, Titemprop, Ttfconfig, Titemconfig]) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *ResourceFabricItemConfigDefinitionProperties[Ttfprop, Titemprop, Ttfconfig, Titemconfig]) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	tflog.Debug(ctx, "MODIFY PLAN", map[string]any{
		"action": "start",
	})

	if !req.State.Raw.IsNull() && !req.Plan.Raw.IsNull() {
		var plan, state ResourceFabricItemConfigDefinitionPropertiesModel[Ttfprop, Titemprop, Ttfconfig, Titemconfig]

		resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

		if resp.Diagnostics.HasError() {
			return
		}

		var reqUpdateDefinition requestUpdateFabricItemDefinition

		doUpdateDefinition, diags := fabricItemCheckUpdateDefinition(
			ctx,
			plan.Definition,
			state.Definition,
			plan.Format,
			plan.DefinitionUpdateEnabled,
			r.DefinitionEmpty,
			r.DefinitionFormats,
			&reqUpdateDefinition,
		)
		if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
			return
		}

		if doUpdateDefinition {
			resp.Diagnostics.AddWarning(
				common.WarningItemDefinitionUpdateHeader,
				fmt.Sprintf(common.WarningItemDefinitionUpdateDetails, r.TypeInfo.Name),
			)
		}
	}

	tflog.Debug(ctx, "MODIFY PLAN", map[string]any{
		"action": "end",
	})
}

func (r *ResourceFabricItemConfigDefinitionProperties[Ttfprop, Titemprop, Ttfconfig, Titemconfig]) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = getResourceFabricItemConfigDefinitionPropertiesSchema(ctx, *r)
}

func (r *ResourceFabricItemConfigDefinitionProperties[Ttfprop, Titemprop, Ttfconfig, Titemconfig]) ConfigValidators(
	_ context.Context,
) []resource.ConfigValidator {
	result := []resource.ConfigValidator{}

	if r.ConfigOrDefinitionRequired {
		result = append(result, resourcevalidator.ExactlyOneOf(
			path.MatchRoot("configuration"),
			path.MatchRoot("definition"),
		))
	}

	result = append(result, resourcevalidator.Conflicting(
		path.MatchRoot("configuration"),
		path.MatchRoot("definition"),
	))

	return result
}

func (r *ResourceFabricItemConfigDefinitionProperties[Ttfprop, Titemprop, Ttfconfig, Titemconfig]) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
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

	if resp.Diagnostics.Append(IsPreviewMode(r.TypeInfo.Name, r.TypeInfo.IsPreview, r.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourceFabricItemConfigDefinitionProperties[Ttfprop, Titemprop, Ttfconfig, Titemconfig]) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan, config ResourceFabricItemConfigDefinitionPropertiesModel[Ttfprop, Titemprop, Ttfconfig, Titemconfig]

	if resp.Diagnostics.Append(req.Config.Get(ctx, &config)...); resp.Diagnostics.HasError() {
		return
	}

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
	reqCreate.setType(r.FabricItemType)

	if resp.Diagnostics.Append(reqCreate.setDefinition(ctx, plan.Definition, plan.Format, plan.DefinitionUpdateEnabled, r.DefinitionFormats)...); resp.Diagnostics.HasError() {
		return
	}

	creationPayload, diags := getCreationPayload(ctx, config.Configuration, r.CreationPayloadSetter)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	reqCreate.setCreationPayload(creationPayload)

	respCreate, err := CreateItem(ctx, r.client, plan.WorkspaceID.ValueString(), reqCreate.CreateItemRequest)
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

func (r *ResourceFabricItemConfigDefinitionProperties[Ttfprop, Titemprop, Ttfconfig, Titemconfig]) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state ResourceFabricItemConfigDefinitionPropertiesModel[Ttfprop, Titemprop, Ttfconfig, Titemconfig]

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

func (r *ResourceFabricItemConfigDefinitionProperties[Ttfprop, Titemprop, Ttfconfig, Titemconfig]) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	var plan, state ResourceFabricItemConfigDefinitionPropertiesModel[Ttfprop, Titemprop, Ttfconfig, Titemconfig]

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
		tflog.Trace(ctx, fmt.Sprintf("updating %s (WorkspaceID: %s ItemID: %s)", r.TypeInfo.Name, plan.WorkspaceID.ValueString(), plan.ID.ValueString()))

		_, err := UpdateItem(ctx, r.client, plan.WorkspaceID.ValueString(), plan.ID.ValueString(), reqUpdatePlan.UpdateItemRequest)
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
			return
		}

		if resp.Diagnostics.Append(r.get(ctx, &plan)...); resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	}

	var reqUpdateDefinition requestUpdateFabricItemDefinition

	doUpdateDefinition, diags := fabricItemCheckUpdateDefinition(
		ctx,
		plan.Definition,
		state.Definition,
		plan.Format,
		plan.DefinitionUpdateEnabled,
		r.DefinitionEmpty,
		r.DefinitionFormats,
		&reqUpdateDefinition,
	)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	if doUpdateDefinition {
		tflog.Trace(ctx, fmt.Sprintf("updating %s definition", r.TypeInfo.Name))

		_, err := r.client.UpdateItemDefinition(ctx, plan.WorkspaceID.ValueString(), plan.ID.ValueString(), reqUpdateDefinition.UpdateItemDefinitionRequest, nil)
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
			return
		}

		if resp.Diagnostics.Append(r.get(ctx, &plan)...); resp.Diagnostics.HasError() {
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

func (r *ResourceFabricItemConfigDefinitionProperties[Ttfprop, Titemprop, Ttfconfig, Titemconfig]) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	var state ResourceFabricItemConfigDefinitionPropertiesModel[Ttfprop, Titemprop, Ttfconfig, Titemconfig]

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

func (r *ResourceFabricItemConfigDefinitionProperties[Ttfprop, Titemprop, Ttfconfig, Titemconfig]) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
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
				fmt.Sprintf("WorkspaceID/%sID", string(r.FabricItemType)),
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

	var configuration supertypes.SingleNestedObjectValueOf[Ttfconfig]
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("configuration"), &configuration)...); resp.Diagnostics.HasError() {
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

	state := ResourceFabricItemConfigDefinitionPropertiesModel[Ttfprop, Titemprop, Ttfconfig, Titemconfig]{
		FabricItemPropertiesModel: FabricItemPropertiesModel[Ttfprop, Titemprop]{
			ID:          uuidFabricItemID,
			WorkspaceID: uuidWorkspaceID,
		},
		Configuration:           configuration,
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

func (r *ResourceFabricItemConfigDefinitionProperties[Ttfprop, Titemprop, Ttfconfig, Titemconfig]) get(
	ctx context.Context,
	model *ResourceFabricItemConfigDefinitionPropertiesModel[Ttfprop, Titemprop, Ttfconfig, Titemconfig],
) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s by ID: %s", r.TypeInfo.Name, model.ID.ValueString()))

	var fabricItem FabricItemProperties[Titemprop]

	err := r.ItemGetter(ctx, *r.pConfigData.FabricClient, *model, &fabricItem)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		return diags
	}

	model.set(fabricItem)

	return r.PropertiesSetter(ctx, fabricItem.Properties, model)
}
