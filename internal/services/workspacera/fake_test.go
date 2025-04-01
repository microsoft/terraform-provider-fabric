// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspacera_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func fakeWorkspaceRoleAssignment(
	exampleResp fabcore.WorkspaceRoleAssignment,
) func(ctx context.Context, workspaceID, workspaceRoleAssignmentID string, options *fabcore.WorkspacesClientGetWorkspaceRoleAssignmentOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetWorkspaceRoleAssignmentResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _ string, _ *fabcore.WorkspacesClientGetWorkspaceRoleAssignmentOptions) (resp azfake.Responder[fabcore.WorkspacesClientGetWorkspaceRoleAssignmentResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.WorkspacesClientGetWorkspaceRoleAssignmentResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.WorkspacesClientGetWorkspaceRoleAssignmentResponse{WorkspaceRoleAssignment: exampleResp}, nil)

		return
	}
}

func NewRandomWorkspaceRoleAssignment() fabcore.WorkspaceRoleAssignment {
	return fabcore.WorkspaceRoleAssignment{
		ID: azto.Ptr(testhelp.RandomUUID()),
		Principal: &fabcore.Principal{
			ID:          azto.Ptr(testhelp.RandomUUID()),
			Type:        azto.Ptr(fabcore.PrincipalTypeUser),
			DisplayName: azto.Ptr(testhelp.RandomName()),
			UserDetails: &fabcore.PrincipalUserDetails{
				UserPrincipalName: azto.Ptr(testhelp.RandomName()),
			},
		},
		Role: azto.Ptr(fabcore.WorkspaceRoleAdmin),
	}
}

func fakeWorkspaceRoleAssignments(
	exampleResp fabcore.WorkspaceRoleAssignments,
) func(workspaceID string, options *fabcore.WorkspacesClientListWorkspaceRoleAssignmentsOptions) (resp azfake.PagerResponder[fabcore.WorkspacesClientListWorkspaceRoleAssignmentsResponse]) {
	return func(_ string, _ *fabcore.WorkspacesClientListWorkspaceRoleAssignmentsOptions) (resp azfake.PagerResponder[fabcore.WorkspacesClientListWorkspaceRoleAssignmentsResponse]) {
		resp = azfake.PagerResponder[fabcore.WorkspacesClientListWorkspaceRoleAssignmentsResponse]{}
		resp.AddPage(http.StatusOK, fabcore.WorkspacesClientListWorkspaceRoleAssignmentsResponse{WorkspaceRoleAssignments: exampleResp}, nil)

		return
	}
}

func NewRandomWorkspaceRoleAssignments() fabcore.WorkspaceRoleAssignments {
	principal0ID := testhelp.RandomUUID()
	principal1ID := testhelp.RandomUUID()
	principal2ID := testhelp.RandomUUID()
	principal3ID := testhelp.RandomUUID()

	return fabcore.WorkspaceRoleAssignments{
		Value: []fabcore.WorkspaceRoleAssignment{
			{
				ID:   azto.Ptr(principal0ID),
				Role: azto.Ptr(fabcore.WorkspaceRoleAdmin),
				Principal: &fabcore.Principal{
					ID:          azto.Ptr(principal0ID),
					Type:        azto.Ptr(fabcore.PrincipalTypeGroup),
					DisplayName: azto.Ptr(testhelp.RandomName()),
					GroupDetails: &fabcore.PrincipalGroupDetails{
						GroupType: azto.Ptr(fabcore.GroupTypeSecurityGroup),
					},
				},
			},
			{
				ID:   azto.Ptr(principal1ID),
				Role: azto.Ptr(fabcore.WorkspaceRoleMember),
				Principal: &fabcore.Principal{
					ID:          azto.Ptr(principal1ID),
					Type:        azto.Ptr(fabcore.PrincipalTypeUser),
					DisplayName: azto.Ptr(testhelp.RandomName()),
					UserDetails: &fabcore.PrincipalUserDetails{
						UserPrincipalName: azto.Ptr(testhelp.RandomName()),
					},
				},
			},
			{
				ID:   azto.Ptr(principal2ID),
				Role: azto.Ptr(fabcore.WorkspaceRoleMember),
				Principal: &fabcore.Principal{
					ID:          azto.Ptr(principal2ID),
					Type:        azto.Ptr(fabcore.PrincipalTypeServicePrincipal),
					DisplayName: azto.Ptr(testhelp.RandomName()),
					ServicePrincipalDetails: &fabcore.PrincipalServicePrincipalDetails{
						AADAppID: azto.Ptr(testhelp.RandomUUID()),
					},
				},
			},
			{
				ID:   azto.Ptr(principal3ID),
				Role: azto.Ptr(fabcore.WorkspaceRoleViewer),
				Principal: &fabcore.Principal{
					ID:          azto.Ptr(principal3ID),
					Type:        azto.Ptr(fabcore.PrincipalTypeServicePrincipalProfile),
					DisplayName: azto.Ptr(testhelp.RandomName()),
					ServicePrincipalProfileDetails: &fabcore.PrincipalServicePrincipalProfileDetails{
						ParentPrincipal: &fabcore.Principal{
							ID:          azto.Ptr(principal2ID),
							Type:        azto.Ptr(fabcore.PrincipalTypeServicePrincipal),
							DisplayName: azto.Ptr(testhelp.RandomName()),
							ServicePrincipalDetails: &fabcore.PrincipalServicePrincipalDetails{
								AADAppID: azto.Ptr(testhelp.RandomUUID()),
							},
						},
					},
				},
			},
		},
	}
}
