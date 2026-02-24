// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package deploymentpipelinera_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// Returns a fake pager function that simulates listing deployment pipeline role assignments with a provided example response.
func fakeListDeploymentPipelineRoleAssignments(
	exampleResp fabcore.DeploymentPipelineRoleAssignments,
) func(deploymentPipelineID string, options *fabcore.DeploymentPipelinesClientListDeploymentPipelineRoleAssignmentsOptions) (resp azfake.PagerResponder[fabcore.DeploymentPipelinesClientListDeploymentPipelineRoleAssignmentsResponse]) {
	return func(_ string, _ *fabcore.DeploymentPipelinesClientListDeploymentPipelineRoleAssignmentsOptions) (resp azfake.PagerResponder[fabcore.DeploymentPipelinesClientListDeploymentPipelineRoleAssignmentsResponse]) {
		resp = azfake.PagerResponder[fabcore.DeploymentPipelinesClientListDeploymentPipelineRoleAssignmentsResponse]{}
		resp.AddPage(http.StatusOK, fabcore.DeploymentPipelinesClientListDeploymentPipelineRoleAssignmentsResponse{DeploymentPipelineRoleAssignments: exampleResp}, nil)

		return resp
	}
}

// Returns a fake function that simulates creating a deployment pipeline role assignment, echoing back the request as the response.
func fakeCreateDeploymentPipelineRoleAssignment() func(ctx context.Context, deploymentPipelineID string, deploymentPipelineRoleAssignmentRequest fabcore.AddDeploymentPipelineRoleAssignmentRequest, options *fabcore.DeploymentPipelinesClientAddDeploymentPipelineRoleAssignmentOptions) (resp azfake.Responder[fabcore.DeploymentPipelinesClientAddDeploymentPipelineRoleAssignmentResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, body fabcore.AddDeploymentPipelineRoleAssignmentRequest, _ *fabcore.DeploymentPipelinesClientAddDeploymentPipelineRoleAssignmentOptions) (resp azfake.Responder[fabcore.DeploymentPipelinesClientAddDeploymentPipelineRoleAssignmentResponse], errResp azfake.ErrorResponder) {
		// Return a response that matches the request
		response := fabcore.DeploymentPipelineRoleAssignment{
			ID:   body.Principal.ID,
			Role: body.Role,
			Principal: &fabcore.Principal{
				ID:   body.Principal.ID,
				Type: body.Principal.Type,
			},
		}
		resp = azfake.Responder[fabcore.DeploymentPipelinesClientAddDeploymentPipelineRoleAssignmentResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.DeploymentPipelinesClientAddDeploymentPipelineRoleAssignmentResponse{DeploymentPipelineRoleAssignment: response}, nil)

		return resp, errResp
	}
}

// Returns a fake function that simulates deleting a deployment pipeline role assignment, always returning a successful response.
func fakeDeleteDeploymentPipelineRoleAssignment() func(ctx context.Context, deploymentPipelineID, principalID string, options *fabcore.DeploymentPipelinesClientDeleteDeploymentPipelineRoleAssignmentOptions) (resp azfake.Responder[fabcore.DeploymentPipelinesClientDeleteDeploymentPipelineRoleAssignmentResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _ string, _ *fabcore.DeploymentPipelinesClientDeleteDeploymentPipelineRoleAssignmentOptions) (resp azfake.Responder[fabcore.DeploymentPipelinesClientDeleteDeploymentPipelineRoleAssignmentResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.DeploymentPipelinesClientDeleteDeploymentPipelineRoleAssignmentResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.DeploymentPipelinesClientDeleteDeploymentPipelineRoleAssignmentResponse{}, nil)

		return resp, errResp
	}
}

func NewRandomDeploymentPipelineRoleAssignments() fabcore.DeploymentPipelineRoleAssignments {
	entityID := testhelp.RandomUUID()

	return fabcore.DeploymentPipelineRoleAssignments{
		Value: []fabcore.DeploymentPipelineRoleAssignment{
			{
				ID:   azto.Ptr(entityID),
				Role: azto.Ptr(fabcore.DeploymentPipelineRoleAdmin),
				Principal: &fabcore.Principal{
					ID:   azto.Ptr(entityID),
					Type: azto.Ptr(fabcore.PrincipalTypeGroup),
				},
			},
		},
	}
}
