// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sqldatabase

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceSQLDatabase() datasource.DataSource {
	config := fabricitem.DataSourceFabricItem{
		Type:   ItemType,
		Name:   ItemName,
		TFName: ItemTFName,
		MarkdownDescription: "Get a Fabric " + ItemName + ".\n\n" +
			"Use this data source to fetch a [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
			ItemDocsSPNSupport,
		IsDisplayNameUnique: true,
		IsPreview:           ItemPreview,
	}

	return fabricitem.NewDataSourceFabricItem(config)
}
