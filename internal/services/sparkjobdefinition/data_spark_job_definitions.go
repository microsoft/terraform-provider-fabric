// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sparkjobdefinition

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabsparkjobdefinition "github.com/microsoft/fabric-sdk-go/fabric/sparkjobdefinition"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewDataSourceSparkJobDefinitions(ctx context.Context) datasource.DataSource {
	propertiesSchema := schema.SingleNestedAttribute{
		MarkdownDescription: "The " + ItemName + " properties.",
		Computed:            true,
		CustomType:          supertypes.NewSingleNestedObjectTypeOf[sparkJobDefinitionPropertiesModel](ctx),
		Attributes: map[string]schema.Attribute{
			"onelake_root_path": schema.StringAttribute{
				MarkdownDescription: "OneLake path to the Spark Job Definition root directory.",
				Computed:            true,
			},
		},
	}

	propertiesSetter := func(ctx context.Context, from *fabsparkjobdefinition.Properties, to *fabricitem.FabricItemPropertiesModel[sparkJobDefinitionPropertiesModel, fabsparkjobdefinition.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[sparkJobDefinitionPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &sparkJobDefinitionPropertiesModel{}
			propertiesModel.set(from)

			diags := properties.Set(ctx, propertiesModel)
			if diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemListGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.DataSourceFabricItemsPropertiesModel[sparkJobDefinitionPropertiesModel, fabsparkjobdefinition.Properties], fabricItems *[]fabricitem.FabricItemProperties[fabsparkjobdefinition.Properties]) error {
		client := fabsparkjobdefinition.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		fabItems := make([]fabricitem.FabricItemProperties[fabsparkjobdefinition.Properties], 0)

		respList, err := client.ListSparkJobDefinitions(ctx, model.WorkspaceID.ValueString(), nil)
		if err != nil {
			return err
		}

		for _, entity := range respList {
			var fabricItem fabricitem.FabricItemProperties[fabsparkjobdefinition.Properties]

			fabricItem.Set(entity)

			fabItems = append(fabItems, fabricItem)
		}

		*fabricItems = fabItems

		return nil
	}

	config := fabricitem.DataSourceFabricItemsProperties[sparkJobDefinitionPropertiesModel, fabsparkjobdefinition.Properties]{
		DataSourceFabricItems: fabricitem.DataSourceFabricItems{
			Type:   ItemType,
			Name:   ItemName,
			Names:  ItemsName,
			TFName: ItemsTFName,
			MarkdownDescription: "List a Fabric " + ItemsName + ".\n\n" +
				"Use this data source to list [" + ItemsName + "](" + ItemDocsURL + ").\n\n" +
				ItemDocsSPNSupport,
		},
		PropertiesSchema: propertiesSchema,
		PropertiesSetter: propertiesSetter,
		ItemListGetter:   itemListGetter,
	}

	return fabricitem.NewDataSourceFabricItemsProperties(config)
}
