// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabvariablelibrary "github.com/microsoft/fabric-sdk-go/fabric/variablelibrary"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsVariableLibrary struct{}

// ConvertItemToEntity implements itemConverter.
func (o *operationsVariableLibrary) ConvertItemToEntity(item fabcore.Item) fabvariablelibrary.VariableLibrary {
	return fabvariablelibrary.VariableLibrary{
		ID:          item.ID,
		DisplayName: item.DisplayName,
		Description: item.Description,
		WorkspaceID: item.WorkspaceID,
		FolderID:    item.FolderID,
		Type:        to.Ptr(fabvariablelibrary.ItemTypeVariableLibrary),
		Properties:  NewRandomVariableLibrary().Properties,
	}
}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsVariableLibrary) CreateDefinition(data fabvariablelibrary.CreateVariableLibraryRequest) *fabvariablelibrary.PublicDefinition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsVariableLibrary) TransformDefinition(entity *fabvariablelibrary.PublicDefinition) fabvariablelibrary.ItemsClientGetVariableLibraryDefinitionResponse {
	return fabvariablelibrary.ItemsClientGetVariableLibraryDefinitionResponse{
		DefinitionResponse: fabvariablelibrary.DefinitionResponse{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsVariableLibrary) UpdateDefinition(_ *fabvariablelibrary.PublicDefinition, data fabvariablelibrary.UpdateVariableLibraryDefinitionRequest) *fabvariablelibrary.PublicDefinition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsVariableLibrary) CreateWithParentID(parentID string, data fabvariablelibrary.CreateVariableLibraryRequest) fabvariablelibrary.VariableLibrary {
	entity := NewRandomVariableLibraryWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description
	entity.FolderID = data.FolderID

	return entity
}

// Filter implements concreteOperations.
func (o *operationsVariableLibrary) Filter(entities []fabvariablelibrary.VariableLibrary, parentID string) []fabvariablelibrary.VariableLibrary {
	ret := make([]fabvariablelibrary.VariableLibrary, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsVariableLibrary) GetID(entity fabvariablelibrary.VariableLibrary) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsVariableLibrary) TransformCreate(entity fabvariablelibrary.VariableLibrary) fabvariablelibrary.ItemsClientCreateVariableLibraryResponse {
	return fabvariablelibrary.ItemsClientCreateVariableLibraryResponse{
		VariableLibrary: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsVariableLibrary) TransformGet(entity fabvariablelibrary.VariableLibrary) fabvariablelibrary.ItemsClientGetVariableLibraryResponse {
	return fabvariablelibrary.ItemsClientGetVariableLibraryResponse{
		VariableLibrary: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsVariableLibrary) TransformList(entities []fabvariablelibrary.VariableLibrary) fabvariablelibrary.ItemsClientListVariableLibrariesResponse {
	return fabvariablelibrary.ItemsClientListVariableLibrariesResponse{
		VariableLibraries: fabvariablelibrary.VariableLibraries{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsVariableLibrary) TransformUpdate(entity fabvariablelibrary.VariableLibrary) fabvariablelibrary.ItemsClientUpdateVariableLibraryResponse {
	return fabvariablelibrary.ItemsClientUpdateVariableLibraryResponse{
		VariableLibrary: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsVariableLibrary) Update(base fabvariablelibrary.VariableLibrary, data fabvariablelibrary.UpdateVariableLibraryRequest) fabvariablelibrary.VariableLibrary {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsVariableLibrary) Validate(newEntity fabvariablelibrary.VariableLibrary, existing []fabvariablelibrary.VariableLibrary) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureVariableLibrary(server *fakeServer) fabvariablelibrary.VariableLibrary {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabvariablelibrary.VariableLibrary,
			fabvariablelibrary.ItemsClientGetVariableLibraryResponse,
			fabvariablelibrary.ItemsClientUpdateVariableLibraryResponse,
			fabvariablelibrary.ItemsClientCreateVariableLibraryResponse,
			fabvariablelibrary.ItemsClientListVariableLibrariesResponse,
			fabvariablelibrary.CreateVariableLibraryRequest,
			fabvariablelibrary.UpdateVariableLibraryRequest]
	}
	type concreteDefinitionOperations interface {
		definitionOperations[
			fabvariablelibrary.PublicDefinition,
			fabvariablelibrary.CreateVariableLibraryRequest,
			fabvariablelibrary.UpdateVariableLibraryDefinitionRequest,
			fabvariablelibrary.ItemsClientGetVariableLibraryDefinitionResponse,
			fabvariablelibrary.ItemsClientUpdateVariableLibraryDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsVariableLibrary{}
	var definitionOperations concreteDefinitionOperations = &operationsVariableLibrary{}

	var converter itemConverter[fabvariablelibrary.VariableLibrary] = &operationsVariableLibrary{}

	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.VariableLibrary.ItemsServer.GetVariableLibrary,
		&server.ServerFactory.VariableLibrary.ItemsServer.UpdateVariableLibrary,
		&server.ServerFactory.VariableLibrary.ItemsServer.BeginCreateVariableLibrary,
		&server.ServerFactory.VariableLibrary.ItemsServer.NewListVariableLibrariesPager,
		&server.ServerFactory.VariableLibrary.ItemsServer.DeleteVariableLibrary)
	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.VariableLibrary.ItemsServer.BeginCreateVariableLibrary,
		&server.ServerFactory.VariableLibrary.ItemsServer.BeginGetVariableLibraryDefinition,
		&server.ServerFactory.VariableLibrary.ItemsServer.BeginUpdateVariableLibraryDefinition)

	return fabvariablelibrary.VariableLibrary{}
}

func NewRandomVariableLibrary() fabvariablelibrary.VariableLibrary {
	return fabvariablelibrary.VariableLibrary{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		FolderID:    to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fabvariablelibrary.ItemTypeVariableLibrary),
		Properties: &fabvariablelibrary.Properties{
			ActiveValueSetName: to.Ptr(testhelp.RandomName()),
		},
	}
}

func NewRandomVariableLibraryWithWorkspace(workspaceID string) fabvariablelibrary.VariableLibrary {
	result := NewRandomVariableLibrary()
	result.WorkspaceID = &workspaceID

	return result
}

func NewRandomVariableLibraryDefinition() fabvariablelibrary.PublicDefinition {
	defPart1 := fabvariablelibrary.PublicDefinitionPart{
		PayloadType: to.Ptr(fabvariablelibrary.PayloadTypeInlineBase64),
		Path:        to.Ptr("variables.json"),
		Payload: to.Ptr(
			"eyJjb250ZW50IjoiSGVsbG8gV29ybGQifQ==",
		),
	}

	defPart2 := fabvariablelibrary.PublicDefinitionPart{
		PayloadType: to.Ptr(fabvariablelibrary.PayloadTypeInlineBase64),
		Path:        to.Ptr("valueSet/valueSet1.json"),
		Payload: to.Ptr(
			"eyJjb250ZW50IjoiSGVsbG8gV29ybGQifQ==",
		),
	}

	defPart3 := fabvariablelibrary.PublicDefinitionPart{
		PayloadType: to.Ptr(fabvariablelibrary.PayloadTypeInlineBase64),
		Path:        to.Ptr("valueSet/valueSet2.json"),
		Payload: to.Ptr(
			"eyJjb250ZW50IjoiSGVsbG8gV29ybGQifQ==",
		),
	}

	defPart4 := fabvariablelibrary.PublicDefinitionPart{
		PayloadType: to.Ptr(fabvariablelibrary.PayloadTypeInlineBase64),
		Path:        to.Ptr("settings.json"),
		Payload: to.Ptr(
			"eyJjb250ZW50IjoiSGVsbG8gV29ybGQifQ==",
		),
	}

	var defParts []fabvariablelibrary.PublicDefinitionPart

	defParts = append(defParts, defPart1)
	defParts = append(defParts, defPart2)
	defParts = append(defParts, defPart3)
	defParts = append(defParts, defPart4)

	return fabvariablelibrary.PublicDefinition{
		Format: to.Ptr("json"),
		Parts:  defParts,
	}
}
