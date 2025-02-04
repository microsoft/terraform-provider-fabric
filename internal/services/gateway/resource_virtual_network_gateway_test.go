// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway_test

import (
	"errors"
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testResourceVirtualNetworkGatewayFQN    = testhelp.ResourceFQN("fabric", VirtualNetworkItemTFName, "test")
	testResourceVirtualNetworkGatewayHeader = at.ResourceHeader(testhelp.TypeName("fabric", VirtualNetworkItemTFName), "test")
)

func TestUnit_VirtualNetworkGatewayResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(
		t,
		&testResourceVirtualNetworkGatewayFQN,
		fakes.FakeServer.ServerFactory,
		nil,
		[]resource.TestStep{
			// Error: Missing required attribute "display_name"
			{
				ResourceName: testResourceVirtualNetworkGatewayFQN,
				Config: at.CompileConfig(
					testResourceVirtualNetworkGatewayHeader,
					map[string]any{
						"capacity_id":                     "123e4567-e89b-12d3-a456-426614174000",
						"inactivity_minutes_before_sleep": 30,
						"number_of_member_gateways":       3,
						"virtual_network_azure_resource": map[string]any{
							"subscription_id":      "123e4567-e89b-12d3-a456-426614174001",
							"resource_group_name":  "test-rg",
							"virtual_network_name": "test-vnet",
							"subnet_name":          "test-subnet",
						},
					},
				),
				ExpectError: regexp.MustCompile(`The argument "display_name" is required`),
			},
			// Error: Unexpected attribute provided.
			{
				ResourceName: testResourceVirtualNetworkGatewayFQN,
				Config: at.CompileConfig(
					testResourceVirtualNetworkGatewayHeader,
					map[string]any{
						"display_name":                    "test gateway",
						"capacity_id":                     "123e4567-e89b-12d3-a456-426614174000",
						"inactivity_minutes_before_sleep": 30,
						"number_of_member_gateways":       3,
						"virtual_network_azure_resource": map[string]any{
							"subscription_id":      "123e4567-e89b-12d3-a456-426614174001",
							"resource_group_name":  "test-rg",
							"virtual_network_name": "test-vnet",
							"subnet_name":          "test-subnet",
						},
						"unexpected_attr": "test",
					},
				),
				ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
			},
			// Error: Invalid UUID for "capacity_id"
			{
				ResourceName: testResourceVirtualNetworkGatewayFQN,
				Config: at.CompileConfig(
					testResourceVirtualNetworkGatewayHeader,
					map[string]any{
						"display_name":                    "test gateway",
						"capacity_id":                     "not-a-valid-uuid",
						"inactivity_minutes_before_sleep": 30,
						"number_of_member_gateways":       3,
						"virtual_network_azure_resource": map[string]any{
							"subscription_id":      "123e4567-e89b-12d3-a456-426614174001",
							"resource_group_name":  "test-rg",
							"virtual_network_name": "test-vnet",
							"subnet_name":          "test-subnet",
						},
					},
				),
				ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
			},
			// Add a test step for a successful creation/read with all required attributes.
			{
				ResourceName: testResourceVirtualNetworkGatewayFQN,
				Config: at.CompileConfig(
					testResourceVirtualNetworkGatewayHeader,
					map[string]any{
						"display_name":                    "test gateway",
						"capacity_id":                     "123e4567-e89b-12d3-a456-426614174000",
						"inactivity_minutes_before_sleep": 30,
						"number_of_member_gateways":       3,
						"virtual_network_azure_resource": map[string]any{
							"subscription_id":      "123e4567-e89b-12d3-a456-426614174001",
							"resource_group_name":  "test-rg",
							"virtual_network_name": "test-vnet",
							"subnet_name":          "test-subnet",
						},
					},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "display_name", "test gateway"),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "capacity_id", "123e4567-e89b-12d3-a456-426614174000"),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "inactivity_minutes_before_sleep", "30"),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "number_of_member_gateways", "3"),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.subscription_id", "123e4567-e89b-12d3-a456-426614174001"),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.resource_group_name", "test-rg"),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.virtual_network_name", "test-vnet"),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.subnet_name", "test-subnet"),
				),
			},
		},
	))
}

