// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package mlexperiment

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceMLExperiment() datasource.DataSource {
	config := fabricitem.DataSourceFabricItem{
		TypeInfo:            ItemTypeInfo,
		FabricItemType:      FabricItemType,
		IsDisplayNameUnique: true,
	}

	return fabricitem.NewDataSourceFabricItem(config)
}
