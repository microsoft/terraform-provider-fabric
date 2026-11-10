// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package ontology_test

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
	"path": testhelp.GetFixturesDirPath("ontology"),
})

func TestUnit_OntologyResource_Attributes(t *testing.T) {
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
						"unexpected_attr": "test",
					},
				),
			),
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
					},
				),
			),
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
					},
				),
			),
			ExpectError: regexp.MustCompile(`The argument "display_name" is required, but no definition was found.`),
		},
	}))
}

func TestUnit_OntologyResource_ImportState(t *testing.T) {
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

func TestUnit_OntologyResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entityExist := fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID)
	entityBefore := fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID)
	entityAfter := fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID)

	fakes.FakeServer.Upsert(fakes.NewRandomItemWithWorkspace(fabricItemType, workspaceID))
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(entityAfter)
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
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
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
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

func TestAcc_OntologyResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	entityCreateDisplayName := testhelp.RandomName()
	entityUpdateDisplayName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()
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
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
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
						"display_name": entityUpdateDisplayName,
						"description":  entityUpdateDescription,
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
			),
		},
	},
	))
}

func TestAcc_OntologyDefinitionResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil,
		ontologyDefinitionSteps(workspaceID, testhelp.RandomName()),
	))
}

func TestUnit_OntologyDefinitionResource_CRUD(t *testing.T) {
	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil,
		ontologyDefinitionSteps(testhelp.RandomUUID(), testhelp.RandomName()),
	))
}

func ontologyDefinitionSteps(workspaceID, displayName string) []resource.TestStep {
	config := func(namespace string, updateEnabled bool) string {
		return at.JoinConfigs(
			testHelperLocals,
			at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":              workspaceID,
					"display_name":              displayName,
					"format":                    "TMDL",
					"definition_update_enabled": updateEnabled,
					"definition": map[string]any{
						`".platform"`: map[string]any{
							"source": "${local.path}/platform.json",
						},
						`"database.tmdl"`: map[string]any{
							"source": "${local.path}/database.tmdl",
						},
						`"model.tmdl"`: map[string]any{
							"source": "${local.path}/model.tmdl.tmpl",
							"tokens": map[string]any{"Namespace": namespace},
						},
						fmt.Sprintf(`"namespaces/%s.tmdl"`, namespace): map[string]any{
							"source": "${local.path}/namespace.tmdl.tmpl",
							"tokens": map[string]any{"Namespace": namespace},
						},
					},
				},
			),
			fmt.Sprintf(`
data "fabric_ontology" "definition" {
  workspace_id      = %q
  id                = %s.id
  format            = "TMDL"
  output_definition = true
  depends_on        = [%s]
}
output "database_tmdl" {
  value = trimspace(provider::fabric::content_decode(data.fabric_ontology.definition.definition["database.tmdl"].content))
}
output "model_tmdl" {
  value = trimspace(provider::fabric::content_decode(data.fabric_ontology.definition.definition["model.tmdl"].content))
}
`, workspaceID, testResourceItemFQN, testResourceItemFQN),
		)
	}

	check := func(namespace string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", displayName),
			resource.TestCheckResourceAttr(testResourceItemFQN, "format", "TMDL"),
			resource.TestMatchOutput("database_tmdl", regexp.MustCompile(`database\r?\n\tcompatibilityLevel: 1000000`)),
			resource.TestMatchOutput("model_tmdl", regexp.MustCompile(`ref namespace `+regexp.QuoteMeta(namespace))),
			resource.TestCheckResourceAttrSet("data.fabric_ontology.definition", "definition.namespaces/"+namespace+".tmdl.content"),
		)
	}

	return []resource.TestStep{
		{
			Config: config("operations", false),
			Check:  check("operations"),
		},
		{
			Config: config("maintenance", false),
			Check:  check("operations"),
		},
		{
			Config: config("maintenance", true),
			Check:  check("operations"),
		},
		{
			Config: config("production", true),
			Check:  check("production"),
		},
		{
			Config: config("manufacturing", false),
			Check:  check("production"),
		},
		{
			Config: config("logistics", true),
			Check:  check("logistics"),
		},
		{
			Config: at.JoinConfigs(
				at.CompileConfig(testResourceItemHeader, map[string]any{
					"workspace_id": workspaceID,
					"display_name": displayName,
					"format":       "TMDL",
				}),
				fmt.Sprintf(`
data "fabric_ontology" "definition" {
  workspace_id      = %q
  id                = %s.id
  format            = "TMDL"
  output_definition = true
  depends_on        = [%s]
}
output "database_tmdl" {
  value = trimspace(provider::fabric::content_decode(data.fabric_ontology.definition.definition["database.tmdl"].content))
}
`, workspaceID, testResourceItemFQN, testResourceItemFQN),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestMatchOutput("database_tmdl", regexp.MustCompile(`database\r?\n\tcompatibilityLevel: 1000000`)),
				resource.TestCheckNoResourceAttr("data.fabric_ontology.definition", "definition.namespaces/logistics.tmdl.content"),
			),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(testResourceItemHeader, map[string]any{
				"workspace_id": workspaceID,
				"display_name": displayName,
			}),
			ImportState:             true,
			ImportStateIdPrefix:     workspaceID + "/",
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"definition", "format"},
		},
	}
}

func TestUnit_OntologyResource_DefinitionPaths(t *testing.T) {
	validPaths := []string{
		".platform", "database.tmdl", "model.tmdl", "tables/Sensor.tmdl",
		"metrics/Temperature.tmdl", "relationships.tmdl", "expressions.tmdl",
		"entities/factory#Sensor.tmdl", "entityRelationships.tmdl",
		"namespaces/default.tmdl", "rules/SpeedLimit.tmdl",
	}
	definition := make(map[string]any, len(validPaths))

	for _, part := range validPaths {
		definition[fmt.Sprintf("%q", part)] = map[string]any{"source": "${local.path}/database.tmdl"}
	}

	invalidPaths := []string{
		"definition.json", "EntityTypes/123/definition.json", "RelationshipTypes/123",
		"tables/Sensor.json", "definition/model.tmdl", "unknown.tmdl",
	}
	steps := make([]resource.TestStep, 0, len(invalidPaths)+2)
	steps = append(steps, resource.TestStep{
		PlanOnly:           true,
		ExpectNonEmptyPlan: true,
		Config: at.JoinConfigs(testHelperLocals, at.CompileConfig(testResourceItemHeader, map[string]any{
			"workspace_id": testhelp.RandomUUID(),
			"display_name": "OntologyTMDL",
			"format":       "TMDL",
			"definition":   definition,
		})),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(testResourceItemFQN, "format", "TMDL"),
		),
	})

	for _, part := range invalidPaths {
		steps = append(steps, resource.TestStep{
			PlanOnly: true,
			Config: at.JoinConfigs(testHelperLocals, at.CompileConfig(testResourceItemHeader, map[string]any{
				"workspace_id": testhelp.RandomUUID(),
				"display_name": "OntologyTMDL",
				"format":       "TMDL",
				"definition": map[string]any{
					fmt.Sprintf("%q", part): map[string]any{"source": "${local.path}/database.tmdl"},
				},
			})),
			ExpectError: regexp.MustCompile(`Definition path must\s+match`),
		})
	}

	steps = append(steps, resource.TestStep{
		PlanOnly: true,
		Config: at.CompileConfig(testResourceItemHeader, map[string]any{
			"workspace_id": testhelp.RandomUUID(),
			"display_name": "OntologyTMDL",
			"format":       "Default",
		}),
		ExpectError: regexp.MustCompile(`value must be one of`),
	})

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, steps))
}
