package itemjobscheduler

import (
	"context"
	"fmt"
	"slices"
	"strings"

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
	_ resource.ResourceWithConfigure   = (*resourceItemJobScheduler)(nil)
	_ resource.ResourceWithImportState = (*resourceItemJobScheduler)(nil)
)

type resourceItemJobScheduler struct {
	pConfigData  *pconfig.ProviderData
	client       *fabcore.JobSchedulerClient
	fabricClient *fabcore.ItemsClient
	TypeInfo     tftypeinfo.TFTypeInfo
}

func NewResourceItemJobScheduler() resource.Resource {
	return &resourceItemJobScheduler{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceItemJobScheduler) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceItemJobScheduler) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema(false).GetResource(ctx)
}

func (r *resourceItemJobScheduler) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = (*fabcore.JobSchedulerClient)(fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewJobSchedulerClient())
	r.fabricClient = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewItemsClient()
}

func (r *resourceItemJobScheduler) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan, state resourceJobScheduleModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	state.Timeouts = plan.Timeouts

	var reqCreate requestCreateJobSchedule
	if resp.Diagnostics.Append(reqCreate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	respItem, err := r.fabricClient.GetItem(ctx, plan.WorkspaceID.ValueString(), plan.ItemID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	// Validate job type
	if err := r.validateJobType(respItem.Type, plan.JobType.ValueString()); err != nil {
		resp.Diagnostics.AddError("Invalid Job Type", err.Error())

		return
	}

	respCreate, err := r.client.CreateItemSchedule(ctx, plan.WorkspaceID.ValueString(), plan.ItemID.ValueString(), plan.JobType.ValueString(), reqCreate.CreateScheduleRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	state.set(ctx, plan.WorkspaceID.ValueString(), plan.ItemID.ValueString(), plan.JobType.ValueString(), respCreate.ItemSchedule)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceItemJobScheduler) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceJobScheduleModel

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

func (r *resourceItemJobScheduler) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	var plan resourceJobScheduleModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateJobSchedule

	if resp.Diagnostics.Append(reqUpdate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	respUpdate, err := r.client.UpdateItemSchedule(
		ctx,
		plan.WorkspaceID.ValueString(),
		plan.ItemID.ValueString(),
		plan.JobType.ValueString(),
		plan.ID.ValueString(),
		reqUpdate.UpdateScheduleRequest,
		nil,
	)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.set(ctx, plan.WorkspaceID.ValueString(), plan.ItemID.ValueString(), plan.JobType.ValueString(), respUpdate.ItemSchedule)...); resp.Diagnostics.HasError() {
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

func (r *resourceItemJobScheduler) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	var state resourceJobScheduleModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.DeleteItemSchedule(ctx, state.WorkspaceID.ValueString(), state.ItemID.ValueString(), state.JobType.ValueString(), state.ID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceItemJobScheduler) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
			fmt.Sprintf(common.ErrorImportIdentifierDetails, "WorkspaceID/ItemID/JobType/ScheduleID"),
		)

		return
	}

	workspaceID, itemID, jobType, scheduleID := parts[0], parts[1], parts[2], parts[3]

	uuidWorkspaceID, diags := customtypes.NewUUIDValueMust(workspaceID)
	resp.Diagnostics.Append(diags...)

	uuitemID, diags := customtypes.NewUUIDValueMust(itemID)
	resp.Diagnostics.Append(diags...)

	uuScheduleID, diags := customtypes.NewUUIDValueMust(scheduleID)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var timeout timeouts.Value
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeout)...); resp.Diagnostics.HasError() {
		return
	}

	state := resourceJobScheduleModel{
		baseItemJobSchedulerModel: baseItemJobSchedulerModel{
			ItemID:      uuitemID,
			WorkspaceID: uuidWorkspaceID,
			JobType:     types.StringValue(jobType),
			ID:          uuScheduleID,
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

// validateJobType validates that the job type is supported for the given item type.
func (r *resourceItemJobScheduler) validateJobType(itemType *fabcore.ItemType, jobType string) error {
	if itemType == nil {
		return fmt.Errorf("item type is nil")
	}

	itemTypeLowercase := strings.ToLower(string(*itemType))
	validJobTypes, exists := JobTypeActions[itemTypeLowercase]

	if !exists {
		return fmt.Errorf("item type '%s' does not support job scheduling. Supported types are: %v",
			*itemType, getMapKeys(JobTypeActions))
	}

	//	jobTypeLowercase := strings.ToLower(jobType)
	if !slices.Contains(validJobTypes, jobType) {
		return fmt.Errorf("job type '%s' is not valid for item type '%s'. Valid job types for '%s' are: %v",
			jobType, *itemType, *itemType, validJobTypes)
	}

	return nil
}

func (r *resourceItemJobScheduler) get(ctx context.Context, model *resourceJobScheduleModel) diag.Diagnostics {
	respGet, err := r.client.GetItemSchedule(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), model.JobType.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		return diags
	}

	return model.set(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), model.JobType.ValueString(), respGet.ItemSchedule)
}

func getMapKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
