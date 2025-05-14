package onelakeshortcut_test

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

var fakeOneLakeShortcutStore = map[string]fabcore.Shortcut{}

func fakeOneLakeShortcutsFunc(
	onelakeshortcuts fabcore.Shortcuts,
) func(workspaceID, itemID string, options *fabcore.OneLakeShortcutsClientListShortcutsOptions) (resp azfake.PagerResponder[fabcore.OneLakeShortcutsClientListShortcutsResponse]) {
	return func(_, _ string, _ *fabcore.OneLakeShortcutsClientListShortcutsOptions) (resp azfake.PagerResponder[fabcore.OneLakeShortcutsClientListShortcutsResponse]) {
		resp = azfake.PagerResponder[fabcore.OneLakeShortcutsClientListShortcutsResponse]{}
		resp.AddPage(http.StatusOK, fabcore.OneLakeShortcutsClientListShortcutsResponse{Shortcuts: onelakeshortcuts}, nil)

		return
	}
}

func fakeGetOneLakeShortcutFunc() func(ctx context.Context, workspaceID, itemID, path, name string, options *fabcore.OneLakeShortcutsClientGetShortcutOptions) (resp azfake.Responder[fabcore.OneLakeShortcutsClientGetShortcutResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID, path, name string, _ *fabcore.OneLakeShortcutsClientGetShortcutOptions) (resp azfake.Responder[fabcore.OneLakeShortcutsClientGetShortcutResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeShortcutsClientGetShortcutResponse]{}
		errItemNotFound := fabcore.ErrItem.ItemNotFound.Error()
		id := fmt.Sprintf("%s/%s/%s/%s", workspaceID, itemID, path, name)
		if shortcut, ok := fakeOneLakeShortcutStore[id]; ok {
			resp.SetResponse(http.StatusOK, fabcore.OneLakeShortcutsClientGetShortcutResponse{Shortcut: shortcut}, nil)
		} else {

			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, errItemNotFound, "Item not found"))
			resp.SetResponse(http.StatusNotFound, fabcore.OneLakeShortcutsClientGetShortcutResponse{}, nil)
		}

		return
	}
}

func fakeCreateOneLakeShortcutFunc() func(ctx context.Context, workspaceID, itemID string, createShortcutRequest fabcore.CreateShortcutRequest, options *fabcore.OneLakeShortcutsClientCreateShortcutOptions) (resp azfake.Responder[fabcore.OneLakeShortcutsClientCreateShortcutResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID string, createShortcutRequest fabcore.CreateShortcutRequest, _ *fabcore.OneLakeShortcutsClientCreateShortcutOptions) (resp azfake.Responder[fabcore.OneLakeShortcutsClientCreateShortcutResponse], errResp azfake.ErrorResponder) {
		resp = azfake.Responder[fabcore.OneLakeShortcutsClientCreateShortcutResponse]{}

		errItemAlreadyExists := fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error()
		created := fabcore.Shortcut{
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

		id := fmt.Sprintf("%s/%s/%s/%s", workspaceID, itemID, *createShortcutRequest.Path, *createShortcutRequest.Name)

		if existing, ok := fakeOneLakeShortcutStore[id]; ok {
			// Check if the target details also match
			if existing.Target != nil && existing.Target.OneLake != nil &&
				createShortcutRequest.Target != nil && createShortcutRequest.Target.OneLake != nil &&
				*existing.Target.OneLake.ItemID == *createShortcutRequest.Target.OneLake.ItemID &&
				*existing.Target.OneLake.WorkspaceID == *createShortcutRequest.Target.OneLake.WorkspaceID &&
				*existing.Target.OneLake.Path == *createShortcutRequest.Target.OneLake.Path {

				// Only then: return conflict
				errResp.SetError(fabfake.SetResponseError(http.StatusConflict, errItemAlreadyExists, "Item Display Name Already In Use"))
				resp.SetResponse(http.StatusConflict, fabcore.OneLakeShortcutsClientCreateShortcutResponse{Shortcut: existing}, nil)
				return
			}
			fakeOneLakeShortcutStore[id] = created
			resp.SetResponse(http.StatusOK, fabcore.OneLakeShortcutsClientCreateShortcutResponse{Shortcut: created}, nil)
			return
		}

		fakeOneLakeShortcutStore[id] = created

		// No match, return 200 OK with first existing shortcut as placeholder response
		resp.SetResponse(http.StatusOK, fabcore.OneLakeShortcutsClientCreateShortcutResponse{Shortcut: created}, nil)
		return
	}
}

func fakeDeleteOneLakeShortcutFunc() func(ctx context.Context, workspaceID, itemID, shortcutPath, shortcutName string, options *fabcore.OneLakeShortcutsClientDeleteShortcutOptions) (resp azfake.Responder[fabcore.OneLakeShortcutsClientDeleteShortcutResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID, path, name string, _ *fabcore.OneLakeShortcutsClientDeleteShortcutOptions) (resp azfake.Responder[fabcore.OneLakeShortcutsClientDeleteShortcutResponse], errResp azfake.ErrorResponder) {
		id := fmt.Sprintf("%s/%s/%s/%s", workspaceID, itemID, path, name)

		if _, ok := fakeOneLakeShortcutStore[id]; ok {
			delete(fakeOneLakeShortcutStore, id)
			resp.SetResponse(http.StatusOK, struct{}{}, nil)
		} else {
			errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, "ItemNotFound", "Item not found"))
			resp.SetResponse(http.StatusNotFound, struct{}{}, nil)
		}

		return
	}
}

func NewRandomOneLakeShortcuts(shortcuts []fabcore.Shortcut) fabcore.Shortcuts {
	copied := make([]fabcore.Shortcut, len(shortcuts))

	for i, shortcut := range shortcuts {
		copied[i] = fabcore.Shortcut{
			Name: shortcut.Name,
			Path: shortcut.Path,
			Target: &fabcore.Target{
				OneLake: &fabcore.OneLake{
					ItemID:      shortcut.Target.OneLake.ItemID,
					WorkspaceID: shortcut.Target.OneLake.WorkspaceID,
					Path:        shortcut.Target.OneLake.Path,
				},
			},
		}
	}

	return fabcore.Shortcuts{
		Value: copied,
	}
}

func NewRandomOnelakeShortcutWithWorkspaceIDAndItemID(workspaceID, itemID string) fabcore.Shortcut {
	entity := NewRandomOnelakeShortcut()
	id := fmt.Sprintf("%s/%s/%s/%s", workspaceID, itemID, *entity.Path, *entity.Name)
	fakeOneLakeShortcutStore[id] = entity
	return entity
}

func NewRandomOnelakeShortcut() fabcore.Shortcut {
	return fabcore.Shortcut{
		Name:   to.Ptr(testhelp.RandomName()),
		Path:   to.Ptr(testhelp.RandomName()),
		Target: NewRandomOnelakeShortcutTarget(),
	}
}

func NewRandomOnelakeShortcutTarget() *fabcore.Target {
	return &fabcore.Target{
		OneLake: NewRandomOneLakeShortcutTargetOneLake(),
	}
}

func NewRandomOneLakeShortcutTargetOneLake() *fabcore.OneLake {
	return &fabcore.OneLake{
		ItemID:      to.Ptr(testhelp.RandomUUID()),
		Path:        to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
	}
}
