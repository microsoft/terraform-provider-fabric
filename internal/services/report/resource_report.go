// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package report

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func NewResourceReport() resource.Resource {
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
		FormatTypeDefault:     ItemFormatTypeDefault,
		FormatTypes:           ItemFormatTypes,
		DefinitionPathDocsURL: ItemDefinitionPathDocsURL,
		DefinitionPathKeys:    ItemDefinitionPathsPBIRLegacy,
		DefinitionPathKeysValidator: []validator.Map{
			mapvalidator.SizeAtLeast(3),
			mapvalidator.KeysAre(
				stringvalidator.RegexMatches(
					regexp.MustCompile(`^(report\.json|definition\.pbir|StaticResources/RegisteredResources/.*|StaticResources/SharedResources/.*)$`),
					"Definition path must match one of the following: "+utils.ConvertStringSlicesToString(ItemDefinitionPathsPBIRLegacy, true, false),
				),
			),
		},
		DefinitionRequired: true,
		DefinitionEmpty:    "",
	}

	return fabricitem.NewResourceFabricItemDefinition(config)
}