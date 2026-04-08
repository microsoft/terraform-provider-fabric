// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package operationsagent_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	faboperationsagent "github.com/microsoft/fabric-sdk-go/fabric/operationsagent"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

// var (
// 	KQLDatabaseID = testhelp.WellKnown()["KQLDatabase"].(map[string]any)["id"].(string)
// 	SQLDatabaseID = testhelp.WellKnown()["SQLDatabase"].(map[string]any)["id"].(string)
// )

var testHelperLocals = at.CompileLocalsConfig(map[string]any{
	"path": testhelp.GetFixturesDirPath("operations_agent"),
})

var testHelperDefinition = map[string]any{
	`"Configurations.json"`: map[string]any{
		"source": "${local.path}/Configurations.json.tmpl",
		// "tokens": map[string]any{
		// 	"DATASOURCE": map[string]any{"kqlDatabase": KQLDatabaseID},
		// },
	},
}

func TestUnit_OperationsAgentResource_Attributes(t *testing.T) {
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
						"display_name": testhelp.RandomName(),
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
						"workspace_id":    testhelp.RandomUUID(),
						"display_name":    testhelp.RandomName(),
						"unexpected_attr": "test",
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
						"display_name": testhelp.RandomName(),
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
						"workspace_id": testhelp.RandomUUID(),
					},
				)),
			ExpectError: regexp.MustCompile(`The argument "display_name" is required, but no definition was found.`),
		},
		// error - invalid format
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": testhelp.RandomUUID(),
						"display_name": testhelp.RandomName(),
						"format":       testhelp.RandomName(),
						"definition":   testHelperDefinition,
					},
				)),
			ExpectError: regexp.MustCompile(`Attribute format value must be one of: \["OperationsAgentV1"\]`),
		},
	}))
}

func TestUnit_OperationsAgentResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomOperationsAgentWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomOperationsAgentWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomOperationsAgentWithWorkspace(workspaceID))

	testCase := at.JoinConfigs(
		testHelperLocals,
		at.CompileConfig(
			testResourceItemHeader,
			map[string]any{
				"workspace_id": *entity.WorkspaceID,
				"display_name": *entity.DisplayName,
				"format":       string(faboperationsagent.DefinitionFormatOperationsAgentV1),
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

func TestUnit_OperationsAgentResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entityExist := fakes.NewRandomOperationsAgentWithWorkspace(workspaceID)
	entityBefore := fakes.NewRandomOperationsAgentWithWorkspace(workspaceID)
	entityAfter := fakes.NewRandomOperationsAgentWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomOperationsAgentWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(entityAfter)
	fakes.FakeServer.Upsert(fakes.NewRandomOperationsAgentWithWorkspace(workspaceID))

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
						"format":       string(faboperationsagent.DefinitionFormatOperationsAgentV1),
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
						"format":       string(faboperationsagent.DefinitionFormatOperationsAgentV1),
						"definition":   testHelperDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "folder_id", entityBefore.FolderID),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.state"),
			),
		},
		// Update and Read
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
						"folder_id":    *entityBefore.FolderID,
						"format":       string(faboperationsagent.DefinitionFormatOperationsAgentV1),
						"definition":   testHelperDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "folder_id", entityBefore.FolderID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.state"),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_OperationsAgentResource_CRUD(t *testing.T) {
	if testhelp.ShouldSkipTest(t) {
		t.Skip("No SPN support")
	}

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
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "folder_id", folderResourceFQN, "id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.state"),
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
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "folder_id", folderResourceFQN, "id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.state"),
			),
		},
	}))
}

func TestAcc_OperationsAgentDefinitionResource_CRUD(t *testing.T) {
	if testhelp.ShouldSkipTest(t) {
		t.Skip("No SPN support")
	}

	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entityCreateDisplayName := testhelp.RandomName()
	// entityUpdateDisplayName := testhelp.RandomName()
	// entityUpdateDescription := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityCreateDisplayName,
						"format":       string(faboperationsagent.DefinitionFormatOperationsAgentV1),
						"definition":   testHelperDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.state"),
			),
		},
		// Update and Read
		// {
		// 	ResourceName: testResourceItemFQN,
		// 	Config: at.JoinConfigs(
		// 		testHelperLocals,
		// 		at.CompileConfig(
		// 			testResourceItemHeader,
		// 			map[string]any{
		// 				"workspace_id": workspaceID,
		// 				"display_name": entityUpdateDisplayName,
		// 				"description":  entityUpdateDescription,
		// 				"format":       string(faboperationsagent.DefinitionFormatOperationsAgentV1),
		// 				"definition":   testHelperDefinition,
		// 			},
		// 		)),
		// 	Check: resource.ComposeAggregateTestCheckFunc(
		// 		resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
		// 		resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
		// 		resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
		// 		resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.state"),
		// 	),
		// },
	}))
}
