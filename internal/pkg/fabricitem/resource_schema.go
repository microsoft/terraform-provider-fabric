// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"
	"fmt"
	"maps"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	supermapvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/mapvalidator"
	supersetvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/setvalidator"
	superstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/planmodifiers"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/params"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func getResourceFabricItemSchema(ctx context.Context, r ResourceFabricItem) schema.Schema {
	attributes := getResourceFabricItemBaseAttributes(ctx, r.TypeInfo.Name, r.DisplayNameMaxLength, r.DescriptionMaxLength, r.NameRenameAllowed)

	return schema.Schema{
		MarkdownDescription: NewResourceMarkdownDescription(r.TypeInfo, false),
		Attributes:          attributes,
	}
}

func getResourceFabricItemDefinitionSchema(ctx context.Context, r ResourceFabricItemDefinition) schema.Schema {
	attributes := getResourceFabricItemBaseAttributes(ctx, r.TypeInfo.Name, r.DisplayNameMaxLength, r.DescriptionMaxLength, r.NameRenameAllowed)

	maps.Copy(attributes, getResourceFabricItemDefinitionAttributes(ctx, r.TypeInfo.Name, r.DefinitionPathDocsURL, r.DefinitionFormats, r.DefinitionPathKeysValidator, r.DefinitionRequired, false))

	return schema.Schema{
		MarkdownDescription: NewResourceMarkdownDescription(r.TypeInfo, false),
		Attributes:          attributes,
	}
}

func getResourceFabricItemPropertiesSchema[Ttfprop, Titemprop any](ctx context.Context, r ResourceFabricItemProperties[Ttfprop, Titemprop]) schema.Schema {
	attributes := getResourceFabricItemBaseAttributes(ctx, r.TypeInfo.Name, r.DisplayNameMaxLength, r.DescriptionMaxLength, r.NameRenameAllowed)
	attributes["properties"] = getResourceFabricItemPropertiesNestedAttr[Ttfprop](ctx, r.TypeInfo.Name, r.PropertiesAttributes)

	return schema.Schema{
		MarkdownDescription: NewResourceMarkdownDescription(r.TypeInfo, false),
		Attributes:          attributes,
	}
}

func getResourceFabricItemDefinitionPropertiesSchema[Ttfprop, Titemprop any](ctx context.Context, r ResourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) schema.Schema {
	attributes := getResourceFabricItemBaseAttributes(ctx, r.TypeInfo.Name, r.DisplayNameMaxLength, r.DescriptionMaxLength, r.NameRenameAllowed)
	attributes["properties"] = getResourceFabricItemPropertiesNestedAttr[Ttfprop](ctx, r.TypeInfo.Name, r.PropertiesAttributes)

	maps.Copy(attributes, getResourceFabricItemDefinitionAttributes(ctx, r.TypeInfo.Name, r.DefinitionPathDocsURL, r.DefinitionFormats, r.DefinitionPathKeysValidator, r.DefinitionRequired, false))

	return schema.Schema{
		MarkdownDescription: NewResourceMarkdownDescription(r.TypeInfo, false),
		Attributes:          attributes,
	}
}

func getResourceFabricItemConfigPropertiesSchema[Ttfprop, Titemprop, Ttfconfig, Titemconfig any](
	ctx context.Context,
	r ResourceFabricItemConfigProperties[Ttfprop, Titemprop, Ttfconfig, Titemconfig],
) schema.Schema {
	attributes := getResourceFabricItemBaseAttributes(ctx, r.TypeInfo.Name, r.DisplayNameMaxLength, r.DescriptionMaxLength, r.NameRenameAllowed)
	attributes["configuration"] = getResourceFabricItemConfigNestedAttr[Ttfconfig](ctx, r.TypeInfo.Name, r.ConfigRequired, r.ConfigAttributes)
	attributes["properties"] = getResourceFabricItemPropertiesNestedAttr[Ttfprop](ctx, r.TypeInfo.Name, r.PropertiesAttributes)

	return schema.Schema{
		MarkdownDescription: NewResourceMarkdownDescription(r.TypeInfo, false),
		Attributes:          attributes,
	}
}

