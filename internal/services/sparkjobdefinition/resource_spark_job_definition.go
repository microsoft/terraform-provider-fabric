// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sparkjobdefinition

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabsparkjobdefinition "github.com/microsoft/fabric-sdk-go/fabric/sparkjobdefinition"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceSparkJobDefinition() resource.Resource {
	propertiesSetter := func(ctx context.Context, from *fabsparkjobdefinition.Properties, to *fabricitem.ResourceFabricItemDefinitionPropertiesModel[sparkJobDefinitionPropertiesModel, fabsparkjobdefinition.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[sparkJobDefinitionPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &sparkJobDefinitionPropertiesModel{}
			propertiesModel.set(*from)

			diags := properties.Set(ctx, propertiesModel)
			if diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemDefinitionPropertiesModel[sparkJobDefinitionPropertiesModel, fabsparkjobdefinition.Properties], fabricItem *fabricitem.FabricItemProperties[fabsparkjobdefinition.Properties]) error {
		client := fabsparkjobdefinition.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetSparkJobDefinition(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.SparkJobDefinition)

		return nil
	}

	config := fabricitem.ResourceFabricItemDefinitionProperties[sparkJobDefinitionPropertiesModel, fabsparkjobdefinition.Properties]{
		ResourceFabricItemDefinition: fabricitem.ResourceFabricItemDefinition{
			TypeInfo:              ItemTypeInfo,
			FabricItemType:        FabricItemType,
			NameRenameAllowed:     true,
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
		},
		PropertiesAttributes: getResourceSparkJobDefinitionPropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
	}

	return fabricitem.NewResourceFabricItemDefinitionProperties(config)
}
