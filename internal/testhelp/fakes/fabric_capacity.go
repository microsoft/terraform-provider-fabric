// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// operationsCapacity implements SimpleIDOperations.
type operationsCapacity struct{}

// TransformList implements concreteOperations.
func (o *operationsCapacity) TransformList(entities []fabcore.Capacity) fabcore.CapacitiesClientListCapacitiesResponse {
	return fabcore.CapacitiesClientListCapacitiesResponse{
		Capacities: fabcore.Capacities{
			Value: entities,
		},
	}
}

func (o *operationsCapacity) GetID(entity fabcore.Capacity) string {
	return *entity.ID
}

func configureCapacity(server *fakeServer) fabcore.Capacity {
	type concreteEntityOperations interface {
		identifier[fabcore.Capacity]
		listTransformer[fabcore.Capacity, fabcore.CapacitiesClientListCapacitiesResponse]
	}

	var entityOperations concreteEntityOperations = &operationsCapacity{}

	handler := newTypedHandler(server, entityOperations)

	handleListPager(
		handler,
		entityOperations,
		&handler.ServerFactory.Core.CapacitiesServer.NewListCapacitiesPager)

	return fabcore.Capacity{}
}

func NewRandomCapacity() fabcore.Capacity {
	return fabcore.Capacity{
		ID:          to.Ptr(testhelp.RandomUUID()),
		DisplayName: to.Ptr(testhelp.RandomName()),
		Region:      to.Ptr(testhelp.RandomName()),
		SKU:         to.Ptr(testhelp.RandomName()),
		State:       to.Ptr(fabcore.CapacityStateActive),
	}
}
