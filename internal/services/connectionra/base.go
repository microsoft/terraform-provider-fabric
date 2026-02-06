// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package connectionra

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Connection Role Assignment",
	Type:           "connection_role_assignment",
	Names:          "Connection Role Assignments",
	Types:          "connection_role_assignments",
	DocsURL:        "https://learn.microsoft.com/fabric/data-factory/data-source-management",
	IsPreview:      true,
	IsSPNSupported: true,
}
