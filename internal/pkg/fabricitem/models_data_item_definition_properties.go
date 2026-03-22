// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type DataSourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop any] struct {
	FabricItemPropertiesModel[Ttfprop, Titemprop]

	SensitivityLabel supertypes.SingleNestedObjectValueOf[sensitivityLabelModel]                `tfsdk:"sensitivity_label"`
	Format           types.String                                                               `tfsdk:"format"`
	OutputDefinition types.Bool                                                                 `tfsdk:"output_definition"`
	Definition       supertypes.MapNestedObjectValueOf[dataSourceFabricItemDefinitionPartModel] `tfsdk:"definition"`
	Timeouts         timeouts.Value                                                             `tfsdk:"timeouts"`
}

func (to *DataSourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]) set(ctx context.Context, from FabricItemProperties[Titemprop]) diag.Diagnostics {
	to.FabricItemPropertiesModel.set(from)

	sl, diags := newSensitivityLabelFromAPI(ctx, from.SensitivityLabel)
	if diags.HasError() {
		return diags
	}

	to.SensitivityLabel = sl

	return nil
}

func (to *DataSourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]) setDefinition(v supertypes.MapNestedObjectValueOf[dataSourceFabricItemDefinitionPartModel]) {
	to.Definition = v
}
