// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	OnPremisesItemTFType        = "on_premises_gateway"
	OnPremisesItemsTFType       = "on_premises_gateways"
	OnPremisesPersonalItemType  = "on_premises_personal_gateway"
	OnPremisesPersonalItemsType = "on_premises_personal_gateways"
	VirtualNetworkItemTFType    = "virtual_network_gateway"
	VirtualNetworkItemsTFType   = "virtual_network_gateways"
	ItemName                    = "Gateway"
	ItemsName                   = "Gateways"
	ItemsTFName                 = "gateways"
	ItemDocsSPNSupport          = common.DocsSPNSupported
	ItemDocsURL                 = "https://learn.microsoft.com/en-us/fabric/data-factory/how-to-access-on-premises-data"
)

var (
	PossibleInactivityMinutesBeforeSleepValues = []int32{30, 60, 90, 120, 150, 240, 360, 480, 720, 1440}

	MinNumberOfMemberGatewaysValues = int32(1)

	MaxNumberOfMemberGatewaysValues = int32(7)
)
