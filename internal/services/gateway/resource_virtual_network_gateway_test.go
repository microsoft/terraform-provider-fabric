// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway_test

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/services/gateway"
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
			"display_name":                    *entity.DisplayName,
			"capacity_id":                     *entity.CapacityID,
			"inactivity_minutes_before_sleep": 30,
			"number_of_member_gateways":       3,
			"virtual_network_azure_resource": map[string]any{
				"subscription_id":      "123e4567-e89b-12d3-a456-426614174001",
				"resource_group_name":  "test-rg",
				"virtual_network_name": "test-vnet",
				"subnet_name":          "test-subnet",
			},
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

					// Convert inactivity_minutes_before_sleep from string to int32.
					inactivityStr := is[0].Attributes["inactivity_minutes_before_sleep"]
					inactivityVal, err := strconv.ParseInt(inactivityStr, 10, 32)
					if err != nil {
						return fmt.Errorf("failed to parse inactivity_minutes_before_sleep: %w", err)
					}
					if int32(inactivityVal) != *entity.InactivityMinutesBeforeSleep {
						return errors.New(testResourceVirtualNetworkGatewayFQN + ": unexpected inactivity_minutes_before_sleep")
					}

					// Convert number_of_member_gateways from string to int32.
					memberGatewaysStr := is[0].Attributes["number_of_member_gateways"]
					memberGatewaysVal, err := strconv.ParseInt(memberGatewaysStr, 10, 32)
					if err != nil {
						return fmt.Errorf("failed to parse number_of_member_gateways: %w", err)
					}
					if int32(memberGatewaysVal) != *entity.NumberOfMemberGateways {
						return errors.New(testResourceVirtualNetworkGatewayFQN + ": unexpected number_of_member_gateways")
					}

					if is[0].Attributes["virtual_network_azure_resource.subscription_id"] != *entity.VirtualNetworkAzureResource.SubscriptionID {
						return errors.New(testResourceVirtualNetworkGatewayFQN + ": unexpected virtual_network_azure_resource.subscription_id")
					}

					if is[0].Attributes["virtual_network_azure_resource.resource_group_name"] != *entity.VirtualNetworkAzureResource.ResourceGroupName {
						return errors.New(testResourceVirtualNetworkGatewayFQN + ": unexpected virtual_network_azure_resource.resource_group_name")
					}

					if is[0].Attributes["virtual_network_azure_resource.virtual_network_name"] != *entity.VirtualNetworkAzureResource.VirtualNetworkName {
						return errors.New(testResourceVirtualNetworkGatewayFQN + ": unexpected virtual_network_azure_resource.virtual_network_name")
					}

					if is[0].Attributes["virtual_network_azure_resource.subnet_name"] != *entity.VirtualNetworkAzureResource.SubnetName {
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
						"display_name":                    *entityExist.DisplayName,
						"capacity_id":                     *entityBefore.CapacityID,
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
				ExpectError: regexp.MustCompile(common.ErrorCreateHeader),
			},
			// Create and Read
			{
				ResourceName: testResourceVirtualNetworkGatewayFQN,
				Config: at.CompileConfig(
					testResourceVirtualNetworkGatewayHeader,
					map[string]any{
						"display_name":                    *entityBefore.DisplayName,
						"capacity_id":                     *entityBefore.CapacityID,
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
					resource.TestCheckResourceAttrPtr(testResourceVirtualNetworkGatewayFQN, "display_name", entityBefore.DisplayName),
					resource.TestCheckResourceAttrPtr(testResourceVirtualNetworkGatewayFQN, "capacity_id", entityBefore.CapacityID),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "inactivity_minutes_before_sleep", "30"),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "number_of_member_gateways", "3"),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.subscription_id", "123e4567-e89b-12d3-a456-426614174001"),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.resource_group_name", "test-rg"),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.virtual_network_name", "test-vnet"),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.subnet_name", "test-subnet"),
				),
			},
			// Update and Read
			{
				ResourceName: testResourceVirtualNetworkGatewayFQN,
				Config: at.CompileConfig(
					testResourceVirtualNetworkGatewayHeader,
					map[string]any{
						"display_name":                    *entityAfter.DisplayName,
						"capacity_id":                     *entityAfter.CapacityID,
						"inactivity_minutes_before_sleep": int(*entityAfter.InactivityMinutesBeforeSleep),
						"number_of_member_gateways":       int(*entityAfter.NumberOfMemberGateways),
						"virtual_network_azure_resource": map[string]any{
							"subscription_id":      *entityAfter.VirtualNetworkAzureResource.SubscriptionID,
							"resource_group_name":  *entityAfter.VirtualNetworkAzureResource.ResourceGroupName,
							"virtual_network_name": *entityAfter.VirtualNetworkAzureResource.VirtualNetworkName,
							"subnet_name":          *entityAfter.VirtualNetworkAzureResource.SubnetName,
						},
					},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPtr(testResourceVirtualNetworkGatewayFQN, "display_name", entityAfter.DisplayName),
					resource.TestCheckResourceAttrPtr(testResourceVirtualNetworkGatewayFQN, "capacity_id", entityAfter.CapacityID),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "inactivity_minutes_before_sleep", strconv.Itoa(int(*entityAfter.InactivityMinutesBeforeSleep))),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "number_of_member_gateways", strconv.Itoa(int(*entityAfter.NumberOfMemberGateways))),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.subscription_id", *entityAfter.VirtualNetworkAzureResource.SubscriptionID),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.resource_group_name", *entityAfter.VirtualNetworkAzureResource.ResourceGroupName),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.virtual_network_name", *entityAfter.VirtualNetworkAzureResource.VirtualNetworkName),
					resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.subnet_name", *entityAfter.VirtualNetworkAzureResource.SubnetName),
				),
			},
		},
	))
}

