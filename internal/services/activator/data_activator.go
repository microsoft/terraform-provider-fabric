// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package activator

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceActivator() datasource.DataSource {
	config := fabricitem.DataSourceFabricItemDefinition{
		FabricItemType:      FabricItemType,
		TypeInfo:            ItemTypeInfo,
		IsDisplayNameUnique: true,
		DefinitionFormats:   itemDefinitionFormats,
	}

	return fabricitem.NewDataSourceFabricItemDefinition(config)
}
