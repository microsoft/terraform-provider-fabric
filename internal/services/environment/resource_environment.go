// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package environment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceEnvironment(ctx context.Context) resource.Resource {
	propertiesSetter := func(ctx context.Context, from *fabenvironment.Properties, to *fabricitem.ResourceFabricItemPropertiesModel[environmentPropertiesModel, fabenvironment.Properties]) diag.Diagnostics {
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

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemPropertiesModel[environmentPropertiesModel, fabenvironment.Properties], fabricItem *fabricitem.FabricItemProperties[fabenvironment.Properties]) error {
		client := fabenvironment.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetEnvironment(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.Environment)

		return nil
	}

	config := fabricitem.ResourceFabricItemProperties[environmentPropertiesModel, fabenvironment.Properties]{
		ResourceFabricItem: fabricitem.ResourceFabricItem{
			TypeInfo:             ItemTypeInfo,
			FabricItemType:       FabricItemType,
			NameRenameAllowed:    true,
			DisplayNameMaxLength: 123,
			DescriptionMaxLength: 256,
		},
		PropertiesAttributes: getResourceEnvironmentPropertiesAttributes(ctx),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
	}

	return fabricitem.NewResourceFabricItemProperties(config)
}
