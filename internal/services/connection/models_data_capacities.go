// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

type dataSourceConnectionsModel struct {
	Values   supertypes.ListNestedObjectValueOf[baseConnectionModel] `tfsdk:"values"`
	Timeouts timeouts.Value                                          `tfsdk:"timeouts"`
}

func (to *dataSourceConnectionsModel) setValues(ctx context.Context, from []fabcore.Connection) diag.Diagnostics {
	slice := make([]*baseConnectionModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseConnectionModel
		diags := entityModel.set(ctx, entity)
		if diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}
