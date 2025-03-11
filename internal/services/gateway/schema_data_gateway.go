// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getDataSourceGatewayAttributes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The " + ItemName + " ID.",
			Optional:            true,
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "The " + ItemName + " display name.",
			Optional:            true,
			Computed:            true,
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "The " + ItemName + " type. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleGatewayTypeValues(), true, true),
			Computed:            true,
		},
		"capacity_id": schema.StringAttribute{
			MarkdownDescription: "The " + ItemName + " capacity ID.",
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
		},
		"inactivity_minutes_before_sleep": schema.Int32Attribute{
			MarkdownDescription: "The " + ItemName + " inactivity minutes before sleep. Possible values: " + utils.ConvertStringSlicesToString(PossibleInactivityMinutesBeforeSleepValues, true, true),
			Computed:            true,
		},
		"number_of_member_gateways": schema.Int32Attribute{
			MarkdownDescription: fmt.Sprintf("The number of member gateways. Possible values: %d to %d.", MinNumberOfMemberGatewaysValues, MaxNumberOfMemberGatewaysValues),
			Computed:            true,
		},
		"virtual_network_azure_resource": schema.SingleNestedAttribute{
			MarkdownDescription: "The Azure virtual network resource.",
			Computed:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[virtualNetworkAzureResourceModel](ctx),
			Attributes: map[string]schema.Attribute{
				"resource_group_name": schema.StringAttribute{
					MarkdownDescription: "The resource group name.",
					Computed:            true,
				},
				"subnet_name": schema.StringAttribute{
					MarkdownDescription: "The subnet name.",
					Computed:            true,
				},
				"subscription_id": schema.StringAttribute{
					MarkdownDescription: "The subscription ID.",
					Computed:            true,
					CustomType:          customtypes.UUIDType{},
				},
				"virtual_network_name": schema.StringAttribute{
					MarkdownDescription: "The virtual network name.",
					Computed:            true,
				},
			},
		},
		"allow_cloud_connection_refresh": schema.BoolAttribute{
			MarkdownDescription: "Allow cloud connection refresh.",
			Computed:            true,
		},
		"allow_custom_connectors": schema.BoolAttribute{
			MarkdownDescription: "Allow custom connectors.",
			Computed:            true,
		},
		"load_balancing_setting": schema.StringAttribute{
			MarkdownDescription: "The load balancing setting. Possible values: " + utils.ConvertStringSlicesToString(fabcore.PossibleLoadBalancingSettingValues(), true, true),
			Computed:            true,
		},
		"public_key": schema.SingleNestedAttribute{
			MarkdownDescription: "The public key of the primary gateway member. Used to encrypt the credentials for creating and updating connections.",
			Computed:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[publicKeyModel](ctx),
			Attributes: map[string]schema.Attribute{
				"exponent": schema.StringAttribute{
					MarkdownDescription: "The exponent.",
					Computed:            true,
				},
				"modulus": schema.StringAttribute{
					MarkdownDescription: "The modulus.",
					Computed:            true,
				},
			},
		},
		"version": schema.StringAttribute{
			MarkdownDescription: "The " + ItemName + " version.",
			Computed:            true,
		},
	}
}
