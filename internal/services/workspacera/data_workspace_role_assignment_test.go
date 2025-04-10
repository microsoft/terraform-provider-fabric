// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspacera_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_WorkspaceRoleAssignmentDataSource(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	entity := NewRandomWorkspaceRoleAssignment()
	fakes.FakeServer.ServerFactory.Core.WorkspacesServer.GetWorkspaceRoleAssignment = fakeWorkspaceRoleAssignment(entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
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
				testDataSourceItemHeader,
				map[string]any{
					"id":           *entity.ID,
					"workspace_id": workspaceID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "workspace_id", workspaceID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "principal.id", entity.Principal.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "principal.type", (*string)(entity.Principal.Type)),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "role", (*string)(entity.Role)),
			),
		},
	}))
}

func TestAcc_WorkspaceRoleAssignmentDataSource(t *testing.T) {
	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, capacityID)

	principal := testhelp.WellKnown()["Principal"].(map[string]any)
	principalID := principal["id"].(string)
	principalType := principal["type"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"principal": map[string]any{
							"id":   principalID,
							"type": principalType,
						},
						"role": "Member",
					},
				),
				at.CompileConfig(
					testDataSourceItemHeader,
					map[string]any{
						"id":           principalID,
						"workspace_id": testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"depends_on" : []string{testResourceItemFQN},
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", principalID),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "workspace_id"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "principal.id", principalID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "principal.type", principalType),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "role", "Member"),
			),
		},
	},
	))
}
