// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var PossibleInactivityMinutesBeforeSleepValues = []int32{30, 60, 90, 120, 150, 240, 360, 480, 720, 1440} //nolint:gochecknoglobals

const (
	MinNumberOfMemberGatewaysValues int32 = 1
	MaxNumberOfMemberGatewaysValues int32 = 7
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Gateway",
	Type:           "gateway",
	Names:          "Gateways",
	Types:          "gateways",
	DocsURL:        "https://learn.microsoft.com/power-bi/guidance/powerbi-implementation-planning-data-gateways",
	IsPreview:      true,
	IsSPNSupported: true,
}
