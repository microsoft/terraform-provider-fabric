// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceiar

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Workspace Inbound Azure Resource Rules",
	Type:           "workspace_inbound_azure_resource_rules",
	DocsURL:        "https://learn.microsoft.com/fabric/onelake/onelake-manage-inbound-access-trusted-resources",
	IsPreview:      true,
	IsSPNSupported: true,
}
