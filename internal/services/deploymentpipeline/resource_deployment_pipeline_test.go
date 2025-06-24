// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package deploymentpipeline_test

import (
	"errors"
	"regexp"
	"strconv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_DeploymentPipelineResource_Attributes(t *testing.T) {
	entity := fakes.NewRandomDeploymentPipelineWithStages()
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
		// error - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":    "test",
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attribute - display_name
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"description": "test",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "display_name" is required, but no definition was found.`),
		},
		// error - no required attribute - stages
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": *entity.DisplayName,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "stages" is required, but no definition was found.`),
		},
		// success
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": *entity.DisplayName,
					"stages": []map[string]any{
						{
							"display_name": *entity.Stages[0].DisplayName,
							"description":  *entity.Stages[0].Description,
							"is_public":    *entity.Stages[0].IsPublic,
						},
						{
							"display_name": *entity.Stages[1].DisplayName,
							"description":  *entity.Stages[1].Description,
							"is_public":    *entity.Stages[1].IsPublic,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entity.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "stages.0.display_name", entity.Stages[0].DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "stages.0.description", entity.Stages[0].Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.is_public", strconv.FormatBool(*entity.Stages[0].IsPublic)),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "stages.1.display_name", entity.Stages[1].DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "stages.1.description", entity.Stages[1].Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(*entity.Stages[1].IsPublic)),
			),
		},
	}))
}

func TestUnit_DeploymentPipelineResource_ImportState(t *testing.T) {
	entity := fakes.NewRandomDeploymentPipelineWithStages()

	fakes.FakeServer.Upsert(fakes.NewRandomDeploymentPipelineWithStages())
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomDeploymentPipelineWithStages())

	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"display_name": *entity.DisplayName,
			"description":  *entity.Description,
			"stages": []map[string]any{
				{
					"display_name": *entity.Stages[0].DisplayName,
					"description":  *entity.Stages[0].Description,
					"is_public":    *entity.Stages[0].IsPublic,
				},
				{
					"display_name": *entity.Stages[1].DisplayName,
					"description":  *entity.Stages[1].Description,
					"is_public":    *entity.Stages[1].IsPublic,
				},
			},
		},
	)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// Import state testing
		{
			ResourceName:       testResourceItemFQN,
			Config:             testCase,
			ImportStateId:      *entity.ID,
			ImportState:        true,
			ImportStatePersist: true,
			ImportStateCheck: func(is []*terraform.InstanceState) error {
				if len(is) != 1 {
					return errors.New("expected one instance state")
				}

				if is[0].ID != *entity.ID {
					return errors.New(testResourceItemFQN + ": unexpected ID")
				}

				if is[0].Attributes["display_name"] != *entity.DisplayName {
					return errors.New(testResourceItemFQN + ": unexpected display_name")
				}

				if is[0].Attributes["description"] != *entity.Description {
					return errors.New(testResourceItemFQN + ": unexpected description")
				}

				return nil
			},
		},
	}))
}

