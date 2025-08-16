package jobscheduler

import (
	"context"
	"fmt"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	//revive:disable-line:import-alias-naming
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

type baseJobScheduleModel struct {
	ID              customtypes.UUID                                                        `tfsdk:"id"`
	Enabled         types.Bool                                                              `tfsdk:"enabled"`
	ItemID          customtypes.UUID                                                        `tfsdk:"item_id"`
	WorkspaceID     customtypes.UUID                                                        `tfsdk:"workspace_id"`
	JobType         types.String                                                            `tfsdk:"job_type"`
	CreatedDateTime timetypes.RFC3339                                                       `tfsdk:"created_date_time"`
	Owner           supertypes.SingleNestedObjectValueOf[baseOwnerModel]                    `tfsdk:"owner"`
	Configuration   supertypes.SingleNestedObjectValueOf[baseJobScheduleConfigurationModel] `tfsdk:"configuration"`
}

type baseJobScheduleConfigurationModel struct {
	StartDateTime   timetypes.RFC3339                   `tfsdk:"start_date_time"`
	EndDateTime     timetypes.RFC3339                   `tfsdk:"end_date_time"`
	LocalTimeZoneId types.String                        `tfsdk:"local_time_zone"`
	Type            types.String                        `tfsdk:"type"`
	Interval        types.Int32                         `tfsdk:"interval"`
	Times           supertypes.SetValueOf[types.String] `tfsdk:"times"`
	Weekdays        supertypes.SetValueOf[types.String] `tfsdk:"weekdays"`
}

type baseOwnerModel struct {
	ID                             customtypes.UUID                                                          `tfsdk:"id"`
	DisplayName                    types.String                                                              `tfsdk:"display_name"`
	GroupDetails                   supertypes.SingleNestedObjectValueOf[groupDetailsModel]                   `tfsdk:"group_details"`
	ServicePrincipalDetails        supertypes.SingleNestedObjectValueOf[servicePrincipalDetailsModel]        `tfsdk:"service_principal_details"`
	ServicePrincipalProfileDetails supertypes.SingleNestedObjectValueOf[servicePrincipalProfileDetailsModel] `tfsdk:"service_principal_profile_details"`
	Type                           types.String                                                              `tfsdk:"type"`
	UserDetails                    supertypes.SingleNestedObjectValueOf[userDetailsModel]                    `tfsdk:"user_details"`
}

type groupDetailsModel struct {
	GroupType types.String `tfsdk:"group_type"`
}

type servicePrincipalDetailsModel struct {
	AadAppId types.String `tfsdk:"aad_app_id"`
}

type servicePrincipalProfileDetailsModel struct {
	ParentPrincipal supertypes.SingleNestedObjectValueOf[servicePrincipalBaseOwnerModel] `tfsdk:"parent_principal"`
}

type servicePrincipalBaseOwnerModel struct {
	ID                      customtypes.UUID                                                   `tfsdk:"id"`
	DisplayName             types.String                                                       `tfsdk:"display_name"`
	ServicePrincipalDetails supertypes.SingleNestedObjectValueOf[servicePrincipalDetailsModel] `tfsdk:"service_principal_details"`
	Type                    types.String                                                       `tfsdk:"type"`
}

type userDetailsModel struct {
	UserPrincipalName types.String `tfsdk:"user_principal_name"`
}

func (to *baseJobScheduleModel) set(ctx context.Context, workspaceID, itemID, jobType string, from fabcore.ItemSchedule) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Enabled = types.BoolPointerValue(from.Enabled)
	to.ItemID = customtypes.NewUUIDValue(itemID)
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.JobType = types.StringValue(jobType)
	to.CreatedDateTime = timetypes.NewRFC3339TimePointerValue(from.CreatedDateTime)
	to.Configuration = supertypes.NewSingleNestedObjectValueOfNull[baseJobScheduleConfigurationModel](ctx)
	configuration := supertypes.NewSingleNestedObjectValueOfNull[baseJobScheduleConfigurationModel](ctx)
	owner := supertypes.NewSingleNestedObjectValueOfNull[baseOwnerModel](ctx)
	to.Owner = owner

	if from.Configuration != nil {
		baseJobScheduleConfigurationModel := &baseJobScheduleConfigurationModel{}
		if diags := baseJobScheduleConfigurationModel.set(ctx, from.Configuration); diags.HasError() {
			return diags
		}

		if diags := configuration.Set(ctx, baseJobScheduleConfigurationModel); diags.HasError() {
			return diags
		}
	}

	to.Configuration = configuration

	to.Owner = supertypes.NewSingleNestedObjectValueOfNull[baseOwnerModel](ctx)

	if from.Owner != nil {
		ownerModel := &baseOwnerModel{}
		if diags := ownerModel.set(ctx, from.Owner); diags.HasError() {
			return diags
		}

		if diags := owner.Set(ctx, ownerModel); diags.HasError() {
			return diags
		}
	}

	to.Owner = owner

	return nil
}

