// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_OneLakeDataAccessSecurityResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":    testhelp.RandomUUID(),
					"item_id":         testhelp.RandomUUID(),
					"role_name":       "role",
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - no required attribute - item_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(`The argument "item_id" is required, but no definition was found.`),
		},
		// error - no required attribute - workspace_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"item_id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// error - no required attribute - role_name
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": testhelp.RandomUUID(),
					"item_id":      testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(`The argument "role_name" is required, but no definition was found.`),
		},
	}))
}

func TestUnit_OneLakeDataAccessSecurityResource_ImportState(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()

	entity := fabcore.DataAccessRoleListItem{
		ID:   new(testhelp.RandomUUID()),
		Name: new("example"),
		Kind: to.Ptr(fabcore.DataAccessRoleKindPolicy),
		DecisionRules: []fabcore.DecisionRule{{
			Effect: to.Ptr(fabcore.EffectPermit),
			Permission: []fabcore.PermissionScope{
				{AttributeName: to.Ptr(fabcore.AttributeNamePath), AttributeValueIncludedIn: []string{"*"}},
				{AttributeName: to.Ptr(fabcore.AttributeNameAction), AttributeValueIncludedIn: []string{"Read"}},
			},
		}},
		Members: &fabcore.Members{
			FabricItemMembers: []fabcore.FabricItemMember{{
				ItemAccess: []fabcore.ItemAccess{fabcore.ItemAccessReadAll},
				SourcePath: new(workspaceID + "/" + itemID),
			}},
		},
	}

	UpsertIntoOneLakeDataAccessRoleStore(workspaceID, itemID, entity)

	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.GetDataAccessRole = fakeGetDataAccessRoleFunc()

	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"workspace_id": workspaceID,
			"item_id":      itemID,
			"role_name":    "example",
			"decision_rules": []map[string]any{
				{
					"effect": "Permit",
					"permission": []map[string]any{
						{"attribute_name": "Path", "attribute_value_included_in": []string{"*"}},
						{"attribute_name": "Action", "attribute_value_included_in": []string{"Read"}},
					},
				},
			},
			"members": map[string]any{
				"fabric_item_members": []map[string]any{
					{"item_access": []string{"ReadAll"}, "source_path": workspaceID + "/" + itemID},
				},
			},
		},
	)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(fmt.Sprintf(common.ErrorImportIdentifierDetails, "WorkspaceID/ItemID/RoleName")),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "test/id/name",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_OneLakeDataAccessSecurityResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()

	for key := range fakeOneLakeDataAccessRoleStore {
		delete(fakeOneLakeDataAccessRoleStore, key)
	}

	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.CreateOrUpdateSingleDataAccessRole = fakeCreateOrUpdateSingleDataAccessRoleFunc()
	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.GetDataAccessRole = fakeGetDataAccessRoleFunc()
	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.DeleteDataAccessRole = fakeDeleteDataAccessRoleFunc()

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"role_name":    "example",
					"decision_rules": []map[string]any{
						{
							"effect": "Permit",
							"permission": []map[string]any{
								{"attribute_name": "Path", "attribute_value_included_in": []string{"*"}},
								{"attribute_name": "Action", "attribute_value_included_in": []string{"Read"}},
							},
						},
					},
					"members": map[string]any{
						"fabric_item_members": []map[string]any{
							{"item_access": []string{"ReadAll"}, "source_path": workspaceID + "/" + itemID},
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "role_name", "example"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "item_id", itemID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "decision_rules.0.effect", "Permit"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"role_name":    "example",
					"decision_rules": []map[string]any{
						{
							"effect": "Permit",
							"permission": []map[string]any{
								{"attribute_name": "Path", "attribute_value_included_in": []string{"*"}},
								{"attribute_name": "Action", "attribute_value_included_in": []string{"ReadWrite"}},
							},
						},
					},
					"members": map[string]any{
						"fabric_item_members": []map[string]any{
							{"item_access": []string{"ReadAll", "Write"}, "source_path": workspaceID + "/" + itemID},
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "role_name", "example"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "item_id", itemID),
			),
		},
	}))
}

func TestAcc_OneLakeDataAccessSecurityResource_CRUD(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	lakehouse := testhelp.WellKnown()["LakehouseRS"].(map[string]any)
	itemID := lakehouse["id"].(string)

	roleName := testhelp.RandomName()

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"role_name":    roleName,
					"decision_rules": []map[string]any{
						{
							"effect": "Permit",
							"permission": []map[string]any{
								{"attribute_name": "Path", "attribute_value_included_in": []string{"*"}},
								{"attribute_name": "Action", "attribute_value_included_in": []string{"Read"}},
							},
						},
					},
					"members": map[string]any{
						"fabric_item_members": []map[string]any{
							{"item_access": []string{"ReadAll"}, "source_path": workspaceID + "/" + itemID},
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "role_name", roleName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "item_id", itemID),
			),
		},
		// Update and Read
		{
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id": workspaceID,
					"item_id":      itemID,
					"role_name":    roleName,
					"decision_rules": []map[string]any{
						{
							"effect": "Permit",
							"permission": []map[string]any{
								{"attribute_name": "Path", "attribute_value_included_in": []string{"*"}},
								{"attribute_name": "Action", "attribute_value_included_in": []string{"ReadWrite"}},
							},
						},
					},
					"members": map[string]any{
						"fabric_item_members": []map[string]any{
							{"item_access": []string{"ReadAll", "Write"}, "source_path": workspaceID + "/" + itemID},
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "role_name", roleName),
			),
		},
	}))
}
