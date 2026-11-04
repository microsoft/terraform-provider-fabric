// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabplan "github.com/microsoft/fabric-sdk-go/fabric/plan"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsPlan struct{}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsPlan) CreateDefinition(data fabplan.CreatePlanRequest) *fabplan.Definition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsPlan) TransformDefinition(entity *fabplan.Definition) fabplan.ItemsClientGetPlanDefinitionResponse {
	return fabplan.ItemsClientGetPlanDefinitionResponse{
		Definition: entity,
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsPlan) UpdateDefinition(_ *fabplan.Definition, data fabplan.UpdatePlanDefinitionRequest) *fabplan.Definition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsPlan) CreateWithParentID(parentID string, data fabplan.CreatePlanRequest) fabplan.Plan {
	entity := NewRandomPlanWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description

	return entity
}

// Filter implements concreteOperations.
func (o *operationsPlan) Filter(entities []fabplan.Plan, parentID string) []fabplan.Plan {
	ret := make([]fabplan.Plan, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsPlan) GetID(entity fabplan.Plan) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsPlan) TransformCreate(entity fabplan.Plan) fabplan.ItemsClientCreatePlanResponse {
	return fabplan.ItemsClientCreatePlanResponse{
		Plan: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsPlan) TransformGet(entity fabplan.Plan) fabplan.ItemsClientGetPlanResponse {
	return fabplan.ItemsClientGetPlanResponse{
		Plan: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsPlan) TransformList(entities []fabplan.Plan) fabplan.ItemsClientListPlansResponse {
	return fabplan.ItemsClientListPlansResponse{
		Value: entities,
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsPlan) TransformUpdate(entity fabplan.Plan) fabplan.ItemsClientUpdatePlanResponse {
	return fabplan.ItemsClientUpdatePlanResponse{
		Plan: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsPlan) Update(base fabplan.Plan, data fabplan.UpdatePlanRequest) fabplan.Plan {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsPlan) Validate(newEntity fabplan.Plan, existing []fabplan.Plan) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configurePlan(server *fakeServer) fabplan.Plan {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabplan.Plan,
			fabplan.ItemsClientGetPlanResponse,
			fabplan.ItemsClientUpdatePlanResponse,
			fabplan.ItemsClientCreatePlanResponse,
			fabplan.ItemsClientListPlansResponse,
			fabplan.CreatePlanRequest,
			fabplan.UpdatePlanRequest]
	}
	type concreteDefinitionOperations interface {
		definitionOperations[
			fabplan.Definition,
			fabplan.CreatePlanRequest,
			fabplan.UpdatePlanDefinitionRequest,
			fabplan.ItemsClientGetPlanDefinitionResponse,
			fabplan.ItemsClientUpdatePlanDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsPlan{}
	var definitionOperations concreteDefinitionOperations = &operationsPlan{}

	handler := newTypedHandler(server, entityOperations)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.Plan.ItemsServer.GetPlan,
		&server.ServerFactory.Plan.ItemsServer.UpdatePlan,
		&server.ServerFactory.Plan.ItemsServer.BeginCreatePlan,
		&server.ServerFactory.Plan.ItemsServer.NewListPlansPager,
		&server.ServerFactory.Plan.ItemsServer.DeletePlan,
	)
	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.Plan.ItemsServer.BeginCreatePlan,
		&server.ServerFactory.Plan.ItemsServer.BeginGetPlanDefinition,
		&server.ServerFactory.Plan.ItemsServer.BeginUpdatePlanDefinition,
	)

	return fabplan.Plan{}
}

func NewRandomPlan() fabplan.Plan {
	return fabplan.Plan{
		ID:          new(testhelp.RandomUUID()),
		DisplayName: new(testhelp.RandomName()),
		Description: new(testhelp.RandomName()),
		WorkspaceID: new(testhelp.RandomUUID()),
		FolderID:    new(testhelp.RandomUUID()),
		Type:        to.Ptr(fabplan.ItemTypePlan),
	}
}

func NewRandomPlanWithWorkspace(workspaceID string) fabplan.Plan {
	result := NewRandomPlan()
	result.WorkspaceID = &workspaceID

	return result
}

func NewRandomPlanDefinition() fabplan.Definition {
	defPart := fabplan.DefinitionPart{
		PayloadType: to.Ptr(fabplan.PayloadTypeInlineBase64),
		Path:        new("definition.json"),
		Payload: new(
			"eyJjb250ZW50IjoiSGVsbG8gV29ybGQifQ==", // {"content":"Hello World"} in base64
		),
	}

	defParts := make([]fabplan.DefinitionPart, 0, 1)

	defParts = append(defParts, defPart)

	return fabplan.Definition{
		Parts: defParts,
	}
}
