// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakediagnosticsetting_test

import (
	"regexp"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func TestUnit_OneLakeDiagnosticSettingDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	disabledWorkspaceID := testhelp.RandomUUID()
	disabledWithDestinationWorkspaceID := testhelp.RandomUUID()
	unsupportedWorkspaceID := testhelp.RandomUUID()
	missingStatusWorkspaceID := testhelp.RandomUUID()
	missingDestinationWorkspaceID := testhelp.RandomUUID()
	unsupportedReferenceWorkspaceID := testhelp.RandomUUID()
	missingItemIDWorkspaceID := testhelp.RandomUUID()
	missingTypeWorkspaceID := testhelp.RandomUUID()
	omittedDiagnosticsWorkspaceID := testhelp.RandomUUID()
	destinationWorkspaceID := testhelp.RandomUUID()
	destinationItemID := testhelp.RandomUUID()
	disabledWithDestination := newEnabledOneLakeDiagnosticSetting(destinationWorkspaceID, destinationItemID)
	disabledWithDestination.Status = new("Disabled")

	serverFactory, settingsState := newFakeOneLakeSettingsServer(map[string]fabcore.OneLakeDiagnosticSettings{
		workspaceID: newEnabledOneLakeDiagnosticSetting(destinationWorkspaceID, destinationItemID),
		disabledWorkspaceID: {
			Status: new("Disabled"),
		},
		disabledWithDestinationWorkspaceID: disabledWithDestination,
		unsupportedWorkspaceID: {
			Status: new("Enabled"),
			Destination: &fabcore.OneLakeDiagnosticSettingsDestinationInfo{
				Type: new("FutureDestination"),
			},
		},
		missingStatusWorkspaceID: {},
		missingDestinationWorkspaceID: {
			Status: new("Enabled"),
		},
		unsupportedReferenceWorkspaceID: {
			Status: new("Enabled"),
			Destination: &fabcore.LakehouseOneLakeDiagnosticSettingsDestination{
				Type: new("Lakehouse"),
				Lakehouse: &fabcore.ItemReferenceByVariable{
					ReferenceType:     to.Ptr(fabcore.ItemReferenceTypeByVariable),
					VariableReference: new("$(/Library/Name)"),
				},
			},
		},
		missingItemIDWorkspaceID: {
			Status: new("Enabled"),
			Destination: &fabcore.LakehouseOneLakeDiagnosticSettingsDestination{
				Type: new("Lakehouse"),
				Lakehouse: &fabcore.ItemReferenceByID{
					ReferenceType: to.Ptr(fabcore.ItemReferenceTypeByID),
					WorkspaceID:   new(destinationWorkspaceID),
				},
			},
		},
		missingTypeWorkspaceID: {
			Status:      new("Enabled"),
			Destination: &fabcore.OneLakeDiagnosticSettingsDestinationInfo{},
		},
		omittedDiagnosticsWorkspaceID: {
			Status: new("Disabled"),
		},
	})
	settingsState.omitDiagnostics[omittedDiagnosticsWorkspaceID] = true

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, serverFactory, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": "invalid-uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id":    workspaceID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "status", "Enabled"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "destination.lakehouse.workspace_id", destinationWorkspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "destination.lakehouse.item_id", destinationItemID),
			),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": disabledWorkspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "status", "Disabled"),
				resource.TestCheckNoResourceAttr(testDataSourceItemFQN, "destination"),
			),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": disabledWithDestinationWorkspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "status", "Disabled"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "destination.lakehouse.workspace_id", destinationWorkspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "destination.lakehouse.item_id", destinationItemID),
			),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": unsupportedWorkspaceID,
				},
			),
			ExpectError: regexp.MustCompile(`Unsupported OneLake diagnostic destination type`),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": missingStatusWorkspaceID,
				},
			),
			ExpectError: regexp.MustCompile(`omitted the diagnostics status`),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": missingDestinationWorkspaceID,
				},
			),
			ExpectError: regexp.MustCompile(`omitted the enabled diagnostics destination`),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": unsupportedReferenceWorkspaceID,
				},
			),
			ExpectError: regexp.MustCompile(`Unsupported OneLake diagnostic Lakehouse reference type`),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": missingItemIDWorkspaceID,
				},
			),
			ExpectError: regexp.MustCompile(`omitted the destination Lakehouse workspace or`),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": missingTypeWorkspaceID,
				},
			),
			ExpectError: regexp.MustCompile(`Unsupported OneLake diagnostic destination type`),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": omittedDiagnosticsWorkspaceID,
				},
			),
			ExpectError: regexp.MustCompile(`omitted diagnostics`),
		},
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}

func TestUnit_OneLakeDiagnosticSettingDataSource_IgnoresUnrelatedLifecycleType(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	serverFactory, _ := newFakeOneLakeSettingsServer(map[string]fabcore.OneLakeDiagnosticSettings{
		workspaceID: {
			Status: new("Disabled"),
		},
	})

	testCase := testhelp.NewTestUnitCase(t, nil, serverFactory, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "status", "Disabled"),
			),
		},
	})

	transport := &oneLakeSettingsResponseTransport{
		transport: fabfake.NewServerFactoryTransport(serverFactory),
	}
	testCase.ProtoV6ProviderFactories = testhelp.GetTestUnitProtoV6ProviderFactories(
		&fabric.ClientOptions{Transport: transport},
		nil,
	)

	resource.ParallelTest(t, testCase)
}

func TestAcc_OneLakeDiagnosticSettingDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "status"),
			),
		},
	}))
}
