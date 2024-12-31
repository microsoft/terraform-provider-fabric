// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DataSourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop any] struct {
	FabricItemPropertiesModel[Ttfprop, Titemprop]
	Format           types.String                                                               `tfsdk:"format"`
	OutputDefinition types.Bool                                                                 `tfsdk:"output_definition"`
	Definition       supertypes.MapNestedObjectValueOf[dataSourceFabricItemDefinitionPartModel] `tfsdk:"definition"`
	Timeouts         timeouts.Value                                                             `tfsdk:"timeouts"`
}

func (to *DataSourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]) setDefinition(v supertypes.MapNestedObjectValueOf[dataSourceFabricItemDefinitionPartModel]) {
	to.Definition = v
}
