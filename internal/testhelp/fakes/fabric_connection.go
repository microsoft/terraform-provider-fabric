// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fakes

import (
	"net/http"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
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

	// Handle CredentialDetails for specific connection request types
	switch req := data.(type) {
	case *fabcore.CreateCloudConnectionRequest:
		if req.CredentialDetails != nil {
			entity.CredentialDetails = &fabcore.ListCredentialDetails{
				CredentialType:       req.CredentialDetails.Credentials.GetCredentials().CredentialType,
				ConnectionEncryption: req.CredentialDetails.ConnectionEncryption,
				SingleSignOnType:     req.CredentialDetails.SingleSignOnType,
				SkipTestConnection:   req.CredentialDetails.SkipTestConnection,
				// AllowConnectionUsageInGateway: req.CredentialDetails.AllowConnectionUsageInGateway,
			}
		}
	case *fabcore.CreateVirtualNetworkGatewayConnectionRequest:
		entity.GatewayID = req.GatewayID

		if req.CredentialDetails != nil {
			entity.CredentialDetails = &fabcore.ListCredentialDetails{
				CredentialType:       req.CredentialDetails.Credentials.GetCredentials().CredentialType,
				ConnectionEncryption: req.CredentialDetails.ConnectionEncryption,
				SingleSignOnType:     req.CredentialDetails.SingleSignOnType,
				SkipTestConnection:   req.CredentialDetails.SkipTestConnection,
			}
		}
	}

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

	// Handle specific update request types
	switch req := data.(type) {
	case *fabcore.UpdateShareableCloudConnectionRequest:
		base.DisplayName = req.DisplayName
		base.PrivacyLevel = req.PrivacyLevel

		if req.CredentialDetails != nil {
			base.CredentialDetails = &fabcore.ListCredentialDetails{
				CredentialType:       req.CredentialDetails.Credentials.GetCredentials().CredentialType,
				ConnectionEncryption: req.CredentialDetails.ConnectionEncryption,
				SingleSignOnType:     req.CredentialDetails.SingleSignOnType,
				SkipTestConnection:   req.CredentialDetails.SkipTestConnection,
				// AllowConnectionUsageInGateway: req.CredentialDetails.AllowConnectionUsageInGateway,
			}
		}
	case *fabcore.UpdateVirtualNetworkGatewayConnectionRequest:
		base.DisplayName = req.DisplayName
		base.PrivacyLevel = req.PrivacyLevel

		if req.CredentialDetails != nil {
			base.CredentialDetails = &fabcore.ListCredentialDetails{
				CredentialType:       req.CredentialDetails.Credentials.GetCredentials().CredentialType,
				ConnectionEncryption: req.CredentialDetails.ConnectionEncryption,
				SingleSignOnType:     req.CredentialDetails.SingleSignOnType,
				SkipTestConnection:   req.CredentialDetails.SkipTestConnection,
			}
		}
	}

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

	// Configure the NewListSupportedConnectionTypesPager handler
	handler.ServerFactory.Core.ConnectionsServer.NewListSupportedConnectionTypesPager = FakeListSupportedConnectionTypes()

	return fabcore.Connection{}
}

// FakeListSupportedConnectionTypes returns a fake handler for listing supported connection types.
func FakeListSupportedConnectionTypes() func(options *fabcore.ConnectionsClientListSupportedConnectionTypesOptions) azfake.PagerResponder[fabcore.ConnectionsClientListSupportedConnectionTypesResponse] {
	return func(options *fabcore.ConnectionsClientListSupportedConnectionTypesOptions) azfake.PagerResponder[fabcore.ConnectionsClientListSupportedConnectionTypesResponse] {
		var resp azfake.PagerResponder[fabcore.ConnectionsClientListSupportedConnectionTypesResponse]

		// Create the fake response data with comprehensive test parameters
		response := fabcore.ConnectionsClientListSupportedConnectionTypesResponse{
			ListSupportedConnectionTypesResponse: fabcore.ListSupportedConnectionTypesResponse{
				Value: []fabcore.ConnectionCreationMetadata{
					{
						Type: to.Ptr("FTP"),
						CreationMethods: []fabcore.ConnectionCreationMethod{
							{
								Name: to.Ptr("FTP.Contents"),
								Parameters: []fabcore.ConnectionCreationParameter{
									{
										Name:     to.Ptr("server"),
										DataType: to.Ptr(fabcore.DataTypeText),
										Required: to.Ptr(true),
									},
									{
										Name:     to.Ptr("database"),
										DataType: to.Ptr(fabcore.DataTypeText),
										Required: to.Ptr(false),
									},
									{
										Name:     to.Ptr("enable_ssl"),
										DataType: to.Ptr(fabcore.DataTypeBoolean),
										Required: to.Ptr(false),
									},
									{
										Name:     to.Ptr("start_date"),
										DataType: to.Ptr(fabcore.DataTypeDate),
										Required: to.Ptr(false),
									},
									{
										Name:     to.Ptr("created_at"),
										DataType: to.Ptr(fabcore.DataTypeDateTime),
										Required: to.Ptr(false),
									},
									{
										Name:     to.Ptr("backup_time"),
										DataType: to.Ptr(fabcore.DataTypeTime),
										Required: to.Ptr(false),
									},
									{
										Name:     to.Ptr("port"),
										DataType: to.Ptr(fabcore.DataTypeNumber),
										Required: to.Ptr(false),
									},
									{
										Name:     to.Ptr("timeout"),
										DataType: to.Ptr(fabcore.DataTypeDuration),
										Required: to.Ptr(false),
									},
									{
										Name:          to.Ptr("ssl_mode"),
										DataType:      to.Ptr(fabcore.DataTypeText),
										Required:      to.Ptr(false),
										AllowedValues: []string{"required", "optional", "disabled"},
									},
								},
							},
						},
						SupportedCredentialTypes: []fabcore.CredentialType{
							fabcore.CredentialTypeBasic,
							fabcore.CredentialTypeOAuth2,
							fabcore.CredentialTypeKey,
						},
						SupportedConnectionEncryptionTypes: []fabcore.ConnectionEncryption{
							fabcore.ConnectionEncryptionNotEncrypted,
						},
						SupportsSkipTestConnection: to.Ptr(true),
					},
				},
			},
		}

		resp.AddPage(http.StatusOK, response, nil)

		return resp
	}
}

