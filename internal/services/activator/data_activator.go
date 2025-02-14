// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package activator

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceReflex() datasource.DataSource {
	config := fabricitem.DataSourceFabricItemDefinition{
		Type:   ItemType,
		Name:   ItemName,
		TFName: ItemTFName,
		MarkdownDescription: "Get a Fabric " + ItemName + " (also known as Reflex).\n\n" +
			"Use this data source to fetch a [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
			ItemDocsSPNSupport,
		IsDisplayNameUnique: true,
		DefinitionFormats:   itemDefinitionFormats,
		IsPreview:           ItemPreview,
	}

	return fabricitem.NewDataSourceFabricItemDefinition(config)
}
