// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mirroredazuredatabrickscatalog

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType = fabcore.ItemTypeMirroredAzureDatabricksCatalog

	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/mirrored-azuredatabricks-unitycatalog-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Mirrored Azure Databricks Catalog",
	Type:           "mirrored_azure_databricks_catalog",
	Names:          "Mirrored Azure Databricks Catalogs",
	Types:          "mirrored_azure_databricks_catalogs",
	DocsURL:        "https://learn.microsoft.com/fabric/database/mirrored-database/azure-databricks",
	IsPreview:      true,
	IsSPNSupported: false,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"mirroringAzureDatabricksCatalog.json"},
	},
}
