// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

type resourceGatewayModel struct {
	baseGatewayModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
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
		}
	default:
		diags.AddError("Unsupported Gateway type", fmt.Sprintf("The Gateway type '%T' is not supported.", gatewayType))

		return diags
	}

	return nil
}