func TestUnit_DeploymentPipelineResource_CRUD(t *testing.T) {
	entityExist := fakes.NewRandomDeploymentPipelineWithStages()
	entityBefore := fakes.NewRandomDeploymentPipelineWithStages()
	entityAfter := fakes.NewRandomDeploymentPipelineWithStages()

	fakes.FakeServer.Upsert(fakes.NewRandomDeploymentPipelineWithStages())
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(fakes.NewRandomDeploymentPipelineWithStages())

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - stages should be between 2 and 10
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": *entityExist.DisplayName,
					"stages": []map[string]any{
						{
							"display_name": *entityExist.Stages[0].DisplayName,
							"description":  *entityExist.Stages[0].Description,
							"is_public":    *entityExist.Stages[0].IsPublic,
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`stages list must contain at least 2 elements and at most 10`),
		},
		// error - create - existing entity
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": *entityExist.DisplayName,
					"stages": []map[string]any{
						{
							"display_name": *entityBefore.Stages[0].DisplayName,
							"description":  *entityBefore.Stages[0].Description,
							"is_public":    *entityBefore.Stages[0].IsPublic,
						},
						{
							"display_name": *entityBefore.Stages[1].DisplayName,
							"description":  *entityBefore.Stages[1].Description,
							"is_public":    *entityBefore.Stages[1].IsPublic,
						},
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
					"display_name": *entityBefore.DisplayName,
					"stages": []map[string]any{
						{
							"display_name": *entityBefore.Stages[0].DisplayName,
							"description":  *entityBefore.Stages[0].Description,
							"is_public":    *entityBefore.Stages[0].IsPublic,
						},
						{
							"display_name": *entityBefore.Stages[1].DisplayName,
							"description":  *entityBefore.Stages[1].Description,
							"is_public":    *entityBefore.Stages[1].IsPublic,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.display_name", *entityBefore.Stages[0].DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.description", *entityBefore.Stages[0].Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.is_public", strconv.FormatBool(*entityBefore.Stages[0].IsPublic)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", *entityBefore.Stages[1].DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", *entityBefore.Stages[1].Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(*entityBefore.Stages[1].IsPublic)),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": *entityBefore.DisplayName,
					"description":  *entityAfter.Description,
					"stages": []map[string]any{
						{
							"display_name": *entityBefore.Stages[0].DisplayName,
							"description":  *entityBefore.Stages[0].Description,
							"is_public":    *entityBefore.Stages[0].IsPublic,
						},
						{
							"display_name": *entityBefore.Stages[1].DisplayName,
							"description":  *entityBefore.Stages[1].Description,
							"is_public":    *entityBefore.Stages[1].IsPublic,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
			),
		},
	}))
}

func TestUnit_DeploymentPipelineResource_CRUD_Stages(t *testing.T) {
	entityExist := fakes.NewRandomDeploymentPipelineWithStages()
	entityBefore := fakes.NewRandomDeploymentPipelineWithStages()
	entityAfter := fakes.NewRandomDeploymentPipelineWithStages()

	fakes.FakeServer.Upsert(fakes.NewRandomDeploymentPipelineWithStages())
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(fakes.NewRandomDeploymentPipelineWithStages())

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Update and Read - add new stage
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": *entityBefore.DisplayName,
					"description":  *entityAfter.Description,
					"stages": []map[string]any{
						{
							"display_name": *entityBefore.Stages[0].DisplayName,
							"description":  *entityBefore.Stages[0].Description,
							"is_public":    *entityBefore.Stages[0].IsPublic,
						},
						{
							"display_name": *entityBefore.Stages[1].DisplayName,
							"description":  *entityBefore.Stages[1].Description,
							"is_public":    *entityBefore.Stages[1].IsPublic,
						},
						{
							"display_name": *entityAfter.Stages[0].DisplayName,
							"description":  *entityAfter.Stages[0].Description,
							"is_public":    *entityAfter.Stages[0].IsPublic,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.display_name", *entityBefore.Stages[0].DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.description", *entityBefore.Stages[0].Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.is_public", strconv.FormatBool(*entityBefore.Stages[0].IsPublic)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", *entityBefore.Stages[1].DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", *entityBefore.Stages[1].Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(*entityBefore.Stages[1].IsPublic)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.2.display_name", *entityAfter.Stages[0].DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.2.description", *entityAfter.Stages[0].Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.2.is_public", strconv.FormatBool(*entityAfter.Stages[0].IsPublic)),
			),
		},
		// Update and Read - remove new stage
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": *entityBefore.DisplayName,
					"description":  *entityAfter.Description,
					"stages": []map[string]any{
						{
							"display_name": *entityBefore.Stages[0].DisplayName,
							"description":  *entityBefore.Stages[0].Description,
							"is_public":    *entityBefore.Stages[0].IsPublic,
						},
						{
							"display_name": *entityBefore.Stages[1].DisplayName,
							"description":  *entityBefore.Stages[1].Description,
							"is_public":    *entityBefore.Stages[1].IsPublic,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.display_name", *entityBefore.Stages[0].DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.description", *entityBefore.Stages[0].Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.is_public", strconv.FormatBool(*entityBefore.Stages[0].IsPublic)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", *entityBefore.Stages[1].DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", *entityBefore.Stages[1].Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(*entityBefore.Stages[1].IsPublic)),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.2.display_name"),
			),
		},
	}))
}

