// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package itemjob_test

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

var fakeItemJobInstanceStore = map[string]fabcore.ItemJobInstance{}

func GenerateItemJobInstanceID(workspaceID, jobInstanceID string) string {
	return fmt.Sprintf("%s/%s", workspaceID, jobInstanceID)
}

func fakeRunOnDemandItemJobFunc() func(ctx context.Context, workspaceID, itemID, jobType string, options *fabcore.JobSchedulerClientRunOnDemandItemJobOptions) (resp azfake.Responder[fabcore.JobSchedulerClientRunOnDemandItemJobResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID, jobType string, _ *fabcore.JobSchedulerClientRunOnDemandItemJobOptions) (resp azfake.Responder[fabcore.JobSchedulerClientRunOnDemandItemJobResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.JobSchedulerClientRunOnDemandItemJobResponse]{}
		jobInstanceID := testhelp.RandomUUID()

		itemJobInstance := fabcore.ItemJobInstance{
			ID:             new(jobInstanceID),
			ItemID:         new(itemID),
			JobType:        new(jobType),
			InvokeType:     to.Ptr(fabcore.InvokeTypeManual),
			Status:         to.Ptr(fabcore.ItemJobStatusNotStarted),
			RootActivityID: new(testhelp.RandomUUID()),
			StartTimeUTC:   new(time.Now().UTC().Format(time.RFC3339)),
		}

		fakeItemJobInstanceStore[GenerateItemJobInstanceID(workspaceID, jobInstanceID)] = itemJobInstance

		location := fmt.Sprintf("https://api.fabric.microsoft.com/v1/workspaces/%s/items/%s/jobs/instances/%s", workspaceID, itemID, jobInstanceID)

		resp.SetResponse(http.StatusAccepted, fabcore.JobSchedulerClientRunOnDemandItemJobResponse{}, &azfake.SetResponseOptions{
			Header: http.Header{"Location": []string{location}},
		})

		return resp, errResp
	}
}

func fakeGetItemJobInstanceFunc() func(ctx context.Context, workspaceID, itemID, jobInstanceID string, options *fabcore.JobSchedulerClientGetItemJobInstanceOptions) (resp azfake.Responder[fabcore.JobSchedulerClientGetItemJobInstanceResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, _, jobInstanceID string, _ *fabcore.JobSchedulerClientGetItemJobInstanceOptions) (resp azfake.Responder[fabcore.JobSchedulerClientGetItemJobInstanceResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.JobSchedulerClientGetItemJobInstanceResponse]{}
		errItemNotFound := fabcore.ErrItem.ItemNotFound.Error()
		id := GenerateItemJobInstanceID(workspaceID, jobInstanceID)

		if itemJobInstance, ok := fakeItemJobInstanceStore[id]; ok {
			resp.SetResponse(http.StatusOK, fabcore.JobSchedulerClientGetItemJobInstanceResponse{ItemJobInstance: itemJobInstance}, nil)
		} else {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, errItemNotFound, "Item not found"))
			resp.SetResponse(http.StatusNotFound, fabcore.JobSchedulerClientGetItemJobInstanceResponse{}, nil)
		}

		return resp, errResp
	}
}

func fakeCancelItemJobInstanceFunc() func(ctx context.Context, workspaceID, itemID, jobInstanceID string, options *fabcore.JobSchedulerClientCancelItemJobInstanceOptions) (resp azfake.Responder[fabcore.JobSchedulerClientCancelItemJobInstanceResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, _, jobInstanceID string, _ *fabcore.JobSchedulerClientCancelItemJobInstanceOptions) (resp azfake.Responder[fabcore.JobSchedulerClientCancelItemJobInstanceResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.JobSchedulerClientCancelItemJobInstanceResponse]{}

		delete(fakeItemJobInstanceStore, GenerateItemJobInstanceID(workspaceID, jobInstanceID))
		resp.SetResponse(http.StatusAccepted, fabcore.JobSchedulerClientCancelItemJobInstanceResponse{}, nil)

		return resp, errResp
	}
}

func NewRandomItemJobInstance(workspaceID, itemID, jobType string) fabcore.ItemJobInstance {
	jobInstanceID := testhelp.RandomUUID()

	itemJobInstance := fabcore.ItemJobInstance{
		ID:             new(jobInstanceID),
		ItemID:         new(itemID),
		JobType:        new(jobType),
		InvokeType:     to.Ptr(fabcore.InvokeTypeManual),
		Status:         to.Ptr(fabcore.ItemJobStatusCompleted),
		RootActivityID: new(testhelp.RandomUUID()),
		StartTimeUTC:   new(time.Now().UTC().Format(time.RFC3339)),
		EndTimeUTC:     new(time.Now().UTC().Add(time.Minute).Format(time.RFC3339)),
	}

	fakeItemJobInstanceStore[GenerateItemJobInstanceID(workspaceID, jobInstanceID)] = itemJobInstance

	return itemJobInstance
}
