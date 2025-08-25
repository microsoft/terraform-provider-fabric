// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package variablelibrary

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabvariablelibrary "github.com/microsoft/fabric-sdk-go/fabric/variablelibrary"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceVariableLibraries() datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabvariablelibrary.Properties, to *fabricitem.FabricItemPropertiesModel[variableLibraryPropertiesModel, fabvariablelibrary.Properties]) diag.Diagnostics {
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

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemsPropertiesModel[variableLibraryPropertiesModel, fabvariablelibrary.Properties], fabricItems *[]fabricitem.FabricItemProperties[fabvariablelibrary.Properties]) error {
		client := fabvariablelibrary.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		fabItems := make([]fabricitem.FabricItemProperties[fabvariablelibrary.Properties], 0)

		respList, err := client.ListVariableLibraries(ctx, model.WorkspaceID.ValueString(), nil)
		if err != nil {
			return err
		}

		for _, entity := range respList {
			var fabricItem fabricitem.FabricItemProperties[fabvariablelibrary.Properties]

			fabricItem.Set(entity)

			fabItems = append(fabItems, fabricItem)
		}

		*fabricItems = fabItems

		return nil
	}

	config := fabricitem.DataSourceFabricItemsProperties[variableLibraryPropertiesModel, fabvariablelibrary.Properties]{
		DataSourceFabricItems: fabricitem.DataSourceFabricItems{
			TypeInfo:       ItemTypeInfo,
			FabricItemType: FabricItemType,
		},
		PropertiesAttributes: getDataSourceVariableLibraryPropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemsProperties(config)
}
