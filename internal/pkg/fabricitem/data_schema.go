// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getDataSourceFabricItemSchema(ctx context.Context, d DataSourceFabricItem) schema.Schema {
	attributes := getDataSourceFabricItemBaseAttributes(ctx, d.Name, d.IsDisplayNameUnique)

	return schema.Schema{
		MarkdownDescription: d.MarkdownDescription,
		Attributes:          attributes,
	}
}

func getDataSourceFabricItemDefinitionSchema(ctx context.Context, d DataSourceFabricItemDefinition) schema.Schema {
	attributes := getDataSourceFabricItemBaseAttributes(ctx, d.Name, d.IsDisplayNameUnique)

	for key, value := range getDataSourceFabricItemDefinitionAttributes(ctx, d.Name, d.FormatTypes, d.DefinitionPathKeys) {
		attributes[key] = value
	}

	return schema.Schema{
		MarkdownDescription: d.MarkdownDescription,
		Attributes:          attributes,
	}
}

func getDataSourceFabricItemPropertiesSchema[Ttfprop, Titemprop any](ctx context.Context, d DataSourceFabricItemProperties[Ttfprop, Titemprop]) schema.Schema {
	attributes := getDataSourceFabricItemBaseAttributes(ctx, d.Name, d.IsDisplayNameUnique)
	attributes["properties"] = getDataSourceFabricItemPropertiesNestedAttr[Ttfprop](ctx, d.Name, d.PropertiesAttributes)

	return schema.Schema{
		MarkdownDescription: d.MarkdownDescription,
		Attributes:          attributes,
	}
}

func getDataSourceFabricItemDefinitionPropertiesSchema[Ttfprop, Titemprop any](ctx context.Context, d DataSourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) schema.Schema {
	attributes := getDataSourceFabricItemBaseAttributes(ctx, d.Name, d.IsDisplayNameUnique)
	attributes["properties"] = getDataSourceFabricItemPropertiesNestedAttr[Ttfprop](ctx, d.Name, d.PropertiesAttributes)

	for key, value := range getDataSourceFabricItemDefinitionAttributes(ctx, d.Name, d.FormatTypes, d.DefinitionPathKeys) {
		attributes[key] = value
	}

	return schema.Schema{
		MarkdownDescription: d.MarkdownDescription,
		Attributes:          attributes,
	}
}

// Helper function to get base Fabric Item data-source attributes.
func getDataSourceFabricItemBaseAttributes(ctx context.Context, itemName string, isDisplayNameUnique bool) map[string]schema.Attribute { //revive:disable-line:flag-parameter
	attributes := map[string]schema.Attribute{
		"workspace_id": schema.StringAttribute{
			MarkdownDescription: "The Workspace ID.",
			Required:            true,
			CustomType:          customtypes.UUIDType{},
		},
		"description": schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The %s description.", itemName),
			Computed:            true,
		},
		"timeouts": timeouts.Attributes(ctx),
	}

	// id attribute
	attrID := schema.StringAttribute{}
	attrID.MarkdownDescription = fmt.Sprintf("The %s ID.", itemName)
	attrID.CustomType = customtypes.UUIDType{}

	if isDisplayNameUnique {
		attrID.Optional = true
		attrID.Computed = true
	} else {
		attrID.Required = true
	}

	attributes["id"] = attrID

	// display_name attribute
	attrDisplayName := schema.StringAttribute{}
	attrDisplayName.MarkdownDescription = fmt.Sprintf("The %s display name.", itemName)
	attrDisplayName.Computed = true

	if isDisplayNameUnique {
		attrDisplayName.Optional = true
	}

	attributes["display_name"] = attrDisplayName

	return attributes
}

// Helper function to get Fabric Item data-source definition attributes.
func getDataSourceFabricItemDefinitionAttributes(ctx context.Context, name string, formatTypes, definitionPathKeys []string) map[string]schema.Attribute {
	attributes := make(map[string]schema.Attribute)

	if len(formatTypes) > 0 {
		attributes["format"] = schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The %s format. Possible values: %s.", name, utils.ConvertStringSlicesToString(formatTypes, true, false)),
			Computed:            true,
		}
	} else {
		attributes["format"] = schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The %s format. Possible values: `%s`", name, DefinitionFormatNotApplicable),
			Computed:            true,
		}
	}

	attributes["output_definition"] = schema.BoolAttribute{
		MarkdownDescription: "Output definition parts as gzip base64 content? Default: `false`\n\n" +
			"!> Your terraform state file may grow a lot if you output definition content. Only use it when you must use data from the definition.",
		Optional: true,
		Computed: true,
	}

	attrDefinition := schema.MapNestedAttribute{
		Computed:   true,
		CustomType: supertypes.NewMapNestedObjectTypeOf[dataSourceFabricItemDefinitionPartModel](ctx),
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"content": schema.StringAttribute{
					MarkdownDescription: "Gzip base64 content of definition part.\n" +
						"Use [`provider::fabric::content_decode`](../functions/content_decode.md) function to decode content.",
					Computed: true,
				},
			},
		},
	}

	if len(definitionPathKeys) > 0 {
		attrDefinition.MarkdownDescription = "Definition parts. Possible path keys: " + utils.ConvertStringSlicesToString(definitionPathKeys, true, false) + "."
	} else {
		attrDefinition.MarkdownDescription = "Definition parts."
	}

	attributes["definition"] = attrDefinition

	return attributes
}

func getDataSourceFabricItemPropertiesNestedAttr[Ttfprop any](ctx context.Context, name string, attributes map[string]schema.Attribute) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "The " + name + " properties.",
		Computed:            true,
		CustomType:          supertypes.NewSingleNestedObjectTypeOf[Ttfprop](ctx),
		Attributes:          attributes,
	}
}
