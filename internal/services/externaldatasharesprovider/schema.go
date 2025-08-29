package externaldatasharesprovider

import (
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/types"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func itemSchema() superschema.Schema {
	markdownDescriptionR := fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false)
	markdownDescriptionD := fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, true)

	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: markdownDescriptionR,
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: markdownDescriptionD,
		},
		Attributes: map[string]superschema.Attribute{
			"external_data_share_id": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The external data share ID.",
					CustomType:          customtypes.UUIDType{},
					Required:            true,
				},
			},
			"workspace_id": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The Workspace ID.",
					CustomType:          customtypes.UUIDType{},
					Required:            true,
				},
			},
			"item_id": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The item ID.",
					CustomType:          customtypes.UUIDType{},
					Required:            true,
				},
			},
			"value": superschema.SuperSetNestedAttributeOf[externalDataSharesModel]{
				DataSource: &schemaD.SetNestedAttribute{
					MarkdownDescription: "Map of external data share values.",
					Computed:            true,
				},
				Attributes: superschema.Attributes{
					"id": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The name of the Data access role.",
							CustomType:          customtypes.UUIDType{},
							Computed:            true,
						},
					},
					"paths": superschema.SuperSetAttribute{
						DataSource: &schemaD.SetAttribute{
							MarkdownDescription: "Allowed values for this attribute.",
							ElementType:         types.StringType,
							Computed:            true,
						},
					},
					"creator_principal": superschema.SuperSingleNestedAttributeOf[creatorPrincipalModel]{
						DataSource: &schemaD.SingleNestedAttribute{
							MarkdownDescription: "The creator principal of the external data share.",
							Computed:            true,
						},
						Attributes: superschema.Attributes{
							"id": superschema.SuperStringAttribute{
								DataSource: &schemaD.StringAttribute{
									MarkdownDescription: "The ID of the creator principal.",
									CustomType:          customtypes.UUIDType{},
									Computed:            true,
								},
							},
							"display_name": superschema.SuperStringAttribute{
								DataSource: &schemaD.StringAttribute{
									MarkdownDescription: "The display name of the creator principal.",
									Computed:            true,
								},
							},
							"type": superschema.SuperStringAttribute{
								DataSource: &schemaD.StringAttribute{
									MarkdownDescription: "The type of the creator principal.",
									Computed:            true,
								},
							},
							"user_details": superschema.SuperSingleNestedAttributeOf[userDetailsModel]{
								DataSource: &schemaD.SingleNestedAttribute{
									MarkdownDescription: "The user details of the creator principal.",
									Computed:            true,
								},
								Attributes: superschema.Attributes{
									"user_principal_name": superschema.SuperStringAttribute{
										DataSource: &schemaD.StringAttribute{
											MarkdownDescription: "The user principal name of the creator principal.",
											Computed:            true,
										},
									},
								},
							},
						},
					},
					"recipient": superschema.SuperSingleNestedAttributeOf[recipientModel]{
						DataSource: &schemaD.SingleNestedAttribute{
							MarkdownDescription: "The recipient of the external data share.",
							Computed:            true,
						},
						Attributes: superschema.Attributes{
							"user_principal_name": superschema.SuperStringAttribute{
								DataSource: &schemaD.StringAttribute{
									MarkdownDescription: "The user principal name of the recipient.",
									Computed:            true,
								},
							},
						},
					},
					"status": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The status of the external data share.",
							Computed:            true,
						},
					},
					"expiration_time_utc": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The expiration time of the external data share in UTC.",
							Computed:            true,
						},
					},
					"workspace_id": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The Workspace ID.",
							CustomType:          customtypes.UUIDType{},
							Computed:            true,
						},
					},
					"item_id": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The item ID.",
							CustomType:          customtypes.UUIDType{},
							Computed:            true,
						},
					},
					"invitation_url": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The invitation URL for the external data share.",
							Computed:            true,
						},
					},
					"accepted_by_tenant_id": superschema.SuperStringAttribute{
						DataSource: &schemaD.StringAttribute{
							MarkdownDescription: "The tenant ID that accepted the external data share.",
							CustomType:          customtypes.UUIDType{},
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
