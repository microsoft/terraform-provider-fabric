// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceFabricItemsModel struct {
	WorkspaceID customtypes.UUID                                                 `tfsdk:"workspace_id"`
	Values      supertypes.SetNestedObjectValueOf[dataSourceFabricItemBaseModel] `tfsdk:"values"`
	Timeouts    timeouts.Value                                                   `tfsdk:"timeouts"`
}

func (to *dataSourceFabricItemsModel) setValues(ctx context.Context, from []fabcore.Item) diag.Diagnostics {
	slice := make([]*dataSourceFabricItemBaseModel, 0, len(from))

	for _, entity := range from {
		var entityModel dataSourceFabricItemBaseModel
		entityModel.set(ctx, entity)

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}
