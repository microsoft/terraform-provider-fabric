// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeLakehouse
	ItemFormatTypeDefault     = fabricitem.DefinitionFormatDefault
	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/lakehouse-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Lakehouse",
	Type:           "lakehouse",
	Names:          "Lakehouses",
	Types:          "lakehouses",
	DocsURL:        "https://learn.microsoft.com/training/modules/get-started-lakehouses",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"lakehouse.metadata.json", "shortcuts.metadata.json", "data-access-roles.json", "alm.settings.json"},
	},
}
