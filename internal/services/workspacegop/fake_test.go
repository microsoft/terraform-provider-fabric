// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspacegop_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func fakeSetGitOutboundPolicy(
	entity *fabcore.NetworkRules,
) func(ctx context.Context, workspaceID string, workspaceGitOutboundPolicy fabcore.NetworkRules, options *fabcore.WorkspacesClientSetGitOutboundPolicyOptions) (resp azfake.Responder[fabcore.WorkspacesClientSetGitOutboundPolicyResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, workspaceGitOutboundPolicy fabcore.NetworkRules, _ *fabcore.WorkspacesClientSetGitOutboundPolicyOptions) (resp azfake.Responder[fabcore.WorkspacesClientSetGitOutboundPolicyResponse], errResp azfake.ErrorResponder) {
		if workspaceGitOutboundPolicy.DefaultAction != nil {
			entity.DefaultAction = workspaceGitOutboundPolicy.DefaultAction
		} else {
			entity.DefaultAction = azto.Ptr(fabcore.NetworkAccessRuleAllow)
		}

		resp = azfake.Responder[fabcore.WorkspacesClientSetGitOutboundPolicyResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientSetGitOutboundPolicyResponse{
			ETag: new("fake-etag"),
		}, nil)

		return resp, errResp
	}
}

func fakeGetGitOutboundPolicy(
	entity *fabcore.NetworkRules,
) func(ctx context.Context, workspaceID string, options *fabcore.WorkspacesClientGetGitOutboundPolicyOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetGitOutboundPolicyResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.WorkspacesClientGetGitOutboundPolicyOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetGitOutboundPolicyResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.WorkspacesClientGetGitOutboundPolicyResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientGetGitOutboundPolicyResponse{
			NetworkRules: *entity,
			ETag:         new("fake-etag"),
		}, nil)

		return resp, errResp
	}
}

func NewRandomWorkspaceGitOutboundPolicy() fabcore.NetworkRules {
	return fabcore.NetworkRules{
		DefaultAction: new(testhelp.RandomElement(fabcore.PossibleNetworkAccessRuleValues())),
	}
}
