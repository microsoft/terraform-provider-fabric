// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakediagnosticsetting_test

import (
	"errors"
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func TestUnit_OneLakeDiagnosticSettingResource_Attributes(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	destinationWorkspaceID := testhelp.RandomUUID()
	serverFactory, _ := newFakeOneLakeSettingsServer(map[string]fabcore.OneLakeDiagnosticSettings{})

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, serverFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"status":       "Invalid",
				},
			),
			ExpectError: regexp.MustCompile(`status value must be one of`),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"status":       "Enabled",
				},
			),
			ExpectError: regexp.MustCompile(`destination`),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "invalid-uuid",
					"status":       "Disabled",
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
					"status":       "Enabled",
					"destination": map[string]any{
						"lakehouse": map[string]any{
							"workspace_id": destinationWorkspaceID,
							"item_id":      "invalid-uuid",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_OneLakeDiagnosticSettingResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	destinationWorkspaceID := testhelp.RandomUUID()
	destinationItemID := testhelp.RandomUUID()
	serverFactory, state := newFakeOneLakeSettingsServer(map[string]fabcore.OneLakeDiagnosticSettings{
		workspaceID: {
			Status: new("Disabled"),
		},
	})

	testCase := testhelp.NewTestUnitCase(t, &testResourceItemFQN, serverFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"status":       "Enabled",
					"destination": map[string]any{
						"lakehouse": map[string]any{
							"workspace_id": destinationWorkspaceID,
							"item_id":      destinationItemID,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "status", "Enabled"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "destination.lakehouse.workspace_id", destinationWorkspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "destination.lakehouse.item_id", destinationItemID),
			),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"status":       "Disabled",
					"destination": map[string]any{
						"lakehouse": map[string]any{
							"workspace_id": destinationWorkspaceID,
							"item_id":      destinationItemID,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "status", "Disabled"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "destination.lakehouse.workspace_id", destinationWorkspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "destination.lakehouse.item_id", destinationItemID),
			),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"status":       "Disabled",
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "status", "Disabled"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "destination"),
			),
		},
		{
			ResourceName:                         testResourceItemFQN,
			ImportState:                          true,
			ImportStateId:                        workspaceID,
			ImportStateVerify:                    true,
			ImportStateVerifyIdentifierAttribute: "workspace_id",
		},
	})

	testCase.CheckDestroy = func(_ *terraform.State) error {
		settings := state.settings[workspaceID]
		if settings.Status == nil || *settings.Status != "Disabled" {
			return errors.New("expected diagnostics to be disabled after destroy")
		}

		destination, ok := settings.Destination.(*fabcore.LakehouseOneLakeDiagnosticSettingsDestination)
		if !ok || destination.Lakehouse != nil {
			return errors.New("expected destroy to clear the Lakehouse destination")
		}

		if state.noPreviousSettingsErrorCount != 1 {
			return errors.New("expected delete to treat NoPreviousSettingsExist as an idempotent reset")
		}

		return nil
	}

	resource.ParallelTest(t, testCase)
}

func TestUnit_OneLakeDiagnosticSettingResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	serverFactory, _ := newFakeOneLakeSettingsServer(map[string]fabcore.OneLakeDiagnosticSettings{
		workspaceID: {
			Status: new("Disabled"),
		},
	})

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, serverFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"status":       "Disabled",
				},
			),
			ImportState:   true,
			ImportStateId: "invalid-workspace-id",
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"status":       "Disabled",
				},
			),
			ImportState:   true,
			ImportStateId: workspaceID,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "status", "Disabled"),
			),
		},
	}))
}

func TestUnit_OneLakeDiagnosticSettingResource_ReadNotFound(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	serverFactory, state := newFakeOneLakeSettingsServer(map[string]fabcore.OneLakeDiagnosticSettings{
		workspaceID: {
			Status: new("Disabled"),
		},
	})

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, serverFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"status":       "Disabled",
				},
			),
		},
		{
			ResourceName: testResourceItemFQN,
			PreConfig: func() {
				delete(state.settings, workspaceID)
			},
			RefreshState:       true,
			ExpectNonEmptyPlan: true,
			Check: func(state *terraform.State) error {
				if _, exists := state.RootModule().Resources[testResourceItemFQN]; exists {
					return errors.New("expected missing workspace diagnostics to remove the resource from state")
				}

				return nil
			},
		},
	}))
}

func TestAcc_OneLakeDiagnosticSettingResource_CRUD(t *testing.T) {
	sourceWorkspace := testhelp.WellKnown()["WorkspaceOD"].(map[string]any)
	destinationWorkspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	lakehouse := testhelp.WellKnown()["Lakehouse"].(map[string]any)

	sourceWorkspaceID := sourceWorkspace["id"].(string)
	destinationWorkspaceID := destinationWorkspace["id"].(string)
	destinationItemID := lakehouse["id"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": sourceWorkspaceID,
					"status":       "Enabled",
					"destination": map[string]any{
						"lakehouse": map[string]any{
							"workspace_id": destinationWorkspaceID,
							"item_id":      destinationItemID,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "status", "Enabled"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "destination.lakehouse.workspace_id", destinationWorkspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "destination.lakehouse.item_id", destinationItemID),
			),
		},
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": sourceWorkspaceID,
					"status":       "Disabled",
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "status", "Disabled"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "destination"),
			),
		},
	}))
}

func TestUnit_OneLakeDiagnosticSettingResource_CreateError(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	serverFactory, _ := newFakeOneLakeSettingsServer(map[string]fabcore.OneLakeDiagnosticSettings{})

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, serverFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"status":       "Disabled",
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorCreateHeader),
		},
	}))
}
