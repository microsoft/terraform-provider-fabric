// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehouse

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func getDataSourceWarehousePropertiesAttributes() map[string]schema.Attribute {
	result := map[string]schema.Attribute{
		"connection_string": schema.StringAttribute{
			MarkdownDescription: "The SQL connection string connected to the workspace containing this warehouse.",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "The date and time the warehouse was created.",
			Computed:            true,
			CustomType:          timetypes.RFC3339Type{},
		},
		"last_updated_time": schema.StringAttribute{
			MarkdownDescription: "The date and time the warehouse was last updated.",
			Computed:            true,
			CustomType:          timetypes.RFC3339Type{},
		},
	}

	return result
}
