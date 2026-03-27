// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceFabricItemListModel struct {
	WorkspaceID      customtypes.UUID                                            `tfsdk:"workspace_id"`
	ID               customtypes.UUID                                            `tfsdk:"id"`
	DisplayName      types.String                                                `tfsdk:"display_name"`
	Description      types.String                                                `tfsdk:"description"`
	FolderID         customtypes.UUID                                            `tfsdk:"folder_id"`
	SensitivityLabel supertypes.SingleNestedObjectValueOf[sensitivityLabelModel] `tfsdk:"sensitivity_label"`
}

type dataSourceFabricItemsModel struct {
	WorkspaceID customtypes.UUID                                                 `tfsdk:"workspace_id"`
	Values      supertypes.SetNestedObjectValueOf[dataSourceFabricItemListModel] `tfsdk:"values"`
	Timeouts    timeouts.Value                                                   `tfsdk:"timeouts"`
}

func (to *dataSourceFabricItemsModel) setValues(ctx context.Context, from []fabcore.Item) diag.Diagnostics {
	slice := make([]*dataSourceFabricItemListModel, 0, len(from))

	for _, entity := range from {
		sensitivityLabel := supertypes.NewSingleNestedObjectValueOfNull[sensitivityLabelModel](ctx)

		if entity.SensitivityLabel != nil && entity.SensitivityLabel.ID != nil {
			sensitivityLabelModel := &sensitivityLabelModel{}
			sensitivityLabelModel.set(*entity.SensitivityLabel)

			if diags := sensitivityLabel.Set(ctx, sensitivityLabelModel); diags.HasError() {
				return diags
			}
		}

		entityModel := &dataSourceFabricItemListModel{
			WorkspaceID:      customtypes.NewUUIDPointerValue(entity.WorkspaceID),
			ID:               customtypes.NewUUIDPointerValue(entity.ID),
			DisplayName:      types.StringPointerValue(entity.DisplayName),
			Description:      types.StringPointerValue(entity.Description),
			FolderID:         customtypes.NewUUIDPointerValue(entity.FolderID),
			SensitivityLabel: sensitivityLabel,
		}

		slice = append(slice, entityModel)
	}

	return to.Values.Set(ctx, slice)
}
