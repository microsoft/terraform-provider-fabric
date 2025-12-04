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

// ConvertItemToEntity implements itemConverter.
func (o *operationsEventhouse) ConvertItemToEntity(item fabcore.Item) fabeventhouse.Eventhouse {
	return fabeventhouse.Eventhouse{
		ID:          item.ID,
		DisplayName: item.DisplayName,
		Description: item.Description,
		WorkspaceID: item.WorkspaceID,
		FolderID:    item.FolderID,
		Type:        to.Ptr(fabeventhouse.ItemTypeEventhouse),
		Properties:  NewRandomEventhouse().Properties,
	}
}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsEventhouse) CreateDefinition(data fabeventhouse.CreateEventhouseRequest) *fabeventhouse.Definition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsEventhouse) TransformDefinition(entity *fabeventhouse.Definition) fabeventhouse.ItemsClientGetEventhouseDefinitionResponse {
	return fabeventhouse.ItemsClientGetEventhouseDefinitionResponse{
		DefinitionResponse: fabeventhouse.DefinitionResponse{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsEventhouse) UpdateDefinition(_ *fabeventhouse.Definition, data fabeventhouse.UpdateEventhouseDefinitionRequest) *fabeventhouse.Definition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsEventhouse) CreateWithParentID(parentID string, data fabeventhouse.CreateEventhouseRequest) fabeventhouse.Eventhouse {
	entity := NewRandomEventhouseWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description
	entity.FolderID = data.FolderID

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
	base.DisplayName = data.DisplayName
	base.Description = data.Description

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

	type concreteDefinitionOperations interface {
		definitionOperations[
			fabeventhouse.Definition,
			fabeventhouse.CreateEventhouseRequest,
			fabeventhouse.UpdateEventhouseDefinitionRequest,
			fabeventhouse.ItemsClientGetEventhouseDefinitionResponse,
			fabeventhouse.ItemsClientUpdateEventhouseDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsEventhouse{}

	var definitionOperations concreteDefinitionOperations = &operationsEventhouse{}

	var converter itemConverter[fabeventhouse.Eventhouse] = &operationsEventhouse{}

	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.Eventhouse.ItemsServer.GetEventhouse,
		&server.ServerFactory.Eventhouse.ItemsServer.UpdateEventhouse,
		&server.ServerFactory.Eventhouse.ItemsServer.BeginCreateEventhouse,
		&server.ServerFactory.Eventhouse.ItemsServer.NewListEventhousesPager,
		&server.ServerFactory.Eventhouse.ItemsServer.DeleteEventhouse)

	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.Eventhouse.ItemsServer.BeginCreateEventhouse,
		&server.ServerFactory.Eventhouse.ItemsServer.BeginGetEventhouseDefinition,
		&server.ServerFactory.Eventhouse.ItemsServer.BeginUpdateEventhouseDefinition)

	return fabeventhouse.Eventhouse{}
}

func NewRandomEventhouse() fabeventhouse.Eventhouse {
	return fabeventhouse.Eventhouse{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		FolderID:    to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fabeventhouse.ItemTypeEventhouse),
		Properties: &fabeventhouse.Properties{
			IngestionServiceURI:     to.Ptr(testhelp.RandomURI()),
			QueryServiceURI:         to.Ptr(testhelp.RandomURI()),
			DatabasesItemIDs:        []string{testhelp.RandomUUID()},
			MinimumConsumptionUnits: to.Ptr(0.0),
		},
	}
}

func NewRandomEventhouseWithWorkspace(workspaceID string) fabeventhouse.Eventhouse {
	result := NewRandomEventhouse()
	result.WorkspaceID = &workspaceID

	return result
}

func NewRandomEventhouseDefinition() fabeventhouse.Definition {
	defPart := fabeventhouse.DefinitionPart{
		PayloadType: to.Ptr(fabeventhouse.PayloadTypeInlineBase64),
		Path:        to.Ptr("EventhouseProperties.json"),
		Payload:     to.Ptr("e30="),
	}

	var defParts []fabeventhouse.DefinitionPart

	defParts = append(defParts, defPart)

	return fabeventhouse.Definition{
		Parts: defParts,
	}
}
