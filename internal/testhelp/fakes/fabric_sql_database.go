// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabsqldatabase "github.com/microsoft/fabric-sdk-go/fabric/sqldatabase"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsSQLDatabase struct{}

// ConvertItemToEntity implements itemConverter.
func (o *operationsSQLDatabase) ConvertItemToEntity(item fabcore.Item) fabsqldatabase.SQLDatabase {
	return fabsqldatabase.SQLDatabase{
		ID:          item.ID,
		DisplayName: item.DisplayName,
		Description: item.Description,
		WorkspaceID: item.WorkspaceID,
		FolderID:    item.FolderID,
		Type:        to.Ptr(fabsqldatabase.ItemTypeSQLDatabase),
		Properties:  NewRandomSQLDatabase().Properties,
	}
}

// CreateWithParentID implements concreteOperations.
func (o *operationsSQLDatabase) CreateWithParentID(parentID string, data fabsqldatabase.CreateSQLDatabaseRequest) fabsqldatabase.SQLDatabase {
	entity := NewRandomSQLDatabaseWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description
	entity.FolderID = data.FolderID

	return entity
}

// Filter implements concreteOperations.
func (o *operationsSQLDatabase) Filter(entities []fabsqldatabase.SQLDatabase, parentID string) []fabsqldatabase.SQLDatabase {
	ret := make([]fabsqldatabase.SQLDatabase, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsSQLDatabase) GetID(entity fabsqldatabase.SQLDatabase) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsSQLDatabase) TransformCreate(entity fabsqldatabase.SQLDatabase) fabsqldatabase.ItemsClientCreateSQLDatabaseResponse {
	return fabsqldatabase.ItemsClientCreateSQLDatabaseResponse{
		SQLDatabase: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsSQLDatabase) TransformGet(entity fabsqldatabase.SQLDatabase) fabsqldatabase.ItemsClientGetSQLDatabaseResponse {
	return fabsqldatabase.ItemsClientGetSQLDatabaseResponse{
		SQLDatabase: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsSQLDatabase) TransformList(entities []fabsqldatabase.SQLDatabase) fabsqldatabase.ItemsClientListSQLDatabasesResponse {
	return fabsqldatabase.ItemsClientListSQLDatabasesResponse{
		SQLDatabases: fabsqldatabase.SQLDatabases{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsSQLDatabase) TransformUpdate(entity fabsqldatabase.SQLDatabase) fabsqldatabase.ItemsClientUpdateSQLDatabaseResponse {
	return fabsqldatabase.ItemsClientUpdateSQLDatabaseResponse{
		SQLDatabase: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsSQLDatabase) Update(base fabsqldatabase.SQLDatabase, data fabsqldatabase.UpdateSQLDatabaseRequest) fabsqldatabase.SQLDatabase {
	base.Description = data.Description

	return base
}

// Validate implements concreteOperations.
func (o *operationsSQLDatabase) Validate(newEntity fabsqldatabase.SQLDatabase, existing []fabsqldatabase.SQLDatabase) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureSQLDatabase(server *fakeServer) fabsqldatabase.SQLDatabase {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabsqldatabase.SQLDatabase,
			fabsqldatabase.ItemsClientGetSQLDatabaseResponse,
			fabsqldatabase.ItemsClientUpdateSQLDatabaseResponse,
			fabsqldatabase.ItemsClientCreateSQLDatabaseResponse,
			fabsqldatabase.ItemsClientListSQLDatabasesResponse,
			fabsqldatabase.CreateSQLDatabaseRequest,
			fabsqldatabase.UpdateSQLDatabaseRequest]
	}

	var entityOperations concreteEntityOperations = &operationsSQLDatabase{}
	var converter itemConverter[fabsqldatabase.SQLDatabase] = &operationsSQLDatabase{}
	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.SQLDatabase.ItemsServer.GetSQLDatabase,
		&server.ServerFactory.SQLDatabase.ItemsServer.UpdateSQLDatabase,
		&server.ServerFactory.SQLDatabase.ItemsServer.BeginCreateSQLDatabase,
		&server.ServerFactory.SQLDatabase.ItemsServer.NewListSQLDatabasesPager,
		&server.ServerFactory.SQLDatabase.ItemsServer.DeleteSQLDatabase)

	return fabsqldatabase.SQLDatabase{}
}

func NewRandomSQLDatabase() fabsqldatabase.SQLDatabase {
	return fabsqldatabase.SQLDatabase{
		ID:          new(testhelp.RandomUUID()),
		DisplayName: new(testhelp.RandomName()),
		Description: new(testhelp.RandomName()),
		WorkspaceID: new(testhelp.RandomUUID()),
		FolderID:    new(testhelp.RandomUUID()),
		Type:        to.Ptr(fabsqldatabase.ItemTypeSQLDatabase),
		Properties: &fabsqldatabase.Properties{
			ConnectionString:     new(testhelp.RandomName()),
			DatabaseName:         new(testhelp.RandomName()),
			ServerFqdn:           new(testhelp.RandomName()),
			BackupRetentionDays:  to.Ptr(int32(testhelp.RandomIntRange(1, 35))),
			Collation:            to.Ptr("SQL_Latin1_General_CP1_CI_AS"),
			EarliestRestorePoint: to.Ptr(time.Now().Add(-24 * time.Hour).UTC().Truncate(time.Second)),
			LatestRestorePoint:   to.Ptr(time.Now().UTC().Truncate(time.Second)),
		},
	}
}

func NewRandomSQLDatabaseWithWorkspace(workspaceID string) fabsqldatabase.SQLDatabase {
	result := NewRandomSQLDatabase()
	result.WorkspaceID = &workspaceID

	return result
}
