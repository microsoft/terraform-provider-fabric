// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabwarehouse "github.com/microsoft/fabric-sdk-go/fabric/warehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// operationsWarehouse implements ParentIDOperations.
type operationsWarehouse struct{}

func (o *operationsWarehouse) CreateWithParentID(parentID string, data fabwarehouse.CreateWarehouseRequest) fabwarehouse.Warehouse {
	entity := NewRandomWarehouseWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description

	return entity
}

// Filter implements ParentIDOperations.
func (o *operationsWarehouse) Filter(entities []fabwarehouse.Warehouse, parentID string) []fabwarehouse.Warehouse {
	ret := make([]fabwarehouse.Warehouse, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements ParentIDOperations.
func (o *operationsWarehouse) GetID(entity fabwarehouse.Warehouse) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements ParentIDOperations.
func (o *operationsWarehouse) TransformCreate(entity fabwarehouse.Warehouse) fabwarehouse.ItemsClientCreateWarehouseResponse {
	return fabwarehouse.ItemsClientCreateWarehouseResponse{
		Warehouse: entity,
	}
}

// TransformGet implements ParentIDOperations.
func (o *operationsWarehouse) TransformGet(entity fabwarehouse.Warehouse) fabwarehouse.ItemsClientGetWarehouseResponse {
	return fabwarehouse.ItemsClientGetWarehouseResponse{
		Warehouse: entity,
	}
}

// TransformList implements ParentIDOperations.
func (o *operationsWarehouse) TransformList(entities []fabwarehouse.Warehouse) fabwarehouse.ItemsClientListWarehousesResponse {
	return fabwarehouse.ItemsClientListWarehousesResponse{
		Warehouses: fabwarehouse.Warehouses{
			Value: entities,
		},
	}
}

// TransformUpdate implements ParentIDOperations.
func (o *operationsWarehouse) TransformUpdate(entity fabwarehouse.Warehouse) fabwarehouse.ItemsClientUpdateWarehouseResponse {
	return fabwarehouse.ItemsClientUpdateWarehouseResponse{
		Warehouse: entity,
	}
}

// Update implements ParentIDOperations.
func (o *operationsWarehouse) Update(base fabwarehouse.Warehouse, data fabwarehouse.UpdateWarehouseRequest) fabwarehouse.Warehouse {
	base.DisplayName = data.DisplayName
	base.Description = data.Description

	return base
}

// Validate implements ParentIDOperations.
func (o *operationsWarehouse) Validate(newEntity fabwarehouse.Warehouse, existing []fabwarehouse.Warehouse) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

// ConvertItemToEntity implements itemConverter.
func (o *operationsWarehouse) ConvertItemToEntity(entity fabcore.Item) fabwarehouse.Warehouse {
	return fabwarehouse.Warehouse{
		ID:          entity.ID,
		DisplayName: entity.DisplayName,
		Description: entity.Description,
		WorkspaceID: entity.WorkspaceID,
		Type:        to.Ptr(fabwarehouse.ItemTypeWarehouse),
		Properties:  NewRandomWarehouse().Properties,
	}
}

func configureWarehouse(server *fakeServer) fabwarehouse.Warehouse {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabwarehouse.Warehouse,
			fabwarehouse.ItemsClientGetWarehouseResponse,
			fabwarehouse.ItemsClientUpdateWarehouseResponse,
			fabwarehouse.ItemsClientCreateWarehouseResponse,
			fabwarehouse.ItemsClientListWarehousesResponse,
			fabwarehouse.CreateWarehouseRequest,
			fabwarehouse.UpdateWarehouseRequest]
	}

	var entityOperations concreteEntityOperations = &operationsWarehouse{}
	var converter itemConverter[fabwarehouse.Warehouse] = &operationsWarehouse{}

	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.Warehouse.ItemsServer.GetWarehouse,
		&server.ServerFactory.Warehouse.ItemsServer.UpdateWarehouse,
		&server.ServerFactory.Warehouse.ItemsServer.BeginCreateWarehouse,
		&server.ServerFactory.Warehouse.ItemsServer.NewListWarehousesPager,
		&server.ServerFactory.Warehouse.ItemsServer.DeleteWarehouse)

	return fabwarehouse.Warehouse{}
}

func NewRandomWarehouse() fabwarehouse.Warehouse {
	return fabwarehouse.Warehouse{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fabwarehouse.ItemTypeWarehouse),
		Properties: &fabwarehouse.Properties{
			CollationType:    to.Ptr(fabwarehouse.CollationTypeLatin1General100BIN2UTF8),
			ConnectionString: to.Ptr(testhelp.RandomURI()),
			CreatedDate:      to.Ptr(time.Now()),
			LastUpdatedTime:  to.Ptr(time.Now()),
			CollationType:    to.Ptr(fabwarehouse.CollationTypeLatin1General100CIASKSWSSCUTF8),
		},
	}
}

func NewRandomWarehouseWithWorkspace(workspaceID string) fabwarehouse.Warehouse {
	result := NewRandomWarehouse()
	result.WorkspaceID = &workspaceID

	return result
}
