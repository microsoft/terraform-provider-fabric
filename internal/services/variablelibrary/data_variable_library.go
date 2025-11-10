// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package variablelibrary

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabvariablelibrary "github.com/microsoft/fabric-sdk-go/fabric/variablelibrary"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceVariableLibrary() datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabvariablelibrary.Properties, to *fabricitem.DataSourceFabricItemDefinitionPropertiesModel[variableLibraryPropertiesModel, fabvariablelibrary.Properties]) diag.Diagnostics {
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

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[variableLibraryPropertiesModel, fabvariablelibrary.Properties], fabricItem *fabricitem.FabricItemProperties[fabvariablelibrary.Properties]) error {
		client := fabvariablelibrary.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetVariableLibrary(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.VariableLibrary)

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[variableLibraryPropertiesModel, fabvariablelibrary.Properties], errNotFound fabcore.ResponseError, fabricItem *fabricitem.FabricItemProperties[fabvariablelibrary.Properties]) error {
		client := fabvariablelibrary.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		pager := client.NewListVariableLibrariesPager(model.WorkspaceID.ValueString(), nil)
		for pager.More() {
			page, err := pager.NextPage(ctx)
			if err != nil {
				return err
			}

			for _, entity := range page.Value {
				if *entity.DisplayName == model.DisplayName.ValueString() {
					fabricItem.Set(entity)

					return nil
				}
			}
		}

		return &errNotFound
	}

	config := fabricitem.DataSourceFabricItemDefinitionProperties[variableLibraryPropertiesModel, fabvariablelibrary.Properties]{
		DataSourceFabricItemDefinition: fabricitem.DataSourceFabricItemDefinition{
			TypeInfo:            ItemTypeInfo,
			FabricItemType:      FabricItemType,
			IsDisplayNameUnique: true,
			DefinitionFormats:   itemDefinitionFormats,
		},
		PropertiesAttributes: getDataSourceVariableLibraryPropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemDefinitionProperties(config)
}
