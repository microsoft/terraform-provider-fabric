// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type resourceGatewayModel struct {
	baseResourceGatewayModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type baseResourceGatewayModel struct {
	ID   customtypes.UUID `tfsdk:"id"`
	Type types.String     `tfsdk:"type"`

	DisplayName                  types.String                                                           `tfsdk:"display_name"`                    // VirtualNetwork
	CapacityID                   customtypes.UUID                                                       `tfsdk:"capacity_id"`                     // VirtualNetwork
	InactivityMinutesBeforeSleep types.Int32                                                            `tfsdk:"inactivity_minutes_before_sleep"` // VirtualNetwork
	NumberOfMemberGateways       types.Int32                                                            `tfsdk:"number_of_member_gateways"`       // VirtualNetwork
	VirtualNetworkAzureResource  supertypes.SingleNestedObjectValueOf[virtualNetworkAzureResourceModel] `tfsdk:"virtual_network_azure_resource"`  // VirtualNetwork
}

func (to *baseResourceGatewayModel) set(ctx context.Context, from fabcore.GatewayClassification) diag.Diagnostics {
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
			virtualNetworkAzureResourceModel.set(*gateway.VirtualNetworkAzureResource)

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

type requestCreateGateway struct {
	fabcore.CreateGatewayRequestClassification
}

func (to *requestCreateGateway) set(ctx context.Context, from resourceGatewayModel) diag.Diagnostics {
	var diags diag.Diagnostics

	gatewayType := (fabcore.GatewayType)(from.Type.ValueString())

	switch gatewayType {
	case fabcore.GatewayTypeVirtualNetwork:
		virtualNetworkAzureResource, diags := from.VirtualNetworkAzureResource.Get(ctx)
		if diags.HasError() {
			return diags
		}

		to.CreateGatewayRequestClassification = &fabcore.CreateVirtualNetworkGatewayRequest{
			Type:                         &gatewayType,
			CapacityID:                   from.CapacityID.ValueStringPointer(),
			DisplayName:                  from.DisplayName.ValueStringPointer(),
			InactivityMinutesBeforeSleep: from.InactivityMinutesBeforeSleep.ValueInt32Pointer(),
			NumberOfMemberGateways:       from.NumberOfMemberGateways.ValueInt32Pointer(),
			VirtualNetworkAzureResource: &fabcore.VirtualNetworkAzureResource{
				SubscriptionID:     virtualNetworkAzureResource.SubscriptionID.ValueStringPointer(),
				ResourceGroupName:  virtualNetworkAzureResource.ResourceGroupName.ValueStringPointer(),
				VirtualNetworkName: virtualNetworkAzureResource.VirtualNetworkName.ValueStringPointer(),
				SubnetName:         virtualNetworkAzureResource.SubnetName.ValueStringPointer(),
			},
		}
	default:
		diags.AddError("Unsupported Gateway type", fmt.Sprintf("The Gateway type '%T' is not supported.", gatewayType))

		return diags
	}

	return nil
}

type requestUpdateGateway struct {
	fabcore.UpdateGatewayRequestClassification
}

func (to *requestUpdateGateway) set(from resourceGatewayModel) diag.Diagnostics {
	var diags diag.Diagnostics

	gatewayType := (fabcore.GatewayType)(from.Type.ValueString())

	switch gatewayType {
	case fabcore.GatewayTypeVirtualNetwork:
		to.UpdateGatewayRequestClassification = &fabcore.UpdateVirtualNetworkGatewayRequest{
			Type:                         &gatewayType,
			DisplayName:                  from.DisplayName.ValueStringPointer(),
			CapacityID:                   from.CapacityID.ValueStringPointer(),
			InactivityMinutesBeforeSleep: from.InactivityMinutesBeforeSleep.ValueInt32Pointer(),
			NumberOfMemberGateways:       from.NumberOfMemberGateways.ValueInt32Pointer(),
		}
	default:
		diags.AddError("Unsupported Gateway type", fmt.Sprintf("The Gateway type '%T' is not supported.", gatewayType))

		return diags
	}

	return nil
}
