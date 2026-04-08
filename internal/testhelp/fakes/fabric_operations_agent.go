// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	faboperationsagent "github.com/microsoft/fabric-sdk-go/fabric/operationsagent"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsOperationsAgent struct{}

// ConvertItemToEntity implements itemConverter.
func (o *operationsOperationsAgent) ConvertItemToEntity(item fabcore.Item) faboperationsagent.OperationsAgent {
	return faboperationsagent.OperationsAgent{
		ID:          item.ID,
		DisplayName: item.DisplayName,
		Description: item.Description,
		WorkspaceID: item.WorkspaceID,
		FolderID:    item.FolderID,
		Type:        to.Ptr(faboperationsagent.ItemTypeOperationsAgent),
		Properties:  NewRandomOperationsAgent().Properties,
	}
}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsOperationsAgent) CreateDefinition(data faboperationsagent.CreateOperationsAgentRequest) *faboperationsagent.PublicDefinition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsOperationsAgent) TransformDefinition(entity *faboperationsagent.PublicDefinition) faboperationsagent.ItemsClientGetOperationsAgentDefinitionResponse {
	return faboperationsagent.ItemsClientGetOperationsAgentDefinitionResponse{
		DefinitionResponse: faboperationsagent.DefinitionResponse{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsOperationsAgent) UpdateDefinition(_ *faboperationsagent.PublicDefinition, data faboperationsagent.UpdateOperationsAgentDefinitionRequest) *faboperationsagent.PublicDefinition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsOperationsAgent) CreateWithParentID(parentID string, data faboperationsagent.CreateOperationsAgentRequest) faboperationsagent.OperationsAgent {
	entity := NewRandomOperationsAgentWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description
	entity.FolderID = data.FolderID

	return entity
}

// Filter implements concreteOperations.
func (o *operationsOperationsAgent) Filter(entities []faboperationsagent.OperationsAgent, parentID string) []faboperationsagent.OperationsAgent {
	ret := make([]faboperationsagent.OperationsAgent, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsOperationsAgent) GetID(entity faboperationsagent.OperationsAgent) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsOperationsAgent) TransformCreate(entity faboperationsagent.OperationsAgent) faboperationsagent.ItemsClientCreateOperationsAgentResponse {
	return faboperationsagent.ItemsClientCreateOperationsAgentResponse{
		OperationsAgent: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsOperationsAgent) TransformGet(entity faboperationsagent.OperationsAgent) faboperationsagent.ItemsClientGetOperationsAgentResponse {
	return faboperationsagent.ItemsClientGetOperationsAgentResponse{
		OperationsAgent: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsOperationsAgent) TransformList(entities []faboperationsagent.OperationsAgent) faboperationsagent.ItemsClientListOperationsAgentsResponse {
	return faboperationsagent.ItemsClientListOperationsAgentsResponse{
		OperationsAgents: faboperationsagent.OperationsAgents{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsOperationsAgent) TransformUpdate(entity faboperationsagent.OperationsAgent) faboperationsagent.ItemsClientUpdateOperationsAgentResponse {
	return faboperationsagent.ItemsClientUpdateOperationsAgentResponse{
		OperationsAgent: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsOperationsAgent) Update(base faboperationsagent.OperationsAgent, data faboperationsagent.UpdateOperationsAgentRequest) faboperationsagent.OperationsAgent {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsOperationsAgent) Validate(newEntity faboperationsagent.OperationsAgent, existing []faboperationsagent.OperationsAgent) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureOperationsAgent(server *fakeServer) faboperationsagent.OperationsAgent {
	type concreteEntityOperations interface {
		parentIDOperations[
			faboperationsagent.OperationsAgent,
			faboperationsagent.ItemsClientGetOperationsAgentResponse,
			faboperationsagent.ItemsClientUpdateOperationsAgentResponse,
			faboperationsagent.ItemsClientCreateOperationsAgentResponse,
			faboperationsagent.ItemsClientListOperationsAgentsResponse,
			faboperationsagent.CreateOperationsAgentRequest,
			faboperationsagent.UpdateOperationsAgentRequest]
	}
	type concreteDefinitionOperations interface {
		definitionOperations[
			faboperationsagent.PublicDefinition,
			faboperationsagent.CreateOperationsAgentRequest,
			faboperationsagent.UpdateOperationsAgentDefinitionRequest,
			faboperationsagent.ItemsClientGetOperationsAgentDefinitionResponse,
			faboperationsagent.ItemsClientUpdateOperationsAgentDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsOperationsAgent{}
	var definitionOperations concreteDefinitionOperations = &operationsOperationsAgent{}
	var converter itemConverter[faboperationsagent.OperationsAgent] = &operationsOperationsAgent{}

	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.OperationsAgent.ItemsServer.GetOperationsAgent,
		&server.ServerFactory.OperationsAgent.ItemsServer.UpdateOperationsAgent,
		&server.ServerFactory.OperationsAgent.ItemsServer.BeginCreateOperationsAgent,
		&server.ServerFactory.OperationsAgent.ItemsServer.NewListOperationsAgentsPager,
		&server.ServerFactory.OperationsAgent.ItemsServer.DeleteOperationsAgent)
	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.OperationsAgent.ItemsServer.BeginCreateOperationsAgent,
		&server.ServerFactory.OperationsAgent.ItemsServer.BeginGetOperationsAgentDefinition,
		&server.ServerFactory.OperationsAgent.ItemsServer.BeginUpdateOperationsAgentDefinition)

	return faboperationsagent.OperationsAgent{}
}

func NewRandomOperationsAgent() faboperationsagent.OperationsAgent {
	return faboperationsagent.OperationsAgent{
		ID:          new(testhelp.RandomUUID()),
		DisplayName: new(testhelp.RandomName()),
		Description: new(testhelp.RandomName()),
		WorkspaceID: new(testhelp.RandomUUID()),
		FolderID:    new(testhelp.RandomUUID()),
		Type:        to.Ptr(faboperationsagent.ItemTypeOperationsAgent),
		Properties: &faboperationsagent.Properties{
			State: new(faboperationsagent.AgentStateActive),
		},
	}
}

func NewRandomOperationsAgentWithWorkspace(workspaceID string) faboperationsagent.OperationsAgent {
	result := NewRandomOperationsAgent()
	result.WorkspaceID = &workspaceID

	return result
}

func NewRandomOperationsAgentDefinition() faboperationsagent.PublicDefinition {
	defPart := faboperationsagent.PublicDefinitionPart{
		PayloadType: new(faboperationsagent.PayloadTypeInlineBase64),
		Path:        new("OperationsAgentV1.json"),
		Payload: new(
			"eyJjb250ZW50IjoiSGVsbG8gV29ybGQifQ==", // {"content":"Hello World"} in base64
		),
	}

	defParts := make([]faboperationsagent.PublicDefinitionPart, 0, 1)

	defParts = append(defParts, defPart)

	return faboperationsagent.PublicDefinition{
		Format: new(faboperationsagent.DefinitionFormatOperationsAgentV1),
		Parts:  defParts,
	}
}
