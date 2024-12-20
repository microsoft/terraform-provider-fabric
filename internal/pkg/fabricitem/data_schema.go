// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"
	"fmt"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func GetDataSourceFabricItemSchema(ctx context.Context, d DataSourceFabricItem) schema.Schema {
	attributes := getDataSourceFabricItemBaseAttributes(ctx, d.Name, d.IsDisplayNameUnique)

	return schema.Schema{
		MarkdownDescription: d.MarkdownDescription,
		Attributes:          attributes,
	}
}

func GetDataSourceFabricItemDefinitionSchema(ctx context.Context, d DataSourceFabricItemDefinition) schema.Schema {
	attributes := getDataSourceFabricItemBaseAttributes(ctx, d.Name, d.IsDisplayNameUnique)

	for key, value := range getDataSourceFabricItemDefinitionAttributes(ctx, d.Name, d.FormatTypes, d.DefinitionPathKeys) {
		attributes[key] = value
	}

	return schema.Schema{
		MarkdownDescription: d.MarkdownDescription,
		Attributes:          attributes,
	}
}

func GetDataSourceFabricItemPropertiesSchema(ctx context.Context, d DataSourceFabricItem, properties schema.SingleNestedAttribute) schema.Schema {
	attributes := getDataSourceFabricItemBaseAttributes(ctx, d.Name, d.IsDisplayNameUnique)
	attributes["properties"] = properties

	return schema.Schema{
		MarkdownDescription: d.MarkdownDescription,
		Attributes:          attributes,
	}
}

func GetDataSourceFabricItemDefinitionPropertiesSchema(ctx context.Context, d DataSourceFabricItemDefinition, properties schema.SingleNestedAttribute) schema.Schema {
	attributes := getDataSourceFabricItemBaseAttributes(ctx, d.Name, d.IsDisplayNameUnique)
	attributes["properties"] = properties

	for key, value := range getDataSourceFabricItemDefinitionAttributes(ctx, d.Name, d.FormatTypes, d.DefinitionPathKeys) {
		attributes[key] = value
	}

	return schema.Schema{
		MarkdownDescription: d.MarkdownDescription,
		Attributes:          attributes,
	}
}

func GetDataSourceFabricItemDefinitionPropertiesSchema1[Ttfprop, Titemprop any](ctx context.Context, d DataSourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) schema.Schema {
	attributes := getDataSourceFabricItemBaseAttributes(ctx, d.Name, d.IsDisplayNameUnique)
	attributes["properties"] = d.PropertiesSchema

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

	if isDisplayNameUnique {
		attributes["id"] = schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The %s ID.", itemName),
			Optional:            true,
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
		}
		attributes["display_name"] = schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The %s display name.", itemName),
			Optional:            true,
			Computed:            true,
		}
	} else {
		attributes["id"] = schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The %s ID.", itemName),
			Required:            true,
			CustomType:          customtypes.UUIDType{},
		}
		attributes["display_name"] = schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The %s display name.", itemName),
			Computed:            true,
		}
	}

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

	definitionMarkdownDescription := "Definition parts."

	if len(definitionPathKeys) > 0 {
		definitionMarkdownDescription = definitionMarkdownDescription + " Possible path keys: " + utils.ConvertStringSlicesToString(definitionPathKeys, true, false) + "."
	}

	attributes["definition"] = schema.MapNestedAttribute{
		MarkdownDescription: definitionMarkdownDescription,
		Computed:            true,
		CustomType:          supertypes.NewMapNestedObjectTypeOf[DataSourceFabricItemDefinitionPartModel](ctx),
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

	return attributes
}
