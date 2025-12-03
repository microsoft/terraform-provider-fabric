// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package tenantsettings_test

import (
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"

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
