// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package kqldatabase

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

const (
	ItemName                  = "KQL Database"
	ItemTFName                = "kql_database"
	ItemsName                 = "KQL Databases"
	ItemsTFName               = "kql_databases"
	ItemType                  = fabcore.ItemTypeKQLDatabase
	ItemDocsSPNSupport        = common.DocsSPNSupported
	ItemDocsURL               = "https://learn.microsoft.com/fabric/real-time-intelligence/create-database"
	ItemFormatTypeDefault     = fabricitem.DefinitionFormatDefault
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/kql-database-definition"
)

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"DatabaseProperties.json", "DatabaseSchema.kql"},
	},
}
