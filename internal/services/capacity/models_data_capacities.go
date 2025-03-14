// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package capacity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type dataSourceCapacitiesModel struct {
	Values   supertypes.SetNestedObjectValueOf[baseCapacityModel] `tfsdk:"values"`
	Timeouts timeouts.Value                                       `tfsdk:"timeouts"`
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
