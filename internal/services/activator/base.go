// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package activator

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/reflex-definition"
)

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"ReflexEntities.json"},
	},
}

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Activator",
	Type:           "activator",
	Names:          "Activators",
	Types:          "activators",
	DocsURL:        "https://learn.microsoft.com/fabric/real-time-intelligence/event-streams/add-destination-activator",
	IsPreview:      true,
	IsSPNSupported: false,
}

const FabricItemType = fabcore.ItemTypeReflex
