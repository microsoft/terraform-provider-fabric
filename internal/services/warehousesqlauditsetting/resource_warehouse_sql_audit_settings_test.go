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
					"warehouse_id":             testhelp.RandomUUID(),
					"state":                    string(fabwarehouse.AuditSettingsStateEnabled),
					"retention_days":           10,
					"audit_actions_and_groups": []string{"BATCH_COMPLETED_GROUP"},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - warehouse_id - invalid UUID
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":             testhelp.RandomUUID(),
					"warehouse_id":             "invalid uuid",
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
					"warehouse_id":             testhelp.RandomUUID(),
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
					"warehouse_id":             testhelp.RandomUUID(),
					"state":                    string(fabwarehouse.AuditSettingsStateEnabled),
					"retention_days":           10,
					"audit_actions_and_groups": []string{"BATCH_COMPLETED_GROUP"},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// error - missing warehouse_id
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
			ExpectError: regexp.MustCompile(`The argument "warehouse_id" is required, but no definition was found.`),
		},
		// error - invalid state value
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":             testhelp.RandomUUID(),
					"warehouse_id":             testhelp.RandomUUID(),
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
	warehouseID := testhelp.RandomUUID()
	workspaceID := testhelp.RandomUUID()

	entity := NewRandomSQLAuditSettings()
	fakeTestUpsertSQLAuditSettings(warehouseID, entity)

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
					"warehouse_id":             warehouseID,
					"state":                    string(*entity.State),
					"retention_days":           int(*entity.RetentionDays),
					"audit_actions_and_groups": entity.AuditActionsAndGroups,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "warehouse_id", warehouseID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "state", string(*entity.State)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "retention_days", strconv.Itoa(int(*entity.RetentionDays))),
				resource.TestCheckResourceAttr(testResourceItemFQN, "audit_actions_and_groups.#", "1"),
			),
		},
		// Update - omit state, expect default (Disabled)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":             workspaceID,
					"warehouse_id":             warehouseID,
					"retention_days":           10,
					"audit_actions_and_groups": []string{"BATCH_COMPLETED_GROUP"},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "state", string(fabwarehouse.AuditSettingsStateDisabled)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "retention_days", "10"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "audit_actions_and_groups.#", "1"),
			),
		},
		// Update - omit retention_days, expect default (0)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":             workspaceID,
					"warehouse_id":             warehouseID,
					"state":                    string(fabwarehouse.AuditSettingsStateEnabled),
					"audit_actions_and_groups": []string{"BATCH_COMPLETED_GROUP"},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "state", string(fabwarehouse.AuditSettingsStateEnabled)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "retention_days", "0"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "audit_actions_and_groups.#", "1"),
			),
		},
		// Update - omit audit_actions_and_groups, expect defaults (3 groups)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"warehouse_id":   warehouseID,
					"state":          string(fabwarehouse.AuditSettingsStateDisabled),
					"retention_days": 5,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "state", string(fabwarehouse.AuditSettingsStateDisabled)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "retention_days", "5"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "audit_actions_and_groups.#", "3"),
				resource.TestCheckTypeSetElemAttr(testResourceItemFQN, "audit_actions_and_groups.*", "SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP"),
				resource.TestCheckTypeSetElemAttr(testResourceItemFQN, "audit_actions_and_groups.*", "FAILED_DATABASE_AUTHENTICATION_GROUP"),
				resource.TestCheckTypeSetElemAttr(testResourceItemFQN, "audit_actions_and_groups.*", "BATCH_COMPLETED_GROUP"),
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
						"warehouse_id":             testhelp.RefByFQN(warehouseResourceFQN, "id"),
						"state":                    string(fabwarehouse.AuditSettingsStateDisabled),
						"retention_days":           10,
						"audit_actions_and_groups": []string{"BATCH_COMPLETED_GROUP"},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "warehouse_id", warehouseResourceFQN, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "state", string(fabwarehouse.AuditSettingsStateDisabled)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "retention_days", "10"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "audit_actions_and_groups.#", "1"),
				resource.TestCheckTypeSetElemAttr(testResourceItemFQN, "audit_actions_and_groups.*", "BATCH_COMPLETED_GROUP"),
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
						"warehouse_id": testhelp.RefByFQN(warehouseResourceFQN, "id"),
						"state":        string(fabwarehouse.AuditSettingsStateEnabled),
						"audit_actions_and_groups": []string{
							"SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP",
							"BATCH_COMPLETED_GROUP",
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "state", string(fabwarehouse.AuditSettingsStateEnabled)),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "warehouse_id", warehouseResourceFQN, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "retention_days", "0"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "audit_actions_and_groups.#", "2"),
				resource.TestCheckTypeSetElemAttr(testResourceItemFQN, "audit_actions_and_groups.*", "SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP"),
				resource.TestCheckTypeSetElemAttr(testResourceItemFQN, "audit_actions_and_groups.*", "BATCH_COMPLETED_GROUP"),
			),
		},
		// Update - reset all to defaults
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				warehouseResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"warehouse_id": testhelp.RefByFQN(warehouseResourceFQN, "id"),
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "state", string(fabwarehouse.AuditSettingsStateDisabled)),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "warehouse_id", warehouseResourceFQN, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "retention_days", "0"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "audit_actions_and_groups.#", "3"),
				resource.TestCheckTypeSetElemAttr(testResourceItemFQN, "audit_actions_and_groups.*", "SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP"),
				resource.TestCheckTypeSetElemAttr(testResourceItemFQN, "audit_actions_and_groups.*", "FAILED_DATABASE_AUTHENTICATION_GROUP"),
				resource.TestCheckTypeSetElemAttr(testResourceItemFQN, "audit_actions_and_groups.*", "BATCH_COMPLETED_GROUP"),
			),
		},
	}))
}
