// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package kqldatabase

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeKQLDatabase
	ItemFormatTypeDefault     = fabricitem.DefinitionFormatDefault
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/kql-database-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "KQL Database",
	Type:           "kql_database",
	Names:          "KQL Databases",
	Types:          "kql_databases",
	DocsURL:        "https://learn.microsoft.com/fabric/real-time-intelligence/create-database",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"DatabaseProperties.json", "DatabaseSchema.kql"},
	},
}
