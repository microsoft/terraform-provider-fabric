// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabkqldatabase "github.com/microsoft/fabric-sdk-go/fabric/kqldatabase"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsKQLDatabase struct{}

// ConvertItemToEntity implements itemConverter.
func (o *operationsKQLDatabase) ConvertItemToEntity(item fabcore.Item) fabkqldatabase.KQLDatabase {
	return fabkqldatabase.KQLDatabase{
		ID:          item.ID,
		DisplayName: item.DisplayName,
		Description: item.Description,
		WorkspaceID: item.WorkspaceID,
		Type:        to.Ptr(fabkqldatabase.ItemTypeKQLDatabase),
		Properties:  NewRandomKQLDatabase().Properties,
	}
}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsKQLDatabase) CreateDefinition(data fabkqldatabase.CreateKQLDatabaseRequest) *fabkqldatabase.Definition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsKQLDatabase) TransformDefinition(entity *fabkqldatabase.Definition) fabkqldatabase.ItemsClientGetKQLDatabaseDefinitionResponse {
	return fabkqldatabase.ItemsClientGetKQLDatabaseDefinitionResponse{
		DefinitionResponse: fabkqldatabase.DefinitionResponse{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsKQLDatabase) UpdateDefinition(_ *fabkqldatabase.Definition, data fabkqldatabase.UpdateKQLDatabaseDefinitionRequest) *fabkqldatabase.Definition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsKQLDatabase) CreateWithParentID(parentID string, data fabkqldatabase.CreateKQLDatabaseRequest) fabkqldatabase.KQLDatabase {
	entity := NewRandomKQLDatabaseWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description

	return entity
}

// Filter implements concreteOperations.
func (o *operationsKQLDatabase) Filter(entities []fabkqldatabase.KQLDatabase, parentID string) []fabkqldatabase.KQLDatabase {
	ret := make([]fabkqldatabase.KQLDatabase, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsKQLDatabase) GetID(entity fabkqldatabase.KQLDatabase) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsKQLDatabase) TransformCreate(entity fabkqldatabase.KQLDatabase) fabkqldatabase.ItemsClientCreateKQLDatabaseResponse {
	return fabkqldatabase.ItemsClientCreateKQLDatabaseResponse{
		KQLDatabase: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsKQLDatabase) TransformGet(entity fabkqldatabase.KQLDatabase) fabkqldatabase.ItemsClientGetKQLDatabaseResponse {
	return fabkqldatabase.ItemsClientGetKQLDatabaseResponse{
		KQLDatabase: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsKQLDatabase) TransformList(entities []fabkqldatabase.KQLDatabase) fabkqldatabase.ItemsClientListKQLDatabasesResponse {
	return fabkqldatabase.ItemsClientListKQLDatabasesResponse{
		KQLDatabases: fabkqldatabase.KQLDatabases{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsKQLDatabase) TransformUpdate(entity fabkqldatabase.KQLDatabase) fabkqldatabase.ItemsClientUpdateKQLDatabaseResponse {
	return fabkqldatabase.ItemsClientUpdateKQLDatabaseResponse{
		KQLDatabase: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsKQLDatabase) Update(base fabkqldatabase.KQLDatabase, data fabkqldatabase.UpdateKQLDatabaseRequest) fabkqldatabase.KQLDatabase {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsKQLDatabase) Validate(newEntity fabkqldatabase.KQLDatabase, existing []fabkqldatabase.KQLDatabase) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusOK, nil
}

func configureKQLDatabase(server *fakeServer) fabkqldatabase.KQLDatabase {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabkqldatabase.KQLDatabase,
			fabkqldatabase.ItemsClientGetKQLDatabaseResponse,
			fabkqldatabase.ItemsClientUpdateKQLDatabaseResponse,
			fabkqldatabase.ItemsClientCreateKQLDatabaseResponse,
			fabkqldatabase.ItemsClientListKQLDatabasesResponse,
			fabkqldatabase.CreateKQLDatabaseRequest,
			fabkqldatabase.UpdateKQLDatabaseRequest]
	}

	type concreteDefinitionOperations interface {
		definitionOperations[
			fabkqldatabase.Definition,
			fabkqldatabase.CreateKQLDatabaseRequest,
			fabkqldatabase.UpdateKQLDatabaseDefinitionRequest,
			fabkqldatabase.ItemsClientGetKQLDatabaseDefinitionResponse,
			fabkqldatabase.ItemsClientUpdateKQLDatabaseDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsKQLDatabase{}

	var definitionOperations concreteDefinitionOperations = &operationsKQLDatabase{}

	var converter itemConverter[fabkqldatabase.KQLDatabase] = &operationsKQLDatabase{}

	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.KQLDatabase.ItemsServer.GetKQLDatabase,
		&server.ServerFactory.KQLDatabase.ItemsServer.UpdateKQLDatabase,
		&server.ServerFactory.KQLDatabase.ItemsServer.BeginCreateKQLDatabase,
		&server.ServerFactory.KQLDatabase.ItemsServer.NewListKQLDatabasesPager,
		&server.ServerFactory.KQLDatabase.ItemsServer.DeleteKQLDatabase)

	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.KQLDatabase.ItemsServer.BeginCreateKQLDatabase,
		&server.ServerFactory.KQLDatabase.ItemsServer.BeginGetKQLDatabaseDefinition,
		&server.ServerFactory.KQLDatabase.ItemsServer.BeginUpdateKQLDatabaseDefinition)

	return fabkqldatabase.KQLDatabase{}
}

func NewRandomKQLDatabase() fabkqldatabase.KQLDatabase {
	return fabkqldatabase.KQLDatabase{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fabkqldatabase.ItemTypeKQLDatabase),
		Properties: &fabkqldatabase.Properties{
			DatabaseType:           to.Ptr(fabkqldatabase.TypeReadWrite),
			ParentEventhouseItemID: to.Ptr(testhelp.RandomUUID()),
			IngestionServiceURI:    to.Ptr(testhelp.RandomURI()),
			QueryServiceURI:        to.Ptr(testhelp.RandomURI()),
		},
	}
}

func NewRandomKQLDatabaseWithWorkspace(workspaceID string) fabkqldatabase.KQLDatabase {
	result := NewRandomKQLDatabase()
	result.WorkspaceID = &workspaceID

	return result
}

func NewRandomKQLDatabaseDefinition() fabkqldatabase.Definition {
	defPart1 := fabkqldatabase.DefinitionPart{
		PayloadType: to.Ptr(fabkqldatabase.PayloadTypeInlineBase64),
		Path:        to.Ptr("DatabaseProperties.json"),
		Payload: to.Ptr(
			"ew0KICAiZGF0YWJhc2VUeXBlIjogIlJlYWRXcml0ZSIsDQogICJwYXJlbnRFdmVudGhvdXNlSXRlbUlkIjogIjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsIA0KICAib25lTGFrZUNhY2hpbmdQZXJpb2QiOiAiUDM2NTAwRCIsIA0KICAib25lTGFrZVN0YW5kYXJkU3RvcmFnZVBlcmlvZCI6ICJQMzY1MDAwRCIgDQp9",
		),
	}

	defPart2 := fabkqldatabase.DefinitionPart{
		PayloadType: to.Ptr(fabkqldatabase.PayloadTypeInlineBase64),
		Path:        to.Ptr("DatabaseSchema.kql"),
		Payload: to.Ptr(
			"LmNyZWF0ZS1tZXJnZSB0YWJsZSBNeUxvZ3MyIChMZXZlbDpzdHJpbmcsIFRpbWVzdGFtcDpkYXRldGltZSwgVXNlcklkOnN0cmluZywgVHJhY2VJZDpzdHJpbmcsIE1lc3NhZ2U6c3RyaW5nLCBQcm9jZXNzSWQ6aW50KSANCi5jcmVhdGUtbWVyZ2UgdGFibGUgTXlMb2dzMyAoTGV2ZWw6c3RyaW5nLCBUaW1lc3RhbXA6ZGF0ZXRpbWUsIFVzZXJJZDpzdHJpbmcsIFRyYWNlSWQ6c3RyaW5nLCBNZXNzYWdlOnN0cmluZywgUHJvY2Vzc0lkOmludCkgDQouY3JlYXRlLW1lcmdlIHRhYmxlIE15TG9nczcgKExldmVsOnN0cmluZywgVGltZXN0YW1wOmRhdGV0aW1lLCBVc2VySWQ6c3RyaW5nLCBUcmFjZUlkOnN0cmluZywgTWVzc2FnZTpzdHJpbmcsIFByb2Nlc3NJZDppbnQp",
		),
	}

	var defParts []fabkqldatabase.DefinitionPart

	defParts = append(defParts, defPart1)
	defParts = append(defParts, defPart2)

	return fabkqldatabase.Definition{
		Parts: defParts,
	}
}
