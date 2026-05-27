// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type oneLakeDataAccessSecurityState struct {
	currentEntity fabcore.DataAccessRoleBase
}

func newOneLakeDataAccessSecurityState(initialEntity fabcore.DataAccessRoleBase) *oneLakeDataAccessSecurityState {
	return &oneLakeDataAccessSecurityState{currentEntity: initialEntity}
}

func fakeGetDataAccessRole(
	entity fabcore.DataAccessRoleBase,
) func(ctx context.Context, workspaceID, itemID, roleName string, options *fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _, _ string, _ *fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse{
			DataAccessRoleBase: entity,
			ETag:               to.Ptr(testhelp.RandomName()),
		}, nil)

		return resp, errResp
	}
}

func fakeStatefulGetDataAccessRole(
	state *oneLakeDataAccessSecurityState,
) func(ctx context.Context, workspaceID, itemID, roleName string, options *fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _, _ string, _ *fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse{
			DataAccessRoleBase: state.currentEntity,
			ETag:               to.Ptr(testhelp.RandomName()),
		}, nil)

		return resp, errResp
	}
}

func fakeCreateOrUpdateSingleDataAccessRole(
	state *oneLakeDataAccessSecurityState,
) func(ctx context.Context, workspaceID, itemID string, body fabcore.DataAccessRoleBase, options *fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _ string, body fabcore.DataAccessRoleBase, _ *fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleResponse], errResp azfake.ErrorResponder) {
		state.currentEntity = body

		resp = azfake.Responder[fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleResponse{
			ETag: to.Ptr(testhelp.RandomName()),
		}, nil)

		return resp, errResp
	}
}

func fakeDeleteDataAccessRole() func(ctx context.Context, workspaceID, itemID, roleName string, options *fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _, _ string, _ *fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleResponse{}, nil)

		return resp, errResp
	}
}

func fakeListDataAccessRoles(
	entities []fabcore.DataAccessRoleListItem,
) func(ctx context.Context, workspaceID, itemID string, options *fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _ string, _ *fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse{
			DataAccessRoles: fabcore.DataAccessRoles{Value: entities},
			Etag:            to.Ptr(testhelp.RandomName()),
		}, nil)

		return resp, errResp
	}
}

func fakeNotFoundDataAccessRole() func(ctx context.Context, workspaceID, itemID, roleName string, options *fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _, _ string, _ *fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse]{}
		errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, fabcore.ErrCommon.EntityNotFound.Error(), "Entity not found"))

		return resp, errResp
	}
}