func (to *baseOwnerModel) set(ctx context.Context, from *fabcore.Principal) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.Type = types.StringPointerValue((*string)(from.Type))

	to.GroupDetails = supertypes.NewSingleNestedObjectValueOfNull[groupDetailsModel](ctx)
	to.ServicePrincipalDetails = supertypes.NewSingleNestedObjectValueOfNull[servicePrincipalDetailsModel](ctx)
	to.ServicePrincipalProfileDetails = supertypes.NewSingleNestedObjectValueOfNull[servicePrincipalProfileDetailsModel](ctx)
	to.UserDetails = supertypes.NewSingleNestedObjectValueOfNull[userDetailsModel](ctx)

	if from.Type == nil {
		return nil
	}

	switch *from.Type {
	case fabcore.PrincipalTypeGroup:
		return setGroupDetails(ctx, to, from)
	case fabcore.PrincipalTypeUser:
		return setUserDetails(ctx, to, from)
	case fabcore.PrincipalTypeServicePrincipal:
		return setServicePrincipalDetails(ctx, to, from)
	case fabcore.PrincipalTypeServicePrincipalProfile:
		return setServicePrincipalProfileDetails(ctx, to, from)
	default:
		var diags diag.Diagnostics

		diags.AddError(
			"Unsupported Owner type",
			"The owner type is not supported.",
		)

		return diags
	}
}

func setGroupDetails(ctx context.Context, to *baseOwnerModel, from *fabcore.Principal) diag.Diagnostics {
	if from.GroupDetails == nil {
		return nil
	}

	groupDetails := supertypes.NewSingleNestedObjectValueOfNull[groupDetailsModel](ctx)

	groupDetailsModel := &groupDetailsModel{}
	if diags := groupDetailsModel.set(from.GroupDetails); diags.HasError() {
		return diags
	}

	if diags := groupDetails.Set(ctx, groupDetailsModel); diags.HasError() {
		return diags
	}

	to.GroupDetails = groupDetails

	return nil
}

func setUserDetails(ctx context.Context, to *baseOwnerModel, from *fabcore.Principal) diag.Diagnostics {
	if from.UserDetails == nil {
		return nil
	}

	userDetails := supertypes.NewSingleNestedObjectValueOfNull[userDetailsModel](ctx)

	userDetailsModel := &userDetailsModel{}
	if diags := userDetailsModel.set(from.UserDetails); diags.HasError() {
		return diags
	}

	if diags := userDetails.Set(ctx, userDetailsModel); diags.HasError() {
		return diags
	}

	to.UserDetails = userDetails

	return nil
}

func setServicePrincipalDetails(ctx context.Context, to *baseOwnerModel, from *fabcore.Principal) diag.Diagnostics {
	if from.ServicePrincipalDetails == nil {
		return nil
	}

	servicePrincipalDetails := supertypes.NewSingleNestedObjectValueOfNull[servicePrincipalDetailsModel](ctx)

	spDetailsModel := &servicePrincipalDetailsModel{}
	if diags := spDetailsModel.set(from.ServicePrincipalDetails); diags.HasError() {
		return diags
	}

	if diags := servicePrincipalDetails.Set(ctx, spDetailsModel); diags.HasError() {
		return diags
	}

	to.ServicePrincipalDetails = servicePrincipalDetails

	return nil
}

func setServicePrincipalProfileDetails(ctx context.Context, to *baseOwnerModel, from *fabcore.Principal) diag.Diagnostics {
	if from.ServicePrincipalProfileDetails == nil {
		return nil
	}

	servicePrincipalProfileDetails := supertypes.NewSingleNestedObjectValueOfNull[servicePrincipalProfileDetailsModel](ctx)

	spProfileDetailsModel := &servicePrincipalProfileDetailsModel{}
	if diags := spProfileDetailsModel.set(ctx, from.ServicePrincipalProfileDetails); diags.HasError() {
		return diags
	}

	if diags := servicePrincipalProfileDetails.Set(ctx, spProfileDetailsModel); diags.HasError() {
		return diags
	}

	to.ServicePrincipalProfileDetails = servicePrincipalProfileDetails

	return nil
}

func (to *groupDetailsModel) set(from *fabcore.PrincipalGroupDetails) diag.Diagnostics {
	to.GroupType = types.StringPointerValue((*string)(from.GroupType))

	return nil
}

func (to *userDetailsModel) set(from *fabcore.PrincipalUserDetails) diag.Diagnostics {
	to.UserPrincipalName = types.StringPointerValue(from.UserPrincipalName)

	return nil
}

func (to *servicePrincipalDetailsModel) set(from *fabcore.PrincipalServicePrincipalDetails) diag.Diagnostics {
	to.AadAppId = types.StringPointerValue(from.AADAppID)

	return nil
}

