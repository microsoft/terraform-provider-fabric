// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventhouse

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	fabeventhouse "github.com/microsoft/fabric-sdk-go/fabric/eventhouse"
)

type eventhouseConfigurationModel struct {
	MinimumConsumptionUnits types.Float64 `tfsdk:"minimum_consumption_units"`
}

func (to *eventhouseConfigurationModel) set(ctx context.Context, from *fabeventhouse.CreationPayload) {
	to.MinimumConsumptionUnits = types.Float64PointerValue(from.MinimumConsumptionUnits)
}
