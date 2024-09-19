// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type resourceWorkspaceRoleAssignmentModel struct {
	ID            customtypes.UUID `tfsdk:"id"`
	PrincipalID   customtypes.UUID `tfsdk:"principal_id"`
	PrincipalType types.String     `tfsdk:"principal_type"`
	Role          types.String     `tfsdk:"role"`
	WorkspaceID   customtypes.UUID `tfsdk:"workspace_id"`
	Timeouts      timeouts.Value   `tfsdk:"timeouts"`
}

func (to *resourceWorkspaceRoleAssignmentModel) set(from fabcore.WorkspaceRoleAssignment) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.PrincipalID = customtypes.NewUUIDPointerValue(from.Principal.ID)
	to.PrincipalType = types.StringPointerValue((*string)(from.Principal.Type))
	to.Role = types.StringPointerValue((*string)(from.Role))
}

type requestCreateWorkspaceRoleAssignment struct {
	fabcore.AddWorkspaceRoleAssignmentRequest
}

func (to *requestCreateWorkspaceRoleAssignment) set(from resourceWorkspaceRoleAssignmentModel) {
	to.Principal = &fabcore.Principal{ID: from.PrincipalID.ValueStringPointer(), Type: (*fabcore.PrincipalType)(from.PrincipalType.ValueStringPointer())}
	to.Role = (*fabcore.WorkspaceRole)(from.Role.ValueStringPointer())
}

type requestUpdateWorkspaceRoleAssignment struct {
	fabcore.UpdateWorkspaceRoleAssignmentRequest
}

func (to *requestUpdateWorkspaceRoleAssignment) set(from resourceWorkspaceRoleAssignmentModel) {
	to.Role = (*fabcore.WorkspaceRole)(from.Role.ValueStringPointer())
}
