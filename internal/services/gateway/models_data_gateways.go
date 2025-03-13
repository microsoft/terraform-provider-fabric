// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type dataSourceGatewaysModel struct {
	Values   supertypes.ListNestedObjectValueOf[baseDataSourceGatewayModel] `tfsdk:"values"`
	Timeouts timeouts.Value                                                 `tfsdk:"timeouts"`
}

func (to *dataSourceGatewaysModel) setValues(ctx context.Context, from []fabcore.GatewayClassification) diag.Diagnostics {
	slice := make([]*baseDataSourceGatewayModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseDataSourceGatewayModel
		if diags := entityModel.set(ctx, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}
