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
func (o *operationsConnection) GetID(entity fabcore.ConnectionClassification) string {
	return *entity.GetConnection().ID
}

// TransformCreate implements concreteOperations.
func (o *operationsConnection) TransformCreate(entity fabcore.ConnectionClassification) fabcore.ConnectionsClientCreateConnectionResponse {
	return fabcore.ConnectionsClientCreateConnectionResponse{
		ConnectionClassification: entity,
	}
}

// Create implements concreteOperations.
func (o *operationsConnection) Create(data fabcore.CreateConnectionRequestClassification) fabcore.ConnectionClassification {
	d := data.GetCreateConnectionRequest()

	switch req := data.(type) {
	case *fabcore.CreateCloudConnectionRequest:
		entity := &fabcore.ShareableCloudConnection{
			ID:                            to.Ptr(testhelp.RandomUUID()),
			DisplayName:                   d.DisplayName,
			PrivacyLevel:                  d.PrivacyLevel,
			ConnectivityType:              d.ConnectivityType,
			AllowConnectionUsageInGateway: req.AllowConnectionUsageInGateway,
			ConnectionDetails: &fabcore.ListConnectionDetails{
				Path: to.Ptr(testhelp.RandomURI()),
				Type: d.ConnectionDetails.Type,
			},
		}

		if req.CredentialDetails != nil {
			entity.CredentialDetails = &fabcore.ListCredentialDetails{
				CredentialType:       req.CredentialDetails.Credentials.GetCredentials().CredentialType,
				ConnectionEncryption: req.CredentialDetails.ConnectionEncryption,
				SingleSignOnType:     req.CredentialDetails.SingleSignOnType,
				SkipTestConnection:   req.CredentialDetails.SkipTestConnection,
			}
		}

		return entity

	case *fabcore.CreateVirtualNetworkGatewayConnectionRequest:
		entity := &fabcore.VirtualNetworkGatewayConnection{
			ID:               to.Ptr(testhelp.RandomUUID()),
			DisplayName:      d.DisplayName,
			PrivacyLevel:     d.PrivacyLevel,
			ConnectivityType: d.ConnectivityType,
			GatewayID:        req.GatewayID,
			ConnectionDetails: &fabcore.ListConnectionDetails{
				Path: to.Ptr(testhelp.RandomURI()),
				Type: d.ConnectionDetails.Type,
			},
		}

		if req.CredentialDetails != nil {
			entity.CredentialDetails = &fabcore.ListCredentialDetails{
				CredentialType:       req.CredentialDetails.Credentials.GetCredentials().CredentialType,
				ConnectionEncryption: req.CredentialDetails.ConnectionEncryption,
				SingleSignOnType:     req.CredentialDetails.SingleSignOnType,
				SkipTestConnection:   req.CredentialDetails.SkipTestConnection,
			}
		}

		return entity

	default:
		panic("Unsupported Connection type") // lintignore:R009
	}
}

// TransformGet implements concreteOperations.
func (o *operationsConnection) TransformGet(entity fabcore.ConnectionClassification) fabcore.ConnectionsClientGetConnectionResponse {
	return fabcore.ConnectionsClientGetConnectionResponse{
		ConnectionClassification: entity,
	}
}

// TransformList implements concreteOperations.
func (o *operationsConnection) TransformList(entities []fabcore.ConnectionClassification) fabcore.ConnectionsClientListConnectionsResponse {
	return fabcore.ConnectionsClientListConnectionsResponse{
		ListConnectionsResponse: fabcore.ListConnectionsResponse{
			Value: entities,
		},
	}
}

// TransformUpdate implements concreteOperations.
func (o *operationsConnection) TransformUpdate(entity fabcore.ConnectionClassification) fabcore.ConnectionsClientUpdateConnectionResponse {
	return fabcore.ConnectionsClientUpdateConnectionResponse{
		ConnectionClassification: entity,
	}
}

// Update implements concreteOperations.
func (o *operationsConnection) Update(base fabcore.ConnectionClassification, data fabcore.UpdateConnectionRequestClassification) fabcore.ConnectionClassification {
	d := data.GetUpdateConnectionRequest()

	updateCredentialDetails := func(credentialDetails *fabcore.UpdateCredentialDetails) *fabcore.ListCredentialDetails {
		if credentialDetails != nil {
			return &fabcore.ListCredentialDetails{
				CredentialType:       credentialDetails.Credentials.GetCredentials().CredentialType,
				ConnectionEncryption: credentialDetails.ConnectionEncryption,
				SingleSignOnType:     credentialDetails.SingleSignOnType,
				SkipTestConnection:   credentialDetails.SkipTestConnection,
			}
		}

		return nil
	}

	// Handle ShareableCloudConnection specific fields
	switch connection := base.(type) {
	case *fabcore.ShareableCloudConnection:
		if req, ok := data.(*fabcore.UpdateShareableCloudConnectionRequest); ok {
			connection.AllowConnectionUsageInGateway = req.AllowConnectionUsageInGateway
			connection.DisplayName = req.DisplayName
			connection.CredentialDetails = updateCredentialDetails(req.CredentialDetails)
		}

		connection.ConnectivityType = d.ConnectivityType
		connection.PrivacyLevel = d.PrivacyLevel

		return connection

	case *fabcore.VirtualNetworkGatewayConnection:
		if req, ok := data.(*fabcore.UpdateVirtualNetworkGatewayConnectionRequest); ok {
			connection.DisplayName = req.DisplayName
			connection.CredentialDetails = updateCredentialDetails(req.CredentialDetails)
		}

		connection.ConnectivityType = d.ConnectivityType
		connection.PrivacyLevel = d.PrivacyLevel

		return connection
	}

	panic("Unsupported Connection type") // lintignore:R009
}

