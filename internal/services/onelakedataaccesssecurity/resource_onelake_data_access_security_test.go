// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/provider"
	"github.com/microsoft/terraform-provider-fabric/internal/services/onelakedataaccesssecurity"
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
	objectID := testhelp.RandomUUID()
	tenantID := testhelp.RandomUUID()

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
						"microsoft_entra_members": []map[string]any{
							{"object_id": objectID, "tenant_id": tenantID},
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "role_name", "example"),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "members.microsoft_entra_members.*", map[string]string{
					"object_id": objectID,
					"tenant_id": tenantID,
				}),
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

func TestUnit_OneLakeDataAccessSecurity_ObjectTypeSchema(t *testing.T) {
	t.Parallel()

	var resourceSchema tfresource.SchemaResponse
	onelakedataaccesssecurity.NewResourceOneLakeDataAccessSecurity().Schema(t.Context(), tfresource.SchemaRequest{}, &resourceSchema)
	require.False(t, resourceSchema.Diagnostics.HasError())
	require.False(t, resourceSchema.Schema.ValidateImplementation(t.Context()).HasError())

	members := resourceSchema.Schema.Attributes["members"].(schemaR.SingleNestedAttribute)
	fabricMembers := members.Attributes["fabric_item_members"].(schemaR.SetNestedAttribute)
	require.True(t, fabricMembers.Optional)
	require.True(t, fabricMembers.Computed)
	require.NotNil(t, fabricMembers.Default)

	entraMembers := members.Attributes["microsoft_entra_members"].(schemaR.SetNestedAttribute)
	objectType := entraMembers.NestedObject.Attributes["object_type"].(schemaR.StringAttribute)
	require.True(t, objectType.Optional)
	require.False(t, objectType.Required)
	require.False(t, objectType.Computed)
	require.Contains(t, objectType.DeprecationMessage, "Remove it from your configuration")

	for _, factory := range []func() datasource.DataSource{
		onelakedataaccesssecurity.NewDataSourceOneLakeDataAccessSecurity,
		onelakedataaccesssecurity.NewDataSourceOneLakeDataAccessSecurities,
	} {
		var dataSourceSchema datasource.SchemaResponse
		factory().Schema(t.Context(), datasource.SchemaRequest{}, &dataSourceSchema)
		require.False(t, dataSourceSchema.Diagnostics.HasError())
		require.False(t, dataSourceSchema.Schema.ValidateImplementation(t.Context()).HasError())

		attributes := dataSourceSchema.Schema.Attributes
		if values, exists := attributes["values"]; exists {
			attributes = values.(schemaD.SetNestedAttribute).NestedObject.Attributes
		}

		members := attributes["members"].(schemaD.SingleNestedAttribute)
		entraMembers := members.Attributes["microsoft_entra_members"].(schemaD.SetNestedAttribute)
		objectType := entraMembers.NestedObject.Attributes["object_type"].(schemaD.StringAttribute)
		require.True(t, objectType.Computed)
		require.False(t, objectType.Optional)
		require.Contains(t, objectType.DeprecationMessage, "always returns null")
	}
}

func TestUnit_OneLakeDataAccessSecurity_ObjectTypeDeprecation(t *testing.T) {
	server := providerserver.NewProtocol6(provider.New("testUnit"))()
	member := map[string]any{
		"object_id": testhelp.RandomUUID(),
		"tenant_id": testhelp.RandomUUID(),
	}
	config := map[string]any{
		"workspace_id": testhelp.RandomUUID(),
		"item_id":      testhelp.RandomUUID(),
		"role_name":    "example",
		"decision_rules": []any{map[string]any{
			"effect": "Permit",
			"permission": []any{
				map[string]any{"attribute_name": "Path", "attribute_value_included_in": []string{"*"}},
				map[string]any{"attribute_name": "Action", "attribute_value_included_in": []string{"Read"}},
			},
		}},
		"members": map[string]any{"microsoft_entra_members": []any{member}},
	}

	for _, value := range []any{"User", nil} {
		member["object_type"] = value
		raw, err := json.Marshal(config)
		require.NoError(t, err)

		response, err := server.ValidateResourceConfig(t.Context(), &tfprotov6.ValidateResourceConfigRequest{
			TypeName: itemTypeInfo.FullTypeName(false),
			Config:   &tfprotov6.DynamicValue{JSON: raw},
		})
		require.NoError(t, err)

		warnings := 0

		for _, diagnostic := range response.Diagnostics {
			require.NotEqual(t, tfprotov6.DiagnosticSeverityError, diagnostic.Severity, "%s: %s", diagnostic.Summary, diagnostic.Detail)

			if diagnostic.Summary == "Attribute Deprecated" {
				warnings++

				require.Equal(t, tfprotov6.DiagnosticSeverityWarning, diagnostic.Severity)
				require.Contains(t, diagnostic.Detail, "Remove it from your configuration")
			}
		}

		if value == nil {
			require.Zero(t, warnings)

			continue
		}

		require.Equal(t, 1, warnings)

		upgraded, err := server.UpgradeResourceState(t.Context(), &tfprotov6.UpgradeResourceStateRequest{
			TypeName: itemTypeInfo.FullTypeName(false),
			Version:  0,
			RawState: &tfprotov6.RawState{JSON: raw},
		})
		require.NoError(t, err)

		for _, diagnostic := range upgraded.Diagnostics {
			require.NotEqual(t, tfprotov6.DiagnosticSeverityError, diagnostic.Severity, "%s", diagnostic.Detail)
		}

		schemas, err := server.GetProviderSchema(t.Context(), &tfprotov6.GetProviderSchemaRequest{})
		require.NoError(t, err)
		require.NotNil(t, upgraded.UpgradedState)
		state, err := upgraded.UpgradedState.Unmarshal(schemas.ResourceSchemas[itemTypeInfo.FullTypeName(false)].ValueType())
		require.NoError(t, err)
		require.Contains(t, state.String(), `"User"`)
	}
}

