// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedefaulttier_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_OneLakeDefaultTierDataSource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.NewFakeServer().ServerFactory, nil, []resource.TestStep{
		{
			Config:      at.CompileConfig(testDataSourceItemHeader, map[string]any{}),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{"workspace_id": "invalid-uuid"},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_OneLakeDefaultTierDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	tier := fabcore.OneLakeAccessTierCool
	state := &oneLakeDefaultTierState{defaultTier: &tier}
	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(state)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testDataSourceItemFQN,
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{"workspace_id": workspaceID},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "default_tier", string(fabcore.OneLakeAccessTierCool)),
			),
		},
	}))
}

func TestUnit_OneLakeDefaultTierDataSource_AbsentLifecycle(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(&oneLakeDefaultTierState{})

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testDataSourceItemFQN,
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{"workspace_id": workspaceID},
			),
			Check: resource.TestCheckNoResourceAttr(testDataSourceItemFQN, "default_tier"),
		},
	}))
}

func TestUnit_OneLakeDefaultTierDataSource_GetError(t *testing.T) {
	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettingsError()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testDataSourceItemFQN,
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{"workspace_id": testhelp.RandomUUID()},
			),
			ExpectError: regexp.MustCompile(`EntityNotFound`),
		},
	}))
}

func TestAcc_OneLakeDefaultTierDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceOneLakeTierDS"].(map[string]any)
	workspaceID := workspace["id"].(string)
	expectedTier := workspace["defaultTier"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testDataSourceItemFQN, nil, []resource.TestStep{
		{
			ResourceName: testDataSourceItemFQN,
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{"workspace_id": workspaceID},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "default_tier", expectedTier),
			),
		},
	}))
}
