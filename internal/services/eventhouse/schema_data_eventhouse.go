// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventhouse

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func getDataSourceEventhousePropertiesAttributes(ctx context.Context) map[string]schema.Attribute {
	result := map[string]schema.Attribute{
		"ingestion_service_uri": schema.StringAttribute{
			MarkdownDescription: "Ingestion service URI.",
			Computed:            true,
		},
		"query_service_uri": schema.StringAttribute{
			MarkdownDescription: "Query service URI.",
			Computed:            true,
		},
		"database_ids": schema.ListAttribute{
			MarkdownDescription: "List of all KQL Database children IDs.",
			Computed:            true,
			CustomType:          supertypes.NewListTypeOf[string](ctx),
		},
	}

	return result
}
