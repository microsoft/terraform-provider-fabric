// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspacera

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Workspace Role Assignment",
	Type:           "workspace_role_assignment",
	Names:          "Workspace Role Assignments",
	Types:          "workspace_role_assignments",
	DocsURL:        "https://learn.microsoft.com/fabric/fundamentals/roles-workspaces",
	IsPreview:      false,
	IsSPNSupported: true,
}