// Validate implements concreteOperations.
func (o *operationsConnection) Validate(newEntity fabcore.ConnectionClassification, existing []fabcore.ConnectionClassification) (int, error) {
	for _, entity := range existing {
		if *entity.GetConnection().DisplayName == *newEntity.GetConnection().DisplayName {
			return http.StatusConflict, fabfake.SetResponseError(http.StatusConflict, fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error(), fabcore.ErrItem.ItemDisplayNameAlreadyInUse.Error())
		}
	}

	return http.StatusCreated, nil
}

func configureVirtualNetworkGatewayConnection(server *fakeServer) fabcore.VirtualNetworkGatewayConnection {
	configureConnection(server)

	return fabcore.VirtualNetworkGatewayConnection{}
}

func configureShareableCloudConnection(server *fakeServer) fabcore.ShareableCloudConnection {
	configureConnection(server)

	return fabcore.ShareableCloudConnection{}
}

func configureConnection(server *fakeServer) {
	type concreteEntityOperations interface {
		simpleIDOperations[
			fabcore.ConnectionClassification,
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
}

// FakeListSupportedConnectionTypes returns a fake handler for listing supported connection types.
func FakeListSupportedConnectionTypes() func(options *fabcore.ConnectionsClientListSupportedConnectionTypesOptions) azfake.PagerResponder[fabcore.ConnectionsClientListSupportedConnectionTypesResponse] {
	return func(_ *fabcore.ConnectionsClientListSupportedConnectionTypesOptions) azfake.PagerResponder[fabcore.ConnectionsClientListSupportedConnectionTypesResponse] {
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
							fabcore.CredentialTypeAnonymous,
						},
						SupportedConnectionEncryptionTypes: []fabcore.ConnectionEncryption{
							fabcore.ConnectionEncryptionNotEncrypted,
						},
						SupportsSkipTestConnection: to.Ptr(true),
					},
					{
						Type: to.Ptr("ConnectionWithEmptyParametersList"),
						CreationMethods: []fabcore.ConnectionCreationMethod{
							{
								Name:       to.Ptr("ConnectionWithEmptyParametersList.Actions"),
								Parameters: []fabcore.ConnectionCreationParameter{},
							},
						},
						SupportedCredentialTypes: []fabcore.CredentialType{
							fabcore.CredentialTypeAnonymous,
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
func NewRandomConnection() fabcore.ConnectionClassification {
	connectivitytypes := []fabcore.ConnectivityType{
		fabcore.ConnectivityTypeShareableCloud,
		fabcore.ConnectivityTypeVirtualNetworkGateway,
	}

	randomConnectivityType := testhelp.RandomElement(connectivitytypes)

	if randomConnectivityType == fabcore.ConnectivityTypeVirtualNetworkGateway {
		return NewRandomVirtualNetworkGatewayConnection()
	}

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
		panic("Unsupported credential type") // lintignore:R009
	}
}

// ========== Connectivity Type Functions ==========.
func NewRandomShareableCloudConnection() *fabcore.ShareableCloudConnection {
	return &fabcore.ShareableCloudConnection{
		ID:                            to.Ptr(testhelp.RandomUUID()),
		DisplayName:                   to.Ptr(testhelp.RandomName()),
		PrivacyLevel:                  to.Ptr(fabcore.PrivacyLevelPrivate),
		ConnectivityType:              to.Ptr(fabcore.ConnectivityTypeShareableCloud),
		AllowConnectionUsageInGateway: to.Ptr(testhelp.RandomBool()),
		ConnectionDetails: &fabcore.ListConnectionDetails{
			Path: to.Ptr(testhelp.RandomURI()),
			Type: to.Ptr("GitHubSourceControl"),
		},
		CredentialDetails: &fabcore.ListCredentialDetails{
			CredentialType:       to.Ptr(fabcore.CredentialTypeKey),
			SingleSignOnType:     to.Ptr(fabcore.SingleSignOnTypeNone),
			ConnectionEncryption: to.Ptr(fabcore.ConnectionEncryptionNotEncrypted),
			SkipTestConnection:   to.Ptr(testhelp.RandomBool()),
		},
	}
}

func NewRandomVirtualNetworkGatewayConnection() *fabcore.VirtualNetworkGatewayConnection {
	return &fabcore.VirtualNetworkGatewayConnection{
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
			SkipTestConnection:   to.Ptr(testhelp.RandomBool()),
		},
	}
}

// ========== Credential Type Functions ==========

// NewRandomConnectionWithBasicCredentials creates a connection with Basic credentials.
func NewRandomConnectionWithBasicCredentials() *fabcore.ShareableCloudConnection {
	return &fabcore.ShareableCloudConnection{
		ID:                            to.Ptr(testhelp.RandomUUID()),
		DisplayName:                   to.Ptr(testhelp.RandomName()),
		PrivacyLevel:                  to.Ptr(fabcore.PrivacyLevelPrivate),
		ConnectivityType:              to.Ptr(fabcore.ConnectivityTypeShareableCloud),
		AllowConnectionUsageInGateway: to.Ptr(testhelp.RandomBool()),
		ConnectionDetails: &fabcore.ListConnectionDetails{
			Path: to.Ptr(testhelp.RandomURI()),
			Type: to.Ptr("GitHubSourceControl"),
		},
		CredentialDetails: &fabcore.ListCredentialDetails{
			CredentialType:       to.Ptr(fabcore.CredentialTypeBasic),
			SingleSignOnType:     to.Ptr(fabcore.SingleSignOnTypeNone),
			ConnectionEncryption: to.Ptr(fabcore.ConnectionEncryptionNotEncrypted),
			SkipTestConnection:   to.Ptr(testhelp.RandomBool()),
		},
	}
}

// NewRandomConnectionWithKeyCredentials creates a connection with Key credentials.
func NewRandomConnectionWithKeyCredentials() *fabcore.ShareableCloudConnection {
	return &fabcore.ShareableCloudConnection{
		ID:                            to.Ptr(testhelp.RandomUUID()),
		DisplayName:                   to.Ptr(testhelp.RandomName()),
		PrivacyLevel:                  to.Ptr(fabcore.PrivacyLevelPrivate),
		ConnectivityType:              to.Ptr(fabcore.ConnectivityTypeShareableCloud),
		AllowConnectionUsageInGateway: to.Ptr(testhelp.RandomBool()),
		ConnectionDetails: &fabcore.ListConnectionDetails{
			Path: to.Ptr(testhelp.RandomURI()),
			Type: to.Ptr("GitHubSourceControl"),
		},
		CredentialDetails: &fabcore.ListCredentialDetails{
			CredentialType:       to.Ptr(fabcore.CredentialTypeKey),
			SingleSignOnType:     to.Ptr(fabcore.SingleSignOnTypeNone),
			ConnectionEncryption: to.Ptr(fabcore.ConnectionEncryptionNotEncrypted),
			SkipTestConnection:   to.Ptr(testhelp.RandomBool()),
		},
	}
}

// NewRandomConnectionWithServicePrincipalCredentials creates a connection with ServicePrincipal credentials.
func NewRandomConnectionWithServicePrincipalCredentials() *fabcore.ShareableCloudConnection {
	return &fabcore.ShareableCloudConnection{
		ID:                            to.Ptr(testhelp.RandomUUID()),
		DisplayName:                   to.Ptr(testhelp.RandomName()),
		PrivacyLevel:                  to.Ptr(fabcore.PrivacyLevelPrivate),
		ConnectivityType:              to.Ptr(fabcore.ConnectivityTypeShareableCloud),
		AllowConnectionUsageInGateway: to.Ptr(testhelp.RandomBool()),
		ConnectionDetails: &fabcore.ListConnectionDetails{
			Path: to.Ptr(testhelp.RandomURI()),
			Type: to.Ptr("GitHubSourceControl"),
		},
		CredentialDetails: &fabcore.ListCredentialDetails{
			CredentialType:       to.Ptr(fabcore.CredentialTypeServicePrincipal),
			SingleSignOnType:     to.Ptr(fabcore.SingleSignOnTypeNone),
			ConnectionEncryption: to.Ptr(fabcore.ConnectionEncryptionNotEncrypted),
			SkipTestConnection:   to.Ptr(testhelp.RandomBool()),
		},
	}
}

// NewRandomConnectionWithAnonymousCredentials creates a connection with Anonymous credentials.
func NewRandomConnectionWithAnonymousCredentials() *fabcore.ShareableCloudConnection {
	return &fabcore.ShareableCloudConnection{
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
			SkipTestConnection:   to.Ptr(testhelp.RandomBool()),
		},
	}
}

// NewRandomConnectionWithWorkspaceIdentityCredentials creates a connection with WorkspaceIdentity credentials.
func NewRandomConnectionWithWorkspaceIdentityCredentials() *fabcore.ShareableCloudConnection {
	return &fabcore.ShareableCloudConnection{
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
			SkipTestConnection:   to.Ptr(testhelp.RandomBool()),
		},
	}
}

// NewRandomConnectionWithSharedAccessSignatureCredentials creates a connection with SharedAccessSignature credentials.
func NewRandomConnectionWithSharedAccessSignatureCredentials() *fabcore.ShareableCloudConnection {
	return &fabcore.ShareableCloudConnection{
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
			SkipTestConnection:   to.Ptr(testhelp.RandomBool()),
		},
	}
}
