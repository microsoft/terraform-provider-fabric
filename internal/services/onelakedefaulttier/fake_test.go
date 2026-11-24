// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedefaulttier_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
)

type oneLakeDefaultTierState struct {
	defaultTier *fabcore.OneLakeAccessTier
}

func fakeGetOneLakeSettings(
	state *oneLakeDefaultTierState,
) func(ctx context.Context, workspaceID string, options *fabcore.OneLakeSettingsClientGetSettingsOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientGetSettingsResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.OneLakeSettingsClientGetSettingsOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientGetSettingsResponse], errResp azfake.ErrorResponder) {
		response := fabcore.GetOneLakeSettingsResponse{}
		if state.defaultTier != nil {
			response.Lifecycle = &fabcore.GetOneLakeSettingsResponseLifecycle{
				DefaultTier: state.defaultTier,
			}
		}

		resp = azfake.Responder[fabcore.OneLakeSettingsClientGetSettingsResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeSettingsClientGetSettingsResponse{
			GetOneLakeSettingsResponse: response,
		}, nil)

		return resp, errResp
	}
}

func fakeModifyDefaultTier(
	state *oneLakeDefaultTierState,
) func(ctx context.Context, workspaceID string, defaultTier fabcore.OneLakeAccessTier, options *fabcore.OneLakeSettingsClientModifyDefaultTierOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientModifyDefaultTierResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, defaultTier fabcore.OneLakeAccessTier, _ *fabcore.OneLakeSettingsClientModifyDefaultTierOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientModifyDefaultTierResponse], errResp azfake.ErrorResponder) {
		state.defaultTier = &defaultTier

		resp = azfake.Responder[fabcore.OneLakeSettingsClientModifyDefaultTierResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeSettingsClientModifyDefaultTierResponse{
			DefaultTier: &defaultTier,
		}, nil)

		return resp, errResp
	}
}

func fakeModifyDefaultTierMissingResponse() func(ctx context.Context, workspaceID string, defaultTier fabcore.OneLakeAccessTier, options *fabcore.OneLakeSettingsClientModifyDefaultTierOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientModifyDefaultTierResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ fabcore.OneLakeAccessTier, _ *fabcore.OneLakeSettingsClientModifyDefaultTierOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientModifyDefaultTierResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeSettingsClientModifyDefaultTierResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeSettingsClientModifyDefaultTierResponse{}, nil)

		return resp, errResp
	}
}

func fakeModifyDefaultTierError() func(ctx context.Context, workspaceID string, defaultTier fabcore.OneLakeAccessTier, options *fabcore.OneLakeSettingsClientModifyDefaultTierOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientModifyDefaultTierResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ fabcore.OneLakeAccessTier, _ *fabcore.OneLakeSettingsClientModifyDefaultTierOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientModifyDefaultTierResponse], errResp azfake.ErrorResponder) {
		errResp.SetError(fabfake.SetResponseError(http.StatusBadRequest, "InvalidInput", "Invalid OneLake access tier"))

		resp = azfake.Responder[fabcore.OneLakeSettingsClientModifyDefaultTierResponse]{}
		resp.SetResponse(http.StatusBadRequest, fabcore.OneLakeSettingsClientModifyDefaultTierResponse{}, nil)

		return resp, errResp
	}
}

func fakeGetOneLakeSettingsWithoutLifecycle() func(ctx context.Context, workspaceID string, options *fabcore.OneLakeSettingsClientGetSettingsOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientGetSettingsResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.OneLakeSettingsClientGetSettingsOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientGetSettingsResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeSettingsClientGetSettingsResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeSettingsClientGetSettingsResponse{}, nil)

		return resp, errResp
	}
}

func fakeGetOneLakeSettingsError() func(ctx context.Context, workspaceID string, options *fabcore.OneLakeSettingsClientGetSettingsOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientGetSettingsResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.OneLakeSettingsClientGetSettingsOptions) (resp azfake.Responder[fabcore.OneLakeSettingsClientGetSettingsResponse], errResp azfake.ErrorResponder) {
		errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, "EntityNotFound", "Workspace not found"))

		resp = azfake.Responder[fabcore.OneLakeSettingsClientGetSettingsResponse]{}
		resp.SetResponse(http.StatusNotFound, fabcore.OneLakeSettingsClientGetSettingsResponse{}, nil)

		return resp, errResp
	}
}
