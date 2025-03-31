// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domainwa

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func itemSchema() superschema.Schema { //revive:disable-line:flag-parameter
	markdownDescriptionR := fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, true)
	markdownDescriptionD := fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, true)

	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: markdownDescriptionR,
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: markdownDescriptionD,
		},

		Attributes: map[string]superschema.Attribute{
			// "id": superschema.SuperStringAttribute{
			// 	Common: &schemaR.StringAttribute{
			// 		MarkdownDescription: "The " + ItemTypeInfo.Name + " ID.",
			// 		CustomType:          customtypes.UUIDType{},
			// 	},
			// 	Resource: &schemaR.StringAttribute{
			// 		Computed: true,
			// 		PlanModifiers: []planmodifier.String{
			// 			stringplanmodifier.UseStateForUnknown(),
			// 		},
			// 	},
			// 	DataSource: &schemaD.StringAttribute{
			// 		Computed: true,
			// 	},
			// },
			"domain_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Domain ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: true,
				},
			},
			"workspace_ids": superschema.SuperSetAttribute{
				Common: &schemaR.SetAttribute{
					MarkdownDescription: "The set of Workspace IDs.",
					CustomType: customtypes.SetTypeOf[customtypes.UUID]{
						SetType: basetypes.SetType{
							ElemType: customtypes.UUIDType{},
						},
					},
					ElementType: customtypes.UUIDType{},
				},
				Resource: &schemaR.SetAttribute{
					Required: true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
				},
				DataSource: &schemaD.SetAttribute{
					Computed: true,
				},
			},
			"timeouts": superschema.TimeoutAttribute{
				Resource: &superschema.ResourceTimeoutAttribute{
					Create: true,
					Read:   true,
					Update: true,
					Delete: true,
				},
				DataSource: &superschema.DatasourceTimeoutAttribute{
					Read: true,
				},
			},
		},
	}
}
