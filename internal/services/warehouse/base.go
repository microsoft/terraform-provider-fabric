// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehouse

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName           = "Warehouse"
	ItemTFName         = "warehouse"
	ItemsName          = "Warehouses"
	ItemsTFName        = "warehouses"
	ItemType           = fabcore.ItemTypeWarehouse
	ItemDocsSPNSupport = common.DocsSPNNotSupported
	ItemDocsURL        = "https://learn.microsoft.com/fabric/data-warehouse/data-warehousing"
)
