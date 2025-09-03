// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package validators_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	fwvalidators "github.com/microsoft/terraform-provider-fabric/internal/framework/validators"
)

// This test is commented out because we need a full integration test
// with proper terraform configuration to test the validator properly.
// The validator requires access to the config schema to work correctly.

/*
func TestPBIRFormatRequiresPagesValidator(t *testing.T) {
	// This would require a full terraform configuration setup to test properly
	// We'll test this through the integration tests in the report package instead
}
*/