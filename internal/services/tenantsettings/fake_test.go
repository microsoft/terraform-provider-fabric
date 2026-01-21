// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package tenantsettings_test

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

var fakeTenantSettingsStore = map[string]fabadmin.TenantSetting{}

func fakeTenantSettingFunc() func(options *fabadmin.TenantsClientListTenantSettingsOptions) (resp azfake.PagerResponder[fabadmin.TenantsClientListTenantSettingsResponse]) {
	return func(_ *fabadmin.TenantsClientListTenantSettingsOptions) (resp azfake.PagerResponder[fabadmin.TenantsClientListTenantSettingsResponse]) {
		resp = azfake.PagerResponder[fabadmin.TenantsClientListTenantSettingsResponse]{}
		resp.AddPage(http.StatusOK, fabadmin.TenantsClientListTenantSettingsResponse{TenantSettings: fabadmin.TenantSettings{Value: GetAllStoredTenantSettings()}}, nil)

		return resp
	}
}

func NewRandomTenantSettingsWithoutProperties() fabadmin.TenantSetting {
	return fabadmin.TenantSetting{
		SettingName:              to.Ptr(testhelp.RandomName()),
		TenantSettingGroup:       to.Ptr(testhelp.RandomName()),
		Title:                    to.Ptr(testhelp.RandomName()),
		CanSpecifySecurityGroups: to.Ptr(testhelp.RandomBool()),
		Enabled:                  to.Ptr(testhelp.RandomBool()),
		DelegateToCapacity:       to.Ptr(testhelp.RandomBool()),
		DelegateToDomain:         to.Ptr(testhelp.RandomBool()),
		DelegateToWorkspace:      to.Ptr(testhelp.RandomBool()),
	}
}

func NewRandomTenantSetting() fabadmin.TenantSetting {
	return fabadmin.TenantSetting{
		SettingName:              to.Ptr(testhelp.RandomName()),
		TenantSettingGroup:       to.Ptr(testhelp.RandomName()),
		Title:                    to.Ptr(testhelp.RandomName()),
		CanSpecifySecurityGroups: to.Ptr(testhelp.RandomBool()),
		Enabled:                  to.Ptr(testhelp.RandomBool()),
		DelegateToCapacity:       to.Ptr(testhelp.RandomBool()),
		DelegateToDomain:         to.Ptr(testhelp.RandomBool()),
		DelegateToWorkspace:      to.Ptr(testhelp.RandomBool()),
		EnabledSecurityGroups: []fabadmin.TenantSettingSecurityGroup{
			{
				GraphID: to.Ptr(testhelp.RandomUUID()),
				Name:    to.Ptr(testhelp.RandomName()),
			},
			{
				GraphID: to.Ptr(testhelp.RandomUUID()),
				Name:    to.Ptr(testhelp.RandomName()),
			},
		},
		ExcludedSecurityGroups: []fabadmin.TenantSettingSecurityGroup{
			{
				GraphID: to.Ptr(testhelp.RandomUUID()),
				Name:    to.Ptr(testhelp.RandomName()),
			},
			{
				GraphID: to.Ptr(testhelp.RandomUUID()),
				Name:    to.Ptr(testhelp.RandomName()),
			},
		},
		Properties: []fabadmin.TenantSettingProperty{
			{
				Name:  to.Ptr(testhelp.RandomName()),
				Type:  to.Ptr(fabadmin.TenantSettingPropertyTypeBoolean),
				Value: to.Ptr(testhelp.RandomName()),
			},
			{
				Name:  to.Ptr(testhelp.RandomName()),
				Type:  to.Ptr(fabadmin.TenantSettingPropertyTypeFreeText),
				Value: to.Ptr(testhelp.RandomName()),
			},
		},
	}
}

func GetAllStoredTenantSettings() []fabadmin.TenantSetting {
	tenantSettings := make([]fabadmin.TenantSetting, 0, len(fakeTenantSettingsStore))
	for _, tenantSetting := range fakeTenantSettingsStore {
		tenantSettings = append(tenantSettings, tenantSetting)
	}

	return tenantSettings
}

func fakeTestUpsert(entity fabadmin.TenantSetting) {
	fakeTenantSettingsStore[*entity.SettingName] = entity
}

func fakeUpdateTenantSettings() func(ctx context.Context, tenantSettingName string, updateTenantSettingRequest fabadmin.UpdateTenantSettingRequest, options *fabadmin.TenantsClientUpdateTenantSettingOptions) (resp azfake.Responder[fabadmin.TenantsClientUpdateTenantSettingResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, tenantSettingName string, updateTenantSettingRequest fabadmin.UpdateTenantSettingRequest, _ *fabadmin.TenantsClientUpdateTenantSettingOptions) (resp azfake.Responder[fabadmin.TenantsClientUpdateTenantSettingResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabadmin.TenantsClientUpdateTenantSettingResponse]{}

		errItemNotFound := fabcore.ErrItem.ItemNotFound.Error()

		if _, ok := fakeTenantSettingsStore[tenantSettingName]; !ok {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, errItemNotFound, "Tenant Setting not found"))
			resp.SetResponse(http.StatusNotFound, fabadmin.TenantsClientUpdateTenantSettingResponse{}, nil)

			return resp, errResp
		}

		fakeTenantSettingsStore[tenantSettingName] = fabadmin.TenantSetting{
			SettingName:              to.Ptr(tenantSettingName),
			Title:                    fakeTenantSettingsStore[tenantSettingName].Title,
			TenantSettingGroup:       fakeTenantSettingsStore[tenantSettingName].TenantSettingGroup,
			Enabled:                  updateTenantSettingRequest.Enabled,
			CanSpecifySecurityGroups: fakeTenantSettingsStore[tenantSettingName].CanSpecifySecurityGroups,
			DelegateToCapacity:       updateTenantSettingRequest.DelegateToCapacity,
			DelegateToDomain:         updateTenantSettingRequest.DelegateToDomain,
			DelegateToWorkspace:      updateTenantSettingRequest.DelegateToWorkspace,
			EnabledSecurityGroups:    updateTenantSettingRequest.EnabledSecurityGroups,
			ExcludedSecurityGroups:   updateTenantSettingRequest.ExcludedSecurityGroups,
			Properties:               updateTenantSettingRequest.Properties,
		}

		resp.SetResponse(
			http.StatusOK,
			fabadmin.TenantsClientUpdateTenantSettingResponse{UpdateTenantSettingResponse: fabadmin.UpdateTenantSettingResponse{TenantSettings: GetAllStoredTenantSettings()}},
			nil,
		)

		return resp, errResp
	}
}