func TestAcc_DeploymentPipelineResource_CRUD(t *testing.T) {
	entity := testhelp.WellKnown()["DeploymentPipeline"].(map[string]any)
	entityName := entity["displayName"].(string)
	entityDescription := entity["description"].(string)

	entityStages := entity["stages"].([]any)
	// stage 1
	stage1Map := entityStages[0].(map[string]any)
	entityStage1DisplayName := stage1Map["displayName"]
	entityStage1Description := stage1Map["description"]
	entityStage1IsPublic := stage1Map["isPublic"]
	// stage 2
	stage2Map := entityStages[1].(map[string]any)
	entityStage2DisplayName := stage2Map["displayName"]
	entityStage2Description := stage2Map["description"]
	entityStage2IsPublic := stage2Map["isPublic"]
	// stage 3
	stage3Map := entityStages[2].(map[string]any)
	entityStage3DisplayName := stage3Map["displayName"]
	entityStage3Description := stage3Map["description"]
	entityStage3IsPublic := stage3Map["isPublic"]

	entityUpdateDisplayName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()
	entityCreateDisplayName := testhelp.RandomName()

	entityStage1Name := testhelp.RandomName()
	entityStage1IsPublicRandom := testhelp.RandomBool()
	entityStage2Name := testhelp.RandomName()
	entityStage2IsPublicRandom := testhelp.RandomBool()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// error - create - existing entity
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": entityName,
					"description":  entityDescription,
					"stages": []map[string]any{
						{
							"display_name": entityStage1DisplayName,
							"description":  entityStage1Description,
							"is_public":    entityStage1IsPublic,
						},
						{
							"display_name": entityStage2DisplayName,
							"description":  entityStage2Description,
							"is_public":    entityStage2IsPublic,
						},
						{
							"display_name": entityStage3DisplayName,
							"description":  entityStage3Description,
							"is_public":    entityStage3IsPublic,
						},
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
					"display_name": entityCreateDisplayName,
					"stages": []map[string]any{
						{
							"display_name": entityStage1Name,
							"description":  entityStage1Name,
							"is_public":    entityStage1IsPublicRandom,
						},
						{
							"display_name": entityStage2Name,
							"description":  entityStage2Name,
							"is_public":    entityStage2IsPublicRandom,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.display_name", entityStage1Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.description", entityStage1Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.is_public", strconv.FormatBool(entityStage1IsPublicRandom)),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.0.workspace_id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(entityStage2IsPublicRandom)),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.1.workspace_id"),
			),
		},
		// Update and Read - update Deployment Pipeline Name and Description
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": entityUpdateDisplayName,
					"description":  entityUpdateDescription,
					"stages": []map[string]any{
						{
							"display_name": entityStage1Name,
							"description":  entityStage1Name,
							"is_public":    entityStage1IsPublicRandom,
						},
						{
							"display_name": entityStage2Name,
							"description":  entityStage2Name,
							"is_public":    entityStage2IsPublicRandom,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.display_name", entityStage1Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.description", entityStage1Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.is_public", strconv.FormatBool(entityStage1IsPublicRandom)),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.0.workspace_id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(entityStage2IsPublicRandom)),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.1.workspace_id"),
			),
		},
	},
	))
}

