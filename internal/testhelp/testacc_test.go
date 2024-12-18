// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package testhelp_test

import (
	"testing"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
)

func TestDevEnv_WellKnown(t *testing.T) {
	// if the file does not exist, the test will fail
	if !testhelp.IsWellKnownDataAvailable() {
		t.Fatalf("well-known resources file does not exist")
	}
}
