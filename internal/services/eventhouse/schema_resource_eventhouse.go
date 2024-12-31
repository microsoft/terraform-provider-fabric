// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventhouse

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func getResourceEventhousePropertiesAttributes(ctx context.Context) map[string]schema.Attribute {
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
