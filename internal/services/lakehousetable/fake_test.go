// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package lakehousetable_test

import (
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fablakehouse "github.com/microsoft/fabric-sdk-go/fabric/lakehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func fakeLakehouseTablesFunc(
	lakehouseTables fablakehouse.Tables,
) func(workspaceID, lakehouseID string, options *fablakehouse.TablesClientListTablesOptions) (resp azfake.PagerResponder[fablakehouse.TablesClientListTablesResponse]) {
	return func(_, _ string, _ *fablakehouse.TablesClientListTablesOptions) (resp azfake.PagerResponder[fablakehouse.TablesClientListTablesResponse]) {
		resp = azfake.PagerResponder[fablakehouse.TablesClientListTablesResponse]{}
		resp.AddPage(http.StatusOK, fablakehouse.TablesClientListTablesResponse{Tables: lakehouseTables}, nil)

		return resp
	}
}

func NewRandomLakehouseTables(lakehouseID string) fablakehouse.Tables {
	table0Name := testhelp.RandomName()
	table1Name := testhelp.RandomName()
	table2Name := testhelp.RandomName()

	return fablakehouse.Tables{
		Data: []fablakehouse.Table{
			{
				Name:     new(table0Name),
				Type:     azto.Ptr(fablakehouse.TableTypeExternal),
				Format:   new("Delta"),
				Location: new("abfss://" + testhelp.RandomUUID() + "@onelake.dfs.fabric.microsoft.com/" + lakehouseID + "/Tables/" + table0Name),
			},
			{
				Name:     new(table1Name),
				Type:     azto.Ptr(fablakehouse.TableTypeManaged),
				Format:   new("Delta"),
				Location: new("abfss://" + testhelp.RandomUUID() + "@onelake.dfs.fabric.microsoft.com/" + lakehouseID + "/Tables/" + table1Name),
			},
			{
				Name:     new(table2Name),
				Type:     azto.Ptr(fablakehouse.TableTypeManaged),
				Format:   new("Delta"),
				Location: new("abfss://" + testhelp.RandomUUID() + "@onelake.dfs.fabric.microsoft.com/" + lakehouseID + "/Tables/" + table2Name),
			},
		},
	}
}
