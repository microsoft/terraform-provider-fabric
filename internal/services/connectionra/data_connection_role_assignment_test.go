// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connectionra_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_ConnectionRoleAssignmentDataSource(t *testing.T) {
	connectionID := testhelp.RandomUUID()
	entity := NewRandomConnectionRoleAssignment()
	fakes.FakeServer.ServerFactory.Core.ConnectionsServer.GetConnectionRoleAssignment = fakeGetConnectionRoleAssignment(entity)

	resource.Test(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"connection_id":   connectionID,
					"id":              *entity.ID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - invalid UUID - id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"connection_id": connectionID,
					"id":            "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - invalid UUID - connection_id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"connection_id": "invalid uuid",
					"id":            *entity.ID,
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"connection_id": connectionID,
					"id":            *entity.ID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "connection_id", connectionID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "role", (*string)(entity.Role)),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "principal.id", entity.Principal.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "principal.type", (*string)(entity.Principal.Type)),
			),
		},
	}))
}

func TestAcc_ConnectionRoleAssignmentDataSource(t *testing.T) {
	virtualNetworkGatewayConnection := testhelp.WellKnown()["VirtualNetworkGatewayConnection"].(map[string]any)
	virtualNetworkGatewayConnectionID := virtualNetworkGatewayConnection["id"].(string)

	virtualNetworkGatewayConnectionRoleAssignment := testhelp.WellKnown()["VirtualNetworkGatewayConnectionRoleAssignment"].(map[string]any)
	virtualNetworkGatewayConnectionRoleAssignmentID := virtualNetworkGatewayConnectionRoleAssignment["id"].(string)
	principalType := virtualNetworkGatewayConnectionRoleAssignment["principalType"].(string)
	role := virtualNetworkGatewayConnectionRoleAssignment["role"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"connection_id": virtualNetworkGatewayConnectionID,
					"id":            virtualNetworkGatewayConnectionRoleAssignmentID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "connection_id"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", virtualNetworkGatewayConnectionRoleAssignmentID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "role", role),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "principal.id", virtualNetworkGatewayConnectionRoleAssignmentID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "principal.type", principalType),
			),
		},
	},
	))
}
