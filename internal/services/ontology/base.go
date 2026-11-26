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

	ItemDefinitionEmpty       = "database\n\tcompatibilityLevel: 1000000\n"
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
		Type: "TMDL",
		// The API infers TMDL from the parts; an explicit "TMDL" format is rejected.
		API: "",
		Paths: []string{
			"database.tmdl",
			".platform",
			"model.tmdl",
			"tables/*.tmdl",
			"metrics/*.tmdl",
			"relationships.tmdl",
			"expressions.tmdl",
			"entities/*.tmdl",
			"entityRelationships.tmdl",
			"namespaces/*.tmdl",
			"rules/*.tmdl",
		},
	},
}
