// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package lakehouse_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
)

// fakeApplyTagsFunc returns a fake ApplyTags handler that succeeds silently.
// In unit tests, tags operations are verified by checking the test state assertions,
// not by inspecting the fake tag store.
func fakeApplyTagsFunc() func(ctx context.Context, workspaceID, itemID string, request fabcore.ApplyTagsRequest, options *fabcore.TagsClientApplyTagsOptions) (resp azfake.Responder[fabcore.TagsClientApplyTagsResponse], errResp azfake.ErrorResponder) {
	return func(ctx context.Context, workspaceID, itemID string, request fabcore.ApplyTagsRequest, _ *fabcore.TagsClientApplyTagsOptions) (resp azfake.Responder[fabcore.TagsClientApplyTagsResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.TagsClientApplyTagsResponse]{}
		errResp = azfake.ErrorResponder{}
		resp.SetResponse(http.StatusOK, fabcore.TagsClientApplyTagsResponse{}, nil)
		return resp, errResp
	}
}

// fakeUnapplyTagsFunc returns a fake UnapplyTags handler that succeeds silently.
// In unit tests, tags operations are verified by checking the test state assertions,
// not by inspecting the fake tag store.
func fakeUnapplyTagsFunc() func(ctx context.Context, workspaceID, itemID string, request fabcore.UnapplyTagsRequest, options *fabcore.TagsClientUnapplyTagsOptions) (resp azfake.Responder[fabcore.TagsClientUnapplyTagsResponse], errResp azfake.ErrorResponder) {
	return func(ctx context.Context, workspaceID, itemID string, request fabcore.UnapplyTagsRequest, _ *fabcore.TagsClientUnapplyTagsOptions) (resp azfake.Responder[fabcore.TagsClientUnapplyTagsResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.TagsClientUnapplyTagsResponse]{}
		errResp = azfake.ErrorResponder{}
		resp.SetResponse(http.StatusOK, fabcore.TagsClientUnapplyTagsResponse{}, nil)
		return resp, errResp
	}
}
