// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsFolder struct{}

// GetID implements concreteOperations.
func (o *operationsFolder) GetID(entity fabcore.Folder) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsFolder) TransformCreate(entity fabcore.Folder) fabcore.FoldersClientCreateFolderResponse {
	return fabcore.FoldersClientCreateFolderResponse{
		Folder: entity,
	}
}

// CreateWithParentID implements concreteOperations.
func (o *operationsFolder) CreateWithParentID(parentID string, data fabcore.CreateFolderRequest) fabcore.Folder {
	entity := NewRandomFolderWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.ParentFolderID = data.ParentFolderID

	return entity
}

// Filter implements concreteOperations.
func (o *operationsFolder) Filter(entities []fabcore.Folder, _ string) []fabcore.Folder {
	return entities
}

// TransformGet implements concreteOperations.
func (o *operationsFolder) TransformGet(entity fabcore.Folder) fabcore.FoldersClientGetFolderResponse {
	return fabcore.FoldersClientGetFolderResponse{
		Folder: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsFolder) TransformList(entities []fabcore.Folder) fabcore.FoldersClientListFoldersResponse {
	return fabcore.FoldersClientListFoldersResponse{
		Folders: fabcore.Folders{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsFolder) TransformUpdate(entity fabcore.Folder) fabcore.FoldersClientUpdateFolderResponse {
	return fabcore.FoldersClientUpdateFolderResponse{
		Folder: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsFolder) Update(base fabcore.Folder, data fabcore.UpdateFolderRequest) fabcore.Folder {
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsFolder) Validate(newEntity fabcore.Folder, existing []fabcore.Folder) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName && *entity.WorkspaceID == *newEntity.WorkspaceID {
			if (entity.ParentFolderID != nil && newEntity.ParentFolderID != nil && *entity.ParentFolderID == *newEntity.ParentFolderID) ||
				(entity.ParentFolderID == nil && newEntity.ParentFolderID == nil) {
				return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
			}
		}
	}

	return http.StatusCreated, nil
}

// FilterWithOptions filters folders by workspace and optionally by root folder and recursion.
func (o *operationsFolder) FilterWithOptions(entities []fabcore.Folder, workspaceID string, options *fabcore.FoldersClientListFoldersOptions) []fabcore.Folder {
	// Filter by workspace
	workspaceFolders := make([]fabcore.Folder, 0)

	for _, folder := range entities {
		if *folder.WorkspaceID == workspaceID {
			workspaceFolders = append(workspaceFolders, folder)
		}
	}

	// No root folder specified
	if options == nil || options.RootFolderID == nil {
		return workspaceFolders
	}

	recursive := *options.Recursive // Default true
	if recursive {
		return o.getDescendants(*options.RootFolderID, workspaceFolders)
	}

	return o.getChildren(*options.RootFolderID, workspaceFolders)
}

func FakeMoveFolder(
	handler *typedHandler[fabcore.Folder],
) func(ctx context.Context, workspaceID, folderID string, moveFolderRequest fabcore.MoveFolderRequest, options *fabcore.FoldersClientMoveFolderOptions) (resp azfake.Responder[fabcore.FoldersClientMoveFolderResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, folderID string, moveReq fabcore.MoveFolderRequest, _ *fabcore.FoldersClientMoveFolderOptions) (azfake.Responder[fabcore.FoldersClientMoveFolderResponse], azfake.ErrorResponder) {
		moveUpdater := &moveFolderOperations{}
		moveTransformer := &moveFolderOperations{}

		id := generateID(workspaceID, folderID)

		return updateByID(handler, id, moveReq, moveUpdater, moveTransformer)
	}
}

type moveFolderOperations struct{}

func (m *moveFolderOperations) TransformUpdate(entity fabcore.Folder) fabcore.FoldersClientMoveFolderResponse {
	return fabcore.FoldersClientMoveFolderResponse{
		Folder: entity,
	}
}

func (m *moveFolderOperations) Update(base fabcore.Folder, moveReq fabcore.MoveFolderRequest) fabcore.Folder {
	base.ParentFolderID = moveReq.TargetFolderID

	return base
}

// FakeListFolders creates a custom list folders handler that supports options.
func FakeListFolders(
	handler *typedHandler[fabcore.Folder],
) func(workspaceID string, options *fabcore.FoldersClientListFoldersOptions) (resp azfake.PagerResponder[fabcore.FoldersClientListFoldersResponse]) {
	return func(workspaceID string, options *fabcore.FoldersClientListFoldersOptions) (resp azfake.PagerResponder[fabcore.FoldersClientListFoldersResponse]) {
		entityOperations := &operationsFolder{}

		// Get all entities and filter them
		allEntities := handler.Elements()
		filteredEntities := entityOperations.FilterWithOptions(allEntities, workspaceID, options)

		// Transform to response format
		response := entityOperations.TransformList(filteredEntities)

		// Create and return pager
		var pager azfake.PagerResponder[fabcore.FoldersClientListFoldersResponse]
		pager.AddPage(200, response, nil)

		return pager
	}
}

// getChildren returns direct children of a parent folder.
func (o *operationsFolder) getChildren(parentID string, folders []fabcore.Folder) []fabcore.Folder {
	result := make([]fabcore.Folder, 0)

	for _, folder := range folders {
		if folder.ParentFolderID != nil && *folder.ParentFolderID == parentID {
			result = append(result, folder)
		}
	}

	return result
}

// getDescendants returns all descendants of a parent folder recursively.
func (o *operationsFolder) getDescendants(parentID string, folders []fabcore.Folder) []fabcore.Folder {
	result := make([]fabcore.Folder, 0)
	children := o.getChildren(parentID, folders)

	// Add direct children
	result = append(result, children...)

	// Add descendants of each child
	for _, child := range children {
		descendants := o.getDescendants(*child.ID, folders)
		result = append(result, descendants...)
	}

	return result
}

func configureFolder(server *fakeServer) fabcore.Folder {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabcore.Folder,
			fabcore.FoldersClientGetFolderResponse,
			fabcore.FoldersClientUpdateFolderResponse,
			fabcore.FoldersClientCreateFolderResponse,
			fabcore.FoldersClientListFoldersResponse,
			fabcore.CreateFolderRequest,
			fabcore.UpdateFolderRequest]
	}

	var entityOperations concreteEntityOperations = &operationsFolder{}

	handler := newTypedHandler(server, entityOperations)

	configureEntityWithParentIDNoLRO(
		handler,
		entityOperations,
		&server.ServerFactory.Core.FoldersServer.GetFolder,
		&server.ServerFactory.Core.FoldersServer.UpdateFolder,
		&server.ServerFactory.Core.FoldersServer.CreateFolder,
		&server.ServerFactory.Core.FoldersServer.NewListFoldersPager,
		&server.ServerFactory.Core.FoldersServer.DeleteFolder)

	server.ServerFactory.Core.FoldersServer.NewListFoldersPager = FakeListFolders(handler)
	server.ServerFactory.Core.FoldersServer.MoveFolder = FakeMoveFolder(handler)

	return fabcore.Folder{}
}

func NewRandomFolder() fabcore.Folder {
	return fabcore.Folder{
		ID:             to.Ptr(testhelp.RandomUUID()),
		DisplayName:    to.Ptr(testhelp.RandomName()),
		ParentFolderID: to.Ptr(testhelp.RandomUUID()),
	}
}

func NewRandomFolderWithWorkspace(workspaceID string) fabcore.Folder {
	result := NewRandomFolder()
	result.WorkspaceID = &workspaceID

	return result
}

func NewRandomSubfolder(workspaceID, parentFolderID string) fabcore.Folder {
	result := NewRandomFolderWithWorkspace(workspaceID)
	result.ParentFolderID = &parentFolderID

	return result
}
