// Copyright (c) 2026 Microsoft Corporation
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

var testDataSourceItemsFQN, testDataSourceItemsHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Types, "test")

func TestUnit_TenantSettingsDataSource(t *testing.T) {
	fakeTestUpsert(NewRandomTenantSetting())
	fakeTestUpsert(NewRandomTenantSetting())
	fakeTestUpsert(NewRandomTenantSetting())

	fakes.FakeServer.ServerFactory.Admin.TenantsServer.NewListTenantSettingsPager = fakeTenantSettingFunc()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.setting_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.enabled"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.tenant_setting_group"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.title"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.enabled_security_groups.0.graph_id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.enabled_security_groups.0.name"),
			),
		},
	}))
}

func TestAcc_TenantSettingsDataSource(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.setting_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.enabled"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.tenant_setting_group"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.title"),
			),
		},
	},
	))
}
