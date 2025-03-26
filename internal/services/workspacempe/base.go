// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspacempe

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Workspace Managed Private Endpoint",
	Type:           "workspace_managed_private_endpoint",
	Names:          "Workspace Managed Private Endpoints",
	Types:          "workspace_managed_private_endpoints",
	DocsURL:        "https://learn.microsoft.com/fabric/security/security-managed-private-endpoints-overview",
	IsPreview:      true,
	IsSPNSupported: true,
}
