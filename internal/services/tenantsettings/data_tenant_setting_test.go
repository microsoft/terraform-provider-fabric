// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package tenantsettings_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_TenantSettingDataSource(t *testing.T) {
	randomEntity := NewRandomTenantSetting()
	fakeTestUpsert(randomEntity)
	fakeTestUpsert(NewRandomTenantSetting())

	fakes.FakeServer.ServerFactory.Admin.TenantsServer.NewListTenantSettingsPager = fakeTenantSettingFunc()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// read by setting_name
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"setting_name": *randomEntity.SettingName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "setting_name", randomEntity.SettingName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "title", randomEntity.Title),
			),
		},
	}))
}

func TestAcc_TenantSettingDataSource(t *testing.T) {
	entity := testhelp.WellKnown()["TenantSettings"].(map[string]any)
	settingName := entity["settingName"].(string)
	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read by setting_name
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"setting_name": settingName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "setting_name", settingName),
			),
		},
	},
	))
}
