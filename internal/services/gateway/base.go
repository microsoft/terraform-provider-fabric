// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	OnPremisesItemTFType        = "OnPremisesGateway"
	OnPremisesItemsTFType       = "OnPremisesGateways"
	OnPremisesPersonalItemType  = "OnPremisesPersonalGateway"
	OnPremisesPersonalItemsType = "OnPremisesPersonalGateways"
	VirtualNetworkItemTFType    = "VirtualNetworkGateway"
	VirtualNetworkItemsTFType   = "VirtualNetworkGateways"
	ItemName                    = "Gateway"
	ItemsName                   = "Gateways"
	ItemsTFName                 = "gateways"
	ItemDocsSPNSupport          = common.DocsSPNSupported
	ItemDocsURL                 = "https://learn.microsoft.com/en-us/fabric/data-factory/how-to-access-on-premises-data"
)
