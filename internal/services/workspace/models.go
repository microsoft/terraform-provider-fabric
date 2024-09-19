// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"
	"fmt"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type baseWorkspaceInfoModel struct {
	baseWorkspaceModel
	CapacityAssignmentProgress types.String                                                 `tfsdk:"capacity_assignment_progress"`
	Identity                   supertypes.SingleNestedObjectValueOf[workspaceIdentityModel] `tfsdk:"identity"`
}

func (to *baseWorkspaceInfoModel) set(ctx context.Context, from fabcore.WorkspaceInfo) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
	to.Type = types.StringPointerValue((*string)(from.Type))
	to.CapacityID = customtypes.NewUUIDPointerValue(from.CapacityID)
	to.CapacityAssignmentProgress = types.StringPointerValue((*string)(from.CapacityAssignmentProgress))

	workspaceIdentity := supertypes.NewSingleNestedObjectValueOfNull[workspaceIdentityModel](ctx)

	if from.WorkspaceIdentity != nil {
		workspaceIdentityModel := &workspaceIdentityModel{}
		workspaceIdentityModel.set(from.WorkspaceIdentity)
		workspaceIdentity.Set(ctx, workspaceIdentityModel)
	}

	to.Identity = workspaceIdentity
}

type baseWorkspaceModel struct {
	ID          customtypes.UUID `tfsdk:"id"`
	DisplayName types.String     `tfsdk:"display_name"`
	Description types.String     `tfsdk:"description"`
	Type        types.String     `tfsdk:"type"`
	CapacityID  customtypes.UUID `tfsdk:"capacity_id"`
}

func (to *baseWorkspaceModel) set(from fabcore.Workspace) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Description = types.StringPointerValue(from.Description)
	to.Type = types.StringPointerValue((*string)(from.Type))
	to.CapacityID = customtypes.NewUUIDPointerValue(from.CapacityID)
}

type workspaceIdentityModel struct {
	Type               types.String     `tfsdk:"type"`
	ApplicationID      customtypes.UUID `tfsdk:"application_id"`
	ServicePrincipalID customtypes.UUID `tfsdk:"service_principal_id"`
}

func (to *workspaceIdentityModel) set(from *fabcore.WorkspaceIdentity) {
	to.Type = types.StringValue(workspaceIdentityTypes[0])
	to.ApplicationID = customtypes.NewUUIDPointerValue(from.ApplicationID)
	to.ServicePrincipalID = customtypes.NewUUIDPointerValue(from.ServicePrincipalID)
}

func checkWorkspaceType(entity fabcore.WorkspaceInfo) diag.Diagnostics {
	var diags diag.Diagnostics

	switch *entity.Type {
	case fabcore.WorkspaceTypePersonal:
		diags.AddError(
			common.ErrorWorkspaceNotSupportedHeader,
			fmt.Sprintf(common.ErrorWorkspaceNotSupportedDetails, string(fabcore.WorkspaceTypePersonal)),
		)

		return diags
	case fabcore.WorkspaceTypeAdminWorkspace:
		diags.AddError(
			common.ErrorWorkspaceNotSupportedHeader,
			fmt.Sprintf(common.ErrorWorkspaceNotSupportedDetails, string(fabcore.WorkspaceTypeAdminWorkspace)),
		)

		return diags
	default:
		return nil
	}
}
