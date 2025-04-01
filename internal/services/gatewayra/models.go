// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gatewayra

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseGatewayRoleAssignmentModel struct {
	ID        customtypes.UUID                                     `tfsdk:"id"`
	GatewayID customtypes.UUID                                     `tfsdk:"gateway_id"`
	Role      types.String                                         `tfsdk:"role"`
	Principal supertypes.SingleNestedObjectValueOf[principalModel] `tfsdk:"principal"`
}

func (to *baseGatewayRoleAssignmentModel) set(ctx context.Context, gatewayID string, from fabcore.GatewayRoleAssignment) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.GatewayID = customtypes.NewUUIDValue(gatewayID)
	to.Role = types.StringPointerValue((*string)(from.Role))

	principal := supertypes.NewSingleNestedObjectValueOfNull[principalModel](ctx)

	if from.Principal != nil {
		principalModel := &principalModel{}

		principalModel.set(*from.Principal)

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

type dataSourceGatewayRoleAssignmentModel struct {
	baseGatewayRoleAssignmentModel
	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceGatewayRoleAssignmentsModel struct {
	GatewayID customtypes.UUID                                                  `tfsdk:"gateway_id"`
	Values    supertypes.SetNestedObjectValueOf[baseGatewayRoleAssignmentModel] `tfsdk:"values"`
	Timeouts  timeoutsD.Value                                                   `tfsdk:"timeouts"`
}

func (to *dataSourceGatewayRoleAssignmentsModel) setValues(ctx context.Context, gatewayID string, from []fabcore.GatewayRoleAssignment) diag.Diagnostics {
	slice := make([]*baseGatewayRoleAssignmentModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseGatewayRoleAssignmentModel

		if diags := entityModel.set(ctx, gatewayID, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

/*
RESOURCE
*/

type resourceGatewayRoleAssignmentModel struct {
	baseGatewayRoleAssignmentModel
	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateGatewayRoleAssignment struct {
	fabcore.AddGatewayRoleAssignmentRequest
}

func (to *requestCreateGatewayRoleAssignment) set(ctx context.Context, from resourceGatewayRoleAssignmentModel) diag.Diagnostics {
	principal, diags := from.Principal.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.Principal = &fabcore.Principal{
		ID:   principal.ID.ValueStringPointer(),
		Type: (*fabcore.PrincipalType)(principal.Type.ValueStringPointer()),
	}
	to.Role = (*fabcore.GatewayRole)(from.Role.ValueStringPointer())

	return nil
}

type requestUpdateGatewayRoleAssignment struct {
	fabcore.UpdateGatewayRoleAssignmentRequest
}

func (to *requestUpdateGatewayRoleAssignment) set(from resourceGatewayRoleAssignmentModel) {
	to.Role = (*fabcore.GatewayRole)(from.Role.ValueStringPointer())
}

/*
HELPER MODELS
*/

type principalModel struct {
	ID   customtypes.UUID `tfsdk:"id"`
	Type types.String     `tfsdk:"type"`
}

func (to *principalModel) set(from fabcore.Principal) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Type = types.StringPointerValue((*string)(from.Type))
}
