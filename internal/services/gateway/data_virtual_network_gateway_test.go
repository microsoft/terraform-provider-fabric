// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway_test

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

var (
	testDataSourceVirtualNetworkFQN    = testhelp.DataSourceFQN("fabric", VirtualNetworkItemTFName, "test")
	testVirtualNetworkDataSourceHeader = at.DataSourceHeader(testhelp.TypeName("fabric", VirtualNetworkItemTFName), "test")
)

func TestUnit_VirtualNetworkGatewayDataSource(t *testing.T) {
	entity := fakes.NewRandomVirtualNetworkGateway()

	fakes.FakeServer.Upsert(fakes.NewRandomVirtualNetworkGateway())
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomVirtualNetworkGateway())

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testVirtualNetworkDataSourceHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Exactly one of these attributes must be configured: \[id,display_name\]`),
		},
		// error - id - invalid UUID
		{
			Config: at.CompileConfig(
				testVirtualNetworkDataSourceHeader,
				map[string]any{
					"id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			Config: at.CompileConfig(
				testVirtualNetworkDataSourceHeader,
				map[string]any{
					"id":              *entity.ID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - conflicting attributes
		{
			Config: at.CompileConfig(
				testVirtualNetworkDataSourceHeader,
				map[string]any{
					"id":           *entity.ID,
					"display_name": *entity.DisplayName,
				},
			),
			ExpectError: regexp.MustCompile(`These attributes cannot be configured together: \[id,display_name\]`),
		},
		// read by id - not found
		{
			Config: at.CompileConfig(
				testVirtualNetworkDataSourceHeader,
				map[string]any{
					"id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by id
		{
			Config: at.CompileConfig(
				testVirtualNetworkDataSourceHeader,
				map[string]any{
					"id": *entity.ID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceVirtualNetworkFQN, "id", *entity.ID),
				resource.TestCheckResourceAttr(testDataSourceVirtualNetworkFQN, "display_name", *entity.DisplayName),
				resource.TestCheckResourceAttrSet(testDataSourceVirtualNetworkFQN, "inactivity_minutes_before_sleep"),
				resource.TestCheckResourceAttrSet(testDataSourceVirtualNetworkFQN, "capacity_id"),
				resource.TestCheckResourceAttrSet(testDataSourceVirtualNetworkFQN, "number_of_member_gateways"),
				resource.TestCheckResourceAttrSet(testDataSourceVirtualNetworkFQN, "virtual_network_azure_resource.subscription_id"),
				resource.TestCheckResourceAttrSet(testDataSourceVirtualNetworkFQN, "virtual_network_azure_resource.resource_group_name"),
				resource.TestCheckResourceAttrSet(testDataSourceVirtualNetworkFQN, "virtual_network_azure_resource.virtual_network_name"),
				resource.TestCheckResourceAttrSet(testDataSourceVirtualNetworkFQN, "virtual_network_azure_resource.subnet_name"),
			),
		},
		// read by name - not found
		{
			Config: at.CompileConfig(
				testVirtualNetworkDataSourceHeader,
				map[string]any{
					"display_name": testhelp.RandomName(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by name
		{
			Config: at.CompileConfig(
				testVirtualNetworkDataSourceHeader,
				map[string]any{
					"display_name": *entity.DisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceVirtualNetworkFQN, "id", *entity.ID),
				resource.TestCheckResourceAttr(testDataSourceVirtualNetworkFQN, "display_name", *entity.DisplayName),
				resource.TestCheckResourceAttrSet(testDataSourceVirtualNetworkFQN, "inactivity_minutes_before_sleep"),
				resource.TestCheckResourceAttrSet(testDataSourceVirtualNetworkFQN, "capacity_id"),
				resource.TestCheckResourceAttrSet(testDataSourceVirtualNetworkFQN, "number_of_member_gateways"),
				resource.TestCheckResourceAttrSet(testDataSourceVirtualNetworkFQN, "virtual_network_azure_resource.subscription_id"),
				resource.TestCheckResourceAttrSet(testDataSourceVirtualNetworkFQN, "virtual_network_azure_resource.resource_group_name"),
				resource.TestCheckResourceAttrSet(testDataSourceVirtualNetworkFQN, "virtual_network_azure_resource.virtual_network_name"),
				resource.TestCheckResourceAttrSet(testDataSourceVirtualNetworkFQN, "virtual_network_azure_resource.subnet_name"),
			),
		},
	}))
}

func TestAcc_VirtualNetworkGatewayDataSource(t *testing.T) {
	entity := testhelp.WellKnown()["VirtualNetworkGateway"].(map[string]any)
	entityID := entity["id"].(string)
	entityDisplayName := entity["displayName"].(string)
	entityDescription := entity["description"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(
				testVirtualNetworkDataSourceHeader,
				map[string]any{
					"id": entityID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceVirtualNetworkFQN, "id", entityID),
				resource.TestCheckResourceAttr(testDataSourceVirtualNetworkFQN, "display_name", entityDisplayName),
				resource.TestCheckResourceAttr(testDataSourceVirtualNetworkFQN, "description", entityDescription),
			),
		},
		// read by id - not found
		{
			Config: at.CompileConfig(
				testVirtualNetworkDataSourceHeader,
				map[string]any{
					"id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}
