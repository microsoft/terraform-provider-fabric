// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceiar_test

import (
	"fmt"
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

var (
	testRuleFactory = map[string]string{
		"display_name": "data factory",
		"resource_id":  "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/testrg/providers/Microsoft.DataFactory/factories/testadf",
	}
	testRuleSQLServer = map[string]string{
		"display_name": "sql server",
		"resource_id":  "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/testrg/providers/Microsoft.Sql/servers/testsql",
	}
	testRuleDatabricks = map[string]string{
		"display_name": "databricks access connector",
		"resource_id":  "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/testrg/providers/Microsoft.Databricks/accessConnectors/testdbx",
	}
	// ARM emits lowercase segment keywords in some IDs.
	testRuleLowercaseSegments = map[string]string{
		"display_name": "lowercase segments",
		"resource_id":  "/subscriptions/11111111-1111-1111-1111-111111111111/resourcegroups/testrg/providers/Microsoft.Kusto/clusters/testadx",
	} // an odd number of segments after the provider namespace is accepted rather than rejected at plan time
	testRuleOddSegments = map[string]string{
		"display_name": "odd segments",
		"resource_id":  "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/testrg/providers/Microsoft.Sql/servers/testsql/databases",
	}
)

func testResourceConfig(workspaceID string, rules ...map[string]string) string {
	cfgRules := make([]map[string]any, 0, len(rules))

	for _, rule := range rules {
		cfgRules = append(cfgRules, map[string]any{
			"display_name": rule["display_name"],
			"resource_id":  rule["resource_id"],
		})
	}

	return at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id": workspaceID,
			"rules":        cfgRules,
		},
	)
}

