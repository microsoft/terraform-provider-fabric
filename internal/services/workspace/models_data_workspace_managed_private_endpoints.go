// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceWorkspaceManagedPrivateEndpointsModel struct {
	WorkspaceID customtypes.UUID                                                             `tfsdk:"workspace_id"`
	Values      supertypes.ListNestedObjectValueOf[baseWorkspaceManagedPrivateEndpointModel] `tfsdk:"values"`
	Timeouts    timeouts.Value                                                               `tfsdk:"timeouts"`
}

func (to *dataSourceWorkspaceManagedPrivateEndpointsModel) setValues(ctx context.Context, from []fabcore.ManagedPrivateEndpoint) diag.Diagnostics {
	slice := make([]*baseWorkspaceManagedPrivateEndpointModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseWorkspaceManagedPrivateEndpointModel
		if diags := entityModel.set(ctx, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}
