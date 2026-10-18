// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package warehousesqlauditsetting_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabwarehouse "github.com/microsoft/fabric-sdk-go/fabric/warehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var fakeSQLAuditSettingsStore = map[string]fabwarehouse.SQLAuditSettings{}

func fakeGetSQLAuditSettingsFunc() func(ctx context.Context, workspaceID, itemID string, options *fabwarehouse.SQLAuditSettingsClientGetSQLAuditSettingsOptions) (resp azfake.Responder[fabwarehouse.SQLAuditSettingsClientGetSQLAuditSettingsResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, itemID string, _ *fabwarehouse.SQLAuditSettingsClientGetSQLAuditSettingsOptions) (resp azfake.Responder[fabwarehouse.SQLAuditSettingsClientGetSQLAuditSettingsResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabwarehouse.SQLAuditSettingsClientGetSQLAuditSettingsResponse]{}

		if settings, ok := fakeSQLAuditSettingsStore[itemID]; ok {
			resp.SetResponse(http.StatusOK, fabwarehouse.SQLAuditSettingsClientGetSQLAuditSettingsResponse{SQLAuditSettings: settings}, nil)

			return resp, errResp
		}

		errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, fabcore.ErrCommon.EntityNotFound.Error(), "Entity not found"))
		resp.SetResponse(http.StatusNotFound, fabwarehouse.SQLAuditSettingsClientGetSQLAuditSettingsResponse{}, nil)

		return resp, errResp
	}
}

func fakeUpdateSQLAuditSettingsFunc() func(ctx context.Context, workspaceID, itemID string, updateReq fabwarehouse.SQLAuditSettingsUpdate, options *fabwarehouse.SQLAuditSettingsClientUpdateSQLAuditSettingsOptions) (resp azfake.Responder[fabwarehouse.SQLAuditSettingsClientUpdateSQLAuditSettingsResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, itemID string, updateReq fabwarehouse.SQLAuditSettingsUpdate, _ *fabwarehouse.SQLAuditSettingsClientUpdateSQLAuditSettingsOptions) (resp azfake.Responder[fabwarehouse.SQLAuditSettingsClientUpdateSQLAuditSettingsResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabwarehouse.SQLAuditSettingsClientUpdateSQLAuditSettingsResponse]{}

		current, ok := fakeSQLAuditSettingsStore[itemID]
		if !ok {
			current = fabwarehouse.SQLAuditSettings{
				State:         updateReq.State,
				RetentionDays: updateReq.RetentionDays,
			}
		}

		if updateReq.State != nil {
			current.State = updateReq.State
		}

		if updateReq.RetentionDays != nil {
			current.RetentionDays = updateReq.RetentionDays
		}

		fakeSQLAuditSettingsStore[itemID] = current

		resp.SetResponse(http.StatusOK, fabwarehouse.SQLAuditSettingsClientUpdateSQLAuditSettingsResponse{SQLAuditSettings: current}, nil)

		return resp, errResp
	}
}

func fakeSetAuditActionsAndGroupsFunc() func(ctx context.Context, workspaceID, itemID string, setReq []string, options *fabwarehouse.SQLAuditSettingsClientSetAuditActionsAndGroupsOptions) (resp azfake.Responder[fabwarehouse.SQLAuditSettingsClientSetAuditActionsAndGroupsResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, itemID string, setReq []string, _ *fabwarehouse.SQLAuditSettingsClientSetAuditActionsAndGroupsOptions) (resp azfake.Responder[fabwarehouse.SQLAuditSettingsClientSetAuditActionsAndGroupsResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabwarehouse.SQLAuditSettingsClientSetAuditActionsAndGroupsResponse]{}

		current, ok := fakeSQLAuditSettingsStore[itemID]
		if !ok {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, fabcore.ErrCommon.EntityNotFound.Error(), "Entity not found"))
			resp.SetResponse(http.StatusNotFound, fabwarehouse.SQLAuditSettingsClientSetAuditActionsAndGroupsResponse{}, nil)

			return resp, errResp
		}

		current.AuditActionsAndGroups = setReq
		fakeSQLAuditSettingsStore[itemID] = current

		resp.SetResponse(http.StatusOK, fabwarehouse.SQLAuditSettingsClientSetAuditActionsAndGroupsResponse{}, nil)

		return resp, errResp
	}
}

func fakeTestUpsertSQLAuditSettings(warehouseID string, settings fabwarehouse.SQLAuditSettings) {
	fakeSQLAuditSettingsStore[warehouseID] = settings
}

func NewRandomSQLAuditSettings() fabwarehouse.SQLAuditSettings {
	return fabwarehouse.SQLAuditSettings{
		State:         new(fabwarehouse.AuditSettingsStateEnabled),
		RetentionDays: new(testhelp.RandomIntRange(int32(1), int32(90))),
		AuditActionsAndGroups: []string{
			testhelp.RandomName(),
		},
	}
}
