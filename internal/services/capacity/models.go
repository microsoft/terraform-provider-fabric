// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package capacity

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseCapacityModel struct {
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
	SKU         types.String     `tfsdk:"sku"`
	Region      types.String     `tfsdk:"region"`
	State       types.String     `tfsdk:"state"`
}

func (to *baseCapacityModel) set(from fabcore.Capacity) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.SKU = types.StringPointerValue(from.SKU)
	to.Region = types.StringPointerValue(from.Region)
	to.State = types.StringPointerValue((*string)(from.State))
}

/*
DATA-SOURCE
*/

type dataSourceCapacityModel struct {
	baseCapacityModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceCapacitiesModel struct {
	Values   supertypes.SetNestedObjectValueOf[baseCapacityModel] `tfsdk:"values"`
	Timeouts timeoutsD.Value                                      `tfsdk:"timeouts"`
}

func (to *dataSourceCapacitiesModel) setValues(ctx context.Context, from []fabcore.Capacity) diag.Diagnostics {
	slice := make([]*baseCapacityModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseCapacityModel
		entityModel.set(entity)
		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}
