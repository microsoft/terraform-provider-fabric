// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehouse

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabwarehouse "github.com/microsoft/fabric-sdk-go/fabric/warehouse"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceWarehouse() resource.Resource {
	propertiesSetter := func(ctx context.Context, from *fabwarehouse.Properties, to *fabricitem.ResourceFabricItemPropertiesModel[warehousePropertiesModel, fabwarehouse.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[warehousePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &warehousePropertiesModel{}
			propertiesModel.set(from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemPropertiesModel[warehousePropertiesModel, fabwarehouse.Properties], fabricItem *fabricitem.FabricItemProperties[fabwarehouse.Properties]) error {
		client := fabwarehouse.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetWarehouse(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.Warehouse)

		return nil
	}

	config := fabricitem.ResourceFabricItemProperties[warehousePropertiesModel, fabwarehouse.Properties]{
		ResourceFabricItem: fabricitem.ResourceFabricItem{
			Type:              ItemType,
			Name:              ItemName,
			NameRenameAllowed: true,
			TFName:            ItemTFName,
			MarkdownDescription: "Manage a Fabric " + ItemName + ".\n\n" +
				"Use this resource to manage a [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
				ItemDocsSPNSupport,
			DisplayNameMaxLength: 123,
			DescriptionMaxLength: 256,
		},
		PropertiesAttributes: getResourceWarehousePropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
	}

	return fabricitem.NewResourceFabricItemProperties(config)
}
