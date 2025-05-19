// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gatewayra_test

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

func TestUnit_GatewayRoleAssignmentDataSource(t *testing.T) {
	gatewayID := testhelp.RandomUUID()
	entity := NewRandomGatewayRoleAssignment()
	fakes.FakeServer.ServerFactory.Core.GatewaysServer.GetGatewayRoleAssignment = fakeGatewayRoleAssignment(entity)

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
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
				testDataSourceItemHeader,
				map[string]any{
					"id":         *entity.ID,
					"gateway_id": gatewayID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "gateway_id", gatewayID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "id", entity.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "role", (*string)(entity.Role)),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "principal.id", entity.Principal.ID),
				resource.TestCheckResourceAttrPtr(testDataSourceItemFQN, "principal.type", (*string)(entity.Principal.Type)),
			),
		},
	}))
}

func TestAcc_GatewayRoleAssignmentDataSource(t *testing.T) {
	// principal := testhelp.WellKnown()["Principal"].(map[string]any)
	// principalID := principal["id"].(string)
	// principalType := principal["type"].(string)

	group := testhelp.WellKnown()["Group"].(map[string]any)
	groupID := group["id"].(string)

	gw := testhelp.WellKnown()["GatewayVirtualNetwork"].(map[string]any)
	gwID := gw["id"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id":         groupID,
					"gateway_id": gwID,
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// {
		// 	Config: at.CompileConfig(
		// 		testDataSourceItemHeader,
		// 		map[string]any{
		// 			"id":         principalID,
		// 			"gateway_id": gwID,
		// 		},
		// 	),
		// 	Check: resource.ComposeAggregateTestCheckFunc(
		// 		resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "gateway_id"),
		// 		resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", principalID),
		// 		resource.TestCheckResourceAttr(testDataSourceItemFQN, "role", string(fabcore.GatewayRoleConnectionCreator)),
		// 		resource.TestCheckResourceAttr(testDataSourceItemFQN, "principal.id", principalID),
		// 		resource.TestCheckResourceAttr(testDataSourceItemFQN, "principal.type", principalType),
		// 	),
		// },
	},
	))
}
