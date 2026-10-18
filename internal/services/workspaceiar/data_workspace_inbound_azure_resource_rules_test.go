// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceiar_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_WorkspaceInboundAzureResourceRulesDataSource(t *testing.T) {
	entity := NewRandomWorkspaceInboundAzureResourceRules()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.GetInboundAzureResourceRules = fakeGetInboundAzureResourceRules(&entity)

	workspaceID := testhelp.RandomUUID()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - unexpected_attr
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
		// error - invalid workspace_id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": "invalid-uuid",
				},
			),
			ExpectError: regexp.MustCompile(`Invalid UUID String Value`),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "rules.#", "2"),
				resource.TestCheckTypeSetElemNestedAttrs(testDataSourceItemFQN, "rules.*", map[string]string{
					"display_name": *entity.Rules[0].DisplayName,
					"resource_id":  *entity.Rules[0].ResourceID,
				}),
				resource.TestCheckTypeSetElemNestedAttrs(testDataSourceItemFQN, "rules.*", map[string]string{
					"display_name": *entity.Rules[1].DisplayName,
					"resource_id":  *entity.Rules[1].ResourceID,
				}),
			),
		},
	}))
}

func TestUnit_WorkspaceInboundAzureResourceRulesDataSource_Error(t *testing.T) {
	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.GetInboundAzureResourceRules = fakeGetInboundAzureResourceRulesError()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakeServer.ServerFactory, nil, []resource.TestStep{
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

func TestAcc_WorkspaceInboundAzureResourceRulesDataSource(t *testing.T) {
	entity := testhelp.WellKnown()["WorkspaceIAP"].(map[string]any)
	entityID := entity["id"].(string)
	rule := entity["inboundAzureResourceRule"].(map[string]any)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read by id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": entityID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", entityID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "rules.#", "1"),
				resource.TestCheckTypeSetElemNestedAttrs(testDataSourceItemFQN, "rules.*", map[string]string{
					"display_name": rule["displayName"].(string),
					"resource_id":  rule["resourceId"].(string),
				}),
			),
		},
		// read by id - not found
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