func TestAcc_DeploymentPipelineResource_CRUD_Stages(t *testing.T) {
	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, capacityID)
	workspaceID := testhelp.RefByFQN(workspaceResourceFQN, "id")
	entityCreateDisplayName := testhelp.RandomName()
	entityCreateDescription := testhelp.RandomName()

	entityStage1Name := testhelp.RandomName()
	entityStage1Description := testhelp.RandomName()
	entityStage1IsPublicRandom := testhelp.RandomBool()

	entityStage2Name := testhelp.RandomName()
	entityStage2Description := testhelp.RandomName()
	entityStage2IsPublicRandom := testhelp.RandomBool()

	entityStage1UpdateName := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read - with stage assignment
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(workspaceResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"display_name": entityCreateDisplayName,
						"description":  entityCreateDescription,
						"stages": []map[string]any{
							{
								"display_name": entityStage1Name,
								"description":  entityStage1Description,
								"is_public":    entityStage1IsPublicRandom,
								"workspace_id": workspaceID,
							},
							{
								"display_name": entityStage2Name,
								"description":  entityStage2Description,
								"is_public":    entityStage2IsPublicRandom,
							},
						},
					}),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityCreateDescription),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.display_name", entityStage1Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.description", entityStage1Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.is_public", strconv.FormatBool(entityStage1IsPublicRandom)),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "stages.0.workspace_id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", entityStage2Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(entityStage2IsPublicRandom)),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.1.workspace_id"),
			),
		},
		// Update and Read - update stage 1 and unassign to workspace
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(workspaceResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"display_name": entityCreateDisplayName,
						"description":  entityCreateDescription,
						"stages": []map[string]any{
							{
								"display_name": entityStage1UpdateName,
								"description":  entityStage1Description,
								"is_public":    entityStage1IsPublicRandom,
							},
							{
								"display_name": entityStage2Name,
								"description":  entityStage2Description,
								"is_public":    entityStage2IsPublicRandom,
							},
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityCreateDescription),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.display_name", entityStage1UpdateName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.description", entityStage1Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.is_public", strconv.FormatBool(entityStage1IsPublicRandom)),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.0.workspace_id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", entityStage2Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(entityStage2IsPublicRandom)),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.1.workspace_id"),
			),
		},
	}))
}

func TestAcc_DeploymentPipelineResource_CRUD_DPAndStages(t *testing.T) {
	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, capacityID)
	entityCreateDisplayName := testhelp.RandomName()
	entityCreateDescription := testhelp.RandomName()

	entityUpdateDisplayName := testhelp.RandomName()

	entityStage1Name := testhelp.RandomName()
	entityStage1Description := testhelp.RandomName()
	entityStage1IsPublicRandom := testhelp.RandomBool()

	entityStage2Name := testhelp.RandomName()
	entityStage2Description := testhelp.RandomName()
	entityStage2IsPublicRandom := testhelp.RandomBool()

	entityStage3Name := testhelp.RandomName()
	entityStage3Description := testhelp.RandomName()
	entityStage3IsPublicRandom := testhelp.RandomBool()

	entityStage1UpdateName := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": entityCreateDisplayName,
					"description":  entityCreateDescription,
					"stages": []map[string]any{
						{
							"display_name": entityStage1Name,
							"description":  entityStage1Description,
							"is_public":    entityStage1IsPublicRandom,
						},
						{
							"display_name": entityStage2Name,
							"description":  entityStage2Description,
							"is_public":    entityStage2IsPublicRandom,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityCreateDescription),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.display_name", entityStage1Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.description", entityStage1Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.is_public", strconv.FormatBool(entityStage1IsPublicRandom)),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.0.workspace_id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", entityStage2Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(entityStage2IsPublicRandom)),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.1.workspace_id"),
			),
		},
		// Update and Read - update stage 1 and unassign to workspace
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(workspaceResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"display_name": entityUpdateDisplayName,
						"description":  entityUpdateDisplayName,
						"stages": []map[string]any{
							{
								"display_name": entityStage1UpdateName,
								"description":  entityStage1Description,
								"is_public":    entityStage1IsPublicRandom,
								"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
							},
							{
								"display_name": entityStage2Name,
								"description":  entityStage2Description,
								"is_public":    entityStage2IsPublicRandom,
							},
						},
					}),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.display_name", entityStage1UpdateName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.description", entityStage1Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.is_public", strconv.FormatBool(entityStage1IsPublicRandom)),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "stages.0.workspace_id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", entityStage2Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(entityStage2IsPublicRandom)),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.1.workspace_id"),
			),
		},
		// Update and Read - unassign stage
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(workspaceResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"display_name": entityUpdateDisplayName,
						"description":  entityUpdateDisplayName,
						"stages": []map[string]any{
							{
								"display_name": entityStage1UpdateName,
								"description":  entityStage1Description,
								"is_public":    entityStage1IsPublicRandom,
							},
							{
								"display_name": entityStage2Name,
								"description":  entityStage2Description,
								"is_public":    entityStage2IsPublicRandom,
							},
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.display_name", entityStage1UpdateName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.description", entityStage1Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.is_public", strconv.FormatBool(entityStage1IsPublicRandom)),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.0.workspace_id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", entityStage2Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(entityStage2IsPublicRandom)),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.1.workspace_id"),
			),
		},
		// Update and Read - add new stage
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(workspaceResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"display_name": entityUpdateDisplayName,
						"description":  entityUpdateDisplayName,
						"stages": []map[string]any{
							{
								"display_name": entityStage1UpdateName,
								"description":  entityStage1Description,
								"is_public":    entityStage1IsPublicRandom,
							},
							{
								"display_name": entityStage2Name,
								"description":  entityStage2Description,
								"is_public":    entityStage2IsPublicRandom,
							},
							{
								"display_name": entityStage3Name,
								"description":  entityStage3Description,
								"is_public":    entityStage3IsPublicRandom,
							},
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.display_name", entityStage1UpdateName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.description", entityStage1Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.is_public", strconv.FormatBool(entityStage1IsPublicRandom)),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.0.workspace_id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", entityStage2Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(entityStage2IsPublicRandom)),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.1.workspace_id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.2.display_name", entityStage3Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.2.description", entityStage3Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.2.is_public", strconv.FormatBool(entityStage3IsPublicRandom)),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.2.workspace_id"),
			),
		},
	}))
}

