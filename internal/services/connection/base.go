// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

// Define a custom Connection item type since it's not in the SDK yet.
const FabricItemType = fabcore.ItemType("Connection")

// CredentialType defines supported credential types for connections.
type CredentialType string

const (
	ConnectionCredentialTypeAnonymous                   CredentialType = "Anonymous"
	ConnectionCredentialTypeBasic                       CredentialType = "Basic"
	ConnectionCredentialTypeKey                         CredentialType = "Key"
	ConnectionCredentialTypeOAuth2                      CredentialType = "OAuth2"
	ConnectionCredentialTypeServicePrincipal            CredentialType = "ServicePrincipal"
	ConnectionCredentialTypeSharedAccessSignature       CredentialType = "SharedAccessSignature"
	ConnectionCredentialTypeWindows                     CredentialType = "Windows"
	ConnectionCredentialTypeWindowsWithoutImpersonation CredentialType = "WindowsWithoutImpersonation"
	ConnectionCredentialTypeWorkspaceIdentity           CredentialType = "WorkspaceIdentity"
)

// SingleSignOnType defines SSO types for connections.
type SingleSignOnType string

const (
	SingleSignOnTypeNone                       SingleSignOnType = "None"
	SingleSignOnTypeKerberos                   SingleSignOnType = "Kerberos"
	SingleSignOnTypeKerberosDirectQueryRefresh SingleSignOnType = "KerberosDirectQueryAndRefresh"
	SingleSignOnTypeMicrosoftEntraID           SingleSignOnType = "MicrosoftEntraID"
	SingleSignOnTypeSAML                       SingleSignOnType = "SecurityAssertionMarkupLanguage"
)

// Encryption defines encryption options for connections.
type Encryption string

const (
	ConnectionEncryptionAny          Encryption = "Any"
	ConnectionEncryptionEncrypted    Encryption = "Encrypted"
	ConnectionEncryptionNotEncrypted Encryption = "NotEncrypted"
)

// ConnectivityType defines connection connectivity types.
type ConnectivityType string

const (
	ConnectivityTypeShareableCloud        ConnectivityType = "ShareableCloud"
	ConnectivityTypeOnPremisesGateway     ConnectivityType = "OnPremisesGateway"
	ConnectivityTypeVirtualNetworkGateway ConnectivityType = "VirtualNetworkGateway"
)

// PrivacyLevel defines privacy levels for connections.
type PrivacyLevel string

const (
	PrivacyLevelOrganizational PrivacyLevel = "Organizational"
	PrivacyLevelPrivate        PrivacyLevel = "Private"
	PrivacyLevelPublic         PrivacyLevel = "Public"
	PrivacyLevelNone           PrivacyLevel = "None"
)

// ParameterDataType defines data types for connection parameters.
type ParameterDataType string

const (
	ParameterDataTypeText         ParameterDataType = "Text"
	ParameterDataTypeNumber       ParameterDataType = "Number"
	ParameterDataTypeBoolean      ParameterDataType = "Boolean"
	ParameterDataTypeDate         ParameterDataType = "Date"
	ParameterDataTypeDateTime     ParameterDataType = "DateTime"
	ParameterDataTypeTime         ParameterDataType = "Time"
	ParameterDataTypeDuration     ParameterDataType = "Duration"
	ParameterDataTypeDateTimeZone ParameterDataType = "DateTimeZone"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Connection",
	Type:           "connection",
	Names:          "Connections",
	Types:          "connections",
	DocsURL:        "https://learn.microsoft.com/fabric/connections/connection-types",
	IsPreview:      true,
	IsSPNSupported: true,
}
