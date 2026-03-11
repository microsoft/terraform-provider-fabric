// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package deploymentpipeline

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func itemSchema(isList bool) superschema.Schema { //revive:disable-line:flag-parameter
	markdownDescriptionR := fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false)
	markdownDescriptionD := fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, isList)

	var dsTimeout *superschema.DatasourceTimeoutAttribute

	if !isList {
		dsTimeout = &superschema.DatasourceTimeoutAttribute{
			Read: true,
		}
	}

	var stagesAttribute superschema.SuperListNestedAttributeOf[baseDeploymentPipelineStageModel]
	if !isList {
		stagesAttribute = superschema.SuperListNestedAttributeOf[baseDeploymentPipelineStageModel]{
			Common: &schemaR.ListNestedAttribute{
				MarkdownDescription: "The collection of " + ItemTypeInfo.Name + " stages.",
			},
			Resource: &schemaR.ListNestedAttribute{
				Required: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				Validators: []validator.List{
					listvalidator.SizeBetween(2, 10),
				},
			},
			DataSource: &schemaD.ListNestedAttribute{
				Computed: true,
			},
			Attributes: superschema.Attributes{
				"id": superschema.SuperStringAttribute{
					Common: &schemaR.StringAttribute{
						MarkdownDescription: "The ID of the stage.",
						CustomType:          customtypes.UUIDType{},
					},
					Resource: &schemaR.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					DataSource: &schemaD.StringAttribute{
						Optional: true,
						Computed: true,
					},
				},
				"display_name": superschema.SuperStringAttribute{
					Common: &schemaR.StringAttribute{
						MarkdownDescription: "The display name of the stage.",
					},
					DataSource: &schemaD.StringAttribute{
						Computed: true,
					},
					Resource: &schemaR.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.LengthAtMost(256),
						},
					},
				},
				"description": superschema.SuperStringAttribute{
					Common: &schemaR.StringAttribute{
						MarkdownDescription: "The description of the stage.",
					},
					DataSource: &schemaD.StringAttribute{
						Computed: true,
					},
					Resource: &schemaR.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.LengthAtMost(1024),
						},
					},
				},
				"is_public": superschema.SuperBoolAttribute{
					Common: &schemaR.BoolAttribute{
						MarkdownDescription: "Whether the stage is public.",
					},
					DataSource: &schemaD.BoolAttribute{
						Computed: true,
					},
					Resource: &schemaR.BoolAttribute{
						Required: true,
					},
				},
				"workspace_id": superschema.SuperStringAttribute{
					Common: &schemaR.StringAttribute{
						MarkdownDescription: "The assigned workspace ID.",
						CustomType:          customtypes.UUIDType{},
					},
					Resource: &schemaR.StringAttribute{
						Optional: true,
						Computed: false,
					},
					DataSource: &schemaD.StringAttribute{
						Computed: true,
					},
				},
			},
		}
	}

	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: markdownDescriptionR,
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: markdownDescriptionD,
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Optional: !isList,
					Computed: true,
				},
			},
			"display_name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " display name.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.LengthAtMost(246),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Optional: !isList,
					Computed: true,
				},
			},
			"description": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " description.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Computed: true,
					Default:  stringdefault.StaticString(""),
					Validators: []validator.String{
						stringvalidator.LengthAtMost(256),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"stages": stagesAttribute,
			"timeouts": superschema.TimeoutAttribute{
				Resource: &superschema.ResourceTimeoutAttribute{
					Create: true,
					Read:   true,
					Update: true,
					Delete: true,
				},
				DataSource: dsTimeout,
			},
		},
	}
}
