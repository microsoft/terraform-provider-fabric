// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package orgappaudience_test

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

var testHelperLocals = at.CompileLocalsConfig(map[string]any{
	"path": testhelp.GetFixturesDirPath("org_app_audience"),
})

var testHelperDefinition = map[string]any{
	`"definition.json"`: map[string]any{
		"source": "${local.path}/definition.json",
	},
}

func testHelperDefinitionTemplate(parentAppID string, hasAccessToHiddenContent bool, tabOrder int) map[string]any {
	return map[string]any{
		`"definition.json"`: map[string]any{
			"source": "${local.path}/definition.json.tmpl",
			"tokens": map[string]any{
				"ParentAppId":              parentAppID,
				"HasAccessToHiddenContent": strconv.FormatBool(hasAccessToHiddenContent),
				"TabOrder":                 strconv.Itoa(tabOrder),
			},
		},
	}
}

func testHelperParentOrgAppResource(workspaceID string) (resourceHCL, resourceFQN string) {
	resourceHCL = at.CompileConfig(
		at.ResourceHeader(testhelp.TypeName(common.ProviderTypeName, "org_app"), "test"),
		map[string]any{
			"workspace_id": workspaceID,
			"display_name": testhelp.RandomName(),
		},
	)
	resourceFQN = testhelp.ResourceFQN(common.ProviderTypeName, "org_app", "test")

	return resourceHCL, resourceFQN
}

