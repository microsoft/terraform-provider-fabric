// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package externaldatashareprovider

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
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

	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: markdownDescriptionR,
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: markdownDescriptionD,
		},
		Attributes: map[string]superschema.Attribute{
			"workspace_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Workspace ID.",
					CustomType:          customtypes.UUIDType{},
				},
				DataSource: &schemaD.StringAttribute{
					Required: true,
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
			},
			"item_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The item ID.",
					CustomType:          customtypes.UUIDType{},
				},
				DataSource: &schemaD.StringAttribute{
					Required: true,
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
			},
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Data access role.",
					CustomType:          customtypes.UUIDType{},
				},
				DataSource: &schemaD.StringAttribute{
					Required: !isList,
					Computed: isList,
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"paths": superschema.SuperSetAttribute{
				Common: &schemaR.SetAttribute{
					MarkdownDescription: "Allowed values for this attribute.",
					CustomType: supertypes.SetTypeOf[types.String]{
						SetType: basetypes.SetType{
							ElemType: types.StringType,
						},
					},
					ElementType: types.StringType,
				},
				DataSource: &schemaD.SetAttribute{
					Computed: true,
				},
				Resource: &schemaR.SetAttribute{
					Required: true,
					Validators: []validator.Set{
						setvalidator.SizeAtMost(100),
					},
				},
			},
			"principal_model": superschema.SuperSingleNestedAttributeOf[common.PrincipalModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The creator principal of the external data share.",
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Resource: &schemaR.SingleNestedAttribute{
					Required: false,
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"id": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The ID of the Data access role.",
							CustomType:          customtypes.UUIDType{},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
						Resource: &schemaR.StringAttribute{
							Required: false,
							Computed: true,
						},
					},
					"type": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The type of the creator principal.",
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
						Resource: &schemaR.StringAttribute{
							Required: false,
							Computed: true,
						},
					},
				},
			},
			"recipient": superschema.SuperSingleNestedAttributeOf[recipientModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The recipient of the external data share.",
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Resource: &schemaR.SingleNestedAttribute{
					Required: true,
				},
				Attributes: superschema.Attributes{
					"user_principal_name": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The user principal name of the recipient.",
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
					"tenant_id": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The tenant ID of the recipient.",
							CustomType:          customtypes.UUIDType{},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
						},
					},
				},
			},
			"status": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The status of the external data share.",
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleExternalDataShareStatusValues(), true)...),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
				Resource: &schemaR.StringAttribute{
					Required: false,
					Computed: true,
				},
			},
			"expiration_time": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The expiration time of the external data share in UTC.",
					CustomType:          timetypes.RFC3339Type{},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
				Resource: &schemaR.StringAttribute{
					Required: false,
					Computed: true,
				},
			},
			"invitation_url": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The invitation URL for the external data share.",
					CustomType:          customtypes.URLType{},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
				Resource: &schemaR.StringAttribute{
					Required: false,
					Computed: true,
				},
			},
			"accepted_by_tenant_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The tenant ID that accepted the external data share.",
					CustomType:          customtypes.UUIDType{},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
				Resource: &schemaR.StringAttribute{
					Required: false,
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
				DataSource: dsTimeout,
			},
		},
	}
}
