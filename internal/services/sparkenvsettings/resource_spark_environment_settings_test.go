// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sparkenvsettings_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_SparkEnvSettingsResource_Attributes(t *testing.T) {
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Missing required argument`),
		},
		// error - workspace_id - invalid UUID
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":       "invalid uuid",
					"environment_id":     testhelp.RandomUUID(),
					"publication_status": "Published",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - environment_id - invalid UUID
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":       testhelp.RandomUUID(),
					"environment_id":     "invalid uuid",
					"publication_status": "Published",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":       testhelp.RandomUUID(),
					"environment_id":     testhelp.RandomUUID(),
					"publication_status": "Published",
					"unexpected_attr":    "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - missing workspace_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"environment_id":     testhelp.RandomUUID(),
					"publication_status": "Published",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "workspace_id" is required, but no definition was found.`),
		},
		// error - missing environment_id
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":       testhelp.RandomUUID(),
					"publication_status": "Published",
				},
			),
			ExpectError: regexp.MustCompile(`The argument "environment_id" is required, but no definition was found.`),
		},
		// error - missing publication_status
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":   testhelp.RandomUUID(),
					"environment_id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(`The argument "publication_status" is required, but no definition was found.`),
		},
		// error - invalid key format in spark_properties
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":       testhelp.RandomUUID(),
					"environment_id":     testhelp.RandomUUID(),
					"publication_status": "Published",
					"spark_properties": map[string]any{
						`"spark test"`: "12",
					},
				},
			),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
		},
		// error - invalid nil value for key
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":       testhelp.RandomUUID(),
					"environment_id":     testhelp.RandomUUID(),
					"publication_status": "Published",
					"spark_properties": map[string]any{
						`"spark.acls.enable"`: nil,
					},
				},
			),
			ExpectError: regexp.MustCompile(`Error: Null Map Value`),
		},
	}))
}

func TestUnit_SparkEnvironmentSettingsResource_CRUD(t *testing.T) {
	workspaceID := testhelp.RandomUUID()
	environmentID := testhelp.RandomUUID()

	stagingSparkCompute := NewRandomSparkCompute()
	stagingSparkCompute.SparkProperties = []fabenvironment.SparkProperty{
		{
			Key:   new("spark.acls.enable"),
			Value: new("false"),
		},
		{
			Key:   new("spark.acls.groups"),
			Value: new("test"),
		},
	}

	fakeTestUpsertSparkComputeStaging(environmentID, stagingSparkCompute)

	fakes.FakeServer.ServerFactory.Environment.StagingServer.GetSparkCompute = fakeGetStagingSparkComputeFunc()
	fakes.FakeServer.ServerFactory.Environment.StagingServer.UpdateSparkCompute = fakeUpdateStagingSparkComputeFunc()
	fakes.FakeServer.ServerFactory.Environment.PublishedServer.GetSparkCompute = fakeGetPublishedSparkComputeFunc()
	fakes.FakeServer.ServerFactory.Environment.ItemsServer.BeginPublishEnvironment = fakeBeginPublishEnvironmentFunc()

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, &testResourceItemFQN, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":       workspaceID,
					"environment_id":     environmentID,
					"publication_status": "Published",
					"driver_cores":       4,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "id", environmentID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "driver_cores", "4"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":       workspaceID,
					"environment_id":     environmentID,
					"publication_status": "Published",
					"driver_cores":       8,
					"spark_properties": map[string]any{
						`"spark.acls.enable"`: "true",
					},
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "id", environmentID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "driver_cores", "8"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "spark_properties.spark.acls.enable", "true"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "spark_properties.%", "1"),
			),
		},
		// Update and Read - remove spark properties
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"workspace_id":       workspaceID,
					"environment_id":     environmentID,
					"publication_status": "Published",
					"driver_cores":       8,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "id", environmentID),
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "driver_cores", "8"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "spark_properties"),
			),
		},
	}))
}

func TestAcc_SparkEnvironmentSettingsResource_CRUD(t *testing.T) {
	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, capacityID)
	environmentResourceHCL, environmentResourceFQN := environmentResource(t, testhelp.RefByFQN(workspaceResourceFQN, "id"))

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				environmentResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":       testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"environment_id":     testhelp.RefByFQN(environmentResourceFQN, "id"),
						"publication_status": "Published",
						"driver_cores":       4,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "driver_cores", "4"),
			),
		},
		// Update and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				environmentResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":       testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"environment_id":     testhelp.RefByFQN(environmentResourceFQN, "id"),
						"publication_status": "Published",
						"driver_cores":       8,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "driver_cores", "8"),
			),
		},
	},
	))
}

func TestAcc_SparkEnvironmentSettingsSparkPropertiesResource_CRUD(t *testing.T) {
	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, capacityID)
	environmentResourceHCL, environmentResourceFQN := environmentResource(t, testhelp.RefByFQN(workspaceResourceFQN, "id"))

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read (Spark properties)
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				environmentResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":       testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"environment_id":     testhelp.RefByFQN(environmentResourceFQN, "id"),
						"publication_status": "Staging",
						"driver_cores":       8,
						"spark_properties": map[string]any{
							`"spark.acls.enable"`:       "true",
							`"spark.admin.acls.groups"`: "test",
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "driver_cores", "8"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "spark_properties.spark.acls.enable", "true"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "spark_properties.spark.admin.acls.groups", "test"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "spark_properties.%", "2"),
			),
		},
		// Update and Read (test Spark properties sync)
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				environmentResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":       testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"environment_id":     testhelp.RefByFQN(environmentResourceFQN, "id"),
						"publication_status": "Staging",
						"driver_cores":       8,
						"spark_properties": map[string]any{
							`"spark.cores.max"`: "12",
						},
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "driver_cores", "8"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "spark_properties.spark.cores.max", "12"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "spark_properties.%", "1"),
			),
		},
		// Create and Read (remove Spark properties)
		{
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				workspaceResourceHCL,
				environmentResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":       testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"environment_id":     testhelp.RefByFQN(environmentResourceFQN, "id"),
						"publication_status": "Published",
						"driver_cores":       8,
					},
				)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testResourceItemFQN, "pool.name", "Starter Pool"),
				resource.TestCheckResourceAttr(testResourceItemFQN, "driver_cores", "8"),
				resource.TestCheckNoResourceAttr(testResourceItemFQN, "spark_properties"),
			),
		},
	},
	))
}
