// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package tenantsetting_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func newFakeTenantSettingsStore() map[string]fabadmin.TenantSetting {
	return make(map[string]fabadmin.TenantSetting)
}

func fakeTenantSettingFunc(store map[string]fabadmin.TenantSetting) func(options *fabadmin.TenantsClientListTenantSettingsOptions) (resp azfake.PagerResponder[fabadmin.TenantsClientListTenantSettingsResponse]) {
	return func(_ *fabadmin.TenantsClientListTenantSettingsOptions) (resp azfake.PagerResponder[fabadmin.TenantsClientListTenantSettingsResponse]) {
		resp = azfake.PagerResponder[fabadmin.TenantsClientListTenantSettingsResponse]{}
		resp.AddPage(http.StatusOK, fabadmin.TenantsClientListTenantSettingsResponse{TenantSettings: fabadmin.TenantSettings{Value: getAllStoredTenantSettings(store)}}, nil)

		return resp
	}
}

func NewRandomTenantSettingsWithoutProperties() fabadmin.TenantSetting {
	return fabadmin.TenantSetting{
		SettingName:              new(testhelp.RandomName()),
		TenantSettingGroup:       new(testhelp.RandomName()),
		Title:                    new(testhelp.RandomName()),
		CanSpecifySecurityGroups: new(testhelp.RandomBool()),
		Enabled:                  new(testhelp.RandomBool()),
		DelegateToCapacity:       new(testhelp.RandomBool()),
		DelegateToDomain:         new(testhelp.RandomBool()),
		DelegateToWorkspace:      new(testhelp.RandomBool()),
	}
}

func NewRandomTenantSetting() fabadmin.TenantSetting {
	return fabadmin.TenantSetting{
		SettingName:              new(testhelp.RandomName()),
		TenantSettingGroup:       new(testhelp.RandomName()),
		Title:                    new(testhelp.RandomName()),
		CanSpecifySecurityGroups: new(testhelp.RandomBool()),
		Enabled:                  new(testhelp.RandomBool()),
		DelegateToCapacity:       new(testhelp.RandomBool()),
		DelegateToDomain:         new(testhelp.RandomBool()),
		DelegateToWorkspace:      new(testhelp.RandomBool()),
		EnabledSecurityGroups: []fabadmin.TenantSettingSecurityGroup{
			{
				GraphID: new(testhelp.RandomUUID()),
				Name:    new(testhelp.RandomName()),
			},
			{
				GraphID: new(testhelp.RandomUUID()),
				Name:    new(testhelp.RandomName()),
			},
		},
		ExcludedSecurityGroups: []fabadmin.TenantSettingSecurityGroup{
			{
				GraphID: new(testhelp.RandomUUID()),
				Name:    new(testhelp.RandomName()),
			},
			{
				GraphID: new(testhelp.RandomUUID()),
				Name:    new(testhelp.RandomName()),
			},
		},
		Properties: []fabadmin.TenantSettingProperty{
			{
				Name:  new(testhelp.RandomName()),
				Type:  to.Ptr(fabadmin.TenantSettingPropertyTypeBoolean),
				Value: new(testhelp.RandomName()),
			},
			{
				Name:  new(testhelp.RandomName()),
				Type:  to.Ptr(fabadmin.TenantSettingPropertyTypeFreeText),
				Value: new(testhelp.RandomName()),
			},
		},
	}
}

func getAllStoredTenantSettings(store map[string]fabadmin.TenantSetting) []fabadmin.TenantSetting {
	tenantSettings := make([]fabadmin.TenantSetting, 0, len(store))
	for _, tenantSetting := range store {
		tenantSettings = append(tenantSettings, tenantSetting)
	}

	return tenantSettings
}

func fakeTestUpsert(store map[string]fabadmin.TenantSetting, entity fabadmin.TenantSetting) {
	store[*entity.SettingName] = entity
}

func fakeUpdateTenantSettings(store map[string]fabadmin.TenantSetting) func(ctx context.Context, tenantSettingName string, updateTenantSettingRequest fabadmin.UpdateTenantSettingRequest, options *fabadmin.TenantsClientUpdateTenantSettingOptions) (resp azfake.Responder[fabadmin.TenantsClientUpdateTenantSettingResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, tenantSettingName string, updateTenantSettingRequest fabadmin.UpdateTenantSettingRequest, _ *fabadmin.TenantsClientUpdateTenantSettingOptions) (resp azfake.Responder[fabadmin.TenantsClientUpdateTenantSettingResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabadmin.TenantsClientUpdateTenantSettingResponse]{}

		errItemNotFound := fabcore.ErrItem.ItemNotFound.Error()

		if _, ok := store[tenantSettingName]; !ok {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, errItemNotFound, "Tenant Setting not found"))
			resp.SetResponse(http.StatusNotFound, fabadmin.TenantsClientUpdateTenantSettingResponse{}, nil)

			return resp, errResp
		}

		store[tenantSettingName] = fabadmin.TenantSetting{
			SettingName:              new(tenantSettingName),
			Title:                    store[tenantSettingName].Title,
			TenantSettingGroup:       store[tenantSettingName].TenantSettingGroup,
			Enabled:                  updateTenantSettingRequest.Enabled,
			CanSpecifySecurityGroups: store[tenantSettingName].CanSpecifySecurityGroups,
			DelegateToCapacity:       updateTenantSettingRequest.DelegateToCapacity,
			DelegateToDomain:         updateTenantSettingRequest.DelegateToDomain,
			DelegateToWorkspace:      updateTenantSettingRequest.DelegateToWorkspace,
			EnabledSecurityGroups:    updateTenantSettingRequest.EnabledSecurityGroups,
			ExcludedSecurityGroups:   updateTenantSettingRequest.ExcludedSecurityGroups,
			Properties:               updateTenantSettingRequest.Properties,
		}

		resp.SetResponse(
			http.StatusOK,
			fabadmin.TenantsClientUpdateTenantSettingResponse{UpdateTenantSettingResponse: fabadmin.UpdateTenantSettingResponse{TenantSettings: getAllStoredTenantSettings(store)}},
			nil,
		)

		return resp, errResp
	}
}
