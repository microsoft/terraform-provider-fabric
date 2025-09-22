// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package itemjobscheduler

import (
	"context"
	"fmt"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseItemJobSchedulerModel struct {
	ID              customtypes.UUID                                         `tfsdk:"id"`
	Enabled         types.Bool                                               `tfsdk:"enabled"`
	ItemID          customtypes.UUID                                         `tfsdk:"item_id"`
	WorkspaceID     customtypes.UUID                                         `tfsdk:"workspace_id"`
	JobType         types.String                                             `tfsdk:"job_type"`
	CreatedDateTime timetypes.RFC3339                                        `tfsdk:"created_date_time"`
	Owner           supertypes.SingleNestedObjectValueOf[principalModel]     `tfsdk:"owner"`
	Configuration   supertypes.SingleNestedObjectValueOf[configurationModel] `tfsdk:"configuration"`
}

type configurationModel struct {
	StartDateTime timetypes.RFC3339                                     `tfsdk:"start_date_time"`
	EndDateTime   timetypes.RFC3339                                     `tfsdk:"end_date_time"`
	Type          types.String                                          `tfsdk:"type"`
	Interval      types.Int32                                           `tfsdk:"interval"`   // Cron
	Times         supertypes.SetValueOf[types.String]                   `tfsdk:"times"`      // Daily, Weekly and Monthly
	Weekdays      supertypes.SetValueOf[types.String]                   `tfsdk:"weekdays"`   // Weekly
	Occurrence    supertypes.SingleNestedObjectValueOf[occurrenceModel] `tfsdk:"occurrence"` // Monthly
	Recurrence    types.Int32                                           `tfsdk:"recurrence"` // Monthly
}

type occurrenceModel struct {
	DayOfMonth     types.Int32  `tfsdk:"day_of_month"`
	OccurrenceType types.String `tfsdk:"occurrence_type"`
	Weekday        types.String `tfsdk:"weekday"`
	WeekIndex      types.String `tfsdk:"week_index"`
}

func (to *baseItemJobSchedulerModel) set(ctx context.Context, workspaceID, itemID, jobType string, from fabcore.ItemSchedule) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Enabled = types.BoolPointerValue(from.Enabled)
	to.ItemID = customtypes.NewUUIDValue(itemID)
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.JobType = types.StringValue(jobType)
	to.CreatedDateTime = timetypes.NewRFC3339TimePointerValue(from.CreatedDateTime)

	configuration := supertypes.NewSingleNestedObjectValueOfNull[configurationModel](ctx)
	principal := supertypes.NewSingleNestedObjectValueOfNull[principalModel](ctx)

	if from.Owner != nil {
		principalModel := &principalModel{}

		principalModel.set(*from.Owner)

		if diags := principal.Set(ctx, principalModel); diags.HasError() {
			return diags
		}
	}

	to.Owner = principal

	if from.Configuration != nil {
		baseJobScheduleConfigurationModel := &configurationModel{}
		if diags := baseJobScheduleConfigurationModel.set(ctx, from.Configuration); diags.HasError() {
			return diags
		}

		if diags := configuration.Set(ctx, baseJobScheduleConfigurationModel); diags.HasError() {
			return diags
		}
	}

	to.Configuration = configuration

	return nil
}

