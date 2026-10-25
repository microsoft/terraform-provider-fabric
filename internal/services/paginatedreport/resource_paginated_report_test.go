// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package paginatedreport_test

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
	"path": testhelp.GetFixturesDirPath("paginated_report"),
})

func testDefinition(displayName string) map[string]any {
	return map[string]any{
		fmt.Sprintf("%q", displayName+".rdl"): map[string]any{
			"source": "${local.path}/report.rdl",
		},
	}
}

func TestUnit_PaginatedReportResource_Attributes(t *testing.T) {
	validDefinition := testDefinition("test")

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config:       at.CompileConfig(testResourceItemHeader, map[string]any{}),
			ExpectError:  regexp.MustCompile(`Missing required argument`),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(testHelperLocals, at.CompileConfig(testResourceItemHeader, map[string]any{
				"workspace_id": "invalid uuid",
				"display_name": "test",
				"format":       "Default",
				"definition":   validDefinition,
			})),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(testHelperLocals, at.CompileConfig(testResourceItemHeader, map[string]any{
				"workspace_id":    "00000000-0000-0000-0000-000000000000",
				"display_name":    "test",
				"unexpected_attr": "test",
				"format":          "Default",
				"definition":      validDefinition,
			})),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(testHelperLocals, at.CompileConfig(testResourceItemHeader, map[string]any{
				"display_name": "test",
				"format":       "Default",
				"definition":   validDefinition,
			})),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(testHelperLocals, at.CompileConfig(testResourceItemHeader, map[string]any{
				"workspace_id": "00000000-0000-0000-0000-000000000000",
				"format":       "Default",
				"definition":   validDefinition,
			})),
			ExpectError: regexp.MustCompile(`The argument "display_name" is required, but no definition was found.`),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(testResourceItemHeader, map[string]any{
				"workspace_id": "00000000-0000-0000-0000-000000000000",
				"display_name": "test",
			}),
			ExpectError: regexp.MustCompile(`The argument "definition" is required, but no definition was found.`),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(testHelperLocals, at.CompileConfig(testResourceItemHeader, map[string]any{
				"workspace_id": "00000000-0000-0000-0000-000000000000",
				"display_name": "test",
				"format":       "Default",
				"definition": map[string]any{
					`"test.rdl"`:  map[string]any{"source": "${local.path}/report.rdl"},
					`"other.rdl"`: map[string]any{"source": "${local.path}/report.rdl"},
				},
			})),
			ExpectError: regexp.MustCompile(`map must contain at least 1 elements and at most 1\s+elements`),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(testHelperLocals, at.CompileConfig(testResourceItemHeader, map[string]any{
				"workspace_id": "00000000-0000-0000-0000-000000000000",
				"display_name": "test",
				"format":       "Default",
				"definition": map[string]any{
					`"test.txt"`: map[string]any{"source": "${local.path}/report.rdl"},
				},
			})),
			ExpectError: regexp.MustCompile(`Definition path must match`),
		},
	}))
}

func TestUnit_PaginatedReportResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID)
	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID))

	testCase := at.JoinConfigs(testHelperLocals, at.CompileConfig(testResourceItemHeader, map[string]any{
		"workspace_id": *entity.WorkspaceID,
		"display_name": *entity.DisplayName,
		"format":       "Default",
		"definition":   testDefinition(*entity.DisplayName),
	}))

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
		{
			ResourceName:       testResourceItemFQN,
			Config:             testCase,
			ImportStateId:      fmt.Sprintf("%s/%s", *entity.WorkspaceID, *entity.ID),
			ImportState:        true,
			ImportStatePersist: true,
			ImportStateCheck: func(states []*terraform.InstanceState) error {
				if len(states) != 1 || states[0].ID != *entity.ID {
					return errors.New(testResourceItemFQN + ": unexpected imported state")
				}

				return nil
			},
		},
	}))
}

func TestUnit_PaginatedReportResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entityExist := fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID)
	entityBefore := fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID)
	entityAfter := fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID)
	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID))
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(entityAfter)
	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID))

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(testHelperLocals, at.CompileConfig(testResourceItemHeader, map[string]any{
				"workspace_id": *entityExist.WorkspaceID,
				"display_name": *entityExist.DisplayName,
				"format":       "Default",
				"definition":   testDefinition(*entityExist.DisplayName),
			})),
			ExpectError: regexp.MustCompile(common.ErrorCreateHeader),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(testHelperLocals, at.CompileConfig(testResourceItemHeader, map[string]any{
				"workspace_id": *entityBefore.WorkspaceID,
				"display_name": *entityBefore.DisplayName,
				"folder_id":    *entityBefore.FolderID,
				"format":       "Default",
				"definition":   testDefinition(*entityBefore.DisplayName),
			})),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "folder_id", entityBefore.FolderID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "format", "Default"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
			),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(testHelperLocals, at.CompileConfig(testResourceItemHeader, map[string]any{
				"workspace_id": *entityBefore.WorkspaceID,
				"display_name": *entityAfter.DisplayName,
				"description":  *entityAfter.Description,
				"folder_id":    *entityBefore.FolderID,
				"format":       "Default",
				"definition":   testDefinition(*entityAfter.DisplayName),
			})),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "folder_id", entityBefore.FolderID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
			),
		},
	}))
}

func TestAcc_PaginatedReportResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)
	folder1 := testhelp.WellKnown()["FolderRS1"].(map[string]any)
	folder1ID := folder1["id"].(string)
	folder2 := testhelp.WellKnown()["FolderRS2"].(map[string]any)
	folder2ID := folder2["id"].(string)
	createName := testhelp.RandomName()
	updateName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(testResourceItemHeader, map[string]any{
					"workspace_id": workspaceID,
					"display_name": createName,
					"folder_id":    folder1ID,
					"format":       "Default",
					"definition":   testDefinition(createName),
				}),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", createName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttr(testResourceItemFQN, "folder_id", folder1ID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "format", "Default"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
			),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(testResourceItemHeader, map[string]any{
					"workspace_id": workspaceID,
					"display_name": updateName,
					"description":  entityUpdateDescription,
					"folder_id":    folder2ID,
					"format":       "Default",
					"definition":   testDefinition(updateName),
				}),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", updateName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttr(testResourceItemFQN, "folder_id", folder2ID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
			),
		},
	}))
}
