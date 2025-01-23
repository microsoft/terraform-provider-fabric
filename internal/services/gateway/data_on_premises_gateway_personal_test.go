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
	testDataSourceOnPremisesPresonalItemFabricFQN = testhelp.DataSourceFQN("fabric", OnPremisesItemTFName, "test")
	testDataSourceOnPremisesPersonalItemHeader    = at.DataSourceHeader(testhelp.TypeName("fabric", OnPremisesItemTFName), "test")
)

func TestUnit_OnPremisesGatewayPersonalDataSource(t *testing.T) {
	entity := fakes.NewRandomOnPermisesGatewayPersonal()

	fakes.FakeServer.Upsert(fakes.NewRandomOnPermisesGatewayPersonal())
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomOnPermisesGatewayPersonal())

	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, fakes.FakeServer.ServerFactory, nil, []resource.TestStep{
		// error - no attributes
		{
			Config: at.CompileConfig(
				testDataSourceOnPremisesPersonalItemHeader,
				map[string]any{},
			),
			ExpectError: regexp.MustCompile(`Exactly one of these attributes must be configured: \[id,display_name\]`),
		},
		// error - id - invalid UUID
		{
			Config: at.CompileConfig(
				testDataSourceOnPremisesPersonalItemHeader,
				map[string]any{
					"id": "invalid uuid",
				},
			),
			ExpectError: regexp.MustCompile(customtypes.UUIDTypeErrorInvalidStringHeader),
		},
		// error - unexpected attribute
		{
			Config: at.CompileConfig(
				testDataSourceOnPremisesPersonalItemHeader,
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
				testDataSourceOnPremisesPersonalItemHeader,
				map[string]any{
					"id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by id
		{
			Config: at.CompileConfig(
				testDataSourceOnPremisesPersonalItemHeader,
				map[string]any{
					"id": *entity.ID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceOnPremisesPresonalItemFabricFQN, "id", *entity.ID),
				resource.TestCheckResourceAttrSet(testDataSourceOnPremisesPresonalItemFabricFQN, "public_key.exponent"),
				resource.TestCheckResourceAttrSet(testDataSourceOnPremisesPresonalItemFabricFQN, "public_key.modulus"),
				resource.TestCheckResourceAttrSet(testDataSourceOnPremisesPresonalItemFabricFQN, "version"),
			),
		},
	}))
}

func TestAcc_OnPremisesGatewayPersonalDataSource(t *testing.T) {
	entity := testhelp.WellKnown()["OnPermisesGatewayPersonal"].(map[string]any)
	entityID := entity["id"].(string)
	entityDescription := entity["description"].(string)

	resource.ParallelTest(t, testhelp.NewTestAccCase(t, nil, nil, []resource.TestStep{
		{
			Config: at.CompileConfig(
				testDataSourceOnPremisesPersonalItemHeader,
				map[string]any{
					"id": entityID,
				},
			),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(testDataSourceOnPremisesPresonalItemFabricFQN, "id", entityID),
				resource.TestCheckResourceAttr(testDataSourceOnPremisesPresonalItemFabricFQN, "description", entityDescription),
			),
		},
		// read by id - not found
		{
			Config: at.CompileConfig(
				testDataSourceOnPremisesPersonalItemHeader,
				map[string]any{
					"id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
	}))
}
