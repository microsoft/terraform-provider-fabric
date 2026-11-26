// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package connection_test

import (
	"testing"

	at "github.com/dcarbone/terraform-plugin-framework-utils/v3/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

func TestUnit_ConnectionDataSources_Classifications(t *testing.T) {
	server := fakes.NewFakeServer()
	entities := make([]fabcore.ConnectionClassification, 0, 7)

	for _, connectivityType := range []fabcore.ConnectivityType{
		fabcore.ConnectivityTypeShareableCloud, fabcore.ConnectivityTypePersonalCloud,
		fabcore.ConnectivityTypeOnPremisesGateway, fabcore.ConnectivityTypeOnPremisesGatewayPersonal,
		fabcore.ConnectivityTypeVirtualNetworkGateway, fabcore.ConnectivityTypeStreamingVirtualNetworkGateway,
		fabcore.ConnectivityTypeAutomatic,
	} {
		base := fabcore.Connection{
			ID: new(testhelp.RandomUUID()), DisplayName: new(string(connectivityType)), ConnectivityType: new(connectivityType),
			ConnectionDetails: &fabcore.ListConnectionDetails{Type: new("SQL"), Path: new("existing.example.com;database")},
		}
		base.GatewayID = new(testhelp.RandomUUID())
		var entity fabcore.ConnectionClassification

		switch connectivityType {
		case fabcore.ConnectivityTypeShareableCloud:
			entity = &fabcore.ShareableCloudConnection{
				ID: base.ID, DisplayName: base.DisplayName, ConnectivityType: base.ConnectivityType, ConnectionDetails: base.ConnectionDetails,
				GatewayID: base.GatewayID, AllowConnectionUsageInGateway: new(true), AllowUsageInUserControlledCode: new(true),
			}
		case fabcore.ConnectivityTypePersonalCloud:
			entity = &fabcore.PersonalCloudConnection{
				ID: base.ID, DisplayName: base.DisplayName, ConnectivityType: base.ConnectivityType, ConnectionDetails: base.ConnectionDetails,
				GatewayID: base.GatewayID, AllowConnectionUsageInGateway: new(true),
				CredentialDetails: &fabcore.ListCredentialDetails{CredentialType: new(fabcore.CredentialTypeOAuth2)},
			}
		case fabcore.ConnectivityTypeOnPremisesGateway:
			base.CredentialDetails = &fabcore.ListCredentialDetails{CredentialType: new(fabcore.CredentialTypeWindows)}
			connection := fabcore.OnPremisesGatewayConnection(base)
			entity = &connection
		case fabcore.ConnectivityTypeOnPremisesGatewayPersonal:
			base.CredentialDetails = &fabcore.ListCredentialDetails{CredentialType: new(fabcore.CredentialTypeWindowsWithoutImpersonation)}
			connection := fabcore.OnPremisesGatewayPersonalConnection(base)
			entity = &connection
		case fabcore.ConnectivityTypeVirtualNetworkGateway:
			connection := fabcore.VirtualNetworkGatewayConnection(base)
			entity = &connection
		case fabcore.ConnectivityTypeStreamingVirtualNetworkGateway:
			connection := fabcore.StreamingVirtualNetworkGatewayConnection(base)
			entity = &connection
		default:
			entity = &base
		}

		server.Upsert(entity)
		entities = append(entities, entity)
	}

	steps := make([]resource.TestStep, 0, len(entities)+1)
	listChecks := make([]knownvalue.Check, 0, len(entities))

	for _, entity := range entities {
		connection := entity.GetConnection()

		checks := map[string]knownvalue.Check{
			"id": knownvalue.StringExact(*connection.ID), "display_name": knownvalue.StringExact(*connection.DisplayName),
			"connectivity_type": knownvalue.StringExact(string(*connection.ConnectivityType)),
			"gateway_id":        knownvalue.Null(), "privacy_level": knownvalue.Null(),
			"connection_details": knownvalue.ObjectExact(map[string]knownvalue.Check{
				"type": knownvalue.StringExact(*connection.ConnectionDetails.Type), "path": knownvalue.StringExact(*connection.ConnectionDetails.Path),
			}),
			"allow_connection_usage_in_gateway": knownvalue.Null(), "allow_usage_in_user_controlled_code": knownvalue.Null(),
			"credential_details": knownvalue.Null(),
		}
		if connection.GatewayID != nil {
			checks["gateway_id"] = knownvalue.StringExact(*connection.GatewayID)
		}

		if *connection.ConnectivityType == fabcore.ConnectivityTypeShareableCloud || *connection.ConnectivityType == fabcore.ConnectivityTypePersonalCloud {
			checks["allow_connection_usage_in_gateway"] = knownvalue.Bool(true)
		}

		if *connection.ConnectivityType == fabcore.ConnectivityTypeShareableCloud {
			checks["allow_usage_in_user_controlled_code"] = knownvalue.Bool(true)
		}

		if connection.CredentialDetails != nil {
			checks["credential_details"] = knownvalue.ObjectExact(map[string]knownvalue.Check{
				"credential_type":       knownvalue.StringExact(string(*connection.CredentialDetails.CredentialType)),
				"connection_encryption": knownvalue.Null(), "single_sign_on_type": knownvalue.Null(), "skip_test_connection": knownvalue.Null(),
			})
		}

		listChecks = append(listChecks, knownvalue.ObjectPartial(checks))

		stateChecks := make([]statecheck.StateCheck, 0, len(checks))
		for name, check := range checks {
			stateChecks = append(stateChecks, statecheck.ExpectKnownValue(testDataSourceItemFQN, tfjsonpath.New(name), check))
		}

		steps = append(steps, resource.TestStep{
			Config: at.CompileConfig(testDataSourceItemHeader, map[string]any{"id": *connection.ID}), ConfigStateChecks: stateChecks,
		})
	}

	steps = append(steps, resource.TestStep{
		Config: at.CompileConfig(testDataSourceItemsHeader, map[string]any{}),
		Check:  resource.TestCheckResourceAttr(testDataSourceItemsFQN, "values.#", "7"),
		ConfigStateChecks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue(testDataSourceItemsFQN, tfjsonpath.New("values"), knownvalue.SetExact(listChecks)),
		},
	})
	resource.ParallelTest(t, testhelp.NewTestUnitCase(t, nil, server.ServerFactory, nil, steps))
}
