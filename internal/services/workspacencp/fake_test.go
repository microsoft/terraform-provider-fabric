// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspacencp_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

func fakeSetNetworkCommunicationPolicy(
	entity *fabcore.WorkspaceNetworkingCommunicationPolicy,
) func(ctx context.Context, workspaceID string, setWorkspaceNetworkingCommunicationPolicy fabcore.WorkspaceNetworkingCommunicationPolicy, options *fabcore.WorkspacesClientSetNetworkCommunicationPolicyOptions) (resp azfake.Responder[fabcore.WorkspacesClientSetNetworkCommunicationPolicyResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, setWorkspaceNetworkingCommunicationPolicy fabcore.WorkspaceNetworkingCommunicationPolicy, _ *fabcore.WorkspacesClientSetNetworkCommunicationPolicyOptions) (resp azfake.Responder[fabcore.WorkspacesClientSetNetworkCommunicationPolicyResponse], errResp azfake.ErrorResponder) {
		if setWorkspaceNetworkingCommunicationPolicy.Inbound != nil && setWorkspaceNetworkingCommunicationPolicy.Inbound.PublicAccessRules != nil {
			entity.Inbound.PublicAccessRules.DefaultAction = setWorkspaceNetworkingCommunicationPolicy.Inbound.PublicAccessRules.DefaultAction
		} else {
			entity.Inbound.PublicAccessRules.DefaultAction = azto.Ptr(fabcore.NetworkAccessRuleAllow)
		}
		if setWorkspaceNetworkingCommunicationPolicy.Outbound != nil && setWorkspaceNetworkingCommunicationPolicy.Outbound.PublicAccessRules != nil {
			entity.Outbound.PublicAccessRules.DefaultAction = setWorkspaceNetworkingCommunicationPolicy.Outbound.PublicAccessRules.DefaultAction
		} else {
			entity.Outbound.PublicAccessRules.DefaultAction = azto.Ptr(fabcore.NetworkAccessRuleAllow)
		}

		resp = azfake.Responder[fabcore.WorkspacesClientSetNetworkCommunicationPolicyResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientSetNetworkCommunicationPolicyResponse{
			ETag: azto.Ptr("fake-etag"),
		}, nil)

		return resp, errResp
	}
}

func fakeGetNetworkCommunicationPolicy(
	entity *fabcore.WorkspaceNetworkingCommunicationPolicy,
) func(ctx context.Context, workspaceID string, options *fabcore.WorkspacesClientGetNetworkCommunicationPolicyOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetNetworkCommunicationPolicyResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.WorkspacesClientGetNetworkCommunicationPolicyOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetNetworkCommunicationPolicyResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.WorkspacesClientGetNetworkCommunicationPolicyResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientGetNetworkCommunicationPolicyResponse{
			WorkspaceNetworkingCommunicationPolicy: *entity,
			ETag:                                   azto.Ptr("fake-etag"),
		}, nil)

		return resp, errResp
	}
}

func NewRandomWorkspaceNetworkingCommunicationPolicy() fabcore.WorkspaceNetworkingCommunicationPolicy {
	return fabcore.WorkspaceNetworkingCommunicationPolicy{
		Inbound: &fabcore.InboundRules{
			PublicAccessRules: &fabcore.NetworkRules{
				DefaultAction: azto.Ptr(fabcore.NetworkAccessRuleAllow),
			},
		},
		Outbound: &fabcore.OutboundRules{
			PublicAccessRules: &fabcore.NetworkRules{
				DefaultAction: azto.Ptr(fabcore.NetworkAccessRuleAllow),
			},
		},
	}
}