func TestAcc_VirtualNetworkGatewayResource_CRUD(t *testing.T) {
	// Get well-known test values
	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	// Generate random names for testing
	entityCreateDisplayName := testhelp.RandomName()
	entityUpdateDisplayName := testhelp.RandomName()

	initialVirtualNetworkAzureResource := testhelp.WellKnown()["VirtualNetworkInitial"].(map[string]any)
	initSubscriptionID := initialVirtualNetworkAzureResource["subscriptionId"].(string)
	initResourceGroupName := initialVirtualNetworkAzureResource["resourceGroupName"].(string)
	initVirtualNetworkName := initialVirtualNetworkAzureResource["name"].(string)
	initSubnetName := initialVirtualNetworkAzureResource["subnetName"].(string)

	entityInitialInactivityMinutesBeforeSleep := 30
	entityUpdateInactivityMinutesBeforeSleep := 60
	entityInitialNumberOfMemberGateways := int(testhelp.RandomInt32Range(gateway.MinNumberOfMemberGatewaysValues, gateway.MaxNumberOfMemberGatewaysValues))
	entityUpdateNumberOfMemberGateways := int(testhelp.RandomInt32Range(gateway.MinNumberOfMemberGatewaysValues, gateway.MaxNumberOfMemberGatewaysValues))

	updateVirtualNetworkAzureResource := testhelp.WellKnown()["VirtualNetworkUpdate"].(map[string]any)
	updateVirtualNetworkName := updateVirtualNetworkAzureResource["name"].(string)
	updateResourceGroupName := updateVirtualNetworkAzureResource["resourceGroupName"].(string)
	updateSubnetName := updateVirtualNetworkAzureResource["subnetName"].(string)
	updateSubscriptionID := updateVirtualNetworkAzureResource["subscriptionId"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceVirtualNetworkGatewayFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceVirtualNetworkGatewayFQN,
			Config: at.CompileConfig(
				testResourceVirtualNetworkGatewayHeader,
				map[string]any{
					"display_name":                    entityCreateDisplayName,
					"capacity_id":                     capacityID,
					"inactivity_minutes_before_sleep": entityInitialInactivityMinutesBeforeSleep,
					"number_of_member_gateways":       entityInitialNumberOfMemberGateways,
					"virtual_network_azure_resource": map[string]any{
						"subscription_id":      initSubscriptionID,
						"resource_group_name":  initResourceGroupName,
						"virtual_network_name": initVirtualNetworkName,
						"subnet_name":          initSubnetName,
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "display_name", entityCreateDisplayName),
				resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "capacity_id", capacityID),
				resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "inactivity_minutes_before_sleep", strconv.Itoa(entityInitialInactivityMinutesBeforeSleep)),
				resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "number_of_member_gateways", strconv.Itoa(entityInitialNumberOfMemberGateways)),
				resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.subscription_id", initSubscriptionID),
				resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.resource_group_name", initResourceGroupName),
				resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.virtual_network_name", initVirtualNetworkName),
				resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.subnet_name", initSubnetName),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceVirtualNetworkGatewayFQN,
			Config: at.CompileConfig(
				testResourceVirtualNetworkGatewayHeader,
				map[string]any{
					"display_name":                    entityUpdateDisplayName,
					"capacity_id":                     capacityID,
					"inactivity_minutes_before_sleep": entityUpdateInactivityMinutesBeforeSleep,
					"number_of_member_gateways":       entityUpdateNumberOfMemberGateways,
					"virtual_network_azure_resource": map[string]any{
						"subscription_id":      updateSubscriptionID,
						"resource_group_name":  updateResourceGroupName,
						"virtual_network_name": updateVirtualNetworkName,
						"subnet_name":          updateSubnetName,
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "display_name", entityUpdateDisplayName),
				resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "capacity_id", capacityID),
				resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "inactivity_minutes_before_sleep", strconv.Itoa(entityUpdateInactivityMinutesBeforeSleep)),
				resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "number_of_member_gateways", strconv.Itoa(entityUpdateNumberOfMemberGateways)),
				resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.subscription_id", updateSubscriptionID),
				resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.resource_group_name", updateResourceGroupName),
				resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.virtual_network_name", updateVirtualNetworkName),
				resource.TestCheckResourceAttr(testResourceVirtualNetworkGatewayFQN, "virtual_network_azure_resource.subnet_name", updateSubnetName),
			),
		},
	}))
}
