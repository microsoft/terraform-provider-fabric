// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package report

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

func NewResourceReport() resource.Resource {
	config := fabricitem.ResourceFabricItemDefinition{
		TypeInfo:              ItemTypeInfo,
		FabricItemType:        FabricItemType,
		NameRenameAllowed:     true,
		DisplayNameMaxLength:  123,
		DescriptionMaxLength:  256,
		DefinitionPathDocsURL: ItemDefinitionPathDocsURL,
		DefinitionPathKeysValidator: []validator.Map{
			mapvalidator.SizeAtLeast(3),
			mapvalidator.KeysAre(
				fwvalidators.PatternsIfAttributeIsOneOf(
					path.MatchRoot("format"),
					[]attr.Value{types.StringValue("PBIR-Legacy")},
					fabricitem.GetDefinitionFormatPaths(itemDefinitionFormats, "PBIR-Legacy"),
					"Definition path must match one of the following: "+utils.ConvertStringSlicesToString(fabricitem.GetDefinitionFormatPaths(itemDefinitionFormats, "PBIR-Legacy"), true, false),
				),
				fwvalidators.PatternsIfAttributeIsOneOf(
					path.MatchRoot("format"),
					[]attr.Value{types.StringValue("PBIR")},
					fabricitem.GetDefinitionFormatPaths(itemDefinitionFormats, "PBIR"),
					"Definition path must match one of the following: "+utils.ConvertStringSlicesToString(fabricitem.GetDefinitionFormatPaths(itemDefinitionFormats, "PBIR"), true, false),
				),
			),
		},
		DefinitionRequired: true,
		DefinitionEmpty:    "",
		DefinitionFormats:  itemDefinitionFormats,
	}

	return fabricitem.NewResourceFabricItemDefinition(config)
}
