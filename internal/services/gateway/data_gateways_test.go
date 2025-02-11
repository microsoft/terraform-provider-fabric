// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testDataSourceItemsFQN    = testhelp.DataSourceFQN("fabric", itemsTFName, "test")
	testDataSourceItemsHeader = at.DataSourceHeader(testhelp.TypeName("fabric", itemsTFName), "test")
)

func TestUnit_GatewaysDataSource(t *testing.T) {
	virtualNetworkGateway := fakes.NewRandomVirtualNetworkGateway()
	onPremisesGateway := fakes.NewRandomOnPremisesGateway()
	onPremisesGatewayPersonalGateway := fakes.NewRandomOnPremisesGatewayPersonal()

	fakes.FakeServer.Upsert(fakes.NewRandomGateway())
	fakes.FakeServer.Upsert(virtualNetworkGateway)
	fakes.FakeServer.Upsert(onPremisesGateway)
	fakes.FakeServer.Upsert(onPremisesGatewayPersonalGateway)
	fakes.FakeServer.Upsert(fakes.NewRandomGateway())

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - unexpected_attr
		// {
		// 	Config: at.CompileConfig(
		// 		testDataSourceItemsHeader,
		// 		map[string]any{
		// 			"unexpected_attr": "test",
		// 		},
		// 	),
		// 	ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		// },
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.2.id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.3.id"),
			),
		},
	}))
}

// func TestAcc_GatewaysDataSource(t *testing.T) {
// 	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
// 		// read
// 		{
// 			Config: at.CompileConfig(
// 				testDataSourceItemsHeader,
// 				map[string]any{},
// 			),
// 			Check: resource.ComposeAggregateTestCheckFunc(
// 				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.id"),
// 			),
// 		},
// 	},
// 	))
// }
