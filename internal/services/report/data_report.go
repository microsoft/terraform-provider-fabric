// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package report

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceReport() datasource.DataSource {
	config := fabricitem.DataSourceFabricItemDefinition{
		Type:   ItemType,
		Name:   ItemName,
		TFName: ItemTFName,
		MarkdownDescription: "Get a Fabric " + ItemName + ".\n\n" +
			"Use this data source to fetch a [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
			ItemDocsSPNSupport,
		IsDisplayNameUnique: false,
		DefinitionFormats:   itemDefinitionFormats,
	}

	return fabricitem.NewDataSourceFabricItemDefinition(config)
}
