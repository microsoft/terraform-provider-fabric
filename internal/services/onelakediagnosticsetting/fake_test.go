// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakediagnosticsetting_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

type oneLakeSettingsState struct {
	settings                     map[string]fabcore.OneLakeDiagnosticSettings
	omitDiagnostics              map[string]bool
	noPreviousSettingsErrorCount int
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

		request := *options.ModifyOneLakeDiagnosticSettingRequest
		if isClearedSettings(request) && isClearedSettings(state.settings[workspaceID]) {
			state.noPreviousSettingsErrorCount++

			errResp.SetError(fakeNoPreviousSettingsError())
			resp.SetTerminalResponse(http.StatusBadRequest, fabcore.OneLakeSettingsClientModifyDiagnosticsResponse{}, nil)

			return resp, errResp
		}

		state.settings[workspaceID] = request

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

func isClearedSettings(settings fabcore.OneLakeDiagnosticSettings) bool {
	if settings.Status == nil || *settings.Status != "Disabled" {
		return false
	}

	if settings.Destination == nil {
		return true
	}

	destination, ok := settings.Destination.(*fabcore.LakehouseOneLakeDiagnosticSettingsDestination)

	return ok && destination.Lakehouse == nil
}

func fakeNoPreviousSettingsError() error {
	const errorResponse = `{"errorCode":"BadRequest","message":"The request could not be processed due to missing or invalid information","moreDetails":[{"errorCode":"NoPreviousSettingsExist","message":"No previous OneLake diagnostic settings exist"}]}`

	return fabcore.NewResponseError(&http.Response{
		Status:     "400 Bad Request",
		StatusCode: http.StatusBadRequest,
		Header: http.Header{
			"x-ms-public-api-error-code": []string{"BadRequest"},
		},
		Body: io.NopCloser(strings.NewReader(errorResponse)),
	})
}

type oneLakeSettingsResponseTransport struct {
	transport policy.Transporter
}

func (t *oneLakeSettingsResponseTransport) Do(req *http.Request) (*http.Response, error) {
	response, err := t.transport.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending fake Fabric request: %w", err)
	}

	if response == nil {
		return nil, errors.New("fake Fabric transport returned a nil response")
	}

	if req.Method != http.MethodGet || !strings.HasSuffix(req.URL.Path, "/onelake/settings") || response.StatusCode != http.StatusOK {
		return response, nil
	}

	if response.Body == nil {
		return nil, errors.New("fake Fabric transport returned an empty settings response")
	}

	body, err := io.ReadAll(response.Body)
	closeErr := response.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("reading fake OneLake settings response: %w", err)
	}

	if closeErr != nil {
		return nil, fmt.Errorf("closing fake OneLake settings response: %w", closeErr)
	}

	var responseBody map[string]json.RawMessage

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling fake OneLake settings response: %w", err)
	}

	responseBody["lifecycle"] = json.RawMessage(`{"accessTimeTracking":true}`)

	body, err = json.Marshal(responseBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling fake OneLake settings response: %w", err)
	}

	response.Body = io.NopCloser(bytes.NewReader(body))
	response.ContentLength = int64(len(body))

	if response.Header == nil {
		response.Header = make(http.Header)
	}

	response.Header.Set("Content-Length", strconv.Itoa(len(body)))

	return response, nil
}
