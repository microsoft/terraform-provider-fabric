// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

func getDataSourceLakehousePropertiesAttributes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"onelake_files_path": schema.StringAttribute{
			MarkdownDescription: "OneLake path to the Lakehouse files directory",
			Computed:            true,
		},
		"onelake_tables_path": schema.StringAttribute{
			MarkdownDescription: "OneLake path to the Lakehouse tables directory.",
			Computed:            true,
		},
		"sql_endpoint_properties": schema.SingleNestedAttribute{
			MarkdownDescription: "An object containing the properties of the SQL endpoint.",
			Computed:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[lakehouseSQLEndpointPropertiesModel](ctx),
			Attributes: map[string]schema.Attribute{
				"provisioning_status": schema.StringAttribute{
					MarkdownDescription: "The SQL endpoint provisioning status.",
					Computed:            true,
				},
				"connection_string": schema.StringAttribute{
					MarkdownDescription: "SQL endpoint connection string.",
					Computed:            true,
				},
				"id": schema.StringAttribute{
					MarkdownDescription: "SQL endpoint ID.",
					Computed:            true,
					CustomType:          customtypes.UUIDType{},
				},
			},
		},
		"default_schema": schema.StringAttribute{
			MarkdownDescription: "Default schema of the Lakehouse. This property is returned only for schema enabled Lakehouse.",
			Computed:            true,
		},
	}
}
