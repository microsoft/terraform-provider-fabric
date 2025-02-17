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
	testDataSourceOnPremisesGatewaysFQN    = testhelp.DataSourceFQN("fabric", OnPremisesItemsTFName, "test")
	testDataSourceOnPremisesGatewaysHeader = at.DataSourceHeader(testhelp.TypeName("fabric", OnPremisesItemsTFName), "test")
)

func TestUnit_OnPremisesGatewaysDataSource(t *testing.T) {
	entity := fakes.NewRandomOnPremisesGateway()

	fakes.FakeServer.Upsert(fakes.NewRandomOnPremisesGateway())
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomOnPremisesGateway())

	resource.ParallelTest(t, testhelp.NewTestUnitCase(
		t,
		nil,
		fakes.FakeServer.ServerFactory,
		nil,
		[]resource.TestStep{
			// Step to ensure an unexpected attribute triggers an error.
			{
				Config: at.CompileConfig(
					testDataSourceOnPremisesGatewaysHeader,
					map[string]any{
						"unexpected_attr": "test",
					},
				),
				ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
			},
			// A normal read test with an empty config.
			{
				Config: at.CompileConfig(
					testDataSourceOnPremisesGatewaysHeader,
					map[string]any{},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSourceOnPremisesGatewaysFQN, "values.0.id"),
				),
			},
		},
	))
}
