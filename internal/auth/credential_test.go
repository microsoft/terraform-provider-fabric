// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package auth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/auth"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func TestUnit_NewCredential(t *testing.T) {
	certPass := testhelp.RandomName()
	cert, key, _ := auth.ConvertBase64ToCert(testhelp.RandomP12CertB64(certPass), certPass)

	testCases := map[string]struct {
		cfg         auth.Config
		expected    auth.AuthenticationMethod
		expectError bool
		envs        map[string]string
	}{
		"AzureDevCLIAuth valid": {
			cfg: auth.Config{
				UseDevCLI: true,
			},
			expected:    auth.AzureDevCLIAuth,
			expectError: false,
		},
		"ManagedServiceIdentityUserAuth valid": {
			cfg: auth.Config{
				UseMSI:   true,
				ClientID: testhelp.RandomUUID(),
			},
			expected:    auth.ManagedServiceIdentityUserAuth,
			expectError: false,
		},
		"ManagedServiceIdentitySystemAuth valid": {
			cfg: auth.Config{
				UseMSI: true,
			},
			expected:    auth.ManagedServiceIdentitySystemAuth,
			expectError: false,
		},
		"AzureDevOpsWorkloadIdentityFederationAuth valid": {
			cfg: auth.Config{
				UseOIDC:  true,
				TenantID: testhelp.RandomUUID(),
				ClientID: testhelp.RandomUUID(),
				OIDC: auth.OIDCConfig{
					RequestToken:                   "test-token",
					AzureDevOpsServiceConnectionID: testhelp.RandomUUID(),
				},
			},
			expected:    auth.AzureDevOpsWorkloadIdentityFederationAuth,
			expectError: false,
			envs: map[string]string{
				"SYSTEM_OIDCREQUESTURI": "https://example.com",
			},
		},
		"AzureDevOpsWorkloadIdentityFederationAuth invalid": {
			cfg: auth.Config{
				UseOIDC:  true,
				ClientID: testhelp.RandomUUID(),
				OIDC: auth.OIDCConfig{
					RequestToken:                   "test-token",
					AzureDevOpsServiceConnectionID: testhelp.RandomUUID(),
				},
			},
			expected:    auth.AzureDevOpsWorkloadIdentityFederationAuth,
			expectError: true,
		},
		"ServicePrincipalOIDCAuth valid": {
			cfg: auth.Config{
				UseOIDC:  true,
				TenantID: testhelp.RandomUUID(),
				OIDC: auth.OIDCConfig{
					RequestToken: "test-token",
				},
			},
			expected:    auth.ServicePrincipalOIDCAuth,
			expectError: false,
		},
		"ServicePrincipalOIDCAuth invalid": {
			cfg: auth.Config{
				UseOIDC: true,
				OIDC: auth.OIDCConfig{
					RequestToken: "test-token",
				},
			},
			expected:    auth.ServicePrincipalOIDCAuth,
			expectError: true,
		},
		"ServicePrincipalCertificateAuth valid": {
			cfg: auth.Config{
				TenantID:             testhelp.RandomUUID(),
				ClientID:             testhelp.RandomUUID(),
				ClientCertificate:    cert,
				ClientCertificateKey: key,
			},
			expected:    auth.ServicePrincipalCertificateAuth,
			expectError: false,
		},
		"ServicePrincipalCertificateAuth invalid": {
			cfg: auth.Config{
				ClientID:             testhelp.RandomUUID(),
				ClientCertificate:    cert,
				ClientCertificateKey: key,
			},
			expected:    auth.ServicePrincipalCertificateAuth,
			expectError: true,
		},
		"ServicePrincipalSecretAuth valid": {
			cfg: auth.Config{
				TenantID:     testhelp.RandomUUID(),
				ClientID:     testhelp.RandomUUID(),
				ClientSecret: "test-client-secret",
			},
			expected:    auth.ServicePrincipalSecretAuth,
			expectError: false,
		},
		"ServicePrincipalSecretAuth invalid": {
			cfg: auth.Config{
				ClientID:     testhelp.RandomUUID(),
				ClientSecret: "test-client-secret",
			},
			expected:    auth.ServicePrincipalSecretAuth,
			expectError: true,
		},
		"AzureCLIAuth valid": {
			cfg: auth.Config{
				UseCLI: true,
			},
			expected:    auth.AzureCLIAuth,
			expectError: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			for k, v := range testCase.envs {
				t.Setenv(k, v)
			}

			credResponse, err := auth.NewCredential(testCase.cfg)
			if testCase.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.expected, credResponse.AuthMethod)
				assert.NotNil(t, credResponse.Cred)
			}
		})
	}
}
