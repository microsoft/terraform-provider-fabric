// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package snowflakedatabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

func getResourceSnowflakeDatabasePropertiesAttributes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"connection_id": schema.StringAttribute{
			MarkdownDescription: "The connection ID for the Snowflake Database.",
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
		},
		"default_schema": schema.StringAttribute{
			MarkdownDescription: "The default schema name for the Snowflake Database.",
			Computed:            true,
		},
		"onelake_tables_path": schema.StringAttribute{
			MarkdownDescription: "OneLake path to the Snowflake Database tables directory.",
			Computed:            true,
		},
		"snowflake_account_url": schema.StringAttribute{
			MarkdownDescription: "The Snowflake account URL.",
			Computed:            true,
		},
		"snowflake_database_name": schema.StringAttribute{
			MarkdownDescription: "The Snowflake database name.",
			Computed:            true,
		},
		"snowflake_volume_path": schema.StringAttribute{
			MarkdownDescription: "OneLake path to the Snowflake Database files directory.",
			Computed:            true,
		},
		"sql_endpoint_properties": schema.SingleNestedAttribute{
			MarkdownDescription: "An object containing the properties of the SQL endpoint.",
			Computed:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[snowflakeDatabaseSQLEndpointPropertiesModel](ctx),
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
	}
}

func getResourceSnowflakeDatabaseConfigurationAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"connection_id": schema.StringAttribute{
			MarkdownDescription: "The connection ID for the Snowflake Database.",
			Required:            true,
			CustomType:          customtypes.UUIDType{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"snowflake_database_name": schema.StringAttribute{
			MarkdownDescription: "The Snowflake database name.",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
	}
}
