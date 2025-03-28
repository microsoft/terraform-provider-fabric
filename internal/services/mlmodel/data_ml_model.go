// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mlmodel

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceMLModel() datasource.DataSource {
	config := fabricitem.DataSourceFabricItem{
		TypeInfo:            ItemTypeInfo,
		FabricItemType:      FabricItemType,
		IsDisplayNameUnique: true,
	}

	return fabricitem.NewDataSourceFabricItem(config)
}
