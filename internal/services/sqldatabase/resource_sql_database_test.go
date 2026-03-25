// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sqldatabase_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	fabsqldatabase "github.com/microsoft/fabric-sdk-go/fabric/sqldatabase"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

var testHelperLocals = at.CompileLocalsConfig(map[string]any{ //nolint:gochecknoglobals
	"path": testhelp.GetFixturesDirPath("sql_database"),
})

func TestUnit_SQLDatabaseResource_Attributes(t *testing.T) {
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
		// error - invalid format value
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"format":       "invalid",
					"definition": map[string]any{
						`"test.dacpac"`: map[string]any{
							"source": "test.dacpac",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
		},
		// error - dacpac format with wrong definition path
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "00000000-0000-0000-0000-000000000000",
						"display_name": "test",
						"format":       "dacpac",
						"definition": map[string]any{
							`"definition.sqlproj"`: map[string]any{
								"source": "${local.path}/definition.sqlproj",
							},
						},
					},
				)),
			ExpectError: regexp.MustCompile(`Definition path must match`),
		},
		// error - sqlproj format with wrong definition path
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "00000000-0000-0000-0000-000000000000",
						"display_name": "test",
						"format":       "sqlproj",
						"definition": map[string]any{
							`"test.dacpac"`: map[string]any{
								"source": "${local.path}/definition.sqlproj",
							},
						},
					},
				)),
			ExpectError: regexp.MustCompile(`Definition path must match`),
		},
		// error - configuration conflicts with format
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "00000000-0000-0000-0000-000000000000",
						"display_name": "test",
						"format":       "sqlproj",
						"definition": map[string]any{
							`"definition.sqlproj"`: map[string]any{
								"source": "${local.path}/definition.sqlproj",
							},
						},
						"configuration": map[string]any{
							"creation_mode": "New",
						},
					},
				)),
			ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
		},
	}))
}