func (to *servicePrincipalProfileDetailsModel) set(ctx context.Context, from *fabcore.PrincipalServicePrincipalProfileDetails) diag.Diagnostics {
	to.ParentPrincipal = supertypes.NewSingleNestedObjectValueOfNull[servicePrincipalBaseOwnerModel](ctx)
	owner := supertypes.NewSingleNestedObjectValueOfNull[servicePrincipalBaseOwnerModel](ctx)

	if from.ParentPrincipal != nil {
		ownerModel := &servicePrincipalBaseOwnerModel{}
		if diags := ownerModel.set(ctx, from.ParentPrincipal); diags.HasError() {
			return diags
		}

		if diags := owner.Set(ctx, ownerModel); diags.HasError() {
			return diags
		}
	}

	to.ParentPrincipal = owner

	return nil
}

func (to *servicePrincipalBaseOwnerModel) set(ctx context.Context, from *fabcore.Principal) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.DisplayName = types.StringPointerValue(from.DisplayName)
	to.ServicePrincipalDetails = supertypes.NewSingleNestedObjectValueOfNull[servicePrincipalDetailsModel](ctx)

	if from.ServicePrincipalDetails != nil {
		spDetailsModel := &servicePrincipalDetailsModel{}
		if diags := spDetailsModel.set(from.ServicePrincipalDetails); diags.HasError() {
			return diags
		}

		if diags := to.ServicePrincipalDetails.Set(ctx, spDetailsModel); diags.HasError() {
			return diags
		}
	}

	to.Type = types.StringPointerValue((*string)(from.Type))

	return nil
}

func (to *baseJobScheduleConfigurationModel) set(ctx context.Context, from fabcore.ScheduleConfigClassification) diag.Diagnostics {
	schConfig := from.GetScheduleConfig()
	to.StartDateTime = timetypes.NewRFC3339TimePointerValue(schConfig.StartDateTime)
	to.EndDateTime = timetypes.NewRFC3339TimePointerValue(schConfig.EndDateTime)
	to.LocalTimeZoneId = types.StringPointerValue(schConfig.LocalTimeZoneID)
	to.Type = types.StringPointerValue((*string)(schConfig.Type))
	to.Times = supertypes.NewSetValueOfNull[types.String](ctx)
	to.Weekdays = supertypes.NewSetValueOfNull[types.String](ctx)

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
		timesPtr := make([]*types.String, len(entity.Times))

		for i, t := range entity.Times {
			val := types.StringValue(t)
			timesPtr[i] = &val
		}

		times := make([]types.String, len(timesPtr))

		for i, t := range timesPtr {
			if t != nil {
				times[i] = *t
			} else {
				times[i] = types.StringNull()
			}
		}

		to.Times.Set(ctx, times)

		weekdaysPtr := make([]*types.String, len(entity.Weekdays))

		for i, w := range entity.Weekdays {
			val := types.StringValue(string(w))
			weekdaysPtr[i] = &val
		}

		weekdays := make([]types.String, len(weekdaysPtr))

		for i, w := range weekdaysPtr {
			if w != nil {
				weekdays[i] = *w
			} else {
				weekdays[i] = types.StringNull()
			}
		}

		to.Weekdays.Set(ctx, weekdays)
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

/*
DATA-SOURCE
*/

type dataSourceJobScheduleModel struct {
	baseJobScheduleModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
DATA-SOURCE (list)
*/

type dataSourceJobSchedulesModel struct {
	WorkspaceID customtypes.UUID                                        `tfsdk:"workspace_id"`
	ItemID      customtypes.UUID                                        `tfsdk:"item_id"`
	JobType     types.String                                            `tfsdk:"job_type"`
	Values      supertypes.SetNestedObjectValueOf[baseJobScheduleModel] `tfsdk:"values"`
	Timeouts    timeoutsD.Value                                         `tfsdk:"timeouts"`
}

func (to *dataSourceJobSchedulesModel) setValues(ctx context.Context, workspaceID, itemID, jobType string, from []fabcore.ItemSchedule) diag.Diagnostics {
	slice := make([]*baseJobScheduleModel, 0, len(from))

	for _, entity := range from {
		var entityModel baseJobScheduleModel

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
	baseJobScheduleModel

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
			LocalTimeZoneID: configuration.LocalTimeZoneId.ValueStringPointer(),
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
			LocalTimeZoneID: configuration.LocalTimeZoneId.ValueStringPointer(),
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
			LocalTimeZoneID: configuration.LocalTimeZoneId.ValueStringPointer(),
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
			LocalTimeZoneID: configuration.LocalTimeZoneId.ValueStringPointer(),
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
			LocalTimeZoneID: configuration.LocalTimeZoneId.ValueStringPointer(),
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
			LocalTimeZoneID: configuration.LocalTimeZoneId.ValueStringPointer(),
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

var JobTypeActions = map[string][]string{
	"dataflow": {"Execute", "ApplyChanges"},
}
