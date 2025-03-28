// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package graphqlapi

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceGraphQLApi() resource.Resource {
	config := fabricitem.ResourceFabricItem{
		TypeInfo:             ItemTypeInfo,
		FabricItemType:       FabricItemType,
		NameRenameAllowed:    true,
		DisplayNameMaxLength: 123,
		DescriptionMaxLength: 256,
	}

	return fabricitem.NewResourceFabricItem(config)
}
