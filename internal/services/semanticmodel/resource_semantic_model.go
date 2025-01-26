// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package semanticmodel

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

func NewResourceSemanticModel() resource.Resource {
	config := fabricitem.ResourceFabricItemDefinition{
		Type:              ItemType,
		Name:              ItemName,
		NameRenameAllowed: true,
		TFName:            ItemTFName,
		MarkdownDescription: "Manage a Fabric " + ItemName + ".\n\n" +
			"Use this resource to manage [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
			ItemDocsSPNSupport,
		DisplayNameMaxLength:  123,
		DescriptionMaxLength:  256,
		DefinitionPathDocsURL: ItemDefinitionPathDocsURL,
		DefinitionPathKeysValidator: []validator.Map{
			mapvalidator.SizeAtLeast(2),
			mapvalidator.KeysAre(
				fwvalidators.PatternsIfAttributeIsOneOf(
					path.MatchRoot("format"),
					[]attr.Value{types.StringValue("TMSL")},
					fabricitem.GetDefinitionFormatPaths(itemDefinitionFormats, "TMSL"),
					"Definition path must match one of the following: "+utils.ConvertStringSlicesToString(fabricitem.GetDefinitionFormatPaths(itemDefinitionFormats, "TMSL"), true, false),
				),
				fwvalidators.PatternsIfAttributeIsOneOf(
					path.MatchRoot("format"),
					[]attr.Value{types.StringValue("TMDL")},
					fabricitem.GetDefinitionFormatPaths(itemDefinitionFormats, "TMDL"),
					"Definition path must match one of the following: "+utils.ConvertStringSlicesToString(fabricitem.GetDefinitionFormatPaths(itemDefinitionFormats, "TMDL"), true, false),
				),
			),
		},
		DefinitionRequired: true,
		DefinitionEmpty:    "",
		DefinitionFormats:  itemDefinitionFormats,
	}

	return fabricitem.NewResourceFabricItemDefinition(config)
}
