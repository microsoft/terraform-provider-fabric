// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testDataSourceGatewayRoleAssignments       = testhelp.DataSourceFQN("fabric", gatewayRoleAssignmentsTFName, "test")
	testDataSourceGatewayRoleAssignmentsHeader = at.DataSourceHeader(testhelp.TypeName("fabric", gatewayRoleAssignmentsTFName), "test")
)

func TestUnit_GatewayRoleAssignmentsDataSource(t *testing.T) {
	gatewayID := testhelp.RandomUUID()
	gatewayRoleAssignments := NewRandomGatewayRoleAssignments()
	fakes.FakeServer.ServerFactory.Core.GatewaysServer.NewListGatewayRoleAssignmentsPager = fakeGatewayRoleAssignments(gatewayRoleAssignments)

	entity := gatewayRoleAssignments.Value[1]

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceGatewayRoleAssignmentsHeader,
				map[string]any{
					"gateway_id":      gatewayID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceGatewayRoleAssignmentsHeader,
				map[string]any{
					"gateway_id": gatewayID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceGatewayRoleAssignments, "gateway_id", gatewayID),
				resource.TestCheckResourceAttrPtr(testDataSourceGatewayRoleAssignments, "values.1.id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceGatewayRoleAssignments, "values.1.role", (*string)(entity.Role)),
				// resource.TestCheckResourceAttrPtr(testDataSourceGatewayRoleAssignments, "values.1.display_name", entity.Principal.DisplayName),
				resource.TestCheckResourceAttrPtr(testDataSourceGatewayRoleAssignments, "values.1.type", (*string)(entity.Principal.Type)),
			),
		},
	}))
}

func TestAcc_GatewayRoleAssignmentsDataSource(t *testing.T) {
	gateway := testhelp.WellKnown()["GatewayVirtualNetwork"].(map[string]any)
	gatewayID := gateway["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceGatewayRoleAssignmentsHeader,
				map[string]any{
					"gateway_id": gatewayID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceGatewayRoleAssignments, "gateway_id", gatewayID),
				resource.TestCheckResourceAttrSet(testDataSourceGatewayRoleAssignments, "values.0.id"),
			),
		},
	},
	))
}