// NewRandomConnection creates a connection with a randomly selected credential type.
func NewRandomConnection() fabcore.Connection {
	credentialTypes := []fabcore.CredentialType{
		fabcore.CredentialTypeBasic,
		fabcore.CredentialTypeKey,
		fabcore.CredentialTypeServicePrincipal,
		fabcore.CredentialTypeSharedAccessSignature,
		fabcore.CredentialTypeAnonymous,
		fabcore.CredentialTypeWorkspaceIdentity,
	}

	randomCredentialType := testhelp.RandomElement(credentialTypes)

	switch randomCredentialType {
	case fabcore.CredentialTypeBasic:
		return NewRandomConnectionWithBasicCredentials()
	case fabcore.CredentialTypeKey:
		return NewRandomConnectionWithKeyCredentials()
	case fabcore.CredentialTypeServicePrincipal:
		return NewRandomConnectionWithServicePrincipalCredentials()
	case fabcore.CredentialTypeSharedAccessSignature:
		return NewRandomConnectionWithSharedAccessSignatureCredentials()
	case fabcore.CredentialTypeAnonymous:
		return NewRandomConnectionWithAnonymousCredentials()
	case fabcore.CredentialTypeWorkspaceIdentity:
		return NewRandomConnectionWithWorkspaceIdentityCredentials()
	default:
		return NewRandomConnectionWithKeyCredentials()
	}
}

// ========== Connectivity Type Functions ==========

