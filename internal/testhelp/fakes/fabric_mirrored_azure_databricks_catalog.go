// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabMirroredAzureDatabricksCatalog "github.com/microsoft/fabric-sdk-go/fabric/mirroredazuredatabrickscatalog"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsMirroredAzureDatabricksCatalog struct{}

// ConvertItemToEntity implements itemConverter.
func (o *operationsMirroredAzureDatabricksCatalog) ConvertItemToEntity(item fabcore.Item) fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog {
	return fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog{
		ID:          item.ID,
		DisplayName: item.DisplayName,
		Description: item.Description,
		WorkspaceID: item.WorkspaceID,
		Type:        to.Ptr(fabMirroredAzureDatabricksCatalog.ItemTypeMirroredAzureDatabricksCatalog),
		Properties:  NewRandomMirroredAzureDatabricksCatalog().Properties,
	}
}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsMirroredAzureDatabricksCatalog) CreateDefinition(
	data fabMirroredAzureDatabricksCatalog.CreateMirroredAzureDatabricksCatalogRequest,
) *fabMirroredAzureDatabricksCatalog.PublicDefinition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsMirroredAzureDatabricksCatalog) TransformDefinition(
	entity *fabMirroredAzureDatabricksCatalog.PublicDefinition,
) fabMirroredAzureDatabricksCatalog.ItemsClientGetMirroredAzureDatabricksCatalogDefinitionResponse {
	return fabMirroredAzureDatabricksCatalog.ItemsClientGetMirroredAzureDatabricksCatalogDefinitionResponse{
		DefinitionResponse: fabMirroredAzureDatabricksCatalog.DefinitionResponse{
			Definition: entity,
		},
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsMirroredAzureDatabricksCatalog) UpdateDefinition(
	_ *fabMirroredAzureDatabricksCatalog.PublicDefinition,
	data fabMirroredAzureDatabricksCatalog.UpdatemirroredAzureDatabricksCatalogDefinitionRequest,
) *fabMirroredAzureDatabricksCatalog.PublicDefinition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) CreateWithParentID(
	parentID string,
	data fabMirroredAzureDatabricksCatalog.CreateMirroredAzureDatabricksCatalogRequest,
) fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog {
	entity := NewRandomMirroredAzureDatabricksCatalogWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description

	return entity
}