func TestUnit_OneLakeDataAccessSecurityResource_DeprecatedObjectType(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()
	objectID := "ABCDEFAB-1234-4678-9ABC-0123456789AB"
	tenantID := "12345678-ABCD-4123-ABCD-123456789ABC"
	otherTenantID := "FEDCBA98-ABCD-4123-9876-123456789ABC"
	roleID := GenerateOneLakeDataAccessRoleID(workspaceID, itemID, "deprecated-types")

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.CreateOrUpdateSingleDataAccessRole = fakeCreateOrUpdateSingleDataAccessRoleFunc()
	responseUUIDCase := strings.ToLower
	getRole := fakeGetDataAccessRoleFunc()
	fakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.GetDataAccessRole = func(
		ctx context.Context, workspaceID, itemID, roleName string, options *fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleOptions,
	) (azfake.Responder[fabcore.OneLakeDataAccessSecurityClientGetDataAccessRoleResponse], azfake.ErrorResponder) {
		role, exists := fakeOneLakeDataAccessRoleStore[GenerateOneLakeDataAccessRoleID(workspaceID, itemID, roleName)]
		if exists && role.Members != nil {
			for _, member := range role.Members.MicrosoftEntraMembers {
				if member.ObjectID != nil {
					*member.ObjectID = responseUUIDCase(*member.ObjectID)
				}

				if member.TenantID != nil {
					*member.TenantID = responseUUIDCase(*member.TenantID)
				}
			}
		}

		return getRole(ctx, workspaceID, itemID, roleName, options)
	}
	fakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.DeleteDataAccessRole = fakeDeleteDataAccessRoleFunc()

	config := func(firstType any, secondTenant string, secondType any) string {
		members := []map[string]any{
			{"object_id": objectID, "tenant_id": tenantID},
			{"object_id": objectID, "tenant_id": secondTenant},
		}
		if firstType != nil {
			members[0]["object_type"] = firstType
		}

		if secondType != nil {
			members[1]["object_type"] = secondType
		}

		return at.CompileConfig(testResourceItemHeader, map[string]any{
			"workspace_id": workspaceID,
			"item_id":      itemID,
			"role_name":    "deprecated-types",
			"kind":         "Policy",
			"decision_rules": []map[string]any{{
				"effect": "Permit",
				"permission": []map[string]any{
					{"attribute_name": "Path", "attribute_value_included_in": []string{"*"}},
					{"attribute_name": "Action", "attribute_value_included_in": []string{"Read"}},
				},
			}},
			"members": map[string]any{
				"microsoft_entra_members": members,
			},
		}) + fmt.Sprintf(`
			output "object_types_unset" {
				value = alltrue([for member in %s.members.microsoft_entra_members : member.object_type == null])
			}
		`, testResourceItemFQN)
	}
	checkTypes := func(firstType, secondTenant, secondType string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "members.microsoft_entra_members.*", map[string]string{
				"object_id": objectID, "tenant_id": tenantID, "object_type": firstType,
			}),
			resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "members.microsoft_entra_members.*", map[string]string{
				"object_id": objectID, "tenant_id": secondTenant, "object_type": secondType,
			}),
		)
	}
	reverseMembers := func() {
		role := fakeOneLakeDataAccessRoleStore[roleID]
		slices.Reverse(role.Members.MicrosoftEntraMembers)
	}

	legacyConfig := config("User", otherTenantID, "Group")
	updatedConfig := config("ServicePrincipal", otherTenantID, "Group")
	removedConfig := config(nil, otherTenantID, nil)
	newTenantID := testhelp.RandomUUID()
	replacedConfig := strings.NewReplacer(
		objectID, strings.ToLower(objectID),
		tenantID, strings.ToLower(tenantID),
	).Replace(config(nil, newTenantID, nil))

	testCase := testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			Config: legacyConfig,
			Check:  checkTypes("User", otherTenantID, "Group"),
		},
		{
			PreConfig: func() {
				responseUUIDCase = strings.ToUpper

				reverseMembers()
			},
			Config:   legacyConfig,
			PlanOnly: true,
		},
		{
			Config: updatedConfig,
			Check:  checkTypes("ServicePrincipal", otherTenantID, "Group"),
		},
		{
			PreConfig: func() {
				responseUUIDCase = strings.ToLower

				reverseMembers()
			},
			Config:   updatedConfig,
			PlanOnly: true,
		},
		{
			Config: removedConfig,
			Check:  resource.TestCheckOutput("object_types_unset", "true"),
		},
		{
			PreConfig: reverseMembers,
			Config:    removedConfig,
			PlanOnly:  true,
		},
		{
			Config: replacedConfig,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "members.microsoft_entra_members.*", map[string]string{"tenant_id": newTenantID}),
				resource.TestCheckTypeSetElemNestedAttrs(testResourceItemFQN, "members.microsoft_entra_members.*", map[string]string{
					"object_id": strings.ToLower(objectID), "tenant_id": strings.ToLower(tenantID),
				}),
				resource.TestCheckOutput("object_types_unset", "true"),
			),
		},
		{
			ResourceName:                         testResourceItemFQN,
			ImportState:                          true,
			ImportStateId:                        roleID,
			ImportStateVerify:                    true,
			ImportStateVerifyIdentifierAttribute: "role_name",
		},
	})
	transport := &oneLakeRequestTransport{transport: fabfake.NewServerFactoryTransport(fakeServer.ServerFactory)}
	testCase.ProtoV6ProviderFactories = testhelp.GetTestUnitProtoV6ProviderFactories(&fabric.ClientOptions{Transport: transport}, nil)
	resource.Test(t, testCase)

	transport.mu.Lock()
	defer transport.mu.Unlock()

	require.GreaterOrEqual(t, len(transport.bodies), 4)

	for _, body := range transport.bodies {
		require.NotContains(t, body, `"objectType"`)
		require.Contains(t, body, `"objectId"`)
		require.Contains(t, body, `"tenantId"`)
	}
}

