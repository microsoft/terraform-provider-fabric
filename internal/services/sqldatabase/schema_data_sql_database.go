// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sqldatabase

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func getDataSourceSQLDatabasePropertiesAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"connection_string": schema.StringAttribute{
			MarkdownDescription: "The connection string of the database.",
			Computed:            true,
		},
		"database_name": schema.StringAttribute{
			MarkdownDescription: "The database name.",
			Computed:            true,
		},
		"server_fqdn": schema.StringAttribute{
			MarkdownDescription: "The server fully qualified domain name (FQDN).",
			Computed:            true,
		},
		"backup_retention_days": schema.Int32Attribute{
			MarkdownDescription: "The backup retention period in days.",
			Computed:            true,
		},
		"collation": schema.StringAttribute{
			MarkdownDescription: "The collation of the SQL database.",
			Computed:            true,
		},
		"earliest_restore_point": schema.StringAttribute{
			MarkdownDescription: "The earliest restore point of the database in UTC, using the YYYY-MM-DDTHH:mm:ssZ format.",
			Computed:            true,
			CustomType:          timetypes.RFC3339Type{},
		},
		"latest_restore_point": schema.StringAttribute{
			MarkdownDescription: "The latest restore point of the database in UTC, using the YYYY-MM-DDTHH:mm:ssZ format.",
			Computed:            true,
			CustomType:          timetypes.RFC3339Type{},
		},
	}
}
