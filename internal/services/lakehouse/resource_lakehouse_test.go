// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_LakehouseResource_Attributes(t *testing.T) {
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
					"workspace_id":  "00000000-0000-0000-0000-000000000000",
					"display_name":  "test",
					"configuration": map[string]any{},
				},
			),
			ExpectError: regexp.MustCompile(`Inappropriate value for attribute "configuration".`),
		},
	}))
}

func TestUnit_LakehouseResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomLakehouseWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomLakehouseWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomLakehouseWithWorkspace(workspaceID))

	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id": *entity.WorkspaceID,
			"display_name": *entity.DisplayName,
		},
	)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(fmt.Sprintf(common.ErrorImportIdentifierDetails, "WorkspaceID/LakehouseID")),
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

func TestUnit_LakehouseResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entityExist := fakes.NewRandomLakehouseWithWorkspace(workspaceID)
	entityBefore := fakes.NewRandomLakehouseWithWorkspace(workspaceID)
	entityAfter := fakes.NewRandomLakehouseWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomLakehouseWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(entityAfter)
	fakes.FakeServer.Upsert(fakes.NewRandomLakehouseWithWorkspace(workspaceID))

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - create - existing entity
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": *entityExist.WorkspaceID,
					"display_name": *entityExist.DisplayName,
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
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "folder_id", entityBefore.FolderID),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.onelake_files_path"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.onelake_tables_path"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.default_schema"),
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
					"description":  *entityAfter.Description,
					"folder_id":    *entityBefore.FolderID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "folder_id", entityBefore.FolderID),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.onelake_files_path"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.onelake_tables_path"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.default_schema"),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_LakehouseResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entityCreateDisplayName := testhelp.RandomName()
	entityUpdateDisplayName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()

	folderResourceHCL, folderResourceFQN := testhelp.FolderResource(t, workspaceID)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				folderResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityCreateDisplayName,
						"folder_id":    testhelp.RefByFQN(folderResourceFQN, "id"),
					},
				)),

			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "folder_id"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "configuration"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "properties.default_schema"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				folderResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityUpdateDisplayName,
						"description":  entityUpdateDescription,
						"folder_id":    testhelp.RefByFQN(folderResourceFQN, "id"),
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "folder_id"),

				resource.TestCheckNoResourceAttr(testResourceItemFQN, "configuration"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "properties.default_schema"),
			),
		},
		// Recreate resource and read - goal is to check handling for "ItemDisplayNameNotAvailableYet" error
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				folderResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityUpdateDisplayName,
						"configuration": map[string]any{
							"enable_schemas": true,
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "folder_id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.enable_schemas", "true"),
			),
		},
	},
	))
}

func TestAcc_LakehouseResource_CRUD_Configuration(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entityCreateDisplayName1 := testhelp.RandomName()
	entityUpdateDisplayName1 := testhelp.RandomName()
	entityUpdateDescription1 := testhelp.RandomName()

	entityCreateDisplayName2 := testhelp.RandomName()
	entityUpdateDisplayName2 := testhelp.RandomName()
	entityUpdateDescription2 := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read (configuration)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"display_name": entityCreateDisplayName1,
					"configuration": map[string]any{
						"enable_schemas": true,
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName1),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.enable_schemas", "true"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "properties.default_schema", "dbo"),
			),
		},
		// Update and Read (configuration)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"display_name": entityUpdateDisplayName1,
					"description":  entityUpdateDescription1,
					"configuration": map[string]any{
						"enable_schemas": true,
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName1),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription1),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.enable_schemas", "true"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "properties.default_schema", "dbo"),
			),
		},
		// Create and Read (configuration)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"display_name": entityCreateDisplayName2,
					"configuration": map[string]any{
						"enable_schemas": false,
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName2),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.enable_schemas", "false"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "properties.default_schema"),
			),
		},
		// Update and Read (configuration)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"display_name": entityUpdateDisplayName2,
					"description":  entityUpdateDescription2,
					"configuration": map[string]any{
						"enable_schemas": false,
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName2),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription2),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.enable_schemas", "false"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "properties.default_schema"),
			),
		},
	},
	))
}
