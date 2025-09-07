// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connectionra

import (
	"context"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	//revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseConnectionRoleAssignmentModel struct {
	ID        customtypes.UUID                                     `tfsdk:"id"`
	Role      types.String                                         `tfsdk:"role"`
	Principal supertypes.SingleNestedObjectValueOf[principalModel] `tfsdk:"principal"`
}

func (to *baseConnectionRoleAssignmentModel) set(ctx context.Context, from fabcore.ConnectionRoleAssignment) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Role = types.StringPointerValue((*string)(from.Role))

	if from.Principal != nil {
		principalModel := &principalModel{}

		principalModel.set(*from.Principal)

		if diags := to.Principal.Set(ctx, principalModel); diags.HasError() {
			return diags
		}
	}

	return nil
}

/*
DATA-SOURCE
*/

type dataSourceConnectionRoleAssignmentModel struct {
	baseConnectionRoleAssignmentModel

	ConnectionID               customtypes.UUID `tfsdk:"connection_id"`
	ConnectionRoleAssignmentID customtypes.UUID `tfsdk:"connection_role_assignment_id"`

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

func (to *dataSourceConnectionRoleAssignmentModel) set(ctx context.Context, connectionID, connectionRoleAssignmentID string, from fabcore.ConnectionRoleAssignment) diag.Diagnostics {
	to.ConnectionID = customtypes.NewUUIDPointerValue(&connectionID)
	to.ConnectionRoleAssignmentID = customtypes.NewUUIDPointerValue(&connectionRoleAssignmentID)
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Role = types.StringPointerValue((*string)(from.Role))

	if from.Principal != nil {
		principalModel := &principalModel{}

		principalModel.set(*from.Principal)

		if diags := to.Principal.Set(ctx, principalModel); diags.HasError() {
			return diags
		}
	}

	return nil
}

/*
DATA-SOURCE (list)
*/

type dataSourceConnectionRoleAssignmentsModel struct {
	ConnectionID customtypes.UUID                                                     `tfsdk:"connection_id"`
	Values       supertypes.SetNestedObjectValueOf[baseConnectionRoleAssignmentModel] `tfsdk:"values"`
	Timeouts     timeoutsD.Value                                                      `tfsdk:"timeouts"`
}

func (to *dataSourceConnectionRoleAssignmentsModel) setValues(ctx context.Context, connectionID string, from []fabcore.ConnectionRoleAssignment) diag.Diagnostics {
	to.ConnectionID = customtypes.NewUUIDValue(connectionID)
	slice := make([]*baseConnectionRoleAssignmentModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseConnectionRoleAssignmentModel

		if diags := entityModel.set(ctx, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

/*
RESOURCE
*/

type resourceConnectionRoleAssignmentModel struct {
	baseConnectionRoleAssignmentModel

	ConnectionID customtypes.UUID `tfsdk:"connection_id"`
	Timeouts     timeoutsR.Value  `tfsdk:"timeouts"`
}

type requestCreateConnectionRoleAssignment struct {
	fabcore.AddConnectionRoleAssignmentRequest
}

func (to *requestCreateConnectionRoleAssignment) set(ctx context.Context, from resourceConnectionRoleAssignmentModel) diag.Diagnostics {
	principal, diags := from.Principal.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.Principal = &fabcore.Principal{
		ID:   principal.ID.ValueStringPointer(),
		Type: (*fabcore.PrincipalType)(principal.Type.ValueStringPointer()),
	}
	to.Role = (*fabcore.ConnectionRole)(from.Role.ValueStringPointer())

	return nil
}

type requestUpdateConnectionRoleAssignment struct {
	fabcore.UpdateConnectionRoleAssignmentRequest
}

func (to *requestUpdateConnectionRoleAssignment) set(from resourceConnectionRoleAssignmentModel) {
	to.Role = (*fabcore.ConnectionRole)(from.Role.ValueStringPointer())
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
