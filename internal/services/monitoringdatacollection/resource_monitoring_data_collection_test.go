// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package monitoringdatacollection_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"
	fabmonitoring "github.com/microsoft/fabric-sdk-go/fabric/monitoring"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

const resourceFQN = "fabric_monitoring_data_collection.test"

func collectionConfig(workspaceID, monitoringID string, enabled bool) string {
	return fmt.Sprintf(`
resource "fabric_monitoring_data_collection" "test" {
  workspace_id  = %q
  monitoring_id = %q
  enabled       = %t
}
`, workspaceID, monitoringID, enabled)
}

func TestUnit_MonitoringDataCollectionResource_CRUD(t *testing.T) {
	server := fakes.NewFakeServer()
	entity := fakes.NewRandomMonitoring()
	server.Upsert(entity)

	config := func(enabled bool) string {
		return collectionConfig(*entity.WorkspaceID, *entity.ID, enabled)
	}
	testCase := testhelp.NewTestUnitCase(t, nil, server.ServerFactory, nil, []resource.TestStep{
		{
			Config: config(true),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceFQN, "enabled", "true"),
				resource.TestCheckResourceAttr(resourceFQN, "id", *entity.ID),
			),
		},
		{
			Config:   config(true),
			PlanOnly: true,
		},
		{
			Config: config(false),
			Check:  resource.TestCheckResourceAttr(resourceFQN, "enabled", "false"),
		},
		{
			ResourceName:      resourceFQN,
			ImportState:       true,
			ImportStateId:     *entity.WorkspaceID + "/" + *entity.ID,
			ImportStateVerify: true,
		},
		{
			Config: config(true),
			Check:  resource.TestCheckResourceAttr(resourceFQN, "enabled", "true"),
		},
	})
	testCase.CheckDestroy = func(_ *terraform.State) error {
		options := &fabric.ClientOptions{}
		options.Transport = fabfake.NewServerFactoryTransport(server.ServerFactory)

		factory, err := fabmonitoring.NewClientFactory(&azfake.TokenCredential{}, nil, options)
		if err != nil {
			return err
		}

		response, err := factory.NewItemsClient().GetMonitoring(t.Context(), *entity.WorkspaceID, *entity.ID, nil)
		if err != nil {
			return fmt.Errorf("destroy removed the Monitoring item: %w", err)
		}

		if response.Properties == nil || response.Properties.DataCollectionState == nil ||
			*response.Properties.DataCollectionState != fabmonitoring.DataCollectionStateDisabled {
			return errors.New("destroy did not disable Monitoring collection")
		}

		return nil
	}
	resource.ParallelTest(t, testCase)
}

func TestUnit_MonitoringDataCollectionResource_Attributes(t *testing.T) {
	server := fakes.NewFakeServer()
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, server.ServerFactory, nil, []resource.TestStep{
		{
			Config:      `resource "fabric_monitoring_data_collection" "test" {}`,
			ExpectError: regexp.MustCompile("Missing required argument"),
		},
		{
			Config:      collectionConfig("invalid", testhelp.RandomUUID(), true),
			ExpectError: regexp.MustCompile("UUID"),
		},
		{
			Config:      collectionConfig(testhelp.RandomUUID(), "invalid", true),
			ExpectError: regexp.MustCompile("UUID"),
		},
		{
			Config:      collectionConfig(testhelp.RandomUUID(), testhelp.RandomUUID(), true),
			ExpectError: regexp.MustCompile(common.ErrorCreateHeader),
		},
	}))
}

func TestUnit_MonitoringDataCollectionResource_ImportState(t *testing.T) {
	for _, id := range []string{"invalid", "invalid/invalid", testhelp.RandomUUID() + "/invalid", testhelp.RandomUUID() + "/" + testhelp.RandomUUID() + "/extra"} {
		t.Run(id, func(t *testing.T) {
			server := fakes.NewFakeServer()
			resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, server.ServerFactory, nil, []resource.TestStep{
				{
					Config:        collectionConfig(testhelp.RandomUUID(), testhelp.RandomUUID(), true),
					ResourceName:  resourceFQN,
					ImportState:   true,
					ImportStateId: id,
					ExpectError:   regexp.MustCompile("(?i)(import|UUID)"),
				},
			}))
		})
	}
}

