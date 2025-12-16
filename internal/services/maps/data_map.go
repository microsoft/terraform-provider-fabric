// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package maps

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceMap() datasource.DataSource {
	config := fabricitem.DataSourceFabricItemDefinition{
		TypeInfo:            ItemTypeInfo,
		FabricItemType:      FabricItemType,
		IsDisplayNameUnique: true,
		DefinitionFormats:   itemDefinitionFormats,
	}

	return fabricitem.NewDataSourceFabricItemDefinition(config)
}
