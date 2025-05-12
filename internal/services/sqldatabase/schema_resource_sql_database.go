// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sqldatabase

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func getResourceSQLDatabasePropertiesAttributes() map[string]schema.Attribute {
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
	}
}
