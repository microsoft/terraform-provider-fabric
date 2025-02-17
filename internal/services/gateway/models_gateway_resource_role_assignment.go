// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type resourceGatewayRoleAssignmentModel struct {
	ID            customtypes.UUID `tfsdk:"id"`
	PrincipalID   customtypes.UUID `tfsdk:"principal_id"`
	PrincipalType types.String     `tfsdk:"principal_type"`
	Role          types.String     `tfsdk:"role"`
	GatewayID     customtypes.UUID `tfsdk:"gateway_id"`
	Timeouts      timeouts.Value   `tfsdk:"timeouts"`
}

func (to *resourceGatewayRoleAssignmentModel) set(from fabcore.GatewayRoleAssignment) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.PrincipalID = customtypes.NewUUIDPointerValue(from.Principal.ID)
	to.PrincipalType = types.StringPointerValue((*string)(from.Principal.Type))
	to.Role = types.StringPointerValue((*string)(from.Role))
}

type requestCreateGatewayRoleAssignment struct {
	fabcore.AddGatewayRoleAssignmentRequest
}

func (to *requestCreateGatewayRoleAssignment) set(from resourceGatewayRoleAssignmentModel) {
	to.Principal = &fabcore.Principal{
		ID:   from.PrincipalID.ValueStringPointer(),
		Type: (*fabcore.PrincipalType)(from.PrincipalType.ValueStringPointer()),
	}
	to.Role = (*fabcore.GatewayRole)(from.Role.ValueStringPointer())
}

type requestUpdateGatewayRoleAssignment struct {
	fabcore.UpdateGatewayRoleAssignmentRequest
}

func (to *requestUpdateGatewayRoleAssignment) set(from resourceGatewayRoleAssignmentModel) {
	to.Role = (*fabcore.GatewayRole)(from.Role.ValueStringPointer())
}
