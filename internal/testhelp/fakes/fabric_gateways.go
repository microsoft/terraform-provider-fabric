// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// operationsGateway implements SimpleIDOperations.
type operationsGateway struct{}

// Create implements concreteEntityOperations.
func (o *operationsGateway) Create(data fabcore.CreateGatewayRequestClassification) fabcore.GatewayClassification {
	switch gateway := data.(type) {
	case *fabcore.CreateVirtualNetworkGatewayRequest:
		returnGateway := NewRandomVirtualNetworkGateway()
		returnGateway.DisplayName = gateway.DisplayName
		returnGateway.CapacityID = gateway.CapacityID
		returnGateway.InactivityMinutesBeforeSleep = gateway.InactivityMinutesBeforeSleep
		returnGateway.NumberOfMemberGateways = gateway.NumberOfMemberGateways
		returnGateway.VirtualNetworkAzureResource = gateway.VirtualNetworkAzureResource
		return &returnGateway
	default:
		panic("unimplemented")
	}
}

// GetID implements concreteEntityOperations.
func (o *operationsGateway) GetID(entity fabcore.GatewayClassification) string {
	return *entity.GetGateway().ID
}

// TransformGet implements concreteEntityOperations.
func (o *operationsGateway) TransformGet(entity fabcore.GatewayClassification) fabcore.GatewaysClientGetGatewayResponse {
	return fabcore.GatewaysClientGetGatewayResponse{
		GatewayClassification: entity,
	}
}

// TransformList implements concreteEntityOperations.
func (o *operationsGateway) TransformList(list []fabcore.GatewayClassification) fabcore.GatewaysClientListGatewaysResponse {
	return fabcore.GatewaysClientListGatewaysResponse{
		ListGatewaysResponse: fabcore.ListGatewaysResponse{
			Value: list,
		},
	}
}

// TransformUpdate implements concreteEntityOperations.
func (o *operationsGateway) TransformUpdate(entity fabcore.GatewayClassification) fabcore.GatewaysClientUpdateGatewayResponse {
	return fabcore.GatewaysClientUpdateGatewayResponse{
		GatewayClassification: entity,
	}
}

// Update implements concreteEntityOperations.
func (o *operationsGateway) Update(base fabcore.GatewayClassification, data fabcore.UpdateGatewayRequestClassification) fabcore.GatewayClassification {
	switch request := data.(type) {
	case *fabcore.UpdateVirtualNetworkGatewayRequest:
		gateway, _ := base.(*fabcore.VirtualNetworkGateway)
		gateway.CapacityID = request.CapacityID
		gateway.DisplayName = request.DisplayName
		gateway.InactivityMinutesBeforeSleep = request.InactivityMinutesBeforeSleep
		gateway.NumberOfMemberGateways = request.NumberOfMemberGateways
		return gateway
	case *fabcore.UpdateOnPremisesGatewayRequest:
		gateway, _ := base.(*fabcore.OnPremisesGateway)
		gateway.AllowCloudConnectionRefresh = request.AllowCloudConnectionRefresh
		gateway.AllowCustomConnectors = request.AllowCustomConnectors
		gateway.LoadBalancingSetting = request.LoadBalancingSetting
		gateway.DisplayName = request.DisplayName
		return gateway
	default:
		panic("unimplemented")
	}
}

// Validate implements concreteEntityOperations.
func (o *operationsGateway) Validate(newEntity fabcore.GatewayClassification, existing []fabcore.GatewayClassification) (statusCode int, err error) {
	for _, existingGateway := range existing {
		if existingGateway.GetGateway().Type != newEntity.GetGateway().Type {
			continue
		}
		switch gateway := newEntity.(type) {
		case *fabcore.VirtualNetworkGateway:
			vng := existingGateway.(*fabcore.VirtualNetworkGateway)
			if *vng.DisplayName == *gateway.DisplayName {
				return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrWorkspace.WorkspaceNameAlreadyExists.Error(), fabcore.ErrWorkspace.WorkspaceNameAlreadyExists.Error())
			}
		case *fabcore.OnPremisesGateway:
			opg := existingGateway.(*fabcore.OnPremisesGateway)
			if *opg.DisplayName == *gateway.DisplayName {
				return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrWorkspace.WorkspaceNameAlreadyExists.Error(), fabcore.ErrWorkspace.WorkspaceNameAlreadyExists.Error())
			}
		}
	}

	return http.StatusCreated, nil
}

