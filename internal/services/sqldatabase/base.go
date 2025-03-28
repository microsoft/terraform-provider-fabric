// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sqldatabase

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType = fabcore.ItemTypeSQLDatabase
	ItemDocsURL    = "https://learn.microsoft.com/fabric/database/sql/overview"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "SQL Database",
	Type:           "sql_database",
	Names:          "SQL Databases",
	Types:          "sql_databases",
	DocsURL:        "https://learn.microsoft.com/fabric/database/sql/overview",
	IsPreview:      true,
	IsSPNSupported: true,
}
