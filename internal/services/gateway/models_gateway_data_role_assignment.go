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

		if diags := entityModel.set(entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

type gatewayRoleAssignmentModel struct {
	ID          customtypes.UUID `tfsdk:"id"`
	Role        types.String     `tfsdk:"role"`
	DisplayName types.String     `tfsdk:"display_name"`
	Type        types.String     `tfsdk:"type"`
}

func (to *gatewayRoleAssignmentModel) set(from fabcore.GatewayRoleAssignment) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Role = types.StringPointerValue((*string)(from.Role))

	return nil
}
