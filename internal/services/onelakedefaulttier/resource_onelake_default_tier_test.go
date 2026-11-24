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

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_OneLakeDefaultTierResource_Attributes(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	fakeServer := fakes.NewFakeServer()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config:       at.CompileConfig(testResourceItemHeader, map[string]any{}),
			ExpectError:  regexp.MustCompile(`Missing required argument`),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "invalid-uuid",
					"default_tier": string(fabcore.OneLakeAccessTierHot),
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"default_tier": "Invalid",
				},
			),
			ExpectError: regexp.MustCompile(`Attribute default_tier value must be one of`),
		},
	}))
}

func TestUnit_OneLakeDefaultTierResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	fakeServer := fakes.NewFakeServer()
	state := &oneLakeDefaultTierState{}
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(state)
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.ModifyDefaultTier = fakeModifyDefaultTier(state)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"default_tier": string(fabcore.OneLakeAccessTierCool),
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "default_tier", string(fabcore.OneLakeAccessTierCool)),
			),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"default_tier": string(fabcore.OneLakeAccessTierCold),
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "default_tier", string(fabcore.OneLakeAccessTierCold)),
			),
		},
		{
			ResourceName:      testResourceItemFQN,
			ImportState:       true,
			ImportStateId:     workspaceID,
			ImportStateVerify: true,
		},
	}))
}

func TestUnit_OneLakeDefaultTierResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	fakeServer := fakes.NewFakeServer()
	state := &oneLakeDefaultTierState{}
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(state)
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.ModifyDefaultTier = fakeModifyDefaultTier(state)

	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id": workspaceID,
			"default_tier": string(fabcore.OneLakeAccessTierCool),
		},
	)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_OneLakeDefaultTierResource_ModifyError(t *testing.T) {
	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.ModifyDefaultTier = fakeModifyDefaultTierError()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"default_tier": string(fabcore.OneLakeAccessTierCool),
				},
			),
			ExpectError: regexp.MustCompile(`InvalidInput`),
		},
	}))
}

func TestUnit_OneLakeDefaultTierResource_UpdateError(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	fakeServer := fakes.NewFakeServer()
	state := &oneLakeDefaultTierState{}
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(state)
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.ModifyDefaultTier = fakeModifyDefaultTier(state)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"default_tier": string(fabcore.OneLakeAccessTierCool),
				},
			),
		},
		{
			PreConfig: func() {
				fakeServer.ServerFactory.Core.OneLakeSettingsServer.ModifyDefaultTier = fakeModifyDefaultTierError()
			},
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"default_tier": string(fabcore.OneLakeAccessTierCold),
				},
			),
			ExpectError: regexp.MustCompile(`InvalidInput`),
		},
	}))
}

func TestUnit_OneLakeDefaultTierResource_UpdateMissingModifyResponse(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	fakeServer := fakes.NewFakeServer()
	state := &oneLakeDefaultTierState{}
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(state)
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.ModifyDefaultTier = fakeModifyDefaultTier(state)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"default_tier": string(fabcore.OneLakeAccessTierCool),
				},
			),
		},
		{
			PreConfig: func() {
				fakeServer.ServerFactory.Core.OneLakeSettingsServer.ModifyDefaultTier = fakeModifyDefaultTierMissingResponse()
			},
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"default_tier": string(fabcore.OneLakeAccessTierCold),
				},
			),
			ExpectError: regexp.MustCompile(`Modify Default Tier API did not return the resulting tier`),
		},
	}))
}

func TestUnit_OneLakeDefaultTierResource_MissingModifyResponse(t *testing.T) {
	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.ModifyDefaultTier = fakeModifyDefaultTierMissingResponse()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"default_tier": string(fabcore.OneLakeAccessTierCool),
				},
			),
			ExpectError: regexp.MustCompile(`Modify Default Tier API did not return the resulting tier`),
		},
	}))
}

func TestUnit_OneLakeDefaultTierResource_ReadNotFoundRemovesState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	fakeServer := fakes.NewFakeServer()
	state := &oneLakeDefaultTierState{}
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(state)
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.ModifyDefaultTier = fakeModifyDefaultTier(state)

	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id": workspaceID,
			"default_tier": string(fabcore.OneLakeAccessTierCool),
		},
	)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config:       testCase,
		},
		{
			PreConfig: func() {
				fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettingsError()
			},
			ResourceName:       testResourceItemFQN,
			Config:             testCase,
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	}))
}

func TestUnit_OneLakeDefaultTierResource_ReadUnavailable(t *testing.T) {
	fakeServer := fakes.NewFakeServer()
	state := &oneLakeDefaultTierState{}
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettingsWithoutLifecycle()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.ModifyDefaultTier = fakeModifyDefaultTier(state)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"default_tier": string(fabcore.OneLakeAccessTierCool),
				},
			),
			ExpectError: regexp.MustCompile(`Get Settings did not return lifecycle.defaultTier`),
		},
	}))
}

func TestAcc_OneLakeDefaultTierResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceOneLakeTierRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"default_tier": string(fabcore.OneLakeAccessTierCool),
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "default_tier", string(fabcore.OneLakeAccessTierCool)),
			),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"default_tier": string(fabcore.OneLakeAccessTierCold),
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "default_tier", string(fabcore.OneLakeAccessTierCold)),
			),
		},
		{
			ResourceName:      testResourceItemFQN,
			ImportState:       true,
			ImportStateId:     workspaceID,
			ImportStateVerify: true,
		},
	}))
}
