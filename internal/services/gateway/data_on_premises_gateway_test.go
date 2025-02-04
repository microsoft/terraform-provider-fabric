// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testDataSourceOnPremisesItemFabricFQN = testhelp.DataSourceFQN("fabric", OnPremisesItemTFName, "test")
	testDataSourceOnPremisesItemHeader    = at.DataSourceHeader(testhelp.TypeName("fabric", OnPremisesItemTFName), "test")
)

func TestUnit_OnPremisesGatewayDataSource(t *testing.T) {
	entity := fakes.NewRandomOnPremisesGateway()

	fakes.FakeServer.Upsert(fakes.NewRandomOnPremisesGateway())
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomOnPremisesGateway())

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceOnPremisesItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Exactly one of these attributes must be configured: \[id,display_name\]`),
		},
		// error - id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceOnPremisesItemHeader,
				map[string]any{
					"id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			Config: at.CompileConfig(
				testDataSourceOnPremisesItemHeader,
				map[string]any{
					"id":              *entity.ID,
					"unexpected_attr": "test",
				},
			),
			ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
		},
		// error - conflicting attributes
		{
			Config: at.CompileConfig(
				testDataSourceOnPremisesItemHeader,
				map[string]any{
					"id":           *entity.ID,
					"display_name": *entity.DisplayName,
				},
			),
			ExpectError: regexp.MustCompile(`These attributes cannot be configured together: \[id,display_name\]`),
		},
		// read by id - not found
		{
			Config: at.CompileConfig(
				testDataSourceOnPremisesItemHeader,
				map[string]any{
					"id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by id
		{
			Config: at.CompileConfig(
				testDataSourceOnPremisesItemHeader,
				map[string]any{
					"id": *entity.ID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceOnPremisesItemFabricFQN, "id", *entity.ID),
				resource.TestCheckResourceAttr(testDataSourceOnPremisesItemFabricFQN, "display_name", *entity.DisplayName),
				resource.TestCheckResourceAttrSet(testDataSourceOnPremisesItemFabricFQN, "allow_cloud_connection_refresh"),
				resource.TestCheckResourceAttrSet(testDataSourceOnPremisesItemFabricFQN, "allow_custom_connectors"),
				resource.TestCheckResourceAttrSet(testDataSourceOnPremisesItemFabricFQN, "load_balancing_setting"),
				resource.TestCheckResourceAttrSet(testDataSourceOnPremisesItemFabricFQN, "number_of_member_gateways"),
				resource.TestCheckResourceAttrSet(testDataSourceOnPremisesItemFabricFQN, "public_key.exponent"),
				resource.TestCheckResourceAttrSet(testDataSourceOnPremisesItemFabricFQN, "public_key.modulus"),
				resource.TestCheckResourceAttrSet(testDataSourceOnPremisesItemFabricFQN, "version"),
			),
		},
		// read by name - not found
		{
			Config: at.CompileConfig(
				testDataSourceOnPremisesItemHeader,
				map[string]any{
					"display_name": testhelp.RandomName(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by name
		{
			Config: at.CompileConfig(
				testDataSourceOnPremisesItemHeader,
				map[string]any{
					"display_name": *entity.DisplayName,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceOnPremisesItemFabricFQN, "id", *entity.ID),
				resource.TestCheckResourceAttr(testDataSourceOnPremisesItemFabricFQN, "display_name", *entity.DisplayName),
				resource.TestCheckResourceAttrSet(testDataSourceOnPremisesItemFabricFQN, "allow_cloud_connection_refresh"),
				resource.TestCheckResourceAttrSet(testDataSourceOnPremisesItemFabricFQN, "allow_custom_connectors"),
				resource.TestCheckResourceAttrSet(testDataSourceOnPremisesItemFabricFQN, "load_balancing_setting"),
				resource.TestCheckResourceAttrSet(testDataSourceOnPremisesItemFabricFQN, "number_of_member_gateways"),
				resource.TestCheckResourceAttrSet(testDataSourceOnPremisesItemFabricFQN, "public_key.exponent"),
				resource.TestCheckResourceAttrSet(testDataSourceOnPremisesItemFabricFQN, "public_key.modulus"),
				resource.TestCheckResourceAttrSet(testDataSourceOnPremisesItemFabricFQN, "version"),
			),
		},
	}))
}

func TestAcc_OnPremisesGatewayDataSource(t *testing.T) {
	entity := testhelp.WellKnown()["OnPremisesGateway"].(map[string]any)
	entityID := entity["id"].(string)
	entityDisplayName := entity["displayName"].(string)
	entityDescription := entity["description"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(
				testDataSourceOnPremisesItemHeader,
				map[string]any{
					"id": entityID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceOnPremisesItemFabricFQN, "id", entityID),
				resource.TestCheckResourceAttr(testDataSourceOnPremisesItemFabricFQN, "display_name", entityDisplayName),
				resource.TestCheckResourceAttr(testDataSourceOnPremisesItemFabricFQN, "description", entityDescription),
			),
		},
		// read by id - not found
		{
			Config: at.CompileConfig(
				testDataSourceOnPremisesItemHeader,
				map[string]any{
					"id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}
