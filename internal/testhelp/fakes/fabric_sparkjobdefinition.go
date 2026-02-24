// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabsparkjobdefinition "github.com/microsoft/fabric-sdk-go/fabric/sparkjobdefinition"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsSparkJobDefinition struct{}

// ConvertItemToEntity implements itemConverter.
func (o *operationsSparkJobDefinition) ConvertItemToEntity(item fabcore.Item) fabsparkjobdefinition.SparkJobDefinition {
	return fabsparkjobdefinition.SparkJobDefinition{
		ID:          item.ID,
		DisplayName: item.DisplayName,
		Description: item.Description,
		WorkspaceID: item.WorkspaceID,
		FolderID:    item.FolderID,
		Type:        to.Ptr(fabsparkjobdefinition.ItemTypeSparkJobDefinition),
		Properties:  NewRandomSparkJobDefinition().Properties,
	}
}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsSparkJobDefinition) CreateDefinition(data fabsparkjobdefinition.CreateSparkJobDefinitionRequest) *fabsparkjobdefinition.PublicDefinition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsSparkJobDefinition) TransformDefinition(entity *fabsparkjobdefinition.PublicDefinition) fabsparkjobdefinition.ItemsClientGetSparkJobDefinitionDefinitionResponse {
	return fabsparkjobdefinition.ItemsClientGetSparkJobDefinitionDefinitionResponse{
		Response: fabsparkjobdefinition.Response{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsSparkJobDefinition) UpdateDefinition(
	_ *fabsparkjobdefinition.PublicDefinition,
	data fabsparkjobdefinition.UpdateSparkJobDefinitionDefinitionRequest,
) *fabsparkjobdefinition.PublicDefinition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsSparkJobDefinition) CreateWithParentID(parentID string, data fabsparkjobdefinition.CreateSparkJobDefinitionRequest) fabsparkjobdefinition.SparkJobDefinition {
	entity := NewRandomSparkJobDefinitionWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description
	entity.FolderID = data.FolderID

	return entity
}

// Filter implements concreteOperations.
func (o *operationsSparkJobDefinition) Filter(entities []fabsparkjobdefinition.SparkJobDefinition, parentID string) []fabsparkjobdefinition.SparkJobDefinition {
	ret := make([]fabsparkjobdefinition.SparkJobDefinition, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsSparkJobDefinition) GetID(entity fabsparkjobdefinition.SparkJobDefinition) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsSparkJobDefinition) TransformCreate(entity fabsparkjobdefinition.SparkJobDefinition) fabsparkjobdefinition.ItemsClientCreateSparkJobDefinitionResponse {
	return fabsparkjobdefinition.ItemsClientCreateSparkJobDefinitionResponse{
		SparkJobDefinition: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsSparkJobDefinition) TransformGet(entity fabsparkjobdefinition.SparkJobDefinition) fabsparkjobdefinition.ItemsClientGetSparkJobDefinitionResponse {
	return fabsparkjobdefinition.ItemsClientGetSparkJobDefinitionResponse{
		SparkJobDefinition: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsSparkJobDefinition) TransformList(entities []fabsparkjobdefinition.SparkJobDefinition) fabsparkjobdefinition.ItemsClientListSparkJobDefinitionsResponse {
	return fabsparkjobdefinition.ItemsClientListSparkJobDefinitionsResponse{
		SparkJobDefinitions: fabsparkjobdefinition.SparkJobDefinitions{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsSparkJobDefinition) TransformUpdate(entity fabsparkjobdefinition.SparkJobDefinition) fabsparkjobdefinition.ItemsClientUpdateSparkJobDefinitionResponse {
	return fabsparkjobdefinition.ItemsClientUpdateSparkJobDefinitionResponse{
		SparkJobDefinition: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsSparkJobDefinition) Update(base fabsparkjobdefinition.SparkJobDefinition, data fabsparkjobdefinition.UpdateSparkJobDefinitionRequest) fabsparkjobdefinition.SparkJobDefinition {
	base.DisplayName = data.DisplayName
	base.Description = data.Description

	return base
}

// Validate implements concreteOperations.
func (o *operationsSparkJobDefinition) Validate(newEntity fabsparkjobdefinition.SparkJobDefinition, existing []fabsparkjobdefinition.SparkJobDefinition) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureSparkJobDefinition(server *fakeServer) fabsparkjobdefinition.SparkJobDefinition {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabsparkjobdefinition.SparkJobDefinition,
			fabsparkjobdefinition.ItemsClientGetSparkJobDefinitionResponse,
			fabsparkjobdefinition.ItemsClientUpdateSparkJobDefinitionResponse,
			fabsparkjobdefinition.ItemsClientCreateSparkJobDefinitionResponse,
			fabsparkjobdefinition.ItemsClientListSparkJobDefinitionsResponse,
			fabsparkjobdefinition.CreateSparkJobDefinitionRequest,
			fabsparkjobdefinition.UpdateSparkJobDefinitionRequest]
	}

	type concreteDefinitionOperations interface {
		definitionOperations[
			fabsparkjobdefinition.PublicDefinition,
			fabsparkjobdefinition.CreateSparkJobDefinitionRequest,
			fabsparkjobdefinition.UpdateSparkJobDefinitionDefinitionRequest,
			fabsparkjobdefinition.ItemsClientGetSparkJobDefinitionDefinitionResponse,
			fabsparkjobdefinition.ItemsClientUpdateSparkJobDefinitionDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsSparkJobDefinition{}
	var converter itemConverter[fabsparkjobdefinition.SparkJobDefinition] = &operationsSparkJobDefinition{}
	var definitionOperations concreteDefinitionOperations = &operationsSparkJobDefinition{}

	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.SparkJobDefinition.ItemsServer.GetSparkJobDefinition,
		&server.ServerFactory.SparkJobDefinition.ItemsServer.UpdateSparkJobDefinition,
		&server.ServerFactory.SparkJobDefinition.ItemsServer.BeginCreateSparkJobDefinition,
		&server.ServerFactory.SparkJobDefinition.ItemsServer.NewListSparkJobDefinitionsPager,
		&server.ServerFactory.SparkJobDefinition.ItemsServer.DeleteSparkJobDefinition)

	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.SparkJobDefinition.ItemsServer.BeginCreateSparkJobDefinition,
		&server.ServerFactory.SparkJobDefinition.ItemsServer.BeginGetSparkJobDefinitionDefinition,
		&server.ServerFactory.SparkJobDefinition.ItemsServer.BeginUpdateSparkJobDefinitionDefinition)

	return fabsparkjobdefinition.SparkJobDefinition{}
}

func NewRandomSparkJobDefinition() fabsparkjobdefinition.SparkJobDefinition {
	return fabsparkjobdefinition.SparkJobDefinition{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		FolderID:    to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fabsparkjobdefinition.ItemTypeSparkJobDefinition),
		Properties: &fabsparkjobdefinition.Properties{
			OneLakeRootPath: to.Ptr(testhelp.RandomURI()),
		},
	}
}

func NewRandomSparkJobDefinitionWithWorkspace(workspaceID string) fabsparkjobdefinition.SparkJobDefinition {
	result := NewRandomSparkJobDefinition()
	result.WorkspaceID = &workspaceID

	return result
}

func NewRandomSparkJobDefinitionDefinition() fabsparkjobdefinition.PublicDefinition {
	defPart := fabsparkjobdefinition.PublicDefinitionPart{
		PayloadType: to.Ptr(fabsparkjobdefinition.PayloadTypeInlineBase64),
		Path:        to.Ptr("SparkJobDefinitionV1.json"),
		Payload: to.Ptr(
			"ew0KICAiZXhlY3V0YWJsZUZpbGUiOiBudWxsLA0KICAiZGVmYXVsdExha2Vob3VzZUFydGlmYWN0SWQiOiBudWxsLA0KICAibWFpbkNsYXNzIjogbnVsbCwNCiAgImFkZGl0aW9uYWxMYWtlaG91c2VJZHMiOiBbXSwNCiAgInJldHJ5UG9saWN5IjogbnVsbCwNCiAgImNvbW1hbmRMaW5lQXJndW1lbnRzIjogbnVsbCwNCiAgImFkZGl0aW9uYWxMaWJyYXJ5VXJpcyI6IG51bGwsDQogICJsYW5ndWFnZSI6IG51bGwsDQogICJlbnZpcm9ubWVudEFydGlmYWN0SWQiOiBudWxsDQp9",
		),
	}

	defParts := make([]fabsparkjobdefinition.PublicDefinitionPart, 0, 1)

	defParts = append(defParts, defPart)

	return fabsparkjobdefinition.PublicDefinition{
		Parts: defParts,
	}
}
