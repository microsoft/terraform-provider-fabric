// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceonelakeimmutability_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_WorkspaceOneLakeImmutabilityPolicyResource_Attributes(t *testing.T) {
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
					"workspace_id":   "invalid-uuid",
					"scope":          string(fabcore.ImmutabilityScopeDiagnosticLogs),
					"retention_days": 30,
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"scope":          "UnsupportedScope",
					"retention_days": 30,
				},
			),
			ExpectError: regexp.MustCompile(`Attribute scope value must be one of`),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"scope":          string(fabcore.ImmutabilityScopeDiagnosticLogs),
					"retention_days": 0,
				},
			),
			ExpectError: regexp.MustCompile(`retention_days`),
		},
	}))
}

func TestUnit_WorkspaceOneLakeImmutabilityPolicyResource_CRUD(t *testing.T) {
	store := &oneLakeImmutabilityPolicyStore{}
	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(store)
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.ModifyImmutabilityPolicy = fakeModifyOneLakeImmutabilityPolicy(store)

	workspaceID := testhelp.RandomUUID()
	scope := string(fabcore.ImmutabilityScopeDiagnosticLogs)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"scope":          scope,
					"retention_days": 7,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "scope", scope),
				resource.TestCheckResourceAttr(testResourceItemFQN, "retention_days", "7"),
			),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"scope":          scope,
					"retention_days": 9,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "scope", scope),
				resource.TestCheckResourceAttr(testResourceItemFQN, "retention_days", "9"),
			),
		},
	}))
}

func TestUnit_WorkspaceOneLakeImmutabilityPolicyResource_RetentionCannotDecrease(t *testing.T) {
	store := &oneLakeImmutabilityPolicyStore{}
	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(store)
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.ModifyImmutabilityPolicy = fakeModifyOneLakeImmutabilityPolicy(store)

	workspaceID := testhelp.RandomUUID()
	scope := string(fabcore.ImmutabilityScopeDiagnosticLogs)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"scope":          scope,
					"retention_days": 9,
				},
			),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"scope":          scope,
					"retention_days": 8,
				},
			),
			ExpectError: regexp.MustCompile(`Retention period cannot be reduced`),
		},
	}))
}

func TestUnit_WorkspaceOneLakeImmutabilityPolicyResource_RetentionCannotDecreaseAtApply(t *testing.T) {
	store := &oneLakeImmutabilityPolicyStore{}
	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(store)
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.ModifyImmutabilityPolicy = fakeModifyOneLakeImmutabilityPolicy(store)

	workspaceID := testhelp.RandomUUID()
	scope := string(fabcore.ImmutabilityScopeDiagnosticLogs)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"scope":          scope,
					"retention_days": 7,
				},
			),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"scope":          scope,
					"retention_days": 8,
				},
			),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					setPolicyBeforeApply{store: store, retentionDays: 30},
					plancheck.ExpectResourceAction(testResourceItemFQN, plancheck.ResourceActionUpdate),
				},
			},
			ExpectError: regexp.MustCompile(`Retention period cannot be reduced`),
		},
	}))
}

func TestUnit_WorkspaceOneLakeImmutabilityPolicyResource_ModifyError(t *testing.T) {
	store := &oneLakeImmutabilityPolicyStore{modifyError: true}
	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(store)
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.ModifyImmutabilityPolicy = fakeModifyOneLakeImmutabilityPolicy(store)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   testhelp.RandomUUID(),
					"scope":          string(fabcore.ImmutabilityScopeDiagnosticLogs),
					"retention_days": 7,
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorCreateHeader),
		},
	}))
}

func TestUnit_WorkspaceOneLakeImmutabilityPolicyResource_UpdateError(t *testing.T) {
	store := &oneLakeImmutabilityPolicyStore{modifyErrorOnCall: 2}
	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(store)
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.ModifyImmutabilityPolicy = fakeModifyOneLakeImmutabilityPolicy(store)

	workspaceID := testhelp.RandomUUID()
	scope := string(fabcore.ImmutabilityScopeDiagnosticLogs)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"scope":          scope,
					"retention_days": 7,
				},
			),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"scope":          scope,
					"retention_days": 8,
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorUpdateHeader),
		},
	}))
}

func TestUnit_WorkspaceOneLakeImmutabilityPolicyResource_ReadRemovesMissingPolicy(t *testing.T) {
	store := &oneLakeImmutabilityPolicyStore{}
	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(store)
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.ModifyImmutabilityPolicy = fakeModifyOneLakeImmutabilityPolicy(store)

	workspaceID := testhelp.RandomUUID()
	scope := string(fabcore.ImmutabilityScopeDiagnosticLogs)
	config := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id":   workspaceID,
			"scope":          scope,
			"retention_days": 7,
		},
	)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config:       config,
		},
		{
			PreConfig: func() {
				store.policy = nil
			},
			ResourceName:       testResourceItemFQN,
			Config:             config,
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	}))
}