func TestUnit_OrgAppAudienceResource_Attributes(t *testing.T) {
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
				),
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - folder_id - invalid UUID
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "00000000-0000-0000-0000-000000000000",
						"display_name": "test",
						"folder_id":    "invalid uuid",
						"format":       "Default",
						"definition":   testHelperDefinition,
					},
				),
			),
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
				),
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attributes - workspace_id
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
				),
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// error - no required attributes - display_name
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
				),
			),
			ExpectError: regexp.MustCompile(`The argument "display_name" is required, but no definition was found.`),
		},
		// error - no required attributes - definition
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "00000000-0000-0000-0000-000000000000",
						"display_name": "test",
						"format":       "Default",
					},
				),
			),
			ExpectError: regexp.MustCompile(`The argument "definition" is required, but no definition was found.`),
		},
		// error - no required attributes - format
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "00000000-0000-0000-0000-000000000000",
						"display_name": "test",
						"definition":   testHelperDefinition,
					},
				),
			),
			ExpectError: regexp.MustCompile(`The argument "format" is required, but no definition was found.`),
		},
		// error - format - invalid value
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "00000000-0000-0000-0000-000000000000",
						"display_name": "test",
						"format":       "Invalid",
						"definition":   testHelperDefinition,
					},
				),
			),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
		},
		// error - definition - invalid path key
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": "00000000-0000-0000-0000-000000000000",
						"display_name": "test",
						"format":       "Default",
						"definition": map[string]any{
							`"invalid.json"`: map[string]any{
								"source": "${local.path}/definition.json",
							},
						},
					},
				),
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute definition`),
		},
	}))
}

func TestUnit_OrgAppAudienceResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID))
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID))

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
		),
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

func TestUnit_OrgAppAudienceResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	parentAppID := testhelp.RandomUUID()
	entityExist := fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID)
	entityBefore := fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID)
	entityAfter := fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID)
	tag1ID := testhelp.RandomUUID()
	tag2ID := testhelp.RandomUUID()
	tag3ID := testhelp.RandomUUID()

	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID))
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID))

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
				),
			),
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
						"definition":   testHelperDefinitionTemplate(parentAppID, false, 1),
						"tags":         []string{tag1ID, tag2ID},
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "workspace_id", entityBefore.WorkspaceID),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "folder_id", entityBefore.FolderID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "format", "Default"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "tags.#", "2"),
				resource.TestCheckTypeSetElemAttr(testResourceItemFQN, "tags.*", tag1ID),
				resource.TestCheckTypeSetElemAttr(testResourceItemFQN, "tags.*", tag2ID),
			),
		},
		// Remove and add a new tag and Read
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
						"definition":   testHelperDefinitionTemplate(parentAppID, false, 1),
						"tags":         []string{tag2ID, tag3ID},
					},
				),
			),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(testResourceItemFQN, plancheck.ResourceActionUpdate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "tags.#", "2"),
				resource.TestCheckTypeSetElemAttr(testResourceItemFQN, "tags.*", tag2ID),
				resource.TestCheckTypeSetElemAttr(testResourceItemFQN, "tags.*", tag3ID),
			),
		},
		// Update (display name, description and definition content) and Read
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
						"format":       "Default",
						"definition":   testHelperDefinitionTemplate(parentAppID, true, 2),
					},
				),
			),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(testResourceItemFQN, plancheck.ResourceActionUpdate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "folder_id", entityBefore.FolderID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition.definition.json.tokens.TabOrder", "2"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition.definition.json.tokens.HasAccessToHiddenContent", "true"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "tags.#", "0"),
			),
		},
		// Update folder (forces replacement) and Read
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
						"definition":   testHelperDefinitionTemplate(parentAppID, true, 2),
					},
				),
			),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(testResourceItemFQN, plancheck.ResourceActionDestroyBeforeCreate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "folder_id", entityAfter.FolderID),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_OrgAppAudienceResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)
	folder1 := testhelp.WellKnown()["FolderRS1"].(map[string]any)
	folder1ID := folder1["id"].(string)
	folder2 := testhelp.WellKnown()["FolderRS2"].(map[string]any)
	folder2ID := folder2["id"].(string)

	orgAppResourceHCL, orgAppResourceFQN := testHelperParentOrgAppResource(workspaceID)
	definition := testHelperDefinitionTemplate(testhelp.RefByFQN(orgAppResourceFQN, "id"), false, 1)

	entityCreateDisplayName := testhelp.RandomName()
	entityUpdateDisplayName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				orgAppResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityCreateDisplayName,
						"folder_id":    folder1ID,
						"format":       "Default",
						"definition":   definition,
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttr(testResourceItemFQN, "folder_id", folder1ID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "format", "Default"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "tags.#", "0"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				orgAppResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityUpdateDisplayName,
						"description":  entityUpdateDescription,
						"folder_id":    folder1ID,
						"format":       "Default",
						"definition":   definition,
					},
				),
			),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(testResourceItemFQN, plancheck.ResourceActionUpdate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttr(testResourceItemFQN, "folder_id", folder1ID),
			),
		},
		// Update folder (forces replacement) and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				orgAppResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityUpdateDisplayName,
						"description":  entityUpdateDescription,
						"folder_id":    folder2ID,
						"format":       "Default",
						"definition":   definition,
					},
				),
			),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(testResourceItemFQN, plancheck.ResourceActionDestroyBeforeCreate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttr(testResourceItemFQN, "folder_id", folder2ID),
			),
		},
	},
	))
}

func TestAcc_OrgAppAudienceResource_CRUD_Definition(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	orgAppResourceHCL, orgAppResourceFQN := testHelperParentOrgAppResource(workspaceID)
	parentAppID := testhelp.RefByFQN(orgAppResourceFQN, "id")

	entityCreateDisplayName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				orgAppResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": workspaceID,
						"display_name": entityCreateDisplayName,
						"format":       "Default",
						"definition":   testHelperDefinitionTemplate(parentAppID, false, 1),
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttr(testResourceItemFQN, "format", "Default"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition.definition.json.tokens.TabOrder", "1"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition.definition.json.tokens.HasAccessToHiddenContent", "false"),
			),
		},
		// Update definition content and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				testHelperLocals,
				orgAppResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":              workspaceID,
						"display_name":              entityCreateDisplayName,
						"description":               entityUpdateDescription,
						"format":                    "Default",
						"definition_update_enabled": true,
						"definition":                testHelperDefinitionTemplate(parentAppID, true, 2),
					},
				),
			),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(testResourceItemFQN, plancheck.ResourceActionUpdate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition_update_enabled", "true"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition.definition.json.tokens.TabOrder", "2"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "definition.definition.json.tokens.HasAccessToHiddenContent", "true"),
			),
		},
	},
	))
}
