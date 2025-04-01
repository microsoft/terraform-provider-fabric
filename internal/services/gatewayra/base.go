// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gatewayra

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Gateway Role Assignment",
	Type:           "gateway_role_assignment",
	Names:          "Gateway Role Assignments",
	Types:          "gateway_role_assignments",
	DocsURL:        "https://learn.microsoft.com/power-bi/guidance/powerbi-implementation-planning-data-gateways",
	IsPreview:      true,
	IsSPNSupported: true,
}
