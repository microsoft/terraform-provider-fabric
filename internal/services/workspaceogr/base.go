// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceogr

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Workspace Outbound Gateway Rules",
	Type:           "workspace_outbound_gateway_rules",
	DocsURL:        "https://learn.microsoft.com/en-us/fabric/security/workspace-outbound-access-protection-overview",
	IsPreview:      true,
	IsSPNSupported: true,
}
