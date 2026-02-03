// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package connection

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Connection",
	Type:           "connection",
	Names:          "Connections",
	Types:          "connections",
	DocsURL:        "https://learn.microsoft.com/fabric/data-factory/data-source-management",
	IsPreview:      true,
	IsSPNSupported: true,
}
