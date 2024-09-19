// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package capacity

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
)

type dataSourceCapacityModel struct {
	baseCapacityModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
