// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity_test

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func checkEntraMemberObjectType(resourceName, objectID, want string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		instance, ok := state.RootModule().Resources[resourceName]
		if !ok || instance.Primary == nil {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}
		for key, value := range instance.Primary.Attributes {
			if strings.Contains(key, "members.microsoft_entra_members.") && strings.HasSuffix(key, ".object_id") && value == objectID {
				typeKey := strings.TrimSuffix(key, ".object_id") + ".object_type"
				if got := instance.Primary.Attributes[typeKey]; got != want {
					return fmt.Errorf("%s: got %q, want %q", typeKey, got, want)
				}
				return nil
			}
		}
		return fmt.Errorf("member %s not found in %s", objectID, resourceName)
	}
}

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
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s/%s", workspaceID, itemID, *entity.Name),
			ImportState:   true,
			ImportStateCheck: func(is []*terraform.InstanceState) error {
				if len(is) != 1 {
					return errors.New("expected one instance state")
				}

				if got := is[0].Attributes["workspace_id"]; got != workspaceID {
					return fmt.Errorf("%s: unexpected workspace_id — got %q, want %q", testResourceItemFQN, got, workspaceID)
				}

				if got := is[0].Attributes["item_id"]; got != itemID {
					return fmt.Errorf("%s: unexpected item_id — got %q, want %q", testResourceItemFQN, got, itemID)
				}

				if got := is[0].Attributes["role_name"]; got != *entity.Name {
					return fmt.Errorf("%s: unexpected role_name — got %q, want %q", testResourceItemFQN, got, *entity.Name)
				}

				return nil
			},
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

func TestUnit_OneLakeDataAccessSecurityResource_DeprecatedObjectType(t *testing.T) {
	workspaceID, itemID, tenantID := testhelp.RandomUUID(), testhelp.RandomUUID(), testhelp.RandomUUID()
	objectID1, objectID2 := testhelp.RandomUUID(), testhelp.RandomUUID()
	config := func(members []map[string]any) string {
		return at.CompileConfig(testResourceItemHeader, map[string]any{
			"workspace_id": workspaceID,
			"item_id":      itemID,
			"role_name":    "deprecated-object-type",
			"decision_rules": []map[string]any{{
				"effect": "Permit",
				"permission": []map[string]any{
					{"attribute_name": "Path", "attribute_value_included_in": []string{"*"}},
					{"attribute_name": "Action", "attribute_value_included_in": []string{"Read"}},
				},
			}},
			"members": map[string]any{"microsoft_entra_members": members},
		})
	}
	legacy := map[string]any{"object_id": objectID1, "tenant_id": tenantID, "object_type": "ManagedIdentity"}
	second := map[string]any{"object_id": objectID2, "tenant_id": tenantID}
	secondLegacy := map[string]any{"object_id": objectID2, "tenant_id": tenantID, "object_type": "User"}
	noType := map[string]any{"object_id": objectID1, "tenant_id": tenantID}

	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.CreateOrUpdateSingleDataAccessRole = fakeCreateOrUpdateSingleDataAccessRoleFunc()
	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.GetDataAccessRole = fakeGetDataAccessRoleFunc()
	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.DeleteDataAccessRole = fakeDeleteDataAccessRoleFunc()

	checkRequestOmitsObjectType := func(_ *terraform.State) error {
		id := GenerateOneLakeDataAccessRoleID(workspaceID, itemID, "deprecated-object-type")
		role, ok := fakeOneLakeDataAccessRoleStore[id]
		if !ok || role.Members == nil {
			return errors.New("created role or its members are missing")
		}
		for _, member := range role.Members.MicrosoftEntraMembers {
			if member.ObjectType != nil {
				return errors.New("object_type was sent to Fabric")
			}
		}
		return nil
	}

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			Config: config([]map[string]any{legacy}),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "members.microsoft_entra_members.*", map[string]string{
					"object_id": objectID1, "tenant_id": tenantID, "object_type": "ManagedIdentity",
				}),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "members.fabric_item_members.#"),
				checkRequestOmitsObjectType,
			),
		},
		{
			Config: config([]map[string]any{legacy, secondLegacy}),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "members.microsoft_entra_members.#", "2"),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "members.microsoft_entra_members.*", map[string]string{
					"object_id": objectID1, "tenant_id": tenantID, "object_type": "ManagedIdentity",
				}),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "members.microsoft_entra_members.*", map[string]string{
					"object_id": objectID2, "tenant_id": tenantID, "object_type": "User",
				}),
				checkRequestOmitsObjectType,
			),
		},
		{
			Config: config([]map[string]any{legacy, second}),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkEntraMemberObjectType(testResourceItemFQN, objectID1, "ManagedIdentity"),
				checkEntraMemberObjectType(testResourceItemFQN, objectID2, ""),
			),
		},
		{
			Config: config([]map[string]any{noType}),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "members.microsoft_entra_members.#", "1"),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "members.microsoft_entra_members.*", map[string]string{
					"object_id": objectID1, "tenant_id": tenantID,
				}),
				checkEntraMemberObjectType(testResourceItemFQN, objectID1, ""),
			),
		},
		{
			ResourceName:                         testResourceItemFQN,
			Config:                               config([]map[string]any{noType}),
			ImportState:                          true,
			ImportStateVerify:                    true,
			ImportStateId:                        fmt.Sprintf("%s/%s/%s", workspaceID, itemID, "deprecated-object-type"),
			ImportStateVerifyIdentifierAttribute: "role_name",
			ImportStateVerifyIgnore: []string{
				"members.fabric_item_members",
			},
		},
		{
			PreConfig: func() {
				id := GenerateOneLakeDataAccessRoleID(workspaceID, itemID, "deprecated-object-type")
				role := fakeOneLakeDataAccessRoleStore[id]
				role.Members.MicrosoftEntraMembers = []fabcore.MicrosoftEntraMember{{ObjectID: &objectID2, TenantID: &tenantID}}
				fakeOneLakeDataAccessRoleStore[id] = role
			},
			RefreshState:       true,
			ExpectNonEmptyPlan: true,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "members.microsoft_entra_members.#", "1"),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "members.microsoft_entra_members.*", map[string]string{
					"object_id": objectID2, "tenant_id": tenantID,
				}),
				checkEntraMemberObjectType(testResourceItemFQN, objectID2, ""),
			),
		},
	}))
}