func TestUnit_WorkspaceOneLakeImmutabilityPolicyResource_ReadRemovesMissingWorkspace(t *testing.T) {
	store := &oneLakeImmutabilityPolicyStore{}
	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(store)
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.ModifyImmutabilityPolicy = fakeModifyOneLakeImmutabilityPolicy(store)

	workspaceID := testhelp.RandomUUID()
	scope := string(fabcore.ImmutabilityScopeDiagnosticLogs)
	config := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id":   workspaceID,
			"scope":          scope,
			"retention_days": 7,
		},
	)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config:       config,
		},
		{
			PreConfig: func() {
				store.getErrorStatus = 404
				store.getErrorCode = "WorkspaceNotFound"
			},
			ResourceName:       testResourceItemFQN,
			Config:             config,
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	}))
}

func TestUnit_WorkspaceOneLakeImmutabilityPolicyResource_ImportStateInvalid(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	config := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id":   workspaceID,
			"scope":          string(fabcore.ImmutabilityScopeDiagnosticLogs),
			"retention_days": 30,
		},
	)
	fakeServer := fakes.NewFakeServer()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        config,
			ImportState:   true,
			ImportStateId: "invalid",
			ExpectError:   regexp.MustCompile(`Expected import ID`),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        config,
			ImportState:   true,
			ImportStateId: "invalid-uuid/DiagnosticLogs",
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        config,
			ImportState:   true,
			ImportStateId: workspaceID + "/UnsupportedScope",
			ExpectError:   regexp.MustCompile(`Unsupported OneLake immutability policy scope`),
		},
	}))
}

func TestUnit_WorkspaceOneLakeImmutabilityPolicyResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	scope := string(fabcore.ImmutabilityScopeDiagnosticLogs)
	store := &oneLakeImmutabilityPolicyStore{policy: ptrPolicy(30)}

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeSettingsServer.GetSettings = fakeGetOneLakeSettings(store)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:       testResourceItemFQN,
			Config:             at.CompileConfig(testResourceItemHeader, map[string]any{"workspace_id": workspaceID, "scope": scope, "retention_days": 30}),
			ImportState:        true,
			ImportStateId:      workspaceID + "/" + scope,
			ImportStatePersist: true,
			ImportStateCheck:   checkImportedPolicyState(workspaceID, scope, 30),
		},
	}))
}

func TestAcc_WorkspaceOneLakeImmutabilityPolicyResource_CRUD(t *testing.T) {
	workspaceID, scope, retentionDays := oneLakeImmutabilityPolicyFromWellKnown(t)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"scope":          scope,
					"retention_days": int(retentionDays),
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "scope", scope),
				resource.TestCheckResourceAttr(testResourceItemFQN, "retention_days", strconv.FormatInt(retentionDays, 10)),
			),
		},
	}))
}

func TestAcc_WorkspaceOneLakeImmutabilityPolicyResource_Import(t *testing.T) {
	workspaceID, scope, retentionDays := oneLakeImmutabilityPolicyFromWellKnown(t)
	config := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id":   workspaceID,
			"scope":          scope,
			"retention_days": int(retentionDays),
		},
	)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		{
			ResourceName:       testResourceItemFQN,
			Config:             config,
			ImportState:        true,
			ImportStateId:      workspaceID + "/" + scope,
			ImportStatePersist: true,
			ImportStateCheck:   checkImportedPolicyState(workspaceID, scope, retentionDays),
		},
	}))
}

func checkImportedPolicyState(workspaceID, scope string, retentionDays int64) func([]*terraform.InstanceState) error {
	return func(states []*terraform.InstanceState) error {
		if len(states) != 1 {
			return fmt.Errorf("expected one instance state, got %d", len(states))
		}

		attributes := states[0].Attributes
		if attributes["workspace_id"] != workspaceID {
			return fmt.Errorf("unexpected imported workspace_id %q", attributes["workspace_id"])
		}

		if attributes["scope"] != scope {
			return fmt.Errorf("unexpected imported scope %q", attributes["scope"])
		}

		if attributes["retention_days"] != strconv.FormatInt(retentionDays, 10) {
			return fmt.Errorf("unexpected imported retention_days %q", attributes["retention_days"])
		}

		return nil
	}
}

func ptrPolicy(retentionDays int32) *fabcore.ImmutabilityPolicy {
	policy := newOneLakeImmutabilityPolicy(retentionDays)

	return &policy
}
