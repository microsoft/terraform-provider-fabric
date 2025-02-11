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

	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/services/gateway"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testResourceItemFQN    = testhelp.ResourceFQN("fabric", itemTFName, "test")
	testResourceItemHeader = at.ResourceHeader(testhelp.TypeName("fabric", itemTFName), "test")
)

func TestUnit_GatewayResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - missing attributes
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
					"type":            string(fabcore.GatewayTypeVirtualNetwork),
					"display_name":    "test",
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - missing required attributes - type
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "type" is required, but no definition was found.`),
		},
		// error - missing required attributes - display_name
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"type":                            string(fabcore.GatewayTypeVirtualNetwork),
					"inactivity_minutes_before_sleep": (int)(gateway.PossibleInactivityMinutesBeforeSleepValues[0]),
					"number_of_member_gateways":       (int)(gateway.MinNumberOfMemberGatewaysValues),
					"virtual_network_azure_resource": map[string]any{
						"resource_group_name":  "test",
						"virtual_network_name": "test",
						"subnet_name":          "test",
						"subscription_id":      "00000000-0000-0000-0000-000000000000",
					},
					"capacity_id": "00000000-0000-0000-0000-000000000000",
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute display_name`),
		},
		// error - missing required attributes - inactivity_minutes_before_sleep
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"type":                      string(fabcore.GatewayTypeVirtualNetwork),
					"display_name":              "test",
					"number_of_member_gateways": (int)(gateway.MinNumberOfMemberGatewaysValues),
					"virtual_network_azure_resource": map[string]any{
						"resource_group_name":  "test",
						"virtual_network_name": "test",
						"subnet_name":          "test",
						"subscription_id":      "00000000-0000-0000-0000-000000000000",
					},
					"capacity_id": "00000000-0000-0000-0000-000000000000",
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute inactivity_minutes_before_sleep`),
		},
		// error - missing required attributes - number_of_member_gateways
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"type":                            string(fabcore.GatewayTypeVirtualNetwork),
					"display_name":                    "test",
					"inactivity_minutes_before_sleep": (int)(gateway.PossibleInactivityMinutesBeforeSleepValues[0]),
					"virtual_network_azure_resource": map[string]any{
						"resource_group_name":  "test",
						"virtual_network_name": "test",
						"subnet_name":          "test",
						"subscription_id":      "00000000-0000-0000-0000-000000000000",
					},
					"capacity_id": "00000000-0000-0000-0000-000000000000",
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute number_of_member_gateways`),
		},
		// error - missing required attributes - virtual_network_azure_resource
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"type":                            string(fabcore.GatewayTypeVirtualNetwork),
					"display_name":                    "test",
					"inactivity_minutes_before_sleep": (int)(gateway.PossibleInactivityMinutesBeforeSleepValues[0]),
					"number_of_member_gateways":       (int)(gateway.MinNumberOfMemberGatewaysValues),
					"capacity_id":                     "00000000-0000-0000-0000-000000000000",
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute virtual_network_azure_resource`),
		},
		// error - missing required attributes - capacity_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"type":                            string(fabcore.GatewayTypeVirtualNetwork),
					"display_name":                    "test",
					"inactivity_minutes_before_sleep": (int)(gateway.PossibleInactivityMinutesBeforeSleepValues[0]),
					"number_of_member_gateways":       (int)(gateway.MinNumberOfMemberGatewaysValues),
					"virtual_network_azure_resource": map[string]any{
						"resource_group_name":  "test",
						"virtual_network_name": "test",
						"subnet_name":          "test",
						"subscription_id":      "00000000-0000-0000-0000-000000000000",
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid configuration for attribute capacity_id`),
		},
		// error - invalid attribute value - inactivity_minutes_before_sleep
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"type":                            string(fabcore.GatewayTypeVirtualNetwork),
					"display_name":                    "test",
					"inactivity_minutes_before_sleep": (int)(gateway.PossibleInactivityMinutesBeforeSleepValues[0]) - 1,
					"number_of_member_gateways":       (int)(gateway.MinNumberOfMemberGatewaysValues),
					"virtual_network_azure_resource": map[string]any{
						"resource_group_name":  "test",
						"virtual_network_name": "test",
						"subnet_name":          "test",
						"subscription_id":      "00000000-0000-0000-0000-000000000000",
					},
					"capacity_id": "00000000-0000-0000-0000-000000000000",
				},
			),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
		// error - invalid attribute value - number_of_member_gateways
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"type":                            string(fabcore.GatewayTypeVirtualNetwork),
					"display_name":                    "test",
					"inactivity_minutes_before_sleep": (int)(gateway.PossibleInactivityMinutesBeforeSleepValues[0]),
					"number_of_member_gateways":       (int)(gateway.MinNumberOfMemberGatewaysValues) - 1,
					"virtual_network_azure_resource": map[string]any{
						"resource_group_name":  "test",
						"virtual_network_name": "test",
						"subnet_name":          "test",
						"subscription_id":      "00000000-0000-0000-0000-000000000000",
					},
					"capacity_id": "00000000-0000-0000-0000-000000000000",
				},
			),
			ExpectError: regexp.MustCompile(fmt.Sprintf(`Attribute number_of_member_gateways value must be between %d and %d`, gateway.MinNumberOfMemberGatewaysValues, gateway.MaxNumberOfMemberGatewaysValues)),
		},
		// error - invalid uuid - capacity_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"type":                            string(fabcore.GatewayTypeVirtualNetwork),
					"display_name":                    "test",
					"inactivity_minutes_before_sleep": (int)(gateway.PossibleInactivityMinutesBeforeSleepValues[0]),
					"number_of_member_gateways":       (int)(gateway.MinNumberOfMemberGatewaysValues),
					"virtual_network_azure_resource": map[string]any{
						"resource_group_name":  "test",
						"virtual_network_name": "test",
						"subnet_name":          "test",
						"subscription_id":      "00000000-0000-0000-0000-000000000000",
					},
					"capacity_id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unsupported gateway type - OnPremises
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"type": string(fabcore.GatewayTypeOnPremises),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
		},
		// error - unsupported gateway type - OnPremisesPersonal
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"type": string(fabcore.GatewayTypeOnPremisesPersonal),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorAttValueMatch),
		},
	}))
}

