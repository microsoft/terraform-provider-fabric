// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceocr_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func fakeSetOutboundCloudConnectionRules(
	entity *fabcore.WorkspaceOutboundConnections,
) func(ctx context.Context, workspaceID string, workspaceOutboundConnections fabcore.WorkspaceOutboundConnections, options *fabcore.WorkspacesClientSetOutboundCloudConnectionRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientSetOutboundCloudConnectionRulesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, workspaceOutboundConnections fabcore.WorkspaceOutboundConnections, _ *fabcore.WorkspacesClientSetOutboundCloudConnectionRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientSetOutboundCloudConnectionRulesResponse], errResp azfake.ErrorResponder) {
		if workspaceOutboundConnections.DefaultAction != nil {
			entity.DefaultAction = workspaceOutboundConnections.DefaultAction
		} else {
			entity.DefaultAction = azto.Ptr(fabcore.ConnectionAccessActionTypeDeny)
		}

		if workspaceOutboundConnections.Rules != nil {
			entity.Rules = workspaceOutboundConnections.Rules
		} else {
			entity.Rules = []fabcore.OutboundConnectionRule{}
		}

		resp = azfake.Responder[fabcore.WorkspacesClientSetOutboundCloudConnectionRulesResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientSetOutboundCloudConnectionRulesResponse{
			ETag: new("fake-etag"),
		}, nil)

		return resp, errResp
	}
}

func fakeGetOutboundCloudConnectionRules(
	entity *fabcore.WorkspaceOutboundConnections,
) func(ctx context.Context, workspaceID string, options *fabcore.WorkspacesClientGetOutboundCloudConnectionRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetOutboundCloudConnectionRulesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.WorkspacesClientGetOutboundCloudConnectionRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetOutboundCloudConnectionRulesResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.WorkspacesClientGetOutboundCloudConnectionRulesResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientGetOutboundCloudConnectionRulesResponse{
			WorkspaceOutboundConnections: *entity,
			ETag:                         new("fake-etag"),
		}, nil)

		return resp, errResp
	}
}

func NewRandomWorkspaceOutboundConnections() fabcore.WorkspaceOutboundConnections {
	workspaceID := testhelp.RandomUUID()

	return fabcore.WorkspaceOutboundConnections{
		DefaultAction: new(fabcore.ConnectionAccessActionTypeDeny),
		Rules: []fabcore.OutboundConnectionRule{
			{
				ConnectionType: new("SQL"),
				DefaultAction:  new(fabcore.ConnectionAccessActionTypeDeny),
				AllowedEndpoints: []fabcore.ConnectionRuleEndpointMetadata{
					{
						HostnamePattern: new("*.microsoft.com"),
					},
				},
				AllowedWorkspaces: []fabcore.ConnectionRuleWorkspaceMetadata{},
			},
			{
				ConnectionType:   new("LAKEHOUSE"),
				DefaultAction:    new(fabcore.ConnectionAccessActionTypeDeny),
				AllowedEndpoints: []fabcore.ConnectionRuleEndpointMetadata{},
				AllowedWorkspaces: []fabcore.ConnectionRuleWorkspaceMetadata{
					{
						WorkspaceID: new(workspaceID),
					},
				},
			},
		},
	}
}
