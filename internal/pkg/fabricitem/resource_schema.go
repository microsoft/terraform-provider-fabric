// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"
	"fmt"
	"regexp"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/planmodifiers"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func GetResourceFabricItemSchema(ctx context.Context, r ResourceFabricItem) schema.Schema {
	attributes := getResourceFabricItemBaseAttributes(ctx, r.Name, r.DisplayNameMaxLength, r.DescriptionMaxLength, r.NameRenameAllowed)

	return schema.Schema{
		MarkdownDescription: r.MarkdownDescription,
		Attributes:          attributes,
	}
}

func GetResourceFabricItemDefinitionSchema(ctx context.Context, r ResourceFabricItemDefinition) schema.Schema {
	attributes := getResourceFabricItemBaseAttributes(ctx, r.Name, r.DisplayNameMaxLength, r.DescriptionMaxLength, r.NameRenameAllowed)

	for key, value := range getResourceFabricItemDefinitionAttributes(ctx, r) {
		attributes[key] = value
	}

	return schema.Schema{
		MarkdownDescription: r.MarkdownDescription,
		Attributes:          attributes,
	}
}

func GetResourceFabricItemPropertiesSchema(ctx context.Context, itemName, markdownDescription string, displayNameMaxLength, descriptionMaxLength int, nameRenameAllowed bool, properties schema.SingleNestedAttribute) schema.Schema {
	attributes := getResourceFabricItemBaseAttributes(ctx, itemName, displayNameMaxLength, descriptionMaxLength, nameRenameAllowed)
	attributes["properties"] = properties

	return schema.Schema{
		MarkdownDescription: markdownDescription,
		Attributes:          attributes,
	}
}

func GetResourceFabricItemPropertiesDefinitionSchema(ctx context.Context, r ResourceFabricItemDefinition, properties schema.SingleNestedAttribute) schema.Schema {
	attributes := getResourceFabricItemBaseAttributes(ctx, r.Name, r.DisplayNameMaxLength, r.DescriptionMaxLength, r.NameRenameAllowed)
	attributes["properties"] = properties

	for key, value := range getResourceFabricItemDefinitionAttributes(ctx, r) {
		attributes[key] = value
	}

	return schema.Schema{
		MarkdownDescription: r.MarkdownDescription,
		Attributes:          attributes,
	}
}

func GetResourceFabricItemPropertiesCreationSchema(ctx context.Context, itemName, markdownDescription string, displayNameMaxLength, descriptionMaxLength int, nameRenameAllowed bool, properties, configuration schema.SingleNestedAttribute) schema.Schema {
	attributes := getResourceFabricItemBaseAttributes(ctx, itemName, displayNameMaxLength, descriptionMaxLength, nameRenameAllowed)
	attributes["properties"] = properties
	attributes["configuration"] = configuration

	return schema.Schema{
		MarkdownDescription: markdownDescription,
		Attributes:          attributes,
	}
}

// Helper function to get base Fabric Item resource attributes.
func getResourceFabricItemBaseAttributes(ctx context.Context, name string, displayNameMaxLength, descriptionMaxLength int, nameRenameAllowed bool) map[string]schema.Attribute { //revive:disable-line:flag-parameter
	displayNamePlanModifiers := []planmodifier.String{}

	if !nameRenameAllowed {
		displayNamePlanModifiers = []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		}
	}

	attributes := map[string]schema.Attribute{
		"workspace_id": schema.StringAttribute{
			MarkdownDescription: "The Workspace ID.",
			Required:            true,
			CustomType:          customtypes.UUIDType{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"id": schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The %s ID.", name),
			Computed:            true,
			CustomType:          customtypes.UUIDType{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The %s display name.", name),
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(displayNameMaxLength),
			},
			PlanModifiers: displayNamePlanModifiers,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The %s description.", name),
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			Validators: []validator.String{
				stringvalidator.LengthAtMost(descriptionMaxLength),
			},
		},
		"timeouts": timeouts.AttributesAll(ctx),
	}

	return attributes
}

// Helper function to get Fabric Item definition attributes.
func getResourceFabricItemDefinitionAttributes(ctx context.Context, r ResourceFabricItemDefinition) map[string]schema.Attribute {
	attributes := make(map[string]schema.Attribute)

	attributes["definition_update_enabled"] = schema.BoolAttribute{
		MarkdownDescription: "Update definition on change of source content. Default: `true`.",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
	}

	if len(r.FormatTypes) > 0 {
		attributes["format"] = schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The %s format. Possible values: %s.", r.Name, utils.ConvertStringSlicesToString(r.FormatTypes, true, false)),
			Computed:            true,
			Default:             stringdefault.StaticString(r.FormatTypeDefault),
		}
	} else {
		attributes["format"] = schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("The %s format. Possible values: `%s`", r.Name, DefinitionFormatNotApplicable),
			Computed:            true,
			Default:             stringdefault.StaticString(DefinitionFormatNotApplicable),
		}
	}

	if r.DefinitionRequired {
		attributes["definition"] = schema.MapNestedAttribute{
			MarkdownDescription: fmt.Sprintf("Definition parts. Accepted path keys: %s. Read more about [%s definition part paths](%s).", utils.ConvertStringSlicesToString(r.DefinitionPathKeys, true, false), r.Name, r.DefinitionPathDocsURL),
			Required:            true,
			CustomType:          supertypes.NewMapNestedObjectTypeOf[ResourceFabricItemDefinitionPartModel](ctx),
			Validators:          r.DefinitionPathKeysValidator,
			NestedObject:        getResourceFabricItemDefinitionPartSchema(ctx),
		}
	} else {
		attributes["definition"] = schema.MapNestedAttribute{
			MarkdownDescription: fmt.Sprintf("Definition parts. Accepted path keys: %s. Read more about [%s definition part paths](%s).", utils.ConvertStringSlicesToString(r.DefinitionPathKeys, true, false), r.Name, r.DefinitionPathDocsURL),
			Optional:            true,
			CustomType:          supertypes.NewMapNestedObjectTypeOf[ResourceFabricItemDefinitionPartModel](ctx),
			Validators:          r.DefinitionPathKeysValidator,
			NestedObject:        getResourceFabricItemDefinitionPartSchema(ctx),
		}
	}

	return attributes
}

// Helper function to get Fabric Item data-source definition part attributes.
func getResourceFabricItemDefinitionPartSchema(ctx context.Context) schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"source": schema.StringAttribute{
				MarkdownDescription: "Path to the file with source of the definition part.\n\n" +
					"The source content may include placeholders for token substitution. Use the dot with the token name `{{ .TokenName }}`.",
				Required: true,
			},
			"tokens": schema.MapAttribute{
				MarkdownDescription: "A map of key/value pairs of tokens substitutes in the source.",
				Optional:            true,
				CustomType:          supertypes.NewMapTypeOf[string](ctx),
				Validators: []validator.Map{
					mapvalidator.KeysAre(stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9]+([_]?[a-zA-Z0-9]+)*$`),
						"Token key:\n"+
							"- cannot contains special characters\n"+
							"- cannot contains any white spaces\n"+
							"- underscore '_' is allowed but not at the start or end of the token key",
					)),
				},
			},
			"source_content_sha256": schema.StringAttribute{
				MarkdownDescription: "SHA256 of source's content of definition part.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.DefinitionContentSha256(path.MatchRelative().AtParent().AtName("source"), path.MatchRelative().AtParent().AtName("tokens")),
				},
			},
		},
	}
}
