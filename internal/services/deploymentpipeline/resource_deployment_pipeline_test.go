// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package deploymentpipeline_test

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
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

func TestAcc_DeploymentPipelineResource_CRUD(t *testing.T) {
	// entity := testhelp.WellKnown()["DeploymentPipeline"].(map[string]any)
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceDS := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)
	workspaceID_DS := workspaceDS["id"].(string)

	//rawStagesAny := entity["stages"].([]any)
	// stage 1
	/*stage1Map := parseStageEntry(rawStagesAny[0].(string))
	entityStage1DisplayName := stage1Map["displayName"]
	entityStage1Description := stage1Map["description"]
	entityStage1IsPublic := strings.EqualFold(stage1Map["isPublic"], "True")
	// stage 2
	stage2Map := parseStageEntry(rawStagesAny[1].(string))
	entityStage2DisplayName := stage2Map["displayName"]
	entityStage2Description := stage2Map["description"]
	entityStage2IsPublic := strings.EqualFold(stage2Map["isPublic"], "True")
	// stage 3
	stage3Map := parseStageEntry(rawStagesAny[2].(string))
	entityStage3DisplayName := stage3Map["displayName"]
	entityStage3Description := stage3Map["description"]
	entityStage3IsPublic := strings.EqualFold(stage3Map["isPublic"], "True")*/

	entityUpdateDisplayName := testhelp.RandomName()
	entityUpdateDescription := testhelp.RandomName()
	entityCreateDisplayName := testhelp.RandomName()

	entityStage1Name := testhelp.RandomName()
	entityStage1IsPublicRandom := testhelp.RandomBool()
	entityStage2Name := testhelp.RandomName()
	entityStage2IsPublicRandom := testhelp.RandomBool()

	entityStage1NameRandom := testhelp.RandomName()

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
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.order", "0"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.0.workspace_id"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.0.workspace_name"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(entityStage2IsPublicRandom)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.order", "1"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.1.workspace_id"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.1.workspace_name"),
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
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.order", "0"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.0.workspace_id"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.0.workspace_name"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(entityStage2IsPublicRandom)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.order", "1"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.1.workspace_id"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.1.workspace_name"),
			),
		},
		// Update and Read - assign stage 1 to workspace
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
							"workspace_id": workspaceID,
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
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.order", "0"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(entityStage2IsPublicRandom)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.order", "1"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.1.workspace_id"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.1.workspace_name"),
			),
		},
		// Update and Read - unassign stage 1 assign stage 2
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
							"workspace_id": workspaceID,
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
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.order", "0"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.0.workspace_id"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(entityStage2IsPublicRandom)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.order", "1"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.workspace_id", workspaceID),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "stages.1.workspace_name"),
			),
		},
		// Update and Read - update stage name
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name": entityUpdateDisplayName,
					"description":  entityUpdateDescription,
					"stages": []map[string]any{
						{
							"display_name": entityStage1NameRandom,
							"description":  entityStage1NameRandom,
							"is_public":    entityStage1IsPublicRandom,
							"workspace_id": workspaceID_DS,
						},
						{
							"display_name": entityStage2Name,
							"description":  entityStage2Name,
							"is_public":    entityStage2IsPublicRandom,
							"workspace_id": workspaceID,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.display_name", entityStage1NameRandom),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.description", entityStage1NameRandom),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.is_public", strconv.FormatBool(entityStage1IsPublicRandom)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.order", "0"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.0.workspace_id", workspaceID_DS),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.display_name", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.description", entityStage2Name),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.is_public", strconv.FormatBool(entityStage2IsPublicRandom)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.order", "1"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "stages.1.workspace_id", workspaceID),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "stages.1.workspace_name"),
			),
		},
	},
	))
}

func parseStageEntry(raw string) map[string]string {
	m := map[string]string{}
	// strip the "@{" prefix and "}" suffix
	s := strings.TrimPrefix(raw, "@{")
	s = strings.TrimSuffix(s, "}")

	for _, pair := range strings.Split(s, "; ") {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			m[kv[0]] = kv[1]
		}
	}

	return m
}
