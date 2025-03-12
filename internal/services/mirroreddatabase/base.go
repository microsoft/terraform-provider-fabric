// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mirroreddatabase

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

const (
	ItemName                  = "Mirrored Database"
	ItemTFName                = "mirrored_database"
	ItemsName                 = "Mirrored Databases"
	ItemsTFName               = "mirrored_databases"
	ItemType                  = fabcore.ItemTypeMirroredDatabase
	ItemDocsSPNSupport        = common.DocsSPNSupported
	ItemDocsURL               = "https://learn.microsoft.com/fabric/database/mirrored-database/overview"
	ItemDefinitionEmpty       = `{"properties":{"source":{"type":"GenericMirror"},"target":{"type":"MountedRelationalDatabase","typeProperties":{"defaultSchema":"dbo","format":"Delta"}}}}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/mirrored-database-definition"
)

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"mirroring.json"},
	},
}