func (to *configurationModel) set(ctx context.Context, from fabcore.ScheduleConfigClassification) diag.Diagnostics {
	schConfig := from.GetScheduleConfig()
	to.StartDateTime = timetypes.NewRFC3339TimePointerValue(schConfig.StartDateTime)
	to.EndDateTime = timetypes.NewRFC3339TimePointerValue(schConfig.EndDateTime)
	to.Type = types.StringPointerValue((*string)(schConfig.Type))
	to.Times = supertypes.NewSetValueOfNull[types.String](ctx)
	to.Weekdays = supertypes.NewSetValueOfNull[types.String](ctx)
	to.Occurrence = supertypes.NewSingleNestedObjectValueOfNull[occurrenceModel](ctx)
	to.Recurrence = types.Int32Null()
	to.Interval = types.Int32Null()

	switch entity := from.(type) {
	case *fabcore.CronScheduleConfig:
		to.Interval = types.Int32PointerValue(entity.Interval)
	case *fabcore.DailyScheduleConfig:
		times := make([]types.String, len(entity.Times))

		for i, t := range entity.Times {
			times[i] = types.StringValue(t)
		}

		to.Times.Set(ctx, times)

	case *fabcore.WeeklyScheduleConfig:
		times := make([]types.String, len(entity.Times))
		weekdays := make([]types.String, len(entity.Weekdays))

		for i, t := range entity.Times {
			times[i] = types.StringValue(t)
		}

		to.Times.Set(ctx, times)

		for i, w := range entity.Weekdays {
			weekdays[i] = types.StringValue(string(w))
		}

		to.Weekdays.Set(ctx, weekdays)
	case *fabcore.MonthlyScheduleConfig:
		to.Recurrence = types.Int32PointerValue(entity.Recurrence)
		times := make([]types.String, len(entity.Times))

		for i, t := range entity.Times {
			times[i] = types.StringValue(t)
		}

		to.Times.Set(ctx, times)

		occurrence := supertypes.NewSingleNestedObjectValueOfNull[occurrenceModel](ctx)
		if entity.Occurrence != nil {
			occurrenceModel := &occurrenceModel{}
			if diags := occurrenceModel.set(entity.Occurrence); diags.HasError() {
				return diags
			}

			if diags := occurrence.Set(ctx, occurrenceModel); diags.HasError() {
				return diags
			}
		}

		to.Occurrence = occurrence

	default:
		var diags diag.Diagnostics

		diags.AddError(
			"Unsupported Configuration type",
			fmt.Sprintf("The Configuration type '%T' is not supported.", entity),
		)

		return diags
	}

	return nil
}

func (to *occurrenceModel) set(from fabcore.MonthlyOccurrenceClassification) diag.Diagnostics {
	monthlyOcc := from.GetMonthlyOccurrence()
	to.OccurrenceType = types.StringPointerValue((*string)(monthlyOcc.OccurrenceType))

	switch entity := from.(type) {
	case *fabcore.DayOfMonth:
		to.DayOfMonth = types.Int32PointerValue(entity.DayOfMonth)
	case *fabcore.OrdinalWeekday:
		to.Weekday = types.StringPointerValue((*string)(entity.Weekday))
		to.WeekIndex = types.StringPointerValue((*string)(entity.WeekIndex))
	default:
		var diags diag.Diagnostics

		diags.AddError(
			"Unsupported Monthly Occurrence type",
			fmt.Sprintf("The Monthly Occurrence type '%T' is not supported.", entity),
		)

		return diags
	}

	return nil
}

/*
DATA-SOURCE
*/

