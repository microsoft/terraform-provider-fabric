// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mounteddatafactory

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType      = fabcore.ItemTypeMountedDataFactory
	ItemDefinitionEmpty = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/mounted-data-factory-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Mounted Data Factory",
	Type:           "mounted_data_factory",
	Names:          "Mounted Data Factories",
	Types:          "mounted_data_factories",
	DocsURL:        "https://learn.microsoft.com/fabric/data-factory/data-factory-overview",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"mountedDataFactory-content.json"},
	},
}
