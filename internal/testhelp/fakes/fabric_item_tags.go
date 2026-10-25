// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"context"
	"net/http"
	"reflect"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
)

func configureItemTags(server *fakeServer) {
	server.ServerFactory.Core.TagsServer.ApplyTags = fakeApplyTags(server)
	server.ServerFactory.Core.TagsServer.UnapplyTags = fakeUnapplyTags(server)
}

func fakeApplyTags(
	server *fakeServer,
) func(ctx context.Context, workspaceID, itemID string, applyTagsRequest fabcore.ApplyTagsRequest, options *fabcore.TagsClientApplyTagsOptions) (resp azfake.Responder[fabcore.TagsClientApplyTagsResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID string, req fabcore.ApplyTagsRequest, _ *fabcore.TagsClientApplyTagsOptions) (azfake.Responder[fabcore.TagsClientApplyTagsResponse], azfake.ErrorResponder) {
		var resp azfake.Responder[fabcore.TagsClientApplyTagsResponse]
		var errResp azfake.ErrorResponder

		idx, item := findItemByID(server, workspaceID, itemID)
		if item == nil {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, fabcore.ErrItem.ItemNotFound.Error(), fabcore.ErrItem.ItemNotFound.Error()))

			return resp, errResp
		}

		// Add tags to the item
		for _, tagID := range req.Tags {
			displayName := "Tag " + tagID

			item.Tags = append(item.Tags, fabcore.ItemTag{
				ID:          &tagID,
				DisplayName: &displayName,
			})
		}

		updateItemTags(server, idx, *item)

		resp.SetResponse(http.StatusOK, fabcore.TagsClientApplyTagsResponse{}, nil)

		return resp, errResp
	}
}

func fakeUnapplyTags(
	server *fakeServer,
) func(ctx context.Context, workspaceID, itemID string, unapplyTagsRequest fabcore.UnapplyTagsRequest, options *fabcore.TagsClientUnapplyTagsOptions) (resp azfake.Responder[fabcore.TagsClientUnapplyTagsResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID string, req fabcore.UnapplyTagsRequest, _ *fabcore.TagsClientUnapplyTagsOptions) (azfake.Responder[fabcore.TagsClientUnapplyTagsResponse], azfake.ErrorResponder) {
		var resp azfake.Responder[fabcore.TagsClientUnapplyTagsResponse]
		var errResp azfake.ErrorResponder

		idx, item := findItemByID(server, workspaceID, itemID)
		if item == nil {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, fabcore.ErrItem.ItemNotFound.Error(), fabcore.ErrItem.ItemNotFound.Error()))

			return resp, errResp
		}

		// Remove specified tags from the item
		toRemove := make(map[string]bool, len(req.Tags))
		for _, tagID := range req.Tags {
			toRemove[tagID] = true
		}

		remaining := make([]fabcore.ItemTag, 0, len(item.Tags))

		for _, tag := range item.Tags {
			if tag.ID != nil && !toRemove[*tag.ID] {
				remaining = append(remaining, tag)
			}
		}

		item.Tags = remaining

		updateItemTags(server, idx, *item)

		resp.SetResponse(http.StatusOK, fabcore.TagsClientUnapplyTagsResponse{}, nil)

		return resp, errResp
	}
}

// findItemByID finds an item in the fake server's elements by workspaceID and itemID.
// Returns the index and a fabcore.Item representation. Uses asFabricItem for non-core types.
func findItemByID(server *fakeServer, workspaceID, itemID string) (int, *fabcore.Item) {
	for i, element := range server.elements {
		// Try direct fabcore.Item first
		if item, ok := element.(fabcore.Item); ok {
			if item.ID != nil && *item.ID == itemID && item.WorkspaceID != nil && *item.WorkspaceID == workspaceID {
				return i, &item
			}

			continue
		}

		// Use asFabricItem for other entity types (e.g., fabwarehouse.Warehouse)
		item := asFabricItem(element)
		if item.ID != nil && *item.ID == itemID && item.WorkspaceID != nil && *item.WorkspaceID == workspaceID {
			return i, &item
		}
	}

	return -1, nil
}

// updateItemTags updates the Tags field on the element at the given index.
func updateItemTags(server *fakeServer, idx int, updated fabcore.Item) {
	element := server.elements[idx]

	// For direct fabcore.Item, replace entirely
	if _, ok := element.(fabcore.Item); ok {
		server.elements[idx] = updated

		return
	}

	// For other types, create an addressable copy and use setReflectedTagsPropertyValue
	v := reflect.ValueOf(element)
	newVal := reflect.New(v.Type())
	newVal.Elem().Set(v)
	setReflectedTagsPropertyValue(newVal.Interface(), "Tags", updated.Tags)
	server.elements[idx] = newVal.Elem().Interface()
}
