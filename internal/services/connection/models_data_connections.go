// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type dataSourceConnectionsModel struct {
	Values   supertypes.SetNestedObjectValueOf[baseConnectionModel[dsConnectionDetailsModel, dsCredentialDetailsModel]] `tfsdk:"values"`
	Timeouts timeouts.Value                                                                                             `tfsdk:"timeouts"`
}

func (to *dataSourceConnectionsModel) setValues(ctx context.Context, from []fabcore.Connection) diag.Diagnostics {
	slice := make([]*baseConnectionModel[dsConnectionDetailsModel, dsCredentialDetailsModel], 0, len(from))

	for _, entity := range from {
		var entityModel baseConnectionModel[dsConnectionDetailsModel, dsCredentialDetailsModel]

		if diags := entityModel.set(ctx, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}