func TestUnit_OneLakeDataAccessSecurityResource_ObjectTypeUUIDCasing(t *testing.T) {
	workspaceID, itemID := testhelp.RandomUUID(), testhelp.RandomUUID()
	objectID, tenantID := strings.ToUpper(testhelp.RandomUUID()), strings.ToUpper(testhelp.RandomUUID())
	config := oneLakeMembersConfig(workspaceID, itemID, "uuid-casing", map[string]any{
		"microsoft_entra_members": []map[string]any{{"object_id": objectID, "tenant_id": tenantID, "object_type": "User"}},
	})

	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.CreateOrUpdateSingleDataAccessRole = fakeCreateOrUpdateSingleDataAccessRoleFunc()
	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.GetDataAccessRole = fakeGetDataAccessRoleFuncWithMemberTransform(func(members *fabcore.Members) {
		for i := range members.MicrosoftEntraMembers {
			member := &members.MicrosoftEntraMembers[i]
			member.ObjectID = new(strings.ToLower(*member.ObjectID))
			member.TenantID = new(strings.ToLower(*member.TenantID))
		}
	})
	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.DeleteDataAccessRole = fakeDeleteDataAccessRoleFunc()

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			Config: config,
			Check:  resource.TestCheckResourceAttr(testResourceItemFQN, "members.microsoft_entra_members.0.object_type", "User"),
		},
		{Config: config, PlanOnly: true},
	}))
}

func TestUnit_OneLakeDataAccessSecurityResource_ObjectTypeSameObjectDifferentTenant(t *testing.T) {
	workspaceID, itemID := testhelp.RandomUUID(), testhelp.RandomUUID()
	objectID, tenantID1, tenantID2 := testhelp.RandomUUID(), testhelp.RandomUUID(), testhelp.RandomUUID()
	config := oneLakeMembersConfig(workspaceID, itemID, "different-tenants", map[string]any{
		"microsoft_entra_members": []map[string]any{
			{"object_id": objectID, "tenant_id": tenantID1, "object_type": "User"},
			{"object_id": objectID, "tenant_id": tenantID2, "object_type": "Group"},
		},
	})

	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.CreateOrUpdateSingleDataAccessRole = fakeCreateOrUpdateSingleDataAccessRoleFunc()
	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.GetDataAccessRole = fakeGetDataAccessRoleFunc()
	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.DeleteDataAccessRole = fakeDeleteDataAccessRoleFunc()

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			Config: config,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "members.microsoft_entra_members.#", "2"),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "members.microsoft_entra_members.*", map[string]string{
					"object_id": objectID, "tenant_id": tenantID1, "object_type": "User",
				}),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "members.microsoft_entra_members.*", map[string]string{
					"object_id": objectID, "tenant_id": tenantID2, "object_type": "Group",
				}),
			),
		},
		{Config: config, PlanOnly: true},
	}))
}

