// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"

	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
)

type ResourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop any] struct {
	// FabricItemPropertiesModel[Titemprop, Titemprop]
	baseFabricItemModel
	Properties              supertypes.SingleNestedObjectValueOf[Ttfprop]                            `tfsdk:"properties"`
	Format                  types.String                                                             `tfsdk:"format"`
	DefinitionUpdateEnabled types.Bool                                                               `tfsdk:"definition_update_enabled"`
	Definition              supertypes.MapNestedObjectValueOf[ResourceFabricItemDefinitionPartModel] `tfsdk:"definition"`
	Timeouts                timeouts.Value                                                           `tfsdk:"timeouts"`
}

func (to *ResourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]) set(from FabricItemProperties[Titemprop]) { //revive:disable-line:confusing-naming
	to.WorkspaceID = customtypes.NewUUIDPointerValue(from.WorkspaceID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
}

type FabricItemDefinitionProperties[Ttfprop, Titemprop any] struct { //revive:disable-line:exported
	fabcore.ItemDefinition
}

func (to *FabricItemDefinitionProperties[Ttfprop, Titemprop]) set(ctx context.Context, from ResourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop], update bool, definitionEmpty string, definitionPaths []string) diag.Diagnostics { //revive:disable-line:flag-parameter,confusing-naming
	if from.Format.ValueString() != DefinitionFormatNotApplicable {
		to.Format = from.Format.ValueStringPointer()
	}

	to.Parts = []fabcore.ItemDefinitionPart{}

	defParts, diags := from.Definition.Get(ctx)
	if diags.HasError() {
		return diags
	}

	if (len(defParts) == 0) && len(definitionPaths) > 0 && update {
		content := definitionEmpty

		if err := transforms.Base64Encode(&content); err != nil {
			diags.AddError(
				common.ErrorBase64EncodeHeader,
				err.Error(),
			)

			return diags
		}

		to.Parts = append(to.Parts, fabcore.ItemDefinitionPart{
			Path:        azto.Ptr(definitionPaths[0]),
			Payload:     &content,
			PayloadType: azto.Ptr(fabcore.PayloadTypeInlineBase64),
		})

		return nil
	}

	for defPartKey, defPartValue := range defParts {
		if !update || (update && from.DefinitionUpdateEnabled.ValueBool()) {
			payloadB64, _, diags := transforms.SourceFileToPayload(ctx, defPartValue.Source, defPartValue.Tokens)
			if diags.HasError() {
				return diags
			}

			to.Parts = append(to.Parts, fabcore.ItemDefinitionPart{
				Path:        azto.Ptr(defPartKey),
				Payload:     payloadB64,
				PayloadType: azto.Ptr(fabcore.PayloadTypeInlineBase64),
			})
		}
	}

	return nil
}

type requestCreateFabricItemDefinitionProperties[Ttfprop, Titemprop any] struct {
	fabcore.CreateItemRequest
}

func (to *requestCreateFabricItemDefinitionProperties[Ttfprop, Titemprop]) set(ctx context.Context, from ResourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop], itemType fabcore.ItemType) diag.Diagnostics { //revive:disable-line:confusing-naming
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
	to.Type = azto.Ptr(itemType)

	if !from.Definition.IsNull() && !from.Definition.IsUnknown() {
		var def FabricItemDefinitionProperties[Ttfprop, Titemprop]

		if diags := def.set(ctx, from, false, "", []string{}); diags.HasError() {
			return diags
		}

		to.Definition = &def.ItemDefinition
	}

	return nil
}

type requestUpdateFabricItemDefinitionProperties[Ttfprop, Titemprop any] struct {
	fabcore.UpdateItemRequest
}

func (to *requestUpdateFabricItemDefinitionProperties[Ttfprop, Titemprop]) set(from ResourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]) { //revive:disable-line:confusing-naming
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
}

type requestUpdateFabricItemDefinitionPropertiesDefinition[Ttfprop, Titemprop any] struct {
	fabcore.UpdateItemDefinitionRequest
}

func (to *requestUpdateFabricItemDefinitionPropertiesDefinition[Ttfprop, Titemprop]) set(ctx context.Context, from ResourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop], definitionEmpty string, definitionPaths []string) diag.Diagnostics { //revive:disable-line:confusing-naming
	var def FabricItemDefinitionProperties[Ttfprop, Titemprop]

	if diags := def.set(ctx, from, true, definitionEmpty, definitionPaths); diags.HasError() {
		return diags
	}

	to.Definition = &def.ItemDefinition

	return nil
}
