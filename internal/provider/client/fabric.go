// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/microsoft/fabric-sdk-go/fabric"

	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

type ProviderWithFabricClient interface {
	provider.Provider
	ConfigureCreateClient(createClient func(ctx context.Context, cfg *pconfig.ProviderConfig) (*fabric.Client, error))
	ConfigureAffirmProviderConfig(validateConfig func(cfg *pconfig.ProviderConfig))
}
