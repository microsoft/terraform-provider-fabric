// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"crypto"
	"crypto/x509"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
)

// Config represents the authentication configuration.
type Config struct {
	UseCLI               bool
	UseDevCLI            bool
	UseOIDC              bool
	UseMSI               bool
	Environment          cloud.Configuration
	AuxiliaryTenantIDs   []string
	TenantID             string
	ClientID             string
	ClientSecret         string
	ClientCertificate    []*x509.Certificate
	ClientCertificateKey crypto.PrivateKey
	OIDC                 OIDCConfig
}

// OIDCConfig represents the OpenID Connect configuration.
type OIDCConfig struct {
	RequestURL                     string
	RequestToken                   string
	Token                          string
	AzureDevOpsServiceConnectionID string
}

type AuthenticationMethod string

// Supported authentication methods.
const (
	ServicePrincipalSecretAuth                AuthenticationMethod = "ServicePrincipalSecret"
	ServicePrincipalCertificateAuth           AuthenticationMethod = "ServicePrincipalCertificate"
	ServicePrincipalOIDCAuth                  AuthenticationMethod = "ServicePrincipalOIDC"
	AzureDevOpsWorkloadIdentityFederationAuth AuthenticationMethod = "AzureDevOpsWorkloadIdentityFederation"
	ManagedServiceIdentityUserAuth            AuthenticationMethod = "ManagedServiceIdentityUser"
	ManagedServiceIdentitySystemAuth          AuthenticationMethod = "ManagedServiceIdentitySystem"
	AzureCLIAuth                              AuthenticationMethod = "AzureCLI"
	AzureDevCLIAuth                           AuthenticationMethod = "AzureDeveloperCLI"
)
