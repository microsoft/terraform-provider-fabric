// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

// operationsConnection implements SimpleIDOperations.
type operationsConnection struct{}

// TransformList implements concreteOperations.
func (o *operationsConnection) TransformList(entities []fabcore.Connection) fabcore.ConnectionsClientListConnectionsResponse {
	return fabcore.ConnectionsClientListConnectionsResponse{
		ListConnectionsResponse: fabcore.ListConnectionsResponse{
			Value: entities,
		},
	}
}

func (o *operationsConnection) GetID(entity fabcore.Connection) string {
	return *entity.ID
}

func configureConnection(server *fakeServer) fabcore.Connection {
	type concreteEntityOperations interface {
		identifier[fabcore.Connection]
		listTransformer[fabcore.Connection, fabcore.ConnectionsClientListConnectionsResponse]
	}

	var entityOperations concreteEntityOperations = &operationsConnection{}

	handler := newTypedHandler(server, entityOperations)

	handleListPager(
		handler,
		entityOperations,
		&handler.ServerFactory.Core.ConnectionsServer.NewListConnectionsPager)

	return fabcore.Connection{}
}

func NewRandomConnection() fabcore.Connection {
	return fabcore.Connection{
		ID:               to.Ptr(testhelp.RandomUUID()),
		DisplayName:      to.Ptr(testhelp.RandomName()),
		GatewayID:        to.Ptr(testhelp.RandomUUID()),
		PrivacyLevel:     to.Ptr(fabcore.PrivacyLevelPrivate),
		ConnectivityType: to.Ptr(fabcore.ConnectivityTypeShareableCloud),
		ConnectionDetails: &fabcore.ListConnectionDetails{
			Path: to.Ptr(testhelp.RandomURI()),
			Type: to.Ptr("GitHubSourceControl"),
		},
		CredentialDetails: &fabcore.ListCredentialDetails{
			CredentialType:       to.Ptr(fabcore.CredentialTypeKey),
			SingleSignOnType:     to.Ptr(fabcore.SingleSignOnTypeNone),
			ConnectionEncryption: to.Ptr(fabcore.ConnectionEncryptionNotEncrypted),
			SkipTestConnection:   to.Ptr(false),
		},
	}
}
