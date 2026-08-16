// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspacefwr

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Workspace Firewall Rules",
	Type:           "workspace_firewall_rules",
	DocsURL:        "https://learn.microsoft.com/fabric/security/security-workspace-level-firewall-overview",
	IsPreview:      true,
	IsSPNSupported: true,
}
