// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package operationsagent

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	faboperationsagent "github.com/microsoft/fabric-sdk-go/fabric/operationsagent"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceOperationsAgents() datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *faboperationsagent.Properties, to *fabricitem.FabricItemPropertiesModel[operationsAgentPropertiesModel, faboperationsagent.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[operationsAgentPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &operationsAgentPropertiesModel{}
			propertiesModel.set(*from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemsPropertiesModel[operationsAgentPropertiesModel, faboperationsagent.Properties], fabricItems *[]fabricitem.FabricItemProperties[faboperationsagent.Properties]) error {
		client := faboperationsagent.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		fabItems := make([]fabricitem.FabricItemProperties[faboperationsagent.Properties], 0)

		respList, err := client.ListOperationsAgents(ctx, model.WorkspaceID.ValueString(), nil)
		if err != nil {
			return err
		}

		for _, entity := range respList {
			var fabricItem fabricitem.FabricItemProperties[faboperationsagent.Properties]

			fabricItem.Set(entity)

			fabItems = append(fabItems, fabricItem)
		}

		*fabricItems = fabItems

		return nil
	}

	config := fabricitem.DataSourceFabricItemsProperties[operationsAgentPropertiesModel, faboperationsagent.Properties]{
		DataSourceFabricItems: fabricitem.DataSourceFabricItems{
			TypeInfo:       ItemTypeInfo,
			FabricItemType: FabricItemType,
		},
		PropertiesAttributes: getDataSourceOperationsAgentPropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemsProperties(config)
}
