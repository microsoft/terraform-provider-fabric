// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceiar_test

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

const (
	errCodeSetInboundAzureResourceRules = "InvalidInboundAzureResourceRule"
	errMsgSetInboundAzureResourceRules  = "Inbound Azure resource rule set operation rejected"
	errCodeGetInboundAzureResourceRules = "WorkspaceNotFound"
	errMsgGetInboundAzureResourceRules  = "Inbound Azure resource rule get operation rejected"
)

func fakeSetInboundAzureResourceRules(
	entity *fabcore.WorkspaceInboundAzureResourceRules,
) func(ctx context.Context, workspaceID string, workspaceInboundAzureResourceRules fabcore.WorkspaceInboundAzureResourceRules, options *fabcore.WorkspacesClientSetInboundAzureResourceRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientSetInboundAzureResourceRulesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, workspaceInboundAzureResourceRules fabcore.WorkspaceInboundAzureResourceRules, _ *fabcore.WorkspacesClientSetInboundAzureResourceRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientSetInboundAzureResourceRulesResponse], errResp azfake.ErrorResponder) {
		entity.Rules = workspaceInboundAzureResourceRules.Rules

		resp = azfake.Responder[fabcore.WorkspacesClientSetInboundAzureResourceRulesResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientSetInboundAzureResourceRulesResponse{}, nil)

		return resp, errResp
	}
}

func fakeGetInboundAzureResourceRules(
	entity *fabcore.WorkspaceInboundAzureResourceRules,
) func(ctx context.Context, workspaceID string, options *fabcore.WorkspacesClientGetInboundAzureResourceRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetInboundAzureResourceRulesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.WorkspacesClientGetInboundAzureResourceRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetInboundAzureResourceRulesResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.WorkspacesClientGetInboundAzureResourceRulesResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientGetInboundAzureResourceRulesResponse{
			WorkspaceInboundAzureResourceRules: *entity,
		}, nil)

		return resp, errResp
	}
}

// fakeSetInboundAzureResourceRulesFailOnCall behaves like fakeSetInboundAzureResourceRules except that the nth call fails,
// which lets a single test target the Create, Update or Delete error path in isolation.
func fakeSetInboundAzureResourceRulesFailOnCall(
	entity *fabcore.WorkspaceInboundAzureResourceRules,
	failOnCall int,
) func(ctx context.Context, workspaceID string, workspaceInboundAzureResourceRules fabcore.WorkspaceInboundAzureResourceRules, options *fabcore.WorkspacesClientSetInboundAzureResourceRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientSetInboundAzureResourceRulesResponse], errResp azfake.ErrorResponder) {
	calls := 0

	return func(_ context.Context, _ string, workspaceInboundAzureResourceRules fabcore.WorkspaceInboundAzureResourceRules, _ *fabcore.WorkspacesClientSetInboundAzureResourceRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientSetInboundAzureResourceRulesResponse], errResp azfake.ErrorResponder) {
		calls++

		resp = azfake.Responder[fabcore.WorkspacesClientSetInboundAzureResourceRulesResponse]{}

		if calls == failOnCall {
			errResp.SetError(fabfake.SetResponseError(http.StatusBadRequest, errCodeSetInboundAzureResourceRules, errMsgSetInboundAzureResourceRules))
			resp.SetResponse(http.StatusBadRequest, fabcore.WorkspacesClientSetInboundAzureResourceRulesResponse{}, nil)

			return resp, errResp
		}

		entity.Rules = workspaceInboundAzureResourceRules.Rules

		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientSetInboundAzureResourceRulesResponse{}, nil)

		return resp, errResp
	}
}

func fakeGetInboundAzureResourceRulesError() func(ctx context.Context, workspaceID string, options *fabcore.WorkspacesClientGetInboundAzureResourceRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetInboundAzureResourceRulesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.WorkspacesClientGetInboundAzureResourceRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetInboundAzureResourceRulesResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.WorkspacesClientGetInboundAzureResourceRulesResponse]{}

		errResp.SetError(fabfake.SetResponseError(http.StatusBadRequest, errCodeGetInboundAzureResourceRules, errMsgGetInboundAzureResourceRules))
		resp.SetResponse(http.StatusBadRequest, fabcore.WorkspacesClientGetInboundAzureResourceRulesResponse{}, nil)

		return resp, errResp
	}
}

// fakeGetInboundAzureResourceRulesReversed returns the stored rules in the opposite order to the one they were
// sent in, so that any ordering sensitivity in the provider surfaces as a diff.
func fakeGetInboundAzureResourceRulesReversed(
	entity *fabcore.WorkspaceInboundAzureResourceRules,
) func(ctx context.Context, workspaceID string, options *fabcore.WorkspacesClientGetInboundAzureResourceRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetInboundAzureResourceRulesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.WorkspacesClientGetInboundAzureResourceRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetInboundAzureResourceRulesResponse], errResp azfake.ErrorResponder) {
		reversed := slices.Clone(entity.Rules)
		slices.Reverse(reversed)

		resp = azfake.Responder[fabcore.WorkspacesClientGetInboundAzureResourceRulesResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientGetInboundAzureResourceRulesResponse{
			WorkspaceInboundAzureResourceRules: fabcore.WorkspaceInboundAzureResourceRules{Rules: reversed},
		}, nil)

		return resp, errResp
	}
}

func NewRandomWorkspaceInboundAzureResourceRules() fabcore.WorkspaceInboundAzureResourceRules {
	return fabcore.WorkspaceInboundAzureResourceRules{
		Rules: []fabcore.WorkspaceInboundAzureResourceRule{
			{
				DisplayName: new(testhelp.RandomName()),
				ResourceID:  new(NewRandomAzureResourceID("Microsoft.DataFactory/factories")),
			},
			{
				DisplayName: new(testhelp.RandomName()),
				ResourceID:  new(NewRandomAzureResourceID("Microsoft.Sql/servers")),
			},
		},
	}
}

func NewRandomAzureResourceID(resourceType string) string {
	return fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/%s/%s",
		testhelp.RandomUUID(),
		testhelp.RandomName(),
		resourceType,
		testhelp.RandomName(),
	)
}
