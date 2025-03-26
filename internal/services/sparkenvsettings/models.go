// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sparkenvsettings

import (
	"context"
	"encoding/json"

	timeoutsd "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsr "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseSparkEnvironmentSettingsModel struct {
	ID                        customtypes.UUID                                                               `tfsdk:"id"`
	WorkspaceID               customtypes.UUID                                                               `tfsdk:"workspace_id"`
	EnvironmentID             customtypes.UUID                                                               `tfsdk:"environment_id"`
	PublicationStatus         types.String                                                                   `tfsdk:"publication_status"`
	DriverCores               types.Int32                                                                    `tfsdk:"driver_cores"`
	DriverMemory              types.String                                                                   `tfsdk:"driver_memory"`
	DynamicExecutorAllocation supertypes.SingleNestedObjectValueOf[dynamicExecutorAllocationPropertiesModel] `tfsdk:"dynamic_executor_allocation"`
	ExecutorCores             types.Int32                                                                    `tfsdk:"executor_cores"`
	ExecutorMemory            types.String                                                                   `tfsdk:"executor_memory"`
	Pool                      supertypes.SingleNestedObjectValueOf[instancePoolPropertiesModel]              `tfsdk:"pool"`
	RuntimeVersion            types.String                                                                   `tfsdk:"runtime_version"`
	SparkProperties           customtypes.MapOfString                                                        `tfsdk:"spark_properties"`
}

func (to *baseSparkEnvironmentSettingsModel) set(ctx context.Context, from fabenvironment.SparkCompute) diag.Diagnostics {
	var diags diag.Diagnostics

	to.DriverCores = types.Int32PointerValue(from.DriverCores)
	to.DriverMemory = types.StringPointerValue(from.DriverMemory)
	to.ExecutorCores = types.Int32PointerValue(from.ExecutorCores)
	to.ExecutorMemory = types.StringPointerValue(from.ExecutorMemory)
	to.RuntimeVersion = types.StringPointerValue(from.RuntimeVersion)

	var sparkProperties map[string]string

	sparkPropertiesBytes, err := json.Marshal(from.SparkProperties)
	if err != nil {
		diags.AddError(
			"failed to marshal Spark properties",
			err.Error(),
		)

		return diags
	}

	err = json.Unmarshal(sparkPropertiesBytes, &sparkProperties)
	if err != nil {
		diags.AddError(
			"failed to unmarshal Spark properties",
			err.Error(),
		)

		return diags
	}

	sparkPropertiesMap := customtypes.NewMapValueOfNull[types.String](ctx)

	if len(sparkProperties) > 0 {
		sparkPropertiesTF := make(map[string]types.String)

		for k, v := range sparkProperties {
			sparkPropertiesTF[k] = types.StringValue(v)
		}

		sparkPropertiesMap, diags = customtypes.NewMapValueOf(ctx, sparkPropertiesTF)
		if diags.HasError() {
			return diags
		}
	}

	to.SparkProperties = sparkPropertiesMap

	dynamicExecutorAllocation := supertypes.NewSingleNestedObjectValueOfNull[dynamicExecutorAllocationPropertiesModel](ctx)

	if from.DynamicExecutorAllocation != nil {
		dynamicExecutorAllocationModel := &dynamicExecutorAllocationPropertiesModel{}
		dynamicExecutorAllocationModel.set(*from.DynamicExecutorAllocation)

		if diags := dynamicExecutorAllocation.Set(ctx, dynamicExecutorAllocationModel); diags.HasError() {
			return diags
		}
	}

	to.DynamicExecutorAllocation = dynamicExecutorAllocation

	instancePool := supertypes.NewSingleNestedObjectValueOfNull[instancePoolPropertiesModel](ctx)

	if from.InstancePool != nil {
		instancePoolModel := &instancePoolPropertiesModel{}
		instancePoolModel.set(*from.InstancePool)

		if diags := instancePool.Set(ctx, instancePoolModel); diags.HasError() {
			return diags
		}
	}

	to.Pool = instancePool

	return nil
}

/*
DATA-SOURCE
*/

type dataSourceSparkEnvironmentSettingsModel struct {
	baseSparkEnvironmentSettingsModel
	Timeouts timeoutsd.Value `tfsdk:"timeouts"`
}

/*
RESOURCE
*/

type resourceSparkEnvironmentSettingsModel struct {
	baseSparkEnvironmentSettingsModel
	Timeouts timeoutsr.Value `tfsdk:"timeouts"`
}

type requestUpdateSparkEnvironmentSettings struct {
	fabenvironment.UpdateEnvironmentSparkComputeRequest
}

func (to *requestUpdateSparkEnvironmentSettings) set(ctx context.Context, from resourceSparkEnvironmentSettingsModel) diag.Diagnostics { //nolint:gocognit, gocyclo
	if !from.DriverCores.IsNull() && !from.DriverCores.IsUnknown() {
		to.DriverCores = from.DriverCores.ValueInt32Pointer()
	}

	if !from.DriverMemory.IsNull() && !from.DriverMemory.IsUnknown() {
		to.DriverMemory = from.DriverMemory.ValueStringPointer()
	}

	if !from.ExecutorCores.IsNull() && !from.ExecutorCores.IsUnknown() {
		to.ExecutorCores = from.ExecutorCores.ValueInt32Pointer()
	}

	if !from.ExecutorMemory.IsNull() && !from.ExecutorMemory.IsUnknown() {
		to.ExecutorMemory = from.ExecutorMemory.ValueStringPointer()
	}

	if !from.RuntimeVersion.IsNull() && !from.RuntimeVersion.IsUnknown() {
		to.RuntimeVersion = from.RuntimeVersion.ValueStringPointer()
	}

	if !from.DynamicExecutorAllocation.IsNull() && !from.DynamicExecutorAllocation.IsUnknown() {
		dynamicExecutorAllocation, diags := from.DynamicExecutorAllocation.Get(ctx)
		if diags.HasError() {
			return diags
		}

		var reqDynamicExecutorAllocation fabenvironment.DynamicExecutorAllocationProperties

		if !dynamicExecutorAllocation.Enabled.IsNull() && !dynamicExecutorAllocation.Enabled.IsUnknown() {
			reqDynamicExecutorAllocation.Enabled = dynamicExecutorAllocation.Enabled.ValueBoolPointer()
		}

		if !dynamicExecutorAllocation.MinExecutors.IsNull() && !dynamicExecutorAllocation.MinExecutors.IsUnknown() {
			reqDynamicExecutorAllocation.MinExecutors = dynamicExecutorAllocation.MinExecutors.ValueInt32Pointer()
		}

		if !dynamicExecutorAllocation.MaxExecutors.IsNull() && !dynamicExecutorAllocation.MaxExecutors.IsUnknown() {
			reqDynamicExecutorAllocation.MaxExecutors = dynamicExecutorAllocation.MaxExecutors.ValueInt32Pointer()
		}

		if reqDynamicExecutorAllocation != (fabenvironment.DynamicExecutorAllocationProperties{}) {
			to.DynamicExecutorAllocation = &reqDynamicExecutorAllocation
		}
	}

	if !from.Pool.IsNull() && !from.Pool.IsUnknown() {
		pool, diags := from.Pool.Get(ctx)
		if diags.HasError() {
			return diags
		}

		var reqPool fabenvironment.InstancePool

		if !pool.ID.IsNull() && !pool.ID.IsUnknown() {
			reqPool.ID = pool.ID.ValueStringPointer()
		}

		if !pool.Name.IsNull() && !pool.Name.IsUnknown() {
			reqPool.Name = pool.Name.ValueStringPointer()
		}

		if !pool.Type.IsNull() && !pool.Type.IsUnknown() {
			reqPool.Type = (*fabenvironment.CustomPoolType)(pool.Type.ValueStringPointer())
		}

		if reqPool != (fabenvironment.InstancePool{}) {
			to.InstancePool = &reqPool
		}
	}

	if !from.SparkProperties.IsNull() && !from.SparkProperties.IsUnknown() {
		sparkProperties, diags := from.SparkProperties.Get(ctx)
		if diags.HasError() {
			return diags
		}

		sparkPropertiesMap := make(map[string]string)

		for k, v := range sparkProperties {
			if !v.IsNull() && !v.IsUnknown() {
				sparkPropertiesMap[k] = v.ValueString()
			}
		}

		to.SparkProperties = sparkPropertiesMap
	}

	return nil
}

/*
HELPER MODELS
*/

type dynamicExecutorAllocationPropertiesModel struct {
	Enabled      types.Bool  `tfsdk:"enabled"`
	MinExecutors types.Int32 `tfsdk:"min_executors"`
	MaxExecutors types.Int32 `tfsdk:"max_executors"`
}

func (to *dynamicExecutorAllocationPropertiesModel) set(from fabenvironment.DynamicExecutorAllocationProperties) {
	to.Enabled = types.BoolPointerValue(from.Enabled)
	to.MinExecutors = types.Int32PointerValue(from.MinExecutors)
	to.MaxExecutors = types.Int32PointerValue(from.MaxExecutors)
}

type instancePoolPropertiesModel struct {
	ID   customtypes.UUID `tfsdk:"id"`
	Name types.String     `tfsdk:"name"`
	Type types.String     `tfsdk:"type"`
}

func (to *instancePoolPropertiesModel) set(from fabenvironment.InstancePool) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Name = types.StringPointerValue(from.Name)
	to.Type = types.StringPointerValue((*string)(from.Type))
}
