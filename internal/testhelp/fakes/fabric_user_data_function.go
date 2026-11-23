// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabuserdatafunction "github.com/microsoft/fabric-sdk-go/fabric/userdatafunction"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsUserDataFunction struct{}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsUserDataFunction) CreateDefinition(data fabuserdatafunction.CreateUserDataFunctionRequest) *fabuserdatafunction.PublicDefinition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsUserDataFunction) TransformDefinition(entity *fabuserdatafunction.PublicDefinition) fabuserdatafunction.ItemsClientGetUserDataFunctionDefinitionResponse {
	return fabuserdatafunction.ItemsClientGetUserDataFunctionDefinitionResponse{
		Definition: entity,
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsUserDataFunction) UpdateDefinition(_ *fabuserdatafunction.PublicDefinition, data fabuserdatafunction.UpdateUserDataFunctionDefinitionRequest) *fabuserdatafunction.PublicDefinition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsUserDataFunction) CreateWithParentID(parentID string, data fabuserdatafunction.CreateUserDataFunctionRequest) fabuserdatafunction.UserDataFunction {
	entity := NewRandomUserDataFunctionWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description

	return entity
}

// Filter implements concreteOperations.
func (o *operationsUserDataFunction) Filter(entities []fabuserdatafunction.UserDataFunction, parentID string) []fabuserdatafunction.UserDataFunction {
	ret := make([]fabuserdatafunction.UserDataFunction, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsUserDataFunction) GetID(entity fabuserdatafunction.UserDataFunction) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsUserDataFunction) TransformCreate(entity fabuserdatafunction.UserDataFunction) fabuserdatafunction.ItemsClientCreateUserDataFunctionResponse {
	return fabuserdatafunction.ItemsClientCreateUserDataFunctionResponse{
		UserDataFunction: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsUserDataFunction) TransformGet(entity fabuserdatafunction.UserDataFunction) fabuserdatafunction.ItemsClientGetUserDataFunctionResponse {
	return fabuserdatafunction.ItemsClientGetUserDataFunctionResponse{
		UserDataFunction: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsUserDataFunction) TransformList(entities []fabuserdatafunction.UserDataFunction) fabuserdatafunction.ItemsClientListUserDataFunctionsResponse {
	return fabuserdatafunction.ItemsClientListUserDataFunctionsResponse{
		Value: entities,
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsUserDataFunction) TransformUpdate(entity fabuserdatafunction.UserDataFunction) fabuserdatafunction.ItemsClientUpdateUserDataFunctionResponse {
	return fabuserdatafunction.ItemsClientUpdateUserDataFunctionResponse{
		UserDataFunction: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsUserDataFunction) Update(base fabuserdatafunction.UserDataFunction, data fabuserdatafunction.UpdateUserDataFunctionRequest) fabuserdatafunction.UserDataFunction {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsUserDataFunction) Validate(newEntity fabuserdatafunction.UserDataFunction, existing []fabuserdatafunction.UserDataFunction) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureUserDataFunction(server *fakeServer) fabuserdatafunction.UserDataFunction {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabuserdatafunction.UserDataFunction,
			fabuserdatafunction.ItemsClientGetUserDataFunctionResponse,
			fabuserdatafunction.ItemsClientUpdateUserDataFunctionResponse,
			fabuserdatafunction.ItemsClientCreateUserDataFunctionResponse,
			fabuserdatafunction.ItemsClientListUserDataFunctionsResponse,
			fabuserdatafunction.CreateUserDataFunctionRequest,
			fabuserdatafunction.UpdateUserDataFunctionRequest]
	}

	type concreteDefinitionOperations interface {
		definitionOperations[
			fabuserdatafunction.PublicDefinition,
			fabuserdatafunction.CreateUserDataFunctionRequest,
			fabuserdatafunction.UpdateUserDataFunctionDefinitionRequest,
			fabuserdatafunction.ItemsClientGetUserDataFunctionDefinitionResponse,
			fabuserdatafunction.ItemsClientUpdateUserDataFunctionDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsUserDataFunction{}

	var definitionOperations concreteDefinitionOperations = &operationsUserDataFunction{}

	handler := newTypedHandler(server, entityOperations)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.UserDataFunction.ItemsServer.GetUserDataFunction,
		&server.ServerFactory.UserDataFunction.ItemsServer.UpdateUserDataFunction,
		&server.ServerFactory.UserDataFunction.ItemsServer.BeginCreateUserDataFunction,
		&server.ServerFactory.UserDataFunction.ItemsServer.NewListUserDataFunctionsPager,
		&server.ServerFactory.UserDataFunction.ItemsServer.DeleteUserDataFunction,
	)

	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.UserDataFunction.ItemsServer.BeginCreateUserDataFunction,
		&server.ServerFactory.UserDataFunction.ItemsServer.BeginGetUserDataFunctionDefinition,
		&server.ServerFactory.UserDataFunction.ItemsServer.BeginUpdateUserDataFunctionDefinition,
	)

	return fabuserdatafunction.UserDataFunction{}
}

func NewRandomUserDataFunction() fabuserdatafunction.UserDataFunction {
	return fabuserdatafunction.UserDataFunction{
		ID:          new(testhelp.RandomUUID()),
		DisplayName: new(testhelp.RandomName()),
		Description: new(testhelp.RandomName()),
		WorkspaceID: new(testhelp.RandomUUID()),
		FolderID:    new(testhelp.RandomUUID()),
		Type:        to.Ptr(fabuserdatafunction.ItemTypeUserDataFunction),
	}
}

func NewRandomUserDataFunctionWithWorkspace(workspaceID string) fabuserdatafunction.UserDataFunction {
	result := NewRandomUserDataFunction()
	result.WorkspaceID = &workspaceID

	return result
}

func NewRandomUserDataFunctionDefinition() fabuserdatafunction.PublicDefinition {
	defPart1 := fabuserdatafunction.PublicDefinitionPart{
		PayloadType: to.Ptr(fabuserdatafunction.PayloadTypeInlineBase64),
		Path:        new("definition.json"),
		Payload: new(
			"eyJydW50aW1lIjoiUFlUSE9OIiwiY29ubmVjdGVkRGF0YVNvdXJjZXMiOltdLCJmdW5jdGlvbnMiOltdLCJsaWJyYXJpZXMiOnsicHVibGljIjpbXSwicHJpdmF0ZSI6W119fQ==", // {"runtime":"PYTHON","connectedDataSources":[],"functions":[],"libraries":{"public":[],"private":[]}} in base64
		),
	}

	defPart2 := fabuserdatafunction.PublicDefinitionPart{
		PayloadType: to.Ptr(fabuserdatafunction.PayloadTypeInlineBase64),
		Path:        new("function_app.py"),
		Payload: new(
			"aW1wb3J0IGZhYnJpYy5mdW5jdGlvbnMgYXMgZm4KdWRmID0gZm4uVXNlckRhdGFGdW5jdGlvbnMoKQo=", // import fabric.functions as fn\nudf = fn.UserDataFunctions()\n in base64
		),
	}

	defPart3 := fabuserdatafunction.PublicDefinitionPart{
		PayloadType: to.Ptr(fabuserdatafunction.PayloadTypeInlineBase64),
		Path:        new("resources/functions.json"),
		Payload: new(
			"e30=", // {} in base64
		),
	}

	defParts := make([]fabuserdatafunction.PublicDefinitionPart, 0, 3)

	defParts = append(defParts, defPart1)
	defParts = append(defParts, defPart2)
	defParts = append(defParts, defPart3)

	return fabuserdatafunction.PublicDefinition{
		Parts: defParts,
	}
}
