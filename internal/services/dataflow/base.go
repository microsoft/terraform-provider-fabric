// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package dataflow

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeDataflow
	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/fabric/data-factory/data-source-management"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Dataflow",
	Type:           "dataflow",
	Names:          "Dataflows",
	Types:          "dataflows",
	DocsURL:        "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/dataflow-definition",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"queryMetadata.json", "mashup.pq"},
	},
}
