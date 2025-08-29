package externaldatasharesprovider_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, "external_data_shares_provider", "test")

func TestUnit_ExternalDataSharesResource(t *testing.T) {
	fakes.FakeServer.ServerFactory.Admin.ExternalDataSharesProviderServer.RevokeExternalDataShare = fakeRevokeExternalDataSharesProvider()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Revoke
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":           testhelp.RandomUUID(),
					"item_id":                testhelp.RandomUUID(),
					"external_data_share_id": testhelp.RandomUUID(),
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(),
		},
	}))
}

func TestAcc_ExternalDataSharesResource(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":           testhelp.RandomUUID(),
					"item_id":                testhelp.RandomUUID(),
					"external_data_share_id": testhelp.RandomUUID(),
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(),
		},
	},
	))
}
