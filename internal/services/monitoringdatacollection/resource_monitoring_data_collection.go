// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package monitoringdatacollection

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabmonitoring "github.com/microsoft/fabric-sdk-go/fabric/monitoring"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var (
	_ resource.ResourceWithConfigure   = (*resourceMonitoringDataCollection)(nil)
	_ resource.ResourceWithImportState = (*resourceMonitoringDataCollection)(nil)
)

type resourceMonitoringDataCollection struct {
	pConfigData *pconfig.ProviderData
	client      *fabmonitoring.DataCollectionClient
	itemsClient *fabmonitoring.ItemsClient
}

func NewResourceMonitoringDataCollection() resource.Resource {
	return &resourceMonitoringDataCollection{}
}

func (r *resourceMonitoringDataCollection) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ItemTypeInfo.FullTypeName(false)
}

func (r *resourceMonitoringDataCollection) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema().GetResource(ctx)
}

func (r *resourceMonitoringDataCollection) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	pConfigData, ok := req.ProviderData.(*pconfig.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(common.ErrorResourceConfigType, fmt.Sprintf(common.ErrorFabricClientType, req.ProviderData))

		return
	}

	r.pConfigData = pConfigData
	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(ItemTypeInfo.Name, ItemTypeInfo.IsPreview, pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}

	factory := fabmonitoring.NewClientFactoryWithClient(*pConfigData.FabricClient)
	r.client = factory.NewDataCollectionClient()
	r.itemsClient = factory.NewItemsClient()
}

func (r *resourceMonitoringDataCollection) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceMonitoringDataCollectionModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(r.setEnabled(ctx, plan, utils.OperationCreate)...); resp.Diagnostics.HasError() {
		return
	}

	// Track the successful mutation even if the subsequent read fails.
	plan.ID = plan.MonitoringID
	if resp.Diagnostics.Append(resp.State.Set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(r.get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *resourceMonitoringDataCollection) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceMonitoringDataCollectionModel
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
	if utils.IsErrNotFound(state.MonitoringID.ValueString(), &diags, fabcore.ErrCommon.EntityNotFound) {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.Append(diags...)

		return
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceMonitoringDataCollection) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resourceMonitoringDataCollectionModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(r.setEnabled(ctx, plan, utils.OperationUpdate)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(r.get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *resourceMonitoringDataCollection) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceMonitoringDataCollectionModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	tflog.Debug(ctx, "Disabling Monitoring data collection before removing resource from state")
	_, err := r.client.Disable(ctx, state.WorkspaceID.ValueString(), state.MonitoringID.ValueString(), nil)
	diags = utils.GetDiagsFromError(ctx, err, utils.OperationDelete, fabcore.ErrCommon.EntityNotFound)
	utils.IsErrNotFound(state.MonitoringID.ValueString(), &diags, fabcore.ErrCommon.EntityNotFound)
	resp.Diagnostics.Append(diags...)
}

func (r *resourceMonitoringDataCollection) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	workspaceID, monitoringID, found := strings.Cut(req.ID, "/")
	if !found {
		resp.Diagnostics.AddError(common.ErrorImportIdentifierHeader, fmt.Sprintf(common.ErrorImportIdentifierDetails, "WorkspaceID/MonitoringID"))

		return
	}

	workspaceUUID, diags := customtypes.NewUUIDValueMust(workspaceID)
	resp.Diagnostics.Append(diags...)
	monitoringUUID, diags := customtypes.NewUUIDValueMust(monitoringID)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_id"), workspaceUUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("monitoring_id"), monitoringUUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), monitoringUUID)...)
}

func (r *resourceMonitoringDataCollection) setEnabled(ctx context.Context, model resourceMonitoringDataCollectionModel, operation utils.Operation) diag.Diagnostics {
	if model.Enabled.IsNull() || model.Enabled.IsUnknown() {
		var diags diag.Diagnostics
		diags.AddAttributeError(path.Root("enabled"), "Invalid collection state", "The enabled value must be known before changing data collection.")

		return diags
	}

	tflog.Debug(ctx, "Updating Monitoring data collection", map[string]any{"enabled": model.Enabled.ValueBool()})

	var err error
	if model.Enabled.ValueBool() {
		_, err = r.client.Enable(ctx, model.WorkspaceID.ValueString(), model.MonitoringID.ValueString(), nil)
	} else {
		_, err = r.client.Disable(ctx, model.WorkspaceID.ValueString(), model.MonitoringID.ValueString(), nil)
	}

	return utils.GetDiagsFromError(ctx, err, operation, nil)
}

func (r *resourceMonitoringDataCollection) get(ctx context.Context, model *resourceMonitoringDataCollectionModel) diag.Diagnostics {
	response, err := r.itemsClient.GetMonitoring(ctx, model.WorkspaceID.ValueString(), model.MonitoringID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		return diags
	}

	var diags diag.Diagnostics
	if response.Properties == nil || response.Properties.DataCollectionState == nil {
		diags.AddError(common.ErrorReadHeader, "The Monitoring response did not contain properties.dataCollectionState.")

		return diags
	}

	switch *response.Properties.DataCollectionState {
	case fabmonitoring.DataCollectionStateEnabled:
		model.Enabled = types.BoolValue(true)
	case fabmonitoring.DataCollectionStateDisabled:
		model.Enabled = types.BoolValue(false)
	default:
		diags.AddError(common.ErrorReadHeader, fmt.Sprintf("Unsupported Monitoring data collection state %q.", *response.Properties.DataCollectionState))

		return diags
	}

	model.ID = model.MonitoringID

	return nil
}
