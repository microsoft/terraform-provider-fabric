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

type operationsItem struct{}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsItem) CreateDefinition(data fabcore.CreateItemRequest) *fabcore.ItemDefinition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsItem) TransformDefinition(entity *fabcore.ItemDefinition) fabcore.ItemsClientGetItemDefinitionResponse {
	return fabcore.ItemsClientGetItemDefinitionResponse{
		ItemDefinitionResponse: fabcore.ItemDefinitionResponse{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsItem) UpdateDefinition(_ *fabcore.ItemDefinition, data fabcore.UpdateItemDefinitionRequest) *fabcore.ItemDefinition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsItem) CreateWithParentID(parentID string, data fabcore.CreateItemRequest) fabcore.Item {
	result := NewRandomItemWithWorkspace(*data.Type, parentID)
	result.DisplayName = data.DisplayName
	result.Description = data.Description
	result.FolderID = data.FolderID

	return result
}

// Filter implements concreteOperations.
func (o *operationsItem) Filter(entities []fabcore.Item, parentID string) []fabcore.Item {
	ret := make([]fabcore.Item, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsItem) GetID(entity fabcore.Item) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsItem) TransformCreate(entity fabcore.Item) fabcore.ItemsClientCreateItemResponse {
	return fabcore.ItemsClientCreateItemResponse{
		Item: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsItem) TransformGet(entity fabcore.Item) fabcore.ItemsClientGetItemResponse {
	return fabcore.ItemsClientGetItemResponse{
		Item: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsItem) TransformList(entities []fabcore.Item) fabcore.ItemsClientListItemsResponse {
	return fabcore.ItemsClientListItemsResponse{
		Items: fabcore.Items{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsItem) TransformUpdate(entity fabcore.Item) fabcore.ItemsClientUpdateItemResponse {
	return fabcore.ItemsClientUpdateItemResponse{
		Item: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsItem) Update(base fabcore.Item, data fabcore.UpdateItemRequest) fabcore.Item {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsItem) Validate(newEntity fabcore.Item, existing []fabcore.Item) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureItem(server *fakeServer) fabcore.Item {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabcore.Item,
			fabcore.ItemsClientGetItemResponse,
			fabcore.ItemsClientUpdateItemResponse,
			fabcore.ItemsClientCreateItemResponse,
			fabcore.ItemsClientListItemsResponse,
			fabcore.CreateItemRequest,
			fabcore.UpdateItemRequest]
	}

	type concreteDefinitionOperations interface {
		definitionOperations[
			fabcore.ItemDefinition,
			fabcore.CreateItemRequest,
			fabcore.UpdateItemDefinitionRequest,
			fabcore.ItemsClientGetItemDefinitionResponse,
			fabcore.ItemsClientUpdateItemDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsItem{}

	var definitionOperations concreteDefinitionOperations = &operationsItem{}

	handler := newTypedHandler(server, entityOperations)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.Core.ItemsServer.GetItem,
		&server.ServerFactory.Core.ItemsServer.UpdateItem,
		&server.ServerFactory.Core.ItemsServer.BeginCreateItem,
		&server.ServerFactory.Core.ItemsServer.NewListItemsPager,
		&server.ServerFactory.Core.ItemsServer.DeleteItem)

	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.Core.ItemsServer.BeginCreateItem,
		&server.ServerFactory.Core.ItemsServer.BeginGetItemDefinition,
		&server.ServerFactory.Core.ItemsServer.BeginUpdateItemDefinition)
	server.ServerFactory.Core.ItemsServer.MoveItem = FakeMoveItem(handler)

	return fabcore.Item{}
}

type moveItemOperations struct{}

func (m *moveItemOperations) TransformUpdate(entity fabcore.Item) fabcore.ItemsClientMoveItemResponse {
	return fabcore.ItemsClientMoveItemResponse{
		MovedItems: fabcore.MovedItems{
			Value: []fabcore.Item{
				entity,
			},
		},
	}
}

func (m *moveItemOperations) Update(base fabcore.Item, moveReq fabcore.MoveItemRequest) fabcore.Item {
	base.FolderID = moveReq.TargetFolderID

	return base
}

func FakeMoveItem(
	handler *typedHandler[fabcore.Item],
) func(ctx context.Context, workspaceID, itemID string, moveItemRequest fabcore.MoveItemRequest, options *fabcore.ItemsClientMoveItemOptions) (resp azfake.Responder[fabcore.ItemsClientMoveItemResponse], errResp azfake.ErrorResponder) {
	return func(_ context.Context, workspaceID, itemID string, moveReq fabcore.MoveItemRequest, _ *fabcore.ItemsClientMoveItemOptions) (azfake.Responder[fabcore.ItemsClientMoveItemResponse], azfake.ErrorResponder) {
		moveUpdater := &moveItemOperations{}
		moveTransformer := &moveItemOperations{}

		id := generateID(workspaceID, itemID)

		return updateByID(handler, id, moveReq, moveUpdater, moveTransformer)
	}
}

func NewRandomItem(itemType fabcore.ItemType) fabcore.Item {
	return fabcore.Item{
		Type:        &itemType,
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		FolderID:    to.Ptr(testhelp.RandomUUID()),
		Tags:        []fabcore.ItemTag{},
	}
}

func NewRandomItemWithWorkspace(itemType fabcore.ItemType, workspaceID string) fabcore.Item {
	result := NewRandomItem(itemType)
	result.WorkspaceID = &workspaceID

	return result
}
