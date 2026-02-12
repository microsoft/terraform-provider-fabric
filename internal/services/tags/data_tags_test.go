// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package tags_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemsFQN, testDataSourceItemsHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Types, "test")

func TestUnit_TagsDataSource(t *testing.T) {
	fakeTestUpsert(NewRandomTag())
	fakeTestUpsert(NewRandomTag())

	fakes.FakeServer.ServerFactory.Admin.TagsServer.NewListTagsPager = fakeTagsFunc()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.display_name"),
			),
		},
	}))
}

func TestAcc_TagsDataSource(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read
		{
			Config: at.CompileConfig(
				testDataSourceItemsHeader,
				map[string]any{},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.id"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.display_name"),
				resource.TestCheckResourceAttrSet(testDataSourceItemsFQN, "values.0.scope.type"),
			),
		},
	},
	))
}
