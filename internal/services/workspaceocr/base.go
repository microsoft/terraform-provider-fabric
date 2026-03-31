// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceocr

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Workspace Outbound Cloud Connection Rules",
	Type:           "workspace_outbound_cloud_connection_rules",
	DocsURL:        "https://learn.microsoft.com/fabric/security/workspace-outbound-access-protection-allow-list-connector",
	IsPreview:      false,
	IsSPNSupported: true,
}
