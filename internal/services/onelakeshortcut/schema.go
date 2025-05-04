// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0
package onelakeshortcut

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func itemSchema(isList bool) superschema.Schema { //revive:disable-line:flag-parameter
	var dsTimeout *superschema.DatasourceTimeoutAttribute

	if !isList {
		dsTimeout = &superschema.DatasourceTimeoutAttribute{
			Read: true,
		}
	}

	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false),
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, isList),
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " ID.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
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
					Required: true,

					Validators: []validator.String{
						stringvalidator.LengthAtMost(200),
						stringvalidator.LengthAtLeast(1),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
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
					Required: !isList,
					Computed: isList,
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
					Required: true,
					Computed: isList,
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
					MarkdownDescription: "A string representing the full path where the shortcut is created, including either \"Files\" or \"Tables\".",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.LengthAtMost(200),
						stringvalidator.LengthAtLeast(1),
						stringvalidator.LengthAtLeast(1),
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[^/].*`), // Regex to ensure the string does not start with a slash
							"Path must not start with a leading slash ('/').",
						),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
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
					Required: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"onelake": superschema.SuperSingleNestedAttributeOf[oneLakeModel]{
						Common: &schemaR.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "An object containing the properties of the target OneLake data source.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Optional: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: map[string]superschema.Attribute{
							"item_id": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target item ID",
									CustomType:          customtypes.UUIDType{},
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
					"adls_gen2": superschema.SuperSingleNestedAttributeOf[adlsGen2]{
						Common: &schemaR.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "An object containing the properties of the target ADLS Gen2 data source.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Computed: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: map[string]superschema.Attribute{
							"connection_id": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target connection ID",
								},

								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
							"location": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target location",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
							"subpath": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target subpath",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
						},
					},
					"amazon_s3": superschema.SuperSingleNestedAttributeOf[amazonS3]{
						Common: &schemaR.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "An object containing the properties of the target Amazon S3 data source.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Computed: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: map[string]superschema.Attribute{
							"connection_id": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target connection ID",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
							"location": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target location",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
							"subpath": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target subpath",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
							"bucket": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target bucket",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
						},
					},
					"google_cloud_storage": superschema.SuperSingleNestedAttributeOf[googleCloudStorage]{
						Common: &schemaR.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "An object containing the properties of the target Google Cloud Storage data source.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Computed: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: map[string]superschema.Attribute{
							"connection_id": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target connection ID",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
							"location": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target location",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
							"subpath": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target subpath",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
							"bucket": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target bucket",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
						},
					},
					"s3_compatible": superschema.SuperSingleNestedAttributeOf[s3Compatible]{
						Common: &schemaR.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "An object containing the properties of the target S3 compatible data source.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Computed: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: map[string]superschema.Attribute{
							"connection_id": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target connection ID",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
							"location": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target location",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
							"subpath": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target subpath",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
							"bucket": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target bucket",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
						},
					},
					"external_data_share": superschema.SuperSingleNestedAttributeOf[externalDataShare]{
						Common: &schemaR.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "An object containing the properties of the target external data share.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Computed: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: map[string]superschema.Attribute{
							"connection_id": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target connection ID",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
						},
					},
					"dataverse": superschema.SuperSingleNestedAttributeOf[dataverse]{
						Common: &schemaR.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "An object containing the properties of the target Dataverse data source.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Computed: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: map[string]superschema.Attribute{
							"connection_id": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target connection ID",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
							"environment_domain": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target environment domain",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
							"table_name": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target table name",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
							},
							"deltalake_folder": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Target delta lake folder",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										stringvalidator.LengthAtMost(200),
									},
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
					Delete: true,
					Update: true,
				},
				DataSource: dsTimeout,
			},
		},
	}
}
