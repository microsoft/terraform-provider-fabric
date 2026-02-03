// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package eventhouse_test

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
	"path": testhelp.GetFixturesDirPath("eventhouse"),
})

var testHelperDefinition = map[string]any{
	`"EventhouseProperties.json"`: map[string]any{
		"source": "${local.path}/EventhouseProperties.json.tmpl",
	},
}

func TestUnit_EventhouseResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{},
				),
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - workspace_id - invalid UUID
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "invalid uuid",
						"display_name": "test",
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				)),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":    "00000000-0000-0000-0000-000000000000",
						"display_name":    "test",
						"unexpected_attr": "test",
						"format":          "Default",
						"definition":      testHelperDefinition,
					},
				)),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"display_name": "test",
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				)),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "00000000-0000-0000-0000-000000000000",
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				)),
			ExpectError: regexp.MustCompile(`The argument "display_name" is required, but no definition was found.`),
		},
		// error - no required attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"display_name": "test",
						"workspace_id": "00000000-0000-0000-0000-000000000000",
						"definition":   testHelperDefinition,
						"configuration": map[string]any{
							"minimum_consumption_units": "2.25",
						},
					},
				)),
			ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
		},
	}))
}

func TestUnit_EventhouseResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomEventhouseWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomEventhouseWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomEventhouseWithWorkspace(workspaceID))

	testCase := at.JoinConfigs(
		testHelperLocals,
		at.CompileConfig(
			testResourceItemHeader,
			map[string]any{
				"workspace_id": *entity.WorkspaceID,
				"display_name": *entity.DisplayName,
				"format":       "Default",
				"definition":   testHelperDefinition,
			},
		))

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

func TestUnit_EventhouseResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entityExist := fakes.NewRandomEventhouseWithWorkspace(workspaceID)
	entityBefore := fakes.NewRandomEventhouseWithWorkspace(workspaceID)
	entityAfter := fakes.NewRandomEventhouseWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomEventhouseWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(entityAfter)
	fakes.FakeServer.Upsert(fakes.NewRandomEventhouseWithWorkspace(workspaceID))

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - create - existing entity
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": *entityExist.WorkspaceID,
						"display_name": *entityExist.DisplayName,
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				)),
			ExpectError: regexp.MustCompile(common.ErrorCreateHeader),
		},
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": *entityBefore.WorkspaceID,
						"display_name": *entityBefore.DisplayName,
						"folder_id":    *entityBefore.FolderID,
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "folder_id", entityBefore.FolderID),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.query_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.ingestion_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_ids.0"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.minimum_consumption_units"),
			),
		},
		// Update, Move and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": *entityBefore.WorkspaceID,
						"display_name": *entityAfter.DisplayName,
						"description":  *entityAfter.Description,
						"folder_id":    *entityAfter.FolderID,
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "folder_id", entityAfter.FolderID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.query_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.ingestion_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_ids.0"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.minimum_consumption_units"),
			),
		},
		// Move and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": *entityBefore.WorkspaceID,
						"display_name": *entityAfter.DisplayName,
						"description":  *entityAfter.Description,
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "folder_id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.query_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.ingestion_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_ids.0"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.minimum_consumption_units"),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_EventhouseResource_CRUD(t *testing.T) {
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
				folderResourceHCL1,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityCreateDisplayName,
						"folder_id":    testhelp.RefByFQN(folderResourceFQN1, "id"),
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "folder_id", folderResourceFQN1, "id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.query_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.ingestion_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_ids.0"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.minimum_consumption_units"),
			),
		},
		// Update, Move and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				folderResourceHCL1,
				folderResourceHCL2,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityUpdateDisplayName,
						"description":  entityUpdateDescription,
						"folder_id":    testhelp.RefByFQN(folderResourceFQN2, "id"),
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "folder_id", folderResourceFQN2, "id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.query_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.ingestion_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_ids.0"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.minimum_consumption_units"),
			),
		},
		// Move and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				folderResourceHCL1,
				folderResourceHCL2,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityUpdateDisplayName,
						"description":  entityUpdateDescription,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "folder_id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.query_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.ingestion_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_ids.0"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.minimum_consumption_units"),
			),
		},
	},
	))
}

func TestAcc_EventhouseDefinitionResource_CRUD(t *testing.T) {
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
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.query_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.ingestion_service_uri"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "properties.database_ids.0"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.minimum_consumption_units"),
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
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.query_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.ingestion_service_uri"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "properties.database_ids.0"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.minimum_consumption_units"),
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
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.query_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.ingestion_service_uri"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "properties.database_ids.0"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.minimum_consumption_units"),
			),
		},
	},
	))
}

func TestAcc_EventhouseConfigurationResource_CRUD(t *testing.T) {
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
						"configuration": map[string]any{
							"minimum_consumption_units": "2.25",
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "folder_id", folderResourceFQN1, "id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.query_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.ingestion_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_ids.0"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.minimum_consumption_units"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.minimum_consumption_units", "2.25"),
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
						"configuration": map[string]any{
							"minimum_consumption_units": "2.25",
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "folder_id", folderResourceFQN2, "id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.query_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.ingestion_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_ids.0"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.minimum_consumption_units"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.minimum_consumption_units", "2.25"),
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
						"configuration": map[string]any{
							"minimum_consumption_units": "2.25",
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "folder_id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.query_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.ingestion_service_uri"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_ids.0"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.minimum_consumption_units"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "configuration.minimum_consumption_units", "2.25"),
			),
		},
	},
	))
}
