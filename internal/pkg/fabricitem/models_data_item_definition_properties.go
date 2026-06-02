// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
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
