// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GatewayModelBase struct {
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
	Type        types.String     `tfsdk:"type"`
}

type OnPremisesGatewayModel struct {
	GatewayModelBase

	AllowCloudConnectionRefresh *bool `tfsdk:"allow_cloud_connection_refresh"`

	AllowCustomConnectors *bool `tfsdk:"allow_custom_connectors"`

	LoadBalancingSetting types.String `tfsdk:"load_balancing_setting"`

	NumberOfMemberGateways *int32 `tfsdk:"number_of_member_gateways"`

	PublicKey supertypes.SingleNestedObjectValueOf[PublicKeyModel] `tfsdk:"public_key"`

	Version types.String `tfsdk:"version"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type VirtualNetworkGatewayModel struct {
	GatewayModelBase

	CapacityId customtypes.UUID `tfsdk:"capacity_id"`

	InactivityMinutesBeforeSleep *int32 `tfsdk:"inactivity_minutes_before_sleep"`

	NumberOfMemberGateways *int32 `tfsdk:"number_of_member_gateways"`

	VirtualNetworkAzureResource supertypes.SingleNestedObjectValueOf[VirtualNetworkAzureResourceModel] `tfsdk:"virtual_network_azure_resource"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type OnPremisesGatewayPersonalModel struct {
	ID        customtypes.UUID                                     `tfsdk:"id"`
	PublicKey supertypes.SingleNestedObjectValueOf[PublicKeyModel] `tfsdk:"public_key"`
	Type      types.String                                         `tfsdk:"type"`
	Version   types.String                                         `tfsdk:"version"`
}

type VirtualNetworkAzureResourceModel struct {
	SubscriptionID customtypes.UUID `tfsdk:"subscription_id"`

	ResourceGroupName types.String `tfsdk:"resource_group_name"`

	VirtualNetworkName types.String `tfsdk:"virtual_network_name"`

	SubnetName types.String `tfsdk:"subnet_name"`
}

type PublicKeyModel struct {
	Exponent types.String `tfsdk:"exponent"`

	Modulus types.String `tfsdk:"modulus"`
}

func (to *OnPremisesGatewayPersonalModel) set(ctx context.Context, from fabcore.OnPremisesGateway) diag.Diagnostics {
	var diags diag.Diagnostics
	to.ID = customtypes.NewUUIDPointerValue(from.ID)

	publicKey := supertypes.NewSingleNestedObjectValueOfNull[PublicKeyModel](ctx)
	if from.PublicKey != nil {
		publicKeyModel := &PublicKeyModel{}
		publicKeyModel.set(*from.PublicKey)

		if pkDiags := publicKey.Set(ctx, publicKeyModel); pkDiags.HasError() {
			diags.Append(pkDiags...)
			return diags
		}
	}
	to.PublicKey = publicKey

	to.Type = types.StringPointerValue((*string)(from.Type))
	to.Version = types.StringPointerValue(from.Version)

	return diags
}

func (to *OnPremisesGatewayModel) set(ctx context.Context, from fabcore.OnPremisesGateway) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Type = types.StringPointerValue((*string)(from.Type))
	to.AllowCloudConnectionRefresh = from.AllowCloudConnectionRefresh
	to.AllowCustomConnectors = from.AllowCustomConnectors
	to.LoadBalancingSetting = types.StringPointerValue((*string)(from.LoadBalancingSetting))
	to.NumberOfMemberGateways = from.NumberOfMemberGateways
	to.Version = types.StringPointerValue(from.Version)

	publicKey := supertypes.NewSingleNestedObjectValueOfNull[PublicKeyModel](ctx)

	if from.PublicKey != nil {
		publicKeyModel := &PublicKeyModel{}
		publicKeyModel.set(*from.PublicKey)

		if diags := publicKey.Set(ctx, publicKeyModel); diags.HasError() {
			return diags
		}
	}

	to.PublicKey = publicKey

	return nil
}

func (to *VirtualNetworkGatewayModel) set(ctx context.Context, from fabcore.VirtualNetworkGateway) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Type = types.StringPointerValue((*string)(from.Type))
	to.CapacityId = customtypes.NewUUIDPointerValue(from.CapacityID)
	to.InactivityMinutesBeforeSleep = from.InactivityMinutesBeforeSleep
	to.NumberOfMemberGateways = from.NumberOfMemberGateways

	virtualNetworkAzureResource := supertypes.NewSingleNestedObjectValueOfNull[VirtualNetworkAzureResourceModel](ctx)

	if from.VirtualNetworkAzureResource != nil {
		virtualNetworkAzureResourceModel := &VirtualNetworkAzureResourceModel{}
		virtualNetworkAzureResourceModel.set(*from.VirtualNetworkAzureResource)

		if diags := virtualNetworkAzureResource.Set(ctx, virtualNetworkAzureResourceModel); diags.HasError() {
			return diags
		}
	}

	to.VirtualNetworkAzureResource = virtualNetworkAzureResource

	return nil
}

func (to *VirtualNetworkAzureResourceModel) set(from fabcore.VirtualNetworkAzureResource) {
	to.SubscriptionID = customtypes.NewUUIDPointerValue(from.SubscriptionID)
	to.ResourceGroupName = types.StringPointerValue(from.ResourceGroupName)
	to.VirtualNetworkName = types.StringPointerValue(from.VirtualNetworkName)
	to.SubnetName = types.StringPointerValue(from.SubnetName)
}

// should I change my blablabla?
// func (to *GatewayModelBase) set(from fabcore.Gateway) {
// 	to.ID = customtypes.NewUUIDPointerValue(from.ID)
// 	to.Type = types.StringPointerValue((*string)(from.Type))
// }

func (to *PublicKeyModel) set(from fabcore.PublicKey) {
	to.Exponent = types.StringPointerValue(from.Exponent)
	to.Modulus = types.StringPointerValue(from.Modulus)
}
