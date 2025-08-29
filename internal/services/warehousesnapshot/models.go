// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehousesnapshot

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabwarehousesnapshot "github.com/microsoft/fabric-sdk-go/fabric/warehousesnapshot"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type warehouseSnapshotConfigurationModel struct {
	ParentWarehouseID customtypes.UUID  `tfsdk:"parent_warehouse_id"`
	SnapshotDateTime  timetypes.RFC3339 `tfsdk:"snapshot_date_time"`
}

type warehouseSnapshotPropertiesModel struct {
	ConnectionString  types.String      `tfsdk:"connection_string"`
	ParentWarehouseID customtypes.UUID  `tfsdk:"parent_warehouse_id"`
	SnapshotDateTime  timetypes.RFC3339 `tfsdk:"snapshot_date_time"`
}

func (to *warehouseSnapshotPropertiesModel) set(from fabwarehousesnapshot.Properties) {
	to.ConnectionString = types.StringPointerValue(from.ConnectionString)
	to.ParentWarehouseID = customtypes.NewUUIDPointerValue(from.ParentWarehouseID)
	to.SnapshotDateTime = timetypes.NewRFC3339TimePointerValue(from.SnapshotDateTime)
}
