// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

// type dataSourceWorkspaceRoleAssignmentModel struct {
// 	WorkspaceID customtypes.UUID `tfsdk:"workspace_id"`
// 	baseWorkspaceRoleAssignmentModel
// 	Timeouts timeouts.Value `tfsdk:"timeouts"`
// }

type dataSourceWorkspaceRoleAssignmentsModel struct {
	WorkspaceID customtypes.UUID                                                     `tfsdk:"workspace_id"`
	Values      supertypes.ListNestedObjectValueOf[baseWorkspaceRoleAssignmentModel] `tfsdk:"values"`
	Timeouts    timeouts.Value                                                       `tfsdk:"timeouts"`
}

func (to *dataSourceWorkspaceRoleAssignmentsModel) setValues(ctx context.Context, from []fabcore.WorkspaceRoleAssignment) diag.Diagnostics {
	slice := make([]*baseWorkspaceRoleAssignmentModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseWorkspaceRoleAssignmentModel

		if diags := entityModel.set(ctx, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

type baseWorkspaceRoleAssignmentModel struct {
	ID                   customtypes.UUID                                            `tfsdk:"id"`
	PrincipalID          customtypes.UUID                                            `tfsdk:"principal_id"`
	Role                 types.String                                                `tfsdk:"role"`
	PrincipalDisplayName types.String                                                `tfsdk:"principal_display_name"`
	PrincipalType        types.String                                                `tfsdk:"principal_type"`
	PrincipalDetails     supertypes.SingleNestedObjectValueOf[principalDetailsModel] `tfsdk:"principal_details"`
}

func (to *baseWorkspaceRoleAssignmentModel) set(ctx context.Context, from fabcore.WorkspaceRoleAssignment) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Role = types.StringPointerValue((*string)(from.Role))

	to.PrincipalID = customtypes.NewUUIDNull()
	to.PrincipalDisplayName = types.StringNull()
	to.PrincipalType = types.StringNull()

	principalDetails := supertypes.NewSingleNestedObjectValueOfNull[principalDetailsModel](ctx)

	if from.Principal != nil {
		to.PrincipalID = customtypes.NewUUIDPointerValue(from.Principal.ID)
		to.PrincipalDisplayName = types.StringPointerValue(from.Principal.DisplayName)
		to.PrincipalType = types.StringPointerValue((*string)(from.Principal.Type))

		principalDetailsModel := &principalDetailsModel{}

		principalDetailsModel.set(*from.Principal)

		if diags := principalDetails.Set(ctx, principalDetailsModel); diags.HasError() {
			return diags
		}
	}

	to.PrincipalDetails = principalDetails

	return nil
}

type principalDetailsModel struct {
	UserPrincipalName types.String     `tfsdk:"user_principal_name"`
	GroupType         types.String     `tfsdk:"group_type"`
	AppID             customtypes.UUID `tfsdk:"app_id"`
	ParentPrincipalID customtypes.UUID `tfsdk:"parent_principal_id"`
}

func (to *principalDetailsModel) set(from fabcore.Principal) {
	to.UserPrincipalName = types.StringNull()
	to.GroupType = types.StringNull()
	to.AppID = customtypes.NewUUIDNull()
	to.ParentPrincipalID = customtypes.NewUUIDNull()

	switch *from.Type {
	case fabcore.PrincipalTypeUser:
		if from.UserDetails != nil {
			to.UserPrincipalName = types.StringPointerValue(from.UserDetails.UserPrincipalName)
		}
	case fabcore.PrincipalTypeGroup:
		if from.GroupDetails != nil {
			to.GroupType = types.StringPointerValue((*string)(from.GroupDetails.GroupType))
		}
	case fabcore.PrincipalTypeServicePrincipal:
		if from.ServicePrincipalDetails != nil {
			to.AppID = customtypes.NewUUIDPointerValue(from.ServicePrincipalDetails.AADAppID)
		}
	case fabcore.PrincipalTypeServicePrincipalProfile:
		if from.ServicePrincipalProfileDetails != nil {
			to.ParentPrincipalID = customtypes.NewUUIDPointerValue(from.ServicePrincipalProfileDetails.ParentPrincipal.ID)
		}
	}
}
