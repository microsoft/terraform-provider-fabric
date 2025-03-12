// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type virtualNetworkAzureResourceModel struct {
	ResourceGroupName  types.String     `tfsdk:"resource_group_name"`
	SubnetName         types.String     `tfsdk:"subnet_name"`
	SubscriptionID     customtypes.UUID `tfsdk:"subscription_id"`
	VirtualNetworkName types.String     `tfsdk:"virtual_network_name"` // Rename to just 'name' or 'display_name'?
}

func (to *virtualNetworkAzureResourceModel) set(from fabcore.VirtualNetworkAzureResource) {
	to.ResourceGroupName = types.StringPointerValue(from.ResourceGroupName)
	to.SubnetName = types.StringPointerValue(from.SubnetName)
	to.SubscriptionID = customtypes.NewUUIDPointerValue(from.SubscriptionID)
	to.VirtualNetworkName = types.StringPointerValue(from.VirtualNetworkName)
}

type baseGatewayRoleAssignmentModel struct {
	ID        customtypes.UUID                                                          `tfsdk:"id"`
	Role      types.String                                                              `tfsdk:"role"`
	Principal supertypes.SingleNestedObjectValueOf[gatewayRoleAssignmentPrincipalModel] `tfsdk:"principal"`
}

func (to *baseGatewayRoleAssignmentModel) set(ctx context.Context, from fabcore.GatewayRoleAssignment) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Role = types.StringPointerValue((*string)(from.Role))

	principal := supertypes.NewSingleNestedObjectValueOfNull[gatewayRoleAssignmentPrincipalModel](ctx)

	if from.Principal != nil {
		principalModel := &gatewayRoleAssignmentPrincipalModel{}

		principalModel.set(*from.Principal)

		if diags := principal.Set(ctx, principalModel); diags.HasError() {
			return diags
		}
	}

	to.Principal = principal

	return nil
}

type gatewayRoleAssignmentPrincipalModel struct {
	ID   customtypes.UUID `tfsdk:"id"`
	Type types.String     `tfsdk:"type"`
}

func (to *gatewayRoleAssignmentPrincipalModel) set(from fabcore.Principal) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Type = types.StringPointerValue((*string)(from.Type))
}
