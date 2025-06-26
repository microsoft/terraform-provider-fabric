// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package deploymentpipeline_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

func fakeWorkspaceAssignmentStage() func(ctx context.Context, deploymentPipelineID, stageID string, deploymentPipelineAssignWorkspaceRequest fabcore.DeploymentPipelineAssignWorkspaceRequest, options *fabcore.DeploymentPipelinesClientAssignWorkspaceToStageOptions) (resp azfake.Responder[fabcore.DeploymentPipelinesClientAssignWorkspaceToStageResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _ string, _ fabcore.DeploymentPipelineAssignWorkspaceRequest, _ *fabcore.DeploymentPipelinesClientAssignWorkspaceToStageOptions) (resp azfake.Responder[fabcore.DeploymentPipelinesClientAssignWorkspaceToStageResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.DeploymentPipelinesClientAssignWorkspaceToStageResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.DeploymentPipelinesClientAssignWorkspaceToStageResponse{}, nil)

		return
	}
}

func fakeWorkspaceUnassignmentStage() func(ctx context.Context, deploymentPipelineID, stageID string, options *fabcore.DeploymentPipelinesClientUnassignWorkspaceFromStageOptions) (resp azfake.Responder[fabcore.DeploymentPipelinesClientUnassignWorkspaceFromStageResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _ string, _ *fabcore.DeploymentPipelinesClientUnassignWorkspaceFromStageOptions) (resp azfake.Responder[fabcore.DeploymentPipelinesClientUnassignWorkspaceFromStageResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.DeploymentPipelinesClientUnassignWorkspaceFromStageResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.DeploymentPipelinesClientUnassignWorkspaceFromStageResponse{}, nil)

		return
	}
}

func fakeGetDeploymentPipeline(
	exampleResp fabcore.DeploymentPipelineExtendedInfo,
) func(ctx context.Context, deploymentPipelineID string, options *fabcore.DeploymentPipelinesClientGetDeploymentPipelineOptions) (resp azfake.Responder[fabcore.DeploymentPipelinesClientGetDeploymentPipelineResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.DeploymentPipelinesClientGetDeploymentPipelineOptions) (resp azfake.Responder[fabcore.DeploymentPipelinesClientGetDeploymentPipelineResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.DeploymentPipelinesClientGetDeploymentPipelineResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.DeploymentPipelinesClientGetDeploymentPipelineResponse{DeploymentPipelineExtendedInfo: exampleResp}, nil)

		return
	}
}
