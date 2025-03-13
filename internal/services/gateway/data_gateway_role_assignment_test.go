// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/services/gateway"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testDataSourceGatewayRoleAssignmentFQN    = testhelp.DataSourceFQN("fabric", gatewayRoleAssignmentTFName, "test")
	testDataSourceGatewayRoleAssignmentHeader = at.DataSourceHeader(testhelp.TypeName("fabric", gatewayRoleAssignmentTFName), "test")
)

func TestUnit_GatewayRoleAssignmentDataSource(t *testing.T) {
	gatewayID := testhelp.RandomUUID()
	entity := NewRandomGatewayRoleAssignment()
	fakes.FakeServer.ServerFactory.Core.GatewaysServer.GetGatewayRoleAssignment = fakeGatewayRoleAssignment(entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceGatewayRoleAssignmentHeader,
				map[string]any{
					"gateway_id":      gatewayID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceGatewayRoleAssignmentHeader,
				map[string]any{
					"id":         *entity.ID,
					"gateway_id": gatewayID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceGatewayRoleAssignmentFQN, "gateway_id", gatewayID),
				resource.TestCheckResourceAttrPtr(testDataSourceGatewayRoleAssignmentFQN, "id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceGatewayRoleAssignmentFQN, "role", (*string)(entity.Role)),
				resource.TestCheckResourceAttrPtr(testDataSourceGatewayRoleAssignmentFQN, "principal.id", entity.Principal.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceGatewayRoleAssignmentFQN, "principal.type", (*string)(entity.Principal.Type)),
			),
		},
	}))
}

func TestAcc_GatewayRoleAssignmentDataSource(t *testing.T) {
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

	principal := testhelp.WellKnown()["Principal"].(map[string]any)
	principalID := principal["id"].(string)
	principalType := principal["type"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.JoinConfigs(
				gatewayResourceHCL,
				at.CompileConfig(
					testResourceGatewayRoleAssignmentHeader,
					map[string]any{
						"gateway_id": testhelp.RefByFQN(gatewayResourceFQN, "id"),
						"principal": map[string]any{
							"id":   principalID,
							"type": principalType,
						},
						"role": "ConnectionCreatorWithResharing",
					},
				),
				at.CompileConfig(
					testDataSourceGatewayRoleAssignmentHeader,
					map[string]any{
						"id":         principalID,
						"gateway_id": testhelp.RefByFQN(testResourceGatewayRoleAssignment, "gateway_id"),
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceGatewayRoleAssignmentFQN, "gateway_id"),
				resource.TestCheckResourceAttr(testDataSourceGatewayRoleAssignmentFQN, "id", principalID),
				resource.TestCheckResourceAttr(testDataSourceGatewayRoleAssignmentFQN, "role", "ConnectionCreatorWithResharing"),
				resource.TestCheckResourceAttr(testDataSourceGatewayRoleAssignmentFQN, "principal.id", principalID),
				resource.TestCheckResourceAttr(testDataSourceGatewayRoleAssignmentFQN, "principal.type", principalType),
			),
		},
	},
	))
}
