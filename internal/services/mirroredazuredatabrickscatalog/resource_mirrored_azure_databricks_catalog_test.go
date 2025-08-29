// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mirroredazuredatabrickscatalog_test

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
	"path": testhelp.GetFixturesDirPath("mirrored_azure_databricks_catalog"),
})

var testHelperDefinition = map[string]any{
	`"mirroringAzureDatabricksCatalog.json"`: map[string]any{
		"source": "${local.path}/mirroringAzureDatabricksCatalog.json.tmpl",
		"tokens": map[string]any{},
	},
}

func TestUnit_MirroredAzureDatabricksCatalogResource_Attributes(t *testing.T) {
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
							"catalog_name": testhelp.RandomName(),
						},
					},
				)),
			ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
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
						"configuration": map[string]any{
							"catalog_name":   testhelp.RandomName(),
							"mirroring_mode": "invalid mirroring mode",
						},
					},
				)),
			ExpectError: regexp.MustCompile(`Attribute configuration.mirroring_mode value must be one of`),
		},
	}))
}

func TestUnit_MirroredAzureDatabricksCatalogResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomMirroredAzureDatabricksCatalogWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomMirroredAzureDatabricksCatalogWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomMirroredAzureDatabricksCatalogWithWorkspace(workspaceID))

	testCase := at.JoinConfigs(
		testHelperLocals,
		at.CompileConfig(
			testResourceItemHeader,
			map[string]any{
				"workspace_id": *entity.WorkspaceID,
				"display_name": *entity.DisplayName,
			},
		))

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(common.ErrorImportIdentifierHeader),
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

func TestUnit_MirroredAzureDatabricksCatalogResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entityExist := fakes.NewRandomMirroredAzureDatabricksCatalogWithWorkspace(workspaceID)
	entityBefore := fakes.NewRandomMirroredAzureDatabricksCatalogWithWorkspace(workspaceID)
	entityAfter := fakes.NewRandomMirroredAzureDatabricksCatalogWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomMirroredAzureDatabricksCatalogWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(entityAfter)
	fakes.FakeServer.Upsert(fakes.NewRandomMirroredAzureDatabricksCatalogWithWorkspace(workspaceID))

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
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.auto_sync"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.catalog_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.databricks_workspace_connection_id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.mirror_status"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.mirroring_mode"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.onelake_tables_path"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.storage_connection_id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.sync_details.status"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.sync_details.last_sync_date_time"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.sql_endpoint_properties.connection_string"),
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
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.auto_sync"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.catalog_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.databricks_workspace_connection_id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.mirror_status"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.mirroring_mode"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.onelake_tables_path"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.storage_connection_id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.sync_details.status"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.sync_details.last_sync_date_time"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.sql_endpoint_properties.connection_string"),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_MirroredAzureDatabricksCatalogResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entityCreateDisplayName := testhelp.RandomName()
	entityUpdateDisplayName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"display_name": entityCreateDisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.auto_sync"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "properties.catalog_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.databricks_workspace_connection_id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.mirror_status"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.mirroring_mode"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.onelake_tables_path"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "properties.sync_details"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.sql_endpoint_properties.connection_string"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"display_name": entityUpdateDisplayName,
					"description":  entityUpdateDescription,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.auto_sync"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "properties.catalog_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.databricks_workspace_connection_id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.mirror_status"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.mirroring_mode"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.onelake_tables_path"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "properties.sync_details"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.sql_endpoint_properties.connection_string"),
			),
		},
	},
	))
}

// func TestAcc_MirroredAzureDatabricksCatalogDefinitionResource_CRUD(t *testing.T) {
// 	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
// 	workspaceID := workspace["id"].(string)

// 	entityCreateDisplayName := testhelp.RandomName()
// 	entityUpdateDisplayName := testhelp.RandomName()
// 	entityUpdateDescription := testhelp.RandomName()

// 	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
// 		// Create and Read
// 		{
// 			ResourceName: testResourceItemFQN,
// 			Config: at.JoinConfigs(
// 				testHelperLocals,
// 				at.CompileConfig(
// 					testResourceItemHeader,
// 					map[string]any{
// 						"workspace_id": workspaceID,
// 						"display_name": entityCreateDisplayName,
// 						"format":       "Default",
// 						"definition":   testHelperDefinition,
// 					},
// 				)),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.auto_sync"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.databricks_workspace_connection_id"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.mirror_status"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.mirroring_mode"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.onelake_tables_path"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.sql_endpoint_properties.connection_string"),
// 			),
// 		},
// 		// Update and Read
// 		{
// 			ResourceName: testResourceItemFQN,
// 			Config: at.JoinConfigs(
// 				testHelperLocals,
// 				at.CompileConfig(
// 					testResourceItemHeader,
// 					map[string]any{
// 						"workspace_id": workspaceID,
// 						"display_name": entityUpdateDisplayName,
// 						"description":  entityUpdateDescription,
// 						"format":       "Default",
// 						"definition":   testHelperDefinition,
// 					},
// 				)),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.auto_sync"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.databricks_workspace_connection_id"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.mirror_status"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.mirroring_mode"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.onelake_tables_path"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.sql_endpoint_properties.connection_string"),
// 			),
// 		},
// 	},
// 	))
// }

// func TestAcc_MirroredAzureDatabricksCatalogConfigurationResource_CRUD(t *testing.T) {
// 	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
// 	workspaceID := workspace["id"].(string)

// 	entityCreateDisplayName := testhelp.RandomName()
// 	entityUpdateDisplayName := testhelp.RandomName()
// 	entityUpdateDescription := testhelp.RandomName()

// 	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
// 		// Create and Read
// 		{
// 			ResourceName: testResourceItemFQN,
// 			Config: at.JoinConfigs(
// 				testHelperLocals,
// 				at.CompileConfig(
// 					testResourceItemHeader,
// 					map[string]any{
// 						"workspace_id": workspaceID,
// 						"display_name": entityCreateDisplayName,
// 						"configuration": map[string]any{
// 							"catalog_name":                       "",
// 							"databricks_workspace_connection_id": "00000000-0000-0000-0000-000000000000",
// 							"mirroring_mode":                     "Partial",
// 						},
// 					},
// 				)),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.auto_sync"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.databricks_workspace_connection_id"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.mirror_status"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.mirroring_mode"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.onelake_tables_path"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.sql_endpoint_properties.connection_string"),
// 			),
// 		},
// 		// Update and Read
// 		{
// 			ResourceName: testResourceItemFQN,
// 			Config: at.JoinConfigs(
// 				testHelperLocals,
// 				at.CompileConfig(
// 					testResourceItemHeader,
// 					map[string]any{
// 						"workspace_id": workspaceID,
// 						"display_name": entityUpdateDisplayName,
// 						"description":  entityUpdateDescription,
// 						"configuration": map[string]any{
// 							"catalog_name":                       "",
// 							"databricks_workspace_connection_id": "00000000-0000-0000-0000-000000000000",
// 							"mirroring_mode":                     "Partial",
// 						},
// 					},
// 				)),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.auto_sync"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.databricks_workspace_connection_id"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.mirror_status"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.mirroring_mode"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.onelake_tables_path"),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.sql_endpoint_properties.connection_string"),
// 			),
// 		},
// 	},
// 	))
// }
