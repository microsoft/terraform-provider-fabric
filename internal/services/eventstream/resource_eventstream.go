// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventstream

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceEventstream() resource.Resource {
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
			mapvalidator.SizeAtMost(1),
			mapvalidator.KeysAre(fabricitem.DefinitionPathKeysValidator(itemDefinitionFormats)...),
		},
		DefinitionRequired: false,
		DefinitionEmpty:    ItemDefinitionEmpty,
		DefinitionFormats:  itemDefinitionFormats,
		IsPreview:          ItemPreview,
	}

	return fabricitem.NewResourceFabricItemDefinition(config)
}
