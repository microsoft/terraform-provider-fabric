// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package warehousesnapshot

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

func getDataSourceWarehouseSnapshotPropertiesAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"connection_string": schema.StringAttribute{
			MarkdownDescription: "The SQL connection string connected to the workspace containing this warehouse.",
			Computed:            true,
		},
		"parent_warehouse_id": schema.StringAttribute{
			MarkdownDescription: "The parent Warehouse ID.",
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
		},
		"snapshot_date_time": schema.StringAttribute{
			MarkdownDescription: "The current warehouse snapshot date and time in UTC, using the YYYY-MM-DDTHH:mm:ssZ format.",
			Computed:            true,
			CustomType:          timetypes.RFC3339Type{},
		},
	}
}