func TestUnit_OneLakeDataAccessSecurityResource_FabricItemMembersDefaults(t *testing.T) {
	for name, attribute := range map[string]string{
		"omitted": "",
		"null":    "fabric_item_members = null",
		"empty":   "fabric_item_members = []",
	} {
		t.Run(name, func(t *testing.T) {
			workspaceID := testhelp.RandomUUID()
			itemID := testhelp.RandomUUID()
			objectID := testhelp.RandomUUID()
			tenantID := testhelp.RandomUUID()
			roleID := GenerateOneLakeDataAccessRoleID(workspaceID, itemID, "member-defaults")

			config := func(fabricMembers string) string {
				return fmt.Sprintf(`
					%s {
						workspace_id = "%s"
						item_id      = "%s"
						role_name    = "member-defaults"
						kind         = "Policy"
						decision_rules = [{
							effect = "Permit"
							permission = [
								{ attribute_name = "Path", attribute_value_included_in = ["*"] },
								{ attribute_name = "Action", attribute_value_included_in = ["Read"] }
							]
						}]
						members = {
							%s
							microsoft_entra_members = [{
								object_id = "%s"
								tenant_id = "%s"
							}]
						}
					}
				`, testResourceItemHeader, workspaceID, itemID, fabricMembers, objectID, tenantID)
			}
			emptyConfig := config(attribute)
			populatedConfig := config(fmt.Sprintf(`
				fabric_item_members = [{
					item_access = ["ReadAll"]
					source_path = "%s/%s"
				}]
			`, workspaceID, itemID))

			fakeServer := fakes.NewFakeServer()
			fakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.CreateOrUpdateSingleDataAccessRole = fakeCreateOrUpdateSingleDataAccessRoleFunc()
			fakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.GetDataAccessRole = fakeGetDataAccessRoleFunc()
			fakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.DeleteDataAccessRole = fakeDeleteDataAccessRoleFunc()

			resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
				{
					Config: emptyConfig,
					Check:  resource.TestCheckResourceAttr(testResourceItemFQN, "members.fabric_item_members.#", "0"),
				},
				{
					Config:   emptyConfig,
					PlanOnly: true,
				},
				{
					Config: populatedConfig,
					Check:  resource.TestCheckResourceAttr(testResourceItemFQN, "members.fabric_item_members.#", "1"),
				},
				{
					Config: emptyConfig,
					Check:  resource.TestCheckResourceAttr(testResourceItemFQN, "members.fabric_item_members.#", "0"),
				},
				{
					Config:   emptyConfig,
					PlanOnly: true,
				},
				{
					ResourceName:                         testResourceItemFQN,
					ImportState:                          true,
					ImportStateId:                        roleID,
					ImportStateVerify:                    true,
					ImportStateVerifyIdentifierAttribute: "role_name",
				},
			}))
		})
	}
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
