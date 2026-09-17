// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package semanticmodelcb

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	fabsemanticmodel "github.com/microsoft/fabric-sdk-go/fabric/semanticmodel"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func itemSchema() superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false),
		},
		Attributes: map[string]superschema.Attribute{
			"workspace_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Workspace ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"semantic_model_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + "Semantic Model" + " ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"connectivity_type": superschema.StringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The connectivity type of the connection. Use `None` to unbind a data source reference.",
					Required:            true,
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabsemanticmodel.PossibleConnectivityTypeValues(), true)...),
					},
				},
			},
			"connection_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The object ID of the connection to bind. Not required when `connectivity_type` is `None`.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
			},
			"connection_details": superschema.SuperSingleNestedAttributeOf[connectionDetailsModel]{
				Resource: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The connection details identifying the data source reference of the semantic model to bind.",
					Required:            true,
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.RequiresReplace(),
					},
				},
				Attributes: map[string]superschema.Attribute{
					"path": superschema.StringAttribute{
						Resource: &schemaR.StringAttribute{
							MarkdownDescription: "The path of the connection (matches the data source reference path in the semantic model definition).",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
					"type": superschema.StringAttribute{
						Resource: &schemaR.StringAttribute{
							MarkdownDescription: "The type of the connection (e.g. `SQL`, `Web`).",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
			},
			"timeouts": superschema.TimeoutAttribute{
				Resource: &superschema.ResourceTimeoutAttribute{
					Create: true,
					Read:   true,
					Update: true,
					Delete: true,
				},
			},
		},
	}
}
