// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package variablelibrary

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func getDataSourceVariableLibraryPropertiesAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"active_value_set_name": schema.StringAttribute{
			MarkdownDescription: "The VariableLibrary current active value set.",
			Computed:            true,
		},
	}
}
