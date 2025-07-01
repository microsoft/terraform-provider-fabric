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

// operationsWorkspace implements SimpleIDOperations.
type operationsWorkspace struct{}

func (o *operationsWorkspace) Create(data fabcore.CreateWorkspaceRequest) fabcore.WorkspaceInfo {
	entity := NewRandomWorkspaceInfo(data.CapacityID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description

	return entity
}

func (o *operationsWorkspace) TransformCreate(entity fabcore.WorkspaceInfo) fabcore.WorkspacesClientCreateWorkspaceResponse {
	return fabcore.WorkspacesClientCreateWorkspaceResponse{
		Workspace: transformWorkspace(entity),
	}
}

func (o *operationsWorkspace) TransformGet(entity fabcore.WorkspaceInfo) fabcore.WorkspacesClientGetWorkspaceResponse {
	return fabcore.WorkspacesClientGetWorkspaceResponse{
		WorkspaceInfo: entity,
	}
}

func (o *operationsWorkspace) TransformList(entities []fabcore.WorkspaceInfo) fabcore.WorkspacesClientListWorkspacesResponse {
	list := make([]fabcore.Workspace, len(entities))
	for i, entity := range entities {
		list[i] = transformWorkspace(entity)
	}

	return fabcore.WorkspacesClientListWorkspacesResponse{
		Workspaces: fabcore.Workspaces{
			Value: list,
		},
	}
}

func (o *operationsWorkspace) TransformUpdate(entity fabcore.WorkspaceInfo) fabcore.WorkspacesClientUpdateWorkspaceResponse {
	return fabcore.WorkspacesClientUpdateWorkspaceResponse{
		Workspace: transformWorkspace(entity),
	}
}

func (o *operationsWorkspace) Update(base fabcore.WorkspaceInfo, data fabcore.UpdateWorkspaceRequest) fabcore.WorkspaceInfo {
	base.DisplayName = data.DisplayName
	base.Description = data.Description

	return base
}

func (o *operationsWorkspace) Validate(newEntity fabcore.WorkspaceInfo, existing []fabcore.WorkspaceInfo) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrWorkspace.WorkspaceNameAlreadyExists.Error(), fabcore.ErrWorkspace.WorkspaceNameAlreadyExists.Error())
		}
	}

	return http.StatusCreated, nil
}

func (o *operationsWorkspace) GetID(entity fabcore.WorkspaceInfo) string {
	return *entity.ID
}

func transformWorkspace(entity fabcore.WorkspaceInfo) fabcore.Workspace {
	return fabcore.Workspace{
		ID:          entity.ID,
		DisplayName: entity.DisplayName,
		Description: entity.Description,
		CapacityID:  entity.CapacityID,
		Type:        entity.Type,
	}
}

func configureWorkspace(server *fakeServer) fabcore.WorkspaceInfo {
	type concreteEntityOperations interface {
		simpleIDOperations[
			fabcore.WorkspaceInfo,
			fabcore.WorkspacesClientGetWorkspaceResponse,
			fabcore.WorkspacesClientUpdateWorkspaceResponse,
			fabcore.WorkspacesClientCreateWorkspaceResponse,
			fabcore.WorkspacesClientListWorkspacesResponse,
			fabcore.CreateWorkspaceRequest,
			fabcore.UpdateWorkspaceRequest]
	}

	var entityOperations concreteEntityOperations = &operationsWorkspace{}

	handler := newTypedHandler(server, entityOperations)

	configureEntityPagerWithSimpleID(
		handler,
		entityOperations,
		&handler.ServerFactory.Core.WorkspacesServer.GetWorkspace,
		&handler.ServerFactory.Core.WorkspacesServer.UpdateWorkspace,
		&handler.ServerFactory.Core.WorkspacesServer.CreateWorkspace,
		&handler.ServerFactory.Core.WorkspacesServer.NewListWorkspacesPager,
		&handler.ServerFactory.Core.WorkspacesServer.DeleteWorkspace)

	return fabcore.WorkspaceInfo{}
}

func NewRandomWorkspaceInfo(capacityID *string) fabcore.WorkspaceInfo {
	return fabcore.WorkspaceInfo{
		ID:                         to.Ptr(testhelp.RandomUUID()),
		DisplayName:                to.Ptr(testhelp.RandomName()),
		Description:                to.Ptr(testhelp.RandomName()),
		Type:                       to.Ptr(fabcore.WorkspaceTypeWorkspace),
		CapacityID:                 capacityID,
		CapacityRegion:             to.Ptr(fabcore.CapacityRegionWestUS2),
		CapacityAssignmentProgress: to.Ptr(fabcore.CapacityAssignmentProgressCompleted),
		OneLakeEndpoints: &fabcore.OneLakeEndpoints{
			BlobEndpoint: to.Ptr(testhelp.RandomURI()),
			DfsEndpoint:  to.Ptr(testhelp.RandomURI()),
		},
	}
}

func NewRandomWorkspaceInfoWithType(entityType fabcore.WorkspaceType, capacityID *string) fabcore.WorkspaceInfo {
	entity := NewRandomWorkspaceInfo(capacityID)
	entity.Type = &entityType

	return entity
}

func NewRandomWorkspaceInfoWithIdentity(capacityID *string) fabcore.WorkspaceInfo {
	entity := NewRandomWorkspaceInfo(capacityID)
	entity.WorkspaceIdentity = &fabcore.WorkspaceIdentity{
		ApplicationID:      to.Ptr(testhelp.RandomUUID()),
		ServicePrincipalID: to.Ptr(testhelp.RandomUUID()),
	}

	return entity
}
