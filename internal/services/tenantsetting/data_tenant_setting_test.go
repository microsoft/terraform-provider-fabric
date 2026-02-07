// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package tenantsetting_test

import (
	"strconv"
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
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "enabled", strconv.FormatBool(*randomEntity.Enabled)),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "tenant_setting_group", randomEntity.TenantSettingGroup),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "enabled_security_groups.0.graph_id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "enabled_security_groups.0.name"),
			),
		},
	}))
}

func TestAcc_TenantSettingDataSource(t *testing.T) {
	entity := testhelp.WellKnown()["TenantSettings"].(map[string]any)
	settingName := entity["settingName"].(string)
	tenantSettingGroup := entity["tenantSettingGroup"].(string)
	enabled := entity["enabled"].(bool)
	title := entity["title"].(string)
	canSpecifySecurityGroups := entity["canSpecifySecurityGroups"].(bool)

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
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "tenant_setting_group", tenantSettingGroup),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "enabled", strconv.FormatBool(enabled)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "title", title),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "can_specify_security_groups", strconv.FormatBool(canSpecifySecurityGroups)),
				resource.TestCheckNoResourceAttr(testDataSourceItemFQN, "delegate_to_capacity"),
				resource.TestCheckNoResourceAttr(testDataSourceItemFQN, "delegate_to_workspace"),
				resource.TestCheckNoResourceAttr(testDataSourceItemFQN, "delegate_to_domain"),
			),
		},
	},
	))
}
