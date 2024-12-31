// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package spark

import (
	"context"

	timeoutsd "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsr "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabspark "github.com/microsoft/fabric-sdk-go/fabric/spark"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceSparkWorkspaceSettingsModel struct {
	baseSparkWorkspaceSettingsModel
	Timeouts timeoutsd.Value `tfsdk:"timeouts"`
}

type resourceSparkWorkspaceSettingsModel struct {
	baseSparkWorkspaceSettingsModel
	Timeouts timeoutsr.Value `tfsdk:"timeouts"`
}

type baseSparkWorkspaceSettingsModel struct {
	ID              customtypes.UUID                                                     `tfsdk:"id"`
	WorkspaceID     customtypes.UUID                                                     `tfsdk:"workspace_id"`
	AutomaticLog    supertypes.SingleNestedObjectValueOf[automaticLogPropertiesModel]    `tfsdk:"automatic_log"`
	Environment     supertypes.SingleNestedObjectValueOf[environmentPropertiesModel]     `tfsdk:"environment"`
	HighConcurrency supertypes.SingleNestedObjectValueOf[highConcurrencyPropertiesModel] `tfsdk:"high_concurrency"`
	Pool            supertypes.SingleNestedObjectValueOf[poolPropertiesModel]            `tfsdk:"pool"`
}

func (to *baseSparkWorkspaceSettingsModel) set(ctx context.Context, from fabspark.WorkspaceSparkSettings) diag.Diagnostics {
	automaticLog := supertypes.NewSingleNestedObjectValueOfNull[automaticLogPropertiesModel](ctx)

	if from.AutomaticLog != nil {
		automaticLogModel := &automaticLogPropertiesModel{}

		automaticLogModel.set(from.AutomaticLog)

		if diags := automaticLog.Set(ctx, automaticLogModel); diags.HasError() {
			return diags
		}
	}

	to.AutomaticLog = automaticLog

	environment := supertypes.NewSingleNestedObjectValueOfNull[environmentPropertiesModel](ctx)

	if from.Environment != nil {
		environmentModel := &environmentPropertiesModel{}
		environmentModel.set(from.Environment)

		if diags := environment.Set(ctx, environmentModel); diags.HasError() {
			return diags
		}
	}

	to.Environment = environment

	highConcurrency := supertypes.NewSingleNestedObjectValueOfNull[highConcurrencyPropertiesModel](ctx)

	if from.HighConcurrency != nil {
		highConcurrencyModel := &highConcurrencyPropertiesModel{}
		highConcurrencyModel.set(from.HighConcurrency)

		if diags := highConcurrency.Set(ctx, highConcurrencyModel); diags.HasError() {
			return diags
		}
	}

	to.HighConcurrency = highConcurrency

	pool := supertypes.NewSingleNestedObjectValueOfNull[poolPropertiesModel](ctx)

	if from.Pool != nil {
		poolModel := &poolPropertiesModel{}

		if diags := poolModel.set(ctx, from.Pool); diags.HasError() {
			return diags
		}

		if diags := pool.Set(ctx, poolModel); diags.HasError() {
			return diags
		}
	}

	to.Pool = pool

	return nil
}

type automaticLogPropertiesModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

func (to *automaticLogPropertiesModel) set(from *fabspark.AutomaticLogProperties) {
	to.Enabled = types.BoolPointerValue(from.Enabled)
}

type environmentPropertiesModel struct {
	Name           types.String `tfsdk:"name"`
	RuntimeVersion types.String `tfsdk:"runtime_version"`
}

func (to *environmentPropertiesModel) set(from *fabspark.EnvironmentProperties) {
	to.Name = types.StringPointerValue(from.Name)
	to.RuntimeVersion = types.StringPointerValue(from.RuntimeVersion)
}

type highConcurrencyPropertiesModel struct {
	NotebookInteractiveRunEnabled types.Bool `tfsdk:"notebook_interactive_run_enabled"`
}

func (to *highConcurrencyPropertiesModel) set(from *fabspark.HighConcurrencyProperties) {
	to.NotebookInteractiveRunEnabled = types.BoolPointerValue(from.NotebookInteractiveRunEnabled)
}

type poolPropertiesModel struct {
	CustomizeComputeEnabled types.Bool                                                       `tfsdk:"customize_compute_enabled"`
	DefaultPool             supertypes.SingleNestedObjectValueOf[defaultPoolPropertiesModel] `tfsdk:"default_pool"`
	StarterPool             supertypes.SingleNestedObjectValueOf[starterPoolPropertiesModel] `tfsdk:"starter_pool"`
}

func (to *poolPropertiesModel) set(ctx context.Context, from *fabspark.PoolProperties) diag.Diagnostics {
	to.CustomizeComputeEnabled = types.BoolPointerValue(from.CustomizeComputeEnabled)

	defaultPool := supertypes.NewSingleNestedObjectValueOfNull[defaultPoolPropertiesModel](ctx)

	if from.DefaultPool != nil {
		defaultPoolModel := &defaultPoolPropertiesModel{}
		defaultPoolModel.set(from.DefaultPool)

		if diags := defaultPool.Set(ctx, defaultPoolModel); diags.HasError() {
			return diags
		}
	}

	to.DefaultPool = defaultPool

	starterPool := supertypes.NewSingleNestedObjectValueOfNull[starterPoolPropertiesModel](ctx)

	if from.StarterPool != nil {
		starterPoolModel := &starterPoolPropertiesModel{}
		starterPoolModel.set(from.StarterPool)

		if diags := starterPool.Set(ctx, starterPoolModel); diags.HasError() {
			return diags
		}
	}

	to.StarterPool = starterPool

	return nil
}

type defaultPoolPropertiesModel struct {
	ID   customtypes.UUID `tfsdk:"id"`
	Name types.String     `tfsdk:"name"`
	Type types.String     `tfsdk:"type"`
}

func (to *defaultPoolPropertiesModel) set(from *fabspark.InstancePool) {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Name = types.StringPointerValue(from.Name)
	to.Type = types.StringPointerValue((*string)(from.Type))
}

type starterPoolPropertiesModel struct {
	MaxNodeCount types.Int32 `tfsdk:"max_node_count"`
	MaxExecutors types.Int32 `tfsdk:"max_executors"`
}

func (to *starterPoolPropertiesModel) set(from *fabspark.StarterPoolProperties) {
	to.MaxNodeCount = types.Int32PointerValue(from.MaxNodeCount)
	to.MaxExecutors = types.Int32PointerValue(from.MaxExecutors)
}

type requestUpdateSparkWorkspaceSettings struct {
	fabspark.UpdateWorkspaceSparkSettingsRequest
}

func (to *requestUpdateSparkWorkspaceSettings) set(ctx context.Context, from resourceSparkWorkspaceSettingsModel) diag.Diagnostics { //nolint:gocognit, gocyclo
	if !from.AutomaticLog.IsNull() && !from.AutomaticLog.IsUnknown() {
		automaticLog, diags := from.AutomaticLog.Get(ctx)
		if diags.HasError() {
			return diags
		}

		if !automaticLog.Enabled.IsNull() && !automaticLog.Enabled.IsUnknown() {
			to.AutomaticLog = &fabspark.AutomaticLogProperties{
				Enabled: automaticLog.Enabled.ValueBoolPointer(),
			}
		}
	}

	if !from.Environment.IsNull() && !from.Environment.IsUnknown() {
		environment, diags := from.Environment.Get(ctx)
		if diags.HasError() {
			return diags
		}

		var reqEnvironment fabspark.EnvironmentProperties

		if !environment.Name.IsNull() && !environment.Name.IsUnknown() {
			reqEnvironment.Name = environment.Name.ValueStringPointer()
		}

		if !environment.RuntimeVersion.IsNull() && !environment.RuntimeVersion.IsUnknown() {
			reqEnvironment.RuntimeVersion = environment.RuntimeVersion.ValueStringPointer()
		}

		if reqEnvironment != (fabspark.EnvironmentProperties{}) {
			to.Environment = &reqEnvironment
		}
	}

	if !from.HighConcurrency.IsNull() && !from.HighConcurrency.IsUnknown() {
		highConcurrency, diags := from.HighConcurrency.Get(ctx)
		if diags.HasError() {
			return diags
		}

		if !highConcurrency.NotebookInteractiveRunEnabled.IsNull() && !highConcurrency.NotebookInteractiveRunEnabled.IsUnknown() {
			to.HighConcurrency = &fabspark.HighConcurrencyProperties{
				NotebookInteractiveRunEnabled: highConcurrency.NotebookInteractiveRunEnabled.ValueBoolPointer(),
			}
		}
	}

	if !from.Pool.IsNull() && !from.Pool.IsUnknown() { //nolint:nestif
		pool, diags := from.Pool.Get(ctx)
		if diags.HasError() {
			return diags
		}

		var reqPool fabspark.PoolProperties

		if !pool.CustomizeComputeEnabled.IsNull() && !pool.CustomizeComputeEnabled.IsUnknown() {
			reqPool.CustomizeComputeEnabled = pool.CustomizeComputeEnabled.ValueBoolPointer()
		}

		if !pool.DefaultPool.IsNull() && !pool.DefaultPool.IsUnknown() {
			defaultPool, diags := pool.DefaultPool.Get(ctx)
			if diags.HasError() {
				return diags
			}

			if !defaultPool.Name.IsNull() && !defaultPool.Name.IsUnknown() {
				reqPool.DefaultPool = &fabspark.InstancePool{
					Name: defaultPool.Name.ValueStringPointer(),
					Type: (*fabspark.CustomPoolType)(defaultPool.Type.ValueStringPointer()),
				}
			}
		}

		if !pool.StarterPool.IsNull() && !pool.StarterPool.IsUnknown() {
			starterPool, diags := pool.StarterPool.Get(ctx)
			if diags.HasError() {
				return diags
			}

			var reqStarterPool fabspark.StarterPoolProperties

			if !starterPool.MaxNodeCount.IsNull() && !starterPool.MaxNodeCount.IsUnknown() {
				reqStarterPool.MaxNodeCount = starterPool.MaxNodeCount.ValueInt32Pointer()
			}

			if !starterPool.MaxExecutors.IsNull() && !starterPool.MaxExecutors.IsUnknown() {
				reqStarterPool.MaxExecutors = starterPool.MaxExecutors.ValueInt32Pointer()
			}

			if reqStarterPool != (fabspark.StarterPoolProperties{}) {
				reqPool.StarterPool = &reqStarterPool
			}
		}

		if reqPool != (fabspark.PoolProperties{}) {
			to.Pool = &reqPool
		}
	}

	return nil
}
