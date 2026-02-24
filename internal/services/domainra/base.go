// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package domainra

import (
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"

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

const (
	// DomainRoleAdmins- Domain admins.
	DomainRoleAdmins fabadmin.DomainRole = "Admins"
	// DomainRoleContributors - Domain contributors.
	DomainRoleContributors fabadmin.DomainRole = "Contributors"
)

// PossibleDomainRoleValues returns the possible values for the DomainRole const type.
func PossibleDomainRoleValues() []fabadmin.DomainRole {
	return []fabadmin.DomainRole{
		DomainRoleAdmins,
		DomainRoleContributors,
	}
}
