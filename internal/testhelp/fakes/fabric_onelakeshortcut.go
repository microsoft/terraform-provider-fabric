// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsOneLakeShortcut struct{}

func (o *operationsOneLakeShortcut) GetID(entity fabcore.Shortcut) string {
	return *entity.Path + "/" + *entity.Name
}

func (o *operationsOneLakeShortcut) CreateWithWorkspaceIDAndItemID(_, _ string, request fabcore.CreateShortcutRequest) fabcore.Shortcut {
	entity := NewRandomOnelakeShortcut()
	entity.Name = request.Name
	entity.Path = request.Path
	entity.Target = &fabcore.Target{
		OneLake: &fabcore.OneLake{
			ItemID:      request.Target.OneLake.ItemID,
			WorkspaceID: request.Target.OneLake.WorkspaceID,
			Path:        request.Target.OneLake.Path,
		},
	}

	return entity
}

// TransformCreate implements concreteOperations.
func (o *operationsOneLakeShortcut) TransformCreate(entity fabcore.Shortcut) fabcore.OneLakeShortcutsClientCreateShortcutResponse {
	return fabcore.OneLakeShortcutsClientCreateShortcutResponse{
		Shortcut: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsOneLakeShortcut) TransformGet(entity fabcore.Shortcut) fabcore.OneLakeShortcutsClientGetShortcutResponse {
	return fabcore.OneLakeShortcutsClientGetShortcutResponse{
		Shortcut: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsOneLakeShortcut) TransformList(entities []fabcore.Shortcut) fabcore.OneLakeShortcutsClientListShortcutsResponse {
	{
		list := make([]fabcore.Shortcut, len(entities))
		for i, entity := range entities {
			list[i] = transformShortcut(entity)
		}

		return fabcore.OneLakeShortcutsClientListShortcutsResponse{
			Shortcuts: fabcore.Shortcuts{
				Value: list,
			},
		}
	}
}

func transformShortcut(entity fabcore.Shortcut) fabcore.Shortcut {
	return fabcore.Shortcut{
		Name:   entity.Name,
		Target: entity.Target,
		Path:   entity.Path,
	}
}

func (o *operationsOneLakeShortcut) Validate(newEntity fabcore.Shortcut, existing []fabcore.Shortcut) (int, error) {
	for _, entity := range existing {
		if *entity.Name == *newEntity.Name &&
			*entity.Path == *newEntity.Path &&
			*entity.Target.OneLake.ItemID == *newEntity.Target.OneLake.ItemID &&
			*entity.Target.OneLake.WorkspaceID == *newEntity.Target.OneLake.WorkspaceID &&
			*entity.Target.OneLake.Path == *newEntity.Target.OneLake.Path {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}

		if entity.Name == newEntity.Name && entity.Path == newEntity.Path {
			return http.StatusUpgradeRequired, nil
		}
	}

	return http.StatusCreated, nil
}

func configureOneLakeShortcut(server *fakeServer) fabcore.Shortcut {
	type concreteEntityOperations interface {
		onelakeOperations[
			fabcore.Shortcut,
			fabcore.OneLakeShortcutsClientGetShortcutResponse,
			fabcore.OneLakeShortcutsClientCreateShortcutResponse,
			fabcore.OneLakeShortcutsClientListShortcutsResponse,
			fabcore.CreateShortcutRequest]
	}

	var entityOperations concreteEntityOperations = &operationsOneLakeShortcut{}

	handler := newTypedHandler(server, entityOperations)

	configureOneLakeShortcutHandler(
		handler,
		entityOperations,
		&handler.ServerFactory.Core.OneLakeShortcutsServer.GetShortcut,
		&handler.ServerFactory.Core.OneLakeShortcutsServer.CreateShortcut,
		&handler.ServerFactory.Core.OneLakeShortcutsServer.NewListShortcutsPager,
		&handler.ServerFactory.Core.OneLakeShortcutsServer.DeleteShortcut)

	return fabcore.Shortcut{}
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
