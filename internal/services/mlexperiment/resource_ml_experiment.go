// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package mlexperiment

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceMLExperiment() resource.Resource {
	config := fabricitem.ResourceFabricItem{
		TypeInfo:             ItemTypeInfo,
		FabricItemType:       FabricItemType,
		NameRenameAllowed:    true,
		DisplayNameMaxLength: 123,
		DescriptionMaxLength: 256,
	}

	return fabricitem.NewResourceFabricItem(config)
}
