// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
)

type dataSourceFabricItemModel struct {
	fabricItemModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
