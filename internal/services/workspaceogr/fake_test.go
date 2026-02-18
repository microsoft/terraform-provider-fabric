// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceogr_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

func fakeSetOutboundGatewayRules(
	entity *fabcore.WorkspaceOutboundGateways,
) func(ctx context.Context, workspaceID string, workspaceOutboundGateways fabcore.WorkspaceOutboundGateways, options *fabcore.WorkspacesClientSetOutboundGatewayRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientSetOutboundGatewayRulesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, workspaceOutboundGateways fabcore.WorkspaceOutboundGateways, _ *fabcore.WorkspacesClientSetOutboundGatewayRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientSetOutboundGatewayRulesResponse], errResp azfake.ErrorResponder) {
		if workspaceOutboundGateways.DefaultAction != nil {
			entity.DefaultAction = workspaceOutboundGateways.DefaultAction
		}
		if workspaceOutboundGateways.AllowedGateways != nil {
			entity.AllowedGateways = workspaceOutboundGateways.AllowedGateways
		}

		resp = azfake.Responder[fabcore.WorkspacesClientSetOutboundGatewayRulesResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientSetOutboundGatewayRulesResponse{
			ETag: azto.Ptr("fake-etag"),
		}, nil)

		return resp, errResp
	}
}

func fakeGetOutboundGatewayRules(
	entity *fabcore.WorkspaceOutboundGateways,
) func(ctx context.Context, workspaceID string, options *fabcore.WorkspacesClientGetOutboundGatewayRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetOutboundGatewayRulesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.WorkspacesClientGetOutboundGatewayRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetOutboundGatewayRulesResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.WorkspacesClientGetOutboundGatewayRulesResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientGetOutboundGatewayRulesResponse{
			WorkspaceOutboundGateways: *entity,
			ETag:                      azto.Ptr("fake-etag"),
		}, nil)

		return resp, errResp
	}
}

func NewRandomWorkspaceOutboundGateways() fabcore.WorkspaceOutboundGateways {
	return fabcore.WorkspaceOutboundGateways{
		DefaultAction:   azto.Ptr(fabcore.GatewayAccessActionTypeAllow),
		AllowedGateways: []fabcore.GatewayAccessRuleMetadata{},
	}
}
