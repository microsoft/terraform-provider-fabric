// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package monitoring

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabmonitoring "github.com/microsoft/fabric-sdk-go/fabric/monitoring"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceMonitoringItems() datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabmonitoring.Properties, to *fabricitem.FabricItemPropertiesModel[monitoringPropertiesModel, fabmonitoring.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[monitoringPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &monitoringPropertiesModel{}

			if diags := propertiesModel.set(ctx, *from); diags.HasError() {
				return diags
			}

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemsPropertiesModel[monitoringPropertiesModel, fabmonitoring.Properties], fabricItems *[]fabricitem.FabricItemProperties[fabmonitoring.Properties]) error {
		client := fabmonitoring.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		fabItems := make([]fabricitem.FabricItemProperties[fabmonitoring.Properties], 0)

		pager := client.NewListMonitoringItemsPager(model.WorkspaceID.ValueString(), nil)
		for pager.More() {
			page, err := pager.NextPage(ctx)
			if err != nil {
				return err
			}

			for _, entity := range page.Value {
				var fabricItem fabricitem.FabricItemProperties[fabmonitoring.Properties]
				fabricItem.Set(entity)
				fabItems = append(fabItems, fabricItem)
			}
		}

		*fabricItems = fabItems

		return nil
	}

	config := fabricitem.DataSourceFabricItemsProperties[monitoringPropertiesModel, fabmonitoring.Properties]{
		TypeInfo:             ItemTypeInfo,
		FabricItemType:       FabricItemType,
		PropertiesAttributes: getDataSourceMonitoringPropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemsProperties(config)
}
