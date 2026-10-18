// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspacefwr_test

import (
	"context"
	"net/http"
	"slices"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

const (
	errCodeSetFirewallRules = "InvalidFirewallRule"
	errMsgSetFirewallRules  = "Firewall rule set operation rejected"
	errCodeGetFirewallRules = "WorkspaceNotFound"
	errMsgGetFirewallRules  = "Firewall rule get operation rejected"
)

func fakeSetFirewallRules(
	entity *fabcore.InboundFirewallConfiguration,
) func(ctx context.Context, workspaceID string, firewallRulesRequest fabcore.InboundFirewallConfiguration, options *fabcore.WorkspacesClientSetFirewallRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientSetFirewallRulesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, firewallRulesRequest fabcore.InboundFirewallConfiguration, _ *fabcore.WorkspacesClientSetFirewallRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientSetFirewallRulesResponse], errResp azfake.ErrorResponder) {
		entity.Rules = firewallRulesRequest.Rules

		resp = azfake.Responder[fabcore.WorkspacesClientSetFirewallRulesResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientSetFirewallRulesResponse{}, nil)

		return resp, errResp
	}
}

func fakeGetFirewallRules(
	entity *fabcore.InboundFirewallConfiguration,
) func(ctx context.Context, workspaceID string, options *fabcore.WorkspacesClientGetFirewallRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetFirewallRulesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.WorkspacesClientGetFirewallRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetFirewallRulesResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.WorkspacesClientGetFirewallRulesResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientGetFirewallRulesResponse{
			InboundFirewallConfiguration: *entity,
		}, nil)

		return resp, errResp
	}
}

func NewRandomInboundFirewallConfiguration() fabcore.InboundFirewallConfiguration {
	return fabcore.InboundFirewallConfiguration{
		Rules: []fabcore.FirewallRule{
			{
				DisplayName: new(testhelp.RandomName()),
				Value:       new("203.0.113.0/24"),
			},
			{
				DisplayName: new(testhelp.RandomName()),
				Value:       new("198.51.100.10-198.51.100.20"),
			},
		},
	}
}

// fakeSetFirewallRulesFailOnCall behaves like fakeSetFirewallRules except that the nth call fails,
// which lets a single test target the Create, Update or Delete error path in isolation.
func fakeSetFirewallRulesFailOnCall(
	entity *fabcore.InboundFirewallConfiguration,
	failOnCall int,
) func(ctx context.Context, workspaceID string, firewallRulesRequest fabcore.InboundFirewallConfiguration, options *fabcore.WorkspacesClientSetFirewallRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientSetFirewallRulesResponse], errResp azfake.ErrorResponder) {
	calls := 0

	return func(_ context.Context, _ string, firewallRulesRequest fabcore.InboundFirewallConfiguration, _ *fabcore.WorkspacesClientSetFirewallRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientSetFirewallRulesResponse], errResp azfake.ErrorResponder) {
		calls++

		resp = azfake.Responder[fabcore.WorkspacesClientSetFirewallRulesResponse]{}

		if calls == failOnCall {
			errResp.SetError(fabfake.SetResponseError(http.StatusBadRequest, errCodeSetFirewallRules, errMsgSetFirewallRules))
			resp.SetResponse(http.StatusBadRequest, fabcore.WorkspacesClientSetFirewallRulesResponse{}, nil)

			return resp, errResp
		}

		entity.Rules = firewallRulesRequest.Rules

		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientSetFirewallRulesResponse{}, nil)

		return resp, errResp
	}
}

func fakeGetFirewallRulesError() func(ctx context.Context, workspaceID string, options *fabcore.WorkspacesClientGetFirewallRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetFirewallRulesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.WorkspacesClientGetFirewallRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetFirewallRulesResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.WorkspacesClientGetFirewallRulesResponse]{}

		errResp.SetError(fabfake.SetResponseError(http.StatusBadRequest, errCodeGetFirewallRules, errMsgGetFirewallRules))
		resp.SetResponse(http.StatusBadRequest, fabcore.WorkspacesClientGetFirewallRulesResponse{}, nil)

		return resp, errResp
	}
}

// fakeGetFirewallRulesReversed returns the stored rules in the opposite order to the one they were
// sent in, so that any ordering sensitivity in the provider surfaces as a diff.
func fakeGetFirewallRulesReversed(
	entity *fabcore.InboundFirewallConfiguration,
) func(ctx context.Context, workspaceID string, options *fabcore.WorkspacesClientGetFirewallRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetFirewallRulesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.WorkspacesClientGetFirewallRulesOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetFirewallRulesResponse], errResp azfake.ErrorResponder) {
		reversed := slices.Clone(entity.Rules)
		slices.Reverse(reversed)

		resp = azfake.Responder[fabcore.WorkspacesClientGetFirewallRulesResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientGetFirewallRulesResponse{
			Rules: reversed,
		}, nil)

		return resp, errResp
	}
}
