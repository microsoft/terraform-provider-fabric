// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mirroredazuredatabrickscatalog

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

func getDataSourceMirroredAzureDatabricksCatalogPropertiesAttributes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"auto_sync": schema.StringAttribute{
			MarkdownDescription: "Auto sync the catalog.",
			Computed:            true,
		},
		"catalog_name": schema.StringAttribute{
			MarkdownDescription: "Azure databricks catalog name.",
			Computed:            true,
		},
		"databricks_workspace_connection_id": schema.StringAttribute{
			MarkdownDescription: "The Azure databricks workspace connection id.",
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
		},
		"mirror_status": schema.StringAttribute{
			MarkdownDescription: "The MirroredAzureDatabricksCatalog sync status.",
			Computed:            true,
		},
		"mirroring_mode": schema.StringAttribute{
			MarkdownDescription: "Mirroring mode.",
			Computed:            true,
		},
		"onelake_tables_path": schema.StringAttribute{
			MarkdownDescription: "OneLake path to the MirroredAzureDatabricksCatalog tables directory.",
			Computed:            true,
			CustomType:          customtypes.URLType{},
		},
		"sql_endpoint_properties": schema.SingleNestedAttribute{
			MarkdownDescription: "An object containing the properties of the SQL endpoint.",
			Computed:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[sqlEndpointPropertiesModel](ctx),
			Attributes: map[string]schema.Attribute{
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
		"storage_connection_id": schema.StringAttribute{
			MarkdownDescription: "The storage connection id.",
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
		},
		"sync_details": schema.SingleNestedAttribute{
			MarkdownDescription: "The MirroredAzureDatabricksCatalog sync status.",
			Computed:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[syncDetailsModel](ctx),
			Attributes: map[string]schema.Attribute{
				"error_info": schema.SingleNestedAttribute{
					MarkdownDescription: "The error information.",
					Computed:            true,
					CustomType:          supertypes.NewSingleNestedObjectTypeOf[errorInfoModel](ctx),
					Attributes: map[string]schema.Attribute{
						"error_code": schema.StringAttribute{
							MarkdownDescription: "The error code.",
							Computed:            true,
						},
						"error_details": schema.StringAttribute{
							MarkdownDescription: "The error details.",
							Computed:            true,
						},
						"error_message": schema.StringAttribute{
							MarkdownDescription: "The error message.",
							Computed:            true,
						},
					},
				},
				"last_sync_date_time": schema.StringAttribute{
					MarkdownDescription: "The last sync date time in UTC, using the YYYY-MM-DDTHH:mm:ssZ format.",
					Computed:            true,
					CustomType:          timetypes.RFC3339Type{},
				},
				"status": schema.StringAttribute{
					MarkdownDescription: "The sync status.",
					Computed:            true,
				},
			},
		},
	}
}
