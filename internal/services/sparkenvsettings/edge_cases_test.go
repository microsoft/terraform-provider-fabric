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

// TestPoolConfigurationEdgeCases tests various edge cases for the pool configuration validation
func TestPoolConfigurationEdgeCases(t *testing.T) {
	ctx := context.Background()
	resource := &resourceSparkEnvironmentSettings{}

	t.Run("Type only mismatch should be detected", func(t *testing.T) {
		// User requests pool with only type (no name), API returns default pool
		capacityType := environment.CustomPoolTypeCapacity
		request := &requestUpdateSparkEnvironmentSettings{}
		request.InstancePool = &environment.InstancePool{
			Type: &capacityType, // Only type specified
			// Name is nil
		}

		// API returns default Starter Pool
		poolModel := &instancePoolPropertiesModel{
			ID:   customtypes.NewUUIDValue("00000000-0000-0000-0000-000000000000"),
			Name: types.StringValue("Starter Pool"),
			Type: types.StringValue("Workspace"),
		}

		poolValue := supertypes.NewSingleNestedObjectValueOfNull[instancePoolPropertiesModel](ctx)
		diags := poolValue.Set(ctx, poolModel)
		assert.False(t, diags.HasError())

		model := &resourceSparkEnvironmentSettingsModel{}
		model.Pool = poolValue

		// Validation should detect type mismatch
		diags = resource.validatePoolConfiguration(ctx, request, model)
		assert.True(t, diags.HasError())
		
		errors := diags.Errors()
		assert.Len(t, errors, 1) // Only type error, no name error since name wasn't requested
		assert.Contains(t, errors[0].Detail(), "pool type 'Capacity'")
		assert.Contains(t, errors[0].Detail(), "returned pool type 'Workspace'")
	})

	t.Run("Name only mismatch should be detected", func(t *testing.T) {
		// User requests pool with only name (no type)
		request := &requestUpdateSparkEnvironmentSettings{}
		request.InstancePool = &environment.InstancePool{
			Name: stringPtrHelper("My Custom Pool"),
			// Type is nil
		}

		// API returns default Starter Pool
		poolModel := &instancePoolPropertiesModel{
			ID:   customtypes.NewUUIDValue("00000000-0000-0000-0000-000000000000"),
			Name: types.StringValue("Starter Pool"),
			Type: types.StringValue("Workspace"),
		}

		poolValue := supertypes.NewSingleNestedObjectValueOfNull[instancePoolPropertiesModel](ctx)
		diags := poolValue.Set(ctx, poolModel)
		assert.False(t, diags.HasError())

		model := &resourceSparkEnvironmentSettingsModel{}
		model.Pool = poolValue

		// Validation should detect name mismatch
		diags = resource.validatePoolConfiguration(ctx, request, model)
		assert.True(t, diags.HasError())
		
		errors := diags.Errors()
		assert.Len(t, errors, 1) // Only name error, no type error since type wasn't requested
		assert.Contains(t, errors[0].Detail(), "pool name 'My Custom Pool'")
		assert.Contains(t, errors[0].Detail(), "returned pool name 'Starter Pool'")
	})

	t.Run("Empty pool request should pass validation", func(t *testing.T) {
		// User requests empty pool (this shouldn't normally happen due to validation)
		request := &requestUpdateSparkEnvironmentSettings{}
		request.InstancePool = &environment.InstancePool{
			// Both Name and Type are nil
		}

		// API returns some pool
		poolModel := &instancePoolPropertiesModel{
			ID:   customtypes.NewUUIDValue("00000000-0000-0000-0000-000000000000"),
			Name: types.StringValue("Starter Pool"),
			Type: types.StringValue("Workspace"),
		}

		poolValue := supertypes.NewSingleNestedObjectValueOfNull[instancePoolPropertiesModel](ctx)
		diags := poolValue.Set(ctx, poolModel)
		assert.False(t, diags.HasError())

		model := &resourceSparkEnvironmentSettingsModel{}
		model.Pool = poolValue

		// Validation should pass since no specific configuration was requested
		diags = resource.validatePoolConfiguration(ctx, request, model)
		assert.False(t, diags.HasError())
	})

	t.Run("API returning null pool fields should be handled", func(t *testing.T) {
		// User requests pool configuration
		capacityType := environment.CustomPoolTypeCapacity
		request := &requestUpdateSparkEnvironmentSettings{}
		request.InstancePool = &environment.InstancePool{
			Name: stringPtrHelper("My Pool"),
			Type: &capacityType,
		}

		// API returns pool with null fields
		poolModel := &instancePoolPropertiesModel{
			ID:   customtypes.NewUUIDValue("12345678-1234-1234-1234-123456789012"),
			Name: types.StringNull(),
			Type: types.StringNull(),
		}

		poolValue := supertypes.NewSingleNestedObjectValueOfNull[instancePoolPropertiesModel](ctx)
		diags := poolValue.Set(ctx, poolModel)
		assert.False(t, diags.HasError())

		model := &resourceSparkEnvironmentSettingsModel{}
		model.Pool = poolValue

		// Validation should detect mismatches and handle null values gracefully
		diags = resource.validatePoolConfiguration(ctx, request, model)
		assert.True(t, diags.HasError())
		
		errors := diags.Errors()
		assert.Len(t, errors, 2) // Both type and name errors
		
		// Check that null values are reported correctly
		typeError := errors[0].Detail()
		assert.Contains(t, typeError, "pool type 'Capacity'")
		assert.Contains(t, typeError, "returned pool type 'null'")
		
		nameError := errors[1].Detail()
		assert.Contains(t, nameError, "pool name 'My Pool'")
		assert.Contains(t, nameError, "returned pool name 'null'")
	})

	t.Run("Successful validation with all fields matching", func(t *testing.T) {
		// User requests specific pool configuration
		capacityType := environment.CustomPoolTypeCapacity
		request := &requestUpdateSparkEnvironmentSettings{}
		request.InstancePool = &environment.InstancePool{
			Name: stringPtrHelper("Production Capacity Pool"),
			Type: &capacityType,
		}

		// API returns exactly what was requested
		poolModel := &instancePoolPropertiesModel{
			ID:   customtypes.NewUUIDValue("12345678-1234-1234-1234-123456789012"),
			Name: types.StringValue("Production Capacity Pool"),
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
}