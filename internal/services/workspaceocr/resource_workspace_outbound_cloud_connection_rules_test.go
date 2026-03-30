// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceocr_test

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

func TestUnit_WorkspaceOutboundCloudConnectionRulesResource_Attributes(t *testing.T) {
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
		// error - missing required attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "default_action" is required, but no definition was found.`),
		},
		// error - invalid workspace_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   "invalid-uuid",
					"default_action": string(fabcore.ConnectionAccessActionTypeAllow),
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
		// error - rule missing required connection_type
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"default_action": string(fabcore.ConnectionAccessActionTypeDeny),
					"rules": []map[string]any{
						{
							"default_action": string(fabcore.ConnectionAccessActionTypeDeny),
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`"connection_type" is required`),
		},
		// error - rule missing required default_action
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"default_action": string(fabcore.ConnectionAccessActionTypeDeny),
					"rules": []map[string]any{
						{
							"connection_type": "SQL",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`"default_action" is required`),
		},
		// error - invalid rule default_action
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"default_action": string(fabcore.ConnectionAccessActionTypeDeny),
					"rules": []map[string]any{
						{
							"connection_type": "SQL",
							"default_action":  "Invalid",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Attribute rules\[0\].default_action value must be one of`),
		},
		// error - invalid workspace_id in allowed_workspaces
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"default_action": string(fabcore.ConnectionAccessActionTypeDeny),
					"rules": []map[string]any{
						{
							"connection_type": "LakeHouse",
							"default_action":  string(fabcore.ConnectionAccessActionTypeDeny),
							"allowed_workspaces": []map[string]any{
								{
									"workspace_id": "invalid-uuid",
								},
							},
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - allowed_endpoints and allowed_workspaces are conflicting
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"default_action": string(fabcore.ConnectionAccessActionTypeDeny),
					"rules": []map[string]any{
						{
							"connection_type": "SQL",
							"default_action":  string(fabcore.ConnectionAccessActionTypeDeny),
							"allowed_endpoints": []map[string]any{
								{
									"hostname_pattern": "*.microsoft.com",
								},
							},
							"allowed_workspaces": []map[string]any{
								{
									"workspace_id": workspaceID,
								},
							},
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Attribute "rules\[0\].allowed_(endpoints|workspaces)" cannot be specified when`),
		},
	}))
}

func TestUnit_WorkspaceOutboundCloudConnectionRulesResource_CRUD(t *testing.T) {
	entity := NewRandomWorkspaceOutboundConnections()
	workspaceID := testhelp.RandomUUID()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.SetOutboundCloudConnectionRules = fakeSetOutboundCloudConnectionRules(&entity)
	fakeServer.ServerFactory.Core.WorkspacesServer.GetOutboundCloudConnectionRules = fakeGetOutboundCloudConnectionRules(&entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		// create and read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"default_action": string(fabcore.ConnectionAccessActionTypeDeny),
					"rules": []map[string]any{
						{
							"connection_type": "SQL",
							"default_action":  string(fabcore.ConnectionAccessActionTypeDeny),
							"allowed_endpoints": []map[string]any{
								{
									"hostname_pattern": "*.microsoft.com",
								},
							},
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "default_action", string(fabcore.ConnectionAccessActionTypeDeny)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "1"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.0.connection_type", "SQL"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.0.default_action", string(fabcore.ConnectionAccessActionTypeDeny)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.0.allowed_endpoints.#", "1"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.0.allowed_endpoints.0.hostname_pattern", "*.microsoft.com"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.0.allowed_workspaces.#", "0"),
			),
		},
		// update and read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"default_action": string(fabcore.ConnectionAccessActionTypeDeny),
					"rules": []map[string]any{
						{
							"connection_type": "SQL",
							"default_action":  string(fabcore.ConnectionAccessActionTypeDeny),
							"allowed_workspaces": []map[string]any{
								{
									"workspace_id": workspaceID,
								},
							},
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "default_action", string(fabcore.ConnectionAccessActionTypeDeny)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "1"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.0.connection_type", "SQL"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.0.default_action", string(fabcore.ConnectionAccessActionTypeDeny)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.0.allowed_endpoints.#", "0"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.0.allowed_workspaces.#", "1"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.0.allowed_workspaces.0.workspace_id", workspaceID),
			),
		},
		// update and read - set to default values
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   workspaceID,
					"default_action": string(fabcore.ConnectionAccessActionTypeAllow),
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "default_action", string(fabcore.ConnectionAccessActionTypeAllow)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "0"),
			),
		},
	}))
}

func TestAcc_WorkspaceSetOutboundCloudConnectionRules_CRUD(t *testing.T) {
	entity := testhelp.WellKnown()["WorkspaceOAP"].(map[string]any)
	workspaceDS := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceDSID := workspaceDS["id"].(string)
	entityID := entity["id"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   entityID,
					"default_action": "Deny",
					"rules": []map[string]any{
						{
							"connection_type": "SQL",
							"default_action":  string(fabcore.ConnectionAccessActionTypeDeny),
							"allowed_endpoints": []map[string]any{
								{
									"hostname_pattern": "*.microsoft.com",
								},
							},
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "default_action", string(fabcore.ConnectionAccessActionTypeDeny)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "1"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.0.connection_type", "SQL"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.0.default_action", string(fabcore.ConnectionAccessActionTypeDeny)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.0.allowed_endpoints.#", "1"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.0.allowed_endpoints.0.hostname_pattern", "*.microsoft.com"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   entityID,
					"default_action": string(fabcore.ConnectionAccessActionTypeDeny),
					"rules": []map[string]any{
						{
							"connection_type": "lakehouse",
							"default_action":  string(fabcore.ConnectionAccessActionTypeDeny),
							"allowed_workspaces": []map[string]any{
								{
									"workspace_id": workspaceDSID,
								},
							},
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "default_action", string(fabcore.ConnectionAccessActionTypeDeny)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.0.allowed_workspaces.#", "1"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.0.allowed_workspaces.0.workspace_id", workspaceDSID),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   entityID,
					"default_action": string(fabcore.ConnectionAccessActionTypeAllow),
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "default_action", string(fabcore.ConnectionAccessActionTypeAllow)),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "0"),
			),
		},
	}))
}
