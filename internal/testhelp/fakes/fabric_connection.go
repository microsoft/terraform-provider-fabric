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

// operationsConnection implements SimpleIDOperations.
type operationsConnection struct{}

// GetID implements concreteOperations.
func (o *operationsConnection) GetID(entity fabcore.Connection) string {
	return *entity.ID
}

// TransformCreate implements concreteOperations.
func (o *operationsConnection) TransformCreate(entity fabcore.Connection) fabcore.ConnectionsClientCreateConnectionResponse {
	return fabcore.ConnectionsClientCreateConnectionResponse{
		Connection: entity,
	}
}

// Create implements concreteOperations.
func (o *operationsConnection) Create(data fabcore.CreateConnectionRequestClassification) fabcore.Connection {
	d := data.GetCreateConnectionRequest()

	entity := NewRandomConnection()
	entity.DisplayName = d.DisplayName
	entity.PrivacyLevel = d.PrivacyLevel
	entity.ConnectivityType = d.ConnectivityType
	entity.ConnectionDetails.Type = d.ConnectionDetails.Type

	return entity
}

// TransformGet implements concreteOperations.
func (o *operationsConnection) TransformGet(entity fabcore.Connection) fabcore.ConnectionsClientGetConnectionResponse {
	return fabcore.ConnectionsClientGetConnectionResponse{
		Connection: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsConnection) TransformList(entities []fabcore.Connection) fabcore.ConnectionsClientListConnectionsResponse {
	return fabcore.ConnectionsClientListConnectionsResponse{
		ListConnectionsResponse: fabcore.ListConnectionsResponse{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsConnection) TransformUpdate(entity fabcore.Connection) fabcore.ConnectionsClientUpdateConnectionResponse {
	return fabcore.ConnectionsClientUpdateConnectionResponse{
		Connection: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsConnection) Update(base fabcore.Connection, data fabcore.UpdateConnectionRequestClassification) fabcore.Connection {
	d := data.GetUpdateConnectionRequest()

	base.ConnectivityType = d.ConnectivityType
	base.PrivacyLevel = d.PrivacyLevel

	return base
}

// Validate implements concreteOperations.
func (o *operationsConnection) Validate(newEntity fabcore.Connection, existing []fabcore.Connection) (int, error) {
	for _, entity := range existing {
		if *entity.DisplayName == *newEntity.DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureConnection(server *fakeServer) fabcore.Connection {
	type concreteEntityOperations interface {
		simpleIDOperations[
			fabcore.Connection,
			fabcore.ConnectionsClientGetConnectionResponse,
			fabcore.ConnectionsClientUpdateConnectionResponse,
			fabcore.ConnectionsClientCreateConnectionResponse,
			fabcore.ConnectionsClientListConnectionsResponse,
			fabcore.CreateConnectionRequestClassification,
			fabcore.UpdateConnectionRequestClassification,
		]
	}

	var entityOperations concreteEntityOperations = &operationsConnection{}

	handler := newTypedHandler(server, entityOperations)

	configureEntityPagerWithSimpleID(
		handler,
		entityOperations,
		&handler.ServerFactory.Core.ConnectionsServer.GetConnection,
		&handler.ServerFactory.Core.ConnectionsServer.UpdateConnection,
		&handler.ServerFactory.Core.ConnectionsServer.CreateConnection,
		&handler.ServerFactory.Core.ConnectionsServer.NewListConnectionsPager,
		&handler.ServerFactory.Core.ConnectionsServer.DeleteConnection,
	)

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
