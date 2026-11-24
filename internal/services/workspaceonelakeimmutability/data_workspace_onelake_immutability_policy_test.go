// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceonelakeimmutability_test

import (
	"net/http"
	"regexp"
	"strconv"
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

func TestUnit_WorkspaceOneLakeImmutabilityPolicyDataSource(t *testing.T) {
	store := &oneLakeImmutabilityPolicyStore{policy: ptrPolicy(30)}
	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(store)

	workspaceID := testhelp.RandomUUID()
	scope := string(fabcore.ImmutabilityScopeDiagnosticLogs)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			Config:      at.CompileConfig(testDataSourceItemHeader, map[string]any{}),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": "invalid-uuid",
					"scope":        scope,
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"scope":        "UnsupportedScope",
				},
			),
			ExpectError: regexp.MustCompile(`Attribute scope value must be one of`),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"scope":        scope,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "scope", scope),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "retention_days", "30"),
			),
		},
	}))
}

func TestUnit_WorkspaceOneLakeImmutabilityPolicyDataSource_NotFound(t *testing.T) {
	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(&oneLakeImmutabilityPolicyStore{})

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"scope":        string(fabcore.ImmutabilityScopeDiagnosticLogs),
				},
			),
			ExpectError: regexp.MustCompile(`No OneLake immutability policy`),
		},
	}))
}

func TestUnit_WorkspaceOneLakeImmutabilityPolicyDataSource_GetError(t *testing.T) {
	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(&oneLakeImmutabilityPolicyStore{
		getErrorStatus: http.StatusForbidden,
		getErrorCode:   "InsufficientPrivileges",
	})

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"scope":        string(fabcore.ImmutabilityScopeDiagnosticLogs),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}

func TestUnit_WorkspaceOneLakeImmutabilityPolicyDataSource_InvalidResponse(t *testing.T) {
	scope := fabcore.ImmutabilityScopeDiagnosticLogs
	invalidPolicy := fabcore.ImmutabilityPolicy{
		Scope: &scope,
	}

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(&oneLakeImmutabilityPolicyStore{policy: &invalidPolicy})

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"scope":        string(fabcore.ImmutabilityScopeDiagnosticLogs),
				},
			),
			ExpectError: regexp.MustCompile(`Invalid OneLake immutability policy response`),
		},
	}))
}

func TestUnit_WorkspaceOneLakeImmutabilityPolicyDataSource_InvalidRetention(t *testing.T) {
	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(&oneLakeImmutabilityPolicyStore{policy: ptrPolicy(0)})

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"scope":        string(fabcore.ImmutabilityScopeDiagnosticLogs),
				},
			),
			ExpectError: regexp.MustCompile(`Invalid OneLake immutability policy response`),
		},
	}))
}

func TestAcc_WorkspaceOneLakeImmutabilityPolicyDataSource(t *testing.T) {
	workspaceID, scope, retentionDays := oneLakeImmutabilityPolicyFromWellKnown(t)

	resource.Test(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"scope":        scope,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "scope", scope),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "retention_days", strconv.FormatInt(retentionDays, 10)),
			),
		},
	}))
}
