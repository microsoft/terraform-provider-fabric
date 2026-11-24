// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakediagnosticsetting_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

type oneLakeSettingsState struct {
	settings        map[string]fabcore.OneLakeDiagnosticSettings
	omitDiagnostics map[string]bool
}

func newFakeOneLakeSettingsServer(settings map[string]fabcore.OneLakeDiagnosticSettings) (*fabfake.ServerFactory, *oneLakeSettingsState) {
	state := &oneLakeSettingsState{
		settings:        settings,
		omitDiagnostics: make(map[string]bool),
	}

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(state)
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.BeginModifyDiagnostics = fakeModifyOneLakeDiagnostics(state)

	return fakeServer.ServerFactory, state
}

func fakeGetOneLakeSettings(
	state *oneLakeSettingsState,
) func(ctx context.Context, workspaceID string, options *fabcore.OneLakeSettingsClientGetSettingsOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientGetSettingsResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID string, _ *fabcore.OneLakeSettingsClientGetSettingsOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientGetSettingsResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeSettingsClientGetSettingsResponse]{}

		settings, ok := state.settings[workspaceID]
		if !ok {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, "WorkspaceNotFound", "Workspace not found"))
			resp.SetResponse(http.StatusNotFound, fabcore.OneLakeSettingsClientGetSettingsResponse{}, nil)

			return resp, errResp
		}

		if state.omitDiagnostics[workspaceID] {
			resp.SetResponse(http.StatusOK, fabcore.OneLakeSettingsClientGetSettingsResponse{}, nil)

			return resp, errResp
		}

		resp.SetResponse(http.StatusOK, fabcore.OneLakeSettingsClientGetSettingsResponse{
			Diagnostics: &settings,
		}, nil)

		return resp, errResp
	}
}

func fakeModifyOneLakeDiagnostics(
	state *oneLakeSettingsState,
) func(ctx context.Context, workspaceID string, options *fabcore.OneLakeSettingsClientBeginModifyDiagnosticsOptions) (resp azfake.PollerResponder[fabcore.OneLakeSettingsClientModifyDiagnosticsResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID string, options *fabcore.OneLakeSettingsClientBeginModifyDiagnosticsOptions) (resp azfake.PollerResponder[fabcore.OneLakeSettingsClientModifyDiagnosticsResponse], errResp azfake.ErrorResponder) {
		resp = azfake.PollerResponder[fabcore.OneLakeSettingsClientModifyDiagnosticsResponse]{}

		if _, ok := state.settings[workspaceID]; !ok {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, "WorkspaceNotFound", "Workspace not found"))
			resp.SetTerminalResponse(http.StatusNotFound, fabcore.OneLakeSettingsClientModifyDiagnosticsResponse{}, nil)

			return resp, errResp
		}

		if options == nil || options.ModifyOneLakeDiagnosticSettingRequest == nil {
			errResp.SetError(fabfake.SetResponseError(http.StatusBadRequest, "BadRequest", "Diagnostics request is required"))
			resp.SetTerminalResponse(http.StatusBadRequest, fabcore.OneLakeSettingsClientModifyDiagnosticsResponse{}, nil)

			return resp, errResp
		}

		state.settings[workspaceID] = *options.ModifyOneLakeDiagnosticSettingRequest

		resp.SetTerminalResponse(http.StatusOK, fabcore.OneLakeSettingsClientModifyDiagnosticsResponse{}, nil)

		return resp, errResp
	}
}

func newEnabledOneLakeDiagnosticSetting(destinationWorkspaceID, destinationItemID string) fabcore.OneLakeDiagnosticSettings {
	return fabcore.OneLakeDiagnosticSettings{
		Status: new("Enabled"),
		Destination: &fabcore.LakehouseOneLakeDiagnosticSettingsDestination{
			Type: new("Lakehouse"),
			Lakehouse: &fabcore.ItemReferenceByID{
				ReferenceType: to.Ptr(fabcore.ItemReferenceTypeByID),
				ItemID:        new(destinationItemID),
				WorkspaceID:   new(destinationWorkspaceID),
			},
		},
	}
}
