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
	testDataSourceVirtualNetworkGatewaysFQN    = testhelp.DataSourceFQN("fabric", "virtual_network_gateways", "test")
	testDataSourceVirtualNetworkGatewaysHeader = at.DataSourceHeader(testhelp.TypeName("fabric", "virtual_network_gateways"), "test")
)

func TestUnit_VirtualNetworkGatewaysDataSource(t *testing.T) {
	entity := fakes.NewRandomVirtualNetworkGateway()

	fakes.FakeServer.Upsert(fakes.NewRandomVirtualNetworkGateway())
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomVirtualNetworkGateway())

	resource.ParallelTest(t, testhelp.NewTestUnitCase(
		t,
		nil,
		fakes.FakeServer.ServerFactory,
		nil,
		[]resource.TestStep{
			// Check that using an unexpected attribute fails.
			{
				Config: at.CompileConfig(
					testDataSourceVirtualNetworkGatewaysHeader,
					map[string]any{
						"unexpected_attr": "test",
					},
				),
				ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
			},
			// A normal read test with an empty config (no filter).
			{
				Config: at.CompileConfig(
					testDataSourceVirtualNetworkGatewaysHeader,
					map[string]any{},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSourceVirtualNetworkGatewaysFQN, "values.0.id"),
				),
			},
		},
	))
}

func TestAcc_VirtualNetworkGatewaysDataSource(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestAccCase(
		t,
		nil,
		nil,
		[]resource.TestStep{
			// read test.
			{
				Config: at.CompileConfig(
					testDataSourceVirtualNetworkGatewaysHeader,
					map[string]any{},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSourceVirtualNetworkGatewaysFQN, "values.0.id"),
				),
			},
		},
	))
}
