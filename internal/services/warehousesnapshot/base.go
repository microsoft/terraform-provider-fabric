// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehousesnapshot

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const FabricItemType = fabcore.ItemTypeWarehouseSnapshot

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Warehouse Snapshot",
	Type:           "warehouse_snapshot",
	Names:          "Warehouse Snapshots",
	Types:          "warehouse_snapshots",
	DocsURL:        "https://learn.microsoft.com/fabric/data-warehouse/warehouse-snapshot",
	IsPreview:      true,
	IsSPNSupported: true,
}
