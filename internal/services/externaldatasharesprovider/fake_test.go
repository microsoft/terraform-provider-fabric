// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package externaldatasharesprovider_test

import (
	"context"
	"net/http"
	"time"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var fakeExternalDataShareStore = map[string]fabcore.ExternalDataShare{}

func fakeGetExternalDataShareProvider() func(ctx context.Context, workspaceID, itemID, id string, options *fabcore.ExternalDataSharesProviderClientGetExternalDataShareOptions) (resp azfake.Responder[fabcore.ExternalDataSharesProviderClientGetExternalDataShareResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _, id string, _ *fabcore.ExternalDataSharesProviderClientGetExternalDataShareOptions) (resp azfake.Responder[fabcore.ExternalDataSharesProviderClientGetExternalDataShareResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.ExternalDataSharesProviderClientGetExternalDataShareResponse]{}
		errItemNotFound := fabcore.ErrItem.ItemNotFound.Error()

		if externalDataShare, ok := fakeExternalDataShareStore[id]; ok {
			resp.SetResponse(http.StatusOK, fabcore.ExternalDataSharesProviderClientGetExternalDataShareResponse{ExternalDataShare: externalDataShare}, nil)
		} else {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, errItemNotFound, "Item not found"))
			resp.SetResponse(http.StatusNotFound, fabcore.ExternalDataSharesProviderClientGetExternalDataShareResponse{}, nil)
		}

		return resp, errResp
	}
}

func fakeListExternalDataSharesProvider() func(workspaceID, itemID string, options *fabcore.ExternalDataSharesProviderClientListExternalDataSharesInItemOptions) (resp azfake.PagerResponder[fabcore.ExternalDataSharesProviderClientListExternalDataSharesInItemResponse]) {
	return func(_, _ string, _ *fabcore.ExternalDataSharesProviderClientListExternalDataSharesInItemOptions) (resp azfake.PagerResponder[fabcore.ExternalDataSharesProviderClientListExternalDataSharesInItemResponse]) {
		resp = azfake.PagerResponder[fabcore.ExternalDataSharesProviderClientListExternalDataSharesInItemResponse]{}
		resp.AddPage(http.StatusOK, fabcore.ExternalDataSharesProviderClientListExternalDataSharesInItemResponse{
			ExternalDataShares: fabcore.ExternalDataShares{
				Value: GetAllStoredExternalDataShares(),
			},
		}, nil)

		return resp
	}
}

func fakeCreateExternalDataShareProvider() func(ctx context.Context, workspaceID, itemID string, createExternalDataShareRequest fabcore.CreateExternalDataShareRequest, options *fabcore.ExternalDataSharesProviderClientCreateExternalDataShareOptions) (resp azfake.Responder[fabcore.ExternalDataSharesProviderClientCreateExternalDataShareResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID string, createExternalDataShareRequest fabcore.CreateExternalDataShareRequest, _ *fabcore.ExternalDataSharesProviderClientCreateExternalDataShareOptions) (resp azfake.Responder[fabcore.ExternalDataSharesProviderClientCreateExternalDataShareResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.ExternalDataSharesProviderClientCreateExternalDataShareResponse]{}

		entity := NewRandomExternalDataShare(workspaceID)
		entity.Paths = createExternalDataShareRequest.Paths
		entity.Recipient = createExternalDataShareRequest.Recipient
		entity.ItemID = to.Ptr(itemID)

		fakeTestUpsert(entity)
		resp.SetResponse(http.StatusCreated, fabcore.ExternalDataSharesProviderClientCreateExternalDataShareResponse{ExternalDataShare: entity}, nil)

		return resp, errResp
	}
}

func fakeDeleteExternalDataShareProvider() func(ctx context.Context, workspaceID, itemID, id string, options *fabcore.ExternalDataSharesProviderClientDeleteExternalDataShareOptions) (resp azfake.Responder[fabcore.ExternalDataSharesProviderClientDeleteExternalDataShareResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _, id string, _ *fabcore.ExternalDataSharesProviderClientDeleteExternalDataShareOptions) (resp azfake.Responder[fabcore.ExternalDataSharesProviderClientDeleteExternalDataShareResponse], errResp azfake.ErrorResponder) {
		if _, ok := fakeExternalDataShareStore[id]; ok {
			delete(fakeExternalDataShareStore, id)
			resp.SetResponse(http.StatusOK, struct{}{}, nil)
		} else {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, "ItemNotFound", "Item not found"))
			resp.SetResponse(http.StatusNotFound, struct{}{}, nil)
		}

		return resp, errResp
	}
}

func GetAllStoredExternalDataShares() []fabcore.ExternalDataShare {
	externalDataShares := make([]fabcore.ExternalDataShare, 0, len(fakeExternalDataShareStore))
	for _, externalDataShare := range fakeExternalDataShareStore {
		externalDataShares = append(externalDataShares, externalDataShare)
	}

	return externalDataShares
}

func fakeTestUpsert(entity fabcore.ExternalDataShare) {
	fakeExternalDataShareStore[*entity.ID] = entity
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
		ExpirationTimeUTC:  to.Ptr(time.Now()),
		ItemID:             to.Ptr(testhelp.RandomUUID()),
		InvitationURL:      to.Ptr(testhelp.RandomName()),
		AcceptedByTenantID: to.Ptr(testhelp.RandomUUID()),
	}
}
