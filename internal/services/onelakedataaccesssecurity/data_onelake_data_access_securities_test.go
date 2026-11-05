// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity_test

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

var testDataSourceItemsFQN, testDataSourceItemsHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Types, "test")

func TestUnit_OneLakeDataAccessSecuritiesDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()

	entity1 := fakes.NewRandomOneLakeDataAccessRoleListItem()
	entity1.Members.MicrosoftEntraMembers = []fabcore.MicrosoftEntraMember{{
		ObjectID: new(testhelp.RandomUUID()),
		TenantID: new(testhelp.RandomUUID()),
	}}
	entity2 := fakes.NewRandomOneLakeDataAccessRoleListItem()
	entity3 := fakes.NewRandomOneLakeDataAccessRoleListItem()

	for key := range fakeOneLakeDataAccessRoleStore {
		delete(fakeOneLakeDataAccessRoleStore, key)
	}

	UpsertIntoOneLakeDataAccessRoleStore(workspaceID, itemID, entity1)
	UpsertIntoOneLakeDataAccessRoleStore(workspaceID, itemID, entity2)
	UpsertIntoOneLakeDataAccessRoleStore(workspaceID, itemID, entity3)

	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.ListDataAccessRoles = fakeListDataAccessRolesFunc()

	resource.Test(t, testhelp.NewTestUnitCase(t, &testDataSourceItemsFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - workspace_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id": "invalid uuid",
					"item_id":      itemID,
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - item_id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
				},
			) + fmt.Sprintf(`
				output "object_types_unset" {
					value = alltrue([
						for member in one([for role in %s.values : role if role.role_name == "%s"]).members.microsoft_entra_members :
						member.object_type == null
					])
				}
			`, testDataSourceItemsFQN, *entity1.Name),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "item_id", itemID),
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "values.#", "3"),
				resource.TestCheckTypeSetElemNestedAttrs(testDataSourceItemsFQN, "values.*", map[string]string{
					"role_name": *entity1.Name, "members.microsoft_entra_members.#": "1",
				}),
				resource.TestCheckOutput("object_types_unset", "true"),
			),
		},
	}))
}

func TestAcc_OneLakeDataAccessSecuritiesDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	lakehouse := testhelp.WellKnown()["LakehouseRS"].(map[string]any)
	itemID := lakehouse["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testDataSourceItemsFQN, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "item_id", itemID),
			),
		},
	}))
}
