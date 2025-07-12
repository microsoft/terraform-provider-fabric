// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

func fakeCreateOrUpdateOneLakeDataAccessSecurity() func(ctx context.Context, workspaceID, itemID string, createOrUpdateDataAccessRolesRequest fabcore.CreateOrUpdateDataAccessRolesRequest, options *fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateDataAccessRolesOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateDataAccessRolesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _ string, _ fabcore.CreateOrUpdateDataAccessRolesRequest, _ *fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateDataAccessRolesOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateDataAccessRolesResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateDataAccessRolesResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateDataAccessRolesResponse{Etag: to.Ptr("123")}, nil)

		return
	}
}

func fakeListOneLakeDataAccessSecurity(
	exampleResp fabcore.DataAccessRoles,
) func(ctx context.Context, workspaceID, itemID string, options *fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _ string, _ *fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse{
			DataAccessRoles: exampleResp,
			Etag:            to.Ptr("123"),
		}, nil)

		return
	}
}
