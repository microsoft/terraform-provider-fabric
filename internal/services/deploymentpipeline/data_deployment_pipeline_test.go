// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package deploymentpipeline_test

import (
	"regexp"
	"strconv"
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
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// missing required "id" attribute
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"display_name": *entity.DisplayName,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found.`),
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
			),
		},
	}))
}

func TestAcc_DeploymentPipelineDataSource(t *testing.T) {
	entity := testhelp.WellKnown()["DeploymentPipeline"].(map[string]any)
	entityID := entity["id"].(string)
	entityDisplayName := entity["displayName"].(string)
	entityDescription := entity["description"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
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
	}))
}
