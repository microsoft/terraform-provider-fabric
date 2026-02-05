// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// operationsWorkspaceManagedPrivateEndpoint implements ParentIDOperations.
type operationsWorkspaceManagedPrivateEndpoint struct{}

// GetID implements concreteOperations.
func (o *operationsWorkspaceManagedPrivateEndpoint) GetID(entity fabcore.ManagedPrivateEndpoint) string {
	return *entity.ID
}

func (o *operationsWorkspaceManagedPrivateEndpoint) GetIDWithParentID(_ string, entity fabcore.ManagedPrivateEndpoint) string {
	return *entity.ID
}

func (o *operationsWorkspaceManagedPrivateEndpoint) Create(data fabcore.CreateManagedPrivateEndpointRequest) fabcore.ManagedPrivateEndpoint {
	entity := NewRandomWorkspaceManagedPrivateEndpoint()
	entity.Name = data.Name
	entity.TargetPrivateLinkResourceID = data.TargetPrivateLinkResourceID
	entity.TargetSubresourceType = data.TargetSubresourceType

	return entity
}

// CreateWithParentID implements concreteOperations.
func (o *operationsWorkspaceManagedPrivateEndpoint) CreateWithParentID(_ string, data fabcore.CreateManagedPrivateEndpointRequest) fabcore.ManagedPrivateEndpoint {
	entity := NewRandomWorkspaceManagedPrivateEndpoint()
	entity.Name = data.Name
	entity.TargetPrivateLinkResourceID = data.TargetPrivateLinkResourceID
	entity.TargetSubresourceType = data.TargetSubresourceType

	return entity
}

// Filter implements concreteOperations.
func (o *operationsWorkspaceManagedPrivateEndpoint) Filter(entities []fabcore.ManagedPrivateEndpoint, _ string) []fabcore.ManagedPrivateEndpoint {
	ret := make([]fabcore.ManagedPrivateEndpoint, 0, len(entities))

	ret = append(ret, entities...)

	return ret
}

// TransformCreate implements concreteOperations.
func (o *operationsWorkspaceManagedPrivateEndpoint) TransformCreate(entity fabcore.ManagedPrivateEndpoint) fabcore.ManagedPrivateEndpointsClientCreateWorkspaceManagedPrivateEndpointResponse {
	return fabcore.ManagedPrivateEndpointsClientCreateWorkspaceManagedPrivateEndpointResponse{
		ManagedPrivateEndpoint: entity,
	}
}

// TransformGet implements concreteOperations.
func (o *operationsWorkspaceManagedPrivateEndpoint) TransformGet(entity fabcore.ManagedPrivateEndpoint) fabcore.ManagedPrivateEndpointsClientGetWorkspaceManagedPrivateEndpointResponse {
	return fabcore.ManagedPrivateEndpointsClientGetWorkspaceManagedPrivateEndpointResponse{
		ManagedPrivateEndpoint: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsWorkspaceManagedPrivateEndpoint) TransformList(entities []fabcore.ManagedPrivateEndpoint) fabcore.ManagedPrivateEndpointsClientListWorkspaceManagedPrivateEndpointsResponse {
	return fabcore.ManagedPrivateEndpointsClientListWorkspaceManagedPrivateEndpointsResponse{
		ManagedPrivateEndpoints: fabcore.ManagedPrivateEndpoints{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsWorkspaceManagedPrivateEndpoint) TransformUpdate(entity fabcore.ManagedPrivateEndpoint) fabcore.ManagedPrivateEndpoint {
	return entity
}

// Update implements concreteOperations.
func (o *operationsWorkspaceManagedPrivateEndpoint) Update(_, data fabcore.ManagedPrivateEndpoint) fabcore.ManagedPrivateEndpoint {
	return data
}

// Validate implements concreteOperations.
func (o *operationsWorkspaceManagedPrivateEndpoint) Validate(newEntity fabcore.ManagedPrivateEndpoint, existing []fabcore.ManagedPrivateEndpoint) (int, error) {
	for _, entity := range existing {
		if *entity.Name == *newEntity.Name {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureWorkspaceManagedPrivateEndpoint(server *fakeServer) fabcore.ManagedPrivateEndpoint {
	type concreteEntityOperations interface {
		parentIDOperations[
			fabcore.ManagedPrivateEndpoint,
			fabcore.ManagedPrivateEndpointsClientGetWorkspaceManagedPrivateEndpointResponse,
			fabcore.ManagedPrivateEndpoint,
			fabcore.ManagedPrivateEndpointsClientCreateWorkspaceManagedPrivateEndpointResponse,
			fabcore.ManagedPrivateEndpointsClientListWorkspaceManagedPrivateEndpointsResponse,
			fabcore.CreateManagedPrivateEndpointRequest,
			fabcore.ManagedPrivateEndpoint,
		]
	}

	var entityOperations concreteEntityOperations = &operationsWorkspaceManagedPrivateEndpoint{}

	handler := newTypedHandler(server, entityOperations)

	configureEntityWithParentIDNoLRONoUpdate(
		handler,
		entityOperations,
		&handler.ServerFactory.Core.ManagedPrivateEndpointsServer.GetWorkspaceManagedPrivateEndpoint,
		&handler.ServerFactory.Core.ManagedPrivateEndpointsServer.CreateWorkspaceManagedPrivateEndpoint,
		&handler.ServerFactory.Core.ManagedPrivateEndpointsServer.NewListWorkspaceManagedPrivateEndpointsPager,
		&handler.ServerFactory.Core.ManagedPrivateEndpointsServer.DeleteWorkspaceManagedPrivateEndpoint,
	)

	return fabcore.ManagedPrivateEndpoint{}
}

func NewRandomWorkspaceManagedPrivateEndpoint() fabcore.ManagedPrivateEndpoint {
	return fabcore.ManagedPrivateEndpoint{
		ID:                          to.Ptr(testhelp.RandomUUID()),
		Name:                        to.Ptr(testhelp.RandomName()),
		ProvisioningState:           to.Ptr(fabcore.PrivateEndpointProvisioningStateSucceeded),
		TargetPrivateLinkResourceID: to.Ptr(testhelp.RandomUUID()),
		TargetSubresourceType:       to.Ptr(testhelp.RandomName()),
		ConnectionState: &fabcore.PrivateEndpointConnectionState{
			ActionsRequired: to.Ptr("None"),
			Description:     to.Ptr(testhelp.RandomName()),
			Status:          to.Ptr(fabcore.ConnectionStatusApproved),
		},
	}
}
