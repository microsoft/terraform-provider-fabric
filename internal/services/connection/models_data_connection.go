// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
)

type dataSourceConnectionModel struct {
	baseConnectionModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