func TestUnit_GatewayResource_ImportState(t *testing.T) {
	entity := fakes.NewRandomVirtualNetworkGateway()

	fakes.FakeServer.Upsert(fakes.NewRandomVirtualNetworkGateway())
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomVirtualNetworkGateway())

	testCase := at.CompileConfig(
		testResourceItemHeader,
		map[string]any{
			"type":                            (string)(*entity.Type),
			"display_name":                    *entity.DisplayName,
			"inactivity_minutes_before_sleep": (int)(*entity.InactivityMinutesBeforeSleep),
			"number_of_member_gateways":       (int)(*entity.NumberOfMemberGateways),
			"virtual_network_azure_resource": map[string]any{
				"resource_group_name":  *entity.VirtualNetworkAzureResource.ResourceGroupName,
				"virtual_network_name": *entity.VirtualNetworkAzureResource.VirtualNetworkName,
				"subnet_name":          *entity.VirtualNetworkAzureResource.SubnetName,
				"subscription_id":      *entity.VirtualNetworkAzureResource.SubscriptionID,
			},
			"capacity_id": *entity.CapacityID,
		},
	)

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceItemFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// Import state testing
		{
			ResourceName:       testResourceItemFQN,
			Config:             testCase,
			ImportStateId:      *entity.ID,
			ImportState:        true,
			ImportStatePersist: true,
			ImportStateCheck: func(is []*terraform.InstanceState) error {
				if len(is) != 1 {
					return errors.New("expected one instance state")
				}

				if is[0].ID != *entity.ID {
					return errors.New(testResourceItemFQN + ": unexpected ID")
				}

				return nil
			},
		},
	}))
}

