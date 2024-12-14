// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceWorkspaceRoleAssignmentsModel struct {
	WorkspaceID customtypes.UUID                                                 `tfsdk:"workspace_id"`
	Values      supertypes.ListNestedObjectValueOf[workspaceRoleAssignmentModel] `tfsdk:"values"`
	Timeouts    timeouts.Value                                                   `tfsdk:"timeouts"`
}

func (to *dataSourceWorkspaceRoleAssignmentsModel) setValues(ctx context.Context, from []fabcore.WorkspaceRoleAssignment) diag.Diagnostics {
	slice := make([]*workspaceRoleAssignmentModel, 0, len(from))

	for _, entity := range from {
		var entityModel workspaceRoleAssignmentModel

		if diags := entityModel.set(ctx, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

type workspaceRoleAssignmentModel struct {
	ID          customtypes.UUID                                            `tfsdk:"id"`
	Role        types.String                                                `tfsdk:"role"`
	DisplayName types.String                                                `tfsdk:"display_name"`
	Type        types.String                                                `tfsdk:"type"`
	Details     supertypes.SingleNestedObjectValueOf[principalDetailsModel] `tfsdk:"details"`
}

func (to *workspaceRoleAssignmentModel) set(ctx context.Context, from fabcore.WorkspaceRoleAssignment) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Role = types.StringPointerValue((*string)(from.Role))

	detailsModel := &principalDetailsModel{}
	detailsModel.set(from.Principal, to)

	if diags := to.Details.Set(ctx, detailsModel); diags.HasError() {
		return diags
	}

	return nil
}

type principalDetailsModel struct {
	UserPrincipalName types.String     `tfsdk:"user_principal_name"`
	GroupType         types.String     `tfsdk:"group_type"`
	AppID             customtypes.UUID `tfsdk:"app_id"`
	ParentPrincipalID customtypes.UUID `tfsdk:"parent_principal_id"`
}

func (to *principalDetailsModel) set(from *fabcore.Principal, roleAssignment *workspaceRoleAssignmentModel) {
	to.UserPrincipalName = types.StringNull()
	to.GroupType = types.StringNull()
	to.AppID = customtypes.NewUUIDNull()
	to.ParentPrincipalID = customtypes.NewUUIDNull()

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
