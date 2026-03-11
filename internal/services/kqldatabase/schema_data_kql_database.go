// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package kqldatabase

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	fabkqldatabase "github.com/microsoft/fabric-sdk-go/fabric/kqldatabase"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getDataSourceKQLDatabasePropertiesAttributes() map[string]schema.Attribute {
	result := map[string]schema.Attribute{
		"database_type": schema.StringAttribute{
			MarkdownDescription: "The type of the database. Possible values:" + utils.ConvertStringSlicesToString(fabkqldatabase.PossibleKqlDatabaseTypeValues(), true, true) + ".",
			Computed:            true,
		},
		"eventhouse_id": schema.StringAttribute{
			MarkdownDescription: "Parent Eventhouse ID.",
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
		},
		"ingestion_service_uri": schema.StringAttribute{
			MarkdownDescription: "Ingestion service URI.",
			Computed:            true,
			CustomType:          customtypes.URLType{},
		},
		"query_service_uri": schema.StringAttribute{
			MarkdownDescription: "Query service URI.",
			Computed:            true,
			CustomType:          customtypes.URLType{},
		},
	}

	return result
}
