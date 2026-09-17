// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabsnowflakedatabase "github.com/microsoft/fabric-sdk-go/fabric/snowflakedatabase"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsSnowflakeDatabase struct{}

// ConvertItemToEntity implements itemConverter.
func (o *operationsSnowflakeDatabase) ConvertItemToEntity(item fabcore.Item) fabsnowflakedatabase.SnowflakeDatabase {
	return fabsnowflakedatabase.SnowflakeDatabase{
		ID:          item.ID,
		DisplayName: item.DisplayName,
		Description: item.Description,
		WorkspaceID: item.WorkspaceID,
		FolderID:    item.FolderID,
		Type:        to.Ptr(fabsnowflakedatabase.ItemTypeSnowflakeDatabase),
		Properties:  NewRandomSnowflakeDatabase().Properties,
	}
}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsSnowflakeDatabase) CreateDefinition(data fabsnowflakedatabase.CreateSnowflakeDatabaseRequest) *fabsnowflakedatabase.Definition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsSnowflakeDatabase) TransformDefinition(entity *fabsnowflakedatabase.Definition) fabsnowflakedatabase.ItemsClientGetSnowflakeDatabaseDefinitionResponse {
	return fabsnowflakedatabase.ItemsClientGetSnowflakeDatabaseDefinitionResponse{
		DefinitionResponse: fabsnowflakedatabase.DefinitionResponse{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsSnowflakeDatabase) UpdateDefinition(_ *fabsnowflakedatabase.Definition, data fabsnowflakedatabase.UpdateSnowflakeDatabaseDefinitionRequest) *fabsnowflakedatabase.Definition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsSnowflakeDatabase) CreateWithParentID(parentID string, data fabsnowflakedatabase.CreateSnowflakeDatabaseRequest) fabsnowflakedatabase.SnowflakeDatabase {
	entity := NewRandomSnowflakeDatabaseWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description

	return entity
}

// Filter implements concreteOperations.
func (o *operationsSnowflakeDatabase) Filter(entities []fabsnowflakedatabase.SnowflakeDatabase, parentID string) []fabsnowflakedatabase.SnowflakeDatabase {
	ret := make([]fabsnowflakedatabase.SnowflakeDatabase, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsSnowflakeDatabase) GetID(entity fabsnowflakedatabase.SnowflakeDatabase) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsSnowflakeDatabase) TransformCreate(entity fabsnowflakedatabase.SnowflakeDatabase) fabsnowflakedatabase.ItemsClientCreateSnowflakeDatabaseResponse {
	return fabsnowflakedatabase.ItemsClientCreateSnowflakeDatabaseResponse{
		SnowflakeDatabase: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsSnowflakeDatabase) TransformGet(entity fabsnowflakedatabase.SnowflakeDatabase) fabsnowflakedatabase.ItemsClientGetSnowflakeDatabaseResponse {
	return fabsnowflakedatabase.ItemsClientGetSnowflakeDatabaseResponse{
		SnowflakeDatabase: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsSnowflakeDatabase) TransformList(entities []fabsnowflakedatabase.SnowflakeDatabase) fabsnowflakedatabase.ItemsClientListSnowflakeDatabasesResponse {
	return fabsnowflakedatabase.ItemsClientListSnowflakeDatabasesResponse{
		SnowflakeDatabases: fabsnowflakedatabase.SnowflakeDatabases{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsSnowflakeDatabase) TransformUpdate(entity fabsnowflakedatabase.SnowflakeDatabase) fabsnowflakedatabase.ItemsClientUpdateSnowflakeDatabaseResponse {
	return fabsnowflakedatabase.ItemsClientUpdateSnowflakeDatabaseResponse{
		SnowflakeDatabase: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsSnowflakeDatabase) Update(base fabsnowflakedatabase.SnowflakeDatabase, data fabsnowflakedatabase.UpdateSnowflakeDatabaseRequest) fabsnowflakedatabase.SnowflakeDatabase {
	base.DisplayName = data.DisplayName
	base.Description = data.Description

	return base
}

// Validate implements concreteOperations.
func (o *operationsSnowflakeDatabase) Validate(newEntity fabsnowflakedatabase.SnowflakeDatabase, existing []fabsnowflakedatabase.SnowflakeDatabase) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureSnowflakeDatabase(server *fakeServer) fabsnowflakedatabase.SnowflakeDatabase {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabsnowflakedatabase.SnowflakeDatabase,
			fabsnowflakedatabase.ItemsClientGetSnowflakeDatabaseResponse,
			fabsnowflakedatabase.ItemsClientUpdateSnowflakeDatabaseResponse,
			fabsnowflakedatabase.ItemsClientCreateSnowflakeDatabaseResponse,
			fabsnowflakedatabase.ItemsClientListSnowflakeDatabasesResponse,
			fabsnowflakedatabase.CreateSnowflakeDatabaseRequest,
			fabsnowflakedatabase.UpdateSnowflakeDatabaseRequest]
	}

	type concreteDefinitionOperations interface {
		definitionOperations[
			fabsnowflakedatabase.Definition,
			fabsnowflakedatabase.CreateSnowflakeDatabaseRequest,
			fabsnowflakedatabase.UpdateSnowflakeDatabaseDefinitionRequest,
			fabsnowflakedatabase.ItemsClientGetSnowflakeDatabaseDefinitionResponse,
			fabsnowflakedatabase.ItemsClientUpdateSnowflakeDatabaseDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsSnowflakeDatabase{}

	var definitionOperations concreteDefinitionOperations = &operationsSnowflakeDatabase{}

	var converter itemConverter[fabsnowflakedatabase.SnowflakeDatabase] = &operationsSnowflakeDatabase{}

	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.SnowflakeDatabase.ItemsServer.GetSnowflakeDatabase,
		&server.ServerFactory.SnowflakeDatabase.ItemsServer.UpdateSnowflakeDatabase,
		&server.ServerFactory.SnowflakeDatabase.ItemsServer.BeginCreateSnowflakeDatabase,
		&server.ServerFactory.SnowflakeDatabase.ItemsServer.NewListSnowflakeDatabasesPager,
		&server.ServerFactory.SnowflakeDatabase.ItemsServer.DeleteSnowflakeDatabase)

	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.SnowflakeDatabase.ItemsServer.BeginCreateSnowflakeDatabase,
		&server.ServerFactory.SnowflakeDatabase.ItemsServer.BeginGetSnowflakeDatabaseDefinition,
		&server.ServerFactory.SnowflakeDatabase.ItemsServer.BeginUpdateSnowflakeDatabaseDefinition)

	return fabsnowflakedatabase.SnowflakeDatabase{}
}

func NewRandomSnowflakeDatabase() fabsnowflakedatabase.SnowflakeDatabase {
	return fabsnowflakedatabase.SnowflakeDatabase{
		ID:          new(testhelp.RandomUUID()),
		DisplayName: new(testhelp.RandomName()),
		Description: new(testhelp.RandomName()),
		WorkspaceID: new(testhelp.RandomUUID()),
		FolderID:    new(testhelp.RandomUUID()),
		Type:        to.Ptr(fabsnowflakedatabase.ItemTypeSnowflakeDatabase),
		Properties: &fabsnowflakedatabase.Properties{
			ConnectionID:          new(testhelp.RandomUUID()),
			DefaultSchema:         new("PUBLIC"),
			OnelakeTablesPath:     new(testhelp.RandomName()),
			SnowflakeAccountURL:   new(testhelp.RandomURI()),
			SnowflakeDatabaseName: new(testhelp.RandomName()),
			SnowflakeVolumePath:   new(testhelp.RandomName()),
			SQLEndpointProperties: &fabsnowflakedatabase.SQLEndpointProperties{
				ID:                 new(testhelp.RandomUUID()),
				ProvisioningStatus: to.Ptr(fabsnowflakedatabase.ProvisioningStatusSuccess),
				ConnectionString:   new(testhelp.RandomURI()),
			},
		},
	}
}

func NewRandomSnowflakeDatabaseWithWorkspace(workspaceID string) fabsnowflakedatabase.SnowflakeDatabase {
	result := NewRandomSnowflakeDatabase()
	result.WorkspaceID = &workspaceID

	return result
}

func NewRandomSnowflakeDatabaseDefinition() fabsnowflakedatabase.Definition {
	defPart := fabsnowflakedatabase.DefinitionPart{
		PayloadType: to.Ptr(fabsnowflakedatabase.PayloadTypeInlineBase64),
		Path:        new("SnowflakeDatabaseProperties.json"),
		Payload:     new("e30="),
	}

	defParts := make([]fabsnowflakedatabase.DefinitionPart, 0, 1)

	defParts = append(defParts, defPart)

	return fabsnowflakedatabase.Definition{
		Parts: defParts,
	}
}
