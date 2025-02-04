// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway_test

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
	testResourceGatewayRoleAssignment       = testhelp.ResourceFQN("fabric", "gateway_role_assignment", "test")
	testResourceGatewayRoleAssignmentHeader = at.ResourceHeader(testhelp.TypeName("fabric", "gateway_role_assignment"), "test")
)

func TestUnit_GatewayRoleAssignmentResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceGatewayRoleAssignment, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - missing required attribute: gateway_id
		{
			ResourceName: testResourceGatewayRoleAssignment,
			Config: at.CompileConfig(
				testResourceGatewayRoleAssignmentHeader,
				map[string]any{
					"principal_id":   "00000000-0000-0000-0000-000000000000",
					"principal_type": "User",
					"role":           "Member",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "gateway_id" is required, but no definition was found.`),
		},
		// error - missing required attribute: principal_id
		{
			ResourceName: testResourceGatewayRoleAssignment,
			Config: at.CompileConfig(
				testResourceGatewayRoleAssignmentHeader,
				map[string]any{
					"gateway_id":     "00000000-0000-0000-0000-000000000000",
					"principal_type": "User",
					"role":           "Member",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "principal_id" is required, but no definition was found.`),
		},
		// error - missing required attribute: principal_type
		{
			ResourceName: testResourceGatewayRoleAssignment,
			Config: at.CompileConfig(
				testResourceGatewayRoleAssignmentHeader,
				map[string]any{
					"gateway_id":   "00000000-0000-0000-0000-000000000000",
					"principal_id": "00000000-0000-0000-0000-000000000000",
					"role":         "Member",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "principal_type" is required, but no definition was found.`),
		},
		// error - missing required attribute: role
		{
			ResourceName: testResourceGatewayRoleAssignment,
			Config: at.CompileConfig(
				testResourceGatewayRoleAssignmentHeader,
				map[string]any{
					"gateway_id":     "00000000-0000-0000-0000-000000000000",
					"principal_id":   "00000000-0000-0000-0000-000000000000",
					"principal_type": "User",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "role" is required, but no definition was found.`),
		},
		// error - invalid UUID for gateway_id
		{
			ResourceName: testResourceGatewayRoleAssignment,
			Config: at.CompileConfig(
				testResourceGatewayRoleAssignmentHeader,
				map[string]any{
					"gateway_id":     "invalid uuid",
					"principal_id":   "00000000-0000-0000-0000-000000000000",
					"principal_type": "User",
					"role":           "Member",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - invalid UUID for principal_id
		{
			ResourceName: testResourceGatewayRoleAssignment,
			Config: at.CompileConfig(
				testResourceGatewayRoleAssignmentHeader,
				map[string]any{
					"gateway_id":     "00000000-0000-0000-0000-000000000000",
					"principal_id":   "invalid uuid",
					"principal_type": "User",
					"role":           "Member",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestUnit_GatewayRoleAssignmentResource_ImportState(t *testing.T) {
	testCase := at.CompileConfig(
		testResourceGatewayRoleAssignmentHeader,
		map[string]any{},
	)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceGatewayRoleAssignment, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		{
			ResourceName:  testResourceGatewayRoleAssignment,
			Config:        testCase,
			ImportStateId: "not-valid",
			ImportState:   true,
			ExpectError:   regexp.MustCompile("GatewayID/GatewayRoleAssignmentID"),
		},
		{
			ResourceName:  testResourceGatewayRoleAssignment,
			Config:        testCase,
			ImportStateId: "test/id",
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceGatewayRoleAssignment,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s", "test", "00000000-0000-0000-0000-000000000000"),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		{
			ResourceName:  testResourceGatewayRoleAssignment,
			Config:        testCase,
			ImportStateId: fmt.Sprintf("%s/%s", "00000000-0000-0000-0000-000000000000", "test"),
			ImportState:   true,
			ExpectError:   regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
	}))
}

func TestAcc_GatewayRoleAssignmentResource_CRUD(t *testing.T) {
	// Assume a well-known gateway is defined in the test environment.
	gateway := testhelp.WellKnown()["Gateway"].(map[string]any)
	gatewayID := gateway["id"].(string)

	// Assume a known principal is available.
	principal := testhelp.WellKnown()["Principal"].(map[string]any)
	principalID := principal["id"].(string)
	principalType := principal["type"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceGatewayRoleAssignment, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceGatewayRoleAssignment,
			Config: at.CompileConfig(
				testResourceGatewayRoleAssignmentHeader,
				map[string]any{
					"gateway_id":     gatewayID,
					"principal_id":   principalID,
					"principal_type": principalType,
					"role":           "Member",
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceGatewayRoleAssignment, "gateway_id", gatewayID),
				resource.TestCheckResourceAttr(testResourceGatewayRoleAssignment, "principal_id", principalID),
				resource.TestCheckResourceAttr(testResourceGatewayRoleAssignment, "principal_type", principalType),
				resource.TestCheckResourceAttr(testResourceGatewayRoleAssignment, "role", "Member"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceGatewayRoleAssignment,
			Config: at.CompileConfig(
				testResourceGatewayRoleAssignmentHeader,
				map[string]any{
					"gateway_id":     gatewayID,
					"principal_id":   principalID,
					"principal_type": principalType,
					"role":           "Viewer",
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceGatewayRoleAssignment, "gateway_id", gatewayID),
				resource.TestCheckResourceAttr(testResourceGatewayRoleAssignment, "principal_id", principalID),
				resource.TestCheckResourceAttr(testResourceGatewayRoleAssignment, "principal_type", principalType),
				resource.TestCheckResourceAttr(testResourceGatewayRoleAssignment, "role", "Viewer"),
			),
		},
	}))
}
