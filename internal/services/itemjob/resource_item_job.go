// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package itemjob

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var (
	_ resource.ResourceWithConfigure   = (*resourceItemJob)(nil)
	_ resource.ResourceWithImportState = (*resourceItemJob)(nil)
)

type resourceItemJob struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.JobSchedulerClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewResourceItemJob() resource.Resource {
	return &resourceItemJob{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceItemJob) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceItemJob) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema().GetResource(ctx)
}

func (r *resourceItemJob) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	pConfigData, ok := req.ProviderData.(*pconfig.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			common.ErrorResourceConfigType,
			fmt.Sprintf(common.ErrorFabricClientType, req.ProviderData),
		)

		return
	}

	r.pConfigData = pConfigData

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(r.TypeInfo.Name, r.TypeInfo.IsPreview, r.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}

	r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewJobSchedulerClient()
}

func (r *resourceItemJob) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan, state resourceItemJobModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqRun requestRunOnDemandItemJob
	if resp.Diagnostics.Append(reqRun.set(plan)...); resp.Diagnostics.HasError() {
		return
	}

	var rawResp *http.Response
	ctxCapture := policy.WithCaptureResponse(ctx, &rawResp)

	_, err := r.client.RunOnDemandItemJob(
		ctxCapture,
		plan.WorkspaceID.ValueString(),
		plan.ItemID.ValueString(),
		plan.JobType.ValueString(),
		&fabcore.JobSchedulerClientRunOnDemandItemJobOptions{
			RunOnDemandItemJobRequest: &reqRun.RunOnDemandItemJobRequest,
		},
	)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	jobInstanceID, diags := getJobInstanceIDFromResponse(rawResp)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	state.ExecutionData = plan.ExecutionData
	state.Timeouts = plan.Timeouts
	state.ID = customtypes.NewUUIDValue(jobInstanceID)
	state.WorkspaceID = plan.WorkspaceID
	state.ItemID = plan.ItemID
	state.JobType = plan.JobType

	if resp.Diagnostics.Append(r.get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceItemJob) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceItemJobModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Read(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	diags = r.get(ctx, &state)
	if utils.IsErrNotFound(state.ID.ValueString(), &diags, fabcore.ErrCommon.EntityNotFound) {
		resp.State.RemoveResource(ctx)

		resp.Diagnostics.Append(diags...)

		return
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceItemJob) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	var plan resourceItemJobModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(r.get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceItemJob) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	var state resourceItemJobModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// A job instance cannot be deleted. Best effort: cancel the job if it is still running.
	// Any error (e.g. the job already reached a terminal state, or was not found) is ignored
	// so that the resource can always be removed from the state.
	if _, err := r.client.CancelItemJobInstance(ctx, state.WorkspaceID.ValueString(), state.ItemID.ValueString(), state.ID.ValueString(), nil); err != nil {
		tflog.Debug(ctx, "DELETE", map[string]any{
			"action":  "cancel job instance failed (ignored)",
			"message": err.Error(),
		})
	}

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceItemJob) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "IMPORT", map[string]any{
		"id": req.ID,
	})

	parts := strings.Split(req.ID, "/")
	if len(parts) != 4 {
		resp.Diagnostics.AddError(
			common.ErrorImportIdentifierHeader,
			fmt.Sprintf(common.ErrorImportIdentifierDetails, "WorkspaceID/ItemID/JobType/JobInstanceID"),
		)

		return
	}

	workspaceID, itemID, jobType, jobInstanceID := parts[0], parts[1], parts[2], parts[3]

	uuidWorkspaceID, diags := customtypes.NewUUIDValueMust(workspaceID)
	resp.Diagnostics.Append(diags...)

	uuidItemID, diags := customtypes.NewUUIDValueMust(itemID)
	resp.Diagnostics.Append(diags...)

	uuidJobInstanceID, diags := customtypes.NewUUIDValueMust(jobInstanceID)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var timeout timeouts.Value
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeout)...); resp.Diagnostics.HasError() {
		return
	}

	state := resourceItemJobModel{
		baseItemJobModel: baseItemJobModel{
			ID:          uuidJobInstanceID,
			WorkspaceID: uuidWorkspaceID,
			ItemID:      uuidItemID,
			JobType:     types.StringValue(jobType),
		},
		Timeouts: timeout,
	}

	if resp.Diagnostics.Append(r.get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceItemJob) get(ctx context.Context, model *resourceItemJobModel) diag.Diagnostics {
	respGet, err := r.client.GetItemJobInstance(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		return diags
	}

	return model.set(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), model.JobType.ValueString(), respGet.ItemJobInstance)
}

// getJobInstanceIDFromResponse extracts the created job instance ID from the Location header of the
// RunOnDemandItemJob response. The Location header points to the created job instance, with the job
// instance ID as its last path segment.
func getJobInstanceIDFromResponse(rawResp *http.Response) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if rawResp == nil {
		diags.AddError(
			"Missing job instance response",
			"The response from the on-demand job trigger was empty.",
		)

		return "", diags
	}

	location := rawResp.Header.Get("Location")
	if location == "" {
		diags.AddError(
			"Missing job instance location",
			"The response from the on-demand job trigger did not contain a Location header with the job instance ID.",
		)

		return "", diags
	}

	// Trim any query string or trailing slash before extracting the last path segment.
	if idx := strings.IndexAny(location, "?#"); idx != -1 {
		location = location[:idx]
	}

	location = strings.TrimRight(location, "/")

	jobInstanceID := location[strings.LastIndex(location, "/")+1:]
	if jobInstanceID == "" {
		diags.AddError(
			"Invalid job instance location",
			fmt.Sprintf("Could not determine the job instance ID from the Location header: %q.", rawResp.Header.Get("Location")),
		)

		return "", diags
	}

	return jobInstanceID, diags
}
