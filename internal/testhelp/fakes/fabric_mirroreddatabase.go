// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	"github.com/microsoft/fabric-sdk-go/fabric/mirroreddatabase"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// operationsMirroredDatabase implements the operations for MirroredDatabase
// and satisfies the concreteEntityOperations and concreteDefinitionOperations interfaces.
type operationsMirroredDatabase struct{}

// ConvertItemToEntity implements itemConverter. It converts a generic fabcore.Item into a
// mirroreddatabase.MirroredDatabase.
func (o *operationsMirroredDatabase) ConvertItemToEntity(item fabcore.Item) mirroreddatabase.MirroredDatabase {
	return mirroreddatabase.MirroredDatabase{
		ID:          item.ID,
		DisplayName: item.DisplayName,
		Description: item.Description,
		WorkspaceID: item.WorkspaceID,
		// Set the type to mirrored database.
		Type:       to.Ptr(mirroreddatabase.ItemTypeMirroredDatabase),
		Properties: NewRandomMirroredDatabase().Properties,
	}
}

// CreateDefinition implements concreteDefinitionOperations.
// It simply returns the Definition provided in the Create request.
func (o *operationsMirroredDatabase) CreateDefinition(data mirroreddatabase.CreateMirroredDatabaseRequest) *mirroreddatabase.Definition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
// It wraps the public definition into the GetDefinition response.
func (o *operationsMirroredDatabase) TransformDefinition(entity *mirroreddatabase.Definition) mirroreddatabase.ItemsClientGetMirroredDatabaseDefinitionResponse {
	return mirroreddatabase.ItemsClientGetMirroredDatabaseDefinitionResponse{
		DefinitionResponse: mirroreddatabase.DefinitionResponse{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
// It returns the updated definition as provided in the Update request.
func (o *operationsMirroredDatabase) UpdateDefinition(_ *mirroreddatabase.Definition, data mirroreddatabase.UpdateMirroredDatabaseDefinitionRequest) *mirroreddatabase.Definition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
// It creates a mirrored database entity with the given parent workspace ID.
func (o *operationsMirroredDatabase) CreateWithParentID(parentID string, data mirroreddatabase.CreateMirroredDatabaseRequest) mirroreddatabase.MirroredDatabase {
	entity := NewRandomMirroredDatabaseWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description
	return entity
}

// Filter implements concreteOperations, returning only entities whose WorkspaceID matches parentID.
func (o *operationsMirroredDatabase) Filter(entities []mirroreddatabase.MirroredDatabase, parentID string) []mirroreddatabase.MirroredDatabase {
	var ret []mirroreddatabase.MirroredDatabase
	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}
	return ret
}

// GetID implements concreteOperations.
// It returns a generated unique ID composed of the workspace ID and the entity ID.
func (o *operationsMirroredDatabase) GetID(entity mirroreddatabase.MirroredDatabase) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations. It wraps the created entity in the Create response.
func (o *operationsMirroredDatabase) TransformCreate(entity mirroreddatabase.MirroredDatabase) mirroreddatabase.ItemsClientCreateMirroredDatabaseResponse {
	return mirroreddatabase.ItemsClientCreateMirroredDatabaseResponse{
		MirroredDatabase: entity,
	}
}

// TransformGet implements concreteOperations.
// It returns the get response for the given mirrored database entity.
func (o *operationsMirroredDatabase) TransformGet(entity mirroreddatabase.MirroredDatabase) mirroreddatabase.ItemsClientGetMirroredDatabaseResponse {
	return mirroreddatabase.ItemsClientGetMirroredDatabaseResponse{
		MirroredDatabase: entity,
	}
}

// TransformList implements concreteOperations.
// It wraps a list of entities into the List response.
func (o *operationsMirroredDatabase) TransformList(entities []mirroreddatabase.MirroredDatabase) mirroreddatabase.ItemsClientListMirroredDatabasesResponse {
	return mirroreddatabase.ItemsClientListMirroredDatabasesResponse{
		MirroredDatabases: mirroreddatabase.MirroredDatabases{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
// It wraps the updated entity into the Update response.
func (o *operationsMirroredDatabase) TransformUpdate(entity mirroreddatabase.MirroredDatabase) mirroreddatabase.ItemsClientUpdateMirroredDatabaseResponse {
	return mirroreddatabase.ItemsClientUpdateMirroredDatabaseResponse{
		MirroredDatabase: entity,
	}
}

// Update implements concreteOperations.
// It updates the base entity with the new data.
func (o *operationsMirroredDatabase) Update(base mirroreddatabase.MirroredDatabase, data mirroreddatabase.UpdateMirroredDatabaseRequest) mirroreddatabase.MirroredDatabase {
	base.DisplayName = data.DisplayName
	base.Description = data.Description
	return base
}

// Validate implements concreteOperations.
// It checks for duplicate DisplayName among existing entities.
func (o *operationsMirroredDatabase) Validate(newEntity mirroreddatabase.MirroredDatabase, existing []mirroreddatabase.MirroredDatabase) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(
				http.StatusConflict,
				fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(),
				fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(),
			)
		}
	}
	return http.StatusCreated, nil
}

// Add this function below your operationsMirroredDatabase type and its methods.

func configureMirroredDatabase(server *fakeServer) mirroreddatabase.MirroredDatabase {
	// Define the entity operations interface for MirroredDatabase.
	type concreteEntityOperations interface {
		parentIDOperations[
			mirroreddatabase.MirroredDatabase,
			mirroreddatabase.ItemsClientGetMirroredDatabaseResponse,
			mirroreddatabase.ItemsClientUpdateMirroredDatabaseResponse,
			mirroreddatabase.ItemsClientCreateMirroredDatabaseResponse,
			mirroreddatabase.ItemsClientListMirroredDatabasesResponse,
			mirroreddatabase.CreateMirroredDatabaseRequest,
			mirroreddatabase.UpdateMirroredDatabaseRequest,
		]
	}

	// Define the definition operations interface for MirroredDatabase.
	type concreteDefinitionOperations interface {
		definitionOperations[
			mirroreddatabase.Definition,
			mirroreddatabase.CreateMirroredDatabaseRequest,
			mirroreddatabase.UpdateMirroredDatabaseDefinitionRequest,
			mirroreddatabase.ItemsClientGetMirroredDatabaseDefinitionResponse,
			mirroreddatabase.ItemsClientUpdateMirroredDatabaseDefinitionResponse,
		]
	}

	var entityOperations concreteEntityOperations = &operationsMirroredDatabase{}
	var definitionOperations concreteDefinitionOperations = &operationsMirroredDatabase{}
	var converter itemConverter[mirroreddatabase.MirroredDatabase] = &operationsMirroredDatabase{}

	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	// handleGetWithParentID(handler, entityOperations, &handler.ServerFactory.MirroredDatabase.ItemsServer.GetMirroredDatabase)
	// handleUpdateWithParentID(handler, entityOperations, entityOperations, &handler.ServerFactory.MirroredDatabase.ItemsServer.UpdateMirroredDatabase)
	// handleNonLROCreate(handler, entityOperations, entityOperations, entityOperations, &handler.ServerFactory.MirroredDatabase.ItemsServer.CreateMirroredDatabase)
	// handleListPagerWithParentID(handler, entityOperations, entityOperations, &handler.ServerFactory.MirroredDatabase.ItemsServer.NewListMirroredDatabasesPager)
	// handleDeleteWithParentID(handler, &handler.ServerFactory.MirroredDatabase.ItemsServer.DeleteMirroredDatabase)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.MirroredDatabase.ItemsServer.GetMirroredDatabase,
		&server.ServerFactory.MirroredDatabase.ItemsServer.UpdateMirroredDatabase,
		&server.ServerFactory.MirroredDatabase.ItemsServer.CreateMirroredDatabase,
		&server.ServerFactory.MirroredDatabase.ItemsServer.NewListMirroredDatabasesPager,
		&server.ServerFactory.MirroredDatabase.ItemsServer.DeleteMirroredDatabase,
	)

	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.MirroredDatabase.ItemsServer.CreateMirroredDatabase,
		&server.ServerFactory.MirroredDatabase.ItemsServer.GetMirroredDatabaseDefinition,
		&server.ServerFactory.MirroredDatabase.ItemsServer.UpdateMirroredDatabaseDefinition,
	)

	return mirroreddatabase.MirroredDatabase{}
}

// NewRandomMirroredDatabase returns a random Mirrored Database entity.
func NewRandomMirroredDatabase() mirroreddatabase.MirroredDatabase {
	return mirroreddatabase.MirroredDatabase{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(mirroreddatabase.ItemTypeMirroredDatabase),
		Properties: &mirroreddatabase.Properties{
			DefaultSchema:     to.Ptr(testhelp.RandomName()),
			OneLakeTablesPath: to.Ptr(testhelp.RandomURI()),
			SQLEndpointProperties: &mirroreddatabase.SQLEndpointProperties{
				ProvisioningStatus: to.Ptr(mirroreddatabase.SQLEndpointProvisioningStatusSuccess),
				ConnectionString:   to.Ptr("Server=test;Database=mydb;User Id=test;Password=secret;"),
				ID:                 to.Ptr(testhelp.RandomUUID()),
			},
		},
	}
}

// NewRandomMirroredDatabaseWithWorkspace returns a random Mirrored Database entity with a given workspaceID.
func NewRandomMirroredDatabaseWithWorkspace(workspaceID string) mirroreddatabase.MirroredDatabase {
	db := NewRandomMirroredDatabase()
	db.WorkspaceID = to.Ptr(workspaceID)
	return db
}
