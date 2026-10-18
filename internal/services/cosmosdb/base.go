// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType = fabcore.ItemTypeCosmosDBDatabase

	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/cosmosdb-database-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Cosmos DB",
	Type:           "cosmos_db",
	Names:          "Cosmos DBs",
	Types:          "cosmos_dbs",
	DocsURL:        "https://learn.microsoft.com/fabric/real-time-intelligence/event-streams/add-source-azure-cosmos-db-change-data-capture",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"definition.json"},
	},
}
