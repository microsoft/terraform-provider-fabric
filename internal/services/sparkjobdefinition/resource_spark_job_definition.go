// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sparkjobdefinition

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabsparkjobdefinition "github.com/microsoft/fabric-sdk-go/fabric/sparkjobdefinition"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func NewResourceSparkJobDefinition(ctx context.Context) resource.Resource {
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

	propertiesSetter := func(ctx context.Context, from *fabsparkjobdefinition.Properties, to *fabricitem.ResourceFabricItemDefinitionPropertiesModel[sparkJobDefinitionPropertiesModel, fabsparkjobdefinition.Properties]) diag.Diagnostics {
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

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemDefinitionPropertiesModel[sparkJobDefinitionPropertiesModel, fabsparkjobdefinition.Properties], fabricItem *fabricitem.FabricItemProperties[fabsparkjobdefinition.Properties]) error {
		client := fabsparkjobdefinition.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetSparkJobDefinition(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.SparkJobDefinition)

		return nil
	}

	config := fabricitem.ResourceFabricItemDefinitionProperties[sparkJobDefinitionPropertiesModel, fabsparkjobdefinition.Properties]{
		ResourceFabricItemDefinition: fabricitem.ResourceFabricItemDefinition{
			Type:              ItemType,
			Name:              ItemName,
			NameRenameAllowed: true,
			TFName:            ItemTFName,
			MarkdownDescription: "Manage a Fabric " + ItemName + ".\n\n" +
				"Use this resource to manage a [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
				ItemDocsSPNSupport,
			DisplayNameMaxLength:  123,
			DescriptionMaxLength:  256,
			FormatTypeDefault:     ItemFormatTypeDefault,
			FormatTypes:           ItemFormatTypes,
			DefinitionPathDocsURL: ItemDefinitionPathDocsURL,
			DefinitionPathKeys:    ItemDefinitionPaths,
			DefinitionPathKeysValidator: []validator.Map{
				mapvalidator.SizeAtMost(1),
				mapvalidator.KeysAre(stringvalidator.OneOf(ItemDefinitionPaths...)),
			},
			DefinitionRequired: false,
			DefinitionEmpty:    ItemDefinitionEmpty,
		},
		PropertiesSchema: propertiesSchema,
		PropertiesSetter: propertiesSetter,
		ItemGetter:       itemGetter,
	}

	return fabricitem.NewResourceFabricItemDefinitionProperties(config)
}
