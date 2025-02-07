// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package auth_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-fabric/internal/auth"
)

func TestUnit_TokenCredential_GetToken(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		token       string
		expected    azcore.AccessToken
		expectError bool
	}{
		"valid token": {
			token: "test-token",
			expected: azcore.AccessToken{
				Token: "test-token",
			},
			expectError: false,
		},
		"empty token": {
			token: "",
			expected: azcore.AccessToken{
				Token: "",
			},
			expectError: true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cred, err := auth.NewTokenCredential(testCase.token)
			if testCase.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{})
				require.NoError(t, err)
				assert.Equal(t, testCase.expected, token, "they should be equal")
			}
		})
	}
}
