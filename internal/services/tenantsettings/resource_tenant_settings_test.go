// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package tenantsettings_test

import (
	"regexp"
	"strconv"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_TenantSettingsResource_Attributes(t *testing.T) {
	randomEntity := NewRandomTenantSetting()
	enableUpdate := !*randomEntity.Enabled
	fakes.FakeServer.ServerFactory.Admin.TenantsServer.NewListTenantSettingsPager = fakeTenantSettingFunc()
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
		// success
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"setting_name": *randomEntity.SettingName,
					"enabled":      enableUpdate,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "setting_name", randomEntity.SettingName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "enabled", strconv.FormatBool(enableUpdate)),
			),
		},
	}))
}
