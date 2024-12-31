// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package environment

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceEnvironments(ctx context.Context) datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabenvironment.PublishInfo, to *fabricitem.FabricItemPropertiesModel[environmentPropertiesModel, fabenvironment.PublishInfo]) diag.Diagnostics {
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

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemsPropertiesModel[environmentPropertiesModel, fabenvironment.PublishInfo], fabricItems *[]fabricitem.FabricItemProperties[fabenvironment.PublishInfo]) error {
		client := fabenvironment.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		fabItems := make([]fabricitem.FabricItemProperties[fabenvironment.PublishInfo], 0)

		respList, err := client.ListEnvironments(ctx, model.WorkspaceID.ValueString(), nil)
		if err != nil {
			return err
		}

		for _, entity := range respList {
			var fabricItem fabricitem.FabricItemProperties[fabenvironment.PublishInfo]

			fabricItem.Set(entity)

			fabItems = append(fabItems, fabricItem)
		}

		*fabricItems = fabItems

		return nil
	}

	config := fabricitem.DataSourceFabricItemsProperties[environmentPropertiesModel, fabenvironment.PublishInfo]{
		DataSourceFabricItems: fabricitem.DataSourceFabricItems{
			Type:   ItemType,
			Name:   ItemName,
			Names:  ItemsName,
			TFName: ItemsTFName,
			MarkdownDescription: "List a Fabric " + ItemsName + ".\n\n" +
				"Use this data source to list [" + ItemsName + "](" + ItemDocsURL + ").\n\n" +
				ItemDocsSPNSupport,
		},
		PropertiesAttributes: getDataSourceEnvironmentPropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemsProperties(config)
}
