// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testDataSourceLakehouseTables       = testhelp.DataSourceFQN("fabric", lakehouseTablesTFName, "test")
	testDataSourceLakehouseTablesHeader = at.DataSourceHeader(testhelp.TypeName("fabric", lakehouseTablesTFName), "test")
)

func TestUnit_LakehouseTablesDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	lakehouseID := testhelp.RandomUUID()
	lakehouseTables := NewRandomLakehouseTables(lakehouseID)
	fakes.FakeServer.ServerFactory.Lakehouse.TablesServer.NewListTablesPager = fakeLakehouseTablesFunc(lakehouseTables)

	tableName := lakehouseTables.Data[1].Name

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceLakehouseTablesHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		// error - workspace_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceLakehouseTablesHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
					"lakehouse_id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceLakehouseTablesHeader,
				map[string]any{
					"workspace_id":    workspaceID,
					"lakehouse_id":    lakehouseID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceLakehouseTablesHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"lakehouse_id": lakehouseID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceLakehouseTables, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceLakehouseTables, "lakehouse_id", lakehouseID),
				resource.TestCheckResourceAttrPtr(testDataSourceLakehouseTables, "values.1.name", tableName),
			),
		},
	}))
}

func TestAcc_LakehouseTablesDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entity := testhelp.WellKnown()["Lakehouse"].(map[string]any)
	entityID := entity["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceLakehouseTablesHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"lakehouse_id": entityID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceLakehouseTables, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceLakehouseTables, "lakehouse_id", entityID),
				resource.TestCheckResourceAttrSet(testDataSourceLakehouseTables, "values.0.name"),
			),
		},
	},
	))
}
