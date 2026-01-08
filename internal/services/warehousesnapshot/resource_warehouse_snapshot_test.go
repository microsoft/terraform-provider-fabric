// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehousesnapshot_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_WarehouseSnapshotResource_Attributes(t *testing.T) {
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
					"workspace_id": "invalid uuid",
					"display_name": "test",
					"configuration": map[string]any{
						"parent_warehouse_id": testhelp.RandomUUID(),
					},
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
					"workspace_id":    "00000000-0000-0000-0000-000000000000",
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": "test",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "display_name" is required, but no definition was found.`),
		},
		// error - no required attributes (configuration)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "configuration" is required, but no definition was found.`),
		},
		// error - no required configuration attributes (parent_warehouse_id)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":  "00000000-0000-0000-0000-000000000000",
					"display_name":  "test",
					"configuration": map[string]any{},
				},
			),
			ExpectError: regexp.MustCompile(`Inappropriate value for attribute "configuration".`),
		},
	}))
}

func TestUnit_WarehouseSnapshotResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomWarehouseSnapshotWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomWarehouseSnapshotWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomWarehouseSnapshotWithWorkspace(workspaceID))

	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id": *entity.WorkspaceID,
			"display_name": *entity.DisplayName,
			"configuration": map[string]any{
				"parent_warehouse_id": *entity.Properties.ParentWarehouseID,
			},
		},
	)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(fmt.Sprintf(common.ErrorImportIdentifierDetails, "WorkspaceID/WarehouseSnapshotID")),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "test/id",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s", "test", *entity.ID),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s", *entity.WorkspaceID, "test"),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// Import state testing
		{
			ResourceName:       testResourceItemFQN,
			Config:             testCase,
			ImportStateId:      fmt.Sprintf("%s/%s", *entity.WorkspaceID, *entity.ID),
			ImportState:        true,
			ImportStatePersist: true,
			ImportStateCheck: func(is []*terraform.InstanceState) error {
				if len(is) != 1 {
					return errors.New("expected one instance state")
				}

				if is[0].ID != *entity.ID {
					return errors.New(testResourceItemFQN + ": unexpected ID")
				}

				return nil
			},
		},
	}))
}

func TestUnit_WarehouseSnapshotResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entityExist := fakes.NewRandomWarehouseSnapshotWithWorkspace(workspaceID)
	entityBefore := fakes.NewRandomWarehouseSnapshotWithWorkspace(workspaceID)
	entityAfter := fakes.NewRandomWarehouseSnapshotWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomWarehouseSnapshotWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(entityAfter)
	fakes.FakeServer.Upsert(fakes.NewRandomWarehouseSnapshotWithWorkspace(workspaceID))

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - create - existing entity
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": *entityExist.WorkspaceID,
					"display_name": *entityExist.DisplayName,
					"configuration": map[string]any{
						"parent_warehouse_id": *entityExist.Properties.ParentWarehouseID,
						"snapshot_date_time":  entityExist.Properties.SnapshotDateTime.Format(time.RFC3339),
					},
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorCreateHeader),
		},
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": *entityBefore.WorkspaceID,
					"display_name": *entityBefore.DisplayName,
					"folder_id":    *entityBefore.FolderID,
					"configuration": map[string]any{
						"parent_warehouse_id": *entityBefore.Properties.ParentWarehouseID,
						"snapshot_date_time":  entityBefore.Properties.SnapshotDateTime.Format(time.RFC3339),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "folder_id", entityBefore.FolderID),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_string"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.snapshot_date_time"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.parent_warehouse_id"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": *entityBefore.WorkspaceID,
					"display_name": *entityAfter.DisplayName,
					"folder_id":    *entityBefore.FolderID,
					"description":  *entityAfter.Description,
					"configuration": map[string]any{
						"parent_warehouse_id": *entityBefore.Properties.ParentWarehouseID,
						"snapshot_date_time":  entityBefore.Properties.SnapshotDateTime.Format(time.RFC3339),
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "folder_id", entityBefore.FolderID),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_string"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.snapshot_date_time"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.parent_warehouse_id"),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_WarehouseSnapshotResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entityCreateDisplayName := testhelp.RandomName()
	entityUpdateDisplayName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()

	folderResourceHCL, folderResourceFQN := testhelp.FolderResource(t, workspaceID)
	warehouseResourceHCL := at.CompileConfig(
		at.ResourceHeader(testhelp.TypeName("fabric", "warehouse"), "test"),
		map[string]any{
			"display_name": testhelp.RandomName(),
			"workspace_id": workspaceID,
			"folder_id":    testhelp.RefByFQN(folderResourceFQN, "id"),
		},
	)

	warehouseResourceFQN := testhelp.ResourceFQN("fabric", "warehouse", "test")

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				warehouseResourceHCL,
				folderResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityCreateDisplayName,
						"folder_id":    testhelp.RefByFQN(folderResourceFQN, "id"),
						"configuration": map[string]any{
							"parent_warehouse_id": testhelp.RefByFQN(warehouseResourceFQN, "id"),
							"snapshot_date_time":  testhelp.RefByFQN(warehouseResourceFQN, "properties.last_updated_time"),
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "folder_id", folderResourceFQN, "id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_string"),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "properties.snapshot_date_time", warehouseResourceFQN, "properties.last_updated_time"),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "properties.parent_warehouse_id", warehouseResourceFQN, "id"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				warehouseResourceHCL,
				folderResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityUpdateDisplayName,
						"description":  entityUpdateDescription,
						"folder_id":    testhelp.RefByFQN(folderResourceFQN, "id"),
						"configuration": map[string]any{
							"parent_warehouse_id": testhelp.RefByFQN(warehouseResourceFQN, "id"),
							"snapshot_date_time":  testhelp.RefByFQN(warehouseResourceFQN, "properties.last_updated_time"),
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "folder_id", folderResourceFQN, "id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_string"),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "properties.snapshot_date_time", warehouseResourceFQN, "properties.last_updated_time"),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "properties.parent_warehouse_id", warehouseResourceFQN, "id"),
			),
		},
	},
	))
}