func TestUnit_DeploymentPipelineResource_CRUD_Stage_WorkspaceAssignment(t *testing.T) {
	testState := testhelp.NewTestState()
	workspaceID := testhelp.RandomUUID()

	entityBefore := fakes.NewRandomDeploymentPipelineWithStages()

	entityWithWorkspaceAssigned := entityBefore
	entityWithWorkspaceAssigned.Stages[0].WorkspaceID = to.Ptr(workspaceID)

	fakes.FakeServer.ServerFactory.Core.DeploymentPipelinesServer.AssignWorkspaceToStage = fakeWorkspaceAssignmentStage()
	fakes.FakeServer.ServerFactory.Core.DeploymentPipelinesServer.UnassignWorkspaceFromStage = fakeWorkspaceUnassignmentStage()

	preFakeGet := fakes.FakeServer.ServerFactory.Core.DeploymentPipelinesServer.GetDeploymentPipeline

	fakes.FakeServer.ServerFactory.Core.DeploymentPipelinesServer.GetDeploymentPipeline = fakeGetDeploymentPipeline(entityWithWorkspaceAssigned)

	resource.Test(t, testhelp.NewTestUnitCaseWithState(t, nil, fakes.FakeServer.ServerFactory, testState, nil, []resource.TestStep{
		// Update and Read - assign stage
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": *entityBefore.DisplayName,
					"description":  *entityBefore.Description,
					"stages": []map[string]any{
						{
							"display_name": *entityBefore.Stages[0].DisplayName,
							"description":  *entityBefore.Stages[0].Description,
							"is_public":    *entityBefore.Stages[0].IsPublic,
							"workspace_id": workspaceID,
						},
						{
							"display_name": *entityBefore.Stages[1].DisplayName,
							"description":  *entityBefore.Stages[1].Description,
							"is_public":    *entityBefore.Stages[1].IsPublic,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityBefore.Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.display_name", *entityBefore.Stages[0].DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.description", *entityBefore.Stages[0].Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.is_public", strconv.FormatBool(*entityBefore.Stages[0].IsPublic)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", *entityBefore.Stages[1].DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", *entityBefore.Stages[1].Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(*entityBefore.Stages[1].IsPublic)),
			),
		},
	}))

	fakes.FakeServer.ServerFactory.Core.DeploymentPipelinesServer.GetDeploymentPipeline = preFakeGet
}
