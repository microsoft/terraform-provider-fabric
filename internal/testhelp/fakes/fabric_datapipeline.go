// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabdatapipeline "github.com/microsoft/fabric-sdk-go/fabric/datapipeline"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsDataPipeline struct{}

// CreateWithParentID implements concreteOperations.
func (o *operationsDataPipeline) CreateWithParentID(parentID string, data fabdatapipeline.CreateDataPipelineRequest) fabdatapipeline.DataPipeline {
	entity := NewRandomDataPipelineWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description

	return entity
}

// Filter implements concreteOperations.
func (o *operationsDataPipeline) Filter(entities []fabdatapipeline.DataPipeline, parentID string) []fabdatapipeline.DataPipeline {
	ret := make([]fabdatapipeline.DataPipeline, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsDataPipeline) GetID(entity fabdatapipeline.DataPipeline) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsDataPipeline) TransformCreate(entity fabdatapipeline.DataPipeline) fabdatapipeline.ItemsClientCreateDataPipelineResponse {
	return fabdatapipeline.ItemsClientCreateDataPipelineResponse{
		DataPipeline: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsDataPipeline) TransformGet(entity fabdatapipeline.DataPipeline) fabdatapipeline.ItemsClientGetDataPipelineResponse {
	return fabdatapipeline.ItemsClientGetDataPipelineResponse{
		DataPipeline: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsDataPipeline) TransformList(entities []fabdatapipeline.DataPipeline) fabdatapipeline.ItemsClientListDataPipelinesResponse {
	return fabdatapipeline.ItemsClientListDataPipelinesResponse{
		DataPipelines: fabdatapipeline.DataPipelines{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsDataPipeline) TransformUpdate(entity fabdatapipeline.DataPipeline) fabdatapipeline.ItemsClientUpdateDataPipelineResponse {
	return fabdatapipeline.ItemsClientUpdateDataPipelineResponse{
		DataPipeline: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsDataPipeline) Update(base fabdatapipeline.DataPipeline, data fabdatapipeline.UpdateDataPipelineRequest) fabdatapipeline.DataPipeline {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsDataPipeline) Validate(newEntity fabdatapipeline.DataPipeline, existing []fabdatapipeline.DataPipeline) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureDataPipeline(server *fakeServer) fabdatapipeline.DataPipeline {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabdatapipeline.DataPipeline,
			fabdatapipeline.ItemsClientGetDataPipelineResponse,
			fabdatapipeline.ItemsClientUpdateDataPipelineResponse,
			fabdatapipeline.ItemsClientCreateDataPipelineResponse,
			fabdatapipeline.ItemsClientListDataPipelinesResponse,
			fabdatapipeline.CreateDataPipelineRequest,
			fabdatapipeline.UpdateDataPipelineRequest]
	}

	var entityOperations concreteEntityOperations = &operationsDataPipeline{}

	handler := newTypedHandler(server, entityOperations)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.DataPipeline.ItemsServer.GetDataPipeline,
		&server.ServerFactory.DataPipeline.ItemsServer.UpdateDataPipeline,
		&server.ServerFactory.DataPipeline.ItemsServer.BeginCreateDataPipeline,
		&server.ServerFactory.DataPipeline.ItemsServer.NewListDataPipelinesPager,
		&server.ServerFactory.DataPipeline.ItemsServer.DeleteDataPipeline)

	return fabdatapipeline.DataPipeline{}
}

func NewRandomDataPipeline() fabdatapipeline.DataPipeline {
	return fabdatapipeline.DataPipeline{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fabdatapipeline.ItemTypeDataPipeline),
	}
}

func NewRandomDataPipelineWithWorkspace(workspaceID string) fabdatapipeline.DataPipeline {
	result := NewRandomDataPipeline()
	result.WorkspaceID = &workspaceID

	return result
}
