// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"

	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type onPremisesGatewayModelBase struct {
	ID customtypes.UUID `tfsdk:"id"`

	DisplayName types.String `tfsdk:"display_name"`

	AllowCloudConnectionRefresh types.Bool `tfsdk:"allow_cloud_connection_refresh"`

	AllowCustomConnectors types.Bool `tfsdk:"allow_custom_connectors"`

	LoadBalancingSetting types.String `tfsdk:"load_balancing_setting"`

	NumberOfMemberGateways types.Int32 `tfsdk:"number_of_member_gateways"`

	PublicKey supertypes.SingleNestedObjectValueOf[publicKeyModel] `tfsdk:"public_key"`

	Version types.String `tfsdk:"version"`
}

type virtualNetworkGatewayModelBase struct {
	ID customtypes.UUID `tfsdk:"id"`

	DisplayName types.String `tfsdk:"display_name"`

	CapacityId customtypes.UUID `tfsdk:"capacity_id"`

	InactivityMinutesBeforeSleep types.Int32 `tfsdk:"inactivity_minutes_before_sleep"`

	NumberOfMemberGateways types.Int32 `tfsdk:"number_of_member_gateways"`

	VirtualNetworkAzureResource supertypes.SingleNestedObjectValueOf[virtualNetworkAzureResourceModel] `tfsdk:"virtual_network_azure_resource"`
}

type onPremisesGatewayPersonalModelBase struct {
	ID customtypes.UUID `tfsdk:"id"`

	PublicKey supertypes.SingleNestedObjectValueOf[publicKeyModel] `tfsdk:"public_key"`

	Version types.String `tfsdk:"version"`
}

type virtualNetworkAzureResourceModel struct {
	SubscriptionID customtypes.UUID `tfsdk:"subscription_id"`

	ResourceGroupName types.String `tfsdk:"resource_group_name"`

	VirtualNetworkName types.String `tfsdk:"virtual_network_name"`

	SubnetName types.String `tfsdk:"subnet_name"`
}

type publicKeyModel struct {
	Exponent types.String `tfsdk:"exponent"`

	Modulus types.String `tfsdk:"modulus"`
}

func (to *onPremisesGatewayPersonalModelBase) set(ctx context.Context, from fabcore.OnPremisesGatewayPersonal) diag.Diagnostics {
	var diags diag.Diagnostics
	to.ID = customtypes.NewUUIDPointerValue(from.ID)

	publicKey := supertypes.NewSingleNestedObjectValueOfNull[publicKeyModel](ctx)
	if from.PublicKey != nil {
		publicKeyModel := &publicKeyModel{}
		publicKeyModel.set(*from.PublicKey)

		if pkDiags := publicKey.Set(ctx, publicKeyModel); pkDiags.HasError() {
			diags.Append(pkDiags...)
			return diags
		}
	}
	to.PublicKey = publicKey

	to.Version = types.StringPointerValue(from.Version)

	return diags
}

func (to *onPremisesGatewayModelBase) set(ctx context.Context, from fabcore.OnPremisesGateway) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.AllowCloudConnectionRefresh = types.BoolPointerValue(from.AllowCloudConnectionRefresh)
	to.AllowCustomConnectors = types.BoolPointerValue(from.AllowCustomConnectors)
	to.LoadBalancingSetting = types.StringPointerValue((*string)(from.LoadBalancingSetting))
	to.NumberOfMemberGateways = types.Int32PointerValue(from.NumberOfMemberGateways)
	to.Version = types.StringPointerValue(from.Version)

	publicKey := supertypes.NewSingleNestedObjectValueOfNull[publicKeyModel](ctx)

	if from.PublicKey != nil {
		publicKeyModel := &publicKeyModel{}
		publicKeyModel.set(*from.PublicKey)

		if diags := publicKey.Set(ctx, publicKeyModel); diags.HasError() {
			return diags
		}
	}

	to.PublicKey = publicKey

	return nil
}

func (to *virtualNetworkGatewayModelBase) set(ctx context.Context, from fabcore.VirtualNetworkGateway) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.CapacityId = customtypes.NewUUIDPointerValue(from.CapacityID)
	to.InactivityMinutesBeforeSleep = types.Int32PointerValue(from.InactivityMinutesBeforeSleep)
	to.NumberOfMemberGateways = types.Int32PointerValue(from.NumberOfMemberGateways)

	virtualNetworkAzureResource := supertypes.NewSingleNestedObjectValueOfNull[virtualNetworkAzureResourceModel](ctx)

	if from.VirtualNetworkAzureResource != nil {
		virtualNetworkAzureResourceModel := &virtualNetworkAzureResourceModel{}
		virtualNetworkAzureResourceModel.set(*from.VirtualNetworkAzureResource)

		if diags := virtualNetworkAzureResource.Set(ctx, virtualNetworkAzureResourceModel); diags.HasError() {
			return diags
		}
	}

	to.VirtualNetworkAzureResource = virtualNetworkAzureResource

	return nil
}

func (to *virtualNetworkAzureResourceModel) set(from fabcore.VirtualNetworkAzureResource) {
	to.SubscriptionID = customtypes.NewUUIDPointerValue(from.SubscriptionID)
	to.ResourceGroupName = types.StringPointerValue(from.ResourceGroupName)
	to.VirtualNetworkName = types.StringPointerValue(from.VirtualNetworkName)
	to.SubnetName = types.StringPointerValue(from.SubnetName)
}

func (to *publicKeyModel) set(from fabcore.PublicKey) {
	to.Exponent = types.StringPointerValue(from.Exponent)
	to.Modulus = types.StringPointerValue(from.Modulus)
}
