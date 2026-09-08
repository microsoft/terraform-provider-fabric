// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package orgapp

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType = fabcore.ItemTypeOrgApp

	ItemDefinitionEmpty       = `{"$schema":"https://developer.microsoft.com/json-schemas/fabric/item/orgapp/definition/orgAppDefinition/2.0.0/schema.json","elements":[]}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/orgapp-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Org App",
	Type:           "org_app",
	Names:          "Org Apps",
	Types:          "org_apps",
	DocsURL:        "https://learn.microsoft.com/rest/api/fabric/orgapp/items/create-org-app",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  "JSON",
		API:   "JSON",
		Paths: []string{"definition.json", ".platform"},
	},
}
