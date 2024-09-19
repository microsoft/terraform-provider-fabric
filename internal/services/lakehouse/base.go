// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName                     = "Lakehouse"
	ItemTFName                   = "lakehouse"
	ItemsName                    = "Lakehouses"
	ItemsTFName                  = "lakehouses"
	ItemType                     = fabcore.ItemTypeLakehouse
	ItemDocsSPNSupport           = common.DocsSPNSupported
	ItemDocsURL                  = "https://learn.microsoft.com/training/modules/get-started-lakehouses"
	LakehouseTableName           = "Lakehouse Table"
	LakehouseTableTFName         = "lakehouse_table"
	LakehouseTablesName          = "Lakehouse Tables"
	LakehouseTablesTFName        = "lakehouse_tables"
	LakehouseTableDocsSPNSupport = common.DocsSPNSupported
	LakehouseTableDocsURL        = "https://learn.microsoft.com/fabric/data-engineering/lakehouse-and-delta-tables"
)
