// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type resourceFabricItemModel struct {
	fabricItemModel

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type requestCreateFabricItem struct {
	fabcore.CreateItemRequest
}

type DBPointer[T any] interface {
	*T
	fabcore.CreateItemRequest
}

func (to *requestCreateFabricItem) setDisplayName(v types.String) {
	to.DisplayName = v.ValueStringPointer()
}

func (to *requestCreateFabricItem) setDescription(v types.String) {
	to.Description = v.ValueStringPointer()
}

func (to *requestCreateFabricItem) setType(v fabcore.ItemType) {
	to.Type = &v
}

func (to *requestCreateFabricItem) setDefinition(
	ctx context.Context,
	definition supertypes.MapNestedObjectValueOf[resourceFabricItemDefinitionPartModel],
	format types.String,
	definitionUpdateEnabled types.Bool,
	definitionFormats []DefinitionFormat,
) diag.Diagnostics {
	if !definition.IsNull() && !definition.IsUnknown() {
		var def fabricItemDefinition

		def.setFormat(format, definitionFormats)

		if diags := def.setParts(ctx, definition, "", []string{}, definitionUpdateEnabled, false); diags.HasError() {
			return diags
		}

		to.Definition = &def.ItemDefinition
	}

	return nil
}

func (to *requestCreateFabricItem) setCreationPayload(v any) {
	if v != nil {
		to.CreationPayload = v
	}
}

func getCreationPayload[Ttfconfig, Titemconfig any](
	ctx context.Context,
	configuration supertypes.SingleNestedObjectValueOf[Ttfconfig],
	creationPayloadSetter func(ctx context.Context, from Ttfconfig) (*Titemconfig, diag.Diagnostics),
) (*Titemconfig, diag.Diagnostics) {
	if !configuration.IsNull() && !configuration.IsUnknown() {
		config, diags := configuration.Get(ctx)
		if diags.HasError() {
			return nil, diags
		}

		creationPayload, diags := creationPayloadSetter(ctx, *config)
		if diags.HasError() {
			return nil, diags
		}

		return creationPayload, nil
	}

	return nil, nil
}

type requestUpdateFabricItem struct {
	fabcore.UpdateItemRequest
}

func (to *requestUpdateFabricItem) setDisplayName(v types.String) {
	to.DisplayName = v.ValueStringPointer()
}

func (to *requestUpdateFabricItem) setDescription(v types.String) {
	to.Description = v.ValueStringPointer()
}

func fabricItemCheckUpdate(planDisplayName, planDescription, stateDisplayName, stateDescription types.String, reqUpdatePlan *requestUpdateFabricItem) bool {
	var reqUpdateState requestUpdateFabricItem

	reqUpdatePlan.setDisplayName(planDisplayName)
	reqUpdatePlan.setDescription(planDescription)

	reqUpdateState.setDisplayName(stateDisplayName)
	reqUpdateState.setDescription(stateDescription)

	return !reflect.DeepEqual(reqUpdatePlan.UpdateItemRequest, reqUpdateState.UpdateItemRequest)
}

func fabricItemCheckUpdateDefinition(
	ctx context.Context,
	planDefinition, stateDefinition supertypes.MapNestedObjectValueOf[resourceFabricItemDefinitionPartModel],
	planFormat types.String,
	planDefinitionUpdateEnabled types.Bool,
	definitionEmpty string,
	definitionFormats []DefinitionFormat,
	reqUpdate *requestUpdateFabricItemDefinition,
) (bool, diag.Diagnostics) {
	if !planDefinition.Equal(stateDefinition) && planDefinitionUpdateEnabled.ValueBool() {
		if diags := reqUpdate.setDefinition(ctx, planDefinition, planFormat, planDefinitionUpdateEnabled, definitionEmpty, definitionFormats); diags.HasError() {
			return false, diags
		}

		if len(reqUpdate.Definition.Parts) > 0 && !planDefinition.Equal(stateDefinition) {
			return true, nil
		}
	}

	return false, nil
}
