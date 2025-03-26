// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspacempe_test

import (
	"fmt"
	"strings"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestAcc_WorkspaceManagedPrivateEndpointResource_CRUD(t *testing.T) {
	azure := testhelp.WellKnown()["Azure"].(map[string]any)
	azureSubscriptionID := azure["subscriptionId"].(string)
	azureLocation := azure["location"].(string)

	capacity := testhelp.WellKnown()["Capacity"].(map[string]any)
	capacityID := capacity["id"].(string)

	workspaceResourceHCL, workspaceResourceFQN := testhelp.TestAccWorkspaceResource(t, capacityID)

	rgName := testhelp.RandomName()
	stName := strings.ToLower(rgName)

	t.Setenv("ARM_SUBSCRIPTION_ID", azureSubscriptionID)

	// export ARM_CLIENT_ID="00000000-0000-0000-0000-000000000000"
	// export ARM_SUBSCRIPTION_ID="00000000-0000-0000-0000-000000000000"
	// export ARM_TENANT_ID="00000000-0000-0000-0000-000000000000"

	// provider "azurerm" {
	// 	features {}
	// 	subscription_id = "%s"
	// }

	azurermHCL := fmt.Sprintf(`
		provider "azurerm" {
			features {}
		}

		resource "azurerm_resource_group" "test" {
			name     = "%s"
			location = "%s"
			tags = {
				environment = "testacc"
				solution = "FabricTF"
			}
		}

		resource "azurerm_storage_account" "test" {
			name                     = "%s"
			resource_group_name      = azurerm_resource_group.test.name
			location                 = azurerm_resource_group.test.location
			account_replication_type = "LRS"
			account_tier             = "Standard"
			account_kind             = "StorageV2"
			tags = {
				environment = "testacc"
				solution = "FabricTF"
			}
		}`, rgName, azureLocation, stName,
	)

	entityName := testhelp.RandomName()

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		{
			// ExternalProviders: map[string]resource.ExternalProvider{
			// 	"azurerm": {
			// 		Source: "hashicorp/azurerm",
			// 	},
			// },
			ResourceName: testResourceItemFQN,
			Config: at.JoinConfigs(
				azurermHCL,
				workspaceResourceHCL,
				at.CompileConfig(
					testResourceItemHeader,
					map[string]any{
						"workspace_id":                    testhelp.RefByFQN(workspaceResourceFQN, "id"),
						"name":                            entityName,
						"target_private_link_resource_id": testhelp.RefByFQN("azurerm_storage_account.test", "id"),
						"target_subresource_type":         "blob",
						"request_message":                 rgName + "/" + stName,
					},
				),
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testResourceItemFQN,
					tfjsonpath.New("name"),
					knownvalue.StringExact(entityName),
				),
			},
		},
	}))
}
