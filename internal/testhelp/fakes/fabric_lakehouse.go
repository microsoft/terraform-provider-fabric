// Copyright (c) Microsoft Corporation
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

// CreateWithParentID implements concreteOperations.
func (o *operationsLakehouse) CreateWithParentID(parentID string, data fablakehouse.CreateLakehouseRequest) fablakehouse.Lakehouse {
	entity := NewRandomLakehouseWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description

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

	var entityOperations concreteEntityOperations = &operationsLakehouse{}

	handler := newTypedHandler(server, entityOperations)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.Lakehouse.ItemsServer.GetLakehouse,
		&server.ServerFactory.Lakehouse.ItemsServer.UpdateLakehouse,
		&server.ServerFactory.Lakehouse.ItemsServer.BeginCreateLakehouse,
		&server.ServerFactory.Lakehouse.ItemsServer.NewListLakehousesPager,
		&server.ServerFactory.Lakehouse.ItemsServer.DeleteLakehouse)

	return fablakehouse.Lakehouse{}
}

func NewRandomLakehouse() fablakehouse.Lakehouse {
	return fablakehouse.Lakehouse{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fablakehouse.ItemTypeLakehouse),
		Properties: &fablakehouse.Properties{
			OneLakeFilesPath:  to.Ptr(testhelp.RandomName()),
			OneLakeTablesPath: to.Ptr(testhelp.RandomName()),
			SQLEndpointProperties: &fablakehouse.SQLEndpointProperties{
				ID:                 to.Ptr(testhelp.RandomUUID()),
				ProvisioningStatus: to.Ptr(fablakehouse.SQLEndpointProvisioningStatusSuccess),
				ConnectionString:   to.Ptr(testhelp.RandomName()),
			},
		},
	}
}

func NewRandomLakehouseWithWorkspace(workspaceID string) fablakehouse.Lakehouse {
	result := NewRandomLakehouse()
	result.WorkspaceID = &workspaceID

	return result
}
