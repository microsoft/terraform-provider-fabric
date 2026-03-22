// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package ontology

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType = fabcore.ItemTypeOntology

	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/ontology-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Ontology",
	Type:           "ontology",
	Names:          "Ontologies",
	Types:          "ontologies",
	DocsURL:        "https://learn.microsoft.com/fabric/iq/ontology/overview",
	IsPreview:      true,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type: fabricitem.DefinitionFormatDefault,
		API:  "",
		Paths: []string{
			"definition.json",
			"EntityTypes/*",
			"EntityTypes/*/DataBindings",
			"EntityTypes/*/Documents",
			"EntityTypes/*/Overviews",
			"RelationshipTypes/*",
			"RelationshipTypes/*/Contextualizations",
		},
	},
}
