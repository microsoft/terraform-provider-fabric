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
	creationPayloadSetter := func(_ context.Context, from warehouseConfigurationModel) (*fabwarehouse.CreationPayload, diag.Diagnostics) {
		cp := fabwarehouse.CreationPayload{
			CollationType: (*fabwarehouse.CollationType)(from.CollationType.ValueStringPointer()),
		}

		return &cp, nil
	}

	propertiesSetter := func(ctx context.Context, from *fabwarehouse.Properties, to *fabricitem.ResourceFabricItemConfigPropertiesModel[warehousePropertiesModel, fabwarehouse.Properties, warehouseConfigurationModel, fabwarehouse.CreationPayload]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[warehousePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &warehousePropertiesModel{}
			propertiesModel.set(*from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemConfigPropertiesModel[warehousePropertiesModel, fabwarehouse.Properties, warehouseConfigurationModel, fabwarehouse.CreationPayload], fabricItem *fabricitem.FabricItemProperties[fabwarehouse.Properties]) error {
		client := fabwarehouse.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetWarehouse(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.Warehouse)

		return nil
	}

	config := fabricitem.ResourceFabricItemConfigProperties[warehousePropertiesModel, fabwarehouse.Properties, warehouseConfigurationModel, fabwarehouse.CreationPayload]{
		ResourceFabricItem: fabricitem.ResourceFabricItem{
			TypeInfo:             ItemTypeInfo,
			FabricItemType:       FabricItemType,
			NameRenameAllowed:    true,
			DisplayNameMaxLength: 123,
			DescriptionMaxLength: 256,
		},
		PropertiesAttributes:  getResourceWarehousePropertiesAttributes(),
		PropertiesSetter:      propertiesSetter,
		ItemGetter:            itemGetter,
		ConfigRequired:        false,
		ConfigAttributes:      getResourceWarehouseConfigurationAttributes(),
		CreationPayloadSetter: creationPayloadSetter,
	}

	return fabricitem.NewResourceFabricItemConfigProperties(config)
}
