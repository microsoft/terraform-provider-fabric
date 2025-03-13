// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway_test

import (
	"fmt"
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/services/gateway"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testResourceGatewayRoleAssignment       = testhelp.ResourceFQN("fabric", gatewayRoleAssignmentTFName, "test")
	testResourceGatewayRoleAssignmentHeader = at.ResourceHeader(testhelp.TypeName("fabric", gatewayRoleAssignmentTFName), "test")
)

func TestUnit_GatewayRoleAssignmentResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceGatewayRoleAssignment, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no required attributes - gateway_id
		{
			ResourceName: testResourceGatewayRoleAssignment,
			Config: at.CompileConfig(
				testResourceGatewayRoleAssignmentHeader,
				map[string]any{
					"principal": map[string]any{
						"id":   "00000000-0000-0000-0000-000000000000",
						"type": "User",
					},
					"role": "ConnectionCreatorWithResharing",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "gateway_id" is required, but no definition was found.`),
		},
		// error - no required attributes - principal.id
		{
			ResourceName: testResourceGatewayRoleAssignment,
			Config: at.CompileConfig(
				testResourceGatewayRoleAssignmentHeader,
				map[string]any{
					"gateway_id": "00000000-0000-0000-0000-000000000000",
					"principal": map[string]any{
						"type": "User",
					},
					"role": "ConnectionCreatorWithResharing",
				},
			),
			ExpectError: regexp.MustCompile(`Inappropriate value for attribute "principal": attribute "id" is required.`),
		},
		// error - no required attributes - principal.type
		{
			ResourceName: testResourceGatewayRoleAssignment,
			Config: at.CompileConfig(
				testResourceGatewayRoleAssignmentHeader,
				map[string]any{
					"gateway_id": "00000000-0000-0000-0000-000000000000",
					"principal": map[string]any{
						"id": "00000000-0000-0000-0000-000000000000",
					},
					"role": "ConnectionCreatorWithResharing",
				},
			),
			ExpectError: regexp.MustCompile(`Inappropriate value for attribute "principal": attribute "type" is required.`),
		},
		// error - no required attributes - role
		{
			ResourceName: testResourceGatewayRoleAssignment,
			Config: at.CompileConfig(
				testResourceGatewayRoleAssignmentHeader,
				map[string]any{
					"gateway_id": "00000000-0000-0000-0000-000000000000",
					"principal": map[string]any{
						"id":   "00000000-0000-0000-0000-000000000000",
						"type": "User",
					},
				},
			),
			ExpectError: regexp.MustCompile(`The argument "role" is required, but no definition was found.`),
		},
		// error - invalid UUID - gateway_id
		{
			ResourceName: testResourceGatewayRoleAssignment,
			Config: at.CompileConfig(
				testResourceGatewayRoleAssignmentHeader,
				map[string]any{
					"gateway_id": "invalid uuid",
					"principal": map[string]any{
						"id":   "00000000-0000-0000-0000-000000000000",
						"type": "User",
					},
					"role": "ConnectionCreatorWithResharing",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - invalid UUID - principal.id
		{
			ResourceName: testResourceGatewayRoleAssignment,
			Config: at.CompileConfig(
				testResourceGatewayRoleAssignmentHeader,
				map[string]any{
					"gateway_id": "00000000-0000-0000-0000-000000000000",
					"principal": map[string]any{
						"id":   "invalid uuid",
						"type": "User",
					},
					"role": "ConnectionCreatorWithResharing",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_GatewayRoleAssignmentResource_ImportState(t *testing.T) {
	testCase := at.CompileConfig(
		testResourceGatewayRoleAssignmentHeader,
		map[string]any{},
	)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceGatewayRoleAssignment, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceGatewayRoleAssignment,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile("GatewayID/GatewayRoleAssignmentID"),
		},
		{
			ResourceName:  testResourceGatewayRoleAssignment,
			Config:        testCase,
			ImportStateId: "test/id",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceGatewayRoleAssignment,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s", "test", "00000000-0000-0000-0000-000000000000"),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceGatewayRoleAssignment,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s", "00000000-0000-0000-0000-000000000000", "test"),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestAcc_GatewayRoleAssignmentResource_CRUD(t *testing.T) {
	if testhelp.ShouldSkipTest(t) {
		t.Skip("No SPN support")
	}

	gatewayType := string(fabcore.GatewayTypeVirtualNetwork)
	gatewayCreateDisplayName := testhelp.RandomName()
	gatewayCreateInactivityMinutesBeforeSleep := int(testhelp.RandomElement(gateway.PossibleInactivityMinutesBeforeSleepValues))
	gatewayCreateNumberOfMemberGateways := int(testhelp.RandomIntRange(gateway.MinNumberOfMemberGatewaysValues, gateway.MaxNumberOfMemberGatewaysValues))

	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	virtualNetworkAzureResource := testhelp.WellKnown()["VirtualNetwork01"].(map[string]any)
	virtualNetworkName := virtualNetworkAzureResource["name"].(string)
	resourceGroupName := virtualNetworkAzureResource["resourceGroupName"].(string)
	subnetName := virtualNetworkAzureResource["subnetName"].(string)
	subscriptionID := virtualNetworkAzureResource["subscriptionId"].(string)

	gatewayResourceHCL := at.CompileConfig(
		at.ResourceHeader(testhelp.TypeName("fabric", itemTFName), "test"),
		map[string]any{
			"type":                            gatewayType,
			"display_name":                    gatewayCreateDisplayName,
			"inactivity_minutes_before_sleep": gatewayCreateInactivityMinutesBeforeSleep,
			"number_of_member_gateways":       gatewayCreateNumberOfMemberGateways,
			"virtual_network_azure_resource": map[string]any{
				"virtual_network_name": virtualNetworkName,
				"resource_group_name":  resourceGroupName,
				"subnet_name":          subnetName,
				"subscription_id":      subscriptionID,
			},
			"capacity_id": capacityID,
		},
	)
	gatewayResourceFQN := testhelp.ResourceFQN("fabric", itemTFName, "test")

	entity := testhelp.WellKnown()["Group"].(map[string]any)
	entityID := entity["id"].(string)
	entityType := entity["type"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceGatewayRoleAssignment, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceGatewayRoleAssignment,
			Config: at.JoinConfigs(
				gatewayResourceHCL,
				at.CompileConfig(
					testResourceGatewayRoleAssignmentHeader,
					map[string]any{
						"gateway_id": testhelp.RefByFQN(gatewayResourceFQN, "id"),
						"principal": map[string]any{
							"id":   entityID,
							"type": entityType,
						},
						"role": "ConnectionCreatorWithResharing",
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceGatewayRoleAssignment, "principal.id", entityID),
				resource.TestCheckResourceAttr(testResourceGatewayRoleAssignment, "principal.type", entityType),
				resource.TestCheckResourceAttr(testResourceGatewayRoleAssignment, "role", "ConnectionCreatorWithResharing"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceGatewayRoleAssignment,
			Config: at.JoinConfigs(
				gatewayResourceHCL,
				at.CompileConfig(
					testResourceGatewayRoleAssignmentHeader,
					map[string]any{
						"gateway_id": testhelp.RefByFQN(gatewayResourceFQN, "id"),
						"principal": map[string]any{
							"id":   entityID,
							"type": entityType,
						},
						"role": "ConnectionCreator",
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceGatewayRoleAssignment, "principal.id", entityID),
				resource.TestCheckResourceAttr(testResourceGatewayRoleAssignment, "principal.type", entityType),
				resource.TestCheckResourceAttr(testResourceGatewayRoleAssignment, "role", "ConnectionCreator"),
			),
		},
	}))
}