func TestUnit_VirtualNetworkGatewayResource_ImportState(t *testing.T) {
	// Create a fake Virtual Network Gateway.
	entity := fakes.NewRandomVirtualNetworkGateway()

	fakes.FakeServer.Upsert(fakes.NewRandomVirtualNetworkGateway())
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomVirtualNetworkGateway())

	testConfig := at.CompileConfig(
		testResourceVirtualNetworkGatewayHeader,
		map[string]any{
			"display_name": *entity.DisplayName,
			// Other required attributes will be populated from state.
		},
	)

	resource.Test(t, testhelp.NewTestUnitCase(
		t,
		&testResourceVirtualNetworkGatewayFQN,
		fakes.FakeServer.ServerFactory,
		nil,
		[]resource.TestStep{
			{
				ResourceName:  testResourceVirtualNetworkGatewayFQN,
				Config:        testConfig,
				ImportStateId: "not-a-valid-uuid",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
			},
			{
				ResourceName:       testResourceVirtualNetworkGatewayFQN,
				Config:             testConfig,
				ImportStateId:      *entity.ID,
				ImportState:        true,
				ImportStatePersist: true,
				ImportStateCheck: func(is []*terraform.InstanceState) error {
					// Optionally, add additional state validations here.
					if len(is) != 1 {
						return errors.New("expected one instance state")
					}

					if is[0].ID != *entity.ID {
						return errors.New(testResourceVirtualNetworkGatewayFQN + ": unexpected ID")
					}

					if is[0].Attributes["display_name"] != *entity.DisplayName {
						return errors.New(testResourceVirtualNetworkGatewayFQN + ": unexpected display_name")
					}

					if is[0].Attributes["capacity_id"] != *entity.CapacityID {
						return errors.New(testResourceVirtualNetworkGatewayFQN + ": unexpected capacity_id")
					}

					// // Convert inactivity_minutes_before_sleep from string to int32.
					// inactivityStr := is[0].Attributes["inactivity_minutes_before_sleep"]
					// inactivityVal, err := strconv.ParseInt(inactivityStr, 10, 32)
					// if err != nil {
					// 	return fmt.Errorf("failed to parse inactivity_minutes_before_sleep: %w", err)
					// }
					// if int32(inactivityVal) != *entity.InactivityMinutesBeforeSleep {
					// 	return errors.New(testResourceVirtualNetworkGatewayFQN + ": unexpected inactivity_minutes_before_sleep")
					// }

					// // Convert number_of_member_gateways from string to int32.
					// memberGatewaysStr := is[0].Attributes["number_of_member_gateways"]
					// memberGatewaysVal, err := strconv.ParseInt(memberGatewaysStr, 10, 32)
					// if err != nil {
					// 	return fmt.Errorf("failed to parse number_of_member_gateways: %w", err)
					// }
					// if int32(memberGatewaysVal) != *entity.NumberOfMemberGateways {
					// 	return errors.New(testResourceVirtualNetworkGatewayFQN + ": unexpected number_of_member_gateways")
					// }

					if is[0].Attributes["virtual_network_azure_resource.0.subscription_id"] != *entity.VirtualNetworkAzureResource.SubscriptionID {
						return errors.New(testResourceVirtualNetworkGatewayFQN + ": unexpected virtual_network_azure_resource.subscription_id")
					}

					if is[0].Attributes["virtual_network_azure_resource.0.resource_group_name"] != *entity.VirtualNetworkAzureResource.ResourceGroupName {
						return errors.New(testResourceVirtualNetworkGatewayFQN + ": unexpected virtual_network_azure_resource.resource_group_name")
					}

					if is[0].Attributes["virtual_network_azure_resource.0.virtual_network_name"] != *entity.VirtualNetworkAzureResource.VirtualNetworkName {
						return errors.New(testResourceVirtualNetworkGatewayFQN + ": unexpected virtual_network_azure_resource.virtual_network_name")
					}

					if is[0].Attributes["virtual_network_azure_resource.0.subnet_name"] != *entity.VirtualNetworkAzureResource.SubnetName {
						return errors.New(testResourceVirtualNetworkGatewayFQN + ": unexpected virtual_network_azure_resource.subnet_name")
					}

					return nil
				},
			},
		},
	))
}

func TestUnit_VirtualNetworkGatewayResource_CRUD(t *testing.T) {
	// Create fake entities.
	entityExist := fakes.NewRandomVirtualNetworkGateway()
	entityBefore := fakes.NewRandomVirtualNetworkGateway()
	entityAfter := fakes.NewRandomVirtualNetworkGateway()

	// Upsert some fake virtual network gateways.
	fakes.FakeServer.Upsert(fakes.NewRandomVirtualNetworkGateway())
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(fakes.NewRandomVirtualNetworkGateway())

	resource.Test(t, testhelp.NewTestUnitCase(
		t,
		&testResourceVirtualNetworkGatewayFQN,
		fakes.FakeServer.ServerFactory,
		nil,
		[]resource.TestStep{
			// Error: Attempting to create a duplicate gateway (existing entity).
			{
				ResourceName: testResourceVirtualNetworkGatewayFQN,
				Config: at.CompileConfig(
					testResourceVirtualNetworkGatewayHeader,
					map[string]any{
						"display_name": *entityExist.DisplayName,
					},
				),
				ExpectError: regexp.MustCompile(common.ErrorCreateHeader),
			},
			// Create and Read
			{
				ResourceName: testResourceVirtualNetworkGatewayFQN,
				Config: at.CompileConfig(
					testResourceVirtualNetworkGatewayHeader,
					map[string]any{
						"display_name": *entityBefore.DisplayName,
					},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPtr(testResourceVirtualNetworkGatewayFQN, "display_name", entityBefore.DisplayName),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "capacity_id", *entityBefore.CapacityID),
				),
			},
			// Update and Read
			{
				ResourceName: testResourceVirtualNetworkGatewayFQN,
				Config: at.CompileConfig(
					testResourceVirtualNetworkGatewayHeader,
					map[string]any{
						"display_name": *entityBefore.DisplayName,
						// Simulate update by changing capacity_id.
						"capacity_id": *entityAfter.CapacityID,
					},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPtr(testResourceVirtualNetworkGatewayFQN, "display_name", entityBefore.DisplayName),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "capacity_id", *entityAfter.CapacityID),
				),
			},
		},
	))
}
