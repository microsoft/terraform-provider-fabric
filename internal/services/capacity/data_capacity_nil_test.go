// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package capacity_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

// Test with service principal configuration
func TestUnit_CapacityDataSource_ServicePrincipal(t *testing.T) {
	testState := testhelp.NewTestState()
	
	// Register a fake capacity entity
	testCapacity := fabcore.Capacity{
		ID:          to.Ptr("00000000-0000-0000-0000-000000000000"),
		DisplayName: to.Ptr("Test Capacity"),
		Region:      to.Ptr("westus"),
		SKU:         to.Ptr("F2"),
		State:       to.Ptr(fabcore.CapacityStateActive),
	}
	
	fakes.FakeServer.Upsert(testCapacity)

	// Test a simple data source access with service principal config
	resource.Test(t, testhelp.NewTestUnitCaseWithState(t, nil, fakes.FakeServer.ServerFactory, testState, testhelp.TestUnitPreCheckNoEnvs, []resource.TestStep{
		{
			Config: `
				provider "fabric" {
					tenant_id = "00000000-0000-0000-0000-000000000000"
					client_id = "00000000-0000-0000-0000-000000000000"
					client_secret = "dummy-secret"
				}
				
				data "fabric_capacity" "test" {
					display_name = "Test Capacity"
				}
			`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("data.fabric_capacity.test", "display_name", "Test Capacity"),
				resource.TestCheckResourceAttr("data.fabric_capacity.test", "sku", "F2"),
				resource.TestCheckResourceAttr("data.fabric_capacity.test", "state", "Active"),
			),
		},
	}))
}