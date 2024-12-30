// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type lakehouseConfigurationModel struct {
	EnableSchemas types.Bool `tfsdk:"enable_schemas"`
}
