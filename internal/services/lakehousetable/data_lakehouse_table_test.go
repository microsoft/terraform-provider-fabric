// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package lakehousetable_test

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

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

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
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found`),
		},
		// error - workspace_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
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
				testDataSourceItemHeader,
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
				testDataSourceItemHeader,
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
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"lakehouse_id": lakehouseID,
					"name":         *entity.Name,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "lakehouse_id", lakehouseID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "name", entity.Name),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "location", entity.Location),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "type", (*string)(entity.Type)),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "format", entity.Format),
			),
		},
	}))
}

func TestAcc_LakehouseTableDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entity := testhelp.WellKnown()["Lakehouse"].(map[string]any)
	entityID := entity["id"].(string)
	entityTableName := entity["tableName"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"lakehouse_id": entityID,
					"name":         testhelp.RandomName(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"lakehouse_id": entityID,
					"name":         entityTableName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "lakehouse_id", entityID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "name", entityTableName),
			),
		},
	},
	))
}
