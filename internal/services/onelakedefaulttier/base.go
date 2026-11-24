// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedefaulttier

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "OneLake Default Tier",
	Type:           "workspace_onelake_default_tier",
	DocsURL:        "https://learn.microsoft.com/rest/api/fabric/core/onelake-settings",
	IsPreview:      false,
	IsSPNSupported: true,
}
