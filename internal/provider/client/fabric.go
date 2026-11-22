// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/microsoft/fabric-sdk-go/fabric"

	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

type ProviderWithFabricClient interface {
	provider.Provider
	ConfigureCreateClient(createClient func(ctx context.Context, cfg *pconfig.ProviderConfig) (*fabric.Client, error))
	ConfigureAffirmProviderConfig(validateConfig func(cfg *pconfig.ProviderConfig))
}

// NewFabricClient applies the selected Fabric environment without overriding SDK endpoint defaults.
func NewFabricClient(credential azcore.TokenCredential, cfg *pconfig.ProviderConfig, options *fabric.ClientOptions) (*fabric.Client, error) {
	var clientOptions fabric.ClientOptions
	if options != nil {
		clientOptions = *options
	}

	clientOptions.Environment = cfg.Environment
	clientOptions.UseWorkspacePrivateLinks = cfg.UseWorkspacePrivateLinkEndpoint

	var endpoint *string
	if cfg.Endpoint != "" {
		endpoint = &cfg.Endpoint
	}

	return fabric.NewClient(credential, endpoint, &clientOptions)
}