func TestUnit_WorkspaceInboundAzureResourceRulesResource_Attributes(t *testing.T) {
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
		// error - missing required rules
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			ExpectError: regexp.MustCompile(`The argument "rules" is required, but no definition was found.`),
		},
		// error - unexpected_attr
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":    workspaceID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - invalid workspace_id
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig("invalid-uuid", testRuleFactory),
			ExpectError:  regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - rules - missing resource_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"rules": []map[string]any{
						{
							"display_name": "test",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Inappropriate value for attribute "rules"`),
		},
		// error - rules - resource_id is not an ARM resource ID
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"rules": []map[string]any{
						{
							"display_name": "test",
							"resource_id":  "not-a-resource-id",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`must be a valid Azure Resource Manager \(ARM\) resource ID`),
		},
		// error - rules - resource_id has a provider namespace but no resource type or name
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"rules": []map[string]any{
						{
							"display_name": "test",
							"resource_id":  "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/testrg/providers/Microsoft.DataFactory",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`must be a valid Azure Resource Manager \(ARM\) resource ID`),
		},
		// error - rules - resource_id has a resource type but no resource name
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"rules": []map[string]any{
						{
							"display_name": "test",
							"resource_id":  "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/testrg/providers/Microsoft.DataFactory/factories",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`must be a valid Azure Resource Manager \(ARM\) resource ID`),
		},
	}))
}

func TestUnit_WorkspaceInboundAzureResourceRulesResource_CRUD(t *testing.T) {
	entity := fabcore.WorkspaceInboundAzureResourceRules{}
	workspaceID := testhelp.RandomUUID()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.SetInboundAzureResourceRules = fakeSetInboundAzureResourceRules(&entity)
	fakeServer.ServerFactory.Core.WorkspacesServer.GetInboundAzureResourceRules = fakeGetInboundAzureResourceRules(&entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		// create and read
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleFactory),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "1"),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "rules.*", testRuleFactory),
			),
		},
		// update - add rules
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleFactory, testRuleSQLServer, testRuleDatabricks, testRuleLowercaseSegments, testRuleOddSegments),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "5"),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "rules.*", testRuleFactory),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "rules.*", testRuleSQLServer),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "rules.*", testRuleDatabricks),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "rules.*", testRuleLowercaseSegments),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "rules.*", testRuleOddSegments),
			),
		},
		// update - remove a rule
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleSQLServer),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "1"),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "rules.*", testRuleSQLServer),
			),
		},
		// update - remove all rules
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "0"),
			),
		},
	}))
}

func TestUnit_WorkspaceInboundAzureResourceRulesResource_RuleOrderIsIgnored(t *testing.T) {
	entity := fabcore.WorkspaceInboundAzureResourceRules{}
	workspaceID := testhelp.RandomUUID()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.SetInboundAzureResourceRules = fakeSetInboundAzureResourceRules(&entity)
	fakeServer.ServerFactory.Core.WorkspacesServer.GetInboundAzureResourceRules = fakeGetInboundAzureResourceRulesReversed(&entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleFactory, testRuleSQLServer),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "2"),
			),
		},
		// the API returns the rules in a different order, which must not produce a diff
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleFactory, testRuleSQLServer),
			PlanOnly:     true,
		},
	}))
}

func TestUnit_WorkspaceInboundAzureResourceRulesResource_CreateError(t *testing.T) {
	entity := fabcore.WorkspaceInboundAzureResourceRules{}
	workspaceID := testhelp.RandomUUID()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.SetInboundAzureResourceRules = fakeSetInboundAzureResourceRulesFailOnCall(&entity, 1)
	fakeServer.ServerFactory.Core.WorkspacesServer.GetInboundAzureResourceRules = fakeGetInboundAzureResourceRules(&entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleFactory),
			ExpectError:  regexp.MustCompile(common.ErrorCreateHeader),
		},
	}))
}

func TestUnit_WorkspaceInboundAzureResourceRulesResource_ReadError(t *testing.T) {
	entity := fabcore.WorkspaceInboundAzureResourceRules{}
	workspaceID := testhelp.RandomUUID()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.SetInboundAzureResourceRules = fakeSetInboundAzureResourceRules(&entity)
	fakeServer.ServerFactory.Core.WorkspacesServer.GetInboundAzureResourceRules = fakeGetInboundAzureResourceRulesError()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleFactory),
			ExpectError:  regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}

func TestUnit_WorkspaceInboundAzureResourceRulesResource_UpdateError(t *testing.T) {
	entity := fabcore.WorkspaceInboundAzureResourceRules{}
	workspaceID := testhelp.RandomUUID()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.SetInboundAzureResourceRules = fakeSetInboundAzureResourceRulesFailOnCall(&entity, 2)
	fakeServer.ServerFactory.Core.WorkspacesServer.GetInboundAzureResourceRules = fakeGetInboundAzureResourceRules(&entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleFactory),
		},
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleSQLServer),
			ExpectError:  regexp.MustCompile(common.ErrorUpdateHeader),
		},
	}))
}

func TestUnit_WorkspaceInboundAzureResourceRulesResource_DeleteError(t *testing.T) {
	entity := fabcore.WorkspaceInboundAzureResourceRules{}
	workspaceID := testhelp.RandomUUID()

	fakeServer := fakes.NewFakeServer()
	// 1: create succeeds, 2: destroy step fails, 3: test teardown destroy succeeds
	fakeServer.ServerFactory.Core.WorkspacesServer.SetInboundAzureResourceRules = fakeSetInboundAzureResourceRulesFailOnCall(&entity, 2)
	fakeServer.ServerFactory.Core.WorkspacesServer.GetInboundAzureResourceRules = fakeGetInboundAzureResourceRules(&entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleFactory),
			Check:        resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "1"),
		},
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleFactory),
			Destroy:      true,
			ExpectError:  regexp.MustCompile(errMsgSetInboundAzureResourceRules),
		},
	}))
}

func TestAcc_WorkspaceInboundAzureResourceRulesResource_CRUD(t *testing.T) {
	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	adf := testhelp.WellKnown()["AzureDataFactory"].(map[string]any)
	resourceID := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataFactory/factories/%s",
		adf["subscriptionId"].(string), adf["resourceGroupName"].(string), adf["name"].(string),
	)

	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, capacityID)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"rules": []map[string]any{
							{
								"display_name": "azure data factory",
								"resource_id":  resourceID,
							},
						},
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "1"),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "rules.*", map[string]string{
					"display_name": "azure data factory",
					"resource_id":  resourceID,
				}),
			),
		},
		// Update and Read - remove all rules
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"rules":        []map[string]any{},
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "0"),
			),
		},
	}))
}
