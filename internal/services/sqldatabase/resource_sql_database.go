// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sqldatabase

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceSQLDatabase() resource.Resource {
	config := fabricitem.ResourceFabricItem{
		TypeInfo:             ItemTypeInfo,
		FabricItemType:       FabricItemType,
		DisplayNameMaxLength: 123,
		DescriptionMaxLength: 256,
	}

	return fabricitem.NewResourceFabricItem(config)
}
