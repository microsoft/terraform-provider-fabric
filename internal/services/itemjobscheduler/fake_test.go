// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package itemjobscheduler_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var fakeItemScheduleStore = map[string]fabcore.ItemSchedule{}

func fakeItemSchedulesFunc() func(workspaceID, itemID, jobType string, options *fabcore.JobSchedulerClientListItemSchedulesOptions) (resp azfake.PagerResponder[fabcore.JobSchedulerClientListItemSchedulesResponse]) {
	return func(_, _, _ string, _ *fabcore.JobSchedulerClientListItemSchedulesOptions) (resp azfake.PagerResponder[fabcore.JobSchedulerClientListItemSchedulesResponse]) {
		resp = azfake.PagerResponder[fabcore.JobSchedulerClientListItemSchedulesResponse]{}
		resp.AddPage(http.StatusOK, fabcore.JobSchedulerClientListItemSchedulesResponse{ItemSchedules: fabcore.ItemSchedules{Value: GetAllStoredItemSchedules()}}, nil)

		return resp
	}
}

func fakeGetItemScheduleFunc() func(ctx context.Context, workspaceID, itemID, jobType, scheduleID string, options *fabcore.JobSchedulerClientGetItemScheduleOptions) (resp azfake.Responder[fabcore.JobSchedulerClientGetItemScheduleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, _, _, scheduleID string, _ *fabcore.JobSchedulerClientGetItemScheduleOptions) (resp azfake.Responder[fabcore.JobSchedulerClientGetItemScheduleResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.JobSchedulerClientGetItemScheduleResponse]{}
		errItemNotFound := fabcore.ErrItem.ItemNotFound.Error()
		id := GenerateItemScheduleID(workspaceID, scheduleID)

		if itemSchedule, ok := fakeItemScheduleStore[id]; ok {
			resp.SetResponse(http.StatusOK, fabcore.JobSchedulerClientGetItemScheduleResponse{ItemSchedule: itemSchedule}, nil)
		} else {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, errItemNotFound, "Item not found"))
			resp.SetResponse(http.StatusNotFound, fabcore.JobSchedulerClientGetItemScheduleResponse{}, nil)
		}

		return resp, errResp
	}
}

func fakeGetFabricItem(
	itemType string,
) func(ctx context.Context, workspaceID, itemID string, options *fabcore.ItemsClientGetItemOptions) (resp azfake.Responder[fabcore.ItemsClientGetItemResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, _, _ string, _ *fabcore.ItemsClientGetItemOptions) (resp azfake.Responder[fabcore.ItemsClientGetItemResponse], errResp azfake.ErrorResponder) {
		item := fabcore.Item{
			Type: to.Ptr(fabcore.ItemType(itemType)),
		}
		resp.SetResponse(http.StatusOK, fabcore.ItemsClientGetItemResponse{Item: item}, nil)

		return resp, errResp
	}
}

func fakeCreateItemScheduleFunc() func(ctx context.Context, workspaceID, itemID, jobType string, createScheduleRequest fabcore.CreateScheduleRequest, options *fabcore.JobSchedulerClientCreateItemScheduleOptions) (resp azfake.Responder[fabcore.JobSchedulerClientCreateItemScheduleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, _, _ string, createScheduleRequest fabcore.CreateScheduleRequest, _ *fabcore.JobSchedulerClientCreateItemScheduleOptions) (resp azfake.Responder[fabcore.JobSchedulerClientCreateItemScheduleResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.JobSchedulerClientCreateItemScheduleResponse]{}
		itemScheduleID := testhelp.RandomUUID()

		itemSchedule := fabcore.ItemSchedule{
			ID:              to.Ptr(itemScheduleID),
			Enabled:         createScheduleRequest.Enabled,
			CreatedDateTime: to.Ptr(time.Now()),
			Configuration:   createScheduleRequest.Configuration,
			Owner: &fabcore.Principal{
				ID:   to.Ptr(testhelp.RandomUUID()),
				Type: to.Ptr(fabcore.PrincipalTypeUser),
			},
		}

		fakeTestUpsert(workspaceID, itemSchedule)
		resp.SetResponse(http.StatusCreated, fabcore.JobSchedulerClientCreateItemScheduleResponse{ItemSchedule: itemSchedule}, nil)

		return resp, errResp
	}
}

