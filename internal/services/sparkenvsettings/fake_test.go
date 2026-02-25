// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sparkenvsettings_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var fakeSparkComputeStagingStore = map[string]fabenvironment.SparkCompute{}

func fakeGetStagingSparkComputeFunc() func(ctx context.Context, workspaceID, environmentID string, _ bool, _ *fabenvironment.StagingClientGetSparkComputeOptions) (resp azfake.Responder[fabenvironment.StagingClientGetSparkComputeResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, environmentID string, _ bool, _ *fabenvironment.StagingClientGetSparkComputeOptions) (resp azfake.Responder[fabenvironment.StagingClientGetSparkComputeResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabenvironment.StagingClientGetSparkComputeResponse]{}

		if sparkCompute, ok := fakeSparkComputeStagingStore[environmentID]; ok {
			resp.SetResponse(http.StatusOK, fabenvironment.StagingClientGetSparkComputeResponse{SparkCompute: sparkCompute}, nil)

			return resp, errResp
		}

		errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, fabcore.ErrCommon.EntityNotFound.Error(), "Entity not found"))
		resp.SetResponse(http.StatusNotFound, fabenvironment.StagingClientGetSparkComputeResponse{}, nil)

		return resp, errResp
	}
}

func fakeGetPublishedSparkComputeFunc() func(ctx context.Context, workspaceID, environmentID string, _ bool, _ *fabenvironment.PublishedClientGetSparkComputeOptions) (resp azfake.Responder[fabenvironment.PublishedClientGetSparkComputeResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, environmentID string, _ bool, _ *fabenvironment.PublishedClientGetSparkComputeOptions) (resp azfake.Responder[fabenvironment.PublishedClientGetSparkComputeResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabenvironment.PublishedClientGetSparkComputeResponse]{}

		if sparkCompute, ok := fakeSparkComputeStagingStore[environmentID]; ok {
			resp.SetResponse(http.StatusOK, fabenvironment.PublishedClientGetSparkComputeResponse{SparkCompute: sparkCompute}, nil)

			return resp, errResp
		}

		errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, fabcore.ErrCommon.EntityNotFound.Error(), "Entity not found"))
		resp.SetResponse(http.StatusNotFound, fabenvironment.PublishedClientGetSparkComputeResponse{}, nil)

		return resp, errResp
	}
}

func fakeUpdateStagingSparkComputeFunc() func(ctx context.Context, workspaceID, environmentID string, _ bool, req fabenvironment.UpdateEnvironmentSparkComputeRequest, _ *fabenvironment.StagingClientUpdateSparkComputeOptions) (resp azfake.Responder[fabenvironment.StagingClientUpdateSparkComputeResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, environmentID string, _ bool, req fabenvironment.UpdateEnvironmentSparkComputeRequest, _ *fabenvironment.StagingClientUpdateSparkComputeOptions) (resp azfake.Responder[fabenvironment.StagingClientUpdateSparkComputeResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabenvironment.StagingClientUpdateSparkComputeResponse]{}

		current, ok := fakeSparkComputeStagingStore[environmentID]
		if !ok {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, fabcore.ErrCommon.EntityNotFound.Error(), "Entity not found"))
			resp.SetResponse(http.StatusNotFound, fabenvironment.StagingClientUpdateSparkComputeResponse{}, nil)

			return resp, errResp
		}

		if req.DriverCores != nil {
			current.DriverCores = req.DriverCores
		}

		if req.DriverMemory != nil {
			current.DriverMemory = req.DriverMemory
		}

		if req.DynamicExecutorAllocation != nil {
			current.DynamicExecutorAllocation = req.DynamicExecutorAllocation
		}

		if req.ExecutorCores != nil {
			current.ExecutorCores = req.ExecutorCores
		}

		if req.ExecutorMemory != nil {
			current.ExecutorMemory = req.ExecutorMemory
		}

		if req.InstancePool != nil {
			current.InstancePool = req.InstancePool
		}

		if req.RuntimeVersion != nil {
			current.RuntimeVersion = req.RuntimeVersion
		}

		if req.SparkProperties != nil {
			current.SparkProperties = applySparkPropertiesUpdate(current.SparkProperties, req.SparkProperties)
		}

		fakeSparkComputeStagingStore[environmentID] = current
		resp.SetResponse(http.StatusOK, fabenvironment.StagingClientUpdateSparkComputeResponse{SparkCompute: current}, nil)

		return resp, errResp
	}
}

func fakeBeginPublishEnvironmentFunc() func(ctx context.Context, workspaceID, environmentID string, _ bool, _ *fabenvironment.ItemsClientBeginPublishEnvironmentOptions) (resp azfake.PollerResponder[fabenvironment.ItemsClientPublishEnvironmentResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, environmentID string, _ bool, _ *fabenvironment.ItemsClientBeginPublishEnvironmentOptions) (resp azfake.PollerResponder[fabenvironment.ItemsClientPublishEnvironmentResponse], errResp azfake.ErrorResponder) {
		resp = azfake.PollerResponder[fabenvironment.ItemsClientPublishEnvironmentResponse]{}

		if _, ok := fakeSparkComputeStagingStore[environmentID]; ok {
			resp.SetTerminalResponse(http.StatusOK, fabenvironment.ItemsClientPublishEnvironmentResponse{
				Properties: fabenvironment.Properties{
					PublishDetails: &fabenvironment.PublishDetails{
						State: to.Ptr(fabenvironment.PublishStateSuccess),
					},
				},
			}, nil)

			return resp, errResp
		}

		errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, fabcore.ErrCommon.EntityNotFound.Error(), "Entity not found"))
		resp.SetTerminalResponse(http.StatusNotFound, fabenvironment.ItemsClientPublishEnvironmentResponse{}, nil)

		return resp, errResp
	}
}

func NewRandomSparkCompute() fabenvironment.SparkCompute {
	return fabenvironment.SparkCompute{
		DriverCores:    to.Ptr(int32(4)),
		DriverMemory:   to.Ptr(fabenvironment.CustomPoolMemoryTwentyEightG),
		ExecutorCores:  to.Ptr(int32(4)),
		ExecutorMemory: to.Ptr(fabenvironment.CustomPoolMemoryTwentyEightG),
		RuntimeVersion: to.Ptr("1.3"),
		InstancePool: &fabenvironment.InstancePool{
			Name: to.Ptr("Starter Pool"),
			Type: to.Ptr(fabenvironment.CustomPoolTypeWorkspace),
			ID:   to.Ptr(testhelp.RandomUUID()),
		},
	}
}

func fakeTestUpsertSparkComputeStaging(environmentID string, sparkCompute fabenvironment.SparkCompute) {
	fakeSparkComputeStagingStore[environmentID] = sparkCompute
}

func applySparkPropertiesUpdate(current, updates []fabenvironment.SparkProperty) []fabenvironment.SparkProperty {
	result := make([]fabenvironment.SparkProperty, 0)

	// Keep current properties that are not being updated or deleted.
	for _, c := range current {
		if c.Key != nil && !containsSparkPropertyKey(updates, *c.Key) {
			result = append(result, c)
		}
	}

	// Add non-nil update entries
	for _, u := range updates {
		if u.Key != nil && u.Value != nil {
			result = append(result, u)
		}
	}

	return result
}

func containsSparkPropertyKey(properties []fabenvironment.SparkProperty, key string) bool {
	for _, p := range properties {
		if p.Key != nil && *p.Key == key {
			return true
		}
	}

	return false
}
