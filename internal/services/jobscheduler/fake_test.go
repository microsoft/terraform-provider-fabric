package jobscheduler_test

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

		return
	}
}

func fakeGetItemScheduleFunc() func(ctx context.Context, workspaceID, itemID, jobType, scheduleID string, options *fabcore.JobSchedulerClientGetItemScheduleOptions) (resp azfake.Responder[fabcore.JobSchedulerClientGetItemScheduleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID, jobType, scheduleID string, _ *fabcore.JobSchedulerClientGetItemScheduleOptions) (resp azfake.Responder[fabcore.JobSchedulerClientGetItemScheduleResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.JobSchedulerClientGetItemScheduleResponse]{}
		errItemNotFound := fabcore.ErrItem.ItemNotFound.Error()
		id := GenerateItemScheduleID(workspaceID, scheduleID)

		if itemSchedule, ok := fakeItemScheduleStore[id]; ok {
			resp.SetResponse(http.StatusOK, fabcore.JobSchedulerClientGetItemScheduleResponse{ItemSchedule: itemSchedule}, nil)
		} else {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, errItemNotFound, "Item not found"))
			resp.SetResponse(http.StatusNotFound, fabcore.JobSchedulerClientGetItemScheduleResponse{}, nil)
		}

		return
	}
}

func fakeGetFabricItem(
	itemType string,
) func(ctx context.Context, workspaceID, itemID string, options *fabcore.ItemsClientGetItemOptions) (resp azfake.Responder[fabcore.ItemsClientGetItemResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID string, _ *fabcore.ItemsClientGetItemOptions) (resp azfake.Responder[fabcore.ItemsClientGetItemResponse], errResp azfake.ErrorResponder) {
		item := fabcore.Item{
			Type: to.Ptr(fabcore.ItemType(itemType)),
		}
		resp.SetResponse(http.StatusOK, fabcore.ItemsClientGetItemResponse{Item: item}, nil)

		return
	}
}

func fakeCreateItemScheduleFunc() func(ctx context.Context, workspaceID, itemID, jobType string, createScheduleRequest fabcore.CreateScheduleRequest, options *fabcore.JobSchedulerClientCreateItemScheduleOptions) (resp azfake.Responder[fabcore.JobSchedulerClientCreateItemScheduleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID, jobType string, createScheduleRequest fabcore.CreateScheduleRequest, _ *fabcore.JobSchedulerClientCreateItemScheduleOptions) (resp azfake.Responder[fabcore.JobSchedulerClientCreateItemScheduleResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.JobSchedulerClientCreateItemScheduleResponse]{}
		itemScheduleID := testhelp.RandomUUID()

		itemSchedule := fabcore.ItemSchedule{
			ID:              to.Ptr(itemScheduleID),
			Enabled:         createScheduleRequest.Enabled,
			CreatedDateTime: to.Ptr(time.Now()),
			Configuration:   createScheduleRequest.Configuration,
			Owner: &fabcore.Principal{
				ID:          to.Ptr(testhelp.RandomUUID()),
				DisplayName: to.Ptr(testhelp.RandomName()),
				Type:        to.Ptr(fabcore.PrincipalTypeUser),
				UserDetails: &fabcore.PrincipalUserDetails{
					UserPrincipalName: to.Ptr(testhelp.RandomName()),
				},
			},
		}

		fakeTestUpsert(workspaceID, itemSchedule)
		resp.SetResponse(http.StatusCreated, fabcore.JobSchedulerClientCreateItemScheduleResponse{ItemSchedule: itemSchedule}, nil)

		return
	}
}

func fakeUpdateItemScheduleFunc() func(ctx context.Context, workspaceID, itemID, jobType, scheduleID string, updateScheduleRequest fabcore.UpdateScheduleRequest, options *fabcore.JobSchedulerClientUpdateItemScheduleOptions) (resp azfake.Responder[fabcore.JobSchedulerClientUpdateItemScheduleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID, jobType, scheduleID string, updateScheduleRequest fabcore.UpdateScheduleRequest, _ *fabcore.JobSchedulerClientUpdateItemScheduleOptions) (resp azfake.Responder[fabcore.JobSchedulerClientUpdateItemScheduleResponse], errResp azfake.ErrorResponder) {
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
				ID:          to.Ptr(testhelp.RandomUUID()),
				DisplayName: to.Ptr(testhelp.RandomName()),
				Type:        to.Ptr(fabcore.PrincipalTypeUser),
				UserDetails: &fabcore.PrincipalUserDetails{
					UserPrincipalName: to.Ptr(testhelp.RandomName()),
				},
			},
		}

		fakeTestUpsert(workspaceID, itemSchedule)
		resp.SetResponse(http.StatusOK, fabcore.JobSchedulerClientUpdateItemScheduleResponse{ItemSchedule: itemSchedule}, nil)

		return resp, errResp
	}
}