func TestUnit_OneLakeDataAccessSecurityResource_EmptyMemberCollections(t *testing.T) {
	workspaceID, itemID, objectID, tenantID := testhelp.RandomUUID(), testhelp.RandomUUID(), testhelp.RandomUUID(), testhelp.RandomUUID()
	entra := []map[string]any{{"object_id": objectID, "tenant_id": tenantID}}
	fabric := []map[string]any{{"item_access": []string{"ReadAll"}, "source_path": workspaceID + "/" + itemID}}
	responseEmptyAsNull := false

	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.CreateOrUpdateSingleDataAccessRole = fakeCreateOrUpdateSingleDataAccessRoleFunc()
	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.GetDataAccessRole = fakeGetDataAccessRoleFuncWithMemberTransform(func(members *fabcore.Members) {
		if len(members.FabricItemMembers) == 0 {
			members.FabricItemMembers = nil
			if !responseEmptyAsNull {
				members.FabricItemMembers = []fabcore.FabricItemMember{}
			}
		}
		if len(members.MicrosoftEntraMembers) == 0 {
			members.MicrosoftEntraMembers = nil
			if !responseEmptyAsNull {
				members.MicrosoftEntraMembers = []fabcore.MicrosoftEntraMember{}
			}
		}
	})
	fakes.FakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.DeleteDataAccessRole = fakeDeleteDataAccessRoleFunc()

	checkRequest := func(fabricNull, entraNull bool) resource.TestCheckFunc {
		return func(_ *terraform.State) error {
			role, ok := fakeOneLakeDataAccessRoleStore[GenerateOneLakeDataAccessRoleID(workspaceID, itemID, "empty-members")]
			if !ok || role.Members == nil {
				return errors.New("created role or its members are missing")
			}
			if (role.Members.FabricItemMembers == nil) != fabricNull || (role.Members.MicrosoftEntraMembers == nil) != entraNull {
				return fmt.Errorf("unexpected request collections: fabric_item_members nil=%t, microsoft_entra_members nil=%t",
					role.Members.FabricItemMembers == nil, role.Members.MicrosoftEntraMembers == nil)
			}
			return nil
		}
	}

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			Config: oneLakeMembersConfig(workspaceID, itemID, "empty-members", map[string]any{"microsoft_entra_members": entra}),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkRequest(true, false),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "members.fabric_item_members.#"),
			),
		},
		{
			PreConfig: func() { responseEmptyAsNull = true },
			Config: oneLakeMembersConfig(workspaceID, itemID, "empty-members", map[string]any{
				"fabric_item_members": []map[string]any{}, "microsoft_entra_members": entra,
			}),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkRequest(false, false),
				resource.TestCheckResourceAttr(testResourceItemFQN, "members.fabric_item_members.#", "0"),
			),
		},
		{
			PreConfig: func() { responseEmptyAsNull = false },
			Config:    oneLakeMembersConfig(workspaceID, itemID, "empty-members", map[string]any{"fabric_item_members": fabric}),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkRequest(false, true),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "members.microsoft_entra_members.#"),
			),
		},
		{
			PreConfig: func() { responseEmptyAsNull = true },
			Config: oneLakeMembersConfig(workspaceID, itemID, "empty-members", map[string]any{
				"fabric_item_members": fabric, "microsoft_entra_members": []map[string]any{},
			}),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkRequest(false, false),
				resource.TestCheckResourceAttr(testResourceItemFQN, "members.microsoft_entra_members.#", "0"),
			),
		},
	}))
}

func oneLakeMembersConfig(workspaceID, itemID, roleName string, members map[string]any) string {
	return at.CompileConfig(testResourceItemHeader, map[string]any{
		"workspace_id": workspaceID,
		"item_id":      itemID,
		"role_name":    roleName,
		"decision_rules": []map[string]any{{
			"effect": "Permit",
			"permission": []map[string]any{
				{"attribute_name": "Path", "attribute_value_included_in": []string{"*"}},
				{"attribute_name": "Action", "attribute_value_included_in": []string{"Read"}},
			},
		}},
		"members": members,
	})
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
