// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type baseGatewayModel struct {
	ID   customtypes.UUID `tfsdk:"id"`
	Type types.String     `tfsdk:"type"`

	DisplayName                  types.String                                                           `tfsdk:"display_name"`
	CapacityID                   customtypes.UUID                                                       `tfsdk:"capacity_id"`
	InactivityMinutesBeforeSleep types.Int32                                                            `tfsdk:"inactivity_minutes_before_sleep"`
	NumberOfMemberGateways       types.Int32                                                            `tfsdk:"number_of_member_gateways"`
	VirtualNetworkAzureResource  supertypes.SingleNestedObjectValueOf[virtualNetworkAzureResourceModel] `tfsdk:"virtual_network_azure_resource"`
}

func (to *baseGatewayModel) set(ctx context.Context, from fabcore.GatewayClassification) diag.Diagnostics {
	var diags diag.Diagnostics

	switch gateway := from.(type) {
	case *fabcore.VirtualNetworkGateway:
		to.ID = customtypes.NewUUIDPointerValue(gateway.ID)
		to.Type = types.StringPointerValue((*string)(gateway.Type))
		to.DisplayName = types.StringPointerValue(gateway.DisplayName)
		to.CapacityID = customtypes.NewUUIDPointerValue(gateway.CapacityID)
		to.InactivityMinutesBeforeSleep = types.Int32PointerValue(gateway.InactivityMinutesBeforeSleep)
		to.NumberOfMemberGateways = types.Int32PointerValue(gateway.NumberOfMemberGateways)

		virtualNetworkAzureResource := supertypes.NewSingleNestedObjectValueOfNull[virtualNetworkAzureResourceModel](ctx)
		if gateway.VirtualNetworkAzureResource != nil {
			virtualNetworkAzureResourceModel := &virtualNetworkAzureResourceModel{}
			virtualNetworkAzureResourceModel.set(gateway.VirtualNetworkAzureResource)

			if diags := virtualNetworkAzureResource.Set(ctx, virtualNetworkAzureResourceModel); diags.HasError() {
				return diags
			}
		}

		to.VirtualNetworkAzureResource = virtualNetworkAzureResource

	default:
		diags.AddError("Unsupported Gateway type", fmt.Sprintf("The Gateway type '%T' is not supported.", gateway))

		return diags
	}

	return nil
}

type virtualNetworkAzureResourceModel struct {
	ResourceGroupName  types.String     `tfsdk:"resource_group_name"`
	SubnetName         types.String     `tfsdk:"subnet_name"`
	SubscriptionID     customtypes.UUID `tfsdk:"subscription_id"`
	VirtualNetworkName types.String     `tfsdk:"virtual_network_name"` // Rename to just 'name' or 'display_name'?
}

func (to *virtualNetworkAzureResourceModel) set(from *fabcore.VirtualNetworkAzureResource) {
	to.ResourceGroupName = types.StringPointerValue(from.ResourceGroupName)
	to.SubnetName = types.StringPointerValue(from.SubnetName)
	to.SubscriptionID = customtypes.NewUUIDPointerValue(from.SubscriptionID)
	to.VirtualNetworkName = types.StringPointerValue(from.VirtualNetworkName)
}
