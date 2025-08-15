// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehousesnapshot

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

func getResourceWarehouseSnapshotPropertiesAttributes() map[string]schema.Attribute {
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

func getResourceWarehouseSnapshotConfigurationAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"parent_warehouse_id": schema.StringAttribute{
			MarkdownDescription: "The parent Warehouse ID.",
			CustomType:          customtypes.UUIDType{},
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"snapshot_date_time": schema.StringAttribute{
			MarkdownDescription: "The date and time used for the Warehouse snapshot, if not provided the current date and time will be taken. If given it should be in UTC, using the YYYY-MM-DDTHH:mm:ssZ format.",
			CustomType:          timetypes.RFC3339Type{},
			Optional:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}
}
