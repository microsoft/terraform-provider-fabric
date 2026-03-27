// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type DataSourceFabricItemListPropertiesModel[Ttfprop, Titemprop any] struct { //revive:disable-line:exported
	WorkspaceID      customtypes.UUID                                            `tfsdk:"workspace_id"`
	ID               customtypes.UUID                                            `tfsdk:"id"`
	DisplayName      types.String                                                `tfsdk:"display_name"`
	Description      types.String                                                `tfsdk:"description"`
	FolderID         customtypes.UUID                                            `tfsdk:"folder_id"`
	SensitivityLabel supertypes.SingleNestedObjectValueOf[sensitivityLabelModel] `tfsdk:"sensitivity_label"`
	Properties       supertypes.SingleNestedObjectValueOf[Ttfprop]               `tfsdk:"properties"`
}

type DataSourceFabricItemsPropertiesModel[Ttfprop, Titemprop any] struct {
	WorkspaceID customtypes.UUID                                                                               `tfsdk:"workspace_id"`
	Values      supertypes.SetNestedObjectValueOf[DataSourceFabricItemListPropertiesModel[Ttfprop, Titemprop]] `tfsdk:"values"`
	Timeouts    timeouts.Value                                                                                 `tfsdk:"timeouts"`
}

func (to *DataSourceFabricItemsPropertiesModel[Ttfprop, Titemprop]) setValues(
	ctx context.Context,
	from []FabricItemProperties[Titemprop],
	propertiesSetter func(ctx context.Context, from *Titemprop, to *FabricItemPropertiesModel[Ttfprop, Titemprop]) diag.Diagnostics,
) diag.Diagnostics {
	slice := make([]*DataSourceFabricItemListPropertiesModel[Ttfprop, Titemprop], 0, len(from))

	for _, entity := range from {
		sensitivityLabel := supertypes.NewSingleNestedObjectValueOfNull[sensitivityLabelModel](ctx)

		if entity.SensitivityLabel != nil && entity.SensitivityLabel.ID != nil {
			sensitivityLabelModel := &sensitivityLabelModel{}
			sensitivityLabelModel.set(*entity.SensitivityLabel)

			if diags := sensitivityLabel.Set(ctx, sensitivityLabelModel); diags.HasError() {
				return diags
			}
		}

		entityModel := &DataSourceFabricItemListPropertiesModel[Ttfprop, Titemprop]{
			WorkspaceID:      customtypes.NewUUIDPointerValue(entity.WorkspaceID),
			ID:               customtypes.NewUUIDPointerValue(entity.ID),
			DisplayName:      types.StringPointerValue(entity.DisplayName),
			Description:      types.StringPointerValue(entity.Description),
			FolderID:         customtypes.NewUUIDPointerValue(entity.FolderID),
			SensitivityLabel: sensitivityLabel,
		}

		var propsModel FabricItemPropertiesModel[Ttfprop, Titemprop]

		diags := propertiesSetter(ctx, entity.Properties, &propsModel)
		if diags.HasError() {
			return diags
		}

		entityModel.Properties = propsModel.Properties

		slice = append(slice, entityModel)
	}

	return to.Values.Set(ctx, slice)
}
