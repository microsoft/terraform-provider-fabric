// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection_test

import (
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

func TestAcc_ConnectionResource_ShareableCloud(t *testing.T) {
	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      testhelp.RandomName(),
					"connectivity_type": "ShareableCloud",
					"privacy_level":     "Organizational",
					"connection_details": map[string]any{
						"type":            "FTP",
						"creation_method": "FTP.Contents",
						"parameters": []map[string]any{
							{
								"name":  "server",
								"value": "ftp.example.com",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": "NotEncrypted",
						"single_sign_on_type":   "None",
						"skip_test_connection":  false,
						"credential_type":       "Basic",
						"basic_credentials": map[string]any{
							"username":            "TODO",
							"password_wo":         "TODO",
							"password_wo_version": 1,
						},
					},
				},
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testResourceItemFQN,
					tfjsonpath.New("id"),
					knownvalue.NotNull(),
				),
			},
		},
	},
	))
}

func TestAcc_ConnectionResource_VirtualNetworkGateway(t *testing.T) {
	entityVirtualNetwork := testhelp.WellKnown()["GatewayVirtualNetwork"].(map[string]any)
	entityVirtualNetworkID := entityVirtualNetwork["id"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      testhelp.RandomName(),
					"connectivity_type": "VirtualNetworkGateway",
					"privacy_level":     "Organizational",
					"gateway_id":        entityVirtualNetworkID,
					"connection_details": map[string]any{
						"type":            "AzureBlobs",
						"creation_method": "AzureBlobs",
						"parameters": []map[string]any{
							{
								"name":  "account",
								"value": "TODO",
							},
							{
								"name":  "domain",
								"value": "blob.core.windows.net",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": "NotEncrypted",
						"single_sign_on_type":   "None",
						"skip_test_connection":  false,
						"credential_type":       "Key",
						"key_credentials": map[string]any{
							"key_wo":         "TODO",
							"key_wo_version": 1,
						},
					},
				},
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testResourceItemFQN,
					tfjsonpath.New("id"),
					knownvalue.NotNull(),
				),
			},
		},
	},
	))
}

func TestAcc_ConnectionResource_OnPremisesGateway(t *testing.T) {
	entityVirtualNetwork := testhelp.WellKnown()["GatewayVirtualNetwork"].(map[string]any)
	entityVirtualNetworkID := entityVirtualNetwork["id"].(string)

	resource.Test(t, testhelp.NewTestAccCase(t, &testResourceItemFQN, nil, []resource.TestStep{
		// Create and Read
		{
			ResourceName: testResourceItemFQN,
			Config: at.CompileConfig(
				testResourceItemHeader,
				map[string]any{
					"display_name":      testhelp.RandomName(),
					"connectivity_type": "OnPremisesGateway",
					"privacy_level":     "Organizational",
					"gateway_id":        entityVirtualNetworkID,
					"connection_details": map[string]any{
						"type":            "AzureBlobs",
						"creation_method": "AzureBlobs",
						"parameters": []map[string]any{
							{
								"name":  "account",
								"value": "TODO",
							},
							{
								"name":  "domain",
								"value": "blob.core.windows.net",
							},
						},
					},
					"credential_details": map[string]any{
						"connection_encryption": "NotEncrypted",
						"single_sign_on_type":   "None",
						"skip_test_connection":  false,
						"on_premises_gateway_credentials": []map[string]any{
							{
								"gateway_id":      entityVirtualNetworkID,
								"credential_type": "Windows",
								"windows_credentials": map[string]any{
									"username":            "TODO",
									"password_wo":         "TODO",
									"password_wo_version": 1,
								},
								"public_key": map[string]any{
									"exponent": "TODO",
									"modulus":  "TODO",
								},
							},
						},
					},
				},
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testResourceItemFQN,
					tfjsonpath.New("id"),
					knownvalue.NotNull(),
				),
			},
		},
	},
	))
}
