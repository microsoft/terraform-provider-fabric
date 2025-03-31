// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehousetable

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	//revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fablakehouse "github.com/microsoft/fabric-sdk-go/fabric/lakehouse"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseLakehouseTableModel struct {
	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
	LakehouseID customtypes.UUID `tfsdk:"lakehouse_id"`
	Name        types.String     `tfsdk:"name"`
	Location    types.String     `tfsdk:"location"`
	Type        types.String     `tfsdk:"type"`
	Format      types.String     `tfsdk:"format"`
}

func (to *baseLakehouseTableModel) set(workspaceID, lakehouseID string, from fablakehouse.Table) {
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.LakehouseID = customtypes.NewUUIDValue(lakehouseID)

	to.Name = types.StringPointerValue(from.Name)
	to.Location = types.StringPointerValue(from.Location)
	to.Type = types.StringPointerValue((*string)(from.Type))
	to.Format = types.StringPointerValue(from.Format)
}

/*
DATA-SOURCE
*/

type dataSourceLakehouseTableModel struct {
	baseLakehouseTableModel
	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceLakehouseTablesModel struct {
	WorkspaceID customtypes.UUID                                           `tfsdk:"workspace_id"`
	LakehouseID customtypes.UUID                                           `tfsdk:"lakehouse_id"`
	Values      supertypes.SetNestedObjectValueOf[baseLakehouseTableModel] `tfsdk:"values"`
	Timeouts    timeoutsD.Value                                            `tfsdk:"timeouts"`
}

func (to *dataSourceLakehouseTablesModel) setValues(ctx context.Context, workspaceID, lakehouseID string, from []fablakehouse.Table) diag.Diagnostics {
	slice := make([]*baseLakehouseTableModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseLakehouseTableModel
		entityModel.set(workspaceID, lakehouseID, entity)
		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}