// TransformCreate implements concreteEntityOperations.
func (o *operationsGateway) TransformCreate(entity fabcore.GatewayClassification) fabcore.GatewaysClientCreateGatewayResponse {
	return fabcore.GatewaysClientCreateGatewayResponse{
		GatewayClassification: entity,
	}
}

func configureGateway(server *fakeServer) fabcore.VirtualNetworkGateway {
	type concreteEntityOperations interface {
		simpleIDOperations[
			fabcore.GatewayClassification,
			fabcore.GatewaysClientGetGatewayResponse,
			fabcore.GatewaysClientUpdateGatewayResponse,
			fabcore.GatewaysClientCreateGatewayResponse,
			fabcore.GatewaysClientListGatewaysResponse,
			fabcore.CreateGatewayRequestClassification,
			fabcore.UpdateGatewayRequestClassification,
		]
	}

	var entityOperations concreteEntityOperations = &operationsGateway{}

	handler := newTypedHandler(server, entityOperations)

	handleGetWithSimpleID(handler, entityOperations, &handler.ServerFactory.Core.GatewaysServer.GetGateway)
	handleUpdateWithSimpleID(handler, entityOperations, entityOperations, &handler.ServerFactory.Core.GatewaysServer.UpdateGateway)
	handleCreate(handler, entityOperations, entityOperations, entityOperations, &handler.ServerFactory.Core.GatewaysServer.CreateGateway)
	handleDeleteWithSimpleID(handler, &handler.ServerFactory.Core.GatewaysServer.DeleteGateway)

	handleListPager(
		handler,
		entityOperations,
		&handler.ServerFactory.Core.GatewaysServer.NewListGatewaysPager)

	return fabcore.VirtualNetworkGateway{}
}

func NewRandomVirtualNetworkGateway() fabcore.VirtualNetworkGateway {
	return fabcore.VirtualNetworkGateway{
		ID:                           to.Ptr(testhelp.RandomUUID()),
		DisplayName:                  to.Ptr(testhelp.RandomName()),
		InactivityMinutesBeforeSleep: to.Ptr(testhelp.RandomInt32(100)),
		NumberOfMemberGateways:       to.Ptr(testhelp.RandomInt32(100)),
		CapacityID:                   to.Ptr(testhelp.RandomUUID()),
		Type:                         to.Ptr(fabcore.GatewayTypeVirtualNetwork),
		VirtualNetworkAzureResource: &fabcore.VirtualNetworkAzureResource{
			SubscriptionID:     to.Ptr(testhelp.RandomUUID()),
			ResourceGroupName:  to.Ptr(testhelp.RandomName()),
			VirtualNetworkName: to.Ptr(testhelp.RandomName()),
			SubnetName:         to.Ptr(testhelp.RandomName()),
		},
	}
}

func NewRandomOnPremisesGateway() fabcore.OnPremisesGateway {
	return fabcore.OnPremisesGateway{
		ID:                          to.Ptr(testhelp.RandomUUID()),
		DisplayName:                 to.Ptr(testhelp.RandomName()),
		NumberOfMemberGateways:      to.Ptr(testhelp.RandomInt32(100)),
		Type:                        to.Ptr(fabcore.GatewayTypeOnPremises),
		AllowCloudConnectionRefresh: to.Ptr(true),
		AllowCustomConnectors:       to.Ptr(false),
		LoadBalancingSetting:        to.Ptr(fabcore.LoadBalancingSettingDistributeEvenly),
		PublicKey: &fabcore.PublicKey{
			Exponent: to.Ptr(testhelp.RandomName()),
			Modulus:  to.Ptr(testhelp.RandomName()),
		},
		Version: to.Ptr("1.0"),
	}
}

func NewRandomOnPermisesGatewayPersonal() fabcore.OnPremisesGatewayPersonal {
	return fabcore.OnPremisesGatewayPersonal{
		ID:      to.Ptr(testhelp.RandomUUID()),
		Type:    to.Ptr(fabcore.GatewayTypeOnPremisesPersonal),
		Version: to.Ptr("1.0"),
		PublicKey: &fabcore.PublicKey{
			Exponent: to.Ptr(testhelp.RandomName()),
			Modulus:  to.Ptr(testhelp.RandomName()),
		},
	}
}
