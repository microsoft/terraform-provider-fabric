// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway_test

import (
	"regexp"
	"strconv"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testDataSourceItemFQN    = testhelp.DataSourceFQN("fabric", itemTFName, "test")
	testDataSourceItemHeader = at.DataSourceHeader(testhelp.TypeName("fabric", itemTFName), "test")
)

func TestUnit_GatewayDataSource(t *testing.T) {
	virtualNetworkGateway := fakes.NewRandomVirtualNetworkGateway()
	onPremisesGateway := fakes.NewRandomOnPremisesGateway()
	onPremisesGatewayPersonalGateway := fakes.NewRandomOnPremisesGatewayPersonal()

	fakes.FakeServer.Upsert(fakes.NewRandomGateway())
	fakes.FakeServer.Upsert(virtualNetworkGateway)
	fakes.FakeServer.Upsert(onPremisesGateway)
	fakes.FakeServer.Upsert(onPremisesGatewayPersonalGateway)
	fakes.FakeServer.Upsert(fakes.NewRandomGateway())

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing Attribute Configuration`),
		},
		// error - id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id":              *virtualNetworkGateway.ID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// read by id - virtual network
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": *virtualNetworkGateway.ID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "id", virtualNetworkGateway.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "type", (*string)(virtualNetworkGateway.Type)),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "display_name", virtualNetworkGateway.DisplayName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "capacity_id", virtualNetworkGateway.CapacityID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "inactivity_minutes_before_sleep", strconv.Itoa(int(*virtualNetworkGateway.InactivityMinutesBeforeSleep))),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "number_of_member_gateways", strconv.Itoa(int(*virtualNetworkGateway.NumberOfMemberGateways))),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "virtual_network_azure_resource.resource_group_name", virtualNetworkGateway.VirtualNetworkAzureResource.ResourceGroupName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "virtual_network_azure_resource.subnet_name", virtualNetworkGateway.VirtualNetworkAzureResource.SubnetName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "virtual_network_azure_resource.subscription_id", virtualNetworkGateway.VirtualNetworkAzureResource.SubscriptionID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "virtual_network_azure_resource.virtual_network_name", virtualNetworkGateway.VirtualNetworkAzureResource.VirtualNetworkName),
			),
		},
		// read by name - virtual network
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"display_name": *virtualNetworkGateway.DisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "id", virtualNetworkGateway.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "type", (*string)(virtualNetworkGateway.Type)),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "display_name", virtualNetworkGateway.DisplayName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "capacity_id", virtualNetworkGateway.CapacityID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "inactivity_minutes_before_sleep", strconv.Itoa(int(*virtualNetworkGateway.InactivityMinutesBeforeSleep))),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "number_of_member_gateways", strconv.Itoa(int(*virtualNetworkGateway.NumberOfMemberGateways))),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "virtual_network_azure_resource.resource_group_name", virtualNetworkGateway.VirtualNetworkAzureResource.ResourceGroupName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "virtual_network_azure_resource.subnet_name", virtualNetworkGateway.VirtualNetworkAzureResource.SubnetName),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "virtual_network_azure_resource.subscription_id", virtualNetworkGateway.VirtualNetworkAzureResource.SubscriptionID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "virtual_network_azure_resource.virtual_network_name", virtualNetworkGateway.VirtualNetworkAzureResource.VirtualNetworkName),
			),
		},
		// read by id - on premises
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": *onPremisesGateway.ID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "id", onPremisesGateway.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "type", (*string)(onPremisesGateway.Type)),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "display_name", onPremisesGateway.DisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "allow_cloud_connection_refresh", strconv.FormatBool(*onPremisesGateway.AllowCloudConnectionRefresh)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "allow_custom_connectors", strconv.FormatBool(*onPremisesGateway.AllowCustomConnectors)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "number_of_member_gateways", strconv.Itoa(int(*onPremisesGateway.NumberOfMemberGateways))),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "load_balancing_setting", string(*onPremisesGateway.LoadBalancingSetting)),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "version", onPremisesGateway.Version),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "public_key.exponent", onPremisesGateway.PublicKey.Exponent),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "public_key.modulus", onPremisesGateway.PublicKey.Modulus),
			),
		},
		// read by name - on premises
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"display_name": *onPremisesGateway.DisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "id", onPremisesGateway.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "type", (*string)(onPremisesGateway.Type)),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "display_name", onPremisesGateway.DisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "allow_cloud_connection_refresh", strconv.FormatBool(*onPremisesGateway.AllowCloudConnectionRefresh)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "allow_custom_connectors", strconv.FormatBool(*onPremisesGateway.AllowCustomConnectors)),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "number_of_member_gateways", strconv.Itoa(int(*onPremisesGateway.NumberOfMemberGateways))),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "load_balancing_setting", string(*onPremisesGateway.LoadBalancingSetting)),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "version", onPremisesGateway.Version),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "public_key.exponent", onPremisesGateway.PublicKey.Exponent),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "public_key.modulus", onPremisesGateway.PublicKey.Modulus),
			),
		},
		// read by id - on premises personal
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": *onPremisesGatewayPersonalGateway.ID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "id", onPremisesGatewayPersonalGateway.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "type", (*string)(onPremisesGatewayPersonalGateway.Type)),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "version", onPremisesGatewayPersonalGateway.Version),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "public_key.exponent", onPremisesGatewayPersonalGateway.PublicKey.Exponent),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "public_key.modulus", onPremisesGatewayPersonalGateway.PublicKey.Modulus),
			),
		},
		// read by id - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by name - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"display_name": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}

func TestAcc_GatewayDataSource(t *testing.T) {
	entityVirtualNetwork := testhelp.WellKnown()["GatewayVirtualNetwork"].(map[string]any)
	entityVirtualNetworkID := entityVirtualNetwork["id"].(string)
	entityVirtualNetworkDisplayName := entityVirtualNetwork["displayName"].(string)
	entityVirtualNetworkType := entityVirtualNetwork["type"].(string)

	entityOnPremises := testhelp.WellKnown()["GatewayOnPremises"].(map[string]any)
	entityOnPremisesID := entityOnPremises["id"].(string)
	entityOnPremisesDisplayName := entityOnPremises["displayName"].(string)
	entityOnPremisesType := entityOnPremises["type"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read by id - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by id - virtual network
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": entityVirtualNetworkID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", entityVirtualNetworkID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "display_name", entityVirtualNetworkDisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "type", entityVirtualNetworkType),
			),
		},
		// read by name- virtual network
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"display_name": entityVirtualNetworkDisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", entityVirtualNetworkID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "display_name", entityVirtualNetworkDisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "type", entityVirtualNetworkType),
			),
		},
		// read by id - on premises
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": entityOnPremisesID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", entityOnPremisesID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "display_name", entityOnPremisesDisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "type", entityOnPremisesType),
			),
		},
		// read by name - on premises
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"display_name": entityOnPremisesDisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", entityOnPremisesID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "display_name", entityOnPremisesDisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "type", entityOnPremisesType),
			),
		},
	}))
}
