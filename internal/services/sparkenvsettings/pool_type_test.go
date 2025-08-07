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

// Test the pool type conversion behavior
func TestPoolTypeConversion(t *testing.T) {
	ctx := context.Background()

	t.Run("Pool type Capacity should remain Capacity", func(t *testing.T) {
		// Create a mock InstancePool with Type Capacity
		capacityType := environment.CustomPoolTypeCapacity
		mockPool := environment.InstancePool{
			ID:   stringPtr("test-id"),
			Name: stringPtr("Test Pool"),
			Type: &capacityType,
		}

		// Test the set method
		var model instancePoolPropertiesModel
		model.set(mockPool)

		// Verify the Type is preserved
		assert.Equal(t, "Capacity", model.Type.ValueString())
	})

	t.Run("Pool type Workspace should remain Workspace", func(t *testing.T) {
		// Create a mock InstancePool with Type Workspace
		workspaceType := environment.CustomPoolTypeWorkspace
		mockPool := environment.InstancePool{
			ID:   stringPtr("test-id"),
			Name: stringPtr("Test Pool"),
			Type: &workspaceType,
		}

		// Test the set method
		var model instancePoolPropertiesModel
		model.set(mockPool)

		// Verify the Type is preserved
		assert.Equal(t, "Workspace", model.Type.ValueString())
	})

	t.Run("Request should preserve pool type", func(t *testing.T) {
		// Create a model with pool type Capacity
		poolModel := &instancePoolPropertiesModel{
			ID:   customtypes.NewUUIDValue("12345678-1234-1234-1234-123456789012"),
			Name: types.StringValue("Test Pool"),
			Type: types.StringValue("Capacity"),
		}

		poolValue := supertypes.NewSingleNestedObjectValueOfNull[instancePoolPropertiesModel](ctx)
		diags := poolValue.Set(ctx, poolModel)
		if diags.HasError() {
			t.Logf("Pool value set errors: %v", diags.Errors())
		}
		assert.False(t, diags.HasError())
		
		resourceModel := resourceSparkEnvironmentSettingsModel{}
		resourceModel.Pool = poolValue

		// Test the request conversion
		var request requestUpdateSparkEnvironmentSettings
		diags = request.set(ctx, resourceModel)

		// Log any errors
		if diags.HasError() {
			for _, err := range diags.Errors() {
				t.Logf("Request conversion error: %s", err.Detail())
			}
		}

		// Should not have errors
		assert.False(t, diags.HasError())
		
		// Verify the request has the correct type
		if request.InstancePool == nil {
			t.Error("InstancePool is nil")
			return
		}
		assert.NotNil(t, request.InstancePool.Type)
		assert.Equal(t, "Capacity", string(*request.InstancePool.Type))
	})

	t.Run("Request should handle missing pool name correctly", func(t *testing.T) {
		// Create a model with pool type but no name
		poolModel := &instancePoolPropertiesModel{
			Type: types.StringValue("Capacity"),
		}

		poolValue := supertypes.NewSingleNestedObjectValueOfNull[instancePoolPropertiesModel](ctx)
		diags := poolValue.Set(ctx, poolModel)
		assert.False(t, diags.HasError())
		
		resourceModel := resourceSparkEnvironmentSettingsModel{}
		resourceModel.Pool = poolValue

		// Test the request conversion
		var request requestUpdateSparkEnvironmentSettings
		diags = request.set(ctx, resourceModel)

		// Should not have errors
		assert.False(t, diags.HasError())
		
		// Verify the request has the correct type but no name
		assert.NotNil(t, request.InstancePool)
		assert.NotNil(t, request.InstancePool.Type)
		assert.Equal(t, "Capacity", string(*request.InstancePool.Type))
		assert.Nil(t, request.InstancePool.Name)
	})
}

func stringPtr(s string) *string {
	return &s
}