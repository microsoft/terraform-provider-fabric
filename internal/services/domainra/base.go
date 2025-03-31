// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domainra

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Domain Role Assignment",
	Type:           "domain_role_assignment",
	Names:          "Domain Role Assignments",
	Types:          "domain_role_assignments",
	DocsURL:        "https://learn.microsoft.com/fabric/governance/domains",
	IsPreview:      true,
	IsSPNSupported: true,
}
