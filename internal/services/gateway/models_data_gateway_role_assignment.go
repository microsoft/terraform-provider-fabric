// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceGatewayRoleAssignmentModel struct {
	GatewayID customtypes.UUID `tfsdk:"gateway_id"`
	baseGatewayRoleAssignmentModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

func (to *dataSourceGatewayRoleAssignmentModel) set(ctx context.Context, from fabcore.GatewayRoleAssignment) diag.Diagnostics {
	return to.baseGatewayRoleAssignmentModel.set(ctx, from)
}

type dataSourceGatewayRoleAssignmentsModel struct {
	GatewayID customtypes.UUID                                                   `tfsdk:"gateway_id"`
	Values    supertypes.ListNestedObjectValueOf[baseGatewayRoleAssignmentModel] `tfsdk:"values"`
	Timeouts  timeouts.Value                                                     `tfsdk:"timeouts"`
}

func (to *dataSourceGatewayRoleAssignmentsModel) setValues(ctx context.Context, from []fabcore.GatewayRoleAssignment) diag.Diagnostics {
	slice := make([]*baseGatewayRoleAssignmentModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseGatewayRoleAssignmentModel

		if diags := entityModel.set(ctx, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}
