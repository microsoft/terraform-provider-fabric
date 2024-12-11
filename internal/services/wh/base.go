// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package wh

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName           = "WH"
	ItemTFName         = "wh"
	ItemsName          = "WHs"
	ItemsTFName        = "whss"
	ItemType           = fabcore.ItemTypeWarehouse
	ItemDocsSPNSupport = common.DocsSPNNotSupported
	ItemDocsURL        = "https://learn.microsoft.com/fabric/data-warehouse/data-warehousing"
)
