// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabeventhouse "github.com/microsoft/fabric-sdk-go/fabric/eventhouse"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsEventhouse struct{}

// CreateWithParentID implements concreteOperations.
func (o *operationsEventhouse) CreateWithParentID(parentID string, data fabeventhouse.CreateEventhouseRequest) fabeventhouse.Eventhouse {
	entity := NewRandomEventhouseWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description

	return entity
}

// Filter implements concreteOperations.
func (o *operationsEventhouse) Filter(entities []fabeventhouse.Eventhouse, parentID string) []fabeventhouse.Eventhouse {
	ret := make([]fabeventhouse.Eventhouse, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsEventhouse) GetID(entity fabeventhouse.Eventhouse) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsEventhouse) TransformCreate(entity fabeventhouse.Eventhouse) fabeventhouse.ItemsClientCreateEventhouseResponse {
	return fabeventhouse.ItemsClientCreateEventhouseResponse{
		Eventhouse: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsEventhouse) TransformGet(entity fabeventhouse.Eventhouse) fabeventhouse.ItemsClientGetEventhouseResponse {
	return fabeventhouse.ItemsClientGetEventhouseResponse{
		Eventhouse: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsEventhouse) TransformList(entities []fabeventhouse.Eventhouse) fabeventhouse.ItemsClientListEventhousesResponse {
	return fabeventhouse.ItemsClientListEventhousesResponse{
		Eventhouses: fabeventhouse.Eventhouses{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsEventhouse) TransformUpdate(entity fabeventhouse.Eventhouse) fabeventhouse.ItemsClientUpdateEventhouseResponse {
	return fabeventhouse.ItemsClientUpdateEventhouseResponse{
		Eventhouse: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsEventhouse) Update(base fabeventhouse.Eventhouse, data fabeventhouse.UpdateEventhouseRequest) fabeventhouse.Eventhouse {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsEventhouse) Validate(newEntity fabeventhouse.Eventhouse, existing []fabeventhouse.Eventhouse) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureEventhouse(server *fakeServer) fabeventhouse.Eventhouse {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabeventhouse.Eventhouse,
			fabeventhouse.ItemsClientGetEventhouseResponse,
			fabeventhouse.ItemsClientUpdateEventhouseResponse,
			fabeventhouse.ItemsClientCreateEventhouseResponse,
			fabeventhouse.ItemsClientListEventhousesResponse,
			fabeventhouse.CreateEventhouseRequest,
			fabeventhouse.UpdateEventhouseRequest]
	}

	var entityOperations concreteEntityOperations = &operationsEventhouse{}

	handler := newTypedHandler(server, entityOperations)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.Eventhouse.ItemsServer.GetEventhouse,
		&server.ServerFactory.Eventhouse.ItemsServer.UpdateEventhouse,
		&server.ServerFactory.Eventhouse.ItemsServer.BeginCreateEventhouse,
		&server.ServerFactory.Eventhouse.ItemsServer.NewListEventhousesPager,
		&server.ServerFactory.Eventhouse.ItemsServer.DeleteEventhouse)

	return fabeventhouse.Eventhouse{}
}

func NewRandomEventhouse() fabeventhouse.Eventhouse {
	return fabeventhouse.Eventhouse{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fabeventhouse.ItemTypeEventhouse),
		Properties: &fabeventhouse.Properties{
			IngestionServiceURI: to.Ptr(testhelp.RandomURI()),
			QueryServiceURI:     to.Ptr(testhelp.RandomURI()),
			DatabasesItemIDs:    []string{testhelp.RandomUUID()},
		},
	}
}

func NewRandomEventhouseWithWorkspace(workspaceID string) fabeventhouse.Eventhouse {
	result := NewRandomEventhouse()
	result.WorkspaceID = &workspaceID

	return result
}
