// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"context"
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabmonitoring "github.com/microsoft/fabric-sdk-go/fabric/monitoring"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsMonitoring struct{}

// ConvertItemToEntity implements itemConverter.
func (o *operationsMonitoring) ConvertItemToEntity(item fabcore.Item) fabmonitoring.Monitoring {
	return fabmonitoring.Monitoring{
		ID:          item.ID,
		DisplayName: item.DisplayName,
		Description: item.Description,
		WorkspaceID: item.WorkspaceID,
		FolderID:    item.FolderID,
		Type:        to.Ptr(fabmonitoring.ItemTypeMonitoring),
		Properties:  NewRandomMonitoring().Properties,
		Tags:        convertItemTags[fabmonitoring.ItemTag](item.Tags),
	}
}

// CreateDefinition implements concreteDefinitionOperations.
func (o *operationsMonitoring) CreateDefinition(data fabcore.CreateItemRequest) *fabcore.ItemDefinition {
	return data.Definition
}

// TransformDefinition implements concreteDefinitionOperations.
func (o *operationsMonitoring) TransformDefinition(entity *fabcore.ItemDefinition) fabcore.ItemsClientGetItemDefinitionResponse {
	return fabcore.ItemsClientGetItemDefinitionResponse{
		Definition: entity,
	}
}

// UpdateDefinition implements concreteDefinitionOperations.
func (o *operationsMonitoring) UpdateDefinition(_ *fabcore.ItemDefinition, data fabcore.UpdateItemDefinitionRequest) *fabcore.ItemDefinition {
	return data.Definition
}

// CreateWithParentID implements concreteOperations.
func (o *operationsMonitoring) CreateWithParentID(parentID string, data fabmonitoring.CreateMonitoringRequest) fabmonitoring.Monitoring {
	entity := NewRandomMonitoringWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description
	entity.FolderID = data.FolderID

	return entity
}

