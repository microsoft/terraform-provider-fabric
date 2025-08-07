// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sparkenvsettings

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/fabric-sdk-go/fabric/environment"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	"github.com/stretchr/testify/assert"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

// Test the validation function that detects pool configuration mismatches
func TestValidatePoolConfiguration(t *testing.T) {
	ctx := context.Background()
	resource := &resourceSparkEnvironmentSettings{}

	t.Run("Valid configuration should pass validation", func(t *testing.T) {
		// Create request with Capacity pool
		capacityType := environment.CustomPoolTypeCapacity
		request := &requestUpdateSparkEnvironmentSettings{}
		request.InstancePool = &environment.InstancePool{
			Name: stringPtrHelper("My Capacity Pool"),
			Type: &capacityType,
		}

		// Create model with matching configuration
		poolModel := &instancePoolPropertiesModel{
			ID:   customtypes.NewUUIDValue("12345678-1234-1234-1234-123456789012"),
			Name: types.StringValue("My Capacity Pool"),
			Type: types.StringValue("Capacity"),
		}

		poolValue := supertypes.NewSingleNestedObjectValueOfNull[instancePoolPropertiesModel](ctx)
		diags := poolValue.Set(ctx, poolModel)
		assert.False(t, diags.HasError())

		model := &resourceSparkEnvironmentSettingsModel{}
		model.Pool = poolValue

		// Validation should pass
		diags = resource.validatePoolConfiguration(ctx, request, model)
		assert.False(t, diags.HasError())
	})

	t.Run("Type mismatch should be detected", func(t *testing.T) {
		// Create request with Capacity pool
		capacityType := environment.CustomPoolTypeCapacity
		request := &requestUpdateSparkEnvironmentSettings{}
		request.InstancePool = &environment.InstancePool{
			Name: stringPtrHelper("My Pool"),
			Type: &capacityType,
		}

		// Create model with Workspace type (mismatch)
		poolModel := &instancePoolPropertiesModel{
			ID:   customtypes.NewUUIDValue("12345678-1234-1234-1234-123456789012"),
			Name: types.StringValue("Starter Pool"),
			Type: types.StringValue("Workspace"),
		}

		poolValue := supertypes.NewSingleNestedObjectValueOfNull[instancePoolPropertiesModel](ctx)
		diags := poolValue.Set(ctx, poolModel)
		assert.False(t, diags.HasError())

		model := &resourceSparkEnvironmentSettingsModel{}
		model.Pool = poolValue

		// Validation should fail with clear error
		diags = resource.validatePoolConfiguration(ctx, request, model)
		assert.True(t, diags.HasError())
		
		errors := diags.Errors()
		assert.Len(t, errors, 2) // Both name and type mismatches
		
		// Check that error mentions the type mismatch
		typeError := errors[0].Detail()
		assert.Contains(t, typeError, "Requested pool type 'Capacity'")
		assert.Contains(t, typeError, "API returned pool type 'Workspace'")
		
		// Check that error mentions the name mismatch  
		nameError := errors[1].Detail()
		assert.Contains(t, nameError, "Requested pool name 'My Pool'")
		assert.Contains(t, nameError, "API returned pool name 'Starter Pool'")
	})

	t.Run("Missing pool configuration should be detected", func(t *testing.T) {
		// Create request with pool
		capacityType := environment.CustomPoolTypeCapacity
		request := &requestUpdateSparkEnvironmentSettings{}
		request.InstancePool = &environment.InstancePool{
			Name: stringPtrHelper("My Pool"),
			Type: &capacityType,
		}

		// Create model with null pool (no pool returned by API)
		model := &resourceSparkEnvironmentSettingsModel{}
		model.Pool = supertypes.NewSingleNestedObjectValueOfNull[instancePoolPropertiesModel](ctx)

		// Validation should fail
		diags := resource.validatePoolConfiguration(ctx, request, model)
		assert.True(t, diags.HasError())
		
		errors := diags.Errors()
		assert.Len(t, errors, 1)
		assert.Contains(t, errors[0].Detail(), "pool configuration was specified")
	})

	t.Run("No request should skip validation", func(t *testing.T) {
		// Create request with no pool
		request := &requestUpdateSparkEnvironmentSettings{}
		request.InstancePool = nil

		// Create model - doesn't matter what it contains
		model := &resourceSparkEnvironmentSettingsModel{}

		// Validation should pass (nothing to validate)
		diags := resource.validatePoolConfiguration(ctx, request, model)
		assert.False(t, diags.HasError())
	})
}