// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package tenantsetting

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Tenant Setting",
	Type:           "tenant_setting",
	Names:          "Tenant Settings",
	Types:          "tenant_settings",
	DocsURL:        "https://learn.microsoft.com/fabric/admin/tenant-settings-index",
	IsPreview:      false,
	IsSPNSupported: true,
}
