// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gatewaymbr

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Gateway Member",
	Type:           "gateway_member",
	Names:          "Gateway Members",
	Types:          "gateway_members",
	DocsURL:        "https://learn.microsoft.com/data-integration/gateway/service-gateway-onprem",
	IsPreview:      true,
	IsSPNSupported: true,
}
