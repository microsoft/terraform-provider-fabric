// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventstream

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeEventstream
	ItemDefinitionEmpty       = `{"sources":[],"destinations":[],"streams":[],"operators":[],"compatibilityLevel":"1.0"}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/eventstream-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Eventstream",
	Type:           "eventstream",
	Names:          "Eventstreams",
	Types:          "eventstreams",
	DocsURL:        "https://learn.microsoft.com/fabric/real-time-intelligence/event-streams/overview",
	IsPreview:      true,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"eventstream.json"},
	},
}
