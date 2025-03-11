// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace_test

import (
	"fmt"
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testResourceWorkspaceRoleAssignmentFQN    = testhelp.ResourceFQN("fabric", workspaceRoleAssignmentTFName, "test")
	testResourceWorkspaceRoleAssignmentHeader = at.ResourceHeader(testhelp.TypeName("fabric", workspaceRoleAssignmentTFName), "test")
)

func TestUnit_WorkspaceRoleAssignmentResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceWorkspaceRoleAssignmentFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no required attributes - workspace_id
		{
			ResourceName: testResourceWorkspaceRoleAssignmentFQN,
			Config: at.CompileConfig(
				testResourceWorkspaceRoleAssignmentHeader,
				map[string]any{
					"principal_id":   "00000000-0000-0000-0000-000000000000",
					"principal_type": "User",
					"role":           "Member",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// error - no required attributes - principal_id
		{
			ResourceName: testResourceWorkspaceRoleAssignmentFQN,
			Config: at.CompileConfig(
				testResourceWorkspaceRoleAssignmentHeader,
				map[string]any{
					"workspace_id":   "00000000-0000-0000-0000-000000000000",
					"principal_type": "User",
					"role":           "Member",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "principal_id" is required, but no definition was found.`),
		},
		// error - no required attributes - principal_type
		{
			ResourceName: testResourceWorkspaceRoleAssignmentFQN,
			Config: at.CompileConfig(
				testResourceWorkspaceRoleAssignmentHeader,
				map[string]any{
					"workspace_id": "00000000-0000-0000-0000-000000000000",
					"principal_id": "00000000-0000-0000-0000-000000000000",
					"role":         "Member",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "principal_type" is required, but no definition was found.`),
		},
		// error - no required attributes - role
		{
			ResourceName: testResourceWorkspaceRoleAssignmentFQN,
			Config: at.CompileConfig(
				testResourceWorkspaceRoleAssignmentHeader,
				map[string]any{
					"workspace_id":   "00000000-0000-0000-0000-000000000000",
					"principal_id":   "00000000-0000-0000-0000-000000000000",
					"principal_type": "User",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "role" is required, but no definition was found.`),
		},
		// error - invalid UUID - workspace_id
		{
			ResourceName: testResourceWorkspaceRoleAssignmentFQN,
			Config: at.CompileConfig(
				testResourceWorkspaceRoleAssignmentHeader,
				map[string]any{
					"workspace_id":   "invalid uuid",
					"principal_id":   "00000000-0000-0000-0000-000000000000",
					"principal_type": "User",
					"role":           "Member",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - invalid UUID - principal_id
		{
			ResourceName: testResourceWorkspaceRoleAssignmentFQN,
			Config: at.CompileConfig(
				testResourceWorkspaceRoleAssignmentHeader,
				map[string]any{
					"workspace_id":   "00000000-0000-0000-0000-000000000000",
					"principal_id":   "invalid uuid",
					"principal_type": "User",
					"role":           "Member",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_WorkspaceRoleAssignmentResource_ImportState(t *testing.T) {
	testCase := at.CompileConfig(
		testResourceWorkspaceRoleAssignmentHeader,
		map[string]any{},
	)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceWorkspaceRoleAssignmentFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceWorkspaceRoleAssignmentFQN,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile("WorkspaceID/WorkspaceRoleAssignmentID"),
		},
		{
			ResourceName:  testResourceWorkspaceRoleAssignmentFQN,
			Config:        testCase,
			ImportStateId: "test/id",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceWorkspaceRoleAssignmentFQN,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s", "test", "00000000-0000-0000-0000-000000000000"),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceWorkspaceRoleAssignmentFQN,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s", "00000000-0000-0000-0000-000000000000", "test"),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestAcc_WorkspaceRoleAssignmentResource_CRUD(t *testing.T) {
	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, capacityID)

	entity := testhelp.WellKnown()["Principal"].(map[string]any)
	entityID := entity["id"].(string)
	entityType := entity["type"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceWorkspaceRoleAssignmentFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceWorkspaceRoleAssignmentFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceWorkspaceRoleAssignmentHeader,
					map[string]any{
						"workspace_id":   testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"principal_id":   entityID,
						"principal_type": entityType,
						"role":           "Member",
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceWorkspaceRoleAssignmentFQN, "principal_id", entityID),
				resource.TestCheckResourceAttr(testResourceWorkspaceRoleAssignmentFQN, "principal_type", entityType),
				resource.TestCheckResourceAttr(testResourceWorkspaceRoleAssignmentFQN, "role", "Member"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceWorkspaceRoleAssignmentFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceWorkspaceRoleAssignmentHeader,
					map[string]any{
						"workspace_id":   testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"principal_id":   entityID,
						"principal_type": entityType,
						"role":           "Viewer",
					},
				),
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceWorkspaceRoleAssignmentFQN, "principal_id", entityID),
				resource.TestCheckResourceAttr(testResourceWorkspaceRoleAssignmentFQN, "principal_type", entityType),
				resource.TestCheckResourceAttr(testResourceWorkspaceRoleAssignmentFQN, "role", "Viewer"),
			),
		},
	}))
}
