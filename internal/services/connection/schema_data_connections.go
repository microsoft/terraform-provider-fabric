// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

func DataSourceConnectionsSchema(ctx context.Context, typeInfo tftypeinfo.TFTypeInfo) schema.Schema {
	tflog.Info(ctx, "Building schema for connections data source")
	
	// First get the data source schema for a single connection
	singleSchema := DataSourceConnectionSchema(ctx)

	return schema.Schema{
		MarkdownDescription: "Use this data source to retrieve a list of all available Microsoft Fabric Connections.",
		Attributes: map[string]schema.Attribute{
			"values": schema.SetNestedAttribute{
				MarkdownDescription: "The set of " + typeInfo.Names + ".",
				Computed:            true,
				CustomType:          supertypes.NewSetNestedObjectTypeOf[baseConnectionModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: singleSchema.Attributes,
				},
			},
		},
	}
}