// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

func getResourceLakehousePropertiesAttributes(ctx context.Context) map[string]schema.Attribute {
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

func getResourceLakehouseConfigurationAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"enable_schemas": schema.BoolAttribute{
			MarkdownDescription: "Schema enabled Lakehouse.",
			Required:            true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
		},
	}
}