func TestUnit_GatewayResource_CRUD(t *testing.T) {
	entityExist := fakes.NewRandomVirtualNetworkGateway()
	entityBefore := fakes.NewRandomVirtualNetworkGateway()
	entityAfter := fakes.NewRandomVirtualNetworkGateway()

	fakes.FakeServer.Upsert(fakes.NewRandomVirtualNetworkGateway())
	fakes.FakeServer.Upsert(entityExist)
	fakes.FakeServer.Upsert(fakes.NewRandomVirtualNetworkGateway())

	resource.Test(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - create - existing entity
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"type":                            string(fabcore.GatewayTypeVirtualNetwork),
					"display_name":                    *entityExist.DisplayName,
					"inactivity_minutes_before_sleep": (int)(gateway.PossibleInactivityMinutesBeforeSleepValues[0]),
					"number_of_member_gateways":       (int)(gateway.MinNumberOfMemberGatewaysValues),
					"virtual_network_azure_resource": map[string]any{
						"resource_group_name":  "test",
						"virtual_network_name": "test",
						"subnet_name":          "test",
						"subscription_id":      "00000000-0000-0000-0000-000000000000",
					},
					"capacity_id": "00000000-0000-0000-0000-000000000000",
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorCreateHeader),
		},
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"type":                            (string)(*entityBefore.Type),
					"display_name":                    *entityBefore.DisplayName,
					"inactivity_minutes_before_sleep": (int)(*entityBefore.InactivityMinutesBeforeSleep),
					"number_of_member_gateways":       (int)(*entityBefore.NumberOfMemberGateways),
					"virtual_network_azure_resource": map[string]any{
						"resource_group_name":  *entityBefore.VirtualNetworkAzureResource.ResourceGroupName,
						"virtual_network_name": *entityBefore.VirtualNetworkAzureResource.VirtualNetworkName,
						"subnet_name":          *entityBefore.VirtualNetworkAzureResource.SubnetName,
						"subscription_id":      *entityBefore.VirtualNetworkAzureResource.SubscriptionID,
					},
					"capacity_id": *entityBefore.CapacityID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "type", (*string)(entityBefore.Type)),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityBefore.DisplayName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "capacity_id", entityBefore.CapacityID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "inactivity_minutes_before_sleep", strconv.Itoa(int(*entityBefore.InactivityMinutesBeforeSleep))),
				resource.TestCheckResourceAttr(testResourceItemFQN, "number_of_member_gateways", strconv.Itoa(int(*entityBefore.NumberOfMemberGateways))),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "virtual_network_azure_resource.resource_group_name", entityBefore.VirtualNetworkAzureResource.ResourceGroupName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "virtual_network_azure_resource.subnet_name", entityBefore.VirtualNetworkAzureResource.SubnetName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "virtual_network_azure_resource.subscription_id", entityBefore.VirtualNetworkAzureResource.SubscriptionID),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "virtual_network_azure_resource.virtual_network_name", entityBefore.VirtualNetworkAzureResource.VirtualNetworkName),
			),
		},
		// Update and Read - no replacement
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"type":                            (string)(*entityBefore.Type),
					"display_name":                    *entityAfter.DisplayName,
					"inactivity_minutes_before_sleep": (int)(*entityAfter.InactivityMinutesBeforeSleep),
					"number_of_member_gateways":       (int)(*entityAfter.NumberOfMemberGateways),
					"virtual_network_azure_resource": map[string]any{
						"resource_group_name":  *entityBefore.VirtualNetworkAzureResource.ResourceGroupName,
						"virtual_network_name": *entityBefore.VirtualNetworkAzureResource.VirtualNetworkName,
						"subnet_name":          *entityBefore.VirtualNetworkAzureResource.SubnetName,
						"subscription_id":      *entityBefore.VirtualNetworkAzureResource.SubscriptionID,
					},
					"capacity_id": *entityAfter.CapacityID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "type", string(*entityBefore.Type)),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "inactivity_minutes_before_sleep", strconv.Itoa(int(*entityAfter.InactivityMinutesBeforeSleep))),
				resource.TestCheckResourceAttr(testResourceItemFQN, "number_of_member_gateways", strconv.Itoa(int(*entityAfter.NumberOfMemberGateways))),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "capacity_id", entityAfter.CapacityID),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "virtual_network_azure_resource.resource_group_name", entityBefore.VirtualNetworkAzureResource.ResourceGroupName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "virtual_network_azure_resource.subnet_name", entityBefore.VirtualNetworkAzureResource.SubnetName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "virtual_network_azure_resource.subscription_id", entityBefore.VirtualNetworkAzureResource.SubscriptionID),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "virtual_network_azure_resource.virtual_network_name", entityBefore.VirtualNetworkAzureResource.VirtualNetworkName),
			),
		},
		// Update and Read - with replacement
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"type":                            (string)(*entityBefore.Type),
					"display_name":                    *entityAfter.DisplayName,
					"inactivity_minutes_before_sleep": (int)(*entityAfter.InactivityMinutesBeforeSleep),
					"number_of_member_gateways":       (int)(*entityAfter.NumberOfMemberGateways),
					"virtual_network_azure_resource": map[string]any{
						"resource_group_name":  *entityAfter.VirtualNetworkAzureResource.ResourceGroupName,
						"virtual_network_name": *entityAfter.VirtualNetworkAzureResource.VirtualNetworkName,
						"subnet_name":          *entityAfter.VirtualNetworkAzureResource.SubnetName,
						"subscription_id":      *entityAfter.VirtualNetworkAzureResource.SubscriptionID,
					},
					"capacity_id": *entityAfter.CapacityID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "type", string(*entityBefore.Type)),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "display_name", entityAfter.DisplayName),
				resource.TestCheckResourceAttr(testResourceItemFQN, "inactivity_minutes_before_sleep", strconv.Itoa(int(*entityAfter.InactivityMinutesBeforeSleep))),
				resource.TestCheckResourceAttr(testResourceItemFQN, "number_of_member_gateways", strconv.Itoa(int(*entityAfter.NumberOfMemberGateways))),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "capacity_id", entityAfter.CapacityID),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "virtual_network_azure_resource.resource_group_name", entityAfter.VirtualNetworkAzureResource.ResourceGroupName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "virtual_network_azure_resource.subnet_name", entityAfter.VirtualNetworkAzureResource.SubnetName),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "virtual_network_azure_resource.subscription_id", entityAfter.VirtualNetworkAzureResource.SubscriptionID),
				resource.TestCheckResourceAttrPtr(testResourceItemFQN, "virtual_network_azure_resource.virtual_network_name", entityAfter.VirtualNetworkAzureResource.VirtualNetworkName),
			),
		},
		// Delete testing automatically occurs in TestCase
	}))
}

