// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// Ensure TokenCredential implements azcore.TokenCredential interface.
var _ azcore.TokenCredential = (*TokenCredential)(nil)

// TokenCredential is a TokenCredential that returns a static bearer token.
type TokenCredential struct {
	token string
}

// GetToken returns the bearer token.
func (c *TokenCredential) GetToken(_ context.Context, _ policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token: c.token,
	}, nil
}

// NewTokenCredential creates a new instance of TokenCredential.
func NewTokenCredential(token string) (*TokenCredential, error) {
	return &TokenCredential{token: token}, nil
}