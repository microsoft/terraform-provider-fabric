// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package deploymentpipelinera_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemsFQN, testDataSourceItemsHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Types, "test")

func TestUnit_DeploymentPipelineRoleAssignmentsDataSource(t *testing.T) {
	deploymentPipelineID := testhelp.RandomUUID()

	deploymentPipelineRoleAssignments := NewRandomDeploymentPipelineRoleAssignments()
	fakes.FakeServer.ServerFactory.Core.DeploymentPipelinesServer.NewListDeploymentPipelineRoleAssignmentsPager = fakeListDeploymentPipelineRoleAssignments(
		deploymentPipelineRoleAssignments,
	)

	entity := deploymentPipelineRoleAssignments.Value[0]

	resource.Test(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"deployment_pipeline_id": deploymentPipelineID,
					"unexpected_attr":        "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"deployment_pipeline_id": deploymentPipelineID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "deployment_pipeline_id", deploymentPipelineID),
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "values.0.principal.id", *entity.ID),
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "values.0.principal.type", string(*entity.Principal.Type)),
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "values.0.role", string(fabcore.DeploymentPipelineRoleAdmin)),
			),
		},
	}))
}

func TestAcc_DeploymentPipelineRoleAssignmentsDataSource(t *testing.T) {
	deploymentPipeline := testhelp.WellKnown()["DeploymentPipeline"].(map[string]any)
	deploymentPipelineID := deploymentPipeline["id"].(string)

	group := testhelp.WellKnown()["Group"].(map[string]any)
	groupType := group["type"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"deployment_pipeline_id": deploymentPipelineID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "deployment_pipeline_id", deploymentPipelineID),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.id"),
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "values.0.principal.type", groupType),
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "values.0.role", string(fabcore.DeploymentPipelineRoleAdmin)),
			),
		},
	},
	))
}
