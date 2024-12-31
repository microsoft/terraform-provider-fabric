// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
)

type dataSourceFabricItemDefinitionModel struct {
	fabricItemModel
	Format           types.String                                                               `tfsdk:"format"`
	OutputDefinition types.Bool                                                                 `tfsdk:"output_definition"`
	Definition       supertypes.MapNestedObjectValueOf[dataSourceFabricItemDefinitionPartModel] `tfsdk:"definition"`
	Timeouts         timeouts.Value                                                             `tfsdk:"timeouts"`
}

func (to *dataSourceFabricItemDefinitionModel) setDefinition(v supertypes.MapNestedObjectValueOf[dataSourceFabricItemDefinitionPartModel]) {
	to.Definition = v
}

func getDataSourceDefinitionModel(ctx context.Context, from fabcore.ItemDefinition) (supertypes.MapNestedObjectValueOf[dataSourceFabricItemDefinitionPartModel], diag.Diagnostics) {
	defParts := make(map[string]*dataSourceFabricItemDefinitionPartModel, len(from.Parts))

	result := supertypes.NewMapNestedObjectValueOfNull[dataSourceFabricItemDefinitionPartModel](ctx)

	for _, part := range from.Parts {
		newPart := &dataSourceFabricItemDefinitionPartModel{}

		if diags := newPart.set(*part.Payload); diags.HasError() {
			return result, diags
		}

		defParts[*part.Path] = newPart
	}

	if diags := result.Set(ctx, defParts); diags.HasError() {
		return result, diags
	}

	return result, nil
}

type dataSourceFabricItemDefinitionPartModel struct {
	Content types.String `tfsdk:"content"`
}

func (to *dataSourceFabricItemDefinitionPartModel) set(from string) diag.Diagnostics {
	content := from

	if diags := transforms.PayloadToGzip(&content); diags.HasError() {
		return diags
	}

	to.Content = types.StringPointerValue(&content)

	return nil
}
