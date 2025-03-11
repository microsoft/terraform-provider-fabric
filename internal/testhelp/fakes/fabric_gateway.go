// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/services/gateway"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// operationsGateway implements SimpleIDOperations.
type operationsGateway struct{}

func (o *operationsGateway) Create(data fabcore.CreateGatewayRequestClassification) fabcore.GatewayClassification {
	switch data := data.(type) {
	case *fabcore.CreateVirtualNetworkGatewayRequest:
		entity := NewRandomVirtualNetworkGateway()
		entity.Type = data.Type
		entity.DisplayName = data.DisplayName
		entity.CapacityID = data.CapacityID
		entity.InactivityMinutesBeforeSleep = data.InactivityMinutesBeforeSleep
		entity.NumberOfMemberGateways = data.NumberOfMemberGateways
		entity.VirtualNetworkAzureResource = data.VirtualNetworkAzureResource

		return entity
	default:
		panic("Unsupported Gateway type")
	}
}

func (o *operationsGateway) TransformCreate(entity fabcore.GatewayClassification) fabcore.GatewaysClientCreateGatewayResponse {
	return fabcore.GatewaysClientCreateGatewayResponse{
		GatewayClassification: entity,
	}
}

func (o *operationsGateway) TransformGet(entity fabcore.GatewayClassification) fabcore.GatewaysClientGetGatewayResponse {
	return fabcore.GatewaysClientGetGatewayResponse{
		GatewayClassification: entity,
	}
}

func (o *operationsGateway) TransformList(list []fabcore.GatewayClassification) fabcore.GatewaysClientListGatewaysResponse {
	return fabcore.GatewaysClientListGatewaysResponse{
		ListGatewaysResponse: fabcore.ListGatewaysResponse{
			Value: list,
		},
	}
}

func (o *operationsGateway) TransformUpdate(entity fabcore.GatewayClassification) fabcore.GatewaysClientUpdateGatewayResponse {
	return fabcore.GatewaysClientUpdateGatewayResponse{
		GatewayClassification: entity,
	}
}

func (o *operationsGateway) Update(base fabcore.GatewayClassification, data fabcore.UpdateGatewayRequestClassification) fabcore.GatewayClassification {
	switch base := base.(type) {
	case *fabcore.VirtualNetworkGateway:
		request, ok := data.(*fabcore.UpdateVirtualNetworkGatewayRequest)
		if !ok {
			panic("Invalid update data for VirtualNetworkGateway")
		}

		base.CapacityID = request.CapacityID
		base.DisplayName = request.DisplayName
		base.InactivityMinutesBeforeSleep = request.InactivityMinutesBeforeSleep
		base.NumberOfMemberGateways = request.NumberOfMemberGateways

		return base
	default:
		panic("Unsupported Gateway type")
	}
}

func (o *operationsGateway) Validate(newEntity fabcore.GatewayClassification, existing []fabcore.GatewayClassification) (int, error) {
	for _, existingGateway := range existing {
		if *(existingGateway.GetGateway().Type) != *(newEntity.GetGateway().Type) {
			continue
		}

		if newVNG, ok := newEntity.(*fabcore.VirtualNetworkGateway); ok {
			if existingVNG, ok := existingGateway.(*fabcore.VirtualNetworkGateway); ok && *existingVNG.DisplayName == *newVNG.DisplayName {
				return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrGateway.DuplicateGatewayName.Error(), fabcore.ErrGateway.DuplicateGatewayName.Error())
			}
		}
	}

	return http.StatusCreated, nil
}

func (o *operationsGateway) GetID(entity fabcore.GatewayClassification) string {
	return *entity.GetGateway().ID
}

func configureVirtualNetworkGateway(server *fakeServer) fabcore.VirtualNetworkGateway {
	configureGateway(server)

	return fabcore.VirtualNetworkGateway{}
}

func configureOnPremisesGateway(server *fakeServer) fabcore.OnPremisesGateway {
	configureGateway(server)

	return fabcore.OnPremisesGateway{}
}

func configureOnPremisesGatewayPersonal(server *fakeServer) fabcore.OnPremisesGatewayPersonal {
	configureGateway(server)

	return fabcore.OnPremisesGatewayPersonal{}
}

