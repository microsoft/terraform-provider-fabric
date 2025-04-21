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
	entity := fakes.NewRandomConnection()

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
			ExpectError: regexp.MustCompile(`Exactly one of these attributes must be configured: \[id,display_name\]`),
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
		// error - conflicting attributes
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
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
				testDataSourceItemHeader,
				map[string]any{
					"id": testhelp.RandomUUID(),
				},
			),
			ExpectError: regexp.MustCompile(common.ErrorReadHeader),
		},
		// read by name - not found
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"display_name": testhelp.RandomName(),
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
					tfjsonpath.New("gateway_id"),
					knownvalue.StringExact(*entity.GatewayID),
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
		// read by name
		{
			Config: at.CompileConfig(
				testDataSourceItemHeader,
				map[string]any{
					"display_name": *entity.DisplayName,
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
					tfjsonpath.New("gateway_id"),
					knownvalue.StringExact(*entity.GatewayID),
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
