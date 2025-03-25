// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mounteddatafactory

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceMountedDataFactories() datasource.DataSource {
	config := fabricitem.DataSourceFabricItems{
		Type:   ItemType,
		Name:   ItemName,
		Names:  ItemsName,
		TFName: ItemsTFName,
		MarkdownDescription: "List a Fabric " + ItemsName + ".\n\n" +
			"Use this data source to list [" + ItemsName + "](" + ItemDocsURL + ").\n\n" +
			ItemDocsSPNSupport,
	}

	return fabricitem.NewDataSourceFabricItems(config)
}
