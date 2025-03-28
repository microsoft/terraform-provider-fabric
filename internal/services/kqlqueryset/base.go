// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package kqlqueryset

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeKQLQueryset
	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/kql-queryset-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "KQL Queryset",
	Type:           "kql_queryset",
	Names:          "KQL Querysets",
	Types:          "kql_querysets",
	DocsURL:        "https://learn.microsoft.com/fabric/real-time-intelligence/kusto-query-set",
	IsPreview:      false,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"RealTimeQueryset.json"},
	},
}
