// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceencryption_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func fakeGetWorkspaceEncryption(
	entity *fabcore.WorkspaceEncryptionDetail,
) func(ctx context.Context, workspaceID string, options *fabcore.WorkspacesClientGetWorkspaceEncryptionOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetWorkspaceEncryptionResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.WorkspacesClientGetWorkspaceEncryptionOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetWorkspaceEncryptionResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.WorkspacesClientGetWorkspaceEncryptionResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientGetWorkspaceEncryptionResponse{
			WorkspaceEncryptionDetail: *entity,
		}, nil)

		return resp, errResp
	}
}

func fakeAssignWorkspaceEncryption(
	entity *fabcore.WorkspaceEncryptionDetail,
) func(ctx context.Context, workspaceID string, assignWorkspaceEncryptionRequest fabcore.AssignWorkspaceEncryptionRequest, options *fabcore.WorkspacesClientAssignWorkspaceEncryptionOptions) (resp azfake.Responder[fabcore.WorkspacesClientAssignWorkspaceEncryptionResponse], errResp azfake.ErrorResponder) {
	return fakeAssignWorkspaceEncryptionWithStatus(entity, fabcore.WorkspaceEncryptionStatusActive)
}

// fakeAssignWorkspaceEncryptionWithStatus lets a test drive the status that the subsequent polling GET observes.
func fakeAssignWorkspaceEncryptionWithStatus(
	entity *fabcore.WorkspaceEncryptionDetail,
	status fabcore.WorkspaceEncryptionStatus,
) func(ctx context.Context, workspaceID string, assignWorkspaceEncryptionRequest fabcore.AssignWorkspaceEncryptionRequest, options *fabcore.WorkspacesClientAssignWorkspaceEncryptionOptions) (resp azfake.Responder[fabcore.WorkspacesClientAssignWorkspaceEncryptionResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, assignWorkspaceEncryptionRequest fabcore.AssignWorkspaceEncryptionRequest, _ *fabcore.WorkspacesClientAssignWorkspaceEncryptionOptions) (resp azfake.Responder[fabcore.WorkspacesClientAssignWorkspaceEncryptionResponse], errResp azfake.ErrorResponder) {
		if entity.EncryptionDetail != nil && entity.EncryptionDetail.KeyIdentifier != nil {
			entity.PreviousEncryptionDetail = entity.EncryptionDetail
		}

		entity.EncryptionDetail = &fabcore.EncryptionDetail{
			KeyIdentifier:    assignWorkspaceEncryptionRequest.KeyIdentifier,
			EncryptionStatus: new(status),
		}

		resp = azfake.Responder[fabcore.WorkspacesClientAssignWorkspaceEncryptionResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientAssignWorkspaceEncryptionResponse{}, nil)

		return resp, errResp
	}
}

func fakeResetWorkspaceEncryption(
	entity *fabcore.WorkspaceEncryptionDetail,
) func(ctx context.Context, workspaceID string, options *fabcore.WorkspacesClientResetWorkspaceEncryptionOptions) (resp azfake.Responder[fabcore.WorkspacesClientResetWorkspaceEncryptionResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ *fabcore.WorkspacesClientResetWorkspaceEncryptionOptions) (resp azfake.Responder[fabcore.WorkspacesClientResetWorkspaceEncryptionResponse], errResp azfake.ErrorResponder) {
		if entity.EncryptionDetail != nil && entity.EncryptionDetail.KeyIdentifier != nil {
			entity.PreviousEncryptionDetail = entity.EncryptionDetail
		}

		entity.EncryptionDetail = &fabcore.EncryptionDetail{
			EncryptionStatus: azto.Ptr(fabcore.WorkspaceEncryptionStatusDisabled),
		}

		resp = azfake.Responder[fabcore.WorkspacesClientResetWorkspaceEncryptionResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientResetWorkspaceEncryptionResponse{}, nil)

		return resp, errResp
	}
}

func NewRandomWorkspaceEncryptionDetail() fabcore.WorkspaceEncryptionDetail {
	return fabcore.WorkspaceEncryptionDetail{
		EncryptionDetail: &fabcore.EncryptionDetail{
			KeyIdentifier:    new(NewRandomKeyIdentifier()),
			EncryptionStatus: azto.Ptr(fabcore.WorkspaceEncryptionStatusActive),
		},
		WorkspaceEncryptionItemsDetails: []fabcore.WorkspaceEncryptionItemsDetail{
			{
				EncryptionStatus: azto.Ptr(fabcore.WorkspaceEncryptionStatusActive),
				Items: []fabcore.WorkspaceEncryptionItem{
					{
						ID:          new(testhelp.RandomUUID()),
						DisplayName: new(testhelp.RandomName()),
						Type:        new("Lakehouse"),
					},
				},
			},
		},
	}
}

// fakeGetWorkspaceEncryptionInProgress reports an in-progress status for the first inProgressCalls responses,
// forcing the polling loop to iterate before the entity settles.
func fakeGetWorkspaceEncryptionInProgress(
	entity *fabcore.WorkspaceEncryptionDetail,
	inProgressCalls int,
) func(ctx context.Context, workspaceID string, options *fabcore.WorkspacesClientGetWorkspaceEncryptionOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetWorkspaceEncryptionResponse], errResp azfake.ErrorResponder) {
	var calls int

	return func(_ context.Context, _ string, _ *fabcore.WorkspacesClientGetWorkspaceEncryptionOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetWorkspaceEncryptionResponse], errResp azfake.ErrorResponder) {
		current := *entity

		if calls < inProgressCalls {
			calls++

			var keyIdentifier *string
			if entity.EncryptionDetail != nil {
				keyIdentifier = entity.EncryptionDetail.KeyIdentifier
			}

			current = fabcore.WorkspaceEncryptionDetail{
				EncryptionDetail: &fabcore.EncryptionDetail{
					KeyIdentifier:    keyIdentifier,
					EncryptionStatus: azto.Ptr(fabcore.WorkspaceEncryptionStatusEnableInProgress),
				},
			}
		}

		resp = azfake.Responder[fabcore.WorkspacesClientGetWorkspaceEncryptionResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientGetWorkspaceEncryptionResponse{
			WorkspaceEncryptionDetail: current,
		}, nil)

		return resp, errResp
	}
}

func NewRandomKeyIdentifier() string {
	return fmt.Sprintf("https://%s.vault.azure.net/keys/%s/", strings.ToLower(testhelp.RandomName()), strings.ToLower(testhelp.RandomName()))
}
