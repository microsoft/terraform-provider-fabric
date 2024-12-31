// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testDataSourceWorkspaceRoleAssignments       = testhelp.DataSourceFQN("fabric", workspaceRoleAssignmentsTFName, "test")
	testDataSourceWorkspaceRoleAssignmentsHeader = at.DataSourceHeader(testhelp.TypeName("fabric", workspaceRoleAssignmentsTFName), "test")
)

func TestUnit_WorkspaceRoleAssignmentsDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	workspaceRoleAssignments := NewRandomWorkspaceRoleAssignments()
	fakes.FakeServer.ServerFactory.Core.WorkspacesServer.NewListWorkspaceRoleAssignmentsPager = fakeWorkspaceRoleAssignments(workspaceRoleAssignments)

	entity := workspaceRoleAssignments.Value[1]

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceWorkspaceRoleAssignmentsHeader,
				map[string]any{
					"workspace_id":    workspaceID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceWorkspaceRoleAssignmentsHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceWorkspaceRoleAssignments, "workspace_id", workspaceID),
				resource.TestCheckResourceAttrPtr(testDataSourceWorkspaceRoleAssignments, "values.1.id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceWorkspaceRoleAssignments, "values.1.role", (*string)(entity.Role)),
				resource.TestCheckResourceAttrPtr(testDataSourceWorkspaceRoleAssignments, "values.1.display_name", entity.Principal.DisplayName),
				resource.TestCheckResourceAttrPtr(testDataSourceWorkspaceRoleAssignments, "values.1.type", (*string)(entity.Principal.Type)),
			),
		},
	}))
}

func TestAcc_WorkspaceRoleAssignmentsDataSource(t *testing.T) {
	workspace := testhelp.WellKnown()["WorkspaceDS"].(map[string]any)
	workspaceID := workspace["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceWorkspaceRoleAssignmentsHeader,
				map[string]any{
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceWorkspaceRoleAssignments, "workspace_id", workspaceID),
				resource.TestCheckResourceAttrSet(testDataSourceWorkspaceRoleAssignments, "values.0.id"),
			),
		},
	},
	))
}
