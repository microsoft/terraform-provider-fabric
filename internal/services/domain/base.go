// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Domain",
	Type:           "domain",
	Names:          "Domains",
	Types:          "domains",
	DocsURL:        "https://learn.microsoft.com/fabric/governance/domains",
	IsPreview:      false,
	IsSPNSupported: true,
}
