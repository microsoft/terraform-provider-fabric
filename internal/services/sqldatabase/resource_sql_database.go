// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sqldatabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabsqldatabase "github.com/microsoft/fabric-sdk-go/fabric/sqldatabase"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceSQLDatabase() resource.Resource {
	propertiesSetter := func(ctx context.Context, from *fabsqldatabase.Properties, to *fabricitem.ResourceFabricItemPropertiesModel[sqlDatabasePropertiesModel, fabsqldatabase.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[sqlDatabasePropertiesModel](ctx)

		if from != nil {
			propertiesModel := &sqlDatabasePropertiesModel{}

			propertiesModel.set(*from)

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemPropertiesModel[sqlDatabasePropertiesModel, fabsqldatabase.Properties], fabricItem *fabricitem.FabricItemProperties[fabsqldatabase.Properties]) error {
		client := fabsqldatabase.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetSQLDatabase(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.SQLDatabase)

		return nil
	}

	config := fabricitem.ResourceFabricItemProperties[sqlDatabasePropertiesModel, fabsqldatabase.Properties]{
		ResourceFabricItem: fabricitem.ResourceFabricItem{
			TypeInfo:             ItemTypeInfo,
			FabricItemType:       FabricItemType,
			NameRenameAllowed:    true,
			DisplayNameMaxLength: 123,
			DescriptionMaxLength: 256,
		},
		PropertiesAttributes: getResourceSQLDatabasePropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
	}

	return fabricitem.NewResourceFabricItemProperties(config)
}
