// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

type DataSourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop any] struct {
	FabricItemPropertiesModel[Ttfprop, Titemprop]
	Format           types.String                                                               `tfsdk:"format"`
	OutputDefinition types.Bool                                                                 `tfsdk:"output_definition"`
	Definition       supertypes.MapNestedObjectValueOf[DataSourceFabricItemDefinitionPartModel] `tfsdk:"definition"`
	Timeouts         timeouts.Value                                                             `tfsdk:"timeouts"`
}

func (to *DataSourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]) setDefinition(ctx context.Context, from fabcore.ItemDefinition) diag.Diagnostics {
	defParts := make(map[string]*DataSourceFabricItemDefinitionPartModel, len(from.Parts))

	for _, part := range from.Parts {
		newPart := &DataSourceFabricItemDefinitionPartModel{}

		if diags := newPart.Set(*part.Payload); diags.HasError() {
			return diags
		}

		defParts[*part.Path] = newPart
	}

	return to.Definition.Set(ctx, defParts)
}
