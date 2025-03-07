// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package domain

import (
	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName                         = "Domain"
	ItemTFName                       = "domain"
	ItemsName                        = "Domains"
	ItemsTFName                      = "domains"
	ItemDocsSPNSupport               = common.DocsSPNSupported
	ItemDocsURL                      = "https://learn.microsoft.com/fabric/governance/domains"
	DomainWorkspaceAssignmentsName   = "Domain Workspace Assignments"
	DomainWorkspaceAssignmentsTFName = "domain_workspace_assignments"
	DomainRoleAssignmentsName        = "Domain Role Assignments"
	DomainRoleAssignmentsTFName      = "domain_role_assignments"
	ItemPreview                      = true
)
