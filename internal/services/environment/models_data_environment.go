// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package environment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
)

type dataSourceEnvironmentModel struct {
	baseEnvironmentPropertiesModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
