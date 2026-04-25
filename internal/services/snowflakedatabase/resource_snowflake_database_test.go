// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package snowflakedatabase_test

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

var testHelperLocals = at.CompileLocalsConfig(map[string]any{
	"path": testhelp.GetFixturesDirPath("snowflake_database"),
})

var testHelperDefinition = map[string]any{
	`"SnowflakeDatabaseProperties.json"`: map[string]any{
		"source": "${local.path}/SnowflakeDatabaseProperties.json.tmpl",
	},
}

func TestUnit_SnowflakeDatabaseResource_Attributes(t *testing.T) {
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
					"display_name":    "test",
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
		// error - missing required configuration attribute
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

func TestUnit_SnowflakeDatabaseResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID))

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
			ExpectError:   regexp.MustCompile(fmt.Sprintf(common.ErrorImportIdentifierDetails, fmt.Sprintf("WorkspaceID/%sID", string(fabricItemType)))),
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

func TestUnit_SnowflakeDatabaseResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entityExist := fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID)
	entityBefore := fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID)
	entityAfter := fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(entityAfter)
	fakes.FakeServer.Upsert(fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID))

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
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.snowflake_database_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.sql_endpoint_properties.id"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": *entityBefore.WorkspaceID,
					"display_name": *entityBefore.DisplayName,
					"description":  *entityBefore.Description,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityBefore.Description),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.snowflake_database_name"),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestUnit_SnowflakeDatabaseResource_CRUD_Configuration(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entityCreate := fakes.NewRandomSnowflakeDatabaseWithWorkspace(workspaceID)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Create and Read with configuration
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": *entityCreate.WorkspaceID,
					"display_name": *entityCreate.DisplayName,
					"configuration": map[string]any{
						"connection_id":           testhelp.RandomUUID(),
						"snowflake_database_name": "ExampleDatabase",
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityCreate.DisplayName),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "configuration.connection_id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.snowflake_database_name", "ExampleDatabase"),
			),
		},
	}))
}

func TestAcc_SnowflakeDatabaseResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entityCreateDisplayName := testhelp.RandomName()
	entityUpdateDisplayName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()

	folderResourceHCL1, folderResourceFQN1 := testhelp.FolderResource(t, workspaceID)
	folderResourceHCL2, folderResourceFQN2 := testhelp.FolderResource(t, workspaceID)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				folderResourceHCL1,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityCreateDisplayName,
						"folder_id":    testhelp.RefByFQN(folderResourceFQN1, "id"),
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "folder_id", folderResourceFQN1, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
			),
		},
		// Update, Move and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				folderResourceHCL1,
				folderResourceHCL2,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityUpdateDisplayName,
						"description":  entityUpdateDescription,
						"folder_id":    testhelp.RefByFQN(folderResourceFQN2, "id"),
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "folder_id", folderResourceFQN2, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
			),
		},
		// Move and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				folderResourceHCL2,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityUpdateDisplayName,
						"description":  entityUpdateDescription,
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "folder_id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
			),
		},
	},
	))
}