func fakeUpdateItemScheduleFunc() func(ctx context.Context, workspaceID, itemID, jobType, scheduleID string, updateScheduleRequest fabcore.UpdateScheduleRequest, options *fabcore.JobSchedulerClientUpdateItemScheduleOptions) (resp azfake.Responder[fabcore.JobSchedulerClientUpdateItemScheduleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, _, _, scheduleID string, updateScheduleRequest fabcore.UpdateScheduleRequest, _ *fabcore.JobSchedulerClientUpdateItemScheduleOptions) (resp azfake.Responder[fabcore.JobSchedulerClientUpdateItemScheduleResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.JobSchedulerClientUpdateItemScheduleResponse]{}
		id := GenerateItemScheduleID(workspaceID, scheduleID)
		errItemNotFound := fabcore.ErrItem.ItemNotFound.Error()

		if _, ok := fakeItemScheduleStore[id]; !ok {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, errItemNotFound, "Item not found"))
			resp.SetResponse(http.StatusNotFound, fabcore.JobSchedulerClientUpdateItemScheduleResponse{}, nil)

			return resp, errResp
		}

		itemSchedule := fabcore.ItemSchedule{
			ID:              to.Ptr(scheduleID),
			Enabled:         updateScheduleRequest.Enabled,
			CreatedDateTime: to.Ptr(time.Now()),
			Configuration:   updateScheduleRequest.Configuration,
			Owner: &fabcore.Principal{
				ID:   to.Ptr(testhelp.RandomUUID()),
				Type: to.Ptr(fabcore.PrincipalTypeUser),
			},
		}

		fakeTestUpsert(workspaceID, itemSchedule)
		resp.SetResponse(http.StatusOK, fabcore.JobSchedulerClientUpdateItemScheduleResponse{ItemSchedule: itemSchedule}, nil)

		return resp, errResp
	}
}

func fakeDeleteItemScheduleFunc() func(ctx context.Context, workspaceID, itemID, jobType, scheduleID string, options *fabcore.JobSchedulerClientDeleteItemScheduleOptions) (resp azfake.Responder[fabcore.JobSchedulerClientDeleteItemScheduleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, _, _, scheduleID string, _ *fabcore.JobSchedulerClientDeleteItemScheduleOptions) (resp azfake.Responder[fabcore.JobSchedulerClientDeleteItemScheduleResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.JobSchedulerClientDeleteItemScheduleResponse]{}
		id := GenerateItemScheduleID(workspaceID, scheduleID)
		errItemNotFound := fabcore.ErrItem.ItemNotFound.Error()

		if _, ok := fakeItemScheduleStore[id]; !ok {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, errItemNotFound, "Item not found"))
			resp.SetResponse(http.StatusNotFound, fabcore.JobSchedulerClientDeleteItemScheduleResponse{}, nil)

			return resp, errResp
		}

		delete(fakeItemScheduleStore, id)
		resp.SetResponse(http.StatusOK, fabcore.JobSchedulerClientDeleteItemScheduleResponse{}, nil)

		return resp, errResp
	}
}

func NewRandomItemSchedule(scheduleType fabcore.ScheduleType) fabcore.ItemSchedule {
	return fabcore.ItemSchedule{
		ID:              to.Ptr(testhelp.RandomUUID()),
		Enabled:         to.Ptr(testhelp.RandomBool()),
		CreatedDateTime: to.Ptr(time.Now()),
		Configuration:   NewRandomScheduleConfig(scheduleType),
		Owner:           NewRandomOwner(),
	}
}

func NewRandomScheduleConfig(scheduleType fabcore.ScheduleType) fabcore.ScheduleConfigClassification {
	switch scheduleType {
	case fabcore.ScheduleTypeDaily:
		return NewRandomDailyScheduleConfig()
	case fabcore.ScheduleTypeWeekly:
		return NewRandomWeeklyScheduleConfig()
	case fabcore.ScheduleTypeCron:
		return NewRandomCronScheduleConfig()
	case fabcore.ScheduleTypeMonthly:
		return NewRandomMonthlyScheduleConfig()
	default:
		panic("Unsupported Schedule type") // lintignore:R009
	}
}

