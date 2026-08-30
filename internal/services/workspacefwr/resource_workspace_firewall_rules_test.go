// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspacefwr_test

import (
	"regexp"
	"strings"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

var (
	testRuleCIDR   = map[string]string{"display_name": "cidr", "value": "203.0.113.0/24"}
	testRuleRange  = map[string]string{"display_name": "range", "value": "198.51.100.10-198.51.100.20"}
	testRuleSingle = map[string]string{"display_name": "single", "value": "192.0.2.42"}
)

func testResourceConfig(workspaceID string, rules ...map[string]string) string {
	cfgRules := make([]map[string]any, 0, len(rules))

	for _, rule := range rules {
		cfgRules = append(cfgRules, map[string]any{
			"display_name": rule["display_name"],
			"value":        rule["value"],
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

func TestUnit_WorkspaceFirewallRulesResource_Attributes(t *testing.T) {
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
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": "invalid-uuid",
					"rules": []map[string]any{
						{
							"display_name": "test",
							"value":        "203.0.113.0/24",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - rules - missing value
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
		// error - rules - display_name too long
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"rules": []map[string]any{
						{
							"display_name": strings.Repeat("a", 129),
							"value":        "203.0.113.0/24",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`string length must be between 1 and\s+128`),
		},
		// error - rules - invalid value
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"rules": []map[string]any{
						{
							"display_name": "test",
							"value":        "not-an-ip",
						},
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute`),
		},
	}))
}

func TestUnit_WorkspaceFirewallRulesResource_CRUD(t *testing.T) {
	entity := NewRandomInboundFirewallConfiguration()
	workspaceID := testhelp.RandomUUID()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.SetFirewallRules = fakeSetFirewallRules(&entity)
	fakeServer.ServerFactory.Core.WorkspacesServer.GetFirewallRules = fakeGetFirewallRules(&entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		// create and read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"rules": []map[string]any{
						{
							"display_name": "cidr",
							"value":        "203.0.113.0/24",
						},
						{
							"display_name": "range",
							"value":        "198.51.100.10-198.51.100.20",
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "2"),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "rules.*", map[string]string{"display_name": "cidr", "value": "203.0.113.0/24"}),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "rules.*", map[string]string{"display_name": "range", "value": "198.51.100.10-198.51.100.20"}),
			),
		},
		// update and read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"rules": []map[string]any{
						{
							"display_name": "single",
							"value":        "192.0.2.42",
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "1"),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "rules.*", map[string]string{"display_name": "single", "value": "192.0.2.42"}),
			),
		},
		// update and read - clear all rules
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"rules":        []map[string]any{},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "0"),
			),
		},
	}))
}

func TestUnit_WorkspaceFirewallRulesResource_CreateError(t *testing.T) {
	entity := fabcore.InboundFirewallConfiguration{}
	workspaceID := testhelp.RandomUUID()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.SetFirewallRules = fakeSetFirewallRulesFailOnCall(&entity, 1)
	fakeServer.ServerFactory.Core.WorkspacesServer.GetFirewallRules = fakeGetFirewallRules(&entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleCIDR),
			ExpectError:  regexp.MustCompile(errMsgSetFirewallRules),
		},
	}))
}

func TestUnit_WorkspaceFirewallRulesResource_UpdateError(t *testing.T) {
	entity := fabcore.InboundFirewallConfiguration{}
	workspaceID := testhelp.RandomUUID()

	fakeServer := fakes.NewFakeServer()
	// 1: create succeeds, 2: update fails, 3: destroy succeeds
	fakeServer.ServerFactory.Core.WorkspacesServer.SetFirewallRules = fakeSetFirewallRulesFailOnCall(&entity, 2)
	fakeServer.ServerFactory.Core.WorkspacesServer.GetFirewallRules = fakeGetFirewallRules(&entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleCIDR),
			Check:        resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "1"),
		},
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleSingle),
			ExpectError:  regexp.MustCompile(errMsgSetFirewallRules),
		},
	}))
}

func TestUnit_WorkspaceFirewallRulesResource_DeleteError(t *testing.T) {
	entity := fabcore.InboundFirewallConfiguration{}
	workspaceID := testhelp.RandomUUID()

	fakeServer := fakes.NewFakeServer()
	// 1: create succeeds, 2: destroy step fails, 3: test teardown destroy succeeds
	fakeServer.ServerFactory.Core.WorkspacesServer.SetFirewallRules = fakeSetFirewallRulesFailOnCall(&entity, 2)
	fakeServer.ServerFactory.Core.WorkspacesServer.GetFirewallRules = fakeGetFirewallRules(&entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleCIDR),
			Check:        resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "1"),
		},
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleCIDR),
			Destroy:      true,
			ExpectError:  regexp.MustCompile(errMsgSetFirewallRules),
		},
	}))
}

// The API does not guarantee rule ordering, so a reordered response must not produce a diff.
func TestUnit_WorkspaceFirewallRulesResource_RuleOrderIsIgnored(t *testing.T) {
	entity := fabcore.InboundFirewallConfiguration{}
	workspaceID := testhelp.RandomUUID()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.SetFirewallRules = fakeSetFirewallRules(&entity)
	fakeServer.ServerFactory.Core.WorkspacesServer.GetFirewallRules = fakeGetFirewallRulesReversed(&entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		// create - the response order is reversed, state must still match the configuration
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleCIDR, testRuleRange, testRuleSingle),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "rules.#", "3"),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "rules.*", testRuleCIDR),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "rules.*", testRuleRange),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "rules.*", testRuleSingle),
			),
		},
		// refresh with the same configuration must be a no-op
		{
			ResourceName: testResourceItemFQN,
			Config:       testResourceConfig(workspaceID, testRuleCIDR, testRuleRange, testRuleSingle),
			PlanOnly:     true,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PostApplyPostRefresh: []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				},
			},
		},
	}))
}
