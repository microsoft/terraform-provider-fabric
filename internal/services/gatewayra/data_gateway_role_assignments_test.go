// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gatewayra_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemsFQN, testDataSourceItemsHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Types, "test")

func TestUnit_GatewayRoleAssignmentsDataSource(t *testing.T) {
	gatewayID := testhelp.RandomUUID()
	gatewayRoleAssignments := NewRandomGatewayRoleAssignments()
	fakes.FakeServer.ServerFactory.Core.GatewaysServer.NewListGatewayRoleAssignmentsPager = fakeGatewayRoleAssignments(gatewayRoleAssignments)

	entity := gatewayRoleAssignments.Value[1]

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - unexpected_attr
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
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
				testDataSourceItemsHeader,
				map[string]any{
					"gateway_id": gatewayID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "gateway_id", gatewayID),
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testDataSourceItemsFQN,
					tfjsonpath.New("values"),
					knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"id":         knownvalue.StringExact(*entity.ID),
							"gateway_id": knownvalue.StringExact(gatewayID),
							"role":       knownvalue.StringExact((string)(*entity.Role)),
							"principal": knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"id":   knownvalue.StringExact(*entity.Principal.ID),
								"type": knownvalue.StringExact((string)(*entity.Principal.Type)),
							}),
						}),
					}),
				),
			},
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
				testDataSourceItemsHeader,
				map[string]any{
					"gateway_id": gatewayID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemsFQN, "gateway_id", gatewayID),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.id"),
			),
		},
	},
	))
}
