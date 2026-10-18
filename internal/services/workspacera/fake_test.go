// Copyright Microsoft Corporation 2026
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

		return resp, errResp
	}
}

func NewRandomWorkspaceRoleAssignment() fabcore.WorkspaceRoleAssignment {
	return fabcore.WorkspaceRoleAssignment{
		ID: new(testhelp.RandomUUID()),
		Principal: &fabcore.UserPrincipal{
			ID:          new(testhelp.RandomUUID()),
			Type:        azto.Ptr(fabcore.PrincipalTypeUser),
			DisplayName: new(testhelp.RandomName()),
			UserDetails: &fabcore.UserPrincipalUserDetails{
				UserPrincipalName: new(testhelp.RandomName()),
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

		return resp
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
				ID:   new(principal0ID),
				Role: azto.Ptr(fabcore.WorkspaceRoleAdmin),
				Principal: &fabcore.GroupPrincipal{
					ID:          new(principal0ID),
					Type:        azto.Ptr(fabcore.PrincipalTypeGroup),
					DisplayName: new(testhelp.RandomName()),
					GroupDetails: &fabcore.GroupPrincipalGroupDetails{
						GroupType: azto.Ptr(fabcore.GroupTypeSecurityGroup),
					},
				},
			},
			{
				ID:   new(principal1ID),
				Role: azto.Ptr(fabcore.WorkspaceRoleMember),
				Principal: &fabcore.UserPrincipal{
					ID:          new(principal1ID),
					Type:        azto.Ptr(fabcore.PrincipalTypeUser),
					DisplayName: new(testhelp.RandomName()),
					UserDetails: &fabcore.UserPrincipalUserDetails{
						UserPrincipalName: new(testhelp.RandomName()),
					},
				},
			},
			{
				ID:   new(principal2ID),
				Role: azto.Ptr(fabcore.WorkspaceRoleMember),
				Principal: &fabcore.ServicePrincipal{
					ID:          new(principal2ID),
					Type:        azto.Ptr(fabcore.PrincipalTypeServicePrincipal),
					DisplayName: new(testhelp.RandomName()),
					ServicePrincipalDetails: &fabcore.ServicePrincipalDetails{
						AADAppID: new(testhelp.RandomUUID()),
					},
				},
			},
			{
				ID:   new(principal3ID),
				Role: azto.Ptr(fabcore.WorkspaceRoleViewer),
				Principal: &fabcore.ServicePrincipalProfilePrincipal{
					ID:          new(principal3ID),
					Type:        azto.Ptr(fabcore.PrincipalTypeServicePrincipalProfile),
					DisplayName: new(testhelp.RandomName()),
					ServicePrincipalProfileDetails: &fabcore.ServicePrincipalProfilePrincipalServicePrincipalProfileDetails{
						ParentPrincipal: &fabcore.ServicePrincipal{
							ID:          new(principal2ID),
							Type:        azto.Ptr(fabcore.PrincipalTypeServicePrincipal),
							DisplayName: new(testhelp.RandomName()),
							ServicePrincipalDetails: &fabcore.ServicePrincipalDetails{
								AADAppID: new(testhelp.RandomUUID()),
							},
						},
					},
				},
			},
		},
	}
}
