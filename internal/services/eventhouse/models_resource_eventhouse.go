// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventhouse

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type eventhouseConfigurationModel struct {
	MinimumConsumptionUnits types.Float64 `tfsdk:"minimum_consumption_units"`
}