func fakeDeleteItemScheduleFunc() func(ctx context.Context, workspaceID, itemID, jobType, scheduleID string, options *fabcore.JobSchedulerClientDeleteItemScheduleOptions) (resp azfake.Responder[fabcore.JobSchedulerClientDeleteItemScheduleResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID, jobType, scheduleID string, _ *fabcore.JobSchedulerClientDeleteItemScheduleOptions) (resp azfake.Responder[fabcore.JobSchedulerClientDeleteItemScheduleResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.JobSchedulerClientDeleteItemScheduleResponse]{}
		id := GenerateItemScheduleID(workspaceID, scheduleID)
		errItemNotFound := fabcore.ErrItem.ItemNotFound.Error()

		if _, ok := fakeItemScheduleStore[id]; !ok {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, errItemNotFound, "Item not found"))
			resp.SetResponse(http.StatusNotFound, fabcore.JobSchedulerClientDeleteItemScheduleResponse{}, nil)

			return
		}

		delete(fakeItemScheduleStore, id)
		resp.SetResponse(http.StatusOK, fabcore.JobSchedulerClientDeleteItemScheduleResponse{}, nil)

		return
	}
}

func NewRandomItemSchedule(scheduleType fabcore.ScheduleType) fabcore.ItemSchedule {
	return fabcore.ItemSchedule{
		ID:              to.Ptr(testhelp.RandomUUID()),
		Enabled:         to.Ptr(testhelp.RandomBool()),
		CreatedDateTime: to.Ptr(time.Now()),
		Configuration:   NewRandomScheduleConfig(scheduleType),
		Owner:           NewRadomOwner(),
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
	default:
		panic("Unsupported Schedule type") // lintignore:R009
	}
}

func NewRandomCronScheduleConfig() *fabcore.CronScheduleConfig {
	return &fabcore.CronScheduleConfig{
		StartDateTime:   to.Ptr(time.Now()),
		EndDateTime:     to.Ptr(time.Now()),
		LocalTimeZoneID: to.Ptr(testhelp.RandomName()),
		Type:            to.Ptr(fabcore.ScheduleTypeCron),
		Interval:        to.Ptr(int32(testhelp.RandomIntRange(1, 60))),
	}
}

func NewRandomDailyScheduleConfig() *fabcore.DailyScheduleConfig {
	return &fabcore.DailyScheduleConfig{
		StartDateTime:   to.Ptr(time.Now()),
		EndDateTime:     to.Ptr(time.Now()),
		LocalTimeZoneID: to.Ptr(testhelp.RandomName()),
		Type:            to.Ptr(fabcore.ScheduleTypeDaily),
		Times: []string{
			time.Now().String(),
		},
	}
}

func NewRandomWeeklyScheduleConfig() *fabcore.WeeklyScheduleConfig {
	return &fabcore.WeeklyScheduleConfig{
		StartDateTime:   to.Ptr(time.Now()),
		EndDateTime:     to.Ptr(time.Now()),
		LocalTimeZoneID: to.Ptr(testhelp.RandomName()),
		Type:            to.Ptr(fabcore.ScheduleTypeWeekly),
		Times: []string{
			"09:00:00",
		},
		Weekdays: []fabcore.DayOfWeek{
			fabcore.DayOfWeekMonday,
		},
	}
}

func NewRadomOwner() *fabcore.Principal {
	return &fabcore.Principal{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Type:        to.Ptr(fabcore.PrincipalTypeUser),
		UserDetails: &fabcore.PrincipalUserDetails{
			UserPrincipalName: to.Ptr(testhelp.RandomName()),
		},
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
