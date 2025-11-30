// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabsemanticmodel "github.com/microsoft/fabric-sdk-go/fabric/semanticmodel"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsSemanticModel struct{}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsSemanticModel) CreateDefinition(data fabsemanticmodel.CreateSemanticModelRequest) *fabsemanticmodel.Definition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsSemanticModel) TransformDefinition(entity *fabsemanticmodel.Definition) fabsemanticmodel.ItemsClientGetSemanticModelDefinitionResponse {
	return fabsemanticmodel.ItemsClientGetSemanticModelDefinitionResponse{
		DefinitionResponse: fabsemanticmodel.DefinitionResponse{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsSemanticModel) UpdateDefinition(_ *fabsemanticmodel.Definition, data fabsemanticmodel.UpdateSemanticModelDefinitionRequest) *fabsemanticmodel.Definition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsSemanticModel) CreateWithParentID(parentID string, data fabsemanticmodel.CreateSemanticModelRequest) fabsemanticmodel.SemanticModel {
	entity := NewRandomSemanticModelWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description
	entity.FolderID = data.FolderID

	return entity
}

// Filter implements concreteOperations.
func (o *operationsSemanticModel) Filter(entities []fabsemanticmodel.SemanticModel, parentID string) []fabsemanticmodel.SemanticModel {
	ret := make([]fabsemanticmodel.SemanticModel, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsSemanticModel) GetID(entity fabsemanticmodel.SemanticModel) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsSemanticModel) TransformCreate(entity fabsemanticmodel.SemanticModel) fabsemanticmodel.ItemsClientCreateSemanticModelResponse {
	return fabsemanticmodel.ItemsClientCreateSemanticModelResponse{
		SemanticModel: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsSemanticModel) TransformGet(entity fabsemanticmodel.SemanticModel) fabsemanticmodel.ItemsClientGetSemanticModelResponse {
	return fabsemanticmodel.ItemsClientGetSemanticModelResponse{
		SemanticModel: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsSemanticModel) TransformList(entities []fabsemanticmodel.SemanticModel) fabsemanticmodel.ItemsClientListSemanticModelsResponse {
	return fabsemanticmodel.ItemsClientListSemanticModelsResponse{
		SemanticModels: fabsemanticmodel.SemanticModels{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsSemanticModel) TransformUpdate(entity fabsemanticmodel.SemanticModel) fabsemanticmodel.ItemsClientUpdateSemanticModelResponse {
	return fabsemanticmodel.ItemsClientUpdateSemanticModelResponse{
		SemanticModel: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsSemanticModel) Update(base fabsemanticmodel.SemanticModel, data fabsemanticmodel.UpdateSemanticModelRequest) fabsemanticmodel.SemanticModel {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsSemanticModel) Validate(newEntity fabsemanticmodel.SemanticModel, existing []fabsemanticmodel.SemanticModel) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureSemanticModel(server *fakeServer) fabsemanticmodel.SemanticModel {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabsemanticmodel.SemanticModel,
			fabsemanticmodel.ItemsClientGetSemanticModelResponse,
			fabsemanticmodel.ItemsClientUpdateSemanticModelResponse,
			fabsemanticmodel.ItemsClientCreateSemanticModelResponse,
			fabsemanticmodel.ItemsClientListSemanticModelsResponse,
			fabsemanticmodel.CreateSemanticModelRequest,
			fabsemanticmodel.UpdateSemanticModelRequest]
	}

	type concreteDefinitionOperations interface {
		definitionOperations[
			fabsemanticmodel.Definition,
			fabsemanticmodel.CreateSemanticModelRequest,
			fabsemanticmodel.UpdateSemanticModelDefinitionRequest,
			fabsemanticmodel.ItemsClientGetSemanticModelDefinitionResponse,
			fabsemanticmodel.ItemsClientUpdateSemanticModelDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsSemanticModel{}

	var definitionOperations concreteDefinitionOperations = &operationsSemanticModel{}

	handler := newTypedHandler(server, entityOperations)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.SemanticModel.ItemsServer.GetSemanticModel,
		&server.ServerFactory.SemanticModel.ItemsServer.UpdateSemanticModel,
		&server.ServerFactory.SemanticModel.ItemsServer.BeginCreateSemanticModel,
		&server.ServerFactory.SemanticModel.ItemsServer.NewListSemanticModelsPager,
		&server.ServerFactory.SemanticModel.ItemsServer.DeleteSemanticModel)

	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.SemanticModel.ItemsServer.BeginCreateSemanticModel,
		&server.ServerFactory.SemanticModel.ItemsServer.BeginGetSemanticModelDefinition,
		&server.ServerFactory.SemanticModel.ItemsServer.BeginUpdateSemanticModelDefinition)

	return fabsemanticmodel.SemanticModel{}
}

func NewRandomSemanticModel() fabsemanticmodel.SemanticModel {
	return fabsemanticmodel.SemanticModel{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		FolderID:    to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fabsemanticmodel.ItemTypeSemanticModel),
	}
}

func NewRandomSemanticModelWithWorkspace(workspaceID string) fabsemanticmodel.SemanticModel {
	result := NewRandomSemanticModel()
	result.WorkspaceID = &workspaceID

	return result
}
