// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package notebook

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName                  = "Notebook"
	ItemTFName                = "notebook"
	ItemsName                 = "Notebooks"
	ItemsTFName               = "notebooks"
	ItemType                  = fabcore.ItemTypeNotebook
	ItemDocsSPNSupport        = common.DocsSPNSupported
	ItemDocsURL               = "https://learn.microsoft.com/fabric/data-engineering/how-to-use-notebook"
	ItemFormatTypeDefault     = "ipynb"
	ItemDefinitionEmptyIPYNB  = `{"cells":[{"cell_type":"code","metadata":{},"source":["# Welcome to your notebook"]}],"metadata":{"language_info":{"name":"python"}},"nbformat":4,"nbformat_minor":5}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/notebook-definition"
)

var (
	ItemFormatTypes          = []string{"ipynb"}                  //nolint:gochecknoglobals
	ItemDefinitionPathsIPYNB = []string{"notebook-content.ipynb"} //nolint:gochecknoglobals
)
