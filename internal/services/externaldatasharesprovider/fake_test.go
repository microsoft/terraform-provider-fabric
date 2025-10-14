// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package externaldatasharesprovider_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func fakeGetExternalDataShareProvider(
	exampleResp fabcore.ExternalDataShare,
) func(ctx context.Context, workspaceID, itemID, externalDataShareID string, options *fabcore.ExternalDataSharesProviderClientGetExternalDataShareOptions) (resp azfake.Responder[fabcore.ExternalDataSharesProviderClientGetExternalDataShareResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _, _ string, _ *fabcore.ExternalDataSharesProviderClientGetExternalDataShareOptions) (resp azfake.Responder[fabcore.ExternalDataSharesProviderClientGetExternalDataShareResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.ExternalDataSharesProviderClientGetExternalDataShareResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.ExternalDataSharesProviderClientGetExternalDataShareResponse{ExternalDataShare: exampleResp}, nil)

		return resp, errResp
	}
}

func fakeListExternalDataSharesProvider(
	exampleResp []fabcore.ExternalDataShare,
) func(workspaceID, itemID string, options *fabcore.ExternalDataSharesProviderClientListExternalDataSharesInItemOptions) (resp azfake.PagerResponder[fabcore.ExternalDataSharesProviderClientListExternalDataSharesInItemResponse]) {
	return func(_, _ string, _ *fabcore.ExternalDataSharesProviderClientListExternalDataSharesInItemOptions) (resp azfake.PagerResponder[fabcore.ExternalDataSharesProviderClientListExternalDataSharesInItemResponse]) {
		resp = azfake.PagerResponder[fabcore.ExternalDataSharesProviderClientListExternalDataSharesInItemResponse]{}
		resp.AddPage(http.StatusOK, fabcore.ExternalDataSharesProviderClientListExternalDataSharesInItemResponse{
			ExternalDataShares: fabcore.ExternalDataShares{
				Value: exampleResp,
			},
		}, nil)

		return resp
	}
}

func fakeCreateExternalDataShareProvider(
	exampleResp fabcore.ExternalDataShare,
) func(ctx context.Context, workspaceID, itemID string, createExternalDataShareRequest fabcore.CreateExternalDataShareRequest, options *fabcore.ExternalDataSharesProviderClientCreateExternalDataShareOptions) (resp azfake.Responder[fabcore.ExternalDataSharesProviderClientCreateExternalDataShareResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _ string, _ fabcore.CreateExternalDataShareRequest, _ *fabcore.ExternalDataSharesProviderClientCreateExternalDataShareOptions) (resp azfake.Responder[fabcore.ExternalDataSharesProviderClientCreateExternalDataShareResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.ExternalDataSharesProviderClientCreateExternalDataShareResponse]{}
		resp.SetResponse(http.StatusCreated, fabcore.ExternalDataSharesProviderClientCreateExternalDataShareResponse{ExternalDataShare: exampleResp}, nil)

		return resp, errResp
	}
}

func fakeDeleteExternalDataShareProvider() func(ctx context.Context, workspaceID, itemID, externalDataShareID string, options *fabcore.ExternalDataSharesProviderClientDeleteExternalDataShareOptions) (resp azfake.Responder[fabcore.ExternalDataSharesProviderClientDeleteExternalDataShareResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _, _ string, _ *fabcore.ExternalDataSharesProviderClientDeleteExternalDataShareOptions) (resp azfake.Responder[fabcore.ExternalDataSharesProviderClientDeleteExternalDataShareResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.ExternalDataSharesProviderClientDeleteExternalDataShareResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.ExternalDataSharesProviderClientDeleteExternalDataShareResponse{}, nil)

		return resp, errResp
	}
}

func NewRandomExternalDataShare(workspaceID string) fabcore.ExternalDataShare {
	return fabcore.ExternalDataShare{
		ID:          to.Ptr(testhelp.RandomUUID()),
		Paths:       []string{testhelp.RandomName()},
		WorkspaceID: to.Ptr(workspaceID),
		Recipient: &fabcore.ExternalDataShareRecipient{
			UserPrincipalName: to.Ptr(testhelp.RandomName()),
		},
		CreatorPrincipal: &fabcore.Principal{
			ID:          to.Ptr(testhelp.RandomUUID()),
			DisplayName: to.Ptr(testhelp.RandomName()),
			Type:        to.Ptr(fabcore.PrincipalTypeUser),
			UserDetails: &fabcore.PrincipalUserDetails{
				UserPrincipalName: to.Ptr(testhelp.RandomName()),
			},
		},
		Status:             to.Ptr(fabcore.ExternalDataShareStatusPending),
		ExpirationTimeUTC:  to.Ptr(testhelp.RandomTimeDefault()),
		ItemID:             to.Ptr(testhelp.RandomUUID()),
		InvitationURL:      to.Ptr(testhelp.RandomName()),
		AcceptedByTenantID: to.Ptr(testhelp.RandomUUID()),
	}
}

func NewRandomExternalDataShares(workspaceID string) fabcore.ExternalDataShares {
	return fabcore.ExternalDataShares{
		Value: []fabcore.ExternalDataShare{
			NewRandomExternalDataShare(workspaceID),
			NewRandomExternalDataShare(workspaceID),
			NewRandomExternalDataShare(workspaceID),
		},
	}
}
