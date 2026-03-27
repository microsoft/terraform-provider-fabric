// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sqldatabase

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeSQLDatabase
	ItemDocsURL               = "https://learn.microsoft.com/fabric/database/sql/overview"
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/sql-database-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "SQL Database",
	Type:           "sql_database",
	Names:          "SQL Databases",
	Types:          "sql_databases",
	DocsURL:        "https://learn.microsoft.com/fabric/database/sql/overview",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  "dacpac",
		API:   "dacpac",
		Paths: []string{"*.dacpac"},
	},
	{
		Type:  "sqlproj",
		API:   "sqlproj",
		Paths: []string{"*.sqlproj", ".sharedqueries/*.sql", "**/*.sql"},
	},
}
