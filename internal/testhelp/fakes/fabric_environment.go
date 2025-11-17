// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

type operationsEnvironment struct{}

// ConvertItemToEntity implements itemConverter.
func (o *operationsEnvironment) ConvertItemToEntity(item fabcore.Item) fabenvironment.Environment {
	return fabenvironment.Environment{
		ID:          item.ID,
		DisplayName: item.DisplayName,
		Description: item.Description,
		WorkspaceID: item.WorkspaceID,
		Type:        to.Ptr(fabenvironment.ItemTypeEnvironment),
		Properties:  NewRandomEnvironment().Properties,
	}
}

// CreateWithParentID implements concreteOperations.
func (o *operationsEnvironment) CreateWithParentID(parentID string, data fabenvironment.CreateEnvironmentRequest) fabenvironment.Environment {
	entity := NewRandomEnvironmentWithWorkspace(parentID)
	entity.DisplayName = data.DisplayName
	entity.Description = data.Description

	return entity
}

// Filter implements concreteOperations.
func (o *operationsEnvironment) Filter(entities []fabenvironment.Environment, parentID string) []fabenvironment.Environment {
	ret := make([]fabenvironment.Environment, 0)

	for _, entity := range entities {
		if *entity.WorkspaceID == parentID {
			ret = append(ret, entity)
		}
	}

	return ret
}

// GetID implements concreteOperations.
func (o *operationsEnvironment) GetID(entity fabenvironment.Environment) string {
	return generateID(*entity.WorkspaceID, *entity.ID)
}

// TransformCreate implements concreteOperations.
func (o *operationsEnvironment) TransformCreate(entity fabenvironment.Environment) fabenvironment.ItemsClientCreateEnvironmentResponse {
	return fabenvironment.ItemsClientCreateEnvironmentResponse{
		Environment: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsEnvironment) TransformGet(entity fabenvironment.Environment) fabenvironment.ItemsClientGetEnvironmentResponse {
	return fabenvironment.ItemsClientGetEnvironmentResponse{
		Environment: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsEnvironment) TransformList(entities []fabenvironment.Environment) fabenvironment.ItemsClientListEnvironmentsResponse {
	return fabenvironment.ItemsClientListEnvironmentsResponse{
		Environments: fabenvironment.Environments{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsEnvironment) TransformUpdate(entity fabenvironment.Environment) fabenvironment.ItemsClientUpdateEnvironmentResponse {
	return fabenvironment.ItemsClientUpdateEnvironmentResponse{
		Environment: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsEnvironment) Update(base fabenvironment.Environment, data fabenvironment.UpdateEnvironmentRequest) fabenvironment.Environment {
	base.Description = data.Description
	base.DisplayName = data.DisplayName

	return base
}

// Validate implements concreteOperations.
func (o *operationsEnvironment) Validate(newEntity fabenvironment.Environment, existing []fabenvironment.Environment) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureEnvironment(server *fakeServer) fabenvironment.Environment {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabenvironment.Environment,
			fabenvironment.ItemsClientGetEnvironmentResponse,
			fabenvironment.ItemsClientUpdateEnvironmentResponse,
			fabenvironment.ItemsClientCreateEnvironmentResponse,
			fabenvironment.ItemsClientListEnvironmentsResponse,
			fabenvironment.CreateEnvironmentRequest,
			fabenvironment.UpdateEnvironmentRequest]
	}

	var entityOperations concreteEntityOperations = &operationsEnvironment{}
	var converter itemConverter[fabenvironment.Environment] = &operationsEnvironment{}
	handler := newTypedHandlerWithConverter(server, entityOperations, converter)

	configureEntityWithParentID(
		handler,
		entityOperations,
		&server.ServerFactory.Environment.ItemsServer.GetEnvironment,
		&server.ServerFactory.Environment.ItemsServer.UpdateEnvironment,
		&server.ServerFactory.Environment.ItemsServer.BeginCreateEnvironment,
		&server.ServerFactory.Environment.ItemsServer.NewListEnvironmentsPager,
		&server.ServerFactory.Environment.ItemsServer.DeleteEnvironment)

	return fabenvironment.Environment{}
}

func NewRandomEnvironment() fabenvironment.Environment {
	return fabenvironment.Environment{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Description: to.Ptr(testhelp.RandomName()),
		WorkspaceID: to.Ptr(testhelp.RandomUUID()),
		Type:        to.Ptr(fabenvironment.ItemTypeEnvironment),
		Properties: &fabenvironment.Properties{
			PublishDetails: &fabenvironment.PublishDetails{
				State:         to.Ptr(fabenvironment.PublishStateSuccess),
				TargetVersion: to.Ptr(testhelp.RandomUUID()),
				StartTime:     to.Ptr(time.Now()),
				EndTime:       to.Ptr(time.Now()),
				ComponentPublishInfo: &fabenvironment.ComponentPublishInfo{
					SparkLibraries: &fabenvironment.SparkLibraries{
						State: to.Ptr(fabenvironment.PublishStateSuccess),
					},
					SparkSettings: &fabenvironment.SparkSettings{
						State: to.Ptr(fabenvironment.PublishStateSuccess),
					},
				},
			},
		},
	}
}

func NewRandomEnvironmentWithWorkspace(workspaceID string) fabenvironment.Environment {
	result := NewRandomEnvironment()
	result.WorkspaceID = &workspaceID

	return result
}
