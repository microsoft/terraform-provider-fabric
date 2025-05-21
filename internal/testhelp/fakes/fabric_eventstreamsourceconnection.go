// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

// import (
// 	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
// 	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
// 	fabeventstream "github.com/microsoft/fabric-sdk-go/fabric/eventstream"

// 	"github.com/microsoft/terraform-provider-fabric/internal/services/gateway"
// 	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
// )

// // operationsEventstreamSourceConnection implements simpleEphemeralIDOperations.
// type operationsEventstreamSourceConnection struct{}

// func (o *operationsEventstreamSourceConnection) TransformOpen(entity fabeventstream.SourceConnectionResponse) fabeventstream.TopologyClientGetEventstreamSourceConnectionResponse {
// 	return fabeventstream.TopologyClientGetEventstreamSourceConnectionResponse{
// 		SourceConnectionResponse: entity,
// 	}
// }

// func (o *operationsEventstreamSourceConnection) GetID(entity fabeventstream.SourceConnectionResponse) string {
// 	return *entity.FullyQualifiedNamespace
// }

// func configureEventstreamSourceConnection(server *fakeServer) {
// 	type concreteEntityOperations interface {
// 		simpleEphemeralIDOperations[
// 			fabeventstream.SourceConnectionResponse,
// 			fabeventstream.TopologyClientGetEventstreamSourceConnectionResponse,
// 		]
// 	}

// 	var entityOperations concreteEntityOperations = &operationsEventstreamSourceConnection{}

// 	handler := newTypedHandler(server, entityOperations)
// 	// func(ctx context.Context, workspaceID string, eventstreamID string, sourceID string, options *eventstream.TopologyClientGetEventstreamSourceConnectionOptions) (resp azfake.Responder[eventstream.TopologyClientGetEventstreamSourceConnectionResponse], errResp azfake.ErrorResponder)

// 	configureEphemeralEntityWithSimpleID(
// 		handler,
// 		entityOperations,
// 		&handler.ServerFactory.Eventstream.TopologyServer.GetEventstreamSourceConnection)
// }

// func NewRandomGateway() fabcore.GatewayClassification {
// 	gatewayType := testhelp.RandomElement(fabcore.PossibleGatewayTypeValues())

// 	switch gatewayType {
// 	case fabcore.GatewayTypeOnPremises:
// 		return NewRandomOnPremisesGateway()
// 	case fabcore.GatewayTypeOnPremisesPersonal:
// 		return NewRandomOnPremisesGatewayPersonal()
// 	case fabcore.GatewayTypeVirtualNetwork:
// 		return NewRandomVirtualNetworkGateway()
// 	default:
// 		panic("Unsupported Gateway type") // lintignore:R009
// 	}
// }

// func NewRandomOnPremisesGateway() *fabcore.OnPremisesGateway {
// 	return &fabcore.OnPremisesGateway{
// 		ID:                          to.Ptr(testhelp.RandomUUID()),
// 		Type:                        to.Ptr(fabcore.GatewayTypeOnPremises),
// 		DisplayName:                 to.Ptr(testhelp.RandomName()),
// 		AllowCloudConnectionRefresh: to.Ptr(testhelp.RandomBool()),
// 		AllowCustomConnectors:       to.Ptr(testhelp.RandomBool()),
// 		NumberOfMemberGateways:      to.Ptr(testhelp.RandomIntRange(gateway.MinNumberOfMemberGatewaysValues, gateway.MaxNumberOfMemberGatewaysValues)),
// 		LoadBalancingSetting:        to.Ptr(testhelp.RandomElement(fabcore.PossibleLoadBalancingSettingValues())),
// 		Version:                     to.Ptr(testhelp.RandomName()),
// 		PublicKey:                   NewRadomPublicKey(),
// 	}
// }

// func NewRandomOnPremisesGatewayPersonal() *fabcore.OnPremisesGatewayPersonal {
// 	return &fabcore.OnPremisesGatewayPersonal{
// 		ID:        to.Ptr(testhelp.RandomUUID()),
// 		Type:      to.Ptr(fabcore.GatewayTypeOnPremisesPersonal),
// 		Version:   to.Ptr(testhelp.RandomName()),
// 		PublicKey: NewRadomPublicKey(),
// 	}
// }

// func NewRandomVirtualNetworkGateway() *fabcore.VirtualNetworkGateway {
// 	return &fabcore.VirtualNetworkGateway{
// 		ID:                           to.Ptr(testhelp.RandomUUID()),
// 		Type:                         to.Ptr(fabcore.GatewayTypeVirtualNetwork),
// 		DisplayName:                  to.Ptr(testhelp.RandomName()),
// 		CapacityID:                   to.Ptr(testhelp.RandomUUID()),
// 		InactivityMinutesBeforeSleep: to.Ptr(testhelp.RandomElement(gateway.PossibleInactivityMinutesBeforeSleepValues)),
// 		NumberOfMemberGateways:       to.Ptr(testhelp.RandomIntRange(gateway.MinNumberOfMemberGatewaysValues, gateway.MaxNumberOfMemberGatewaysValues)),
// 		VirtualNetworkAzureResource:  NewRandomVirtualNetworkAzureResource(),
// 	}
// }

// func NewRadomPublicKey() *fabcore.PublicKey {
// 	return &fabcore.PublicKey{
// 		Exponent: to.Ptr(testhelp.RandomName()),
// 		Modulus:  to.Ptr(testhelp.RandomName()),
// 	}
// }

// func NewRandomVirtualNetworkAzureResource() *fabcore.VirtualNetworkAzureResource {
// 	return &fabcore.VirtualNetworkAzureResource{
// 		ResourceGroupName:  to.Ptr(testhelp.RandomName()),
// 		SubnetName:         to.Ptr(testhelp.RandomName()),
// 		SubscriptionID:     to.Ptr(testhelp.RandomUUID()),
// 		VirtualNetworkName: to.Ptr(testhelp.RandomName()),
// 	}
// }
