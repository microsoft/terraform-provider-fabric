// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package tags_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestAcc_TagResource_CRUD(t *testing.T) {
	entity1DisplayName := testhelp.RandomName()
	entity2DisplayName := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"tags": []map[string]any{
						{
							"display_name": entity1DisplayName,
						},
						{
							"display_name": entity2DisplayName,
						},
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "tags.0.display_name", entity1DisplayName),
			),
		},
	},
	))
}