type dataSourceJobScheduleModel struct {
	baseItemJobSchedulerModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceJobSchedulesModel struct {
	WorkspaceID customtypes.UUID                                             `tfsdk:"workspace_id"`
	ItemID      customtypes.UUID                                             `tfsdk:"item_id"`
	JobType     types.String                                                 `tfsdk:"job_type"`
	Values      supertypes.SetNestedObjectValueOf[baseItemJobSchedulerModel] `tfsdk:"values"`
	Timeouts    timeoutsD.Value                                              `tfsdk:"timeouts"`
}

func (to *dataSourceJobSchedulesModel) setValues(ctx context.Context, workspaceID, itemID, jobType string, from []fabcore.ItemSchedule) diag.Diagnostics {
	slice := make([]*baseItemJobSchedulerModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseItemJobSchedulerModel

		if diags := entityModel.set(ctx, workspaceID, itemID, jobType, entity); diags.HasError() {
			return diags
		}

		slice = append(slice, &entityModel)
	}

	return to.Values.Set(ctx, slice)
}

/*
RESOURCE
*/

type resourceJobScheduleModel struct {
	baseItemJobSchedulerModel

	Timeouts timeoutsR.Value `tfsdk:"timeouts"`
}

type requestCreateJobSchedule struct {
	fabcore.CreateScheduleRequest
}

func (to *requestCreateJobSchedule) set(ctx context.Context, from resourceJobScheduleModel) diag.Diagnostics {
	configuration, diags := from.Configuration.Get(ctx)
	if diags.HasError() {
		return diags
	}

	var reqConfiguration fabcore.ScheduleConfigClassification
	localTimeZoneID := "Central Standard Time"
	configurationType := (fabcore.ScheduleType)(configuration.Type.ValueString())
	startDateTime, startDiags := configuration.StartDateTime.ValueRFC3339Time()
	endDateTime, endDiags := configuration.EndDateTime.ValueRFC3339Time()

	diags.Append(startDiags...)
	diags.Append(endDiags...)

	if diags.HasError() {
		return diags
	}

	switch configurationType {
	case fabcore.ScheduleTypeCron:
		reqConfiguration = &fabcore.CronScheduleConfig{
			Interval:        configuration.Interval.ValueInt32Pointer(),
			StartDateTime:   &startDateTime,
			EndDateTime:     &endDateTime,
			LocalTimeZoneID: &localTimeZoneID,
			Type:            &configurationType,
		}
	case fabcore.ScheduleTypeDaily:
		times, diags := configuration.Times.Get(ctx)
		if diags.HasError() {
			return diags
		}

		timesSlice := make([]string, 0, len(times))
		for _, t := range times {
			timesSlice = append(timesSlice, t.ValueString())
		}

		reqConfiguration = &fabcore.DailyScheduleConfig{
			Times:           timesSlice,
			StartDateTime:   &startDateTime,
			EndDateTime:     &endDateTime,
			LocalTimeZoneID: &localTimeZoneID,
			Type:            &configurationType,
		}
	case fabcore.ScheduleTypeWeekly:
		times, diags := configuration.Times.Get(ctx)
		if diags.HasError() {
			return diags
		}

		timesSlice := make([]string, 0, len(times))
		for _, t := range times {
			timesSlice = append(timesSlice, t.ValueString())
		}

		weekdays, diags := configuration.Weekdays.Get(ctx)
		if diags.HasError() {
			return diags
		}

		weekdaysSlice := make([]fabcore.DayOfWeek, 0, len(weekdays))
		for _, w := range weekdays {
			weekdaysSlice = append(weekdaysSlice, fabcore.DayOfWeek(w.ValueString()))
		}

		reqConfiguration = &fabcore.WeeklyScheduleConfig{
			Times:           timesSlice,
			Weekdays:        weekdaysSlice,
			StartDateTime:   &startDateTime,
			EndDateTime:     &endDateTime,
			LocalTimeZoneID: &localTimeZoneID,
			Type:            &configurationType,
		}
	default:
		diags.AddError(
			"Unsupported Configuration type",
			fmt.Sprintf("The Configuration type '%T' is not supported.", configurationType),
		)

		return diags
	}

	to.Configuration = reqConfiguration
	to.Enabled = from.Enabled.ValueBoolPointer()

	return nil
}

type requestUpdateJobSchedule struct {
	fabcore.UpdateScheduleRequest
}

func (to *requestUpdateJobSchedule) set(ctx context.Context, from resourceJobScheduleModel) diag.Diagnostics {
	configuration, diags := from.Configuration.Get(ctx)

	if diags.HasError() {
		return diags
	}

	var reqConfiguration fabcore.ScheduleConfigClassification
	localTimeZoneID := "Central Standard Time"

	configurationType := (fabcore.ScheduleType)(configuration.Type.ValueString())
	startDateTime, startDiags := configuration.StartDateTime.ValueRFC3339Time()
	endDateTime, endDiags := configuration.EndDateTime.ValueRFC3339Time()

	diags.Append(startDiags...)
	diags.Append(endDiags...)

	if diags.HasError() {
		return diags
	}

	switch configurationType {
	case fabcore.ScheduleTypeCron:
		reqConfiguration = &fabcore.CronScheduleConfig{
			Interval:        configuration.Interval.ValueInt32Pointer(),
			StartDateTime:   &startDateTime,
			EndDateTime:     &endDateTime,
			LocalTimeZoneID: &localTimeZoneID,
			Type:            &configurationType,
		}
	case fabcore.ScheduleTypeDaily:
		times, diags := configuration.Times.Get(ctx)
		if diags.HasError() {
			return diags
		}

		timesSlice := make([]string, 0, len(times))
		for _, t := range times {
			timesSlice = append(timesSlice, t.ValueString())
		}

		reqConfiguration = &fabcore.DailyScheduleConfig{
			Times:           timesSlice,
			StartDateTime:   &startDateTime,
			EndDateTime:     &endDateTime,
			LocalTimeZoneID: &localTimeZoneID,
			Type:            &configurationType,
		}
	case fabcore.ScheduleTypeWeekly:
		times, diags := configuration.Times.Get(ctx)
		if diags.HasError() {
			return diags
		}

		timesSlice := make([]string, 0, len(times))
		for _, t := range times {
			timesSlice = append(timesSlice, t.ValueString())
		}

		weekdays, diags := configuration.Weekdays.Get(ctx)
		if diags.HasError() {
			return diags
		}

		weekdaysSlice := make([]fabcore.DayOfWeek, 0, len(weekdays))
		for _, w := range weekdays {
			weekdaysSlice = append(weekdaysSlice, fabcore.DayOfWeek(w.ValueString()))
		}

		reqConfiguration = &fabcore.WeeklyScheduleConfig{
			Times:           timesSlice,
			Weekdays:        weekdaysSlice,
			StartDateTime:   &startDateTime,
			EndDateTime:     &endDateTime,
			LocalTimeZoneID: &localTimeZoneID,
			Type:            &configurationType,
		}
	case fabcore.ScheduleTypeMonthly:
		times, diags := configuration.Times.Get(ctx)
		if diags.HasError() {
			return diags
		}

		timesSlice := make([]string, 0, len(times))
		for _, t := range times {
			timesSlice = append(timesSlice, t.ValueString())
		}

		occurrence, diags := configuration.Occurrence.Get(ctx)
		if diags.HasError() {
			return diags
		}

		occurrenceType := (fabcore.OccurrenceType)(occurrence.OccurrenceType.ValueString())

		var occurrenceConfig fabcore.MonthlyOccurrenceClassification
		switch occurrenceType {
		case fabcore.OccurrenceTypeDayOfMonth:
			occurrenceConfig = &fabcore.DayOfMonth{
				DayOfMonth:     occurrence.DayOfMonth.ValueInt32Pointer(),
				OccurrenceType: &occurrenceType,
			}
		case fabcore.OccurrenceTypeOrdinalWeekday:
			weekday := fabcore.DayOfWeek(occurrence.Weekday.ValueString())
			weekIndex := fabcore.WeekIndex(occurrence.WeekIndex.ValueString())
			occurrenceConfig = &fabcore.OrdinalWeekday{
				Weekday:        &weekday,
				WeekIndex:      &weekIndex,
				OccurrenceType: &occurrenceType,
			}
		default:
			var diags diag.Diagnostics

			diags.AddError(
				"Unsupported Monthly Occurrence type",
				fmt.Sprintf("The Monthly Occurrence type '%T' is not supported.", occurrenceType),
			)

			return diags
		}

		reqConfiguration = &fabcore.MonthlyScheduleConfig{
			Times:           timesSlice,
			Recurrence:      configuration.Recurrence.ValueInt32Pointer(),
			Occurrence:      occurrenceConfig,
			StartDateTime:   &startDateTime,
			EndDateTime:     &endDateTime,
			LocalTimeZoneID: &localTimeZoneID,
			Type:            &configurationType,
		}

	default:
		diags.AddError(
			"Unsupported Configuration type",
			fmt.Sprintf("The Configuration type '%T' is not supported.", configurationType),
		)

		return diags
	}

	to.Configuration = reqConfiguration
	to.Enabled = from.Enabled.ValueBoolPointer()

	return nil
}

type principalModel struct {
	ID   customtypes.UUID `tfsdk:"id"`
	Type types.String     `tfsdk:"type"`
}

func (to *principalModel) set(from fabcore.Principal) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Type = types.StringPointerValue((*string)(from.Type))
}

var JobTypeActions = map[string][]string{ //nolint:gochecknoglobals
	"dataflow": {"Execute", "ApplyChanges"},
}
