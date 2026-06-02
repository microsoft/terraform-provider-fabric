// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false),
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, isList),
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
				DataSource: &schemaD.StringAttribute{
					Required: true,
				},
			},
			"item_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Fabric item.",
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
			"role_name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Data access role.",
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
			"kind": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The kind of the Data access role. Possible values: " + utils.ConvertStringSlicesToString(
						fabcore.PossibleDataAccessRoleKindValues(),
						true,
						true,
					) + ".",
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleDataAccessRoleKindValues(), true)...),
					},
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
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
					"effect": superschema.StringAttribute{
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
					"constraints": superschema.SuperSingleNestedAttributeOf[constraintsModel]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "Any constraints such as row or column level security that are applied to tables as part of this role. If not included, no constraints apply to any tables in the role.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Optional: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: superschema.Attributes{
							"columns": superschema.SuperSetNestedAttributeOf[columnConstraint]{
								Common: &schemaR.SetNestedAttribute{
									MarkdownDescription: "The array of column constraints applied to one or more tables in the data access role.",
								},
								Resource: &schemaR.SetNestedAttribute{
									Optional: true,
								},
								DataSource: &schemaD.SetNestedAttribute{
									Computed: true,
								},
								Attributes: superschema.Attributes{
									"column_action": superschema.SuperSetAttribute{
										Common: &schemaR.SetAttribute{
											ElementType: types.StringType,
											MarkdownDescription: "The array of actions applied to the column names. Possible values: " + utils.ConvertStringSlicesToString(
												fabcore.PossibleColumnActionValues(),
												true,
												true,
											) + ".",
											Validators: []validator.Set{
												setvalidator.ValueStringsAre(stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleColumnActionValues(), true)...)),
											},
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
									"column_effect": superschema.StringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "The effect given to the column names. Possible values: " + utils.ConvertStringSlicesToString(
												fabcore.PossibleColumnEffectValues(),
												true,
												true,
											) + ".",
											Validators: []validator.String{
												stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleColumnEffectValues(), true)...),
											},
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"column_names": superschema.SuperSetAttribute{
										Common: &schemaR.SetAttribute{
											ElementType:         types.StringType,
											MarkdownDescription: "An array of case sensitive column names. Use `*` to indicate all columns in the table.",
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
									"table_path": superschema.StringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "A relative file path specifying which table the column constraint applies to. Should be in the form `/Tables/{optionalSchema}/{tableName}`.",
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
							"rows": superschema.SuperSetNestedAttributeOf[rowConstraint]{
								Common: &schemaR.SetNestedAttribute{
									MarkdownDescription: "The array of row constraints applied to one or more tables in the data access role.",
								},
								Resource: &schemaR.SetNestedAttribute{
									Optional: true,
								},
								DataSource: &schemaD.SetNestedAttribute{
									Computed: true,
								},
								Attributes: superschema.Attributes{
									"table_path": superschema.StringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "A relative file path specifying which table the row constraint applies to. Should be in the form `/Tables/{optionalSchema}/{tableName}`.",
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"value": superschema.StringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "A T-SQL expression that is used to evaluate which rows the role members can see. Only a subset of T-SQL can be used as a predicate.",
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
							"attribute_name": superschema.StringAttribute{
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
			"members": superschema.SuperSingleNestedAttributeOf[memberModel]{
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
					"fabric_item_members": superschema.SuperSetNestedAttributeOf[fabricItemMember]{
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
							"source_path": superschema.StringAttribute{
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
					"microsoft_entra_members": superschema.SuperSetNestedAttributeOf[microsoftEntraMember]{
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
							"object_type": superschema.StringAttribute{
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
								Resource: &schemaR.StringAttribute{
									Required: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"tenant_id": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The tenant id.",
									CustomType:          customtypes.UUIDType{},
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
