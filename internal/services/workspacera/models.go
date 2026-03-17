// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspacera

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseWorkspaceRoleAssignmentModel struct {
	ID          customtypes.UUID                                            `tfsdk:"id"`
	WorkspaceID customtypes.UUID                                            `tfsdk:"workspace_id"`
	Role        types.String                                                `tfsdk:"role"`
	Principal   supertypes.SingleNestedObjectValueOf[common.PrincipalModel] `tfsdk:"principal"`
}

func (to *baseWorkspaceRoleAssignmentModel) set(ctx context.Context, workspaceID string, from fabcore.WorkspaceRoleAssignment) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.Role = types.StringPointerValue((*string)(from.Role))

	principal := supertypes.NewSingleNestedObjectValueOfNull[common.PrincipalModel](ctx)

	if from.Principal != nil {
		principalModel := &common.PrincipalModel{}

		principalModel.Set(*from.Principal)

		if diags := principal.Set(ctx, principalModel); diags.HasError() {
			return diags
		}
	}

	to.Principal = principal

	return nil
}

/*
DATA-SOURCE
*/

type dataSourceWorkspaceRoleAssignmentModel struct {
	baseWorkspaceRoleAssignmentModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceWorkspaceRoleAssignmentsModel struct {
	WorkspaceID customtypes.UUID                                                    `tfsdk:"workspace_id"`
	Values      supertypes.SetNestedObjectValueOf[baseWorkspaceRoleAssignmentModel] `tfsdk:"values"`
	Timeouts    timeoutsD.Value                                                     `tfsdk:"timeouts"`
}

func (to *dataSourceWorkspaceRoleAssignmentsModel) setValues(ctx context.Context, workspaceID string, from []fabcore.WorkspaceRoleAssignment) diag.Diagnostics {
	slice := make([]*baseWorkspaceRoleAssignmentModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseWorkspaceRoleAssignmentModel

		if diags := entityModel.set(ctx, workspaceID, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

/*
RESOURCE
*/

type resourceWorkspaceRoleAssignmentModel struct {
	baseWorkspaceRoleAssignmentModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateWorkspaceRoleAssignment struct {
	fabcore.AddWorkspaceRoleAssignmentRequest
}

func (to *requestCreateWorkspaceRoleAssignment) set(ctx context.Context, from resourceWorkspaceRoleAssignmentModel) diag.Diagnostics {
	principal, diags := from.Principal.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.Principal = &fabcore.Principal{
		ID:   principal.ID.ValueStringPointer(),
		Type: (*fabcore.PrincipalType)(principal.Type.ValueStringPointer()),
	}
	to.Role = (*fabcore.WorkspaceRole)(from.Role.ValueStringPointer())

	return nil
}

type requestUpdateWorkspaceRoleAssignment struct {
	fabcore.UpdateWorkspaceRoleAssignmentRequest
}

func (to *requestUpdateWorkspaceRoleAssignment) set(from resourceWorkspaceRoleAssignmentModel) {
	to.Role = (*fabcore.WorkspaceRole)(from.Role.ValueStringPointer())
}
