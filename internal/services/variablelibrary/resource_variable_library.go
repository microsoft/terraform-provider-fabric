// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package variablelibrary

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabvariablelibrary "github.com/microsoft/fabric-sdk-go/fabric/variablelibrary"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceVariableLibrary() resource.Resource {
	propertiesSetter := func(ctx context.Context, from *fabvariablelibrary.Properties, to *fabricitem.ResourceFabricItemDefinitionPropertiesModel[variableLibraryPropertiesModel, fabvariablelibrary.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[variableLibraryPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &variableLibraryPropertiesModel{}
			propertiesModel.set(*from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemDefinitionPropertiesModel[variableLibraryPropertiesModel, fabvariablelibrary.Properties], fabricItem *fabricitem.FabricItemProperties[fabvariablelibrary.Properties]) error {
		client := fabvariablelibrary.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetVariableLibrary(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.VariableLibrary)

		return nil
	}

	config := fabricitem.ResourceFabricItemDefinitionProperties[variableLibraryPropertiesModel, fabvariablelibrary.Properties]{
		ResourceFabricItemDefinition: fabricitem.ResourceFabricItemDefinition{
			TypeInfo:              ItemTypeInfo,
			FabricItemType:        FabricItemType,
			NameRenameAllowed:     true,
			DisplayNameMaxLength:  123,
			DescriptionMaxLength:  256,
			DefinitionPathDocsURL: ItemDefinitionPathDocsURL,
			DefinitionFormats:     itemDefinitionFormats,
			DefinitionPathKeysValidator: []validator.Map{
				mapvalidator.SizeAtLeast(3),
				mapvalidator.SizeAtMost(3),
				mapvalidator.KeysAre(fabricitem.DefinitionPathKeysValidator(itemDefinitionFormats)...),
			},
			DefinitionRequired: false,
			DefinitionEmpty:    "",
		},
		PropertiesAttributes: getResourceVariableLibraryPropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
	}

	return fabricitem.NewResourceFabricItemDefinitionProperties(config)
}
