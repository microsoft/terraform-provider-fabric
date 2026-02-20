// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspacencp

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Workspace Network Communication Policy",
	Type:           "workspace_network_communication_policy",
	DocsURL:        "https://learn.microsoft.com/fabric/security/",
	IsPreview:      true,
	IsSPNSupported: true,
}
