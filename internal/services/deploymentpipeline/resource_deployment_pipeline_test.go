// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package deploymentpipeline_test

import (
	"errors"
	"regexp"
	"strconv"
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
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entity.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "stages.0.display_name", entity.Stages[0].DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "stages.0.description", entity.Stages[0].Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.is_public", strconv.FormatBool(*entity.Stages[0].IsPublic)),
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
		// error - create - existing entity
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
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.display_name", *entityBefore.Stages[0].DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.description", *entityBefore.Stages[0].Description),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.is_public", strconv.FormatBool(*entityBefore.Stages[0].IsPublic)),
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
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "description", entityAfter.Description),
			),
		},
	}))
}

func TestAcc_DeploymentPipelineResource_CRUD(t *testing.T) {
	entityCreateDisplayName := testhelp.RandomName()
	entityUpdateDisplayName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()

	entityStage1Name := testhelp.RandomName()
	entityStage1IsPublic := testhelp.RandomBool()
	entityStage2Name := testhelp.RandomName()
	entityStage2IsPublic := testhelp.RandomBool()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
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
							"is_public":    entityStage1IsPublic,
						},
						{
							"display_name": entityStage2Name,
							"description":  entityStage2Name,
							"is_public":    entityStage2IsPublic,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.display_name", entityStage1Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.description", entityStage1Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.is_public", strconv.FormatBool(entityStage1IsPublic)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(entityStage2IsPublic)),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": entityUpdateDisplayName,
					"description":  entityUpdateDescription,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
			),
		},
	},
	))
}
