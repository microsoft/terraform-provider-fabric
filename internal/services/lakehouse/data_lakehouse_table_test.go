// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testDataSourceLakehouseTable       = testhelp.DataSourceFQN("fabric", lakehouseTableTFName, "test")
	testDataSourceLakehouseTableHeader = at.DataSourceHeader(testhelp.TypeName("fabric", lakehouseTableTFName), "test")
)

func TestUnit_LakehouseTableDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	lakehouseID := testhelp.RandomUUID()
	lakehouseTables := NewRandomLakehouseTables(lakehouseID)
	fakes.FakeServer.ServerFactory.Lakehouse.TablesServer.NewListTablesPager = fakeLakehouseTablesFunc(lakehouseTables)

	entity := lakehouseTables.Data[1]

	resource.Test(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceLakehouseTableHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		// error - workspace_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceLakehouseTableHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
					"lakehouse_id": "invalid uuid",
					"name":         *entity.Name,
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceLakehouseTableHeader,
				map[string]any{
					"workspace_id":    workspaceID,
					"lakehouse_id":    lakehouseID,
					"name":            *entity.Name,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// read - not found
		{
			Config: at.CompileConfig(
				testDataSourceLakehouseTableHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"lakehouse_id": lakehouseID,
					"name":         testhelp.RandomName(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceLakehouseTableHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"lakehouse_id": lakehouseID,
					"name":         *entity.Name,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceLakehouseTable, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceLakehouseTable, "lakehouse_id", lakehouseID),
				resource.TestCheckResourceAttrPtr(testDataSourceLakehouseTable, "name", entity.Name),
				resource.TestCheckResourceAttrPtr(testDataSourceLakehouseTable, "location", entity.Location),
				resource.TestCheckResourceAttrPtr(testDataSourceLakehouseTable, "type", (*string)(entity.Type)),
				resource.TestCheckResourceAttrPtr(testDataSourceLakehouseTable, "format", entity.Format),
			),
		},
	}))
}

func TestAcc_LakehouseTableDataSource(t *testing.T) {
	workspaceID := *testhelp.WellKnown().Workspace.ID
	entity := testhelp.WellKnown().Lakehouse

	tableName := "publicholidays"

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read - not found
		{
			Config: at.CompileConfig(
				testDataSourceLakehouseTableHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"lakehouse_id": *entity.ID,
					"name":         testhelp.RandomName(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceLakehouseTableHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"lakehouse_id": *entity.ID,
					"name":         tableName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceLakehouseTable, "workspace_id", workspaceID),
				resource.TestCheckResourceAttrPtr(testDataSourceLakehouseTable, "lakehouse_id", entity.ID),
				resource.TestCheckResourceAttr(testDataSourceLakehouseTable, "name", tableName),
			),
		},
	},
	))
}
