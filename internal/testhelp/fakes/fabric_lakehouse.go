// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fablakehouse "github.com/microsoft/fabric-sdk-go/fabric/lakehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsLakehouse struct{}

// ConvertItemToEntity implements itemConverter.
func (o *operationsLakehouse) ConvertItemToEntity(item fabcore.Item) fablakehouse.Lakehouse {
	return fablakehouse.Lakehouse{
		ID:          item.ID,
		DisplayName: item.DisplayName,
		Description: item.Description,
		WorkspaceID: item.WorkspaceID,
		FolderID:    item.FolderID,
		Type:        to.Ptr(fablakehouse.ItemTypeLakehouse),
		Properties:  NewRandomLakehouse().Properties,
	}
}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsLakehouse) CreateDefinition(data fablakehouse.CreateLakehouseRequest) *fablakehouse.Definition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsLakehouse) TransformDefinition(entity *fablakehouse.Definition) fablakehouse.ItemsClientGetLakehouseDefinitionResponse {
	return fablakehouse.ItemsClientGetLakehouseDefinitionResponse{
		DefinitionResponse: fablakehouse.DefinitionResponse{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsLakehouse) UpdateDefinition(_ *fablakehouse.Definition, data fablakehouse.UpdateLakehouseDefinitionRequest) *fablakehouse.Definition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsLakehouse) CreateWithParentID(parentID string, data fablakehouse.CreateLakehouseRequest) fablakehouse.Lakehouse {
	entity := NewRandomLakehouseWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description
	entity.FolderID = data.FolderID

	return entity
}

// Filter implements concreteOperations.
func (o *operationsLakehouse) Filter(entities []fablakehouse.Lakehouse, parentID string) []fablakehouse.Lakehouse {
	ret := make([]fablakehouse.Lakehouse, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsLakehouse) GetID(entity fablakehouse.Lakehouse) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsLakehouse) TransformCreate(entity fablakehouse.Lakehouse) fablakehouse.ItemsClientCreateLakehouseResponse {
	return fablakehouse.ItemsClientCreateLakehouseResponse{
		Lakehouse: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsLakehouse) TransformGet(entity fablakehouse.Lakehouse) fablakehouse.ItemsClientGetLakehouseResponse {
	return fablakehouse.ItemsClientGetLakehouseResponse{
		Lakehouse: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsLakehouse) TransformList(entities []fablakehouse.Lakehouse) fablakehouse.ItemsClientListLakehousesResponse {
	return fablakehouse.ItemsClientListLakehousesResponse{
		Lakehouses: fablakehouse.Lakehouses{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsLakehouse) TransformUpdate(entity fablakehouse.Lakehouse) fablakehouse.ItemsClientUpdateLakehouseResponse {
	return fablakehouse.ItemsClientUpdateLakehouseResponse{
		Lakehouse: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsLakehouse) Update(base fablakehouse.Lakehouse, data fablakehouse.UpdateLakehouseRequest) fablakehouse.Lakehouse {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsLakehouse) Validate(newEntity fablakehouse.Lakehouse, existing []fablakehouse.Lakehouse) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureLakehouse(server *fakeServer) fablakehouse.Lakehouse {
	type concreteEntityOperations interface {
		parentIDOperations[
			fablakehouse.Lakehouse,
			fablakehouse.ItemsClientGetLakehouseResponse,
			fablakehouse.ItemsClientUpdateLakehouseResponse,
			fablakehouse.ItemsClientCreateLakehouseResponse,
			fablakehouse.ItemsClientListLakehousesResponse,
			fablakehouse.CreateLakehouseRequest,
			fablakehouse.UpdateLakehouseRequest]
	}
	type concreteDefinitionOperations interface {
		definitionOperations[
			fablakehouse.Definition,
			fablakehouse.CreateLakehouseRequest,
			fablakehouse.UpdateLakehouseDefinitionRequest,
			fablakehouse.ItemsClientGetLakehouseDefinitionResponse,
			fablakehouse.ItemsClientUpdateLakehouseDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsLakehouse{}
	var definitionOperations concreteDefinitionOperations = &operationsLakehouse{}

	var converter itemConverter[fablakehouse.Lakehouse] = &operationsLakehouse{}

	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.Lakehouse.ItemsServer.GetLakehouse,
		&server.ServerFactory.Lakehouse.ItemsServer.UpdateLakehouse,
		&server.ServerFactory.Lakehouse.ItemsServer.BeginCreateLakehouse,
		&server.ServerFactory.Lakehouse.ItemsServer.NewListLakehousesPager,
		&server.ServerFactory.Lakehouse.ItemsServer.DeleteLakehouse)
	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.Lakehouse.ItemsServer.BeginCreateLakehouse,
		&server.ServerFactory.Lakehouse.ItemsServer.BeginGetLakehouseDefinition,
		&server.ServerFactory.Lakehouse.ItemsServer.BeginUpdateLakehouseDefinition)

	return fablakehouse.Lakehouse{}
}

func NewRandomLakehouse() fablakehouse.Lakehouse {
	return fablakehouse.Lakehouse{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fablakehouse.ItemTypeLakehouse),
		FolderID:    to.Ptr(testhelp.RandomUUID()),
		Properties: &fablakehouse.Properties{
			OneLakeFilesPath:  to.Ptr(testhelp.RandomName()),
			OneLakeTablesPath: to.Ptr(testhelp.RandomName()),
			DefaultSchema:     to.Ptr("dbo"),
			SQLEndpointProperties: &fablakehouse.SQLEndpointProperties{
				ID:                 to.Ptr(testhelp.RandomUUID()),
				ProvisioningStatus: to.Ptr(fablakehouse.SQLEndpointProvisioningStatusSuccess),
				ConnectionString:   to.Ptr(testhelp.RandomURI()),
			},
		},
	}
}

func NewRandomLakehouseWithWorkspace(workspaceID string) fablakehouse.Lakehouse {
	result := NewRandomLakehouse()
	result.WorkspaceID = &workspaceID

	return result
}

func NewRandomLakehouseDefinition() fablakehouse.Definition {
	defPart1 := fablakehouse.DefinitionPart{
		PayloadType: to.Ptr(fablakehouse.PayloadTypeInlineBase64),
		Path:        to.Ptr("lakehouse.metadata.json"),
		Payload: to.Ptr(
			"eyJkZWZhdWx0U2NoZW1hIjoiZGJvIn0==",
		),
	}

	defPart2 := fablakehouse.DefinitionPart{
		PayloadType: to.Ptr(fablakehouse.PayloadTypeInlineBase64),
		Path:        to.Ptr("shortcuts.metadata.json"),
		Payload: to.Ptr(
			"WwogIHsKICAgICJuYW1lIjogIk55Y1RheGkiLAogICAgInBhdGgiOiAiL1RhYmxlcyIsCiAgICAidGFyZ2V0IjogewogICAgICAi",
		),
	}

	defPart3 := fablakehouse.DefinitionPart{
		PayloadType: to.Ptr(fablakehouse.PayloadTypeInlineBase64),
		Path:        to.Ptr("data-access-roles.json"),
		Payload: to.Ptr(
			"Ww0KICAgICAgICAgICAgICB7DQogICAgICAgICAgICAgICAgIm5hbWUiOiAiZGltZW5zaW9ucnVsZXJlbmFtZSIsDQogICAgICA==",
		),
	}

	defPart4 := fablakehouse.DefinitionPart{
		PayloadType: to.Ptr(fablakehouse.PayloadTypeInlineBase64),
		Path:        to.Ptr("alm.settings.json"),
		Payload: to.Ptr(
			"ew0KICAgICAgICAgICJ2ZXJzaW9uIjogIjEuMC4xIiwNCiAgICAgICAgICAib2JqZWN0VHlwZXMiOiBbDQogICAgICAgICAgICB7",
		),
	}

	defParts := make([]fablakehouse.DefinitionPart, 0, 4)

	defParts = append(defParts, defPart1)
	defParts = append(defParts, defPart2)
	defParts = append(defParts, defPart3)
	defParts = append(defParts, defPart4)

	return fablakehouse.Definition{
		Format: to.Ptr("LakehouseDefinitionV1"),
		Parts:  defParts,
	}
}
