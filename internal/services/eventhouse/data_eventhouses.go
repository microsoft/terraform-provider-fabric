// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventhouse

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabeventhouse "github.com/microsoft/fabric-sdk-go/fabric/eventhouse"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceEventhouses(ctx context.Context) datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabeventhouse.Properties, to *fabricitem.FabricItemPropertiesModel[eventhousePropertiesModel, fabeventhouse.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[eventhousePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &eventhousePropertiesModel{}
			propertiesModel.set(ctx, from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemsPropertiesModel[eventhousePropertiesModel, fabeventhouse.Properties], fabricItems *[]fabricitem.FabricItemProperties[fabeventhouse.Properties]) error {
		client := fabeventhouse.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		fabItems := make([]fabricitem.FabricItemProperties[fabeventhouse.Properties], 0)

		respList, err := client.ListEventhouses(ctx, model.WorkspaceID.ValueString(), nil)
		if err != nil {
			return err
		}

		for _, entity := range respList {
			var fabricItem fabricitem.FabricItemProperties[fabeventhouse.Properties]

			fabricItem.Set(entity)

			fabItems = append(fabItems, fabricItem)
		}

		*fabricItems = fabItems

		return nil
	}

	config := fabricitem.DataSourceFabricItemsProperties[eventhousePropertiesModel, fabeventhouse.Properties]{
		DataSourceFabricItems: fabricitem.DataSourceFabricItems{
			Type:   ItemType,
			Name:   ItemName,
			Names:  ItemsName,
			TFName: ItemsTFName,
			MarkdownDescription: "List a Fabric " + ItemsName + ".\n\n" +
				"Use this data source to list [" + ItemsName + "](" + ItemDocsURL + ").\n\n" +
				ItemDocsSPNSupport,
		},
		PropertiesAttributes: getDataSourceEventhousePropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemsProperties(config)
}