// func TestAcc_GatewayResource_CRUD(t *testing.T) {
// 	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
// 	capacityID := capacity["id"].(string)

// 	entityCreateDisplayName := testhelp.RandomName()
// 	entityUpdateDisplayName := testhelp.RandomName()
// 	entityUpdateDescription := testhelp.RandomName()

// 	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
// 		// Create and Read
// 		{
// 			ResourceName: testResourceItemFQN,
// 			Config: at.CompileConfig(
// 				testResourceItemHeader,
// 				map[string]any{
// 					"display_name": entityCreateDisplayName,
// 					"capacity_id":  capacityID,
// 				},
// 			),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "capacity_id", capacityID),
// 			),
// 		},
// 		// Update and Read
// 		{
// 			ResourceName: testResourceItemFQN,
// 			Config: at.CompileConfig(
// 				testResourceItemHeader,
// 				map[string]any{
// 					"display_name": entityUpdateDisplayName,
// 					"description":  entityUpdateDescription,
// 					"capacity_id":  capacityID,
// 				},
// 			),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "capacity_id", capacityID),
// 			),
// 		},
// 		// Update - unassign capacity
// 		{
// 			ResourceName: testResourceItemFQN,
// 			Config: at.CompileConfig(
// 				testResourceItemHeader,
// 				map[string]any{
// 					"display_name": entityUpdateDisplayName,
// 					"description":  entityUpdateDescription,
// 				},
// 			),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
// 				resource.TestCheckNoResourceAttr(testResourceItemFQN, "capacity_id"),
// 			),
// 		},
// 		// Update - assign capacity
// 		{
// 			ResourceName: testResourceItemFQN,
// 			Config: at.CompileConfig(
// 				testResourceItemHeader,
// 				map[string]any{
// 					"display_name": entityUpdateDisplayName,
// 					"description":  entityUpdateDescription,
// 					"capacity_id":  capacityID,
// 				},
// 			),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "capacity_id", capacityID),
// 			),
// 		},
// 		// Update - assign identity
// 		{
// 			ResourceName: testResourceItemFQN,
// 			Config: at.CompileConfig(
// 				testResourceItemHeader,
// 				map[string]any{
// 					"display_name": entityUpdateDisplayName,
// 					"description":  entityUpdateDescription,
// 					"capacity_id":  capacityID,
// 					"identity": map[string]any{
// 						"type": "SystemAssigned",
// 					},
// 				},
// 			),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "identity.application_id"),
// 			),
// 		},
// 		// Update - unassign identity
// 		{
// 			ResourceName: testResourceItemFQN,
// 			Config: at.CompileConfig(
// 				testResourceItemHeader,
// 				map[string]any{
// 					"display_name": entityUpdateDisplayName,
// 					"description":  entityUpdateDescription,
// 					"capacity_id":  capacityID,
// 				},
// 			),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
// 				resource.TestCheckNoResourceAttr(testResourceItemFQN, "identity"),
// 			),
// 		},
// 	},
// 	))
// }