func getResourceFabricItemConfigDefinitionPropertiesSchema[Ttfprop, Titemprop, Ttfconfig, Titemconfig any](
	ctx context.Context,
	r ResourceFabricItemConfigDefinitionProperties[Ttfprop, Titemprop, Ttfconfig, Titemconfig],
) schema.Schema {
	attributes := getResourceFabricItemBaseAttributes(ctx, r.TypeInfo.Name, r.DisplayNameMaxLength, r.DescriptionMaxLength, r.NameRenameAllowed)
	attrConfiguration := getResourceFabricItemConfigNestedAttr[Ttfconfig](ctx, r.TypeInfo.Name, r.ConfigRequired, r.ConfigAttributes)
	attrConfiguration.Validators = []validator.Object{
		objectvalidator.ConflictsWith(
			path.MatchRoot("definition"),
			path.MatchRoot("definition_update_enabled"),
			path.MatchRoot("format"),
		),
	}
	attributes["configuration"] = attrConfiguration
	attributes["properties"] = getResourceFabricItemPropertiesNestedAttr[Ttfprop](ctx, r.TypeInfo.Name, r.PropertiesAttributes)

	maps.Copy(attributes, getResourceFabricItemDefinitionAttributes(ctx, r.TypeInfo.Name, r.DefinitionPathDocsURL, r.DefinitionFormats, r.DefinitionPathKeysValidator, r.DefinitionRequired, true))

	return schema.Schema{
		MarkdownDescription: NewResourceMarkdownDescription(r.TypeInfo, false),
		Attributes:          attributes,
	}
}

func getResourceFabricItemConfigNestedAttr[Ttfconfig any](
	ctx context.Context,
	name string,
	isRequired bool,
	attributes map[string]schema.Attribute,
) schema.SingleNestedAttribute { //revive:disable-line:flag-parameter
	result := schema.SingleNestedAttribute{
		MarkdownDescription: "The " + name + " creation configuration.\n\n" +
			"Any changes to this configuration will result in recreation of the " + name + ".",
		CustomType: supertypes.NewSingleNestedObjectTypeOf[Ttfconfig](ctx),
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
		Attributes: attributes,
	}

	if isRequired {
		result.Required = true
	} else {
		result.Optional = true
	}

	return result
}

func getResourceFabricItemPropertiesNestedAttr[Ttfprop any](ctx context.Context, name string, attributes map[string]schema.Attribute) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "The " + name + " properties.",
		Computed:            true,
		CustomType:          supertypes.NewSingleNestedObjectTypeOf[Ttfprop](ctx),
		Attributes:          attributes,
	}
}

// Helper function to get base Fabric Item resource attributes.
func getResourceFabricItemBaseAttributes(
	ctx context.Context,
	name string,
	displayNameMaxLength, descriptionMaxLength int,
	nameRenameAllowed bool,
) map[string]schema.Attribute { //revive:disable-line:flag-parameter
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
		"folder_id": schema.StringAttribute{
			MarkdownDescription: "The Folder ID.",
			Optional:            true,
			CustomType:          customtypes.UUIDType{},
		},
		"timeouts": timeouts.AttributesAll(ctx),
	}

	return attributes
}

