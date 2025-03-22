// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mounteddatafactory

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

const (
	ItemName                  = "Mounted Data Factory"
	ItemTFName                = "mounted_data_factory"
	ItemsName                 = "Mounted Data Factories"
	ItemsTFName               = "mounted_data_factories"
	ItemType                  = fabcore.ItemTypeMountedDataFactory
	ItemDocsSPNSupport        = common.DocsSPNSupported
	ItemDocsURL               = "TODO"
	ItemDefinitionEmpty       = `TODO`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/mounted-data-factory-definition"
	ItemPreview               = true
)

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"mountedDataFactory-content.json"},
	},
}
