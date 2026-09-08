// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package orgapp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	faborgapp "github.com/microsoft/fabric-sdk-go/fabric/orgapp"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceOrgApps() datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *faborgapp.Properties, to *fabricitem.FabricItemPropertiesModel[orgAppPropertiesModel, faborgapp.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[orgAppPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &orgAppPropertiesModel{}
			propertiesModel.set(*from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemsPropertiesModel[orgAppPropertiesModel, faborgapp.Properties], fabricItems *[]fabricitem.FabricItemProperties[faborgapp.Properties]) error {
		client := faborgapp.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		fabItems := make([]fabricitem.FabricItemProperties[faborgapp.Properties], 0)

		respList, err := client.ListOrgApps(ctx, model.WorkspaceID.ValueString(), nil)
		if err != nil {
			return err
		}

		for _, entity := range respList {
			var fabricItem fabricitem.FabricItemProperties[faborgapp.Properties]

			fabricItem.Set(entity)

			fabItems = append(fabItems, fabricItem)
		}

		*fabricItems = fabItems

		return nil
	}

	config := fabricitem.DataSourceFabricItemsProperties[orgAppPropertiesModel, faborgapp.Properties]{
		TypeInfo:             ItemTypeInfo,
		FabricItemType:       FabricItemType,
		PropertiesAttributes: getDataSourceOrgAppPropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemsProperties(config)
}
