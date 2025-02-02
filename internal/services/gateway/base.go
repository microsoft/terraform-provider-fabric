// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName           = "Gateway"
	ItemTFName         = "gateway"
	ItemsName          = "Gateways"
	ItemsTFName        = "gateways"
	ItemDocsSPNSupport = common.DocsSPNSupported
	ItemDocsURL        = "https://learn.microsoft.com/power-bi/guidance/powerbi-implementation-planning-data-gateways"
)

var (
	PossibleInactivityMinutesBeforeSleepValues = []int32{30, 60, 90, 120, 150, 240, 360, 480, 720, 1440} //nolint:gochecknoglobals
	MinNumberOfMemberGatewaysValues            = int32(1)                                                //nolint:gochecknoglobals
	MaxNumberOfMemberGatewaysValues            = int32(7)                                                //nolint:gochecknoglobals
)
