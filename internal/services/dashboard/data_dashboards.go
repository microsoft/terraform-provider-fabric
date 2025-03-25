// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package dashboard

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceDashboards() datasource.DataSource {
	config := fabricitem.DataSourceFabricItems{
		Type:                FabricItemType,
		Name:                ItemTypeInfo.Name,
		Names:               ItemTypeInfo.Names,
		TFName:              ItemTypeInfo.Types,
		MarkdownDescription: fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, true),
	}

	return fabricitem.NewDataSourceFabricItems(config)
}
