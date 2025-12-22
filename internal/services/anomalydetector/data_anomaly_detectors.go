// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package anomalydetector

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceAnomalyDetectors() datasource.DataSource {
	config := fabricitem.DataSourceFabricItems{
		TypeInfo:       ItemTypeInfo,
		FabricItemType: FabricItemType,
	}

	return fabricitem.NewDataSourceFabricItems(config)
}
