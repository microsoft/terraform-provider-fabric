// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type resourceGatewayRoleAssignmentModel struct {
	baseGatewayRoleAssignmentModel
	GatewayID customtypes.UUID `tfsdk:"gateway_id"`
	Timeouts  timeouts.Value   `tfsdk:"timeouts"`
}

func (to *resourceGatewayRoleAssignmentModel) set(ctx context.Context, from fabcore.GatewayRoleAssignment) diag.Diagnostics {
	return to.baseGatewayRoleAssignmentModel.set(ctx, from)
}

type requestCreateGatewayRoleAssignment struct {
	fabcore.AddGatewayRoleAssignmentRequest
}

func (to *requestCreateGatewayRoleAssignment) set(ctx context.Context, from resourceGatewayRoleAssignmentModel) diag.Diagnostics {
	principal, diags := from.Principal.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.Principal = &fabcore.Principal{ID: principal.ID.ValueStringPointer(), Type: (*fabcore.PrincipalType)(principal.Type.ValueStringPointer())}
	to.Role = (*fabcore.GatewayRole)(from.Role.ValueStringPointer())

	return nil
}

type requestUpdateGatewayRoleAssignment struct {
	fabcore.UpdateGatewayRoleAssignmentRequest
}

func (to *requestUpdateGatewayRoleAssignment) set(from resourceGatewayRoleAssignmentModel) {
	to.Role = (*fabcore.GatewayRole)(from.Role.ValueStringPointer())
}
