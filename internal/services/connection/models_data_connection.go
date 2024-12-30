// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
)

type dataSourceConnectionModel struct {
	baseDataSourceConnectionModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
