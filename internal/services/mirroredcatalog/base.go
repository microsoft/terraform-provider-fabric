// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package mirroredcatalog

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeMirroredCatalog
	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/mirrored-catalog-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Mirrored Catalog",
	Type:           "mirrored_catalog",
	Names:          "Mirrored Catalogs",
	Types:          "mirrored_catalogs",
	DocsURL:        "https://learn.microsoft.com/fabric/database/mirrored-catalog/overview",
	IsPreview:      true,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"MirroredCatalogDefinition.json"},
	},
}
