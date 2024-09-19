// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package mlmodel

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceMLModel() resource.Resource {
	config := fabricitem.ResourceFabricItem{
		Type:              ItemType,
		Name:              ItemName,
		NameRenameAllowed: false,
		TFName:            ItemTFName,
		MarkdownDescription: "Manage a Fabric " + ItemName + ".\n\n" +
			"Use this resource to manage [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
			ItemDocsSPNSupport,
		DisplayNameMaxLength: 123,
		DescriptionMaxLength: 256,
	}

	return fabricitem.NewResourceFabricItem(config)
}
