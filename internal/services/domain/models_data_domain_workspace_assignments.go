// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceDomainWorkspaceAssignmentsModel struct {
	DomainID customtypes.UUID                                   `tfsdk:"domain_id"`
	Values   supertypes.ListNestedObjectValueOf[workspaceModel] `tfsdk:"values"`
	Timeouts timeouts.Value                                     `tfsdk:"timeouts"`
}

func (to *dataSourceDomainWorkspaceAssignmentsModel) setValues(ctx context.Context, from []fabadmin.DomainWorkspace) diag.Diagnostics {
	slice := make([]*workspaceModel, 0, len(from))

	for _, entity := range from {
		var entityModel workspaceModel
		entityModel.set(entity)
		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}