// Helper function to get Fabric Item definition attributes.
func getResourceFabricItemDefinitionAttributes(
	ctx context.Context,
	name, definitionPathDocsURL string,
	definitionFormats []DefinitionFormat,
	definitionPathKeysValidator []validator.Map,
	definitionRequired, alongConfiguration bool,
) map[string]schema.Attribute { //revive:disable-line:flag-parameter,argument-limit
	attributes := make(map[string]schema.Attribute)

	attrDefinitionUpdateEnabled := schema.BoolAttribute{}

	attrDefinitionUpdateEnabled.MarkdownDescription = "Update definition on change of source content. Default: `true`."
	attrDefinitionUpdateEnabled.Optional = true
	attrDefinitionUpdateEnabled.Computed = true
	attrDefinitionUpdateEnabled.Default = booldefault.StaticBool(true)

	if alongConfiguration {
		attrDefinitionUpdateEnabled.Validators = []validator.Bool{
			boolvalidator.ConflictsWith(path.MatchRoot("configuration")),
		}
	}

	attributes["definition_update_enabled"] = attrDefinitionUpdateEnabled

	formatTypes := getDefinitionFormats(definitionFormats)
	definitionFormatsDocs := getDefinitionFormatsPathsDocs(definitionFormats)

	attrFormat := schema.StringAttribute{}

	attrFormat.MarkdownDescription = fmt.Sprintf("The %s format. Possible values: %s", name, utils.ConvertStringSlicesToString(formatTypes, true, true))
	attrFormat.Validators = []validator.String{
		stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(formatTypes, true)...),
		superstringvalidator.RequireIfAttributeIsSet(path.MatchRoot("definition")),
	}

	if definitionRequired {
		attrFormat.Required = true
	} else {
		attrFormat.Optional = true
	}

	if alongConfiguration {
		attrFormat.Validators = append(attrFormat.Validators, stringvalidator.ConflictsWith(path.MatchRoot("configuration")))
	}

	attributes["format"] = attrFormat

	attrDefinition := schema.MapNestedAttribute{}
	attrDefinition.MarkdownDescription = fmt.Sprintf("Definition parts. Read more about [%s definition part paths](%s). Accepted path keys: %s", name, definitionPathDocsURL, definitionFormatsDocs)
	attrDefinition.CustomType = supertypes.NewMapNestedObjectTypeOf[resourceFabricItemDefinitionPartModel](ctx)
	attrDefinition.Validators = definitionPathKeysValidator
	attrDefinition.NestedObject = getResourceFabricItemDefinitionPartSchema(ctx)

	if definitionRequired {
		attrDefinition.Required = true
	} else {
		attrDefinition.Optional = true
	}

	if alongConfiguration {
		definitionPathKeysValidator = append(definitionPathKeysValidator, mapvalidator.ConflictsWith(path.MatchRoot("configuration")))
		attrDefinition.Validators = definitionPathKeysValidator
	}

	attributes["definition"] = attrDefinition

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
			"processing_mode": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf(
					"Processing mode of the tokens/parameters. Possible values: %s. Default `%s`",
					utils.ConvertStringSlicesToString(transforms.PossibleProcessingModeValues(), true, true),
					transforms.ProcessingModeGoTemplate,
				),
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(transforms.ProcessingModeGoTemplate),
				Validators: []validator.String{
					stringvalidator.OneOf(transforms.PossibleProcessingModeValues()...),
				},
			},
			"parameters": schema.SetNestedAttribute{
				MarkdownDescription: "The set of parameters to be passed and processed in the source content.",
				Optional:            true,
				CustomType:          supertypes.NewSetNestedObjectTypeOf[params.ParametersModel](ctx),
				Validators: []validator.Set{
					setvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("processing_mode")),
					supersetvalidator.RequireIfAttributeIsOneOf(
						path.MatchRelative().AtParent().AtName("processing_mode"),
						[]attr.Value{types.StringValue(transforms.ProcessingModeParameters)},
					),
					supersetvalidator.NullIfAttributeIsOneOf(
						path.MatchRelative().AtParent().AtName("processing_mode"),
						[]attr.Value{types.StringValue(transforms.ProcessingModeGoTemplate), types.StringValue(transforms.ProcessingModeNone)},
					),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: fmt.Sprintf(
								"Processing type of the parameters. Possible values: %s.",
								utils.ConvertStringSlicesToString(transforms.PossibleParameterTypeValues(), true, true),
							),
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf(transforms.PossibleParameterTypeValues()...),
							},
						},
						"find": schema.StringAttribute{
							MarkdownDescription: "The find value of the parameter.",
							Required:            true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "The value of the parameter.",
							Required:            true,
						},
					},
				},
			},
			"tokens": schema.MapAttribute{
				MarkdownDescription: "A map of key/value pairs of tokens substitutes in the source.",
				Optional:            true,
				CustomType:          supertypes.MapTypeOf[types.String]{MapType: types.MapType{ElemType: types.StringType}},
				ElementType:         types.StringType,
				Validators: []validator.Map{
					mapvalidator.KeysAre(stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9]+([_]?[a-zA-Z0-9]+)*$`),
						"Token key:\n"+
							"- cannot contains special characters\n"+
							"- cannot contains any white spaces\n"+
							"- underscore '_' is allowed but not at the start or end of the token key",
					)),
					supermapvalidator.RequireIfAttributeIsOneOf(
						path.MatchRelative().AtParent().AtName("processing_mode"),
						[]attr.Value{types.StringValue(transforms.ProcessingModeGoTemplate)},
					),
					supermapvalidator.NullIfAttributeIsOneOf(
						path.MatchRelative().AtParent().AtName("processing_mode"),
						[]attr.Value{types.StringValue(transforms.ProcessingModeParameters), types.StringValue(transforms.ProcessingModeNone)},
					),
				},
			},
			"tokens_delimiter": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("The delimiter for the tokens in the source content. Possible values: %s. Default: `%s`",
					utils.ConvertStringSlicesToString(transforms.PossibleTokensDelimiterValues(), true, true),
					transforms.TokensDelimiterCurlyBraces,
				),
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(transforms.TokensDelimiterCurlyBraces),
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("tokens")),
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("parameters")),
					stringvalidator.OneOf(transforms.PossibleTokensDelimiterValues()...),
				},
			},
			"source_content_sha256": schema.StringAttribute{
				MarkdownDescription: "SHA256 of source's content of definition part.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.DefinitionContentSha256(
						path.MatchRelative().AtParent().AtName("source"),
						path.MatchRelative().AtParent().AtName("processing_mode"),
						path.MatchRelative().AtParent().AtName("tokens"),
						path.MatchRelative().AtParent().AtName("parameters"),
						path.MatchRelative().AtParent().AtName("tokens_delimiter"),
					),
				},
			},
		},
	}
}