func NewRandomCronScheduleConfig() *fabcore.CronScheduleConfig {
	return &fabcore.CronScheduleConfig{
		StartDateTime:   to.Ptr(time.Now().UTC()),
		EndDateTime:     to.Ptr(time.Now().UTC().Add(24 * time.Hour)),
		LocalTimeZoneID: to.Ptr(testhelp.RandomName()),
		Type:            to.Ptr(fabcore.ScheduleTypeCron),
		Interval:        to.Ptr(testhelp.RandomIntRange(int32(1), int32(60))),
	}
}

func NewRandomDailyScheduleConfig() *fabcore.DailyScheduleConfig {
	timeStr := fmt.Sprintf("%02d:%02d", time.Now().Hour(), time.Now().Minute())

	return &fabcore.DailyScheduleConfig{
		StartDateTime:   to.Ptr(time.Now().UTC()),
		EndDateTime:     to.Ptr(time.Now().UTC().Add(24 * time.Hour)),
		LocalTimeZoneID: to.Ptr(testhelp.RandomName()),
		Type:            to.Ptr(fabcore.ScheduleTypeDaily),
		Times: []string{
			timeStr,
		},
	}
}

func NewRandomWeeklyScheduleConfig() *fabcore.WeeklyScheduleConfig {
	timeStr := fmt.Sprintf("%02d:%02d", time.Now().Hour(), time.Now().Minute())

	return &fabcore.WeeklyScheduleConfig{
		StartDateTime:   to.Ptr(time.Now().UTC()),
		EndDateTime:     to.Ptr(time.Now().UTC().Add(24 * time.Hour)),
		LocalTimeZoneID: to.Ptr(testhelp.RandomName()),
		Type:            to.Ptr(fabcore.ScheduleTypeWeekly),
		Times: []string{
			timeStr,
		},
		Weekdays: []fabcore.DayOfWeek{
			fabcore.DayOfWeekMonday,
		},
	}
}

func NewRandomMonthlyScheduleConfig() *fabcore.MonthlyScheduleConfig {
	timeStr := fmt.Sprintf("%02d:%02d", time.Now().Hour(), time.Now().Minute())

	return &fabcore.MonthlyScheduleConfig{
		StartDateTime:   to.Ptr(time.Now().UTC()),
		EndDateTime:     to.Ptr(time.Now().UTC().Add(24 * time.Hour)),
		LocalTimeZoneID: to.Ptr(testhelp.RandomName()),
		Type:            to.Ptr(fabcore.ScheduleTypeMonthly),
		Times: []string{
			timeStr,
		},
		Recurrence: to.Ptr(testhelp.RandomIntRange(int32(1), int32(12))),
		Occurrence: &fabcore.DayOfMonth{
			OccurrenceType: to.Ptr(fabcore.OccurrenceTypeDayOfMonth),
			DayOfMonth:     to.Ptr(testhelp.RandomIntRange(int32(1), int32(31))),
		},
	}
}

func NewRandomOwner() *fabcore.Principal {
	return &fabcore.Principal{
		ID:   to.Ptr(testhelp.RandomUUID()),
		Type: to.Ptr(fabcore.PrincipalTypeServicePrincipalProfile),
	}
}

func GenerateItemScheduleID(workspaceID, entityID string) string {
	return fmt.Sprintf("%s/%s", workspaceID, entityID)
}

func GetAllStoredItemSchedules() []fabcore.ItemSchedule {
	itemSchedules := make([]fabcore.ItemSchedule, 0, len(fakeItemScheduleStore))
	for _, itemSchedule := range fakeItemScheduleStore {
		itemSchedules = append(itemSchedules, itemSchedule)
	}

	return itemSchedules
}

func fakeTestUpsert(workspaceID string, entity fabcore.ItemSchedule) {
	id := GenerateItemScheduleID(workspaceID, *entity.ID)
	fakeItemScheduleStore[id] = entity
}
