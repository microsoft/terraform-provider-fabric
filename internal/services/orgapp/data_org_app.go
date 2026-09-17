// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package orgapp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	faborgapp "github.com/microsoft/fabric-sdk-go/fabric/orgapp"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceOrgApp() datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *faborgapp.Properties, to *fabricitem.DataSourceFabricItemDefinitionPropertiesModel[orgAppPropertiesModel, faborgapp.Properties]) diag.Diagnostics {
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

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[orgAppPropertiesModel, faborgapp.Properties], fabricItem *fabricitem.FabricItemProperties[faborgapp.Properties]) error {
		client := faborgapp.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetOrgApp(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.OrgApp)

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemDefinitionPropertiesModel[orgAppPropertiesModel, faborgapp.Properties], errNotFound fabcore.ResponseError, fabricItem *fabricitem.FabricItemProperties[faborgapp.Properties]) error {
		client := faborgapp.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		pager := client.NewListOrgAppsPager(model.WorkspaceID.ValueString(), nil)
		for pager.More() {
			page, err := pager.NextPage(ctx)
			if err != nil {
				return err
			}

			for _, entity := range page.Value {
				if *entity.DisplayName == model.DisplayName.ValueString() {
					fabricItem.Set(entity)

					return nil
				}
			}
		}

		return &errNotFound
	}

	config := fabricitem.DataSourceFabricItemDefinitionProperties[orgAppPropertiesModel, faborgapp.Properties]{
		TypeInfo:             ItemTypeInfo,
		FabricItemType:       FabricItemType,
		IsDisplayNameUnique:  true,
		DefinitionFormats:    itemDefinitionFormats,
		PropertiesAttributes: getDataSourceOrgAppPropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemDefinitionProperties(config)
}