// Filter implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) Filter(
	entities []fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog,
	parentID string,
) []fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog {
	ret := make([]fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) GetID(entity fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) TransformCreate(
	entity fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog,
) fabMirroredAzureDatabricksCatalog.ItemsClientCreateMirroredAzureDatabricksCatalogResponse {
	return fabMirroredAzureDatabricksCatalog.ItemsClientCreateMirroredAzureDatabricksCatalogResponse{
		MirroredAzureDatabricksCatalog: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) TransformGet(
	entity fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog,
) fabMirroredAzureDatabricksCatalog.ItemsClientGetMirroredAzureDatabricksCatalogResponse {
	return fabMirroredAzureDatabricksCatalog.ItemsClientGetMirroredAzureDatabricksCatalogResponse{
		MirroredAzureDatabricksCatalog: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) TransformList(
	entities []fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog,
) fabMirroredAzureDatabricksCatalog.ItemsClientListMirroredAzureDatabricksCatalogsResponse {
	return fabMirroredAzureDatabricksCatalog.ItemsClientListMirroredAzureDatabricksCatalogsResponse{
		MirroredAzureDatabricksCatalogs: fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalogs{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) TransformUpdate(
	entity fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog,
) fabMirroredAzureDatabricksCatalog.ItemsClientUpdateMirroredAzureDatabricksCatalogResponse {
	return fabMirroredAzureDatabricksCatalog.ItemsClientUpdateMirroredAzureDatabricksCatalogResponse{
		MirroredAzureDatabricksCatalog: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) Update(
	base fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog,
	data fabMirroredAzureDatabricksCatalog.UpdateMirroredAzureDatabricksCatalogRequest,
) fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog {
	base.DisplayName = data.DisplayName
	base.Description = data.Description

	return base
}

// Validate implements concreteOperations.
func (o *operationsMirroredAzureDatabricksCatalog) Validate(
	newEntity fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog,
	existing []fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog,
) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureMirroredAzureDatabricksCatalog(server *fakeServer) fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog,
			fabMirroredAzureDatabricksCatalog.ItemsClientGetMirroredAzureDatabricksCatalogResponse,
			fabMirroredAzureDatabricksCatalog.ItemsClientUpdateMirroredAzureDatabricksCatalogResponse,
			fabMirroredAzureDatabricksCatalog.ItemsClientCreateMirroredAzureDatabricksCatalogResponse,
			fabMirroredAzureDatabricksCatalog.ItemsClientListMirroredAzureDatabricksCatalogsResponse,
			fabMirroredAzureDatabricksCatalog.CreateMirroredAzureDatabricksCatalogRequest,
			fabMirroredAzureDatabricksCatalog.UpdateMirroredAzureDatabricksCatalogRequest]
	}

	type concreteDefinitionOperations interface {
		definitionOperations[
			fabMirroredAzureDatabricksCatalog.PublicDefinition,
			fabMirroredAzureDatabricksCatalog.CreateMirroredAzureDatabricksCatalogRequest,
			fabMirroredAzureDatabricksCatalog.UpdatemirroredAzureDatabricksCatalogDefinitionRequest,
			fabMirroredAzureDatabricksCatalog.ItemsClientGetMirroredAzureDatabricksCatalogDefinitionResponse,
			fabMirroredAzureDatabricksCatalog.ItemsClientUpdateMirroredAzureDatabricksCatalogDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsMirroredAzureDatabricksCatalog{}

	var definitionOperations concreteDefinitionOperations = &operationsMirroredAzureDatabricksCatalog{}

	var converter itemConverter[fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog] = &operationsMirroredAzureDatabricksCatalog{}

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

	return fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog{}
}

func NewRandomMirroredAzureDatabricksCatalog() fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog {
	return fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fabMirroredAzureDatabricksCatalog.ItemTypeMirroredAzureDatabricksCatalog),
		Properties: &fabMirroredAzureDatabricksCatalog.Properties{
			AutoSync:                        to.Ptr(fabMirroredAzureDatabricksCatalog.AutoSyncEnabled),
			CatalogName:                     to.Ptr(testhelp.RandomName()),
			DatabricksWorkspaceConnectionID: to.Ptr(testhelp.RandomUUID()),
			MirrorStatus:                    to.Ptr(fabMirroredAzureDatabricksCatalog.MirrorStatusMirrored),
			MirroringMode:                   to.Ptr(fabMirroredAzureDatabricksCatalog.MirroringModesFull),
			OneLakeTablesPath:               to.Ptr(testhelp.RandomName()),
			SQLEndpointProperties: &fabMirroredAzureDatabricksCatalog.SQLEndpointProperties{
				ID:               to.Ptr(testhelp.RandomUUID()),
				ConnectionString: to.Ptr(testhelp.RandomURI()),
			},
			StorageConnectionID: to.Ptr(testhelp.RandomUUID()),
			SyncDetails: &fabMirroredAzureDatabricksCatalog.SyncDetails{
				LastSyncDateTime: to.Ptr(time.Now()),
				Status:           to.Ptr(fabMirroredAzureDatabricksCatalog.StatusSuccess),
				ErrorInfo: &fabMirroredAzureDatabricksCatalog.ErrorInfo{
					ErrorCode:    to.Ptr(testhelp.RandomName()),
					ErrorDetails: to.Ptr(testhelp.RandomName()),
					ErrorMessage: to.Ptr(testhelp.RandomName()),
				},
			},
		},
	}
}

func NewRandomMirroredAzureDatabricksCatalogWithWorkspace(workspaceID string) fabMirroredAzureDatabricksCatalog.MirroredAzureDatabricksCatalog {
	result := NewRandomMirroredAzureDatabricksCatalog()
	result.WorkspaceID = &workspaceID

	return result
}

func NewRandomMirroredAzureDatabricksCatalogDefinition() fabMirroredAzureDatabricksCatalog.PublicDefinition {
	defPart := fabMirroredAzureDatabricksCatalog.PublicDefinitionPart{
		PayloadType: to.Ptr(fabMirroredAzureDatabricksCatalog.PayloadTypeInlineBase64),
		Path:        to.Ptr("mirroringAzureDatabricksCatalog.json"),
		Payload:     to.Ptr("e30="),
	}

	var defParts []fabMirroredAzureDatabricksCatalog.PublicDefinitionPart

	defParts = append(defParts, defPart)

	return fabMirroredAzureDatabricksCatalog.PublicDefinition{
		Parts: defParts,
	}
}
