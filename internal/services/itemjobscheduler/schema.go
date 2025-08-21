// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0
package itemjobscheduler

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
					Required: true,
				},
			},
			"item_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The item ID.",
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
			"job_type": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The job type.",
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
			"enabled": superschema.BoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: `Whether this schedule is enabled. True - Enabled, False - Disabled.`,
				},
				Resource: &schemaR.BoolAttribute{
					Required: true,
				},
				DataSource: &schemaD.BoolAttribute{
					Computed: true,
				},
			},
			"created_date_time": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The created time stamp of this schedule in UTC, using the YYYY-MM-DDTHH:mm:ssZ format.",
					CustomType:          timetypes.RFC3339Type{},
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"owner":         ownerSchema(),
			"configuration": configurationSchema(),
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

func ownerSchema() superschema.SuperSingleNestedAttributeOf[baseOwnerModel] {
	return superschema.SuperSingleNestedAttributeOf[baseOwnerModel]{
		Common: &schemaR.SingleNestedAttribute{
			MarkdownDescription: "The user identity that created this schedule or last modified.",
		},
		Resource: &schemaR.SingleNestedAttribute{
			Computed: true,
		},
		DataSource: &schemaD.SingleNestedAttribute{
			Computed: true,
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The principal's ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
				},

				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},

			"display_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The principal's display name.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
				},
			},
			"group_details": superschema.SuperSingleNestedAttributeOf[groupDetailsModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Group specific details. Applicable when the principal type is Group.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Computed: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Optional: true,
				},
				Attributes: map[string]superschema.Attribute{
					"group_type": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The type of the group. Additional group types may be added over time.",
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"service_principal_details": superschema.SuperSingleNestedAttributeOf[servicePrincipalDetailsModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Service principal specific details. Applicable when the principal type is ServicePrincipal.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Computed: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Optional: true,
				},
				Attributes: map[string]superschema.Attribute{
					"aad_app_id": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The service principal's Microsoft Entra AppId.",
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"service_principal_profile_details": superschema.SuperSingleNestedAttributeOf[servicePrincipalProfileDetailsModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Service principal profile details. Applicable when the principal type is ServicePrincipalProfile.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Computed: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Optional: true,
				},
				Attributes: map[string]superschema.Attribute{
					"parent_principal": superschema.SuperSingleNestedAttributeOf[servicePrincipalBaseOwnerModel]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "The service principal profile's parent principal.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Computed: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: map[string]superschema.Attribute{
							"id": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The principal's ID.",
									CustomType:          customtypes.UUIDType{},
								},
								Resource: &schemaR.StringAttribute{
									Computed: true,
								},

								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},

							"display_name": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The principal's display name.",
								},
								Resource: &schemaR.StringAttribute{
									Computed: true,
								},
								DataSource: &schemaD.StringAttribute{
									Optional: true,
								},
							},
							"service_principal_details": superschema.SuperSingleNestedAttributeOf[servicePrincipalDetailsModel]{
								Common: &schemaR.SingleNestedAttribute{
									MarkdownDescription: "Service principal specific details. Applicable when the principal type is ServicePrincipal.",
								},
								Resource: &schemaR.SingleNestedAttribute{
									Computed: true,
								},
								DataSource: &schemaD.SingleNestedAttribute{
									Optional: true,
								},
								Attributes: map[string]superschema.Attribute{
									"aad_app_id": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "The service principal's Microsoft Entra AppId.",
										},
										Resource: &schemaR.StringAttribute{
											Computed: true,
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
								},
							},
							"type": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The type of the principal. Additional principal types may be added over time.",
								},
								Resource: &schemaR.StringAttribute{
									Computed: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
			"type": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The type of the principal. Additional principal types may be added over time.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"user_details": superschema.SuperSingleNestedAttributeOf[userDetailsModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "User principal specific details. Applicable when the principal type is User.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Computed: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Optional: true,
				},
				Attributes: map[string]superschema.Attribute{
					"user_principal_name": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The user principal name.",
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func configurationSchema() superschema.SuperSingleNestedAttributeOf[baseConfigurationModel] {
	return superschema.SuperSingleNestedAttributeOf[baseConfigurationModel]{
		Common: &schemaR.SingleNestedAttribute{
			MarkdownDescription: "The actual data contains the time/weekdays of this schedule.",
		},
		Resource: &schemaR.SingleNestedAttribute{
			Required: true,
		},
		DataSource: &schemaD.SingleNestedAttribute{
			Computed: true,
		},
		Attributes: map[string]superschema.Attribute{
			"start_date_time": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The start time for this schedule. If the start time is in the past, it will trigger a job instantly. The time is in UTC, using the YYYY-MM-DDTHH:mm:ssZ format.",
					CustomType:          timetypes.RFC3339Type{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"The datetime must be in UTC format ending with 'Z'. Timezone offsets are not allowed.",
						),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"end_date_time": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The end time for this schedule. The end time must be later than the start time. It has to be in UTC, using the YYYY-MM-DDTHH:mm:ssZ format.",
					CustomType:          timetypes.RFC3339Type{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"The datetime must be in UTC format ending with 'Z'. Timezone offsets are not allowed.",
						),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"type": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A string represents the type of the plan. Additional planType types may be added over time.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"interval": superschema.Int32Attribute{
				Common: &schemaR.Int32Attribute{
					MarkdownDescription: "The time interval in minutes. A number between 1 and 5270400 (10 years).",
				},
				Resource: &schemaR.Int32Attribute{
					Optional: true,
					Validators: []validator.Int32{
						int32validator.Between(1, 5270400),
					},
				},
				DataSource: &schemaD.Int32Attribute{
					Computed: true,
				},
			},
			"times": superschema.SuperSetAttribute{
				Common: &schemaR.SetAttribute{
					MarkdownDescription: "A list of time slots in hh:mm format, at most 100 elements are allowed.",
					ElementType:         types.StringType,
				},
				Resource: &schemaR.SetAttribute{
					Optional: true,
					Validators: []validator.Set{
						setvalidator.SizeAtMost(100),
					},
				},
				DataSource: &schemaD.SetAttribute{
					Computed: true,
				},
			},
			"weekdays": superschema.SuperSetAttribute{
				Common: &schemaR.SetAttribute{
					MarkdownDescription: "A list of weekdays, at most seven elements are allowed.",
					ElementType:         types.StringType,
				},
				Resource: &schemaR.SetAttribute{
					Optional: true,
					Validators: []validator.Set{
						setvalidator.SizeAtMost(10),
					},
				},
				DataSource: &schemaD.SetAttribute{
					Computed: true,
				},
			},
		},
	}
}
