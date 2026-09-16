// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package paginatedreport

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	fwvalidators "github.com/microsoft/terraform-provider-fabric/internal/framework/validators"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func NewResourcePaginatedReport() resource.Resource {
	definitionPaths := fabricitem.GetDefinitionFormatPaths(itemDefinitionFormats, fabricitem.DefinitionFormatDefault)
	config := fabricitem.ResourceFabricItemDefinition{
		TypeInfo:              ItemTypeInfo,
		FabricItemType:        FabricItemType,
		NameRenameAllowed:     true,
		DisplayNameMaxLength:  123,
		DescriptionMaxLength:  256,
		DefinitionPathDocsURL: ItemDefinitionPathDocsURL,
		DefinitionPathKeysValidator: []validator.Map{
			mapvalidator.SizeBetween(1, 1),
			mapvalidator.KeysAre(
				fwvalidators.PatternsIfAttributeIsOneOf(
					path.MatchRoot("format"),
					[]attr.Value{types.StringValue(fabricitem.DefinitionFormatDefault)},
					definitionPaths,
					"Definition path must match: "+utils.ConvertStringSlicesToString(definitionPaths, true, false),
				),
			),
		},
		DefinitionRequired: true,
		DefinitionEmpty:    "",
		DefinitionFormats:  itemDefinitionFormats,
	}

	return fabricitem.NewResourceFabricItemDefinition(config)
}
