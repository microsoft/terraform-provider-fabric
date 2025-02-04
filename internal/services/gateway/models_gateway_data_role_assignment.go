// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceGatewayRoleAssignmentsModel struct {
	GatewayID customtypes.UUID                                               `tfsdk:"gateway_id"`
	Values    supertypes.ListNestedObjectValueOf[gatewayRoleAssignmentModel] `tfsdk:"values"`
	Timeouts  timeouts.Value                                                 `tfsdk:"timeouts"`
}

func (to *dataSourceGatewayRoleAssignmentsModel) setValues(ctx context.Context, from []fabcore.GatewayRoleAssignment) diag.Diagnostics {
	slice := make([]*gatewayRoleAssignmentModel, 0, len(from))

	for _, entity := range from {
		var entityModel gatewayRoleAssignmentModel

		if diags := entityModel.set(ctx, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

// maybe use common infra from WS
type gatewayRoleAssignmentModel struct {
	ID          customtypes.UUID                                            `tfsdk:"id"`
	Role        types.String                                                `tfsdk:"role"`
	DisplayName types.String                                                `tfsdk:"display_name"`
	Type        types.String                                                `tfsdk:"type"`
	Details     supertypes.SingleNestedObjectValueOf[principalDetailsModel] `tfsdk:"details"`
}

func (to *gatewayRoleAssignmentModel) set(ctx context.Context, from fabcore.GatewayRoleAssignment) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Role = types.StringPointerValue((*string)(from.Role))

	// Initialize the details model and set its values.
	detailsModel := &principalDetailsModel{}
	detailsModel.set(from.Principal, to)

	if diags := to.Details.Set(ctx, detailsModel); diags.HasError() {
		return diags
	}

	// Set common attributes from the principal.
	to.DisplayName = types.StringPointerValue(from.Principal.DisplayName)
	to.Type = types.StringPointerValue((*string)(from.Principal.Type))
	return nil
}

type principalDetailsModel struct {
	UserPrincipalName types.String     `tfsdk:"user_principal_name"`
	GroupType         types.String     `tfsdk:"group_type"`
	AppID             customtypes.UUID `tfsdk:"app_id"`
	ParentPrincipalID customtypes.UUID `tfsdk:"parent_principal_id"`
}

func (to *principalDetailsModel) set(from *fabcore.Principal, roleAssignment *gatewayRoleAssignmentModel) {
	to.UserPrincipalName = types.StringNull()
	to.GroupType = types.StringNull()
	to.AppID = customtypes.NewUUIDNull()
	to.ParentPrincipalID = customtypes.NewUUIDNull()

	// Set the DisplayName and Type on the role assignment (if not already set).
	roleAssignment.DisplayName = types.StringPointerValue(from.DisplayName)
	roleAssignment.Type = types.StringPointerValue((*string)(from.Type))

	switch *from.Type {
	case fabcore.PrincipalTypeUser:
		to.UserPrincipalName = types.StringPointerValue(from.UserDetails.UserPrincipalName)
	case fabcore.PrincipalTypeGroup:
		to.GroupType = types.StringPointerValue((*string)(from.GroupDetails.GroupType))
	case fabcore.PrincipalTypeServicePrincipal:
		to.AppID = customtypes.NewUUIDPointerValue(from.ServicePrincipalDetails.AADAppID)
	case fabcore.PrincipalTypeServicePrincipalProfile:
		to.ParentPrincipalID = customtypes.NewUUIDPointerValue(from.ServicePrincipalProfileDetails.ParentPrincipal.ID)
	}
}
