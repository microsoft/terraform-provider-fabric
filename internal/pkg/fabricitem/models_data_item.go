// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
)

type dataSourceFabricItemModel struct {
	baseFabricItemModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type DataSourceFabricItemDefinitionPartModel struct {
	Content types.String `tfsdk:"content"`
}

func (to *DataSourceFabricItemDefinitionPartModel) Set(from string) diag.Diagnostics {
	content := from

	if diags := transforms.PayloadToGzip(&content); diags.HasError() {
		return diags
	}

	to.Content = types.StringPointerValue(&content)

	return nil
}
