// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type dataSourceDomainsModel struct {
	Values   supertypes.ListNestedObjectValueOf[baseDomainModel] `tfsdk:"values"`
	Timeouts timeouts.Value                                      `tfsdk:"timeouts"`
}

func (to *dataSourceDomainsModel) setValues(ctx context.Context, from []fabadmin.Domain) diag.Diagnostics {
	slice := make([]*baseDomainModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseDomainModel
		entityModel.set(entity)
		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}
