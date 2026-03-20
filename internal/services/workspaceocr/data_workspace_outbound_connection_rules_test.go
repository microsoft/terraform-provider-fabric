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
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_WorkspaceOutboundCloudConnectionRulesDataSource(t *testing.T) {
	entity := NewRandomWorkspaceOutboundConnections()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.GetOutboundCloudConnectionRules = fakeGetOutboundCloudConnectionRules(&entity)

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
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "default_action", string(*entity.DefaultAction)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "rules.#", "2"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "rules.0.connection_type", "SQL"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "rules.0.default_action", string(fabcore.ConnectionAccessActionTypeDeny)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "rules.0.allowed_endpoints.#", "1"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "rules.0.allowed_endpoints.0.hostname_pattern", *entity.Rules[0].AllowedEndpoints[0].HostnamePattern),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "rules.0.allowed_workspaces.#", "0"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "rules.1.connection_type", "LAKEHOUSE"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "rules.1.default_action", string(fabcore.ConnectionAccessActionTypeDeny)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "rules.1.allowed_endpoints.#", "0"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "rules.1.allowed_workspaces.#", "1"),
			),
		},
	}))
}

func TestAcc_WorkspaceOutboundCloudConnectionRulesDataSource(t *testing.T) {
	entity := testhelp.WellKnown()["WorkspaceOAP"].(map[string]any)
	entityID := entity["id"].(string)

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
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "default_action", "Allow"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "rules.#", "0"),
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
