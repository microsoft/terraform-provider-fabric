// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelakeshortcut

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseShortcutModel struct {
	ID          customtypes.UUID                                  `tfsdk:"id"`
	Name        types.String                                      `tfsdk:"name"`
	Path        types.String                                      `tfsdk:"path"`
	Target      supertypes.SingleNestedObjectValueOf[targetModel] `tfsdk:"target"`
	WorkspaceID customtypes.UUID                                  `tfsdk:"workspace_id"`
	ItemID      customtypes.UUID                                  `tfsdk:"item_id"`
}

type targetModel struct {
	Onelake supertypes.SingleNestedObjectValueOf[oneLakeModel] `tfsdk:"onelake"`
	Type    types.String                                       `tfsdk:"type"`
}

type oneLakeModel struct {
	ItemId      customtypes.UUID `tfsdk:"item_id"`
	Path        types.String     `tfsdk:"path"`
	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
}

func (to *baseShortcutModel) set(ctx context.Context, workspaceID, itemID string, from fabcore.Shortcut) diag.Diagnostics {
	to.Name = types.StringPointerValue(from.Name)
	to.Path = types.StringPointerValue(from.Path)
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.ItemID = customtypes.NewUUIDValue(itemID)

	target := supertypes.NewSingleNestedObjectValueOfNull[targetModel](ctx)
	if from.Target != nil {
		targetModel := &targetModel{}
		if diags := targetModel.set(ctx, from.Target); diags.HasError() {
			return diags
		}
		if diags := target.Set(ctx, targetModel); diags.HasError() {
			return diags
		}
	}
	to.Target = target

	return nil
}

func (to *targetModel) set(ctx context.Context, from *fabcore.Target) diag.Diagnostics {
	to.Type = types.StringPointerValue((*string)(from.Type))
	onelake := supertypes.NewSingleNestedObjectValueOfNull[oneLakeModel](ctx)
	if from.OneLake != nil {
		onelakeModel := &oneLakeModel{}

		onelakeModel.set(from.OneLake)

		if diags := onelake.Set(ctx, onelakeModel); diags.HasError() {
			return diags
		}
	}

	to.Onelake = onelake

	return nil
}

func (to *oneLakeModel) set(from *fabcore.OneLake) {
	to.ItemId = customtypes.NewUUIDPointerValue(from.ItemID)
	to.Path = types.StringPointerValue(from.Path)
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
}

/*
DATA-SOURCE
*/

type dataSourceOnelakeShortcutModel struct {
	baseShortcutModel
	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceOnelakeShortcutsModel struct {
	Values   supertypes.SetNestedObjectValueOf[baseShortcutModel] `tfsdk:"values"`
	Timeouts timeoutsD.Value                                      `tfsdk:"timeouts"`
}

func (to *dataSourceOnelakeShortcutsModel) setValues(ctx context.Context, workspaceID, itemID string, from []fabcore.Shortcut) diag.Diagnostics {
	slice := make([]*baseShortcutModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseShortcutModel

		if diags := entityModel.set(ctx, workspaceID, itemID, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

type resourceOneLakeShortcutModel struct {
	baseShortcutModel
	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateOnelakeShortcut struct {
	fabcore.CreateShortcutRequest
}

func (to *requestCreateOnelakeShortcut) set(ctx context.Context, from resourceOneLakeShortcutModel) diag.Diagnostics {
	to.Name = from.Name.ValueStringPointer()
	to.Path = from.Path.ValueStringPointer()

	target, diags := from.Target.Get(ctx)
	if diags.HasError() {
		return diags
	}
	to.Target = &fabcore.CreatableShortcutTarget{
		OneLake: func() *fabcore.OneLake {
			onelake, diags := target.Onelake.Get(ctx)
			if diags.HasError() || onelake == nil {
				return nil
			}
			return &fabcore.OneLake{
				ItemID:      onelake.ItemId.ValueStringPointer(),
				Path:        onelake.Path.ValueStringPointer(),
				WorkspaceID: onelake.WorkspaceID.ValueStringPointer(),
			}
		}(),
	}

	return nil
}
