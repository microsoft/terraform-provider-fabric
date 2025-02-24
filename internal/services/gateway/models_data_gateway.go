// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceGatewayModel struct {
	baseDataSourceGatewayModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type baseDataSourceGatewayModel struct {
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

func (to *baseDataSourceGatewayModel) set(ctx context.Context, from fabcore.GatewayClassification) diag.Diagnostics {
	var diags diag.Diagnostics

	virtualNetworkAzureResource := supertypes.NewSingleNestedObjectValueOfNull[virtualNetworkAzureResourceModel](ctx)
	publicKey := supertypes.NewSingleNestedObjectValueOfNull[publicKeyModel](ctx)

	gw := from.GetGateway()
	to.ID = customtypes.NewUUIDPointerValue(gw.ID)
	to.Type = types.StringPointerValue((*string)(gw.Type))

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

		to.PublicKey = publicKey

	case *fabcore.OnPremisesGateway:
		to.DisplayName = types.StringPointerValue(entity.DisplayName)
		to.NumberOfMemberGateways = types.Int32PointerValue(entity.NumberOfMemberGateways)
		to.AllowCloudConnectionRefresh = types.BoolPointerValue(entity.AllowCloudConnectionRefresh)
		to.AllowCustomConnectors = types.BoolPointerValue(entity.AllowCustomConnectors)
		to.LoadBalancingSetting = types.StringPointerValue((*string)(entity.LoadBalancingSetting))
		to.Version = types.StringPointerValue(entity.Version)

		if entity.PublicKey != nil {
			publicKeyModel := &publicKeyModel{}
			publicKeyModel.set(*entity.PublicKey)

			if diags := publicKey.Set(ctx, publicKeyModel); diags.HasError() {
				return diags
			}
		}

		to.PublicKey = publicKey

		to.VirtualNetworkAzureResource = virtualNetworkAzureResource

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

		to.VirtualNetworkAzureResource = virtualNetworkAzureResource

	default:
		diags.AddError("Unsupported Gateway type", fmt.Sprintf("The Gateway type '%s' is not supported.", (string)(*gw.Type)))

		return diags
	}

	return nil
}

type publicKeyModel struct {
	Exponent types.String `tfsdk:"exponent"`
	Modulus  types.String `tfsdk:"modulus"`
}

func (to *publicKeyModel) set(from fabcore.PublicKey) {
	to.Exponent = types.StringPointerValue(from.Exponent)
	to.Modulus = types.StringPointerValue(from.Modulus)
}
