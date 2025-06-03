// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package digitaltwinbuilder

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeDigitalTwinBuilder
	ItemFormatTypeDefault     = fabricitem.DefinitionFormatDefault
	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/digital-twin-builder-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Digital Twin Builder",
	Type:           "digital_twin_builder",
	Names:          "Digital Twin Builders",
	Types:          "digital_twin_builders",
	DocsURL:        "https://learn.microsoft.com/fabric/real-time-intelligence/digital-twin-builder/overview",
	IsPreview:      true,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"digitaltwinbuilder.json"},
	},
}
