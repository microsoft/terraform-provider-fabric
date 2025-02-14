// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package activator

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

const (
	ItemName                  = "Activator"
	ItemTFName                = "activator"
	ItemsName                 = "Activators"
	ItemsTFName               = "activators"
	ItemType                  = fabcore.ItemTypeReflex
	ItemDocsSPNSupport        = common.DocsSPNNotSupported
	ItemDocsURL               = "https://learn.microsoft.com/fabric/real-time-intelligence/event-streams/add-destination-activator"
	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/activator-definition"
	ItemPreview               = true
)

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"ReflexEntities.json"},
	},
}
