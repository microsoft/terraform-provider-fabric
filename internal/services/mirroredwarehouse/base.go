// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package mirroredwarehouse

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const FabricItemType = fabcore.ItemTypeMirroredWarehouse

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Mirrored Warehouse",
	Type:           "mirrored_warehouse",
	Names:          "Mirrored Warehouses",
	Types:          "mirrored_warehouses",
	DocsURL:        "https://learn.microsoft.com/fabric/database/mirrored-database/overview",
	IsPreview:      true,
	IsSPNSupported: false,
}
