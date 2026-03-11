// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"

	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/params"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
)

type resourceFabricItemDefinitionModel struct {
	fabricItemModel

	Format                  types.String                                                             `tfsdk:"format"`
	DefinitionUpdateEnabled types.Bool                                                               `tfsdk:"definition_update_enabled"`
	Definition              supertypes.MapNestedObjectValueOf[resourceFabricItemDefinitionPartModel] `tfsdk:"definition"`
	Timeouts                timeouts.Value                                                           `tfsdk:"timeouts"`
}

type resourceFabricItemDefinitionPartModel struct {
	Source              types.String                                              `tfsdk:"source"`
	Parameters          supertypes.SetNestedObjectValueOf[params.ParametersModel] `tfsdk:"parameters"`
	ProcessingMode      types.String                                              `tfsdk:"processing_mode"`
	Tokens              supertypes.MapValueOf[types.String]                       `tfsdk:"tokens"`
	TokensDelimiter     types.String                                              `tfsdk:"tokens_delimiter"`
	SourceContentSha256 types.String                                              `tfsdk:"source_content_sha256"`
}

type fabricItemDefinition struct {
	fabcore.ItemDefinition
}

func (to *fabricItemDefinition) setFormat(v types.String, definitionFormats []DefinitionFormat) {
	if v.ValueString() != DefinitionFormatDefault && v.ValueString() != "" {
		apiFormat := getDefinitionFormatAPI(definitionFormats, v.ValueString())

		if apiFormat != "" {
			to.Format = &apiFormat
		}
	}
}

func (to *fabricItemDefinition) setParts(
	ctx context.Context,
	definition supertypes.MapNestedObjectValueOf[resourceFabricItemDefinitionPartModel],
	definitionEmpty string,
	definitionPaths []string,
	definitionUpdateEnabled types.Bool,
	update bool,
) diag.Diagnostics { //revive:disable-line:flag-parameter
	to.Parts = []fabcore.ItemDefinitionPart{}

	defParts, diags := definition.Get(ctx)
	if diags.HasError() {
		return diags
	}

	if (len(defParts) == 0) && len(definitionPaths) > 0 && update {
		contentB64, err := transforms.Base64Encode(definitionEmpty)
		if err != nil {
			diags.AddError(
				common.ErrorBase64EncodeHeader,
				err.Error(),
			)

			return diags
		}

		to.Parts = append(to.Parts, fabcore.ItemDefinitionPart{
			Path:        &definitionPaths[0],
			Payload:     &contentB64,
			PayloadType: azto.Ptr(fabcore.PayloadTypeInlineBase64),
		})

		return nil
	}

	for defPartKey, defPartValue := range defParts {
		if !update || (update && definitionUpdateEnabled.ValueBool()) {
			tokens, diags := defPartValue.Tokens.Get(ctx)
			if diags.HasError() {
				return diags
			}

			tokensValue := make(map[string]string)

			for k, v := range tokens {
				if !v.IsNull() && !v.IsUnknown() {
					tokensValue[k] = v.ValueString()
				}
			}

			parameters, diags := defPartValue.Parameters.Get(ctx)
			if diags.HasError() {
				return diags
			}

			payloadB64, _, diags := transforms.SourceFileToPayload(
				defPartValue.Source.ValueString(),
				defPartValue.ProcessingMode.ValueString(),
				tokensValue,
				parameters,
				defPartValue.TokensDelimiter.ValueString(),
			)
			if diags.HasError() {
				return diags
			}

			to.Parts = append(to.Parts, fabcore.ItemDefinitionPart{
				Path:        &defPartKey,
				Payload:     &payloadB64,
				PayloadType: azto.Ptr(fabcore.PayloadTypeInlineBase64),
			})
		}
	}

	return nil
}

type requestUpdateFabricItemDefinition struct {
	fabcore.UpdateItemDefinitionRequest
}

func (to *requestUpdateFabricItemDefinition) setDefinition(
	ctx context.Context,
	definition supertypes.MapNestedObjectValueOf[resourceFabricItemDefinitionPartModel],
	format types.String,
	definitionUpdateEnabled types.Bool,
	definitionEmpty string,
	definitionFormats []DefinitionFormat,
) diag.Diagnostics {
	var def fabricItemDefinition

	def.setFormat(format, definitionFormats)

	definitionPathKeys := GetDefinitionFormatPaths(definitionFormats, format.ValueString())

	if diags := def.setParts(ctx, definition, definitionEmpty, definitionPathKeys, definitionUpdateEnabled, true); diags.HasError() {
		return diags
	}

	to.Definition = &def.ItemDefinition

	return nil
}