// NewRandomShareableCloudConnection creates a connection with ShareableCloud connectivity.
func NewRandomShareableCloudConnection() fabcore.Connection {
	return fabcore.Connection{
		ID:               to.Ptr(testhelp.RandomUUID()),
		DisplayName:      to.Ptr(testhelp.RandomName()),
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

// NewRandomVirtualNetworkGatewayConnection creates a connection with VirtualNetworkGateway connectivity.
func NewRandomVirtualNetworkGatewayConnection() fabcore.Connection {
	return fabcore.Connection{
		ID:               to.Ptr(testhelp.RandomUUID()),
		DisplayName:      to.Ptr(testhelp.RandomName()),
		GatewayID:        to.Ptr(testhelp.RandomUUID()),
		PrivacyLevel:     to.Ptr(fabcore.PrivacyLevelPrivate),
		ConnectivityType: to.Ptr(fabcore.ConnectivityTypeVirtualNetworkGateway),
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

// ========== Credential Type Functions ==========

// NewRandomConnectionWithBasicCredentials creates a connection with Basic credentials.
func NewRandomConnectionWithBasicCredentials() fabcore.Connection {
	return fabcore.Connection{
		ID:               to.Ptr(testhelp.RandomUUID()),
		DisplayName:      to.Ptr(testhelp.RandomName()),
		PrivacyLevel:     to.Ptr(fabcore.PrivacyLevelPrivate),
		ConnectivityType: to.Ptr(fabcore.ConnectivityTypeShareableCloud),
		ConnectionDetails: &fabcore.ListConnectionDetails{
			Path: to.Ptr(testhelp.RandomURI()),
			Type: to.Ptr("GitHubSourceControl"),
		},
		CredentialDetails: &fabcore.ListCredentialDetails{
			CredentialType:       to.Ptr(fabcore.CredentialTypeBasic),
			SingleSignOnType:     to.Ptr(fabcore.SingleSignOnTypeNone),
			ConnectionEncryption: to.Ptr(fabcore.ConnectionEncryptionNotEncrypted),
			SkipTestConnection:   to.Ptr(false),
		},
	}
}

// NewRandomConnectionWithKeyCredentials creates a connection with Key credentials.
func NewRandomConnectionWithKeyCredentials() fabcore.Connection {
	return fabcore.Connection{
		ID:               to.Ptr(testhelp.RandomUUID()),
		DisplayName:      to.Ptr(testhelp.RandomName()),
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

// NewRandomConnectionWithServicePrincipalCredentials creates a connection with ServicePrincipal credentials.
func NewRandomConnectionWithServicePrincipalCredentials() fabcore.Connection {
	return fabcore.Connection{
		ID:               to.Ptr(testhelp.RandomUUID()),
		DisplayName:      to.Ptr(testhelp.RandomName()),
		PrivacyLevel:     to.Ptr(fabcore.PrivacyLevelPrivate),
		ConnectivityType: to.Ptr(fabcore.ConnectivityTypeShareableCloud),
		ConnectionDetails: &fabcore.ListConnectionDetails{
			Path: to.Ptr(testhelp.RandomURI()),
			Type: to.Ptr("GitHubSourceControl"),
		},
		CredentialDetails: &fabcore.ListCredentialDetails{
			CredentialType:       to.Ptr(fabcore.CredentialTypeServicePrincipal),
			SingleSignOnType:     to.Ptr(fabcore.SingleSignOnTypeNone),
			ConnectionEncryption: to.Ptr(fabcore.ConnectionEncryptionNotEncrypted),
			SkipTestConnection:   to.Ptr(false),
		},
	}
}

// NewRandomConnectionWithAnonymousCredentials creates a connection with Anonymous credentials.
func NewRandomConnectionWithAnonymousCredentials() fabcore.Connection {
	return fabcore.Connection{
		ID:               to.Ptr(testhelp.RandomUUID()),
		DisplayName:      to.Ptr(testhelp.RandomName()),
		PrivacyLevel:     to.Ptr(fabcore.PrivacyLevelPrivate),
		ConnectivityType: to.Ptr(fabcore.ConnectivityTypeShareableCloud),
		ConnectionDetails: &fabcore.ListConnectionDetails{
			Path: to.Ptr(testhelp.RandomURI()),
			Type: to.Ptr("GitHubSourceControl"),
		},
		CredentialDetails: &fabcore.ListCredentialDetails{
			CredentialType:       to.Ptr(fabcore.CredentialTypeAnonymous),
			SingleSignOnType:     to.Ptr(fabcore.SingleSignOnTypeNone),
			ConnectionEncryption: to.Ptr(fabcore.ConnectionEncryptionNotEncrypted),
			SkipTestConnection:   to.Ptr(false),
		},
	}
}

// NewRandomConnectionWithWorkspaceIdentityCredentials creates a connection with WorkspaceIdentity credentials.
func NewRandomConnectionWithWorkspaceIdentityCredentials() fabcore.Connection {
	return fabcore.Connection{
		ID:               to.Ptr(testhelp.RandomUUID()),
		DisplayName:      to.Ptr(testhelp.RandomName()),
		PrivacyLevel:     to.Ptr(fabcore.PrivacyLevelPrivate),
		ConnectivityType: to.Ptr(fabcore.ConnectivityTypeShareableCloud),
		ConnectionDetails: &fabcore.ListConnectionDetails{
			Path: to.Ptr(testhelp.RandomURI()),
			Type: to.Ptr("GitHubSourceControl"),
		},
		CredentialDetails: &fabcore.ListCredentialDetails{
			CredentialType:       to.Ptr(fabcore.CredentialTypeWorkspaceIdentity),
			SingleSignOnType:     to.Ptr(fabcore.SingleSignOnTypeNone),
			ConnectionEncryption: to.Ptr(fabcore.ConnectionEncryptionNotEncrypted),
			SkipTestConnection:   to.Ptr(false),
		},
	}
}

// NewRandomConnectionWithSharedAccessSignatureCredentials creates a connection with SharedAccessSignature credentials.
func NewRandomConnectionWithSharedAccessSignatureCredentials() fabcore.Connection {
	return fabcore.Connection{
		ID:               to.Ptr(testhelp.RandomUUID()),
		DisplayName:      to.Ptr(testhelp.RandomName()),
		PrivacyLevel:     to.Ptr(fabcore.PrivacyLevelPrivate),
		ConnectivityType: to.Ptr(fabcore.ConnectivityTypeShareableCloud),
		ConnectionDetails: &fabcore.ListConnectionDetails{
			Path: to.Ptr(testhelp.RandomURI()),
			Type: to.Ptr("GitHubSourceControl"),
		},
		CredentialDetails: &fabcore.ListCredentialDetails{
			CredentialType:       to.Ptr(fabcore.CredentialTypeSharedAccessSignature),
			SingleSignOnType:     to.Ptr(fabcore.SingleSignOnTypeNone),
			ConnectionEncryption: to.Ptr(fabcore.ConnectionEncryptionNotEncrypted),
			SkipTestConnection:   to.Ptr(false),
		},
	}
}
