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

// Test to reproduce the inconsistent pool type issue
func TestPoolTypeInconsistency(t *testing.T) {
	ctx := context.Background()

	t.Run("API returns different pool type than requested", func(t *testing.T) {
		// Simulate user configuration: pool.type = "Capacity", pool.name = "Some Pool"
		poolModel := &instancePoolPropertiesModel{
			ID:   customtypes.NewUUIDValue("12345678-1234-1234-1234-123456789012"),
			Name: types.StringValue("Some Pool"),
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
		assert.False(t, diags.HasError())
		
		// Verify request has correct type
		assert.NotNil(t, request.InstancePool)
		assert.Equal(t, "Capacity", string(*request.InstancePool.Type))
		assert.Equal(t, "Some Pool", *request.InstancePool.Name)

		// Simulate API returning different configuration (e.g., API ignores the request and returns default)
		// This might happen if "Some Pool" doesn't exist or isn't a capacity pool
		workspaceType := environment.CustomPoolTypeWorkspace
		apiResponse := environment.InstancePool{
			ID:   stringPtrHelper("00000000-0000-0000-0000-000000000000"), // Default starter pool ID
			Name: stringPtrHelper("Starter Pool"),                          // API returns default pool
			Type: &workspaceType,                                     // With workspace type, not capacity
		}

		// Process API response
		var responseModel baseSparkEnvironmentSettingsModel
		sparkCompute := environment.SparkCompute{
			InstancePool: &apiResponse,
		}

		diags = responseModel.set(ctx, sparkCompute)
		assert.False(t, diags.HasError())

		// Verify what the response model contains
		poolResult, diags := responseModel.Pool.Get(ctx)
		assert.False(t, diags.HasError())

		// This is the inconsistency - user requested "Capacity" but got "Workspace"
		assert.Equal(t, "Starter Pool", poolResult.Name.ValueString())
		assert.Equal(t, "Workspace", poolResult.Type.ValueString()) // This should be "Capacity" but it's "Workspace"

		t.Logf("User requested: name=%s, type=%s", "Some Pool", "Capacity")
		t.Logf("API returned: name=%s, type=%s", poolResult.Name.ValueString(), poolResult.Type.ValueString())
	})
}

func stringPtrHelper(s string) *string {
	return &s
}