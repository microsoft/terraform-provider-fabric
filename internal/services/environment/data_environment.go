// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package environment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceEnvironment(ctx context.Context) datasource.DataSource {
	propertiesSetter := func(ctx context.Context, from *fabenvironment.PublishInfo, to *fabricitem.DataSourceFabricItemPropertiesModel[environmentPropertiesModel, fabenvironment.PublishInfo]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[environmentPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &environmentPropertiesModel{}

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

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemPropertiesModel[environmentPropertiesModel, fabenvironment.PublishInfo], fabricItem *fabricitem.FabricItemProperties[fabenvironment.PublishInfo]) error {
		client := fabenvironment.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetEnvironment(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.Environment)

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemPropertiesModel[environmentPropertiesModel, fabenvironment.PublishInfo], errNotFound fabcore.ResponseError, fabricItem *fabricitem.FabricItemProperties[fabenvironment.PublishInfo]) error {
		client := fabenvironment.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		pager := client.NewListEnvironmentsPager(model.WorkspaceID.ValueString(), nil)
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

	config := fabricitem.DataSourceFabricItemProperties[environmentPropertiesModel, fabenvironment.PublishInfo]{
		DataSourceFabricItem: fabricitem.DataSourceFabricItem{
			TypeInfo:            ItemTypeInfo,
			FabricItemType:      FabricItemType,
			IsDisplayNameUnique: true,
		},
		PropertiesAttributes: getDataSourceEnvironmentPropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
		ItemListGetter:       itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemProperties(config)
}