func TestUnit_SQLDatabaseResource_ConfigurationAttributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - configuration - empty
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
			ExpectError: regexp.MustCompile(`Inappropriate value for attribute "configuration"`),
		},
		// error - creation_mode - invalid value
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"creation_mode": "Invalid",
					},
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
		},
		// error - New mode with restore_point_in_time (should be null)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"creation_mode":         "New",
						"restore_point_in_time": "2026-03-22T00:00:00Z",
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute`),
		},
		// error - New mode with source_database_reference (should be null)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"creation_mode": "New",
						"source_database_reference": map[string]any{
							"reference_type": "ById",
							"item_id":        "00000000-0000-0000-0000-000000000000",
							"workspace_id":   "00000000-0000-0000-0000-000000000000",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute`),
		},
		// error - Restore mode without restore_point_in_time (required)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"creation_mode": "Restore",
						"source_database_reference": map[string]any{
							"reference_type": "ById",
							"item_id":        "00000000-0000-0000-0000-000000000000",
							"workspace_id":   "00000000-0000-0000-0000-000000000000",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute configuration.restore_point_in_time`),
		},
		// error - Restore mode without source_database_reference (required)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"creation_mode":         "Restore",
						"restore_point_in_time": "2026-03-22T00:00:00Z",
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute configuration.source_database_reference`),
		},
		// error - Restore mode with backup_retention_days (should be null)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"creation_mode":         "Restore",
						"restore_point_in_time": "2026-03-22T00:00:00Z",
						"backup_retention_days": 7,
						"source_database_reference": map[string]any{
							"reference_type": "ById",
							"item_id":        "00000000-0000-0000-0000-000000000000",
							"workspace_id":   "00000000-0000-0000-0000-000000000000",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute configuration.backup_retention_days`),
		},
		// error - Restore mode with collation (should be null)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"creation_mode":         "Restore",
						"restore_point_in_time": "2026-03-22T00:00:00Z",
						"collation":             "Latin1_General_100_BIN2_UTF8",
						"source_database_reference": map[string]any{
							"reference_type": "ById",
							"item_id":        "00000000-0000-0000-0000-000000000000",
							"workspace_id":   "00000000-0000-0000-0000-000000000000",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute configuration.collation`),
		},
		// error - ById reference_type without item_id (required)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"creation_mode":         "Restore",
						"restore_point_in_time": "2026-03-22T00:00:00Z",
						"source_database_reference": map[string]any{
							"reference_type": "ById",
							"workspace_id":   "00000000-0000-0000-0000-000000000000",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute configuration.source_database_reference.item_id`),
		},
		// error - ById reference_type without workspace_id (required)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"creation_mode":         "Restore",
						"restore_point_in_time": "2026-03-22T00:00:00Z",
						"source_database_reference": map[string]any{
							"reference_type": "ById",
							"item_id":        "00000000-0000-0000-0000-000000000000",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute configuration.source_database_reference.workspace_id`),
		},
		// error - ByVariable reference_type without variable_reference (required)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"creation_mode":         "Restore",
						"restore_point_in_time": "2026-03-22T00:00:00Z",
						"source_database_reference": map[string]any{
							"reference_type": "ByVariable",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute configuration.source_database_reference.variable_reference`),
		},
		// error - ByVariable reference_type with item_id (should be null)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"creation_mode":         "Restore",
						"restore_point_in_time": "2026-03-22T00:00:00Z",
						"source_database_reference": map[string]any{
							"reference_type":     "ByVariable",
							"variable_reference": "test_var",
							"item_id":            "00000000-0000-0000-0000-000000000000",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute configuration.source_database_reference.item_id`),
		},
		// error - ByVariable reference_type with workspace_id (should be null)
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"creation_mode":         "Restore",
						"restore_point_in_time": "2026-03-22T00:00:00Z",
						"source_database_reference": map[string]any{
							"reference_type":     "ByVariable",
							"variable_reference": "test_var",
							"workspace_id":       "00000000-0000-0000-0000-000000000000",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute configuration.source_database_reference.workspace_id`),
		},
		// error - invalid reference_type value
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"display_name": "test",
					"configuration": map[string]any{
						"creation_mode":         "Restore",
						"restore_point_in_time": "2026-03-22T00:00:00Z",
						"source_database_reference": map[string]any{
							"reference_type": "Invalid",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
		},
	}))
}

func TestUnit_SQLDatabaseResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomSQLDatabaseWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomSQLDatabaseWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomSQLDatabaseWithWorkspace(workspaceID))

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

func TestUnit_SQLDatabaseResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entityExist := fakes.NewRandomSQLDatabaseWithWorkspace(workspaceID)
	entityBefore := fakes.NewRandomSQLDatabaseWithWorkspace(workspaceID)
	entityAfter := fakes.NewRandomSQLDatabaseWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomSQLDatabaseWithWorkspace(workspaceID))
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(entityAfter)
	fakes.FakeServer.Upsert(fakes.NewRandomSQLDatabaseWithWorkspace(workspaceID))

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
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_string"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.server_fqdn"),
			),
		},
		// Update, Move and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": *entityBefore.WorkspaceID,
					"display_name": *entityBefore.DisplayName,
					"description":  *entityAfter.Description,
					"folder_id":    *entityAfter.FolderID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "folder_id", entityAfter.FolderID),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_string"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.server_fqdn"),
			),
		},
		// Move and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": *entityBefore.WorkspaceID,
					"display_name": *entityBefore.DisplayName,
					"description":  *entityAfter.Description,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "folder_id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_string"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.server_fqdn"),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestUnit_SQLDatabaseResource_CRUD_Configuration(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entityBefore := fakes.NewRandomSQLDatabaseWithWorkspace(workspaceID)
	sourceDB := fakes.NewRandomSQLDatabaseWithWorkspace(workspaceID)

	fakes.FakeServer.Upsert(sourceDB)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Create with New configuration and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": *entityBefore.WorkspaceID,
					"display_name": *entityBefore.DisplayName,
					"configuration": map[string]any{
						"creation_mode":         string(fabsqldatabase.CreationModeNew),
						"backup_retention_days": 7,
						"collation":             "SQL_Latin1_General_CP1_CI_AS",
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_string"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.server_fqdn"),
			),
		},
		// Create with Restore configuration and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": *entityBefore.WorkspaceID,
					"display_name": *entityBefore.DisplayName,
					"configuration": map[string]any{
						"creation_mode":         string(fabsqldatabase.CreationModeRestore),
						"restore_point_in_time": time.Now().Add(-24 * time.Hour).UTC().Format(time.RFC3339),
						"source_database_reference": map[string]any{
							"reference_type": "ById",
							"item_id":        *sourceDB.ID,
							"workspace_id":   *sourceDB.WorkspaceID,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_string"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.server_fqdn"),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_SQLDatabaseResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entityCreateDisplayName := testhelp.RandomName()
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
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_string"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.server_fqdn"),
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
						"display_name": entityCreateDisplayName,
						"description":  entityUpdateDescription,
						"folder_id":    testhelp.RefByFQN(folderResourceFQN2, "id"),
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttrPair(testResourceItemFQN, "folder_id", folderResourceFQN2, "id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_string"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.server_fqdn"),
			),
		},
		//	Move and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				folderResourceHCL2,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityCreateDisplayName,
						"description":  entityUpdateDescription,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "folder_id"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_string"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.server_fqdn"),
			),
		},
	},
	))
}

func TestAcc_SQLDatabaseDefinitionResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entityDisplayName := testhelp.RandomName()

	testHelperDefinition := map[string]any{
		`"definition.sqlproj"`: map[string]any{
			"source": "${local.path}/definition.sqlproj.tmpl",
		},
	}

	testhelperUpdateDefinition := map[string]any{
		`"definition.sqlproj"`: map[string]any{
			"source": "${local.path}/definition.sqlproj.tmpl",
		},
		`"dbo/Tables/TestTable.sql"`: map[string]any{
			"source": "${local.path}/dbo/Tables/TestTable.sql.tmpl",
		},
	}

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
						"display_name": entityDisplayName,
						"format":       "sqlproj",
						"definition":   testHelperDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_string"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.server_fqdn"),
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
						"workspace_id": workspaceID,
						"display_name": entityDisplayName,
						"format":       "sqlproj",
						"definition":   testhelperUpdateDefinition,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_string"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.server_fqdn"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "definition.dbo/Tables/TestTable.sql.source_content_sha256"),
			),
		},
	}))
}

func TestAcc_SQLDatabaseConfigurationResource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)
	entityCreateDisplayName := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create with New creation mode
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"display_name": entityCreateDisplayName,
					"configuration": map[string]any{
						"creation_mode":         string(fabsqldatabase.CreationModeNew),
						"backup_retention_days": 7,
						"collation":             "SQL_Latin1_General_CP1_CI_AS",
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.connection_string"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.database_name"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.server_fqdn"),
			),
		},
	}))
}
