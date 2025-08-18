// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package externaldatasharesprovider_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// Returns a fake pager function that simulates listing deployment pipeline role assignments with a provided example response.
func fakeListExternalDataSharesProvider(
	exampleResp fabadmin.ExternalDataShares,
) func(options *fabadmin.ExternalDataSharesProviderClientListExternalDataSharesOptions) (resp azfake.PagerResponder[fabadmin.ExternalDataSharesProviderClientListExternalDataSharesResponse]) {
	return func(_ *fabadmin.ExternalDataSharesProviderClientListExternalDataSharesOptions) (resp azfake.PagerResponder[fabadmin.ExternalDataSharesProviderClientListExternalDataSharesResponse]) {
		resp = azfake.PagerResponder[fabadmin.ExternalDataSharesProviderClientListExternalDataSharesResponse]{}
		resp.AddPage(http.StatusOK, fabadmin.ExternalDataSharesProviderClientListExternalDataSharesResponse{ExternalDataShares: exampleResp}, nil)

		return
	}
}

func fakeRevokeExternalDataSharesProvider() func(
	ctx context.Context,
	workspaceID, itemID, externalDataShareID string,
	options *fabadmin.ExternalDataSharesProviderClientRevokeExternalDataShareOptions,
) (
	resp azfake.Responder[fabadmin.ExternalDataSharesProviderClientRevokeExternalDataShareResponse],
	errResp azfake.ErrorResponder,
) {
	return func(_ context.Context, _, _, _ string,
		_ *fabadmin.ExternalDataSharesProviderClientRevokeExternalDataShareOptions,
	) (
		resp azfake.Responder[fabadmin.ExternalDataSharesProviderClientRevokeExternalDataShareResponse],
		errResp azfake.ErrorResponder,
	) {
		resp = azfake.Responder[fabadmin.ExternalDataSharesProviderClientRevokeExternalDataShareResponse]{}
		resp.SetResponse(http.StatusOK, fabadmin.ExternalDataSharesProviderClientRevokeExternalDataShareResponse{}, nil)

		return
	}
}

func NewRandomExternalDataShare(workspaceID string) fabadmin.ExternalDataShare {
	return fabadmin.ExternalDataShare{
		ID:          to.Ptr(testhelp.RandomUUID()),
		Paths:       []string{testhelp.RandomName()},
		WorkspaceID: to.Ptr(workspaceID),
		Recipient: &fabadmin.ExternalDataShareRecipient{
			UserPrincipalName: to.Ptr(testhelp.RandomName()),
		},
		CreatorPrincipal: &fabadmin.Principal{
			ID:          to.Ptr(testhelp.RandomUUID()),
			DisplayName: to.Ptr(testhelp.RandomName()),
			Type:        to.Ptr(fabadmin.PrincipalTypeUser),
			UserDetails: &fabadmin.PrincipalUserDetails{
				UserPrincipalName: to.Ptr(testhelp.RandomName()),
			},
		},
		Status:             to.Ptr(fabadmin.ExternalDataShareStatusPending),
		ExpirationTimeUTC:  to.Ptr(testhelp.RandomTimeDefault()),
		ItemID:             to.Ptr(testhelp.RandomUUID()),
		InvitationURL:      to.Ptr(testhelp.RandomName()),
		AcceptedByTenantID: to.Ptr(testhelp.RandomUUID()),
	}
}

func NewRandomExternalDataShares(workspaceID string) fabadmin.ExternalDataShares {
	return fabadmin.ExternalDataShares{
		Value: []fabadmin.ExternalDataShare{
			NewRandomExternalDataShare(workspaceID),
			NewRandomExternalDataShare(workspaceID),
			NewRandomExternalDataShare(workspaceID),
		},
	}
}
