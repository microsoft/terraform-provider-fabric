// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventhouse

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabeventhouse "github.com/microsoft/fabric-sdk-go/fabric/eventhouse"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceEventhouse(ctx context.Context) resource.Resource {
	propertiesSetter := func(ctx context.Context, from *fabeventhouse.Properties, to *fabricitem.ResourceFabricItemDefinitionPropertiesModel[eventhousePropertiesModel, fabeventhouse.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[eventhousePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &eventhousePropertiesModel{}
			propertiesModel.set(ctx, from)

			diags := properties.Set(ctx, propertiesModel)
			if diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemDefinitionPropertiesModel[eventhousePropertiesModel, fabeventhouse.Properties], fabricItem *fabricitem.FabricItemProperties[fabeventhouse.Properties]) error {
		client := fabeventhouse.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetEventhouse(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.Eventhouse)

		return nil
	}

	config := fabricitem.ResourceFabricItemDefinitionProperties[eventhousePropertiesModel, fabeventhouse.Properties]{
		ResourceFabricItemDefinition: fabricitem.ResourceFabricItemDefinition{
			Type:              ItemType,
			Name:              ItemName,
			NameRenameAllowed: true,
			TFName:            ItemTFName,
			MarkdownDescription: "Manage a Fabric " + ItemName + ".\n\n" +
				"Use this resource to manage an [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
				ItemDocsSPNSupport,
			DisplayNameMaxLength:  123,
			DescriptionMaxLength:  256,
			FormatTypeDefault:     ItemFormatTypeDefault,
			FormatTypes:           ItemFormatTypes,
			DefinitionPathDocsURL: ItemDefinitionPathDocsURL,
			DefinitionPathKeys:    ItemDefinitionPaths,
			DefinitionPathKeysValidator: []validator.Map{
				mapvalidator.SizeAtMost(1),
				mapvalidator.KeysAre(stringvalidator.OneOf(ItemDefinitionPaths...)),
			},
			DefinitionRequired: false,
			DefinitionEmpty:    ItemDefinitionEmpty,
		},
		PropertiesAttributes: getResourceEventhousePropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
	}

	return fabricitem.NewResourceFabricItemDefinitionProperties(config)
}
