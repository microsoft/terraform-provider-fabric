// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity_test

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

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_OneLakeDataAccessSecurityDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()
	entity := fakes.NewRandomOneLakeDataAccessRoleListItem()

	for key := range fakeOneLakeDataAccessRoleStore {
		delete(fakeOneLakeDataAccessRoleStore, key)
	}

	UpsertIntoOneLakeDataAccessRoleStore(workspaceID, itemID, entity)

	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.GetDataAccessRole = fakeGetDataAccessRoleFunc()
	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.ListDataAccessRoles = fakeListDataAccessRolesFunc()

	resource.Test(t, testhelp.NewTestUnitCase(t, &testDataSourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - workspace_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
					"item_id":      itemID,
					"role_name":    *entity.Name,
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - item_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      "invalid uuid",
					"role_name":    *entity.Name,
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// read by role_name
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"role_name":    *entity.Name,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "item_id", itemID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "role_name", *entity.Name),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "decision_rules.0.effect", string(*entity.DecisionRules[0].Effect)),
			),
		},
	}))
}

func TestAcc_OneLakeDataAccessSecurityDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	lakehouse := testhelp.WellKnown()["LakehouseRS"].(map[string]any)
	itemID := lakehouse["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testDataSourceItemFQN, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"role_name":    "DefaultReader",
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "item_id", itemID),
			),
		},
	}))
}
