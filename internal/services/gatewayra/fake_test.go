// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gatewayra_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func fakeGatewayRoleAssignment(
	exampleResp fabcore.GatewayRoleAssignment,
) func(ctx context.Context, workspaceID, workspaceRoleAssignmentID string, options *fabcore.GatewaysClientGetGatewayRoleAssignmentOptions) (resp azfake.Responder[fabcore.GatewaysClientGetGatewayRoleAssignmentResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _ string, _ *fabcore.GatewaysClientGetGatewayRoleAssignmentOptions) (resp azfake.Responder[fabcore.GatewaysClientGetGatewayRoleAssignmentResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.GatewaysClientGetGatewayRoleAssignmentResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.GatewaysClientGetGatewayRoleAssignmentResponse{GatewayRoleAssignment: exampleResp}, nil)

		return
	}
}

func NewRandomGatewayRoleAssignment() fabcore.GatewayRoleAssignment {
	return fabcore.GatewayRoleAssignment{
		ID: azto.Ptr(testhelp.RandomUUID()),
		Principal: &fabcore.Principal{
			ID:   azto.Ptr(testhelp.RandomUUID()),
			Type: azto.Ptr(fabcore.PrincipalTypeUser),
		},
		Role: azto.Ptr(fabcore.GatewayRoleAdmin),
	}
}

func fakeGatewayRoleAssignments(
	exampleResp fabcore.GatewayRoleAssignments,
) func(gatewayID string, options *fabcore.GatewaysClientListGatewayRoleAssignmentsOptions) (resp azfake.PagerResponder[fabcore.GatewaysClientListGatewayRoleAssignmentsResponse]) {
	return func(_ string, _ *fabcore.GatewaysClientListGatewayRoleAssignmentsOptions) (resp azfake.PagerResponder[fabcore.GatewaysClientListGatewayRoleAssignmentsResponse]) {
		resp = azfake.PagerResponder[fabcore.GatewaysClientListGatewayRoleAssignmentsResponse]{}
		resp.AddPage(http.StatusOK, fabcore.GatewaysClientListGatewayRoleAssignmentsResponse{GatewayRoleAssignments: exampleResp}, nil)

		return
	}
}

func NewRandomGatewayRoleAssignments() fabcore.GatewayRoleAssignments {
	principal0ID := testhelp.RandomUUID()
	principal1ID := testhelp.RandomUUID()
	principal2ID := testhelp.RandomUUID()

	return fabcore.GatewayRoleAssignments{
		Value: []fabcore.GatewayRoleAssignment{
			{
				ID:   azto.Ptr(principal0ID),
				Role: azto.Ptr(fabcore.GatewayRoleAdmin),
				Principal: &fabcore.Principal{
					ID:   azto.Ptr(principal0ID),
					Type: azto.Ptr(fabcore.PrincipalTypeGroup),
				},
			},
			{
				ID:   azto.Ptr(principal1ID),
				Role: azto.Ptr(fabcore.GatewayRoleConnectionCreator),
				Principal: &fabcore.Principal{
					ID:   azto.Ptr(principal1ID),
					Type: azto.Ptr(fabcore.PrincipalTypeUser),
				},
			},
			{
				ID:   azto.Ptr(principal2ID),
				Role: azto.Ptr(fabcore.GatewayRoleConnectionCreatorWithResharing),
				Principal: &fabcore.Principal{
					ID:   azto.Ptr(principal2ID),
					Type: azto.Ptr(fabcore.PrincipalTypeServicePrincipal),
				},
			},
		},
	}
}
