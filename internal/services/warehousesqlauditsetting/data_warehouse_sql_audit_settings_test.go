// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package warehousesqlauditsetting_test

import (
	"regexp"
	"strconv"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_WarehouseSQLAuditSettingsDataSource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testDataSourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			ResourceName: testDataSourceItemFQN,
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - workspace_id - invalid UUID
		{
			ResourceName: testDataSourceItemFQN,
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
					"item_id":      testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - item_id - invalid UUID
		{
			ResourceName: testDataSourceItemFQN,
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"item_id":      "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_WarehouseSQLAuditSettingsDataSource_CRUD(t *testing.T) {
	itemID := testhelp.RandomUUID()
	workspaceID := testhelp.RandomUUID()
	entity := NewRandomSQLAuditSettings()
	fakeTestUpsertSQLAuditSettings(itemID, entity)

	fakes.FakeServer.ServerFactory.Warehouse.SQLAuditSettingsServer.GetSQLAuditSettings = fakeGetSQLAuditSettingsFunc()

	resource.Test(t, testhelp.NewTestUnitCase(t, &testDataSourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Read
		{
			ResourceName: testDataSourceItemFQN,
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "item_id", itemID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "state", string(*entity.State)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "retention_days", strconv.Itoa(int(*entity.RetentionDays))),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "audit_actions_and_groups.#", "1"),
			),
		},
	}))
}

func TestAcc_WarehouseSQLAuditSettingsDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entity := testhelp.WellKnown()["Warehouse"].(map[string]any)
	entityID := entity["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      entityID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "item_id", entityID),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "state"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "retention_days"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "audit_actions_and_groups.#"),
			),
		},
	},
	))
}
