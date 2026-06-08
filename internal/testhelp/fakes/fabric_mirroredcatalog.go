// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabmirroredcatalog "github.com/microsoft/fabric-sdk-go/fabric/mirroredcatalog"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// operationsMirroredCatalog implements the operations for MirroredCatalog
// and satisfies the concreteEntityOperations and concreteDefinitionOperations interfaces.
type operationsMirroredCatalog struct{}

// ConvertItemToEntity implements itemConverter.
func (o *operationsMirroredCatalog) ConvertItemToEntity(item fabcore.Item) fabmirroredcatalog.MirroredCatalog {
	return fabmirroredcatalog.MirroredCatalog{
		ID:          item.ID,
		DisplayName: item.DisplayName,
		Description: item.Description,
		WorkspaceID: item.WorkspaceID,
		FolderID:    item.FolderID,
		Type:        to.Ptr(fabmirroredcatalog.ItemTypeMirroredCatalog),
		Properties:  NewRandomMirroredCatalog().Properties,
	}
}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsMirroredCatalog) CreateDefinition(data fabcore.CreateItemRequest) *fabcore.ItemDefinition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsMirroredCatalog) TransformDefinition(entity *fabcore.ItemDefinition) fabcore.ItemsClientGetItemDefinitionResponse {
	return fabcore.ItemsClientGetItemDefinitionResponse{
		ItemDefinitionResponse: fabcore.ItemDefinitionResponse{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsMirroredCatalog) UpdateDefinition(_ *fabcore.ItemDefinition, data fabcore.UpdateItemDefinitionRequest) *fabcore.ItemDefinition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsMirroredCatalog) CreateWithParentID(parentID string, data fabmirroredcatalog.CreateMirroredCatalogRequest) fabmirroredcatalog.MirroredCatalog {
	entity := NewRandomMirroredCatalogWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description
	entity.FolderID = data.FolderID

	return entity
}

// Filter implements concreteOperations.
func (o *operationsMirroredCatalog) Filter(entities []fabmirroredcatalog.MirroredCatalog, parentID string) []fabmirroredcatalog.MirroredCatalog {
	var ret []fabmirroredcatalog.MirroredCatalog

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsMirroredCatalog) GetID(entity fabmirroredcatalog.MirroredCatalog) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsMirroredCatalog) TransformCreate(entity fabmirroredcatalog.MirroredCatalog) fabmirroredcatalog.ItemsClientCreateMirroredCatalogResponse {
	return fabmirroredcatalog.ItemsClientCreateMirroredCatalogResponse{
		MirroredCatalog: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsMirroredCatalog) TransformGet(entity fabmirroredcatalog.MirroredCatalog) fabmirroredcatalog.ItemsClientGetMirroredCatalogResponse {
	return fabmirroredcatalog.ItemsClientGetMirroredCatalogResponse{
		MirroredCatalog: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsMirroredCatalog) TransformList(entities []fabmirroredcatalog.MirroredCatalog) fabmirroredcatalog.ItemsClientListMirroredCatalogsResponse {
	return fabmirroredcatalog.ItemsClientListMirroredCatalogsResponse{
		MirroredCatalogs: fabmirroredcatalog.MirroredCatalogs{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsMirroredCatalog) TransformUpdate(entity fabmirroredcatalog.MirroredCatalog) fabmirroredcatalog.ItemsClientUpdateMirroredCatalogResponse {
	return fabmirroredcatalog.ItemsClientUpdateMirroredCatalogResponse{
		MirroredCatalog: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsMirroredCatalog) Update(base fabmirroredcatalog.MirroredCatalog, data fabmirroredcatalog.UpdateMirroredCatalogRequest) fabmirroredcatalog.MirroredCatalog {
	base.DisplayName = data.DisplayName
	base.Description = data.Description

	return base
}

// Validate implements concreteOperations.
func (o *operationsMirroredCatalog) Validate(newEntity fabmirroredcatalog.MirroredCatalog, existing []fabmirroredcatalog.MirroredCatalog) (int, error) {
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

func configureMirroredCatalog(server *fakeServer) fabmirroredcatalog.MirroredCatalog {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabmirroredcatalog.MirroredCatalog,
			fabmirroredcatalog.ItemsClientGetMirroredCatalogResponse,
			fabmirroredcatalog.ItemsClientUpdateMirroredCatalogResponse,
			fabmirroredcatalog.ItemsClientCreateMirroredCatalogResponse,
			fabmirroredcatalog.ItemsClientListMirroredCatalogsResponse,
			fabmirroredcatalog.CreateMirroredCatalogRequest,
			fabmirroredcatalog.UpdateMirroredCatalogRequest,
		]
	}

	type concreteDefinitionOperations interface {
		definitionOperationsNonLROCreation[
			fabcore.ItemDefinition,
			fabcore.UpdateItemDefinitionRequest,
			fabcore.ItemsClientGetItemDefinitionResponse,
			fabcore.ItemsClientUpdateItemDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsMirroredCatalog{}
	var definitionOperations concreteDefinitionOperations = &operationsMirroredCatalog{}
	var converter itemConverter[fabmirroredcatalog.MirroredCatalog] = &operationsMirroredCatalog{}

	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.MirroredCatalog.ItemsServer.GetMirroredCatalog,
		&server.ServerFactory.MirroredCatalog.ItemsServer.UpdateMirroredCatalog,
		&server.ServerFactory.MirroredCatalog.ItemsServer.BeginCreateMirroredCatalog,
		&server.ServerFactory.MirroredCatalog.ItemsServer.NewListMirroredCatalogsPager,
		&server.ServerFactory.MirroredCatalog.ItemsServer.DeleteMirroredCatalog,
	)

	configureDefinitionsNonLROCreation(
		handler,
		definitionOperations,
		&server.ServerFactory.Core.ItemsServer.BeginGetItemDefinition,
		&server.ServerFactory.Core.ItemsServer.BeginUpdateItemDefinition,
	)

	return fabmirroredcatalog.MirroredCatalog{}
}

// NewRandomMirroredCatalog returns a random Mirrored Catalog entity.
func NewRandomMirroredCatalog() fabmirroredcatalog.MirroredCatalog {
	return fabmirroredcatalog.MirroredCatalog{
		ID:          new(testhelp.RandomUUID()),
		DisplayName: new(testhelp.RandomName()),
		Description: new(testhelp.RandomName()),
		WorkspaceID: new(testhelp.RandomUUID()),
		FolderID:    new(testhelp.RandomUUID()),
		Type:        to.Ptr(fabmirroredcatalog.ItemTypeMirroredCatalog),
		Properties: &fabmirroredcatalog.Properties{
			ConnectionID:      new(testhelp.RandomUUID()),
			OneLakeTablesPath: new(testhelp.RandomURI()),
			Scope:             []string{testhelp.RandomName(), testhelp.RandomName()},
			SourceType:        new("DremioIcebergCatalog"),
			SQLEndpointProperties: &fabmirroredcatalog.SQLEndpointProperties{
				ProvisioningStatus: to.Ptr(fabmirroredcatalog.ProvisioningStatusSuccess),
				ConnectionString:   new(testhelp.RandomURI()),
				ID:                 new(testhelp.RandomUUID()),
			},
		},
	}
}

// NewRandomMirroredCatalogWithWorkspace returns a random Mirrored Catalog with a given workspaceID.
func NewRandomMirroredCatalogWithWorkspace(workspaceID string) fabmirroredcatalog.MirroredCatalog {
	mc := NewRandomMirroredCatalog()
	mc.WorkspaceID = &workspaceID

	return mc
}
