// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sparkjobdefinition

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabsparkjobdefinition "github.com/microsoft/fabric-sdk-go/fabric/sparkjobdefinition"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceSparkJobDefinitions() datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabsparkjobdefinition.Properties, to *fabricitem.FabricItemPropertiesModel[sparkJobDefinitionPropertiesModel, fabsparkjobdefinition.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[sparkJobDefinitionPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &sparkJobDefinitionPropertiesModel{}
			propertiesModel.set(*from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemsPropertiesModel[sparkJobDefinitionPropertiesModel, fabsparkjobdefinition.Properties], fabricItems *[]fabricitem.FabricItemProperties[fabsparkjobdefinition.Properties]) error {
		client := fabsparkjobdefinition.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		fabItems := make([]fabricitem.FabricItemProperties[fabsparkjobdefinition.Properties], 0)

		respList, err := client.ListSparkJobDefinitions(ctx, model.WorkspaceID.ValueString(), nil)
		if err != nil {
			return err
		}

		for _, entity := range respList {
			var fabricItem fabricitem.FabricItemProperties[fabsparkjobdefinition.Properties]

			fabricItem.Set(entity)

			fabItems = append(fabItems, fabricItem)
		}

		*fabricItems = fabItems

		return nil
	}

	config := fabricitem.DataSourceFabricItemsProperties[sparkJobDefinitionPropertiesModel, fabsparkjobdefinition.Properties]{
		DataSourceFabricItems: fabricitem.DataSourceFabricItems{
			TypeInfo:       ItemTypeInfo,
			FabricItemType: FabricItemType,
		},
		PropertiesAttributes: getDataSourceSparkJobDefinitionPropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemsProperties(config)
}