func configureGateway(server *fakeServer) {
	type concreteEntityOperations interface {
		simpleIDOperations[
			fabcore.GatewayClassification,
			fabcore.GatewaysClientGetGatewayResponse,
			fabcore.GatewaysClientUpdateGatewayResponse,
			fabcore.GatewaysClientCreateGatewayResponse,
			fabcore.GatewaysClientListGatewaysResponse,
			fabcore.CreateGatewayRequestClassification,
			fabcore.UpdateGatewayRequestClassification]
	}

	var entityOperations concreteEntityOperations = &operationsGateway{}

	handler := newTypedHandler(server, entityOperations)

	configureEntityPagerWithSimpleID(
		handler,
		entityOperations,
		&handler.ServerFactory.Core.GatewaysServer.GetGateway,
		&handler.ServerFactory.Core.GatewaysServer.UpdateGateway,
		&handler.ServerFactory.Core.GatewaysServer.CreateGateway,
		&handler.ServerFactory.Core.GatewaysServer.NewListGatewaysPager,
		&handler.ServerFactory.Core.GatewaysServer.DeleteGateway)
}

func NewRandomGateway() fabcore.GatewayClassification {
	gatewayType := testhelp.RandomElement(fabcore.PossibleGatewayTypeValues())

	switch gatewayType {
	case fabcore.GatewayTypeOnPremises:
		return NewRandomOnPremisesGateway()
	case fabcore.GatewayTypeOnPremisesPersonal:
		return NewRandomOnPremisesGatewayPersonal()
	case fabcore.GatewayTypeVirtualNetwork:
		return NewRandomVirtualNetworkGateway()
	default:
		panic("Unsupported Gateway type")
	}
}

func NewRandomOnPremisesGateway() *fabcore.OnPremisesGateway {
	return &fabcore.OnPremisesGateway{
		ID:                          to.Ptr(testhelp.RandomUUID()),
		Type:                        to.Ptr(fabcore.GatewayTypeOnPremises),
		DisplayName:                 to.Ptr(testhelp.RandomName()),
		AllowCloudConnectionRefresh: to.Ptr(testhelp.RandomBool()),
		AllowCustomConnectors:       to.Ptr(testhelp.RandomBool()),
		NumberOfMemberGateways:      to.Ptr(testhelp.RandomIntRange(gateway.MinNumberOfMemberGatewaysValues, gateway.MaxNumberOfMemberGatewaysValues)),
		LoadBalancingSetting:        to.Ptr(testhelp.RandomElement(fabcore.PossibleLoadBalancingSettingValues())),
		Version:                     to.Ptr(testhelp.RandomName()),
		PublicKey:                   NewRadomPublicKey(),
	}
}

func NewRandomOnPremisesGatewayPersonal() *fabcore.OnPremisesGatewayPersonal {
	return &fabcore.OnPremisesGatewayPersonal{
		ID:        to.Ptr(testhelp.RandomUUID()),
		Type:      to.Ptr(fabcore.GatewayTypeOnPremisesPersonal),
		Version:   to.Ptr(testhelp.RandomName()),
		PublicKey: NewRadomPublicKey(),
	}
}

func NewRandomVirtualNetworkGateway() *fabcore.VirtualNetworkGateway {
	return &fabcore.VirtualNetworkGateway{
		ID:                           to.Ptr(testhelp.RandomUUID()),
		Type:                         to.Ptr(fabcore.GatewayTypeVirtualNetwork),
		DisplayName:                  to.Ptr(testhelp.RandomName()),
		CapacityID:                   to.Ptr(testhelp.RandomUUID()),
		InactivityMinutesBeforeSleep: to.Ptr(testhelp.RandomElement(gateway.PossibleInactivityMinutesBeforeSleepValues)),
		NumberOfMemberGateways:       to.Ptr(testhelp.RandomIntRange(gateway.MinNumberOfMemberGatewaysValues, gateway.MaxNumberOfMemberGatewaysValues)),
		VirtualNetworkAzureResource:  NewRandomVirtualNetworkAzureResource(),
	}
}

func NewRadomPublicKey() *fabcore.PublicKey {
	return &fabcore.PublicKey{
		Exponent: to.Ptr(testhelp.RandomName()),
		Modulus:  to.Ptr(testhelp.RandomName()),
	}
}

func NewRandomVirtualNetworkAzureResource() *fabcore.VirtualNetworkAzureResource {
	return &fabcore.VirtualNetworkAzureResource{
		ResourceGroupName:  to.Ptr(testhelp.RandomName()),
		SubnetName:         to.Ptr(testhelp.RandomName()),
		SubscriptionID:     to.Ptr(testhelp.RandomUUID()),
		VirtualNetworkName: to.Ptr(testhelp.RandomName()),
	}
}