// Filter implements concreteOperations.
func (o *operationsMonitoring) Filter(entities []fabmonitoring.Monitoring, parentID string) []fabmonitoring.Monitoring {
	ret := make([]fabmonitoring.Monitoring, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsMonitoring) GetID(entity fabmonitoring.Monitoring) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsMonitoring) TransformCreate(entity fabmonitoring.Monitoring) fabmonitoring.ItemsClientCreateMonitoringResponse {
	return fabmonitoring.ItemsClientCreateMonitoringResponse{
		Monitoring: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsMonitoring) TransformGet(entity fabmonitoring.Monitoring) fabmonitoring.ItemsClientGetMonitoringResponse {
	return fabmonitoring.ItemsClientGetMonitoringResponse{
		Monitoring: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsMonitoring) TransformList(entities []fabmonitoring.Monitoring) fabmonitoring.ItemsClientListMonitoringItemsResponse {
	return fabmonitoring.ItemsClientListMonitoringItemsResponse{
		Value: entities,
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsMonitoring) TransformUpdate(entity fabmonitoring.Monitoring) fabmonitoring.ItemsClientUpdateMonitoringResponse {
	return fabmonitoring.ItemsClientUpdateMonitoringResponse{
		Monitoring: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsMonitoring) Update(base fabmonitoring.Monitoring, data fabmonitoring.UpdateMonitoringRequest) fabmonitoring.Monitoring {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsMonitoring) Validate(newEntity fabmonitoring.Monitoring, existing []fabmonitoring.Monitoring) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureMonitoring(server *fakeServer) fabmonitoring.Monitoring {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabmonitoring.Monitoring,
			fabmonitoring.ItemsClientGetMonitoringResponse,
			fabmonitoring.ItemsClientUpdateMonitoringResponse,
			fabmonitoring.ItemsClientCreateMonitoringResponse,
			fabmonitoring.ItemsClientListMonitoringItemsResponse,
			fabmonitoring.CreateMonitoringRequest,
			fabmonitoring.UpdateMonitoringRequest]
	}
	type concreteDefinitionOperations interface {
		definitionOperationsNonLROCreation[
			fabcore.ItemDefinition,
			fabcore.UpdateItemDefinitionRequest,
			fabcore.ItemsClientGetItemDefinitionResponse,
			fabcore.ItemsClientUpdateItemDefinitionResponse]
	}

	var entityOperations concreteEntityOperations = &operationsMonitoring{}
	var definitionOperations concreteDefinitionOperations = &operationsMonitoring{}
	var converter itemConverter[fabmonitoring.Monitoring] = &operationsMonitoring{}

	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.Monitoring.ItemsServer.GetMonitoring,
		&server.ServerFactory.Monitoring.ItemsServer.UpdateMonitoring,
		&server.ServerFactory.Monitoring.ItemsServer.BeginCreateMonitoring,
		&server.ServerFactory.Monitoring.ItemsServer.NewListMonitoringItemsPager,
		&server.ServerFactory.Monitoring.ItemsServer.DeleteMonitoring,
	)
	configureDefinitionsNonLROCreation(
		handler,
		definitionOperations,
		&server.ServerFactory.Core.ItemsServer.BeginGetItemDefinition,
		&server.ServerFactory.Core.ItemsServer.BeginUpdateItemDefinition,
	)

	server.ServerFactory.Monitoring.DataCollectionServer.Enable = func(_ context.Context, workspaceID, monitoringID string, _ *fabmonitoring.DataCollectionClientEnableOptions) (azfake.Responder[fabmonitoring.DataCollectionClientEnableResponse], azfake.ErrorResponder) {
		var resp azfake.Responder[fabmonitoring.DataCollectionClientEnableResponse]

		var errResp azfake.ErrorResponder

		err := setMonitoringCollectionState(handler, workspaceID, monitoringID, fabmonitoring.DataCollectionStateEnabled)
		if err != nil {
			errResp.SetError(err)
		} else {
			resp.SetResponse(http.StatusOK, fabmonitoring.DataCollectionClientEnableResponse{}, nil)
		}

		return resp, errResp
	}
	server.ServerFactory.Monitoring.DataCollectionServer.Disable = func(_ context.Context, workspaceID, monitoringID string, _ *fabmonitoring.DataCollectionClientDisableOptions) (azfake.Responder[fabmonitoring.DataCollectionClientDisableResponse], azfake.ErrorResponder) {
		var resp azfake.Responder[fabmonitoring.DataCollectionClientDisableResponse]

		var errResp azfake.ErrorResponder

		err := setMonitoringCollectionState(handler, workspaceID, monitoringID, fabmonitoring.DataCollectionStateDisabled)
		if err != nil {
			errResp.SetError(err)
		} else {
			resp.SetResponse(http.StatusOK, fabmonitoring.DataCollectionClientDisableResponse{}, nil)
		}

		return resp, errResp
	}

	return fabmonitoring.Monitoring{}
}

func NewRandomMonitoring() fabmonitoring.Monitoring {
	return fabmonitoring.Monitoring{
		ID:          new(testhelp.RandomUUID()),
		DisplayName: new(testhelp.RandomName()),
		Description: new(testhelp.RandomName()),
		WorkspaceID: new(testhelp.RandomUUID()),
		FolderID:    new(testhelp.RandomUUID()),
		Type:        to.Ptr(fabmonitoring.ItemTypeMonitoring),
		Properties: &fabmonitoring.Properties{
			DataCollectionState: to.Ptr(fabmonitoring.DataCollectionStateDisabled),
		},
	}
}

func NewRandomMonitoringWithWorkspace(workspaceID string) fabmonitoring.Monitoring {
	result := NewRandomMonitoring()
	result.WorkspaceID = &workspaceID

	return result
}

func setMonitoringCollectionState(handler *typedHandler[fabmonitoring.Monitoring], workspaceID, monitoringID string, state fabmonitoring.DataCollectionState) error {
	id := generateID(workspaceID, monitoringID)
	if !handler.Contains(id) {
		return fabfake.SetResponseError(http.StatusNotFound, "ItemNotFound", "Monitoring item not found")
	}

	entity := handler.Get(id)
	entity.Properties = &fabmonitoring.Properties{DataCollectionState: new(state)}
	handler.Upsert(entity)

	return nil
}
