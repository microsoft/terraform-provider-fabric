// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package deploymentpipeline_test

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_DeploymentPipelineDataSource(t *testing.T) {
	entity := fakes.NewRandomDeploymentPipelineWithStages()

	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomDeploymentPipelineWithStages())
	fakes.FakeServer.Upsert(fakes.NewRandomDeploymentPipelineWithStages())

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Exactly one of these attributes must be configured: \[id,display_name\]`),
		},
		// error - conflicting attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id":           *entity.ID,
					"display_name": *entity.DisplayName,
				},
			),
			ExpectError: regexp.MustCompile(`These attributes cannot be configured together: \[id,display_name\]`),
		},
		// error - id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id":              *entity.ID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// read by id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": *entity.ID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "display_name", entity.DisplayName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "description", entity.Description),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "stages.0.display_name", entity.Stages[0].DisplayName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "stages.0.description", entity.Stages[0].Description),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.0.is_public", strconv.FormatBool(*entity.Stages[0].IsPublic)),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "stages.1.display_name", entity.Stages[1].DisplayName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "stages.1.description", entity.Stages[1].Description),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.1.is_public", strconv.FormatBool(*entity.Stages[1].IsPublic)),
			),
		},
		// read by id - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by name
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"display_name": *entity.DisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "display_name", entity.DisplayName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "description", entity.Description),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "stages.0.display_name", entity.Stages[0].DisplayName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "stages.0.description", entity.Stages[0].Description),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.0.is_public", strconv.FormatBool(*entity.Stages[0].IsPublic)),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "stages.1.display_name", entity.Stages[1].DisplayName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "stages.1.description", entity.Stages[1].Description),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.1.is_public", strconv.FormatBool(*entity.Stages[1].IsPublic)),
			),
		},
		// read by name - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"display_name": testhelp.RandomName(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}

func TestAcc_DeploymentPipelineDataSource(t *testing.T) {
	entity := testhelp.WellKnown()["DeploymentPipeline"].(map[string]any)
	entityID := entity["id"].(string)
	entityDisplayName := entity["displayName"].(string)
	entityDescription := entity["description"].(string)
	rawStagesAny := entity["stages"].([]any)
	// stage 1
	stage1Map := parseStageEntry(rawStagesAny[0].(string))
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
	entityStage3IsPublic := strings.EqualFold(stage3Map["isPublic"], "True")

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read by id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": entityID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", entityID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "display_name", entityDisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "description", entityDescription),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.0.display_name", entityStage1DisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.0.description", entityStage1Description),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.0.is_public", strconv.FormatBool(entityStage1IsPublic)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.1.display_name", entityStage2DisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.1.description", entityStage2Description),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.1.is_public", strconv.FormatBool(entityStage2IsPublic)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.2.display_name", entityStage3DisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.2.description", entityStage3Description),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.2.is_public", strconv.FormatBool(entityStage3IsPublic)),
			),
		},
		// read by id - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by name
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"display_name": entityDisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", entityID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "display_name", entityDisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "description", entityDescription),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.0.display_name", entityStage1DisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.0.description", entityStage1Description),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.0.is_public", strconv.FormatBool(entityStage1IsPublic)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.1.display_name", entityStage2DisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.1.description", entityStage2Description),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.1.is_public", strconv.FormatBool(entityStage2IsPublic)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.2.display_name", entityStage3DisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.2.description", entityStage3Description),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "stages.2.is_public", strconv.FormatBool(entityStage3IsPublic)),
			),
		},
		// read by name - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"display_name": testhelp.RandomName(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}
