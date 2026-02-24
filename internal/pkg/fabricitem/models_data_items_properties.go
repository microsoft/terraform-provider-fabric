// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type DataSourceFabricItemsPropertiesModel[Ttfprop, Titemprop any] struct {
	WorkspaceID customtypes.UUID                                                                 `tfsdk:"workspace_id"`
	Values      supertypes.SetNestedObjectValueOf[FabricItemPropertiesModel[Ttfprop, Titemprop]] `tfsdk:"values"`
	Timeouts    timeouts.Value                                                                   `tfsdk:"timeouts"`
}

func (to *DataSourceFabricItemsPropertiesModel[Ttfprop, Titemprop]) setValues(
	ctx context.Context,
	from []FabricItemProperties[Titemprop],
	propertiesSetter func(ctx context.Context, from *Titemprop, to *FabricItemPropertiesModel[Ttfprop, Titemprop]) diag.Diagnostics,
) diag.Diagnostics {
	slice := make([]*FabricItemPropertiesModel[Ttfprop, Titemprop], 0, len(from))

	for _, entity := range from {
		var entityModel FabricItemPropertiesModel[Ttfprop, Titemprop]
		entityModel.set(entity)

		diags := propertiesSetter(ctx, entity.Properties, &entityModel)
		if diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}
