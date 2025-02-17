// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gateway_test

import (
	"regexp"
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

var (
	testDataSourceOnPremisesPersonalFQN    = testhelp.DataSourceFQN("fabric", OnPremisesPersonalItemTFName, "test")
	testDataSourceOnPremisesPersonalHeader = at.DataSourceHeader(testhelp.TypeName("fabric", OnPremisesPersonalItemTFName), "test")
)

func TestUnit_OnPremisesGatewayPersonalDataSource(t *testing.T) {
	entity := fakes.NewRandomOnPremisesGatewayPersonal()

	fakes.FakeServer.Upsert(fakes.NewRandomOnPremisesGatewayPersonal())
	fakes.FakeServer.Upsert(entity)
	fakes.FakeServer.Upsert(fakes.NewRandomOnPremisesGatewayPersonal())

	resource.ParallelTest(t, testhelp.NewTestUnitCase(
		t,
		nil,
		fakes.FakeServer.ServerFactory,
		nil,
		[]resource.TestStep{
			// Step 1: Unexpected attribute should trigger an error.
			{
				Config: at.CompileConfig(
					testDataSourceOnPremisesPersonalHeader,
					map[string]any{
						"unexpected_attr": "test",
					},
				),
				ExpectError: regexp.MustCompile(`An argument named "unexpected_attr" is not expected here`),
			},
			// Step 2: Missing ID should trigger an error since ID is required for lookup.
			{
				Config: at.CompileConfig(
					testDataSourceOnPremisesPersonalHeader,
					map[string]any{},
				),
				ExpectError: regexp.MustCompile(`The argument "id" is required`),
			},
			// Step 3: Read by id - not found
			{
				Config: at.CompileConfig(
					testDataSourceOnPremisesItemHeader,
					map[string]any{
						"id": testhelp.RandomUUID(),
					},
				),
				ExpectError: regexp.MustCompile(common.ErrorReadHeader),
			},
			// Step 4: Invalid UUID string should trigger an error.
			{
				Config: at.CompileConfig(
					testDataSourceOnPremisesPersonalHeader,
					map[string]any{
						"id": "invalid uuid",
					},
				),
				ExpectError: regexp.MustCompile(`invalid uuid`),
			},
			// Step 5: Valid read test using the entity's ID.
			{
				Config: at.CompileConfig(
					testDataSourceOnPremisesPersonalHeader,
					map[string]any{
						"id": *entity.ID,
					},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testDataSourceOnPremisesPersonalFQN, "id", *entity.ID),
					resource.TestCheckResourceAttr(testDataSourceOnPremisesPersonalFQN, "version", *entity.Version),
					resource.TestCheckResourceAttrSet(testDataSourceOnPremisesPersonalFQN, "public_key.exponent"),
					resource.TestCheckResourceAttrSet(testDataSourceOnPremisesPersonalFQN, "public_key.modulus"),
				),
			},
		},
	))
}
