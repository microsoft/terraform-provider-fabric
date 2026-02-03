// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package mounteddatafactory

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceMountedDataFactory() resource.Resource {
	config := fabricitem.ResourceFabricItemDefinition{
		TypeInfo:              ItemTypeInfo,
		FabricItemType:        FabricItemType,
		NameRenameAllowed:     true,
		DisplayNameMaxLength:  123,
		DescriptionMaxLength:  256,
		DefinitionRequired:    true,
		DefinitionPathDocsURL: ItemDefinitionPathDocsURL,
		DefinitionPathKeysValidator: []validator.Map{
			mapvalidator.SizeAtMost(1),
			mapvalidator.KeysAre(fabricitem.DefinitionPathKeysValidator(itemDefinitionFormats)...),
		},
		DefinitionEmpty:   ItemDefinitionEmpty,
		DefinitionFormats: itemDefinitionFormats,
	}

	return fabricitem.NewResourceFabricItemDefinition(config)
}
