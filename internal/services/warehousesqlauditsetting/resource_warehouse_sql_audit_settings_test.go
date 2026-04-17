// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package warehousesqlauditsetting_test

import (
	"regexp"
	"strconv"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabwarehouse "github.com/microsoft/fabric-sdk-go/fabric/warehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_WarehouseSQLAuditSettingsResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - workspace_id - invalid UUID
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":             "invalid uuid",
					"item_id":                  testhelp.RandomUUID(),
					"state":                    string(fabwarehouse.AuditSettingsStateEnabled),
					"retention_days":           10,
					"audit_actions_and_groups": []string{"BATCH_COMPLETED_GROUP"},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - item_id - invalid UUID
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":             testhelp.RandomUUID(),
					"item_id":                  "invalid uuid",
					"state":                    string(fabwarehouse.AuditSettingsStateEnabled),
					"retention_days":           10,
					"audit_actions_and_groups": []string{"BATCH_COMPLETED_GROUP"},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":             testhelp.RandomUUID(),
					"item_id":                  testhelp.RandomUUID(),
					"state":                    string(fabwarehouse.AuditSettingsStateEnabled),
					"retention_days":           10,
					"audit_actions_and_groups": []string{"BATCH_COMPLETED_GROUP"},
					"unexpected_attr":          "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - missing workspace_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"item_id":                  testhelp.RandomUUID(),
					"state":                    string(fabwarehouse.AuditSettingsStateEnabled),
					"retention_days":           10,
					"audit_actions_and_groups": []string{"BATCH_COMPLETED_GROUP"},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// error - missing item_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":             testhelp.RandomUUID(),
					"state":                    string(fabwarehouse.AuditSettingsStateEnabled),
					"retention_days":           10,
					"audit_actions_and_groups": []string{"BATCH_COMPLETED_GROUP"},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "item_id" is required, but no definition was found.`),
		},
		// error - invalid state value
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":             testhelp.RandomUUID(),
					"item_id":                  testhelp.RandomUUID(),
					"state":                    "InvalidState",
					"retention_days":           10,
					"audit_actions_and_groups": []string{"BATCH_COMPLETED_GROUP"},
				},
			),
			ExpectError: regexp.MustCompile(`Attribute state value must be one of`),
		},
	}))
}

func TestUnit_WarehouseSQLAuditSettingsResource_CRUD(t *testing.T) {
	itemID := testhelp.RandomUUID()
	workspaceID := testhelp.RandomUUID()

	entity := NewRandomSQLAuditSettings()
	fakeTestUpsertSQLAuditSettings(itemID, entity)

	fakes.FakeServer.ServerFactory.Warehouse.SQLAuditSettingsServer.GetSQLAuditSettings = fakeGetSQLAuditSettingsFunc()
	fakes.FakeServer.ServerFactory.Warehouse.SQLAuditSettingsServer.UpdateSQLAuditSettings = fakeUpdateSQLAuditSettingsFunc()
	fakes.FakeServer.ServerFactory.Warehouse.SQLAuditSettingsServer.SetAuditActionsAndGroups = fakeSetAuditActionsAndGroupsFunc()

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":             workspaceID,
					"item_id":                  itemID,
					"state":                    string(*entity.State),
					"retention_days":           int(*entity.RetentionDays),
					"audit_actions_and_groups": entity.AuditActionsAndGroups,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "item_id", itemID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "state", string(*entity.State)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "retention_days", strconv.Itoa(int(*entity.RetentionDays))),
				resource.TestCheckResourceAttr(testResourceItemFQN, "audit_actions_and_groups.#", "1"),
			),
		},
		// Update
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":             workspaceID,
					"item_id":                  itemID,
					"state":                    string(fabwarehouse.AuditSettingsStateDisabled),
					"retention_days":           0,
					"audit_actions_and_groups": []string{testhelp.RandomName(), testhelp.RandomName()},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "state", string(fabwarehouse.AuditSettingsStateDisabled)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "retention_days", "0"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "audit_actions_and_groups.#", "2"),
			),
		},
	}))
}

func TestAcc_WarehouseSQLAuditSettingsResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	warehouseResourceHCL, warehouseResourceFQN := warehouseResource(t, workspaceID)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				warehouseResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":             workspaceID,
						"item_id":                  testhelp.RefByFQN(warehouseResourceFQN, "id"),
						"state":                    string(fabwarehouse.AuditSettingsStateEnabled),
						"retention_days":           10,
						"audit_actions_and_groups": []string{"BATCH_COMPLETED_GROUP"},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "item_id", warehouseResourceFQN, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "state", string(fabwarehouse.AuditSettingsStateEnabled)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "retention_days", "10"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "audit_actions_and_groups.#", "1"),
			),
		},
		// Update
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				warehouseResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"item_id":      testhelp.RefByFQN(warehouseResourceFQN, "id"),
						"state":        string(fabwarehouse.AuditSettingsStateDisabled),
						"audit_actions_and_groups": []string{
							"SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP",
							"FAILED_DATABASE_AUTHENTICATION_GROUP",
							"BATCH_COMPLETED_GROUP",
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "state", string(fabwarehouse.AuditSettingsStateDisabled)),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "item_id", warehouseResourceFQN, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "retention_days", "10"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "audit_actions_and_groups.#", "3"),
			),
		},
	}))
}
