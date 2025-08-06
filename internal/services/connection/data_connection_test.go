// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")

func TestUnit_ConnectionDataSource(t *testing.T) {
	entity := fakes.NewRandomShareableCloudConnection()

	fakes.FakeServer.Upsert(fakes.NewRandomConnection())
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomConnection())

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found`),
		},
		// error - id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id":              *entity.ID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// read by id - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by id
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": *entity.ID,
				},
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					testDataSourceItemFQN,
					tfjsonpath.New("id"),
					knownvalue.StringExact(*entity.ID),
				),
				statecheck.ExpectKnownValue(
					testDataSourceItemFQN,
					tfjsonpath.New("display_name"),
					knownvalue.StringExact(*entity.DisplayName),
				),
				statecheck.ExpectKnownValue(
					testDataSourceItemFQN,
					tfjsonpath.New("connectivity_type"),
					knownvalue.StringExact(string(*entity.ConnectivityType)),
				),
				statecheck.ExpectKnownValue(
					testDataSourceItemFQN,
					tfjsonpath.New("privacy_level"),
					knownvalue.StringExact(string(*entity.PrivacyLevel)),
				),
				statecheck.ExpectKnownValue(
					testDataSourceItemFQN,
					tfjsonpath.New("connection_details"),
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"path": knownvalue.StringExact(*entity.ConnectionDetails.Path),
						"type": knownvalue.StringExact(*entity.ConnectionDetails.Type),
					}),
				),
				statecheck.ExpectKnownValue(
					testDataSourceItemFQN,
					tfjsonpath.New("credential_details"),
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"connection_encryption": knownvalue.StringExact(string(*entity.CredentialDetails.ConnectionEncryption)),
						"credential_type":       knownvalue.StringExact(string(*entity.CredentialDetails.CredentialType)),
						"single_sign_on_type":   knownvalue.StringExact(string(*entity.CredentialDetails.SingleSignOnType)),
						"skip_test_connection":  knownvalue.Bool(*entity.CredentialDetails.SkipTestConnection),
					}),
				),
			},
		},
	}))
}

func TestAcc_ConnectionDataSource(t *testing.T) {
	shareableCloudConnection := testhelp.WellKnown()["ShareableCloudConnection"].(map[string]any)
	shareableCloudConnectionID := shareableCloudConnection["id"].(string)
	shareableCloudConnectionDisplayName := shareableCloudConnection["displayName"].(string)
	virtualNetworkGatewayConnection := testhelp.WellKnown()["VirtualNetworkGatewayConnection"].(map[string]any)
	virtualNetworkGatewayConnectionID := virtualNetworkGatewayConnection["id"].(string)
	virtualNetworkGatewayConnectionDisplayName := virtualNetworkGatewayConnection["displayName"].(string)
	virtualNetworkGatewayConnectionGatewayID := virtualNetworkGatewayConnection["gatewayId"].(string)
	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		// read by id - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by id - ShareableCloudConnection
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": shareableCloudConnectionID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", shareableCloudConnectionID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "display_name", shareableCloudConnectionDisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "connectivity_type", "ShareableCloud"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "privacy_level"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "connection_details.path"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "connection_details.type"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "credential_details.connection_encryption"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "credential_details.credential_type"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "credential_details.single_sign_on_type"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "credential_details.skip_test_connection"),
			),
		},
		// read by id - VirtualNetworkGatewayConnection
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"id": virtualNetworkGatewayConnectionID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "id", virtualNetworkGatewayConnectionID),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "display_name", virtualNetworkGatewayConnectionDisplayName),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "connectivity_type", "VirtualNetworkGateway"),
				resource.TestCheckResourceAttr(testDataSourceItemFQN, "gateway_id", virtualNetworkGatewayConnectionGatewayID),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "privacy_level"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "connection_details.path"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "connection_details.type"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "credential_details.connection_encryption"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "credential_details.credential_type"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "credential_details.single_sign_on_type"),
				resource.TestCheckResourceAttrSet(testDataSourceItemFQN, "credential_details.skip_test_connection"),
			),
		},
	}))
}
