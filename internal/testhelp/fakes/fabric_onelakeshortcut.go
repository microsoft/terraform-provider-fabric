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
	return *entity.Target.OneLake.ItemID
}

func (o *operationsOneLakeShortcut) GetIDWithParentID(_ string, entity fabcore.Shortcut) string {
	return *entity.Target.OneLake.ItemID
}

func (o *operationsOneLakeShortcut) Create(data fabcore.CreateShortcutRequest) fabcore.Shortcut {
	entity := NewRandomOnelakeShortcut()
	entity.Name = data.Name

	return entity
}

// CreateWithParentID implements concreteOperations.
func (o *operationsOneLakeShortcut) CreateWithParentID(_ string, data fabcore.Shortcut) fabcore.Shortcut {
	entity := NewRandomOnelakeShortcut()
	entity.Name = data.Name

	return entity
}

// Filter implements concreteOperations.
func (o *operationsOneLakeShortcut) Filter(entities []fabcore.Shortcut, _ string) []fabcore.Shortcut {
	ret := make([]fabcore.Shortcut, 0)

	ret = append(ret, entities...)

	return ret
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

func (o *operationsOneLakeShortcut) TransformUpdate(entity fabcore.Shortcut) fabcore.OneLakeShortcutsClientCreateShortcutResponse {
	return fabcore.OneLakeShortcutsClientCreateShortcutResponse{
		Shortcut: transformShortcut(entity),
	}
}

func (o *operationsOneLakeShortcut) Update(base fabcore.Shortcut, data fabcore.CreateShortcutRequest) fabcore.Shortcut {
	base.Name = data.Name

	return base
}

func (o *operationsOneLakeShortcut) Validate(newEntity fabcore.Shortcut, existing []fabcore.Shortcut) (int, error) {
	for _, entity := range existing {
		if *entity.Name == *newEntity.Name {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrWorkspace.WorkspaceNameAlreadyExists.Error(), fabcore.ErrWorkspace.WorkspaceNameAlreadyExists.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureOneLakeShortcut(server *fakeServer) fabcore.Shortcut {
	type concreteEntityOperations interface {
		simpleIDOperations[
			fabcore.Shortcut,
			fabcore.OneLakeShortcutsClientGetShortcutResponse,
			fabcore.OneLakeShortcutsClientCreateShortcutResponse,
			fabcore.OneLakeShortcutsClientCreateShortcutResponse,
			fabcore.OneLakeShortcutsClientListShortcutsResponse,
			fabcore.CreateShortcutRequest,
			fabcore.CreateShortcutRequest]
	}

	// var entityOperations concreteEntityOperations = &operationsOneLakeShortcut{}

	// handler := newTypedHandler(server, entityOperations)

	// configureEntityPagerWithSimpleID(
	// 	handler,
	// 	entityOperations,
	// 	&handler.ServerFactory.Core.OneLakeShortcutsServer.GetShortcut,
	// 	&handler.ServerFactory.Core.OneLakeShortcutsServer.CreateShortcut,
	// 	&handler.ServerFactory.Core.OneLakeShortcutsServer.CreateShortcut,
	// 	&handler.ServerFactory.Core.OneLakeShortcutsServer.NewListShortcutsPager,
	// 	&handler.ServerFactory.Core.OneLakeShortcutsServer.DeleteShortcut)

	return fabcore.Shortcut{}
}

func NewRandomOnelakeShortcut() fabcore.Shortcut {
	return fabcore.Shortcut{
		Name:   to.Ptr(testhelp.RandomName()),
		Path:   to.Ptr(testhelp.RandomURI()),
		Target: NewRandomOnelakeShortcutTarget(),
	}
}

func NewRandomOnelakeShortcutTarget() *fabcore.Target {
	return &fabcore.Target{
		Type:    to.Ptr(fabcore.TypeOneLake),
		OneLake: NewRandomOneLakeShortcutTargetOneLake(),
	}
}

func NewRandomOneLakeShortcutTargetOneLake() *fabcore.OneLake {
	return &fabcore.OneLake{
		ItemID:      to.Ptr(testhelp.RandomUUID()),
		Path:        to.Ptr(testhelp.RandomURI()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
	}
}
