// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakediagnosticsetting

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

const (
	diagnosticStatusEnabled  = "Enabled"
	diagnosticStatusDisabled = "Disabled"
	destinationTypeLakehouse = "Lakehouse"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "OneLake Diagnostic Setting",
	Type:           "onelake_diagnostic_setting",
	DocsURL:        "https://learn.microsoft.com/rest/api/fabric/core/onelake-settings/modify-diagnostics",
	IsPreview:      false,
	IsSPNSupported: true,
}