// func TestAcc_GatewayResource_Identity_CRUD(t *testing.T) {
// 	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
// 	capacityID := capacity["id"].(string)

// 	entityCreateDisplayName := testhelp.RandomName()
// 	entityUpdateDisplayName := testhelp.RandomName()
// 	entityUpdateDescription := testhelp.RandomName()

// 	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
// 		// Create and Read
// 		{
// 			ResourceName: testResourceItemFQN,
// 			Config: at.CompileConfig(
// 				testResourceItemHeader,
// 				map[string]any{
// 					"display_name": entityCreateDisplayName,
// 					"capacity_id":  capacityID,
// 					"identity": map[string]any{
// 						"type": "SystemAssigned",
// 					},
// 				},
// 			),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityCreateDisplayName),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "description", ""),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "capacity_id", capacityID),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "identity.application_id"),
// 			),
// 		},
// 		// Update and Read
// 		{
// 			ResourceName: testResourceItemFQN,
// 			Config: at.CompileConfig(
// 				testResourceItemHeader,
// 				map[string]any{
// 					"display_name": entityUpdateDisplayName,
// 					"description":  entityUpdateDescription,
// 					"capacity_id":  capacityID,
// 				},
// 			),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "capacity_id", capacityID),
// 			),
// 		},
// 		// Update - unassign capacity
// 		{
// 			ResourceName: testResourceItemFQN,
// 			Config: at.CompileConfig(
// 				testResourceItemHeader,
// 				map[string]any{
// 					"display_name": entityUpdateDisplayName,
// 					"description":  entityUpdateDescription,
// 				},
// 			),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
// 				resource.TestCheckNoResourceAttr(testResourceItemFQN, "capacity_id"),
// 			),
// 		},
// 		// Update - assign capacity
// 		{
// 			ResourceName: testResourceItemFQN,
// 			Config: at.CompileConfig(
// 				testResourceItemHeader,
// 				map[string]any{
// 					"display_name": entityUpdateDisplayName,
// 					"description":  entityUpdateDescription,
// 					"capacity_id":  capacityID,
// 				},
// 			),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "capacity_id", capacityID),
// 			),
// 		},
// 		// Update - unassign identity
// 		{
// 			ResourceName: testResourceItemFQN,
// 			Config: at.CompileConfig(
// 				testResourceItemHeader,
// 				map[string]any{
// 					"display_name": entityUpdateDisplayName,
// 					"description":  entityUpdateDescription,
// 					"capacity_id":  capacityID,
// 				},
// 			),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
// 				resource.TestCheckNoResourceAttr(testResourceItemFQN, "identity"),
// 			),
// 		},
// 		// Update - assign identity
// 		{
// 			ResourceName: testResourceItemFQN,
// 			Config: at.CompileConfig(
// 				testResourceItemHeader,
// 				map[string]any{
// 					"display_name": entityUpdateDisplayName,
// 					"description":  entityUpdateDescription,
// 					"capacity_id":  capacityID,
// 					"identity": map[string]any{
// 						"type": "SystemAssigned",
// 					},
// 				},
// 			),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "display_name", entityUpdateDisplayName),
// 				resource.TestCheckResourceAttr(testResourceItemFQN, "description", entityUpdateDescription),
// 				resource.TestCheckResourceAttrSet(testResourceItemFQN, "identity.application_id"),
// 			),
// 		},
// 	},
// 	))
// }
