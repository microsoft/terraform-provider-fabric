// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package tags_test

import (
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"

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
