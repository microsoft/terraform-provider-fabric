// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabmirroredazuredatabrickscatalog "github.com/microsoft/fabric-sdk-go/fabric/mirroredazuredatabrickscatalog"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsMirroredAzureDatabricksCatalog struct{}

// ConvertItemToEntity implements itemConverter.
func (o *operationsMirroredAzureDatabricksCatalog) ConvertItemToEntity(item fabcore.Item) fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog {
	return fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog{
		ID:          item.ID,
		DisplayName: item.DisplayName,
		Description: item.Description,
		WorkspaceID: item.WorkspaceID,
		FolderID:    item.FolderID,
		Type:        to.Ptr(fabmirroredazuredatabrickscatalog.ItemTypeMirroredAzureDatabricksCatalog),
		Properties:  NewRandomMirroredAzureDatabricksCatalog().Properties,
	}
}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsMirroredAzureDatabricksCatalog) CreateDefinition(
	data fabmirroredazuredatabrickscatalog.CreateMirroredAzureDatabricksCatalogRequest,
) *fabmirroredazuredatabrickscatalog.PublicDefinition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsMirroredAzureDatabricksCatalog) TransformDefinition(
	entity *fabmirroredazuredatabrickscatalog.PublicDefinition,
) fabmirroredazuredatabrickscatalog.ItemsClientGetMirroredAzureDatabricksCatalogDefinitionResponse {
	return fabmirroredazuredatabrickscatalog.ItemsClientGetMirroredAzureDatabricksCatalogDefinitionResponse{
		DefinitionResponse: fabmirroredazuredatabrickscatalog.DefinitionResponse{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsMirroredAzureDatabricksCatalog) UpdateDefinition(
	_ *fabmirroredazuredatabrickscatalog.PublicDefinition,
	data fabmirroredazuredatabrickscatalog.UpdatemirroredAzureDatabricksCatalogDefinitionRequest,
) *fabmirroredazuredatabrickscatalog.PublicDefinition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) CreateWithParentID(
	parentID string,
	data fabmirroredazuredatabrickscatalog.CreateMirroredAzureDatabricksCatalogRequest,
) fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog {
	entity := NewRandomMirroredAzureDatabricksCatalogWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description

	return entity
}

// Filter implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) Filter(
	entities []fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog,
	parentID string,
) []fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog {
	ret := make([]fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) GetID(entity fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) TransformCreate(
	entity fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog,
) fabmirroredazuredatabrickscatalog.ItemsClientCreateMirroredAzureDatabricksCatalogResponse {
	return fabmirroredazuredatabrickscatalog.ItemsClientCreateMirroredAzureDatabricksCatalogResponse{
		MirroredAzureDatabricksCatalog: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) TransformGet(
	entity fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog,
) fabmirroredazuredatabrickscatalog.ItemsClientGetMirroredAzureDatabricksCatalogResponse {
	return fabmirroredazuredatabrickscatalog.ItemsClientGetMirroredAzureDatabricksCatalogResponse{
		MirroredAzureDatabricksCatalog: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) TransformList(
	entities []fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog,
) fabmirroredazuredatabrickscatalog.ItemsClientListMirroredAzureDatabricksCatalogsResponse {
	return fabmirroredazuredatabrickscatalog.ItemsClientListMirroredAzureDatabricksCatalogsResponse{
		MirroredAzureDatabricksCatalogs: fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalogs{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) TransformUpdate(
	entity fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog,
) fabmirroredazuredatabrickscatalog.ItemsClientUpdateMirroredAzureDatabricksCatalogResponse {
	return fabmirroredazuredatabrickscatalog.ItemsClientUpdateMirroredAzureDatabricksCatalogResponse{
		MirroredAzureDatabricksCatalog: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) Update(
	base fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog,
	data fabmirroredazuredatabrickscatalog.UpdateMirroredAzureDatabricksCatalogRequest,
) fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog {
	base.DisplayName = data.DisplayName
	base.Description = data.Description

	return base
}

// Validate implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) Validate(
	newEntity fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog,
	existing []fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog,
) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureMirroredAzureDatabricksCatalog(server *fakeServer) fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog,
			fabmirroredazuredatabrickscatalog.ItemsClientGetMirroredAzureDatabricksCatalogResponse,
			fabmirroredazuredatabrickscatalog.ItemsClientUpdateMirroredAzureDatabricksCatalogResponse,
			fabmirroredazuredatabrickscatalog.ItemsClientCreateMirroredAzureDatabricksCatalogResponse,
			fabmirroredazuredatabrickscatalog.ItemsClientListMirroredAzureDatabricksCatalogsResponse,
			fabmirroredazuredatabrickscatalog.CreateMirroredAzureDatabricksCatalogRequest,
			fabmirroredazuredatabrickscatalog.UpdateMirroredAzureDatabricksCatalogRequest]
	}

	type concreteDefinitionOperations interface {
		definitionOperations[
			fabmirroredazuredatabrickscatalog.PublicDefinition,
			fabmirroredazuredatabrickscatalog.CreateMirroredAzureDatabricksCatalogRequest,
			fabmirroredazuredatabrickscatalog.UpdatemirroredAzureDatabricksCatalogDefinitionRequest,
			fabmirroredazuredatabrickscatalog.ItemsClientGetMirroredAzureDatabricksCatalogDefinitionResponse,
			fabmirroredazuredatabrickscatalog.ItemsClientUpdateMirroredAzureDatabricksCatalogDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsMirroredAzureDatabricksCatalog{}

	var definitionOperations concreteDefinitionOperations = &operationsMirroredAzureDatabricksCatalog{}

	var converter itemConverter[fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog] = &operationsMirroredAzureDatabricksCatalog{}

	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.MirroredAzureDatabricksCatalog.ItemsServer.GetMirroredAzureDatabricksCatalog,
		&server.ServerFactory.MirroredAzureDatabricksCatalog.ItemsServer.UpdateMirroredAzureDatabricksCatalog,
		&server.ServerFactory.MirroredAzureDatabricksCatalog.ItemsServer.BeginCreateMirroredAzureDatabricksCatalog,
		&server.ServerFactory.MirroredAzureDatabricksCatalog.ItemsServer.NewListMirroredAzureDatabricksCatalogsPager,
		&server.ServerFactory.MirroredAzureDatabricksCatalog.ItemsServer.DeleteMirroredAzureDatabricksCatalog)

	configureDefinitions(
		handler,
		entityOperations,
		definitionOperations,
		&server.ServerFactory.MirroredAzureDatabricksCatalog.ItemsServer.BeginCreateMirroredAzureDatabricksCatalog,
		&server.ServerFactory.MirroredAzureDatabricksCatalog.ItemsServer.BeginGetMirroredAzureDatabricksCatalogDefinition,
		&server.ServerFactory.MirroredAzureDatabricksCatalog.ItemsServer.BeginUpdateMirroredAzureDatabricksCatalogDefinition)

	return fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog{}
}

func NewRandomMirroredAzureDatabricksCatalog() fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog {
	return fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		FolderID:    to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fabmirroredazuredatabrickscatalog.ItemTypeMirroredAzureDatabricksCatalog),
		Properties: &fabmirroredazuredatabrickscatalog.Properties{
			AutoSync:                        to.Ptr(fabmirroredazuredatabrickscatalog.AutoSyncEnabled),
			CatalogName:                     to.Ptr(testhelp.RandomName()),
			DatabricksWorkspaceConnectionID: to.Ptr(testhelp.RandomUUID()),
			MirrorStatus:                    to.Ptr(fabmirroredazuredatabrickscatalog.MirrorStatusMirrored),
			MirroringMode:                   to.Ptr(fabmirroredazuredatabrickscatalog.MirroringModesFull),
			OneLakeTablesPath:               to.Ptr(testhelp.RandomURI()),
			SQLEndpointProperties: &fabmirroredazuredatabrickscatalog.SQLEndpointProperties{
				ID:               to.Ptr(testhelp.RandomUUID()),
				ConnectionString: to.Ptr(testhelp.RandomURI()),
			},
			StorageConnectionID: to.Ptr(testhelp.RandomUUID()),
			SyncDetails: &fabmirroredazuredatabrickscatalog.SyncDetails{
				LastSyncDateTime: to.Ptr(time.Now()),
				Status:           to.Ptr(fabmirroredazuredatabrickscatalog.StatusSuccess),
			},
		},
	}
}

func NewRandomMirroredAzureDatabricksCatalogWithWorkspace(workspaceID string) fabmirroredazuredatabrickscatalog.MirroredAzureDatabricksCatalog {
	result := NewRandomMirroredAzureDatabricksCatalog()
	result.WorkspaceID = &workspaceID

	return result
}

func NewRandomMirroredAzureDatabricksCatalogDefinition() fabmirroredazuredatabrickscatalog.PublicDefinition {
	defPart := fabmirroredazuredatabrickscatalog.PublicDefinitionPart{
		PayloadType: to.Ptr(fabmirroredazuredatabrickscatalog.PayloadTypeInlineBase64),
		Path:        to.Ptr("mirroringAzureDatabricksCatalog.json"),
		Payload:     to.Ptr("e30="),
	}

	var defParts []fabmirroredazuredatabrickscatalog.PublicDefinitionPart

	defParts = append(defParts, defPart)

	return fabmirroredazuredatabrickscatalog.PublicDefinition{
		Parts: defParts,
	}
}
