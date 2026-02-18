// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspacencp_test

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

func TestUnit_WorkspaceNetworkCommunicationPolicyResource_Attributes(t *testing.T) {
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
		// error - invalid outbound default_action
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"outbound": map[string]any{
						"public_access_rules": map[string]any{
							"default_action": "Invalid",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Attribute outbound.public_access_rules.default_action value must be one of`),
		},
		// error - invalid inbound default_action
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"inbound": map[string]any{
						"public_access_rules": map[string]any{
							"default_action": "Invalid",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Attribute inbound.public_access_rules.default_action value must be one of`),
		},
	}))
}

func TestUnit_WorkspaceNetworkCommunicationPolicyResource_CRUD(t *testing.T) {
	entity := NewRandomWorkspaceNetworkingCommunicationPolicy()
	workspaceID := testhelp.RandomUUID()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.SetNetworkCommunicationPolicy = fakeSetNetworkCommunicationPolicy(&entity)
	fakeServer.ServerFactory.Core.WorkspacesServer.GetNetworkCommunicationPolicy = fakeGetNetworkCommunicationPolicy(&entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		// create and read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"outbound": map[string]any{
						"public_access_rules": map[string]any{
							"default_action": "Deny",
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "outbound.public_access_rules.default_action", "Deny"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "inbound.public_access_rules.default_action", string(*entity.Inbound.PublicAccessRules.DefaultAction)),
			),
		},
		// update and read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"inbound": map[string]any{
						"public_access_rules": map[string]any{
							"default_action": "Deny",
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "outbound.public_access_rules.default_action", "Allow"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "inbound.public_access_rules.default_action", "Deny"),
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
				resource.TestCheckResourceAttr(testResourceItemFQN, "outbound.public_access_rules.default_action", "Allow"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "inbound.public_access_rules.default_action", "Allow"),
			),
		},
	}))
}

func TestAcc_WorkspaceSetNetworkCommunicationPolicy_CRUD(t *testing.T) {
	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, capacityID)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// create and read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"outbound": map[string]any{
							"public_access_rules": map[string]any{
								"default_action": "Deny",
							},
						},
						"inbound": map[string]any{
							"public_access_rules": map[string]any{
								"default_action": "Deny",
							},
						},
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "outbound.public_access_rules.default_action", "Deny"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "inbound.public_access_rules.default_action", "Deny"),
			),
		},
		// update and read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"inbound": map[string]any{
							"public_access_rules": map[string]any{
								"default_action": "Deny",
							},
						},
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "outbound.public_access_rules.default_action", "Allow"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "inbound.public_access_rules.default_action", "Deny"),
			),
		},
		// update and read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"inbound": map[string]any{
							"public_access_rules": map[string]any{
								"default_action": "Allow",
							},
						},
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "outbound.public_access_rules.default_action", "Allow"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "inbound.public_access_rules.default_action", "Allow"),
			),
		},
		// update and read - set to default values
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "outbound.public_access_rules.default_action", "Allow"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "inbound.public_access_rules.default_action", "Allow"),
			),
		},
	}))
}
