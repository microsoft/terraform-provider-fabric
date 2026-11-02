// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package userdatafunction

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType = fabcore.ItemTypeUserDataFunction

	ItemDefinitionEmpty       = `{"$schema":"https://developer.microsoft.com/json-schemas/fabric/item/userDataFunction/definition/1.1.0/schema.json","runtime":"PYTHON","connectedDataSources":[],"functions":[],"libraries":{"public":[{"name":"fabric-user-data-functions","type":"PYPI","version":"1.0"}],"private":[]}}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/user-data-function-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "User Data Function",
	Type:           "user_data_function",
	Names:          "User Data Functions",
	Types:          "user_data_functions",
	DocsURL:        "https://learn.microsoft.com/fabric/data-engineering/user-data-functions/user-data-functions-overview",
	IsPreview:      true,
	IsSPNSupported: false,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"definition.json", "function_app.py", "resources/functions.json", "privateLibraries/*.whl"},
	},
}
