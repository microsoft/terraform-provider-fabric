// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package itemjob

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts" //revive:disable-line:import-alias-naming
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
)

/*
BASE MODEL
*/

type baseItemJobModel struct {
	ID             customtypes.UUID                                         `tfsdk:"id"`
	WorkspaceID    customtypes.UUID                                         `tfsdk:"workspace_id"`
	ItemID         customtypes.UUID                                         `tfsdk:"item_id"`
	JobType        types.String                                             `tfsdk:"job_type"`
	InvokeType     types.String                                             `tfsdk:"invoke_type"`
	Status         types.String                                             `tfsdk:"status"`
	RootActivityID types.String                                             `tfsdk:"root_activity_id"`
	StartTimeUTC   types.String                                             `tfsdk:"start_time_utc"`
	EndTimeUTC     types.String                                             `tfsdk:"end_time_utc"`
	FailureReason  supertypes.SingleNestedObjectValueOf[failureReasonModel] `tfsdk:"failure_reason"`
}

type failureReasonModel struct {
	ErrorCode types.String `tfsdk:"error_code"`
	Message   types.String `tfsdk:"message"`
}

func (to *baseItemJobModel) set(ctx context.Context, workspaceID, itemID, jobType string, from fabcore.ItemJobInstance) diag.Diagnostics {
	to.ID = customtypes.NewUUIDPointerValue(from.ID)
	to.WorkspaceID = customtypes.NewUUIDValue(workspaceID)
	to.ItemID = customtypes.NewUUIDValue(itemID)
	to.JobType = types.StringValue(jobType)

	if from.JobType != nil {
		to.JobType = types.StringValue(*from.JobType)
	}

	to.InvokeType = types.StringPointerValue((*string)(from.InvokeType))
	to.Status = types.StringPointerValue((*string)(from.Status))
	to.RootActivityID = types.StringPointerValue(from.RootActivityID)
	to.StartTimeUTC = types.StringPointerValue(from.StartTimeUTC)
	to.EndTimeUTC = types.StringPointerValue(from.EndTimeUTC)

	failureReason := supertypes.NewSingleNestedObjectValueOfNull[failureReasonModel](ctx)

	if from.FailureReason != nil {
		failureReasonModel := &failureReasonModel{}
		failureReasonModel.set(*from.FailureReason)

		if diags := failureReason.Set(ctx, failureReasonModel); diags.HasError() {
			return diags
		}
	}

	to.FailureReason = failureReason

	return nil
}

func (to *failureReasonModel) set(from fabcore.ErrorResponse) {
	to.ErrorCode = types.StringPointerValue(from.ErrorCode)
	to.Message = types.StringPointerValue(from.Message)
}

/*
DATA-SOURCE
*/

type dataSourceItemJobModel struct {
	baseItemJobModel

	Timeouts timeoutsD.Value `tfsdk:"timeouts"`
}

/*
RESOURCE
*/

type resourceItemJobModel struct {
	baseItemJobModel

	ExecutionData jsontypes.Normalized `tfsdk:"execution_data"`
	Timeouts      timeoutsR.Value      `tfsdk:"timeouts"`
}

type requestRunOnDemandItemJob struct {
	fabcore.RunOnDemandItemJobRequest
}

func (to *requestRunOnDemandItemJob) set(from resourceItemJobModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if !from.ExecutionData.IsNull() && !from.ExecutionData.IsUnknown() {
		var executionData any
		if err := json.Unmarshal([]byte(from.ExecutionData.ValueString()), &executionData); err != nil {
			diags.AddError(
				"Invalid execution_data",
				"The execution_data attribute must be a valid JSON object: "+err.Error(),
			)

			return diags
		}

		to.ExecutionData = executionData
	}

	return diags
}
