// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connectionra_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func fakeConnectionRoleAssignment(
	exampleResp fabcore.ConnectionRoleAssignment,
) func(ctx context.Context, connectionID, connectionRoleAssignmentID string, options *fabcore.ConnectionsClientGetConnectionRoleAssignmentOptions) (resp azfake.Responder[fabcore.ConnectionsClientGetConnectionRoleAssignmentResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _ string, _ *fabcore.ConnectionsClientGetConnectionRoleAssignmentOptions) (resp azfake.Responder[fabcore.ConnectionsClientGetConnectionRoleAssignmentResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.ConnectionsClientGetConnectionRoleAssignmentResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.ConnectionsClientGetConnectionRoleAssignmentResponse{ConnectionRoleAssignment: exampleResp}, nil)

		return
	}
}

func NewRandomConnectionRoleAssignment() fabcore.ConnectionRoleAssignment {
	return fabcore.ConnectionRoleAssignment{
		ID: azto.Ptr(testhelp.RandomUUID()),
		Principal: &fabcore.Principal{
			ID:   azto.Ptr(testhelp.RandomUUID()),
			Type: azto.Ptr(fabcore.PrincipalTypeUser),
		},
		Role: azto.Ptr(fabcore.ConnectionRoleOwner),
	}
}

func fakeConnectionRoleAssignments(
	exampleResp fabcore.ConnectionRoleAssignments,
) func(connectionID string, options *fabcore.ConnectionsClientListConnectionRoleAssignmentsOptions) (resp azfake.PagerResponder[fabcore.ConnectionsClientListConnectionRoleAssignmentsResponse]) {
	return func(_ string, _ *fabcore.ConnectionsClientListConnectionRoleAssignmentsOptions) (resp azfake.PagerResponder[fabcore.ConnectionsClientListConnectionRoleAssignmentsResponse]) {
		resp = azfake.PagerResponder[fabcore.ConnectionsClientListConnectionRoleAssignmentsResponse]{}
		resp.AddPage(http.StatusOK, fabcore.ConnectionsClientListConnectionRoleAssignmentsResponse{ConnectionRoleAssignments: exampleResp}, nil)

		return
	}
}

func NewRandomConnectionRoleAssignments() fabcore.ConnectionRoleAssignments {
	principal0ID := testhelp.RandomUUID()
	principal1ID := testhelp.RandomUUID()
	principal2ID := testhelp.RandomUUID()

	return fabcore.ConnectionRoleAssignments{
		Value: []fabcore.ConnectionRoleAssignment{
			{
				ID:   azto.Ptr(principal0ID),
				Role: azto.Ptr(fabcore.ConnectionRoleOwner),
				Principal: &fabcore.Principal{
					ID:   azto.Ptr(principal0ID),
					Type: azto.Ptr(fabcore.PrincipalTypeGroup),
				},
			},
			{
				ID:   azto.Ptr(principal1ID),
				Role: azto.Ptr(fabcore.ConnectionRoleUser),
				Principal: &fabcore.Principal{
					ID:   azto.Ptr(principal1ID),
					Type: azto.Ptr(fabcore.PrincipalTypeUser),
				},
			},
			{
				ID:   azto.Ptr(principal2ID),
				Role: azto.Ptr(fabcore.ConnectionRoleUserWithReshare),
				Principal: &fabcore.Principal{
					ID:   azto.Ptr(principal2ID),
					Type: azto.Ptr(fabcore.PrincipalTypeServicePrincipal),
				},
			},
		},
	}
}
