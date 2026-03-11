// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package datapipeline

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeDataPipeline
	ItemDefinitionEmpty       = `{"properties":{"activities":[]}}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/fabric/data-factory/pipeline-rest-api"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Data Pipeline",
	Type:           "data_pipeline",
	Names:          "Data Pipelines",
	Types:          "data_pipelines",
	DocsURL:        "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/datapipeline-definition",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"pipeline-content.json"},
	},
}
