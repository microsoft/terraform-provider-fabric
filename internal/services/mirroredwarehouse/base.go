// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mirroredwarehouse

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName           = "Mirrored Warehouse"
	ItemTFName         = "mirrored_warehouse"
	ItemsName          = "Mirrored Warehouses"
	ItemsTFName        = "mirrored_warehouses"
	ItemType           = fabcore.ItemTypeMirroredWarehouse
	ItemDocsSPNSupport = common.DocsSPNNotSupported
	ItemDocsURL        = "https://learn.microsoft.com/fabric/database/mirrored-database/overview"
)
