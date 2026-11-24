// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspaceonelakeimmutability

import "github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"

const maxRetentionDays = int64(1<<31 - 1)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Workspace OneLake Immutability Policy",
	Type:           "workspace_onelake_immutability_policy",
	DocsURL:        "https://learn.microsoft.com/rest/api/fabric/core/onelake-settings/modify-immutability-policy",
	IsPreview:      false,
	IsSPNSupported: true,
}
