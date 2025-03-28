// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehouse

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const FabricItemType = fabcore.ItemTypeWarehouse

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Warehouse",
	Type:           "warehouse",
	Names:          "Warehouses",
	Types:          "warehouses",
	DocsURL:        "https://learn.microsoft.com/fabric/data-warehouse/data-warehousing",
	IsPreview:      true,
	IsSPNSupported: true,
}
