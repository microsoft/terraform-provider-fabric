// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
)

type dataSourceWorkspaceModel struct {
	baseWorkspaceInfoModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
