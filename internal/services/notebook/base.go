// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package notebook

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeNotebook
	ItemDefinitionEmptyIPYNB  = `{"cells":[{"cell_type":"code","metadata":{},"source":["# Welcome to your notebook"]}],"metadata":{"language_info":{"name":"python"}},"nbformat":4,"nbformat_minor":5}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/notebook-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Notebook",
	Type:           "notebook",
	Names:          "Notebooks",
	Types:          "notebooks",
	DocsURL:        "https://learn.microsoft.com/fabric/data-engineering/how-to-use-notebook",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  "ipynb",
		API:   "ipynb",
		Paths: []string{"notebook-content.ipynb"},
	},
	{
		Type:  "py",
		API:   "",
		Paths: []string{"notebook-content.py"},
	},
}
