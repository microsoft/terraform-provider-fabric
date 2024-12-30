// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package environment

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceEnvironment(ctx context.Context) resource.Resource {
	propertiesSetter := func(ctx context.Context, from *fabenvironment.PublishInfo, to *fabricitem.ResourceFabricItemPropertiesModel[environmentPropertiesModel, fabenvironment.PublishInfo]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[environmentPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &environmentPropertiesModel{}

			if diags := propertiesModel.set(ctx, from); diags.HasError() {
				return diags
			}

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemPropertiesModel[environmentPropertiesModel, fabenvironment.PublishInfo], fabricItem *fabricitem.FabricItemProperties[fabenvironment.PublishInfo]) error {
		client := fabenvironment.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetEnvironment(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.Environment)

		return nil
	}

	config := fabricitem.ResourceFabricItemProperties[environmentPropertiesModel, fabenvironment.PublishInfo]{
		ResourceFabricItem: fabricitem.ResourceFabricItem{
			Type:              ItemType,
			Name:              ItemName,
			NameRenameAllowed: true,
			TFName:            ItemTFName,
			MarkdownDescription: "Manage a Fabric " + ItemName + ".\n\n" +
				"Use this resource to manage an [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
				ItemDocsSPNSupport,
			DisplayNameMaxLength: 123,
			DescriptionMaxLength: 256,
		},
		PropertiesAttributes: getResourceEnvironmentPropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
	}

	return fabricitem.NewResourceFabricItemProperties(config)
}
