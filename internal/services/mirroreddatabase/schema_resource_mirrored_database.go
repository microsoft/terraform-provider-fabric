// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mirroreddatabase

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func getResourceMirroredDatabasePropertiesAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"default_schema": schema.StringAttribute{
			MarkdownDescription: "Default schema of the mirrored database, returned only for mirrored databases that enable default schema in definition.",
			Computed:            true,
		},
		"onelake_tables_path": schema.StringAttribute{
			MarkdownDescription: "OneLake path to the mirrored database tables directory.",
			Computed:            true,
		},
		"sql_endpoint_properties": schema.SingleNestedAttribute{
			MarkdownDescription: "An object containing the properties of the SQL endpoint.",
			Computed:            true,
			Attributes: map[string]schema.Attribute{
				"provisioning_status": schema.StringAttribute{
					MarkdownDescription: "The SQL endpoint provisioning status.",
					Computed:            true,
				},
				"connection_string": schema.StringAttribute{
					MarkdownDescription: "The SQL endpoint connection string.",
					Computed:            true,
				},
				"id": schema.StringAttribute{
					MarkdownDescription: "The SQL endpoint ID.",
					Computed:            true,
				},
			},
		},
	}
}
