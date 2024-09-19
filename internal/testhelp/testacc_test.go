// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package testhelp_test

import (
	"testing"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func TestDevEnv_Provision(t *testing.T) { //nolint:wsl
	// Useful for debugging
	// os.Setenv("FABRIC_TESTACC_WELLKNOWN_CREATE_RESOURCES", "true")
	// os.Setenv("FABRIC_TESTACC_WELLKNOWN_CAPACITY_NAME", "fabtestaccqmjern")

	if testhelp.ShouldCreateWellKnownResources() {
		// create the resources
		testhelp.CreateWellKnownResources()
	}

	// if the file does not exist, the test will fail
	if !testhelp.IsWellKnownDataAvailable() {
		t.Fatalf("well-known resources file does not exist")
	}
}
