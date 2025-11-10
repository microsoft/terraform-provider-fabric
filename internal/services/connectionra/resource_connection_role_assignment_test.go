// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connectionra_test

import (
	"fmt"
	"regexp"
	"testing"

	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_ConnectionRoleAssignmentResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no required attributes - connection_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"principal": map[string]any{
						"id":   "00000000-0000-0000-0000-000000000000",
						"type": "User",
					},
					"role": "UserWithReshare",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "connection_id" is required, but no definition was found.`),
		},
		// error - no required attributes - principal.id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"connection_id": "00000000-0000-0000-0000-000000000000",
					"principal": map[string]any{
						"type": "User",
					},
					"role": "UserWithReshare",
				},
			),
			ExpectError: regexp.MustCompile(`Inappropriate value for attribute "principal": attribute "id" is required.`),
		},
		// error - no required attributes - principal.type
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"connection_id": "00000000-0000-0000-0000-000000000000",
					"principal": map[string]any{
						"id": "00000000-0000-0000-0000-000000000000",
					},
					"role": "UserWithReshare",
				},
			),
			ExpectError: regexp.MustCompile(`Inappropriate value for attribute "principal": attribute "type" is required.`),
		},
		// error - no required attributes - role
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"connection_id": "00000000-0000-0000-0000-000000000000",
					"principal": map[string]any{
						"id":   "00000000-0000-0000-0000-000000000000",
						"type": "User",
					},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "role" is required, but no definition was found.`),
		},
		// error - invalid UUID - connection_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"connection_id": "invalid uuid",
					"principal": map[string]any{
						"id":   "00000000-0000-0000-0000-000000000000",
						"type": "User",
					},
					"role": "UserWithReshare",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - invalid UUID - principal.id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"connection_id": "00000000-0000-0000-0000-000000000000",
					"principal": map[string]any{
						"id":   "invalid uuid",
						"type": "User",
					},
					"role": "UserWithReshare",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_ConnectionRoleAssignmentResource_ImportState(t *testing.T) {
	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{},
	)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile("ConnectionID/ConnectionRoleAssignmentID"),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "test/id",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s", "test", "00000000-0000-0000-0000-000000000000"),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s", "00000000-0000-0000-0000-000000000000", "test"),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_ConnectionRoleAssignmentResource_CRUD(t *testing.T) {
	connectionID := testhelp.RandomUUID()
	entity := NewRandomConnectionRoleAssignment()

	entityUpdate := entity
	entityUpdate.Role = azto.Ptr(fabcore.ConnectionRoleUserWithReshare)

	fakes.FakeServer.ServerFactory.Core.ConnectionsServer.AddConnectionRoleAssignment = fakeAddConnectionRoleAssignment(entity)
	fakes.FakeServer.ServerFactory.Core.ConnectionsServer.GetConnectionRoleAssignment = fakeGetConnectionRoleAssignment(entity)
	fakes.FakeServer.ServerFactory.Core.ConnectionsServer.DeleteConnectionRoleAssignment = fakeDeleteConnectionRoleAssignment()
	fakes.FakeServer.ServerFactory.Core.ConnectionsServer.UpdateConnectionRoleAssignment = fakeUpdateConnectionRoleAssignment(entityUpdate)

	entityID := *entity.Principal.ID
	entityType := (string)(*entity.Principal.Type)
	role := (string)(*entity.Role)
	updatedRole := (string)(*entityUpdate.Role)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"connection_id": connectionID,
					"principal": map[string]any{
						"id":   entityID,
						"type": entityType,
					},
					"role": role,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "principal.id", entityID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "principal.type", entityType),
				resource.TestCheckResourceAttr(testResourceItemFQN, "role", role),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"id":            *entity.ID,
					"connection_id": connectionID,
					"principal": map[string]any{
						"id":   entityID,
						"type": entityType,
					},
					"role": updatedRole,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "id", *entity.ID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "principal.id", entityID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "principal.type", entityType),
				resource.TestCheckResourceAttr(testResourceItemFQN, "role", updatedRole),
			),
		},
	}))
}

func TestAcc_ConnectionRoleAssignmentResource_CRUD(t *testing.T) {
	entityVirtualNetwork := testhelp.WellKnown()["GatewayVirtualNetwork"].(map[string]any)
	entityVirtualNetworkID := entityVirtualNetwork["id"].(string)

	connectionHCL := at.CompileConfig(
		at.ResourceHeader(testhelp.TypeName(common.ProviderTypeName, "connection"), "test"),
		map[string]any{
			"display_name":      testhelp.RandomName(),
			"connectivity_type": "VirtualNetworkGateway",
			"privacy_level":     "Organizational",
			"gateway_id":        entityVirtualNetworkID,
			"connection_details": map[string]any{
				"type":            "FTP",
				"creation_method": "FTP.Contents",
				"parameters": []map[string]any{
					{
						"name":  "server",
						"value": "ftp.example.com",
					},
				},
			},
			"credential_details": map[string]any{
				"connection_encryption": string(fabcore.ConnectionEncryptionNotEncrypted),
				"single_sign_on_type":   string(fabcore.SingleSignOnTypeNone),
				"skip_test_connection":  false,
				"credential_type":       string(fabcore.CredentialTypeAnonymous),
			},
		},
	)

	connectionFQN := testhelp.ResourceFQN(common.ProviderTypeName, "connection", "test")

	entity := testhelp.WellKnown()["Principal"].(map[string]any)
	entityID := entity["id"].(string)
	entityType := entity["type"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				connectionHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"connection_id": testhelp.RefByFQN(connectionFQN, "id"),
						"principal": map[string]any{
							"id":   entityID,
							"type": entityType,
						},
						"role": "User",
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "principal.id", entityID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "principal.type", entityType),
				resource.TestCheckResourceAttr(testResourceItemFQN, "role", "User"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				connectionHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"connection_id": testhelp.RefByFQN(connectionFQN, "id"),
						"principal": map[string]any{
							"id":   entityID,
							"type": entityType,
						},
						"role": "UserWithReshare",
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "principal.id", entityID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "principal.type", entityType),
				resource.TestCheckResourceAttr(testResourceItemFQN, "role", "UserWithReshare"),
			),
		},
	}))
}
