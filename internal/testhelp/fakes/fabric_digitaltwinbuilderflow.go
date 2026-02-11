// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabdigitaltwinbuilderflow "github.com/microsoft/fabric-sdk-go/fabric/digitaltwinbuilderflow"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsDigitalTwinBuilderFlow struct{}

// ConvertItemToEntity implements itemConverter.
func (o *operationsDigitalTwinBuilderFlow) ConvertItemToEntity(item fabcore.Item) fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow {
	return fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow{
		ID:          item.ID,
		DisplayName: item.DisplayName,
		Description: item.Description,
		WorkspaceID: item.WorkspaceID,
		FolderID:    item.FolderID,
		Type:        to.Ptr(fabdigitaltwinbuilderflow.ItemTypeDigitalTwinBuilderFlow),
		Properties:  NewRandomDigitalTwinBuilderFlow().Properties,
	}
}

// CreateWithParentID implements concreteOperations.
func (o *operationsDigitalTwinBuilderFlow) CreateWithParentID(parentID string, data fabdigitaltwinbuilderflow.CreateDigitalTwinBuilderFlowRequest) fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow {
	entity := NewRandomDigitalTwinBuilderFlowWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description

	return entity
}

// Filter implements concreteOperations.
func (o *operationsDigitalTwinBuilderFlow) Filter(entities []fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow, parentID string) []fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow {
	ret := make([]fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsDigitalTwinBuilderFlow) CreateDefinition(data fabdigitaltwinbuilderflow.CreateDigitalTwinBuilderFlowRequest) *fabdigitaltwinbuilderflow.PublicDefinition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsDigitalTwinBuilderFlow) TransformDefinition(entity *fabdigitaltwinbuilderflow.PublicDefinition) fabdigitaltwinbuilderflow.ItemsClientGetDigitalTwinBuilderFlowDefinitionResponse {
	return fabdigitaltwinbuilderflow.ItemsClientGetDigitalTwinBuilderFlowDefinitionResponse{
		DefinitionResponse: fabdigitaltwinbuilderflow.DefinitionResponse{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsDigitalTwinBuilderFlow) UpdateDefinition(
	_ *fabdigitaltwinbuilderflow.PublicDefinition,
	data fabdigitaltwinbuilderflow.UpdateDigitalTwinBuilderFlowDefinitionRequest,
) *fabdigitaltwinbuilderflow.PublicDefinition {
	return data.Definition
}

// GetID implements concreteOperations.
func (o *operationsDigitalTwinBuilderFlow) GetID(entity fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsDigitalTwinBuilderFlow) TransformCreate(entity fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow) fabdigitaltwinbuilderflow.ItemsClientCreateDigitalTwinBuilderFlowResponse {
	return fabdigitaltwinbuilderflow.ItemsClientCreateDigitalTwinBuilderFlowResponse{
		DigitalTwinBuilderFlow: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsDigitalTwinBuilderFlow) TransformGet(entity fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow) fabdigitaltwinbuilderflow.ItemsClientGetDigitalTwinBuilderFlowResponse {
	return fabdigitaltwinbuilderflow.ItemsClientGetDigitalTwinBuilderFlowResponse{
		DigitalTwinBuilderFlow: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsDigitalTwinBuilderFlow) TransformList(entities []fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow) fabdigitaltwinbuilderflow.ItemsClientListDigitalTwinBuilderFlowsResponse {
	return fabdigitaltwinbuilderflow.ItemsClientListDigitalTwinBuilderFlowsResponse{
		DigitalTwinBuilderFlows: fabdigitaltwinbuilderflow.DigitalTwinBuilderFlows{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsDigitalTwinBuilderFlow) TransformUpdate(entity fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow) fabdigitaltwinbuilderflow.ItemsClientUpdateDigitalTwinBuilderFlowResponse {
	return fabdigitaltwinbuilderflow.ItemsClientUpdateDigitalTwinBuilderFlowResponse{
		DigitalTwinBuilderFlow: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsDigitalTwinBuilderFlow) Update(
	base fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow,
	data fabdigitaltwinbuilderflow.UpdateDigitalTwinBuilderFlowRequest,
) fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsDigitalTwinBuilderFlow) Validate(newEntity fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow, existing []fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureDigitalTwinBuilderFlow(server *fakeServer) fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow,
			fabdigitaltwinbuilderflow.ItemsClientGetDigitalTwinBuilderFlowResponse,
			fabdigitaltwinbuilderflow.ItemsClientUpdateDigitalTwinBuilderFlowResponse,
			fabdigitaltwinbuilderflow.ItemsClientCreateDigitalTwinBuilderFlowResponse,
			fabdigitaltwinbuilderflow.ItemsClientListDigitalTwinBuilderFlowsResponse,
			fabdigitaltwinbuilderflow.CreateDigitalTwinBuilderFlowRequest,
			fabdigitaltwinbuilderflow.UpdateDigitalTwinBuilderFlowRequest]
	}
	type concreteDefinitionOperations interface {
		definitionOperations[
			fabdigitaltwinbuilderflow.PublicDefinition,
			fabdigitaltwinbuilderflow.CreateDigitalTwinBuilderFlowRequest,
			fabdigitaltwinbuilderflow.UpdateDigitalTwinBuilderFlowDefinitionRequest,
			fabdigitaltwinbuilderflow.ItemsClientGetDigitalTwinBuilderFlowDefinitionResponse,
			fabdigitaltwinbuilderflow.ItemsClientUpdateDigitalTwinBuilderFlowDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsDigitalTwinBuilderFlow{}
	var definitionOperations concreteDefinitionOperations = &operationsDigitalTwinBuilderFlow{}
	var converter itemConverter[fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow] = &operationsDigitalTwinBuilderFlow{}

	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.DigitalTwinBuilderFlow.ItemsServer.GetDigitalTwinBuilderFlow,
		&server.ServerFactory.DigitalTwinBuilderFlow.ItemsServer.UpdateDigitalTwinBuilderFlow,
		&server.ServerFactory.DigitalTwinBuilderFlow.ItemsServer.BeginCreateDigitalTwinBuilderFlow,
		&server.ServerFactory.DigitalTwinBuilderFlow.ItemsServer.NewListDigitalTwinBuilderFlowsPager,
		&server.ServerFactory.DigitalTwinBuilderFlow.ItemsServer.DeleteDigitalTwinBuilderFlow)
	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.DigitalTwinBuilderFlow.ItemsServer.BeginCreateDigitalTwinBuilderFlow,
		&server.ServerFactory.DigitalTwinBuilderFlow.ItemsServer.BeginGetDigitalTwinBuilderFlowDefinition,
		&server.ServerFactory.DigitalTwinBuilderFlow.ItemsServer.BeginUpdateDigitalTwinBuilderFlowDefinition)

	return fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow{}
}

func NewRandomDigitalTwinBuilderFlow() fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow {
	return fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		FolderID:    to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fabdigitaltwinbuilderflow.ItemTypeDigitalTwinBuilderFlow),
		Properties: &fabdigitaltwinbuilderflow.Properties{
			DigitalTwinBuilderItemReference: &fabdigitaltwinbuilderflow.ItemReferenceByID{
				ItemID:        to.Ptr(testhelp.RandomUUID()),
				ReferenceType: to.Ptr(fabdigitaltwinbuilderflow.ItemReferenceTypeByID),
				WorkspaceID:   to.Ptr(testhelp.RandomUUID()),
			},
		},
	}
}

func NewRandomDigitalTwinBuilderFlowWithWorkspace(workspaceID string) fabdigitaltwinbuilderflow.DigitalTwinBuilderFlow {
	result := NewRandomDigitalTwinBuilderFlow()
	result.WorkspaceID = &workspaceID

	return result
}

func NewRandomDigitalTwinBuilderFlowDefinition() fabdigitaltwinbuilderflow.PublicDefinition {
	return fabdigitaltwinbuilderflow.PublicDefinition{
		Parts: []fabdigitaltwinbuilderflow.PublicDefinitionPart{
			{
				PayloadType: to.Ptr(fabdigitaltwinbuilderflow.PayloadTypeInlineBase64),
				Path:        to.Ptr("definition.json"),
				Payload: to.Ptr(
					"eyAKICAiRGlnaXRhbFR3aW5CdWlsZGVySWQiOiAiNTZhMGU2Y2EtMTAxZS1iYzA1LTQ2NDktNjAzOTMzYWUxMjcwIiwgCiAgIk9wZXJhdGlvbklkcyI6IFsgCiAgICAiY2U5ZDBlZjktZDhmNi00MzkxLTllMzctOGJkYjkxYjFmYzE2IiAKICBdLCAKICAiSXNPbkRlbWFuZCI6IGZhbHNlIAp9IA==",
				),
			},
		},
	}
}
