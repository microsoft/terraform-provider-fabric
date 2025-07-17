// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"
	"fmt"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseGatewayModel struct {
	ID                           customtypes.UUID                                                       `tfsdk:"id"`
	Type                         types.String                                                           `tfsdk:"type"`
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
	gw := from.GetGateway()
	to.ID = customtypes.NewUUIDPointerValue(gw.ID)
	to.Type = types.StringPointerValue((*string)(gw.Type))

	virtualNetworkAzureResource := supertypes.NewSingleNestedObjectValueOfNull[virtualNetworkAzureResourceModel](ctx)
	publicKey := supertypes.NewSingleNestedObjectValueOfNull[publicKeyModel](ctx)

	to.DisplayName = types.StringNull()
	to.CapacityID = customtypes.NewUUIDNull()
	to.InactivityMinutesBeforeSleep = types.Int32Null()
	to.NumberOfMemberGateways = types.Int32Null()

	to.AllowCloudConnectionRefresh = types.BoolNull()
	to.AllowCustomConnectors = types.BoolNull()
	to.LoadBalancingSetting = types.StringNull()

	to.Version = types.StringNull()
	to.VirtualNetworkAzureResource = virtualNetworkAzureResource
	to.PublicKey = publicKey

	switch entity := from.(type) {
	case *fabcore.VirtualNetworkGateway:
		to.DisplayName = types.StringPointerValue(entity.DisplayName)
		to.CapacityID = customtypes.NewUUIDPointerValue(entity.CapacityID)
		to.InactivityMinutesBeforeSleep = types.Int32PointerValue(entity.InactivityMinutesBeforeSleep)
		to.NumberOfMemberGateways = types.Int32PointerValue(entity.NumberOfMemberGateways)

		if entity.VirtualNetworkAzureResource != nil {
			virtualNetworkAzureResourceModel := &virtualNetworkAzureResourceModel{}
			virtualNetworkAzureResourceModel.set(*entity.VirtualNetworkAzureResource)

			if diags := virtualNetworkAzureResource.Set(ctx, virtualNetworkAzureResourceModel); diags.HasError() {
				return diags
			}
		}

		to.VirtualNetworkAzureResource = virtualNetworkAzureResource

	case *fabcore.OnPremisesGateway:
		to.DisplayName = types.StringPointerValue(entity.DisplayName)
		to.AllowCloudConnectionRefresh = types.BoolPointerValue(entity.AllowCloudConnectionRefresh)
		to.AllowCustomConnectors = types.BoolPointerValue(entity.AllowCustomConnectors)
		to.LoadBalancingSetting = types.StringPointerValue((*string)(entity.LoadBalancingSetting))
		to.NumberOfMemberGateways = types.Int32PointerValue(entity.NumberOfMemberGateways)
		to.Version = types.StringPointerValue(entity.Version)

		if entity.PublicKey != nil {
			publicKeyModel := &publicKeyModel{}
			publicKeyModel.set(*entity.PublicKey)

			if diags := publicKey.Set(ctx, publicKeyModel); diags.HasError() {
				return diags
			}
		}

		to.PublicKey = publicKey

	case *fabcore.OnPremisesGatewayPersonal:
		to.Version = types.StringPointerValue(entity.Version)

		if entity.PublicKey != nil {
			publicKeyModel := &publicKeyModel{}
			publicKeyModel.set(*entity.PublicKey)

			if diags := publicKey.Set(ctx, publicKeyModel); diags.HasError() {
				return diags
			}
		}

		to.PublicKey = publicKey

	default:
		var diags diag.Diagnostics

		diags.AddError(
			"Unsupported Gateway type",
			fmt.Sprintf("The Gateway type '%T' is not supported.", entity),
		)

		return diags
	}

	return nil
}

/*
DATA-SOURCE
*/

type dataSourceGatewayModel struct {
	baseGatewayModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceGatewaysModel struct {
	Values   supertypes.SetNestedObjectValueOf[baseGatewayModel] `tfsdk:"values"`
	Timeouts timeoutsD.Value                                     `tfsdk:"timeouts"`
}

func (to *dataSourceGatewaysModel) setValues(ctx context.Context, from []fabcore.GatewayClassification) diag.Diagnostics {
	slice := make([]*baseGatewayModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseGatewayModel

		if diags := entityModel.set(ctx, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

/*
RESOURCE
*/

type resourceGatewayModel struct {
	baseGatewayModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateGateway struct {
	fabcore.CreateGatewayRequestClassification
}

func (to *requestCreateGateway) set(ctx context.Context, from resourceGatewayModel) diag.Diagnostics {
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
		var diags diag.Diagnostics

		diags.AddError(
			"Unsupported Gateway type",
			fmt.Sprintf("The Gateway type '%T' is not supported.", gatewayType),
		)

		return diags
	}

	return nil
}

type requestUpdateGateway struct {
	fabcore.UpdateGatewayRequestClassification
}

func (to *requestUpdateGateway) set(from resourceGatewayModel) diag.Diagnostics {
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

	case fabcore.GatewayTypeOnPremises:
		to.UpdateGatewayRequestClassification = &fabcore.UpdateOnPremisesGatewayRequest{
			Type:                        &gatewayType,
			DisplayName:                 from.DisplayName.ValueStringPointer(),
			AllowCloudConnectionRefresh: from.AllowCloudConnectionRefresh.ValueBoolPointer(),
			AllowCustomConnectors:       from.AllowCustomConnectors.ValueBoolPointer(),
			LoadBalancingSetting:        (*fabcore.LoadBalancingSetting)(from.LoadBalancingSetting.ValueStringPointer()),
		}

	default:
		var diags diag.Diagnostics

		diags.AddError(
			"Unsupported Gateway type",
			fmt.Sprintf("The Gateway type '%T' is not supported.", gatewayType),
		)

		return diags
	}

	return nil
}

/*
HELPER MODELS
*/

type publicKeyModel struct {
	Exponent types.String `tfsdk:"exponent"`
	Modulus  types.String `tfsdk:"modulus"`
}

func (to *publicKeyModel) set(from fabcore.PublicKey) {
	to.Exponent = types.StringPointerValue(from.Exponent)
	to.Modulus = types.StringPointerValue(from.Modulus)
}

type virtualNetworkAzureResourceModel struct {
	ResourceGroupName  customtypes.CaseInsensitiveString `tfsdk:"resource_group_name"`
	SubnetName         types.String                      `tfsdk:"subnet_name"`
	SubscriptionID     customtypes.UUID                  `tfsdk:"subscription_id"`
	VirtualNetworkName types.String                      `tfsdk:"virtual_network_name"`
}

func (to *virtualNetworkAzureResourceModel) set(from fabcore.VirtualNetworkAzureResource) {
	to.ResourceGroupName = customtypes.NewCaseInsensitiveStringPointerValue(from.ResourceGroupName)
	to.SubnetName = types.StringPointerValue(from.SubnetName)
	to.SubscriptionID = customtypes.NewUUIDPointerValue(from.SubscriptionID)
	to.VirtualNetworkName = types.StringPointerValue(from.VirtualNetworkName)
}
