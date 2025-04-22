// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0
package onelakeshortcut

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

	possibleTargetTypeValues := utils.RemoveSlicesByValues(fabcore.PossibleTypeValues(), []fabcore.Type{
		fabcore.TypeOneLake,
		fabcore.TypeAdlsGen2,
		fabcore.TypeAmazonS3,
		fabcore.TypeDataverse,
		fabcore.TypeExternalDataShare,
		fabcore.TypeGoogleCloudStorage,
		fabcore.TypeS3Compatible,
	})

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
					Computed: true,
				},
			},
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " name.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					Validators: []validator.String{
						stringvalidator.LengthAtMost(200),
						stringvalidator.LengthAtLeast(1),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Optional: !isList,
					Computed: true,
				},
			},
			"workspace_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Workspace ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: !isList,
					Computed: isList,
				},
			},
			"item_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Item ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: !isList,
					Computed: isList,
				},
			},
			"path": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " path.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					Validators: []validator.String{
						stringvalidator.LengthAtMost(200),
						stringvalidator.LengthAtLeast(1),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Optional: !isList,
					Computed: true,
				},
			},

			"target": superschema.SuperSingleNestedAttributeOf[targetModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "An object that contains the target datasource, and it must specify exactly one of the supported destinations: OneLake, Amazon S3, ADLS Gen2, Google Cloud Storage, S3 compatible or Dataverse.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Computed: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"type": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The " + ItemTypeInfo.Name + " type.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,

							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(possibleTargetTypeValues, true)...),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleGatewayTypeValues(), true)...),
							},
						},
					},
					"onelake": superschema.SuperSingleNestedAttributeOf[oneLakeModel]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "The OneLake datasource.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Computed: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: map[string]superschema.Attribute{
							"item_id": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target item ID",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"path": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target path",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"workspace_id": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target Workspace ID",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
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
				DataSource: dsTimeout,
			},
		},
	}
}
