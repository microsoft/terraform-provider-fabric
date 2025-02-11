// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

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
