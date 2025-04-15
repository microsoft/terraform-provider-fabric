// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehouse

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	fabwarehouse "github.com/microsoft/fabric-sdk-go/fabric/warehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getDataSourceWarehousePropertiesAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"collation_type": schema.StringAttribute{
			MarkdownDescription: "The collation type of the warehouse. Possible values:" + utils.ConvertStringSlicesToString(fabwarehouse.PossibleCollationTypeValues(), true, true) + ".",
			Computed:            true,
		},
		"connection_string": schema.StringAttribute{
			MarkdownDescription: "The SQL connection string connected to the workspace containing this warehouse.",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "The date and time the warehouse was created.",
			Computed:            true,
			CustomType:          timetypes.RFC3339Type{},
		},
		"last_updated_time": schema.StringAttribute{
			MarkdownDescription: "The date and time the warehouse was last updated.",
			Computed:            true,
			CustomType:          timetypes.RFC3339Type{},
		},
	}
}
