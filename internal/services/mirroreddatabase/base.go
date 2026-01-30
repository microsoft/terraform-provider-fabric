// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package mirroreddatabase

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeMirroredDatabase
	ItemDefinitionEmpty       = `{"properties":{"source":{"type":"GenericMirror"},"target":{"type":"MountedRelationalDatabase","typeProperties":{"defaultSchema":"dbo","format":"Delta"}}}}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/mirrored-database-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Mirrored Database",
	Type:           "mirrored_database",
	Names:          "Mirrored Databases",
	Types:          "mirrored_databases",
	DocsURL:        "https://learn.microsoft.com/fabric/database/mirrored-database/overview",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"mirroring.json"},
	},
}
