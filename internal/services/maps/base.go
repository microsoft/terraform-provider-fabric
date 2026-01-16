// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package maps

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType = fabcore.ItemTypeMap

	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/map-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Map",
	Type:           "map",
	Names:          "Maps",
	Types:          "maps",
	DocsURL:        "https://learn.microsoft.com/fabric/real-time-intelligence/map/create-map",
	IsPreview:      true,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"map.json"},
	},
}
