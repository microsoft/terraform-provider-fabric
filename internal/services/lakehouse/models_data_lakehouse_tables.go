// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fablakehouse "github.com/microsoft/fabric-sdk-go/fabric/lakehouse"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceLakehouseTablesModel struct {
	LakehouseID customtypes.UUID                                        `tfsdk:"lakehouse_id"`
	WorkspaceID customtypes.UUID                                        `tfsdk:"workspace_id"`
	Values      supertypes.ListNestedObjectValueOf[lakehouseTableModel] `tfsdk:"values"`
	Timeouts    timeouts.Value                                          `tfsdk:"timeouts"`
}

func (to *dataSourceLakehouseTablesModel) setValues(ctx context.Context, from []fablakehouse.Table) diag.Diagnostics {
	slice := make([]*lakehouseTableModel, 0, len(from))

	for _, entity := range from {
		var entityModel lakehouseTableModel
		entityModel.set(entity)
		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}
