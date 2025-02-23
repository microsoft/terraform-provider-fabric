// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type dataSourceOnPremisesGatewaysModel struct {
	Values   supertypes.ListNestedObjectValueOf[onPremisesGatewayModelBase] `tfsdk:"values"`
	Timeouts timeouts.Value                                                 `tfsdk:"timeouts"`
}

func (m *dataSourceOnPremisesGatewaysModel) setValues(ctx context.Context, from []fabcore.GatewayClassification) diag.Diagnostics {
	var diags diag.Diagnostics
	slice := make([]*onPremisesGatewayModelBase, 0, len(from))

	for _, classification := range from {
		gw, ok := classification.(*fabcore.OnPremisesGateway)
		if !ok {
			continue // skip non-OnPremisesGateway types
		}

		var entityModel onPremisesGatewayModelBase
		if setDiags := entityModel.set(ctx, *gw); setDiags.HasError() {
			diags.Append(setDiags...)
			continue
		}
		slice = append(slice, &entityModel)
	}

	if listDiags := m.Values.Set(ctx, slice); listDiags.HasError() {
		diags.Append(listDiags...)
	}
	return diags
}

type dataSourceVirtualNetworkGatewaysModel struct {
	Values   supertypes.ListNestedObjectValueOf[virtualNetworkGatewayModelBase] `tfsdk:"values"`
	Timeouts timeouts.Value                                                     `tfsdk:"timeouts"`
}

func (m *dataSourceVirtualNetworkGatewaysModel) setValues(ctx context.Context, from []fabcore.GatewayClassification) diag.Diagnostics {
	var diags diag.Diagnostics
	slice := make([]*virtualNetworkGatewayModelBase, 0, len(from))

	for _, entity := range from {
		gw, ok := entity.(*fabcore.VirtualNetworkGateway)

		if !ok {
			continue // skip non-VirtualNetworkGateway types
		}

		var entityModel virtualNetworkGatewayModelBase
		if setDiags := entityModel.set(ctx, *gw); setDiags.HasError() {
			diags.Append(setDiags...)
			continue
		}
		slice = append(slice, &entityModel)
	}

	if listDiags := m.Values.Set(ctx, slice); listDiags.HasError() {
		diags.Append(listDiags...)
	}
	return diags
}

type dataSourceOnPremisesGatewayPersonalsModel struct {
	Values   supertypes.ListNestedObjectValueOf[onPremisesGatewayPersonalModelBase] `tfsdk:"values"`
	Timeouts timeouts.Value                                                         `tfsdk:"timeouts"`
}

func (m *dataSourceOnPremisesGatewayPersonalsModel) setValues(ctx context.Context, from []fabcore.GatewayClassification) diag.Diagnostics {
	var diags diag.Diagnostics
	slice := make([]*onPremisesGatewayPersonalModelBase, 0, len(from))

	for _, classification := range from {
		gw, ok := classification.(*fabcore.OnPremisesGatewayPersonal)
		if !ok {
			continue // skip non-OnPremisesGatewayPersonal types
		}

		var entityModel onPremisesGatewayPersonalModelBase
		if setDiags := entityModel.set(ctx, *gw); setDiags.HasError() {
			diags.Append(setDiags...)
			continue
		}
		slice = append(slice, &entityModel)
	}

	if listDiags := m.Values.Set(ctx, slice); listDiags.HasError() {
		diags.Append(listDiags...)
	}
	return diags
}
