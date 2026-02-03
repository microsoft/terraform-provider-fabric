// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// https://github.com/radius-project/radius/blob/main/pkg/azure/armauth/auth.go
// https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azidentity/README.md#credential-types

type CredentialResponse struct {
	Cred       azcore.TokenCredential
	AuthMethod AuthenticationMethod
	Info       string
}

func newCredentialResponse(cred azcore.TokenCredential, method AuthenticationMethod, info string) CredentialResponse {
	return CredentialResponse{
		Cred:       cred,
		AuthMethod: method,
		Info:       info,
	}
}

// NewCredential evaluates the authentication method and returns the appropriate credential.
func NewCredential(cfg Config) (CredentialResponse, error) {
	authMethod, info := getAuthMethod(cfg)

	switch authMethod {
	case AzureDevCLIAuth:
		cred, err := azidentity.NewAzureDeveloperCLICredential(&azidentity.AzureDeveloperCLICredentialOptions{
			AdditionallyAllowedTenants: cfg.AuxiliaryTenantIDs,
			TenantID:                   cfg.TenantID,
		})

		return newCredentialResponse(cred, authMethod, info), err

	case ManagedServiceIdentityUserAuth:
		clientID := azidentity.ClientID(cfg.ClientID)
		cred, err := azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
			ID: clientID,
			ClientOptions: azcore.ClientOptions{
				Cloud: cfg.Environment,
			},
		})

		return newCredentialResponse(cred, authMethod, info), err

	case ManagedServiceIdentitySystemAuth:
		cred, err := azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
			ClientOptions: azcore.ClientOptions{
				Cloud: cfg.Environment,
			},
		})

		return newCredentialResponse(cred, authMethod, info), err

	case AzureDevOpsWorkloadIdentityFederationAuth:
		cred, err := azidentity.NewAzurePipelinesCredential(cfg.TenantID, cfg.ClientID, cfg.OIDC.AzureDevOpsServiceConnectionID, cfg.OIDC.RequestToken, &azidentity.AzurePipelinesCredentialOptions{
			AdditionallyAllowedTenants: cfg.AuxiliaryTenantIDs,
		})

		return newCredentialResponse(cred, authMethod, info), err

	case ServicePrincipalOIDCAuth:
		o := cfg.OIDC
		cred, err := azidentity.NewClientAssertionCredential(cfg.TenantID, cfg.ClientID, o.getAssertion, &azidentity.ClientAssertionCredentialOptions{
			AdditionallyAllowedTenants: cfg.AuxiliaryTenantIDs,
			ClientOptions: azcore.ClientOptions{
				Cloud: cfg.Environment,
			},
		})

		return newCredentialResponse(cred, authMethod, info), err

	case ServicePrincipalCertificateAuth:
		cred, err := azidentity.NewClientCertificateCredential(cfg.TenantID, cfg.ClientID, cfg.ClientCertificate, cfg.ClientCertificateKey, &azidentity.ClientCertificateCredentialOptions{
			AdditionallyAllowedTenants: cfg.AuxiliaryTenantIDs,
			ClientOptions: azcore.ClientOptions{
				Cloud: cfg.Environment,
			},
		})

		return newCredentialResponse(cred, authMethod, info), err

	case ServicePrincipalSecretAuth:
		cred, err := azidentity.NewClientSecretCredential(cfg.TenantID, cfg.ClientID, cfg.ClientSecret, &azidentity.ClientSecretCredentialOptions{
			AdditionallyAllowedTenants: cfg.AuxiliaryTenantIDs,
			ClientOptions: azcore.ClientOptions{
				Cloud: cfg.Environment,
			},
		})

		return newCredentialResponse(cred, authMethod, info), err

	default:
		cred, err := azidentity.NewAzureCLICredential(&azidentity.AzureCLICredentialOptions{
			AdditionallyAllowedTenants: cfg.AuxiliaryTenantIDs,
			TenantID:                   cfg.TenantID,
		})

		return newCredentialResponse(cred, authMethod, info), err
	}
}

// getAuthMethod returns the authentication method to use based on ProviderCredentials data.
func getAuthMethod(cfg Config) (AuthenticationMethod, string) {
	switch {
	case cfg.UseMSI && cfg.ClientID != "":
		return ManagedServiceIdentityUserAuth, "Using User-Assigned Managed Identity (MSI) authentication"

	case cfg.UseMSI:
		return ManagedServiceIdentitySystemAuth, "Using System-Assigned Managed Identity (MSI) authentication"

	case cfg.UseOIDC && cfg.OIDC.AzureDevOpsServiceConnectionID != "":
		return AzureDevOpsWorkloadIdentityFederationAuth, "Using OpenID Connect (OIDC) authentication from the Azure DevOps Workload Identity Federation service connection."

	case cfg.UseOIDC:
		return ServicePrincipalOIDCAuth, "Using OpenID Connect (OIDC) authentication"

	case cfg.UseDevCLI:
		return AzureDevCLIAuth, "Using Azure Developer CLI authentication"

	case cfg.UseCLI:
		return AzureCLIAuth, "Using Azure CLI authentication"

	case cfg.ClientID != "" && cfg.ClientCertificate != nil:
		return ServicePrincipalCertificateAuth, "Using Service Principal Certificate authentication"

	case cfg.ClientID != "" && cfg.ClientSecret != "":
		return ServicePrincipalSecretAuth, "Using Service Principal Secret authentication"

	default:
		return AzureCLIAuth, "Using Azure CLI authentication"
	}
}
