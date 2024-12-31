// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
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