func TestUnit_MonitoringDataCollectionResource_UpdateError(t *testing.T) {
	server := fakes.NewFakeServer()
	entity := fakes.NewRandomMonitoring()
	server.Upsert(entity)
	server.ServerFactory.Monitoring.DataCollectionServer.Enable = func(context.Context, string, string, *fabmonitoring.DataCollectionClientEnableOptions) (azfake.Responder[fabmonitoring.DataCollectionClientEnableResponse], azfake.ErrorResponder) {
		var response azfake.Responder[fabmonitoring.DataCollectionClientEnableResponse]
		var errResponse azfake.ErrorResponder
		errResponse.SetError(fabfake.SetResponseError(http.StatusForbidden, "InsufficientPrivileges", "Workspace Admin role required"))

		return response, errResponse
	}
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, server.ServerFactory, nil, []resource.TestStep{
		{
			Config: collectionConfig(*entity.WorkspaceID, *entity.ID, false),
			Check:  resource.TestCheckResourceAttr(resourceFQN, "enabled", "false"),
		},
		{
			Config:      collectionConfig(*entity.WorkspaceID, *entity.ID, true),
			ExpectError: regexp.MustCompile(common.ErrorUpdateHeader),
		},
	}))
}

func monitoringCollectionConfig(workspaceID, name string, enabled bool) string {
	return fmt.Sprintf(`
resource "fabric_monitoring" "test" {
  workspace_id = %q
  display_name = %q
  format       = "Default"
  definition = {
    "MonitoringProperties.json" = {
      source = %q
    }
  }
}
resource "fabric_monitoring_data_collection" "test" {
  workspace_id  = fabric_monitoring.test.workspace_id
  monitoring_id = fabric_monitoring.test.id
  enabled       = %t
}
`, workspaceID, name, testhelp.GetFixturesDirPath("monitoring")+"/MonitoringProperties.auto.json.tmpl", enabled)
}

func TestUnit_MonitoringDataCollectionResource_WithMonitoring(t *testing.T) {
	server := fakes.NewFakeServer()
	workspaceID := testhelp.RandomUUID()
	name := testhelp.RandomName()
	renamed := testhelp.RandomName()
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, server.ServerFactory, nil, []resource.TestStep{
		{
			Config: monitoringCollectionConfig(workspaceID, name, true),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("fabric_monitoring.test", "format", "Default"),
				resource.TestCheckResourceAttr(resourceFQN, "enabled", "true"),
			),
		},
		{
			Config: monitoringCollectionConfig(workspaceID, renamed, true),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("fabric_monitoring.test", "display_name", renamed),
				resource.TestCheckResourceAttr("fabric_monitoring.test", "properties.data_collection_state", "Enabled"),
				resource.TestCheckResourceAttr(resourceFQN, "enabled", "true"),
			),
		},
		{
			Config: monitoringCollectionConfig(workspaceID, renamed, false),
			Check:  resource.TestCheckResourceAttr(resourceFQN, "enabled", "false"),
		},
		{
			Config: monitoringCollectionConfig(workspaceID, name, false),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("fabric_monitoring.test", "display_name", name),
				resource.TestCheckResourceAttr("fabric_monitoring.test", "properties.data_collection_state", "Disabled"),
				resource.TestCheckResourceAttr(resourceFQN, "enabled", "false"),
			),
		},
		{
			Config:   monitoringCollectionConfig(workspaceID, name, false),
			PlanOnly: true,
		},
	}))
}

func TestAcc_MonitoringDataCollectionResource_CRUD(t *testing.T) {
	workspaceID := testhelp.WellKnown()["WorkspaceRS"].(map[string]any)["id"].(string)
	name := testhelp.RandomName()
	config := func(enabled bool) string {
		return monitoringCollectionConfig(workspaceID, name, enabled)
	}
	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		{
			Config: config(true),
			Check:  resource.TestCheckResourceAttr(resourceFQN, "enabled", "true"),
		},
		{
			Config: config(false),
			Check:  resource.TestCheckResourceAttr(resourceFQN, "enabled", "false"),
		},
	}))
}
