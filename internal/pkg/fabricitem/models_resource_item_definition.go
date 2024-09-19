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
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
)

type resourceFabricItemDefinitionModel struct {
	baseFabricItemModel
	Format                  types.String                                                             `tfsdk:"format"`
	DefinitionUpdateEnabled types.Bool                                                               `tfsdk:"definition_update_enabled"`
	Definition              supertypes.MapNestedObjectValueOf[ResourceFabricItemDefinitionPartModel] `tfsdk:"definition"`
	Timeouts                timeouts.Value                                                           `tfsdk:"timeouts"`
}

type ResourceFabricItemDefinitionPartModel struct {
	Source              types.String                  `tfsdk:"source"`
	Tokens              supertypes.MapValueOf[string] `tfsdk:"tokens"`
	SourceContentSha256 types.String                  `tfsdk:"source_content_sha256"`
}

type fabricItemDefinition struct {
	fabcore.ItemDefinition
}

func (to *fabricItemDefinition) set(ctx context.Context, from resourceFabricItemDefinitionModel, update bool, definitionEmpty string, definitionPaths []string) diag.Diagnostics { //revive:disable-line:flag-parameter
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

type requestCreateFabricItemDefinition struct {
	fabcore.CreateItemRequest
}

func (to *requestCreateFabricItemDefinition) set(ctx context.Context, from resourceFabricItemDefinitionModel, itemType fabcore.ItemType) diag.Diagnostics {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
	to.Type = azto.Ptr(itemType)

	if !from.Definition.IsNull() && !from.Definition.IsUnknown() {
		var def fabricItemDefinition

		if diags := def.set(ctx, from, false, "", []string{}); diags.HasError() {
			return diags
		}

		to.Definition = &def.ItemDefinition
	}

	return nil
}

type requestUpdateFabricItemDefinition struct {
	fabcore.UpdateItemRequest
}

func (to *requestUpdateFabricItemDefinition) set(from resourceFabricItemDefinitionModel) {
	to.DisplayName = from.DisplayName.ValueStringPointer()
	to.Description = from.Description.ValueStringPointer()
}

type requestUpdateFabricItemDefinitionDefinition struct {
	fabcore.UpdateItemDefinitionRequest
}

func (to *requestUpdateFabricItemDefinitionDefinition) set(ctx context.Context, from resourceFabricItemDefinitionModel, definitionEmpty string, definitionPaths []string) diag.Diagnostics {
	var def fabricItemDefinition

	if diags := def.set(ctx, from, true, definitionEmpty, definitionPaths); diags.HasError() {
		return diags
	}

	to.Definition = &def.ItemDefinition

	return nil
}
