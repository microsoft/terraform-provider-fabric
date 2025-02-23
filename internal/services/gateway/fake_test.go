// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway_test

import (
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// Replace "Workspace" with "Gateway" in types and function names.
func fakeGatewayRoleAssignments(exampleResp fabcore.GatewayRoleAssignments) func(gatewayID string, options *fabcore.GatewaysClientListGatewayRoleAssignmentsOptions) (resp azfake.PagerResponder[fabcore.GatewaysClientListGatewayRoleAssignmentsResponse]) {
	return func(_ string, _ *fabcore.GatewaysClientListGatewayRoleAssignmentsOptions) (resp azfake.PagerResponder[fabcore.GatewaysClientListGatewayRoleAssignmentsResponse]) {
		resp = azfake.PagerResponder[fabcore.GatewaysClientListGatewayRoleAssignmentsResponse]{}
		resp.AddPage(http.StatusOK, fabcore.GatewaysClientListGatewayRoleAssignmentsResponse{GatewayRoleAssignments: exampleResp}, nil)
		return
	}
}

// NewRandomGatewayRoleAssignments creates a random GatewayRoleAssignments object for testing.
func NewRandomGatewayRoleAssignments() fabcore.GatewayRoleAssignments {
	// Generate random IDs for two role assignments.
	assignmentID0 := testhelp.RandomUUID()
	assignmentID1 := testhelp.RandomUUID()
	principalID0 := testhelp.RandomUUID()
	principalID1 := testhelp.RandomUUID()

	return fabcore.GatewayRoleAssignments{
		Value: []fabcore.GatewayRoleAssignment{
			{
				ID:   azto.Ptr(assignmentID0),
				Role: azto.Ptr(fabcore.GatewayRoleAdmin), // assuming GatewayRoleAdmin exists
				Principal: &fabcore.Principal{
					ID:          azto.Ptr(principalID0),
					Type:        azto.Ptr(fabcore.PrincipalTypeGroup),
					DisplayName: azto.Ptr(testhelp.RandomName()),
					GroupDetails: &fabcore.PrincipalGroupDetails{
						GroupType: azto.Ptr(fabcore.GroupTypeSecurityGroup),
					},
				},
			},
			{
				ID:   azto.Ptr(assignmentID1),
				Role: azto.Ptr(fabcore.GatewayRoleConnectionCreator),
				Principal: &fabcore.Principal{
					ID:          azto.Ptr(principalID1),
					Type:        azto.Ptr(fabcore.PrincipalTypeUser),
					DisplayName: azto.Ptr(testhelp.RandomName()),
					UserDetails: &fabcore.PrincipalUserDetails{
						UserPrincipalName: azto.Ptr(testhelp.RandomName()),
					},
				},
			},
		},
		ContinuationToken: nil,
		ContinuationURI:   nil,
	}
}
