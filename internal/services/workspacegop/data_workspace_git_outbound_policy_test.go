// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspacegop_test

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

func TestUnit_WorkspaceGitOutboundPolicyDataSource(t *testing.T) {
	entity := NewRandomWorkspaceGitOutboundPolicy()

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.WorkspacesServer.GetGitOutboundPolicy = fakeGetGitOutboundPolicy(&entity)

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
			),
		},
	}))
}

func TestAcc_WorkspaceGitOutboundPolicyDataSource(t *testing.T) {
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
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "default_action"),
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
