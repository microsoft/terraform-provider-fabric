// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace_test

import (
	"context"
	"net/http"
	"time"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// func fakeWorkspaceRoleAssignment(exampleResp fabcore.WorkspaceRoleAssignment) func(ctx context.Context, workspaceID, workspaceRoleAssignmentID string, options *fabcore.WorkspacesClientGetWorkspaceRoleAssignmentOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetWorkspaceRoleAssignmentResponse], errResp azfake.ErrorResponder) {
// 	return func(_ context.Context, _, _ string, _ *fabcore.WorkspacesClientGetWorkspaceRoleAssignmentOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetWorkspaceRoleAssignmentResponse], errResp azfake.ErrorResponder) {
// 		resp = azfake.Responder[fabcore.WorkspacesClientGetWorkspaceRoleAssignmentResponse]{}
// 		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientGetWorkspaceRoleAssignmentResponse{WorkspaceRoleAssignment: exampleResp}, nil)

// 		return
// 	}
// }

// func NewRandomWorkspaceRoleAssignment() fabcore.WorkspaceRoleAssignment {
// 	return fabcore.WorkspaceRoleAssignment{
// 		ID: azto.Ptr(testhelp.RandomUUID()),
// 		Principal: &fabcore.Principal{
// 			ID:          azto.Ptr(testhelp.RandomUUID()),
// 			Type:        azto.Ptr(fabcore.PrincipalTypeUser),
// 			DisplayName: azto.Ptr(testhelp.RandomName()),
// 			UserDetails: &fabcore.PrincipalUserDetails{
// 				UserPrincipalName: azto.Ptr(testhelp.RandomName()),
// 			},
// 		},
// 		Role: azto.Ptr(fabcore.WorkspaceRoleAdmin),
// 	}
// }

func fakeWorkspaceRoleAssignments(exampleResp fabcore.WorkspaceRoleAssignments) func(workspaceID string, options *fabcore.WorkspacesClientListWorkspaceRoleAssignmentsOptions) (resp azfake.PagerResponder[fabcore.WorkspacesClientListWorkspaceRoleAssignmentsResponse]) {
	return func(_ string, _ *fabcore.WorkspacesClientListWorkspaceRoleAssignmentsOptions) (resp azfake.PagerResponder[fabcore.WorkspacesClientListWorkspaceRoleAssignmentsResponse]) {
		resp = azfake.PagerResponder[fabcore.WorkspacesClientListWorkspaceRoleAssignmentsResponse]{}
		resp.AddPage(http.StatusOK, fabcore.WorkspacesClientListWorkspaceRoleAssignmentsResponse{WorkspaceRoleAssignments: exampleResp}, nil)

		return
	}
}

func NewRandomWorkspaceRoleAssignments() fabcore.WorkspaceRoleAssignments {
	principal0ID := testhelp.RandomUUID()
	principal1ID := testhelp.RandomUUID()
	principal2ID := testhelp.RandomUUID()
	principal3ID := testhelp.RandomUUID()

	return fabcore.WorkspaceRoleAssignments{
		Value: []fabcore.WorkspaceRoleAssignment{
			{
				ID:   azto.Ptr(principal0ID),
				Role: azto.Ptr(fabcore.WorkspaceRoleAdmin),
				Principal: &fabcore.Principal{
					ID:          azto.Ptr(principal0ID),
					Type:        azto.Ptr(fabcore.PrincipalTypeGroup),
					DisplayName: azto.Ptr(testhelp.RandomName()),
					GroupDetails: &fabcore.PrincipalGroupDetails{
						GroupType: azto.Ptr(fabcore.GroupTypeSecurityGroup),
					},
				},
			},
			{
				ID:   azto.Ptr(principal1ID),
				Role: azto.Ptr(fabcore.WorkspaceRoleMember),
				Principal: &fabcore.Principal{
					ID:          azto.Ptr(principal1ID),
					Type:        azto.Ptr(fabcore.PrincipalTypeUser),
					DisplayName: azto.Ptr(testhelp.RandomName()),
					UserDetails: &fabcore.PrincipalUserDetails{
						UserPrincipalName: azto.Ptr(testhelp.RandomName()),
					},
				},
			},
			{
				ID:   azto.Ptr(principal2ID),
				Role: azto.Ptr(fabcore.WorkspaceRoleMember),
				Principal: &fabcore.Principal{
					ID:          azto.Ptr(principal2ID),
					Type:        azto.Ptr(fabcore.PrincipalTypeServicePrincipal),
					DisplayName: azto.Ptr(testhelp.RandomName()),
					ServicePrincipalDetails: &fabcore.PrincipalServicePrincipalDetails{
						AADAppID: azto.Ptr(testhelp.RandomUUID()),
					},
				},
			},
			{
				ID:   azto.Ptr(principal3ID),
				Role: azto.Ptr(fabcore.WorkspaceRoleViewer),
				Principal: &fabcore.Principal{
					ID:          azto.Ptr(principal3ID),
					Type:        azto.Ptr(fabcore.PrincipalTypeServicePrincipalProfile),
					DisplayName: azto.Ptr(testhelp.RandomName()),
					ServicePrincipalProfileDetails: &fabcore.PrincipalServicePrincipalProfileDetails{
						ParentPrincipal: &fabcore.Principal{
							ID:          azto.Ptr(principal2ID),
							Type:        azto.Ptr(fabcore.PrincipalTypeServicePrincipal),
							DisplayName: azto.Ptr(testhelp.RandomName()),
							ServicePrincipalDetails: &fabcore.PrincipalServicePrincipalDetails{
								AADAppID: azto.Ptr(testhelp.RandomUUID()),
							},
						},
					},
				},
			},
		},
	}
}

func fakeGitConnect() func(ctx context.Context, workspaceID string, gitConnectRequest fabcore.GitConnectRequest, options *fabcore.GitClientConnectOptions) (resp azfake.Responder[fabcore.GitClientConnectResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ fabcore.GitConnectRequest, _ *fabcore.GitClientConnectOptions) (resp azfake.Responder[fabcore.GitClientConnectResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.GitClientConnectResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.GitClientConnectResponse{}, nil)

		return
	}
}

func NewRandomGitConnectRequest() fabcore.GitConnectRequest {
	return fabcore.GitConnectRequest{
		GitProviderDetails: NewRandomGitConnection().GitProviderDetails,
	}
}

func fakeGitInitializeGitConnection(exampleResp fabcore.InitializeGitConnectionResponse) func(ctx context.Context, workspaceID string, options *fabcore.GitClientBeginInitializeConnectionOptions) (resp azfake.PollerResponder[fabcore.GitClientInitializeConnectionResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.GitClientBeginInitializeConnectionOptions) (resp azfake.PollerResponder[fabcore.GitClientInitializeConnectionResponse], errResp azfake.ErrorResponder) {
		resp = azfake.PollerResponder[fabcore.GitClientInitializeConnectionResponse]{}
		resp.SetTerminalResponse(http.StatusOK, fabcore.GitClientInitializeConnectionResponse{InitializeGitConnectionResponse: exampleResp}, nil)

		return
	}
}

func NewRandomGitInitializeGitConnection() fabcore.InitializeGitConnectionResponse {
	return fabcore.InitializeGitConnectionResponse{
		RemoteCommitHash: azto.Ptr("7d03b2918bf6aa62f96d0a4307293f3853201705"),
		RequiredAction:   azto.Ptr(fabcore.RequiredActionUpdateFromGit),
		WorkspaceHead:    azto.Ptr("eaa737b48cda41b37ffefac772ea48f6fed3eac4"),
	}
}

func fakeGitUpdateFromGit() func(ctx context.Context, workspaceID string, updateFromGitRequest fabcore.UpdateFromGitRequest, options *fabcore.GitClientBeginUpdateFromGitOptions) (resp azfake.PollerResponder[fabcore.GitClientUpdateFromGitResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ fabcore.UpdateFromGitRequest, _ *fabcore.GitClientBeginUpdateFromGitOptions) (resp azfake.PollerResponder[fabcore.GitClientUpdateFromGitResponse], errResp azfake.ErrorResponder) {
		resp = azfake.PollerResponder[fabcore.GitClientUpdateFromGitResponse]{}
		resp.SetTerminalResponse(http.StatusOK, fabcore.GitClientUpdateFromGitResponse{}, nil)

		return
	}
}

func fakeGitCommitToGit() func(ctx context.Context, workspaceID string, commitToGitRequest fabcore.CommitToGitRequest, options *fabcore.GitClientBeginCommitToGitOptions) (resp azfake.PollerResponder[fabcore.GitClientCommitToGitResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ fabcore.CommitToGitRequest, _ *fabcore.GitClientBeginCommitToGitOptions) (resp azfake.PollerResponder[fabcore.GitClientCommitToGitResponse], errResp azfake.ErrorResponder) {
		resp = azfake.PollerResponder[fabcore.GitClientCommitToGitResponse]{}
		resp.SetTerminalResponse(http.StatusOK, fabcore.GitClientCommitToGitResponse{}, nil)

		return
	}
}

func fakeGitDisconnect() func(ctx context.Context, workspaceID string, options *fabcore.GitClientDisconnectOptions) (resp azfake.Responder[fabcore.GitClientDisconnectResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.GitClientDisconnectOptions) (resp azfake.Responder[fabcore.GitClientDisconnectResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.GitClientDisconnectResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.GitClientDisconnectResponse{}, nil)

		return
	}
}

func fakeGitGetStatus(exampleResp fabcore.GitStatusResponse) func(ctx context.Context, workspaceID string, options *fabcore.GitClientBeginGetStatusOptions) (resp azfake.PollerResponder[fabcore.GitClientGetStatusResponse], errResp azfake.ErrorResponder) { //nolint:unused
	return func(_ context.Context, _ string, _ *fabcore.GitClientBeginGetStatusOptions) (resp azfake.PollerResponder[fabcore.GitClientGetStatusResponse], errResp azfake.ErrorResponder) {
		resp = azfake.PollerResponder[fabcore.GitClientGetStatusResponse]{}
		resp.SetTerminalResponse(http.StatusOK, fabcore.GitClientGetStatusResponse{GitStatusResponse: exampleResp}, nil)

		return
	}
}

func fakeGitGetConnection(exampleResp fabcore.GitConnection) func(ctx context.Context, workspaceID string, options *fabcore.GitClientGetConnectionOptions) (resp azfake.Responder[fabcore.GitClientGetConnectionResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.GitClientGetConnectionOptions) (resp azfake.Responder[fabcore.GitClientGetConnectionResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.GitClientGetConnectionResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.GitClientGetConnectionResponse{GitConnection: exampleResp}, nil)

		return
	}
}

func NewRandomGitConnection() fabcore.GitConnection {
	return fabcore.GitConnection{
		GitConnectionState: azto.Ptr(fabcore.GitConnectionStateConnectedAndInitialized),
		GitProviderDetails: &fabcore.AzureDevOpsDetails{
			GitProviderType:  azto.Ptr(fabcore.GitProviderTypeAzureDevOps),
			OrganizationName: azto.Ptr("TestOrganization"),
			ProjectName:      azto.Ptr("TestProject"),
			RepositoryName:   azto.Ptr("TestRepo"),
			BranchName:       azto.Ptr("TestBranch"),
			DirectoryName:    azto.Ptr("/TestDirectory"),
		},
		GitSyncDetails: &fabcore.GitSyncDetails{
			Head:         azto.Ptr("eaa737b48cda41b37ffefac772ea48f6fed3eac4"),
			LastSyncTime: azto.Ptr(time.Now()),
		},
	}
}
