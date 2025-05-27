// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package deploymentpipelinera_test

import (
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func fakeDeploymentPipelineRoleAssignments(
	exampleResp fabcore.DeploymentPipelineRoleAssignments,
) func(deploymentPipelineID string, options *fabcore.DeploymentPipelinesClientListDeploymentPipelineRoleAssignmentsOptions) (resp azfake.PagerResponder[fabcore.DeploymentPipelinesClientListDeploymentPipelineRoleAssignmentsResponse]) {
	return func(_ string, _ *fabcore.DeploymentPipelinesClientListDeploymentPipelineRoleAssignmentsOptions) (resp azfake.PagerResponder[fabcore.DeploymentPipelinesClientListDeploymentPipelineRoleAssignmentsResponse]) {
		resp = azfake.PagerResponder[fabcore.DeploymentPipelinesClientListDeploymentPipelineRoleAssignmentsResponse]{}
		resp.AddPage(http.StatusOK, fabcore.DeploymentPipelinesClientListDeploymentPipelineRoleAssignmentsResponse{DeploymentPipelineRoleAssignments: exampleResp}, nil)

		return
	}
}

func NewRandomDeploymentPipelineRoleAssignments() fabcore.DeploymentPipelineRoleAssignments {
	principal0ID := testhelp.RandomUUID()
	principal1ID := testhelp.RandomUUID()
	principal2ID := testhelp.RandomUUID()

	return fabcore.DeploymentPipelineRoleAssignments{
		Value: []fabcore.DeploymentPipelineRoleAssignment{
			{
				Role: azto.Ptr(fabcore.DeploymentPipelineRoleAdmin),
				Principal: &fabcore.Principal{
					ID:   azto.Ptr(principal0ID),
					Type: azto.Ptr(fabcore.PrincipalTypeGroup),
				},
			},
			{
				Role: azto.Ptr(fabcore.DeploymentPipelineRoleAdmin),
				Principal: &fabcore.Principal{
					ID:   azto.Ptr(principal1ID),
					Type: azto.Ptr(fabcore.PrincipalTypeUser),
				},
			},
			{
				Role: azto.Ptr(fabcore.DeploymentPipelineRoleAdmin),
				Principal: &fabcore.Principal{
					ID:   azto.Ptr(principal2ID),
					Type: azto.Ptr(fabcore.PrincipalTypeServicePrincipal),
				},
			},
		},
	}
}
