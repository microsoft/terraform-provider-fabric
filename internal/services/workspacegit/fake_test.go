// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspacegit_test

import (
	"context"
	"net/http"
	"time"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func fakeGitConnect() func(ctx context.Context, workspaceID string, gitConnectRequest fabcore.GitConnectRequest, options *fabcore.GitClientConnectOptions) (resp azfake.Responder[fabcore.GitClientConnectResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ fabcore.GitConnectRequest, _ *fabcore.GitClientConnectOptions) (resp azfake.Responder[fabcore.GitClientConnectResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.GitClientConnectResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.GitClientConnectResponse{}, nil)

		return resp, errResp
	}
}

func NewRandomGitConnectRequest(t fabcore.GitProviderType) fabcore.GitConnectRequest {
	var gitCredentials fabcore.GitCredentialsClassification

	switch t {
	case fabcore.GitProviderTypeAzureDevOps:
		gitCredentials = NewRandomGitCredentials(fabcore.GitCredentialsSourceAutomatic)
	case fabcore.GitProviderTypeGitHub:
		gitCredentials = NewRandomGitCredentials(fabcore.GitCredentialsSourceConfiguredConnection)
	}

	return fabcore.GitConnectRequest{
		GitProviderDetails: NewRandomGitConnection(t).GitProviderDetails,
		MyGitCredentials:   gitCredentials,
	}
}

func fakeGitInitializeGitConnection(
	exampleResp fabcore.InitializeGitConnectionResponse,
) func(ctx context.Context, workspaceID string, options *fabcore.GitClientBeginInitializeConnectionOptions) (resp azfake.PollerResponder[fabcore.GitClientInitializeConnectionResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.GitClientBeginInitializeConnectionOptions) (resp azfake.PollerResponder[fabcore.GitClientInitializeConnectionResponse], errResp azfake.ErrorResponder) {
		resp = azfake.PollerResponder[fabcore.GitClientInitializeConnectionResponse]{}
		resp.SetTerminalResponse(http.StatusOK, fabcore.GitClientInitializeConnectionResponse{InitializeGitConnectionResponse: exampleResp}, nil)

		return resp, errResp
	}
}

func NewRandomGitInitializeGitConnection() fabcore.InitializeGitConnectionResponse {
	return fabcore.InitializeGitConnectionResponse{
		RemoteCommitHash: azto.Ptr(testhelp.RandomSHA1()),
		RequiredAction:   azto.Ptr(fabcore.RequiredActionUpdateFromGit),
		WorkspaceHead:    azto.Ptr(testhelp.RandomSHA1()),
	}
}

func fakeGitUpdateFromGit() func(ctx context.Context, workspaceID string, updateFromGitRequest fabcore.UpdateFromGitRequest, options *fabcore.GitClientBeginUpdateFromGitOptions) (resp azfake.PollerResponder[fabcore.GitClientUpdateFromGitResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ fabcore.UpdateFromGitRequest, _ *fabcore.GitClientBeginUpdateFromGitOptions) (resp azfake.PollerResponder[fabcore.GitClientUpdateFromGitResponse], errResp azfake.ErrorResponder) {
		resp = azfake.PollerResponder[fabcore.GitClientUpdateFromGitResponse]{}
		resp.SetTerminalResponse(http.StatusOK, fabcore.GitClientUpdateFromGitResponse{}, nil)

		return resp, errResp
	}
}

func fakeGitCommitToGit() func(ctx context.Context, workspaceID string, commitToGitRequest fabcore.CommitToGitRequest, options *fabcore.GitClientBeginCommitToGitOptions) (resp azfake.PollerResponder[fabcore.GitClientCommitToGitResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ fabcore.CommitToGitRequest, _ *fabcore.GitClientBeginCommitToGitOptions) (resp azfake.PollerResponder[fabcore.GitClientCommitToGitResponse], errResp azfake.ErrorResponder) {
		resp = azfake.PollerResponder[fabcore.GitClientCommitToGitResponse]{}
		resp.SetTerminalResponse(http.StatusOK, fabcore.GitClientCommitToGitResponse{}, nil)

		return resp, errResp
	}
}

func fakeGitDisconnect() func(ctx context.Context, workspaceID string, options *fabcore.GitClientDisconnectOptions) (resp azfake.Responder[fabcore.GitClientDisconnectResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.GitClientDisconnectOptions) (resp azfake.Responder[fabcore.GitClientDisconnectResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.GitClientDisconnectResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.GitClientDisconnectResponse{}, nil)

		return resp, errResp
	}
}

func fakeGitGetStatus( //nolint:unused
	exampleResp fabcore.GitStatusResponse,
) func(ctx context.Context, workspaceID string, options *fabcore.GitClientBeginGetStatusOptions) (resp azfake.PollerResponder[fabcore.GitClientGetStatusResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.GitClientBeginGetStatusOptions) (resp azfake.PollerResponder[fabcore.GitClientGetStatusResponse], errResp azfake.ErrorResponder) {
		resp = azfake.PollerResponder[fabcore.GitClientGetStatusResponse]{}
		resp.SetTerminalResponse(http.StatusOK, fabcore.GitClientGetStatusResponse{GitStatusResponse: exampleResp}, nil)

		return resp, errResp
	}
}

func fakeGitGetConnection(
	exampleResp fabcore.GitConnection,
) func(ctx context.Context, workspaceID string, options *fabcore.GitClientGetConnectionOptions) (resp azfake.Responder[fabcore.GitClientGetConnectionResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.GitClientGetConnectionOptions) (resp azfake.Responder[fabcore.GitClientGetConnectionResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.GitClientGetConnectionResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.GitClientGetConnectionResponse{GitConnection: exampleResp}, nil)

		return resp, errResp
	}
}

func fakeGitGetMyGitCredentials(
	exampleResp fabcore.GitClientGetMyGitCredentialsResponse,
) func(ctx context.Context, workspaceID string, options *fabcore.GitClientGetMyGitCredentialsOptions) (resp azfake.Responder[fabcore.GitClientGetMyGitCredentialsResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.GitClientGetMyGitCredentialsOptions) (resp azfake.Responder[fabcore.GitClientGetMyGitCredentialsResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.GitClientGetMyGitCredentialsResponse]{}
		resp.SetResponse(
			http.StatusOK,
			fabcore.GitClientGetMyGitCredentialsResponse{GitCredentialsConfigurationResponseClassification: exampleResp.GitCredentialsConfigurationResponseClassification},
			nil,
		)

		return resp, errResp
	}
}

func NewRandomGitConnection(t fabcore.GitProviderType) fabcore.GitConnection {
	var gitProviderDetails fabcore.GitProviderDetailsClassification

	switch t {
	case fabcore.GitProviderTypeAzureDevOps:
		gitProviderDetails = &fabcore.AzureDevOpsDetails{
			GitProviderType:  azto.Ptr(fabcore.GitProviderTypeAzureDevOps),
			OrganizationName: azto.Ptr("TestOrganization"),
			ProjectName:      azto.Ptr("TestProject"),
			RepositoryName:   azto.Ptr("TestRepo"),
			BranchName:       azto.Ptr("TestBranch"),
			DirectoryName:    azto.Ptr("/TestDirectory"),
		}

	case fabcore.GitProviderTypeGitHub:
		gitProviderDetails = &fabcore.GitHubDetails{
			GitProviderType: azto.Ptr(fabcore.GitProviderTypeGitHub),
			OwnerName:       azto.Ptr("TestOwner"),
			RepositoryName:  azto.Ptr("TestRepo"),
			BranchName:      azto.Ptr("TestBranch"),
			DirectoryName:   azto.Ptr("/TestDirectory"),
		}
	}

	return fabcore.GitConnection{
		GitConnectionState: azto.Ptr(fabcore.GitConnectionStateConnectedAndInitialized),
		GitProviderDetails: gitProviderDetails,
		GitSyncDetails: &fabcore.GitSyncDetails{
			Head:         azto.Ptr(testhelp.RandomSHA1()),
			LastSyncTime: azto.Ptr(time.Now()),
		},
	}
}

func NewRandomGitCredentials(s fabcore.GitCredentialsSource) fabcore.GitCredentialsClassification {
	var r fabcore.GitCredentialsClassification

	switch s {
	case fabcore.GitCredentialsSourceAutomatic:
		r = &fabcore.AutomaticGitCredentials{
			Source: azto.Ptr(fabcore.GitCredentialsSourceAutomatic),
		}
	case fabcore.GitCredentialsSourceConfiguredConnection:
		r = &fabcore.ConfiguredConnectionGitCredentials{
			Source:       azto.Ptr(fabcore.GitCredentialsSourceConfiguredConnection),
			ConnectionID: azto.Ptr(testhelp.RandomUUID()),
		}
	case fabcore.GitCredentialsSourceNone:
		r = &fabcore.GitCredentials{
			Source: azto.Ptr(fabcore.GitCredentialsSourceNone),
		}
	}

	return r
}

func NewRandomGitCredentialsResponse(s fabcore.GitCredentialsSource) fabcore.GitClientGetMyGitCredentialsResponse {
	var r fabcore.GitCredentialsConfigurationResponseClassification

	switch s {
	case fabcore.GitCredentialsSourceAutomatic:
		r = &fabcore.AutomaticGitCredentialsResponse{
			Source: azto.Ptr(fabcore.GitCredentialsSourceAutomatic),
		}
	case fabcore.GitCredentialsSourceConfiguredConnection:
		r = &fabcore.ConfiguredConnectionGitCredentialsResponse{
			Source:       azto.Ptr(fabcore.GitCredentialsSourceConfiguredConnection),
			ConnectionID: azto.Ptr(testhelp.RandomUUID()),
		}
	case fabcore.GitCredentialsSourceNone:
		r = &fabcore.NoneGitCredentialsResponse{
			Source: azto.Ptr(fabcore.GitCredentialsSourceNone),
		}
	}

	return fabcore.GitClientGetMyGitCredentialsResponse{
		GitCredentialsConfigurationResponseClassification: r,
	}
}
