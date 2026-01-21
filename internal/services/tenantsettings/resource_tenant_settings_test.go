// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package tenantsettings_test

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/services/tenantsettings"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_TenantSettingsResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"enabled":         true,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attribute - enabled
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"setting_name": "test",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "enabled" is required, but no definition was found.`),
		},
		// error - no required attribute - setting_name
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"enabled": true,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "setting_name" is required, but no definition was found.`),
		},
	}))
}

func TestUnit_TenantSettingsResource_CRUD(t *testing.T) {
	entity := NewRandomTenantSettingsWithoutProperties()

	fakeTestUpsert(entity)
	fakes.FakeServer.ServerFactory.Admin.TenantsServer.NewListTenantSettingsPager = fakeTenantSettingFunc()
	entityUpdate := entity
	entityUpdate.Enabled = to.Ptr(!*entity.Enabled)
	fakes.FakeServer.ServerFactory.Admin.TenantsServer.UpdateTenantSetting = fakeUpdateTenantSettings()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// success
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"setting_name":          *entity.SettingName,
					"enabled":               !*entity.Enabled,
					"delegate_to_domain":    *entity.DelegateToDomain,
					"delegate_to_capacity":  *entity.DelegateToCapacity,
					"delegate_to_workspace": *entity.DelegateToWorkspace,
					"delete_behaviour":      string(tenantsettings.NoChange),
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "setting_name", entity.SettingName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "tenant_setting_group", entity.TenantSettingGroup),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "title", entity.Title),
				resource.TestCheckResourceAttr(testResourceItemFQN, "enabled", strconv.FormatBool(*entityUpdate.Enabled)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "can_specify_security_groups", strconv.FormatBool(*entity.CanSpecifySecurityGroups)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "delegate_to_domain", strconv.FormatBool(*entity.DelegateToDomain)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "delegate_to_capacity", strconv.FormatBool(*entity.DelegateToCapacity)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "delegate_to_workspace", strconv.FormatBool(*entity.DelegateToWorkspace)),
			),
		},
	}))
}

func TestAcc_TenantSettingsResource_CRUD(t *testing.T) {
	entity := testhelp.WellKnown()["TenantSettings"].(map[string]any)
	settingName := entity["settingName"].(string)
	enabled := entity["enabled"].(bool)
	securityGroupID := entity["securityGroupId"].(string)
	securityGroupName := entity["securityGroupName"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Update
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"setting_name": settingName,
					"enabled":      !enabled,
					"enabled_security_groups": []map[string]any{
						{
							"graph_id": securityGroupID,
							"name":     securityGroupName,
						},
					},
					"delete_behaviour": string(tenantsettings.Disable),
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "setting_name", &settingName),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "tenant_setting_group"),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "title"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "enabled", strconv.FormatBool(!enabled)),
				resource.TestCheckResourceAttrSet(testResourceItemFQN, "can_specify_security_groups"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "enabled_security_groups.0.graph_id", securityGroupID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "enabled_security_groups.0.name", securityGroupName),
			),
		},
	},
	))
}
