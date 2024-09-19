// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
)

type dataSourceLakehouseModel struct {
	baseLakehouseModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
