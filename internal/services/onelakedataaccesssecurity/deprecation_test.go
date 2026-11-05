// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/provider"
	"github.com/microsoft/terraform-provider-fabric/internal/services/onelakedataaccesssecurity"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

func TestUnit_OneLakeDataAccessSecurity_ObjectTypeSchema(t *testing.T) {
	t.Parallel()

	var resourceSchema tfresource.SchemaResponse
	onelakedataaccesssecurity.NewResourceOneLakeDataAccessSecurity().Schema(t.Context(), tfresource.SchemaRequest{}, &resourceSchema)
	require.False(t, resourceSchema.Diagnostics.HasError())
	require.False(t, resourceSchema.Schema.ValidateImplementation(t.Context()).HasError())

	members := resourceSchema.Schema.Attributes["members"].(schemaR.SingleNestedAttribute)
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

type oneLakeRequestTransport struct {
	transport policy.Transporter
	mu        sync.Mutex
	bodies    []string
}

func (t *oneLakeRequestTransport) Do(req *http.Request) (*http.Response, error) {
	if req.Body != nil && req.Body != http.NoBody {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("reading OneLake request body: %w", err)
		}

		err = req.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("closing OneLake request body: %w", err)
		}

		req.Body = io.NopCloser(bytes.NewReader(body))

		t.mu.Lock()
		t.bodies = append(t.bodies, string(body))
		t.mu.Unlock()
	}

	return t.transport.Do(req)
}

func TestUnit_OneLakeDataAccessSecurityResource_DeprecatedObjectType(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	itemID := testhelp.RandomUUID()
	objectID := testhelp.RandomUUID()
	tenantID := testhelp.RandomUUID()
	otherTenantID := testhelp.RandomUUID()
	roleID := GenerateOneLakeDataAccessRoleID(workspaceID, itemID, "deprecated-types")

	fakeServer := fakes.NewFakeServer()
	fakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.CreateOrUpdateSingleDataAccessRole = fakeCreateOrUpdateSingleDataAccessRoleFunc()
	fakeServer.ServerFactory.Core.OneLakeDataAccessSecurityServer.GetDataAccessRole = fakeGetDataAccessRoleFunc()
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
				"fabric_item_members":     []map[string]any{},
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
	replacedConfig := config(nil, newTenantID, nil)

	testCase := testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakeServer.ServerFactory, nil, []resource.TestStep{
		{
			Config: legacyConfig,
			Check:  checkTypes("User", otherTenantID, "Group"),
		},
		{
			PreConfig: reverseMembers,
			Config:    legacyConfig,
			PlanOnly:  true,
		},
		{
			Config: updatedConfig,
			Check:  checkTypes("ServicePrincipal", otherTenantID, "Group"),
		},
		{
			PreConfig: reverseMembers,
			Config:    updatedConfig,
			PlanOnly:  true,
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
