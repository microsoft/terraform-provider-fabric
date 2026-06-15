// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package mirroredcatalog

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

func getDataSourceMirroredCatalogPropertiesAttributes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"connection_id": schema.StringAttribute{
			MarkdownDescription: "The connection ID used for the mirrored catalog.",
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
		},
		"onelake_tables_path": schema.StringAttribute{
			MarkdownDescription: "OneLake path to the mirrored catalog tables directory.",
			Computed:            true,
		},
		"scope": schema.ListAttribute{
			MarkdownDescription: "The namespace hierarchy path that scopes the mirroring.",
			Computed:            true,
			CustomType:          supertypes.NewListTypeOf[string](ctx),
			ElementType:         types.StringType,
		},
		"source_type": schema.StringAttribute{
			MarkdownDescription: "The source type for the underlying connection (e.g. `DremioIcebergCatalog`).",
			Computed:            true,
		},
		"sql_endpoint_properties": schema.SingleNestedAttribute{
			MarkdownDescription: "An object containing the properties of the SQL endpoint.",
			Computed:            true,
			CustomType:          supertypes.NewSingleNestedObjectTypeOf[mirroredCatalogSQLEndpointPropertiesModel](ctx),
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
					CustomType:          customtypes.UUIDType{},
				},
			},
		},
	}
}
