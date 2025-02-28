// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sqldatabase

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName           = "SQL Database"
	ItemTFName         = "sql_database"
	ItemsName          = "SQL Databases"
	ItemsTFName        = "sql_databases"
	ItemType           = fabcore.ItemTypeSQLDatabase
	ItemDocsSPNSupport = common.DocsSPNSupported
	ItemDocsURL        = "https://learn.microsoft.com/fabric/database/sql/overview"
	ItemPreview        = true
)
