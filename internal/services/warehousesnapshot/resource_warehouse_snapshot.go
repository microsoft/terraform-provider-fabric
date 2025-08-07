// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehousesnapshot

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabwarehousesnapshot "github.com/microsoft/fabric-sdk-go/fabric/warehousesnapshot"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceWarehouseSnapshot() resource.Resource {
	creationPayloadSetter := func(_ context.Context, from warehouseSnapshotConfigurationModel) (*fabwarehousesnapshot.CreationPayload, diag.Diagnostics) {
		creationPayload := &fabwarehousesnapshot.CreationPayload{
			ParentWarehouseID: from.ParentWarehouseID.ValueStringPointer(),
			// Handle ValueRFC3339Time() returning (time.Time, diag.Diagnostics)
			SnapshotDateTime: func() *time.Time {
				t, diags := from.SnapshotDateTime.ValueRFC3339Time()
				if diags.HasError() {
					return nil
				}

				return to.Ptr(t)
			}(),
		}

		return creationPayload, nil
	}

	propertiesSetter := func(ctx context.Context, from *fabwarehousesnapshot.Properties, to *fabricitem.ResourceFabricItemConfigPropertiesModel[warehouseSnapshotPropertiesModel, fabwarehousesnapshot.Properties, warehouseSnapshotConfigurationModel, fabwarehousesnapshot.CreationPayload]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[warehouseSnapshotPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &warehouseSnapshotPropertiesModel{}

			propertiesModel.set(*from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemConfigPropertiesModel[warehouseSnapshotPropertiesModel, fabwarehousesnapshot.Properties, warehouseSnapshotConfigurationModel, fabwarehousesnapshot.CreationPayload], fabricItem *fabricitem.FabricItemProperties[fabwarehousesnapshot.Properties]) error {
		client := fabwarehousesnapshot.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetWarehouseSnapshot(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.WarehouseSnapshot)

		return nil
	}

	config := fabricitem.ResourceFabricItemConfigProperties[warehouseSnapshotPropertiesModel, fabwarehousesnapshot.Properties, warehouseSnapshotConfigurationModel, fabwarehousesnapshot.CreationPayload]{
		ResourceFabricItem: fabricitem.ResourceFabricItem{
			TypeInfo:             ItemTypeInfo,
			FabricItemType:       FabricItemType,
			NameRenameAllowed:    true,
			DisplayNameMaxLength: 123,
			DescriptionMaxLength: 256,
		},
		ConfigRequired:        true,
		ConfigAttributes:      getResourceWarehouseSnapshotConfigurationAttributes(),
		CreationPayloadSetter: creationPayloadSetter,
		PropertiesAttributes:  getResourceWarehouseSnapshotPropertiesAttributes(),
		PropertiesSetter:      propertiesSetter,
		ItemGetter:            itemGetter,
	}

	return fabricitem.NewResourceFabricItemConfigProperties(config)
}
