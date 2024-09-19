// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package spark

import (
	"context"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	timeoutsd "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsr "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabspark "github.com/microsoft/fabric-sdk-go/fabric/spark"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

type dataSourceSparkCustomPoolModel struct {
	baseSparkCustomPoolModel
	Timeouts timeoutsd.Value `tfsdk:"timeouts"`
}

type resourceSparkCustomPoolModel struct {
	baseSparkCustomPoolModel
	Timeouts timeoutsr.Value `tfsdk:"timeouts"`
}

type baseSparkCustomPoolModel struct {
	ID                        customtypes.UUID                                                                    `tfsdk:"id"`
	WorkspaceID               customtypes.UUID                                                                    `tfsdk:"workspace_id"`
	Type                      types.String                                                                        `tfsdk:"type"`
	Name                      types.String                                                                        `tfsdk:"name"`
	NodeFamily                types.String                                                                        `tfsdk:"node_family"`
	NodeSize                  types.String                                                                        `tfsdk:"node_size"`
	AutoScale                 supertypes.SingleNestedObjectValueOf[sparkCustomPoolAutoScaleModel]                 `tfsdk:"auto_scale"`
	DynamicExecutorAllocation supertypes.SingleNestedObjectValueOf[sparkCustomPoolDynamicExecutorAllocationModel] `tfsdk:"dynamic_executor_allocation"`
}

func (to *baseSparkCustomPoolModel) set(ctx context.Context, from fabspark.CustomPool) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.Type = types.StringPointerValue((*string)(from.Type))
	to.Name = types.StringPointerValue(from.Name)
	to.NodeFamily = types.StringPointerValue((*string)(from.NodeFamily))
	to.NodeSize = types.StringPointerValue((*string)(from.NodeSize))

	autoScale := supertypes.NewSingleNestedObjectValueOfNull[sparkCustomPoolAutoScaleModel](ctx)

	if from.AutoScale != nil {
		autoScaleModel := &sparkCustomPoolAutoScaleModel{}
		autoScaleModel.set(from.AutoScale)

		if diags := autoScale.Set(ctx, autoScaleModel); diags.HasError() {
			return diags
		}
	}

	to.AutoScale = autoScale

	dynamicExecutorAllocation := supertypes.NewSingleNestedObjectValueOfNull[sparkCustomPoolDynamicExecutorAllocationModel](ctx)

	if from.DynamicExecutorAllocation != nil {
		dynamicExecutorAllocationModel := &sparkCustomPoolDynamicExecutorAllocationModel{}
		dynamicExecutorAllocationModel.set(from.DynamicExecutorAllocation)

		if diags := dynamicExecutorAllocation.Set(ctx, dynamicExecutorAllocationModel); diags.HasError() {
			return diags
		}
	}

	to.DynamicExecutorAllocation = dynamicExecutorAllocation

	return nil
}

type sparkCustomPoolAutoScaleModel struct {
	Enabled      types.Bool  `tfsdk:"enabled"`
	MinNodeCount types.Int32 `tfsdk:"min_node_count"`
	MaxNodeCount types.Int32 `tfsdk:"max_node_count"`
}

func (to *sparkCustomPoolAutoScaleModel) set(from *fabspark.AutoScaleProperties) {
	to.Enabled = types.BoolPointerValue(from.Enabled)
	to.MinNodeCount = types.Int32PointerValue(from.MinNodeCount)
	to.MaxNodeCount = types.Int32PointerValue(from.MaxNodeCount)
}

type sparkCustomPoolDynamicExecutorAllocationModel struct {
	Enabled      types.Bool  `tfsdk:"enabled"`
	MinExecutors types.Int32 `tfsdk:"min_executors"`
	MaxExecutors types.Int32 `tfsdk:"max_executors"`
}

func (to *sparkCustomPoolDynamicExecutorAllocationModel) set(from *fabspark.DynamicExecutorAllocationProperties) {
	to.Enabled = types.BoolPointerValue(from.Enabled)
	to.MinExecutors = types.Int32PointerValue(from.MinExecutors)
	to.MaxExecutors = types.Int32PointerValue(from.MaxExecutors)
}

type requestCreateSparkCustomPool struct {
	fabspark.CreateCustomPoolRequest
}

func (to *requestCreateSparkCustomPool) set(ctx context.Context, from resourceSparkCustomPoolModel) diag.Diagnostics {
	to.Name = from.Name.ValueStringPointer()
	to.NodeFamily = (*fabspark.NodeFamily)(from.NodeFamily.ValueStringPointer())
	to.NodeSize = (*fabspark.NodeSize)(from.NodeSize.ValueStringPointer())

	sparkCustomPoolAutoScaleModel, diags := from.AutoScale.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.AutoScale = &fabspark.AutoScaleProperties{
		Enabled:      sparkCustomPoolAutoScaleModel.Enabled.ValueBoolPointer(),
		MinNodeCount: sparkCustomPoolAutoScaleModel.MinNodeCount.ValueInt32Pointer(),
		MaxNodeCount: sparkCustomPoolAutoScaleModel.MaxNodeCount.ValueInt32Pointer(),
	}

	dynamicExecutorAllocationModel, diags := from.DynamicExecutorAllocation.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.DynamicExecutorAllocation = &fabspark.DynamicExecutorAllocationProperties{
		Enabled:      dynamicExecutorAllocationModel.Enabled.ValueBoolPointer(),
		MinExecutors: dynamicExecutorAllocationModel.MinExecutors.ValueInt32Pointer(),
		MaxExecutors: dynamicExecutorAllocationModel.MaxExecutors.ValueInt32Pointer(),
	}

	return nil
}

type requestUpdateSparkCustomPool struct {
	fabspark.UpdateCustomPoolRequest
}

func (to *requestUpdateSparkCustomPool) set(ctx context.Context, from resourceSparkCustomPoolModel) diag.Diagnostics {
	to.Name = from.Name.ValueStringPointer()
	to.NodeFamily = (*fabspark.NodeFamily)(from.NodeFamily.ValueStringPointer())
	to.NodeSize = (*fabspark.NodeSize)(from.NodeSize.ValueStringPointer())

	sparkCustomPoolAutoScaleModel, diags := from.AutoScale.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.AutoScale = &fabspark.AutoScaleProperties{
		Enabled:      sparkCustomPoolAutoScaleModel.Enabled.ValueBoolPointer(),
		MinNodeCount: sparkCustomPoolAutoScaleModel.MinNodeCount.ValueInt32Pointer(),
		MaxNodeCount: sparkCustomPoolAutoScaleModel.MaxNodeCount.ValueInt32Pointer(),
	}

	dynamicExecutorAllocationModel, diags := from.DynamicExecutorAllocation.Get(ctx)
	if diags.HasError() {
		return diags
	}

	to.DynamicExecutorAllocation = &fabspark.DynamicExecutorAllocationProperties{
		Enabled:      dynamicExecutorAllocationModel.Enabled.ValueBoolPointer(),
		MinExecutors: dynamicExecutorAllocationModel.MinExecutors.ValueInt32Pointer(),
		MaxExecutors: dynamicExecutorAllocationModel.MaxExecutors.ValueInt32Pointer(),
	}

	return nil
}
