// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
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
			"workspace_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Workspace ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Required: true,
				},
			},
			"item_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Fabric item to put the roles.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Required: true,
				},
			},
			"etag": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ETag of the item.",
					CustomType:          supertypes.StringType{},
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
				},
			},
			"value": superschema.SuperSetNestedAttributeOf[dataAccessRole]{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "Map of data access roles.",
				},
				Resource: &schemaR.SetNestedAttribute{
					Required: true,
				},
				DataSource: &schemaD.SetNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"name": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The name of the Data access role.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"decision_rules": superschema.SuperSetNestedAttributeOf[decisionRule]{
						Common: &schemaR.SetNestedAttribute{
							MarkdownDescription: "The array of permissions that make up the Data access role.",
						},
						Resource: &schemaR.SetNestedAttribute{
							Required: true,
						},
						DataSource: &schemaD.SetNestedAttribute{
							Computed: true,
						},
						Attributes: superschema.Attributes{
							"effect": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The effect that a role has on access to the data resource. Possible values: " + utils.ConvertStringSlicesToString(
										fabcore.PossibleEffectValues(),
										true,
										true,
									) + ".",
									Validators: []validator.String{
										stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleEffectValues(), true)...),
									},
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"permission": superschema.SuperSetNestedAttributeOf[permissionScope]{
								Common: &schemaR.SetNestedAttribute{
									MarkdownDescription: "Permissions defined by attribute name and values.",
								},
								Resource: &schemaR.SetNestedAttribute{
									Required: true,
								},
								DataSource: &schemaD.SetNestedAttribute{
									Computed: true,
								},
								Attributes: superschema.Attributes{
									"attribute_name": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "The name of the attribute that is being evaluated for access permissions. Possible values: " + utils.ConvertStringSlicesToString(
												fabcore.PossibleAttributeNameValues(),
												true,
												true,
											) + ".",
											Validators: []validator.String{
												stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleAttributeNameValues(), true)...),
											},
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"attribute_value_included_in": superschema.SuperSetAttribute{
										Common: &schemaR.SetAttribute{
											ElementType:         types.StringType,
											MarkdownDescription: "Allowed values for this attribute.",
										},
										Resource: &schemaR.SetAttribute{
											ElementType: types.StringType,
											Required:    true,
										},
										DataSource: &schemaD.SetAttribute{
											ElementType: types.StringType,
											Computed:    true,
										},
									},
								},
							},
						},
					},
					"members": superschema.SingleNestedAttribute{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "The members object which contains the members of the role as arrays of different member types.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Required: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: superschema.Attributes{
							"fabric_item_members": superschema.SuperSetNestedAttributeOf[FabricItemMember]{
								Common: &schemaR.SetNestedAttribute{
									MarkdownDescription: "Fabric-scoped members with path-based access.",
								},
								Resource: &schemaR.SetNestedAttribute{
									Optional: true,
								},
								DataSource: &schemaD.SetNestedAttribute{
									Computed: true,
								},
								Attributes: superschema.Attributes{
									"item_access": superschema.SuperSetAttribute{
										Common: &schemaR.SetAttribute{
											ElementType:         types.StringType,
											MarkdownDescription: "Permissions for the item.",
										},
										Resource: &schemaR.SetAttribute{
											ElementType: types.StringType,
											Required:    true,
										},
										DataSource: &schemaD.SetAttribute{
											ElementType: types.StringType,
											Computed:    true,
										},
									},
									"source_path": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "The path to Fabric item having the specified item access.",
											Validators: []validator.String{
												stringvalidator.RegexMatches(
													regexp.MustCompile(`^[{]?[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}[}]?/[{]?[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}[}]?$`),
													"Source path must be a combination of two GUIDs separated by a slash.",
												),
											},
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
							"microsoft_entra_members": superschema.SuperSetNestedAttributeOf[MicrosoftEntraMember]{
								Common: &schemaR.SetNestedAttribute{
									MarkdownDescription: "The list of Microsoft Entra ID members.",
								},
								Resource: &schemaR.SetNestedAttribute{
									Optional: true,
								},
								DataSource: &schemaD.SetNestedAttribute{
									Computed: true,
								},
								Attributes: superschema.Attributes{
									"object_id": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "The object id.",
											CustomType:          customtypes.UUIDType{},
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"object_type": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "The type of Microsoft Entra ID object. Possible values: " + utils.ConvertStringSlicesToString(
												fabcore.PossibleObjectTypeValues(),
												true,
												true,
											) + ".",
											Validators: []validator.String{
												stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleObjectTypeValues(), true)...),
											},
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
										},
									},
									"tenant_id": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "The tenant id.",
											CustomType:          customtypes.UUIDType{},
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
