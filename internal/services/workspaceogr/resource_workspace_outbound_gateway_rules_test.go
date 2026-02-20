// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceogr_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_WorkspaceOutboundGatewayRulesResource_Attributes(t *testing.T) {
	workspaceID := testhelp.RandomUUID()

	fakeServer := fakes.NewFakeServer()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - invalid workspace_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "invalid-uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - invalid default_action
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"default_action": "Invalid",
				},
			),
			ExpectError: regexp.MustCompile(`Attribute default_action value must be one of`),
		},
		// error - invalid allowed_gateways id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"default_action": "Allow",
					"allowed_gateways": []map[string]any{
						{
							"id": "invalid-uuid",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_WorkspaceOutboundGatewayRulesResource_CRUD(t *testing.T) {
	entity := NewRandomWorkspaceOutboundGateways()
	gatewayID := testhelp.RandomUUID()
	workspaceID := testhelp.RandomUUID()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.SetOutboundGatewayRules = fakeSetOutboundGatewayRules(&entity)
	fakeServer.ServerFactory.Core.WorkspacesServer.GetOutboundGatewayRules = fakeGetOutboundGatewayRules(&entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		// create and read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"default_action": "Deny",
					"allowed_gateways": []map[string]any{
						{
							"id": gatewayID,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "default_action", "Deny"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "allowed_gateways.#", "1"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "allowed_gateways.0.id", gatewayID),
			),
		},
		// update and read - set to default values
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "default_action", "Allow"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "allowed_gateways.#", "0"),
			),
		},
	}))
}

func TestAcc_WorkspaceSetOutboundGatewayRules_CRUD(t *testing.T) {
	entity := testhelp.WellKnown()["WorkspaceOAP"].(map[string]any)
	entityID := entity["id"].(string)
	entityVirtualNetwork := testhelp.WellKnown()["GatewayVirtualNetwork"].(map[string]any)
	entityVirtualNetworkID := entityVirtualNetwork["id"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   entityID,
					"default_action": "Deny",
					"allowed_gateways": []map[string]any{
						{
							"id": entityVirtualNetworkID,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "default_action", "Deny"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "allowed_gateways.#", "1"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "allowed_gateways.0.id", entityVirtualNetworkID),
			),
		},
		// Update and Read - set allowed_gateways to default
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   entityID,
					"default_action": "Deny",
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "default_action", "Deny"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "allowed_gateways.#", "0"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   entityID,
					"default_action": "Allow",
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "default_action", "Allow"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "allowed_gateways.#", "0"),
			),
		},
		// Update and Read - set all rules to default
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": entityID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "default_action", "Allow"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "allowed_gateways.#", "0"),
			),
		},
	}))
}
