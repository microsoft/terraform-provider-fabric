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

	DisplayName                  types.String                                                           `tfsdk:"display_name"`                    // VirtualNetwork & OnPremises
	CapacityID                   customtypes.UUID                                                       `tfsdk:"capacity_id"`                     // VirtualNetwork
	InactivityMinutesBeforeSleep types.Int32                                                            `tfsdk:"inactivity_minutes_before_sleep"` // VirtualNetwork
	NumberOfMemberGateways       types.Int32                                                            `tfsdk:"number_of_member_gateways"`       // VirtualNetwork & OnPremises
	VirtualNetworkAzureResource  supertypes.SingleNestedObjectValueOf[virtualNetworkAzureResourceModel] `tfsdk:"virtual_network_azure_resource"`  // VirtualNetwork
	AllowCloudConnectionRefresh  types.Bool                                                             `tfsdk:"allow_cloud_connection_refresh"`  // OnPremises
	AllowCustomConnectors        types.Bool                                                             `tfsdk:"allow_custom_connectors"`         // OnPremises
	LoadBalancingSetting         types.String                                                           `tfsdk:"load_balancing_setting"`          // OnPremises
	PublicKey                    supertypes.SingleNestedObjectValueOf[publicKeyModel]                   `tfsdk:"public_key"`                      // OnPremises & OnPremisesPersonal
	Version                      types.String                                                           `tfsdk:"version"`                         // OnPremises & OnPremisesPersonal
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

	case *fabcore.OnPremisesGateway:
		to.ID = customtypes.NewUUIDPointerValue(gateway.ID)
		to.Type = types.StringPointerValue((*string)(gateway.Type))
		to.DisplayName = types.StringPointerValue(gateway.DisplayName)
		to.NumberOfMemberGateways = types.Int32PointerValue(gateway.NumberOfMemberGateways)
		to.AllowCloudConnectionRefresh = types.BoolPointerValue(gateway.AllowCloudConnectionRefresh)
		to.AllowCustomConnectors = types.BoolPointerValue(gateway.AllowCustomConnectors)
		to.LoadBalancingSetting = types.StringPointerValue((*string)(gateway.LoadBalancingSetting))
		to.Version = types.StringPointerValue(gateway.Version)

		publicKey := supertypes.NewSingleNestedObjectValueOfNull[publicKeyModel](ctx)
		if gateway.PublicKey != nil {
			publicKeyModel := &publicKeyModel{}
			publicKeyModel.set(gateway.PublicKey)

			if diags := publicKey.Set(ctx, publicKeyModel); diags.HasError() {
				return diags
			}
		}

		to.PublicKey = publicKey

	case *fabcore.OnPremisesGatewayPersonal:
		to.ID = customtypes.NewUUIDPointerValue(gateway.ID)
		to.Type = types.StringPointerValue((*string)(gateway.Type))
		to.Version = types.StringPointerValue(gateway.Version)

		publicKey := supertypes.NewSingleNestedObjectValueOfNull[publicKeyModel](ctx)
		if gateway.PublicKey != nil {
			publicKeyModel := &publicKeyModel{}
			publicKeyModel.set(gateway.PublicKey)

			if diags := publicKey.Set(ctx, publicKeyModel); diags.HasError() {
				return diags
			}
		}

		to.PublicKey = publicKey

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

type publicKeyModel struct {
	Exponent types.String `tfsdk:"exponent"`
	Modulus  types.String `tfsdk:"modulus"`
}

func (to *publicKeyModel) set(from *fabcore.PublicKey) {
	to.Exponent = types.StringPointerValue(from.Exponent)
	to.Modulus = types.StringPointerValue(from.Modulus)
}
