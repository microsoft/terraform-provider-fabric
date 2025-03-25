// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabnotebook "github.com/microsoft/fabric-sdk-go/fabric/notebook"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsNotebook struct{}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsNotebook) CreateDefinition(data fabnotebook.CreateNotebookRequest) *fabnotebook.Definition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsNotebook) TransformDefinition(entity *fabnotebook.Definition) fabnotebook.ItemsClientGetNotebookDefinitionResponse {
	return fabnotebook.ItemsClientGetNotebookDefinitionResponse{
		DefinitionResponse: fabnotebook.DefinitionResponse{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsNotebook) UpdateDefinition(_ *fabnotebook.Definition, data fabnotebook.UpdateNotebookDefinitionRequest) *fabnotebook.Definition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsNotebook) CreateWithParentID(parentID string, data fabnotebook.CreateNotebookRequest) fabnotebook.Notebook {
	entity := NewRandomNotebookWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description

	return entity
}

// Filter implements concreteOperations.
func (o *operationsNotebook) Filter(entities []fabnotebook.Notebook, parentID string) []fabnotebook.Notebook {
	ret := make([]fabnotebook.Notebook, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsNotebook) GetID(entity fabnotebook.Notebook) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsNotebook) TransformCreate(entity fabnotebook.Notebook) fabnotebook.ItemsClientCreateNotebookResponse {
	return fabnotebook.ItemsClientCreateNotebookResponse{
		Notebook: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsNotebook) TransformGet(entity fabnotebook.Notebook) fabnotebook.ItemsClientGetNotebookResponse {
	return fabnotebook.ItemsClientGetNotebookResponse{
		Notebook: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsNotebook) TransformList(entities []fabnotebook.Notebook) fabnotebook.ItemsClientListNotebooksResponse {
	return fabnotebook.ItemsClientListNotebooksResponse{
		Notebooks: fabnotebook.Notebooks{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsNotebook) TransformUpdate(entity fabnotebook.Notebook) fabnotebook.ItemsClientUpdateNotebookResponse {
	return fabnotebook.ItemsClientUpdateNotebookResponse{
		Notebook: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsNotebook) Update(base fabnotebook.Notebook, data fabnotebook.UpdateNotebookRequest) fabnotebook.Notebook {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsNotebook) Validate(newEntity fabnotebook.Notebook, existing []fabnotebook.Notebook) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureNotebook(server *fakeServer) fabnotebook.Notebook {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabnotebook.Notebook,
			fabnotebook.ItemsClientGetNotebookResponse,
			fabnotebook.ItemsClientUpdateNotebookResponse,
			fabnotebook.ItemsClientCreateNotebookResponse,
			fabnotebook.ItemsClientListNotebooksResponse,
			fabnotebook.CreateNotebookRequest,
			fabnotebook.UpdateNotebookRequest]
	}

	type concreteDefinitionOperations interface {
		definitionOperations[
			fabnotebook.Definition,
			fabnotebook.CreateNotebookRequest,
			fabnotebook.UpdateNotebookDefinitionRequest,
			fabnotebook.ItemsClientGetNotebookDefinitionResponse,
			fabnotebook.ItemsClientUpdateNotebookDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsNotebook{}

	var definitionOperations concreteDefinitionOperations = &operationsNotebook{}

	handler := newTypedHandler(server, entityOperations)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.Notebook.ItemsServer.GetNotebook,
		&server.ServerFactory.Notebook.ItemsServer.UpdateNotebook,
		&server.ServerFactory.Notebook.ItemsServer.BeginCreateNotebook,
		&server.ServerFactory.Notebook.ItemsServer.NewListNotebooksPager,
		&server.ServerFactory.Notebook.ItemsServer.DeleteNotebook)

	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.Notebook.ItemsServer.BeginCreateNotebook,
		&server.ServerFactory.Notebook.ItemsServer.BeginGetNotebookDefinition,
		&server.ServerFactory.Notebook.ItemsServer.BeginUpdateNotebookDefinition)

	return fabnotebook.Notebook{}
}

func NewRandomNotebook() fabnotebook.Notebook {
	return fabnotebook.Notebook{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fabnotebook.ItemTypeNotebook),
	}
}

func NewRandomNotebookWithWorkspace(workspaceID string) fabnotebook.Notebook {
	result := NewRandomNotebook()
	result.WorkspaceID = &workspaceID

	return result
}

func NewRandomNotebookDefinition() fabnotebook.Definition {
	defPart := fabnotebook.DefinitionPart{
		PayloadType: to.Ptr(fabnotebook.PayloadTypeInlineBase64),
		Path:        to.Ptr("notebook-content.ipynb"),
		Payload: to.Ptr(
			"eyJjZWxscyI6W3siY2VsbF90eXBlIjoiY29kZSIsIm1ldGFkYXRhIjoge30sInNvdXJjZSI6WyIjIFdlbGNvbWUgdG8geW91ciBub3RlYm9vayJdfV0sIm1ldGFkYXRhIjp7Imxhbmd1YWdlX2luZm8iOnsibmFtZSI6InB5dGhvbiJ9fSwibmJmb3JtYXQiOjQsIm5iZm9ybWF0X21pbm9yIjo1fQ==",
		),
	}

	var defParts []fabnotebook.DefinitionPart

	defParts = append(defParts, defPart)

	return fabnotebook.Definition{
		Format: to.Ptr("ipynb"),
		Parts:  defParts,
	}
}
