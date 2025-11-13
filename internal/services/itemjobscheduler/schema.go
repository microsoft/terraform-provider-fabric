// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package itemjobscheduler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema" //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	superint32validator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/int32validator"
	superobjectvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/objectvalidator"
	supersetvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/setvalidator"
	superstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

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
					Required:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"item_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The item ID.",
					CustomType:          customtypes.UUIDType{},
					Required:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"job_type": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The job type. " + allowedJobTypesMarkdownDescription(),
					Required:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
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

func ownerSchema() superschema.SuperSingleNestedAttributeOf[principalModel] {
	return superschema.SuperSingleNestedAttributeOf[principalModel]{
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
			"type": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The type of the principal.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func configurationSchema() superschema.SuperSingleNestedAttributeOf[configurationModel] {
	return superschema.SuperSingleNestedAttributeOf[configurationModel]{
		Common: &schemaR.SingleNestedAttribute{
			MarkdownDescription: "The schedule configuration.",
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
							regexp.MustCompile(`Z$`),
							"The time must be in UTC, using the YYYY-MM-DDTHH:mm:ssZ format.",
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
							regexp.MustCompile(`Z$`),
							"The time must be in UTC, using the YYYY-MM-DDTHH:mm:ssZ format.",
						),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"type": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A string represents the type of the plan.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleScheduleTypeValues(), true)...),
					},
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
						superint32validator.RequireIfAttributeIsOneOf(path.MatchRoot("configuration").AtName("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.ScheduleTypeCron)),
							}),
						superint32validator.NullIfAttributeIsOneOf(path.MatchRoot("configuration").AtName("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.ScheduleTypeDaily)),
								types.StringValue(string(fabcore.ScheduleTypeWeekly)),
								types.StringValue(string(fabcore.ScheduleTypeMonthly)),
							}),
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
					Optional:    true,
					ElementType: types.StringType,
					Validators: []validator.Set{
						setvalidator.SizeAtMost(100),
						setvalidator.ValueStringsAre(
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`),
								"Each time entry must be in hh:mm format.",
							),
						),
						supersetvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("configuration").AtName("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.ScheduleTypeDaily)),
								types.StringValue(string(fabcore.ScheduleTypeWeekly)),
								types.StringValue(string(fabcore.ScheduleTypeMonthly)),
							}),
						supersetvalidator.NullIfAttributeIsOneOf(path.MatchRoot("configuration").AtName("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.ScheduleTypeCron)),
							}),
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
						setvalidator.SizeAtMost(7),
						setvalidator.ValueStringsAre(
							stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleDayOfWeekValues(), true)...),
						),
						supersetvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("configuration").AtName("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.ScheduleTypeWeekly)),
							}),
						supersetvalidator.NullIfAttributeIsOneOf(path.MatchRoot("configuration").AtName("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.ScheduleTypeCron)),
								types.StringValue(string(fabcore.ScheduleTypeDaily)),
								types.StringValue(string(fabcore.ScheduleTypeMonthly)),
							}),
					},
				},
				DataSource: &schemaD.SetAttribute{
					Computed: true,
				},
			},
			"recurrence": superschema.Int32Attribute{
				Common: &schemaR.Int32Attribute{
					MarkdownDescription: "Specifies the monthly job repeat interval. For example, when set to 1 the job is triggered every month.",
				},
				Resource: &schemaR.Int32Attribute{
					Optional: true,
					Validators: []validator.Int32{
						int32validator.Between(1, 12),
						superint32validator.RequireIfAttributeIsOneOf(path.MatchRoot("configuration").AtName("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.ScheduleTypeMonthly)),
							}),
						superint32validator.NullIfAttributeIsOneOf(path.MatchRoot("configuration").AtName("type"),
							[]attr.Value{
								types.StringValue(string(fabcore.ScheduleTypeDaily)),
								types.StringValue(string(fabcore.ScheduleTypeWeekly)),
								types.StringValue(string(fabcore.ScheduleTypeCron)),
							}),
					},
				},
				DataSource: &schemaD.Int32Attribute{
					Computed: true,
				},
			},
			"occurrence": occurrenceSchema(),
		},
	}
}

func occurrenceSchema() superschema.SuperSingleNestedAttributeOf[occurrenceModel] {
	return superschema.SuperSingleNestedAttributeOf[occurrenceModel]{
		Common: &schemaR.SingleNestedAttribute{
			MarkdownDescription: "A date for triggering the job.",
		},
		Resource: &schemaR.SingleNestedAttribute{
			Optional: true,
			Validators: []validator.Object{
				superobjectvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("configuration").AtName("type"),
					[]attr.Value{
						types.StringValue(string(fabcore.ScheduleTypeMonthly)),
					}),
				superobjectvalidator.NullIfAttributeIsOneOf(path.MatchRoot("configuration").AtName("type"),
					[]attr.Value{
						types.StringValue(string(fabcore.ScheduleTypeDaily)),
						types.StringValue(string(fabcore.ScheduleTypeWeekly)),
						types.StringValue(string(fabcore.ScheduleTypeCron)),
					}),
			},
		},
		DataSource: &schemaD.SingleNestedAttribute{
			Computed: true,
		},
		Attributes: map[string]superschema.Attribute{
			"occurrence_type": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "An enumerator that lists the day for triggering jobs.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleOccurrenceTypeValues(), true)...),
						superstringvalidator.RequireIfAttributeIsSet(path.MatchRelative().AtParent()),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"day_of_month": superschema.Int32Attribute{
				Common: &schemaR.Int32Attribute{
					MarkdownDescription: "Specifies a date to trigger the job, using a value between 1 and 31. For example, 2 means the second day of the month. The date must be valid. If an invalid date is provided, such as February 31st, it will automatically skip to the month that includes the 31st day.",
				},
				Resource: &schemaR.Int32Attribute{
					Optional: true,
					Validators: []validator.Int32{
						int32validator.Between(1, 31),
						superint32validator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("occurrence_type"),
							[]attr.Value{
								types.StringValue(string(fabcore.OccurrenceTypeDayOfMonth)),
							}),
						superint32validator.NullIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("occurrence_type"),
							[]attr.Value{
								types.StringValue(string(fabcore.OccurrenceTypeOrdinalWeekday)),
							}),
					},
				},
				DataSource: &schemaD.Int32Attribute{
					Computed: true,
				},
			},
			"week_index": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The week of the month.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleWeekIndexValues(), true)...),
						superstringvalidator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("occurrence_type"),
							[]attr.Value{
								types.StringValue(string(fabcore.OccurrenceTypeOrdinalWeekday)),
							}),
						superstringvalidator.NullIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("occurrence_type"),
							[]attr.Value{
								types.StringValue(string(fabcore.OccurrenceTypeDayOfMonth)),
							}),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"weekday": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Week day for triggering jobs.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabcore.PossibleDayOfWeekValues(), true)...),
						superstringvalidator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("occurrence_type"),
							[]attr.Value{
								types.StringValue(string(fabcore.OccurrenceTypeOrdinalWeekday)),
							}),
						superstringvalidator.NullIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("occurrence_type"),
							[]attr.Value{
								types.StringValue(string(fabcore.OccurrenceTypeDayOfMonth)),
							}),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func allowedJobTypesMarkdownDescription() string {
	result := "Allowed job types per item type: "
	for k, v := range AllowedJobTypesByItemType {
		result += fmt.Sprintf("%s: {%s}; ", k, strings.Join(v, ", "))
	}

	return result
}
