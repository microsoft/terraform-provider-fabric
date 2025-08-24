// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package shortcut_test

import (
	"context"
	"fmt"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var fakeShortcutStore = map[string]fabcore.Shortcut{}

func fakeShortcutsFunc() func(workspaceID, itemID string, options *fabcore.OneLakeShortcutsClientListShortcutsOptions) (resp azfake.PagerResponder[fabcore.OneLakeShortcutsClientListShortcutsResponse]) {
	return func(_, _ string, _ *fabcore.OneLakeShortcutsClientListShortcutsOptions) (resp azfake.PagerResponder[fabcore.OneLakeShortcutsClientListShortcutsResponse]) {
		resp = azfake.PagerResponder[fabcore.OneLakeShortcutsClientListShortcutsResponse]{}
		resp.AddPage(http.StatusOK, fabcore.OneLakeShortcutsClientListShortcutsResponse{Shortcuts: fabcore.Shortcuts{Value: GetAllStoredShortcuts()}}, nil)

		return
	}
}

func fakeGetShortcutFunc() func(ctx context.Context, workspaceID, itemID, path, name string, options *fabcore.OneLakeShortcutsClientGetShortcutOptions) (resp azfake.Responder[fabcore.OneLakeShortcutsClientGetShortcutResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID, path, name string, _ *fabcore.OneLakeShortcutsClientGetShortcutOptions) (resp azfake.Responder[fabcore.OneLakeShortcutsClientGetShortcutResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeShortcutsClientGetShortcutResponse]{}
		errItemNotFound := fabcore.ErrItem.ItemNotFound.Error()
		id := GenerateShortcutID(workspaceID, itemID, path, name)

		if shortcut, ok := fakeShortcutStore[id]; ok {
			resp.SetResponse(http.StatusOK, fabcore.OneLakeShortcutsClientGetShortcutResponse{Shortcut: shortcut}, nil)
		} else {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, errItemNotFound, "Item not found"))
			resp.SetResponse(http.StatusNotFound, fabcore.OneLakeShortcutsClientGetShortcutResponse{}, nil)
		}

		return
	}
}

func fakeCreateShortcutFunc() func(ctx context.Context, workspaceID, itemID string, createShortcutRequest fabcore.CreateShortcutRequest, options *fabcore.OneLakeShortcutsClientCreateShortcutOptions) (resp azfake.Responder[fabcore.OneLakeShortcutsClientCreateShortcutResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID string, createShortcutRequest fabcore.CreateShortcutRequest, _ *fabcore.OneLakeShortcutsClientCreateShortcutOptions) (resp azfake.Responder[fabcore.OneLakeShortcutsClientCreateShortcutResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeShortcutsClientCreateShortcutResponse]{}
		errItemAlreadyExists := fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error()
		id := GenerateShortcutID(workspaceID, itemID, *createShortcutRequest.Path, *createShortcutRequest.Name)

		requestShortcut := fabcore.Shortcut{
			Name: createShortcutRequest.Name,
			Path: createShortcutRequest.Path,
			Target: &fabcore.Target{
				OneLake: &fabcore.OneLake{
					ItemID:      createShortcutRequest.Target.OneLake.ItemID,
					WorkspaceID: createShortcutRequest.Target.OneLake.WorkspaceID,
					Path:        createShortcutRequest.Target.OneLake.Path,
				},
			},
		}

		if existing, ok := fakeShortcutStore[id]; ok {
			if *existing.Target.OneLake.ItemID == *createShortcutRequest.Target.OneLake.ItemID &&
				*existing.Target.OneLake.WorkspaceID == *createShortcutRequest.Target.OneLake.WorkspaceID &&
				*existing.Target.OneLake.Path == *createShortcutRequest.Target.OneLake.Path {
				errResp.SetError(fabfake.SetResponseError(http.StatusConflict, errItemAlreadyExists, "Item Display Name Already In Use"))
				resp.SetResponse(http.StatusConflict, fabcore.OneLakeShortcutsClientCreateShortcutResponse{Shortcut: existing}, nil)

				return resp, errResp
			}

			fakeShortcutStore[id] = requestShortcut
			resp.SetResponse(http.StatusOK, fabcore.OneLakeShortcutsClientCreateShortcutResponse{Shortcut: requestShortcut}, nil)

			return resp, errResp
		}

		fakeShortcutStore[id] = requestShortcut
		resp.SetResponse(http.StatusCreated, fabcore.OneLakeShortcutsClientCreateShortcutResponse{Shortcut: requestShortcut}, nil)

		return resp, errResp
	}
}

func fakeDeleteShortcutFunc() func(ctx context.Context, workspaceID, itemID, shortcutPath, shortcutName string, options *fabcore.OneLakeShortcutsClientDeleteShortcutOptions) (resp azfake.Responder[fabcore.OneLakeShortcutsClientDeleteShortcutResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID, path, name string, _ *fabcore.OneLakeShortcutsClientDeleteShortcutOptions) (resp azfake.Responder[fabcore.OneLakeShortcutsClientDeleteShortcutResponse], errResp azfake.ErrorResponder) {
		id := GenerateShortcutID(workspaceID, itemID, path, name)

		if _, ok := fakeShortcutStore[id]; ok {
			delete(fakeShortcutStore, id)
			resp.SetResponse(http.StatusOK, struct{}{}, nil)
		} else {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, "ItemNotFound", "Item not found"))
			resp.SetResponse(http.StatusNotFound, struct{}{}, nil)
		}

		return
	}
}

func NewRandomShortcut() fabcore.Shortcut {
	return fabcore.Shortcut{
		Name:   to.Ptr(testhelp.RandomName()),
		Path:   to.Ptr(testhelp.RandomName()),
		Target: NewRandomShortcutTarget(),
	}
}

func NewRandomShortcutTarget() *fabcore.Target {
	return &fabcore.Target{
		OneLake: NewRandomShortcutTargetOneLake(),
	}
}

func NewRandomShortcutTargetOneLake() *fabcore.OneLake {
	return &fabcore.OneLake{
		ItemID:      to.Ptr(testhelp.RandomUUID()),
		Path:        to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
	}
}

func GenerateShortcutID(workspaceID, itemID, path, name string) string {
	return fmt.Sprintf("%s/%s/%s/%s", workspaceID, itemID, path, name)
}

func GetAllStoredShortcuts() []fabcore.ShortcutTransformFlagged {
	shortcuts := make([]fabcore.ShortcutTransformFlagged, 0, len(fakeShortcutStore))
	for _, shortcut := range fakeShortcutStore {
		shortcuts = append(shortcuts, toShortcutTransformFlagged(shortcut))
	}

	return shortcuts
}

func fakeTestUpsert(workspaceID, itemID string, entity fabcore.Shortcut) {
	id := GenerateShortcutID(workspaceID, itemID, *entity.Path, *entity.Name)
	fakeShortcutStore[id] = entity
}

func toShortcutTransformFlagged(v fabcore.Shortcut) fabcore.ShortcutTransformFlagged {
	return fabcore.ShortcutTransformFlagged{
		Path:      v.Path,
		Name:      v.Name,
		Target:    v.Target,
		Transform: v.Transform,
	}
}
