// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventstream

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

const (
	ItemName                  = "Eventstream"
	ItemTFName                = "eventstream"
	ItemsName                 = "Eventstreams"
	ItemsTFName               = "eventstreams"
	ItemType                  = fabcore.ItemTypeEventstream
	ItemDocsSPNSupport        = common.DocsSPNSupported
	ItemDocsURL               = "https://learn.microsoft.com/fabric/real-time-intelligence/event-streams/overview"
	ItemDefinitionEmpty       = `{"sources":[],"destinations":[],"streams":[],"operators":[],"compatibilityLevel":"1.0"}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/eventstream-definition"
	ItemPreview               = true
)

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"eventstream.json"},
	},
}
