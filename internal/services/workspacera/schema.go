// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspacera

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

func itemSchema(isList bool) superschema.Schema { //revive:disable-line:flag-parameter
	var dsTimeout *superschema.DatasourceTimeoutAttribute

	if !isList {
		dsTimeout = &superschema.DatasourceTimeoutAttribute{
			Read: true,
		}
	}

	possiblePrincipalTypeValues := utils.RemoveSlicesByValues(
		fabcore.PossiblePrincipalTypeValues(),
		[]fabcore.PrincipalType{fabcore.PrincipalTypeEntireTenant},
	)

	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false),
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, isList),
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
					Required: !isList,
					Computed: isList,
				},
			},
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
				DataSource: &schemaD.StringAttribute{
					Required: !isList,
					Computed: isList,
				},
			},
			"role": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The workspace role of the principal.",
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleWorkspaceRoleValues(), true)...),
					},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"principal": superschema.SuperSingleNestedAttributeOf[principalModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The principal.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Required: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"id": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The principal ID.",
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
					"type": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The type of the principal.",
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
			// "principal_id": superschema.SuperStringAttribute{
			// 	Common: &schemaR.StringAttribute{
			// 		MarkdownDescription: "The principal ID.",
			// 		CustomType:          customtypes.UUIDType{},
			// 	},
			// 	Resource: &schemaR.StringAttribute{
			// 		Required: true,
			// 		PlanModifiers: []planmodifier.String{
			// 			stringplanmodifier.RequiresReplace(),
			// 		},
			// 	},
			// 	DataSource: &schemaD.StringAttribute{
			// 		Computed: true,
			// 	},
			// },
			// "principal_type": superschema.SuperStringAttribute{
			// 	Common: &schemaR.StringAttribute{
			// 		MarkdownDescription: "The type of the principal.",
			// 		CustomType:          customtypes.UUIDType{},
			// 		Validators: []validator.String{
			// 			stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossiblePrincipalTypeValues(), true)...),
			// 		},
			// 	},
			// 	Resource: &schemaR.StringAttribute{
			// 		Required: true,
			// 		PlanModifiers: []planmodifier.String{
			// 			stringplanmodifier.RequiresReplace(),
			// 		},
			// 	},
			// 	DataSource: &schemaD.StringAttribute{
			// 		Computed: true,
			// 	},
			// },
			// "principal_display_name": superschema.StringAttribute{
			// 	Common: &schemaR.StringAttribute{
			// 		MarkdownDescription: "The principal's display name.",
			// 	},
			// 	Resource: &schemaR.StringAttribute{
			// 		Computed: true,
			// 	},
			// 	DataSource: &schemaD.StringAttribute{
			// 		Computed: true,
			// 	},
			// },
			// "principal_details": superschema.SuperSingleNestedAttributeOf[principalDetailsModel]{
			// 	Common: &schemaR.SingleNestedAttribute{
			// 		MarkdownDescription: "The principal details.",
			// 	},
			// 	Resource: &schemaR.SingleNestedAttribute{
			// 		Computed: true,
			// 	},
			// 	DataSource: &schemaD.SingleNestedAttribute{
			// 		Computed: true,
			// 	},
			// 	Attributes: map[string]superschema.Attribute{
			// 		"user_principal_name": superschema.StringAttribute{
			// 			Common: &schemaR.StringAttribute{
			// 				MarkdownDescription: "The principal ID.",
			// 			},
			// 			Resource: &schemaR.StringAttribute{
			// 				Computed: true,
			// 			},
			// 			DataSource: &schemaD.StringAttribute{
			// 				Computed: true,
			// 			},
			// 		},
			// 		"group_type": superschema.StringAttribute{
			// 			Common: &schemaR.StringAttribute{
			// 				MarkdownDescription: "The type of the group.",
			// 				Validators: []validator.String{
			// 					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleGroupTypeValues(), true)...),
			// 				},
			// 			},
			// 			Resource: &schemaR.StringAttribute{
			// 				Computed: true,
			// 			},
			// 			DataSource: &schemaD.StringAttribute{
			// 				Computed: true,
			// 			},
			// 		},
			// 		"app_id": superschema.StringAttribute{
			// 			Common: &schemaR.StringAttribute{
			// 				MarkdownDescription: "The service principal's Microsoft Entra App ID.",
			// 			},
			// 			Resource: &schemaR.StringAttribute{
			// 				Computed: true,
			// 			},
			// 			DataSource: &schemaD.StringAttribute{
			// 				Computed: true,
			// 			},
			// 		},
			// 		"parent_principal_id": superschema.StringAttribute{
			// 			Common: &schemaR.StringAttribute{
			// 				MarkdownDescription: "The parent principal ID of Service Principal Profile.",
			// 			},
			// 			Resource: &schemaR.StringAttribute{
			// 				Computed: true,
			// 			},
			// 			DataSource: &schemaD.StringAttribute{
			// 				Computed: true,
			// 			},
			// 		},
			// 	},
			// },
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
