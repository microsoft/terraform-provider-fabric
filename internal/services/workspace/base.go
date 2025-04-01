// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Workspace",
	Type:           "workspace",
	Names:          "Workspaces",
	Types:          "workspaces",
	DocsURL:        "https://learn.microsoft.com/fabric/get-started/workspaces",
	IsPreview:      false,
	IsSPNSupported: true,
}

var workspaceIdentityTypes = []string{"SystemAssigned"} //nolint:gochecknoglobals
