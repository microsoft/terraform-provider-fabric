// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehouse

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
)

type dataSourceWarehouseModel struct {
	baseWarehouseModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
