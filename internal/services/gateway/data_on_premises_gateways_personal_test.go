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
	testDataSourceOnPremisesPersonalsFQN    = testhelp.DataSourceFQN("fabric", OnPremisesPersonalItemsTFName, "test")
	testDataSourceOnPremisesPersonalsHeader = at.DataSourceHeader(testhelp.TypeName("fabric", OnPremisesPersonalItemsTFName), "test")
)

func TestUnit_OnPremisesGatewaysPersonalDataSource(t *testing.T) {
	entity := fakes.NewRandomOnPremisesGatewayPersonal()

	fakes.FakeServer.Upsert(fakes.NewRandomOnPremisesGatewayPersonal())
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomOnPremisesGatewayPersonal())

	resource.ParallelTest(t, testhelp.NewTestUnitCase(
		t,
		nil,
		fakes.FakeServer.ServerFactory,
		nil,
		[]resource.TestStep{
			// Step 1: Use an unexpected attribute to trigger an error.
			{
				Config: at.CompileConfig(
					testDataSourceOnPremisesPersonalsHeader,
					map[string]any{
						"unexpected_attr": "test",
					},
				),
				ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
			},
			// Step 2: Normal read test with empty configuration.
			{
				Config: at.CompileConfig(
					testDataSourceOnPremisesPersonalsHeader,
					map[string]any{},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSourceOnPremisesPersonalsFQN, "values.0.id"),
				),
			},
		},
	))
}
