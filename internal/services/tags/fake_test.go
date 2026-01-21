// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package tags_test

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var fakeTagStore = map[string]fabadmin.TagInfo{}

func fakeTagsFunc() func(options *fabadmin.TagsClientListTagsOptions) (resp azfake.PagerResponder[fabadmin.TagsClientListTagsResponse]) {
	return func(_ *fabadmin.TagsClientListTagsOptions) (resp azfake.PagerResponder[fabadmin.TagsClientListTagsResponse]) {
		resp = azfake.PagerResponder[fabadmin.TagsClientListTagsResponse]{}
		resp.AddPage(http.StatusOK, fabadmin.TagsClientListTagsResponse{TagsInfo: fabadmin.TagsInfo{Value: GetAllStoredTags()}}, nil)

		return resp
	}
}

func fakeDeleteTagFunc() func(ctx context.Context, tagID string, options *fabadmin.TagsClientDeleteTagOptions) (resp azfake.Responder[fabadmin.TagsClientDeleteTagResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, tagID string, _ *fabadmin.TagsClientDeleteTagOptions) (resp azfake.Responder[fabadmin.TagsClientDeleteTagResponse], errResp azfake.ErrorResponder) {
		if _, ok := fakeTagStore[tagID]; ok {
			delete(fakeTagStore, tagID)
			resp.SetResponse(http.StatusOK, struct{}{}, nil)
		} else {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, "ItemNotFound", "Item not found"))
			resp.SetResponse(http.StatusNotFound, struct{}{}, nil)
		}

		return resp, errResp
	}
}

func fakeUpdateTagFunc(
	tag fabadmin.TagInfo,
) func(ctx context.Context, tagID string, request fabadmin.UpdateTagRequest, options *fabadmin.TagsClientUpdateTagOptions) (resp azfake.Responder[fabadmin.TagsClientUpdateTagResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, tagID string, _ fabadmin.UpdateTagRequest, _ *fabadmin.TagsClientUpdateTagOptions) (resp azfake.Responder[fabadmin.TagsClientUpdateTagResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabadmin.TagsClientUpdateTagResponse]{}
		errItemNotFound := fabcore.ErrItem.ItemNotFound.Error()

		if _, ok := fakeTagStore[tagID]; !ok {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, errItemNotFound, "Item not found"))
			resp.SetResponse(http.StatusNotFound, fabadmin.TagsClientUpdateTagResponse{}, nil)

			return resp, errResp
		}

		tag := fabadmin.TagInfo{
			ID:          to.Ptr(tagID),
			DisplayName: to.Ptr(*tag.DisplayName),
			Scope: &fabadmin.TagScope{
				Type: to.Ptr(fabadmin.TagScopeTypeTenant),
			},
		}

		returnTag := fabadmin.Tag{
			ID:          to.Ptr(tagID),
			DisplayName: to.Ptr(*tag.DisplayName),
			Scope: &fabadmin.TagScope{
				Type: to.Ptr(fabadmin.TagScopeTypeTenant),
			},
		}

		fakeTestUpsert(tag)
		resp.SetResponse(http.StatusOK, fabadmin.TagsClientUpdateTagResponse{Tag: returnTag}, nil)

		return resp, errResp
	}
}

func fakeBulkCreateTagsFunc() func(_ context.Context, body fabadmin.CreateTagsRequest, options *fabadmin.TagsClientBulkCreateTagsOptions) (resp azfake.Responder[fabadmin.TagsClientBulkCreateTagsResponse], err azfake.ErrorResponder) {
	return func(_ context.Context, body fabadmin.CreateTagsRequest, _ *fabadmin.TagsClientBulkCreateTagsOptions) (resp azfake.Responder[fabadmin.TagsClientBulkCreateTagsResponse], err azfake.ErrorResponder) {
		resp = azfake.Responder[fabadmin.TagsClientBulkCreateTagsResponse]{}

		var outputTags []fabadmin.Tag

		for _, item := range body.CreateTagsRequest {
			newTag := NewRandomTag()
			newTag.DisplayName = item.DisplayName

			fakeTestUpsert(newTag)

			outputTags = append(outputTags, fabadmin.Tag(newTag))
		}

		resp.SetResponse(http.StatusCreated, fabadmin.TagsClientBulkCreateTagsResponse{CreateTagsResponse: fabadmin.CreateTagsResponse{Tags: outputTags}}, nil)

		return resp, err
	}
}

func NewRandomTag() fabadmin.TagInfo {
	return fabadmin.TagInfo{
		DisplayName: to.Ptr(testhelp.RandomName()),
		ID:          to.Ptr(testhelp.RandomUUID()),
		Scope: &fabadmin.TagScope{
			Type: to.Ptr(fabadmin.TagScopeTypeTenant),
		},
	}
}

func GetAllStoredTags() []fabadmin.TagInfo {
	tags := make([]fabadmin.TagInfo, 0, len(fakeTagStore))
	for _, tag := range fakeTagStore {
		tags = append(tags, tag)
	}

	return tags
}

func fakeTestUpsert(entity fabadmin.TagInfo) {
	fakeTagStore[*entity.ID] = entity
}
