// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domainra

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func itemSchema(isList bool) superschema.Schema { //revive:disable-line:flag-parameter
	markdownDescriptionR := fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, true)
	markdownDescriptionD := fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, isList)

	var dsTimeout *superschema.DatasourceTimeoutAttribute

	if !isList {
		dsTimeout = &superschema.DatasourceTimeoutAttribute{
			Read: true,
		}
	}

	possiblePrincipalTypeValues := utils.RemoveSlicesByValues(
		fabadmin.PossiblePrincipalTypeValues(),
		[]fabadmin.PrincipalType{fabadmin.PrincipalTypeServicePrincipal, fabadmin.PrincipalTypeServicePrincipalProfile},
	)

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
			// 		Required: !isList,
			// 		Computed: isList,
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
			"role": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Role of the principals.",
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(PossibleDomainRoleValues(), true)...),
					},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"principals": superschema.SuperSetNestedAttributeOf[principalModel]{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "The set of Principals.",
				},
				Resource: &schemaR.SetNestedAttribute{
					Required: true,
					PlanModifiers: []planmodifier.Set{
						setplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.SetNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"id": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The Principal ID.",
							CustomType:          customtypes.UUIDType{},
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"type": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The Principal type.",
							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(possiblePrincipalTypeValues, true)...),
							},
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
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
				DataSource: dsTimeout,
			},
		},
	}
}
