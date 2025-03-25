// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domainwa

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Domain Workspace Assignment",
	Type:           "domain_workspace_assignment",
	Names:          "Domain Workspace Assignments",
	Types:          "domain_workspace_assignments",
	DocsURL:        "https://learn.microsoft.com/fabric/governance/domains",
	IsPreview:      true,
	IsSPNSupported: true,
}
