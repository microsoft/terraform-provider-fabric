// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package spark_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var (
	testResourceSparkCustomPoolFQN    = testhelp.ResourceFQN("fabric", sparkCustomPoolTFName, "test")
	testResourceSparkCustomPoolHeader = at.ResourceHeader(testhelp.TypeName("fabric", sparkCustomPoolTFName), "test")
)

func TestAcc_SparkCustomPoolResource_CRUD(t *testing.T) {
	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, *testhelp.WellKnown().Capacity.ID)
	testHelperSparkCustomPoolResource := getSparkCustomPoolResourceAttr(t, testhelp.RefByFQN(workspaceResourceFQN, "id"), "test")

	entityCreateName := testhelp.RandomName()
	testCaseCreate := testhelp.CopyMap(testHelperSparkCustomPoolResource)
	testCaseCreate["name"] = entityCreateName

	entityUpdateName := testhelp.RandomName()
	testCaseUpdate := testhelp.CopyMap(testHelperSparkCustomPoolResource)
	testCaseUpdate["name"] = entityUpdateName

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceSparkCustomPoolFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceSparkCustomPoolFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceSparkCustomPoolHeader,
					testCaseCreate,
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceSparkCustomPoolFQN, "name", entityCreateName),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceSparkCustomPoolFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceSparkCustomPoolHeader,
					testCaseUpdate,
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceSparkCustomPoolFQN, "name", entityUpdateName),
			),
		},
	},
	))
}
