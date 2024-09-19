// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

type dataSourceWorkspacesModel struct {
	Values   supertypes.ListNestedObjectValueOf[baseWorkspaceModel] `tfsdk:"values"`
	Timeouts timeouts.Value                                         `tfsdk:"timeouts"`
}

func (to *dataSourceWorkspacesModel) setValues(ctx context.Context, from []fabcore.Workspace) diag.Diagnostics {
	slice := make([]*baseWorkspaceModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseWorkspaceModel
		entityModel.set(entity)
		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}
