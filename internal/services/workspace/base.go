// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"github.com/microsoft/terraform-provider-fabric/internal/common"
)

const (
	ItemName                       = "Workspace"
	ItemTFName                     = "workspace"
	ItemsName                      = "Workspaces"
	ItemsTFName                    = "workspaces"
	ItemDocsSPNSupport             = common.DocsSPNSupported
	ItemDocsURL                    = "https://learn.microsoft.com/fabric/get-started/workspaces"
	WorkspaceRoleAssignmentName    = "Workspace Role Assignment"
	WorkspaceRoleAssignmentTFName  = "workspace_role_assignment"
	WorkspaceRoleAssignmentsName   = "Workspace Role Assignments"
	WorkspaceRoleAssignmentsTFName = "workspace_role_assignments"
	WorkspaceRoleAssignmentDocsURL = "https://learn.microsoft.com/fabric/fundamentals/roles-workspaces"
	WorkspaceGitName               = "Workspace Git integration"
	WorkspaceGitTFName             = "workspace_git"
	WorkspaceGitDocsURL            = "https://learn.microsoft.com/fabric/cicd/git-integration/intro-to-git-integration"
)

var workspaceIdentityTypes = []string{"SystemAssigned"} //nolint:gochecknoglobals
