// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var fakeOneLakeDataAccessRoleStore = map[string]fabcore.DataAccessRoleListItem{}

func GenerateOneLakeDataAccessRoleID(workspaceID, itemID, roleName string) string {
	return fmt.Sprintf("%s/%s/%s", workspaceID, itemID, roleName)
}

func GetAllStoredOneLakeDataAccessRoles(workspaceID, itemID string) []fabcore.DataAccessRoleListItem {
	prefix := workspaceID + "/" + itemID + "/"
	roles := make([]fabcore.DataAccessRoleListItem, 0, len(fakeOneLakeDataAccessRoleStore))

	for key, role := range fakeOneLakeDataAccessRoleStore {
		if strings.HasPrefix(key, prefix) {
			roles = append(roles, role)
		}
	}

	return roles
}

func UpsertIntoOneLakeDataAccessRoleStore(workspaceID, itemID string, entity fabcore.DataAccessRoleListItem) {
	id := GenerateOneLakeDataAccessRoleID(workspaceID, itemID, *entity.Name)
	fakeOneLakeDataAccessRoleStore[id] = entity
}

func fakeGetDataAccessRoleFunc() func(ctx context.Context, workspaceID, itemID, roleName string, options *fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID, roleName string, _ *fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse]{}
		errEntityNotFound := fabcore.ErrCommon.EntityNotFound.Error()
		id := GenerateOneLakeDataAccessRoleID(workspaceID, itemID, roleName)

		if role, ok := fakeOneLakeDataAccessRoleStore[id]; ok {
			resp.SetResponse(http.StatusOK, fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse{
				DataAccessRoleBase: fabcore.DataAccessRoleBase{
					Name:          role.Name,
					Kind:          role.Kind,
					DecisionRules: role.DecisionRules,
					Members:       role.Members,
				},
				ETag: new(testhelp.RandomName()),
			}, nil)
		} else {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, errEntityNotFound, "Entity not found"))
			resp.SetResponse(http.StatusNotFound, fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse{}, nil)
		}

		return resp, errResp
	}
}

func fakeCreateOrUpdateSingleDataAccessRoleFunc() func(ctx context.Context, workspaceID, itemID string, body fabcore.DataAccessRoleBase, options *fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID string, body fabcore.DataAccessRoleBase, _ *fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleResponse]{}
		id := GenerateOneLakeDataAccessRoleID(workspaceID, itemID, *body.Name)

		roleID := new(testhelp.RandomUUID())
		if existing, ok := fakeOneLakeDataAccessRoleStore[id]; ok {
			roleID = existing.ID
		}

		fakeOneLakeDataAccessRoleStore[id] = fabcore.DataAccessRoleListItem{
			ID:            roleID,
			Name:          body.Name,
			Kind:          body.Kind,
			DecisionRules: body.DecisionRules,
			Members:       body.Members,
		}

		resp.SetResponse(http.StatusOK, fabcore.OneLakeDataAccessSecurityClientCreateOrUpdateSingleDataAccessRoleResponse{
			ETag: new(testhelp.RandomName()),
		}, nil)

		return resp, errResp
	}
}

func fakeDeleteDataAccessRoleFunc() func(ctx context.Context, workspaceID, itemID, roleName string, options *fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID, roleName string, _ *fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleResponse]{}
		errEntityNotFound := fabcore.ErrCommon.EntityNotFound.Error()
		id := GenerateOneLakeDataAccessRoleID(workspaceID, itemID, roleName)

		if _, ok := fakeOneLakeDataAccessRoleStore[id]; !ok {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, errEntityNotFound, "Entity not found"))
			resp.SetResponse(http.StatusNotFound, fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleResponse{}, nil)

			return resp, errResp
		}

		delete(fakeOneLakeDataAccessRoleStore, id)
		resp.SetResponse(http.StatusOK, fabcore.OneLakeDataAccessSecurityClientDeleteDataAccessRoleResponse{}, nil)

		return resp, errResp
	}
}

func fakeListDataAccessRolesFunc() func(ctx context.Context, workspaceID, itemID string, options *fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID string, _ *fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesOptions) (resp azfake.Responder[fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse]{}
		resp.SetResponse(http.StatusOK, fabcore.OneLakeDataAccessSecurityClientListDataAccessRolesResponse{
			DataAccessRoles: fabcore.DataAccessRoles{Value: GetAllStoredOneLakeDataAccessRoles(workspaceID, itemID)},
			Etag:            new(testhelp.RandomName()),
		}, nil)

		return resp, errResp
	}
}
