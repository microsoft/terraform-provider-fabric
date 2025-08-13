// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func fakeCreateOrUpdateOneLakeDataAccessSecurity() func(ctx context.Context, workspaceID, itemID string, createOrUpdateDataAccessRolesRequest fabcore.CreateOrUpdateDataAccessRolesRequest, options *fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateDataAccessRolesOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateDataAccessRolesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _ string, _ fabcore.CreateOrUpdateDataAccessRolesRequest, _ *fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateDataAccessRolesOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateDataAccessRolesResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateDataAccessRolesResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateDataAccessRolesResponse{Etag: to.Ptr(testhelp.RandomName())}, nil)

		return
	}
}
