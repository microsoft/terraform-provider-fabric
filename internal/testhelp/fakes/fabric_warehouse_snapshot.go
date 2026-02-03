// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabwarehousesnapshot "github.com/microsoft/fabric-sdk-go/fabric/warehousesnapshot"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// operationsWarehouseSnapshot implements ParentIDOperations.
type operationsWarehouseSnapshot struct{}

func (o *operationsWarehouseSnapshot) CreateWithParentID(parentID string, data fabwarehousesnapshot.CreateWarehouseSnapshotRequest) fabwarehousesnapshot.WarehouseSnapshot {
	entity := NewRandomWarehouseSnapshotWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description
	entity.FolderID = data.FolderID

	return entity
}

// Filter implements ParentIDOperations.
func (o *operationsWarehouseSnapshot) Filter(entities []fabwarehousesnapshot.WarehouseSnapshot, parentID string) []fabwarehousesnapshot.WarehouseSnapshot {
	ret := make([]fabwarehousesnapshot.WarehouseSnapshot, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements ParentIDOperations.
func (o *operationsWarehouseSnapshot) GetID(entity fabwarehousesnapshot.WarehouseSnapshot) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements ParentIDOperations.
func (o *operationsWarehouseSnapshot) TransformCreate(entity fabwarehousesnapshot.WarehouseSnapshot) fabwarehousesnapshot.ItemsClientCreateWarehouseSnapshotResponse {
	return fabwarehousesnapshot.ItemsClientCreateWarehouseSnapshotResponse{
		WarehouseSnapshot: entity,
	}
}

// TransformGet implements ParentIDOperations.
func (o *operationsWarehouseSnapshot) TransformGet(entity fabwarehousesnapshot.WarehouseSnapshot) fabwarehousesnapshot.ItemsClientGetWarehouseSnapshotResponse {
	return fabwarehousesnapshot.ItemsClientGetWarehouseSnapshotResponse{
		WarehouseSnapshot: entity,
	}
}

// TransformList implements ParentIDOperations.
func (o *operationsWarehouseSnapshot) TransformList(entities []fabwarehousesnapshot.WarehouseSnapshot) fabwarehousesnapshot.ItemsClientListWarehouseSnapshotsResponse {
	return fabwarehousesnapshot.ItemsClientListWarehouseSnapshotsResponse{
		WarehouseSnapshots: fabwarehousesnapshot.WarehouseSnapshots{
			Value: entities,
		},
	}
}

// TransformUpdate implements ParentIDOperations.
func (o *operationsWarehouseSnapshot) TransformUpdate(entity fabwarehousesnapshot.WarehouseSnapshot) fabwarehousesnapshot.ItemsClientUpdateWarehouseSnapshotResponse {
	return fabwarehousesnapshot.ItemsClientUpdateWarehouseSnapshotResponse{
		WarehouseSnapshot: entity,
	}
}

// Update implements ParentIDOperations.
func (o *operationsWarehouseSnapshot) Update(base fabwarehousesnapshot.WarehouseSnapshot, data fabwarehousesnapshot.UpdateWarehouseSnapshotRequest) fabwarehousesnapshot.WarehouseSnapshot {
	base.DisplayName = data.DisplayName
	base.Description = data.Description

	return base
}

// Validate implements ParentIDOperations.
func (o *operationsWarehouseSnapshot) Validate(newEntity fabwarehousesnapshot.WarehouseSnapshot, existing []fabwarehousesnapshot.WarehouseSnapshot) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

// ConvertItemToEntity implements itemConverter.
func (o *operationsWarehouseSnapshot) ConvertItemToEntity(entity fabcore.Item) fabwarehousesnapshot.WarehouseSnapshot {
	return fabwarehousesnapshot.WarehouseSnapshot{
		ID:          entity.ID,
		DisplayName: entity.DisplayName,
		Description: entity.Description,
		WorkspaceID: entity.WorkspaceID,
		FolderID:    entity.FolderID,
		Type:        to.Ptr(fabwarehousesnapshot.ItemTypeWarehouseSnapshot),
		Properties:  NewRandomWarehouseSnapshot().Properties,
	}
}

func configureWarehouseSnapshot(server *fakeServer) fabwarehousesnapshot.WarehouseSnapshot {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabwarehousesnapshot.WarehouseSnapshot,
			fabwarehousesnapshot.ItemsClientGetWarehouseSnapshotResponse,
			fabwarehousesnapshot.ItemsClientUpdateWarehouseSnapshotResponse,
			fabwarehousesnapshot.ItemsClientCreateWarehouseSnapshotResponse,
			fabwarehousesnapshot.ItemsClientListWarehouseSnapshotsResponse,
			fabwarehousesnapshot.CreateWarehouseSnapshotRequest,
			fabwarehousesnapshot.UpdateWarehouseSnapshotRequest]
	}

	var entityOperations concreteEntityOperations = &operationsWarehouseSnapshot{}
	var converter itemConverter[fabwarehousesnapshot.WarehouseSnapshot] = &operationsWarehouseSnapshot{}

	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.WarehouseSnapshot.ItemsServer.GetWarehouseSnapshot,
		&server.ServerFactory.WarehouseSnapshot.ItemsServer.UpdateWarehouseSnapshot,
		&server.ServerFactory.WarehouseSnapshot.ItemsServer.BeginCreateWarehouseSnapshot,
		&server.ServerFactory.WarehouseSnapshot.ItemsServer.NewListWarehouseSnapshotsPager,
		&server.ServerFactory.WarehouseSnapshot.ItemsServer.DeleteWarehouseSnapshot)

	return fabwarehousesnapshot.WarehouseSnapshot{}
}

func NewRandomWarehouseSnapshot() fabwarehousesnapshot.WarehouseSnapshot {
	return fabwarehousesnapshot.WarehouseSnapshot{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		FolderID:    to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fabwarehousesnapshot.ItemTypeWarehouseSnapshot),
		Properties: &fabwarehousesnapshot.Properties{
			ConnectionString:  to.Ptr(testhelp.RandomURI()),
			ParentWarehouseID: to.Ptr(testhelp.RandomUUID()),
			SnapshotDateTime:  to.Ptr(time.Now()),
		},
	}
}

func NewRandomWarehouseSnapshotWithWorkspace(workspaceID string) fabwarehousesnapshot.WarehouseSnapshot {
	result := NewRandomWarehouseSnapshot()
	result.WorkspaceID = &workspaceID

	return result
}
