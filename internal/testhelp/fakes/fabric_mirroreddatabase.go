// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabmirroreddatabase "github.com/microsoft/fabric-sdk-go/fabric/mirroreddatabase"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// operationsMirroredDatabase implements the operations for MirroredDatabase
// and satisfies the concreteEntityOperations and concreteDefinitionOperations interfaces.
type operationsMirroredDatabase struct{}

// ConvertItemToEntity implements itemConverter. It converts a generic fabcore.Item into a
// fabmirroreddatabase.MirroredDatabase.
func (o *operationsMirroredDatabase) ConvertItemToEntity(item fabcore.Item) fabmirroreddatabase.MirroredDatabase {
	return fabmirroreddatabase.MirroredDatabase{
		ID:          item.ID,
		DisplayName: item.DisplayName,
		Description: item.Description,
		WorkspaceID: item.WorkspaceID,
		// Set the type to mirrored database.
		Type:       to.Ptr(fabmirroreddatabase.ItemTypeMirroredDatabase),
		Properties: NewRandomMirroredDatabase().Properties,
	}
}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsMirroredDatabase) CreateDefinition(data fabcore.CreateItemRequest) *fabcore.ItemDefinition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsMirroredDatabase) TransformDefinition(entity *fabcore.ItemDefinition) fabcore.ItemsClientGetItemDefinitionResponse {
	return fabcore.ItemsClientGetItemDefinitionResponse{
		ItemDefinitionResponse: fabcore.ItemDefinitionResponse{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsMirroredDatabase) UpdateDefinition(_ *fabcore.ItemDefinition, data fabcore.UpdateItemDefinitionRequest) *fabcore.ItemDefinition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
// It creates a mirrored database entity with the given parent workspace ID.
func (o *operationsMirroredDatabase) CreateWithParentID(parentID string, data fabmirroreddatabase.CreateMirroredDatabaseRequest) fabmirroreddatabase.MirroredDatabase {
	entity := NewRandomMirroredDatabaseWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description

	return entity
}

// Filter implements concreteOperations, returning only entities whose WorkspaceID matches parentID.
func (o *operationsMirroredDatabase) Filter(entities []fabmirroreddatabase.MirroredDatabase, parentID string) []fabmirroreddatabase.MirroredDatabase {
	var ret []fabmirroreddatabase.MirroredDatabase

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
// It returns a generated unique ID composed of the workspace ID and the entity ID.
func (o *operationsMirroredDatabase) GetID(entity fabmirroreddatabase.MirroredDatabase) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations. It wraps the created entity in the Create response.
func (o *operationsMirroredDatabase) TransformCreate(entity fabmirroreddatabase.MirroredDatabase) fabmirroreddatabase.ItemsClientCreateMirroredDatabaseResponse {
	return fabmirroreddatabase.ItemsClientCreateMirroredDatabaseResponse{
		MirroredDatabase: entity,
	}
}

// TransformGet implements concreteOperations.
// It returns the get response for the given mirrored database entity.
func (o *operationsMirroredDatabase) TransformGet(entity fabmirroreddatabase.MirroredDatabase) fabmirroreddatabase.ItemsClientGetMirroredDatabaseResponse {
	return fabmirroreddatabase.ItemsClientGetMirroredDatabaseResponse{
		MirroredDatabase: entity,
	}
}

// TransformList implements concreteOperations.
// It wraps a list of entities into the List response.
func (o *operationsMirroredDatabase) TransformList(entities []fabmirroreddatabase.MirroredDatabase) fabmirroreddatabase.ItemsClientListMirroredDatabasesResponse {
	return fabmirroreddatabase.ItemsClientListMirroredDatabasesResponse{
		MirroredDatabases: fabmirroreddatabase.MirroredDatabases{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
// It wraps the updated entity into the Update response.
func (o *operationsMirroredDatabase) TransformUpdate(entity fabmirroreddatabase.MirroredDatabase) fabmirroreddatabase.ItemsClientUpdateMirroredDatabaseResponse {
	return fabmirroreddatabase.ItemsClientUpdateMirroredDatabaseResponse{
		MirroredDatabase: entity,
	}
}

// Update implements concreteOperations.
// It updates the base entity with the new data.
func (o *operationsMirroredDatabase) Update(base fabmirroreddatabase.MirroredDatabase, data fabmirroreddatabase.UpdateMirroredDatabaseRequest) fabmirroreddatabase.MirroredDatabase {
	base.DisplayName = data.DisplayName
	base.Description = data.Description

	return base
}

// Validate implements concreteOperations.
// It checks for duplicate DisplayName among existing entities.
func (o *operationsMirroredDatabase) Validate(newEntity fabmirroreddatabase.MirroredDatabase, existing []fabmirroreddatabase.MirroredDatabase) (int, error) {
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

func configureMirroredDatabase(server *fakeServer) fabmirroreddatabase.MirroredDatabase {
	// Define the entity operations interface for MirroredDatabase.
	type concreteEntityOperations interface {
		parentIDOperations[
			fabmirroreddatabase.MirroredDatabase,
			fabmirroreddatabase.ItemsClientGetMirroredDatabaseResponse,
			fabmirroreddatabase.ItemsClientUpdateMirroredDatabaseResponse,
			fabmirroreddatabase.ItemsClientCreateMirroredDatabaseResponse,
			fabmirroreddatabase.ItemsClientListMirroredDatabasesResponse,
			fabmirroreddatabase.CreateMirroredDatabaseRequest,
			fabmirroreddatabase.UpdateMirroredDatabaseRequest,
		]
	}

	// Define the definition operations interface for MirroredDatabase.
	type concreteDefinitionOperations interface {
		definitionOperationsNonLROCreation[
			fabcore.ItemDefinition,
			fabcore.UpdateItemDefinitionRequest,
			fabcore.ItemsClientGetItemDefinitionResponse,
			fabcore.ItemsClientUpdateItemDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsMirroredDatabase{}
	var definitionOperations concreteDefinitionOperations = &operationsMirroredDatabase{}
	var converter itemConverter[fabmirroreddatabase.MirroredDatabase] = &operationsMirroredDatabase{}

	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	configureNonLROEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.MirroredDatabase.ItemsServer.GetMirroredDatabase,
		&server.ServerFactory.MirroredDatabase.ItemsServer.UpdateMirroredDatabase,
		&server.ServerFactory.MirroredDatabase.ItemsServer.CreateMirroredDatabase,
		&server.ServerFactory.MirroredDatabase.ItemsServer.NewListMirroredDatabasesPager,
		&server.ServerFactory.MirroredDatabase.ItemsServer.DeleteMirroredDatabase,
	)

	ConfigureDefinitionsNonLROCreation(
		handler,
		definitionOperations,
		&server.ServerFactory.Core.ItemsServer.BeginGetItemDefinition,
		&server.ServerFactory.Core.ItemsServer.BeginUpdateItemDefinition,
	)

	return fabmirroreddatabase.MirroredDatabase{}
}

// NewRandomMirroredDatabase returns a random Mirrored Database entity.
func NewRandomMirroredDatabase() fabmirroreddatabase.MirroredDatabase {
	return fabmirroreddatabase.MirroredDatabase{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fabmirroreddatabase.ItemTypeMirroredDatabase),
		Properties: &fabmirroreddatabase.Properties{
			DefaultSchema:     to.Ptr(testhelp.RandomName()),
			OneLakeTablesPath: to.Ptr(testhelp.RandomURI()),
			SQLEndpointProperties: &fabmirroreddatabase.SQLEndpointProperties{
				ProvisioningStatus: to.Ptr(fabmirroreddatabase.SQLEndpointProvisioningStatusSuccess),
				ConnectionString:   to.Ptr(testhelp.RandomURI()),
				ID:                 to.Ptr(testhelp.RandomUUID()),
			},
		},
	}
}

// NewRandomMirroredDatabaseWithWorkspace returns a random Mirrored Database entity with a given workspaceID.
func NewRandomMirroredDatabaseWithWorkspace(workspaceID string) fabmirroreddatabase.MirroredDatabase {
	db := NewRandomMirroredDatabase()
	db.WorkspaceID = to.Ptr(workspaceID)

	return db
}
