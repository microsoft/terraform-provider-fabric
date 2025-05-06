// // Copyright (c) Microsoft Corporation
// // SPDX-License-Identifier: MPL-2.0

package onelakeshortcut

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure      = (*resourceOnelakeShortcut)(nil)
	_ resource.ResourceWithImportState    = (*resourceOnelakeShortcut)(nil)
	_ resource.ResourceWithValidateConfig = (*resourceOnelakeShortcut)(nil)
)

type resourceOnelakeShortcut struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.OneLakeShortcutsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func (r *resourceOnelakeShortcut) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "IMPORT", map[string]any{
		"id": req.ID,
	})

	_, diags := customtypes.NewUUIDValueMust(req.ID)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceOnelakeShortcut) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config resourceOneLakeShortcutModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		targetConfig *targetModel
		diags        diag.Diagnostics
	)
	if !config.Target.IsNull() && !config.Target.IsUnknown() {
		targetConfig, diags = config.Target.Get(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		resp.Diagnostics.AddAttributeError(
			path.Root("target"),
			"Missing target block",
			"You must specify exactly one of onelake, adls_gen2, amazon_s3, google_cloud_storage, s3_compatible, external_data_share or dataverse.",
		)
		return
	}

	count := 0

	if !targetConfig.Onelake.IsNull() && !targetConfig.Onelake.IsUnknown() {
		count++
	}
	if !targetConfig.AdlsGen2.IsNull() && !targetConfig.AdlsGen2.IsUnknown() {
		count++
	}
	if !targetConfig.AmazonS3.IsNull() && !targetConfig.AmazonS3.IsUnknown() {
		count++
	}
	if !targetConfig.GoogleCloudStorage.IsNull() && !targetConfig.GoogleCloudStorage.IsUnknown() {
		count++
	}
	if !targetConfig.S3Compatible.IsNull() && !targetConfig.S3Compatible.IsUnknown() {
		count++
	}
	if !targetConfig.ExternalDataShare.IsNull() && !targetConfig.ExternalDataShare.IsUnknown() {
		count++
	}
	if !targetConfig.Dataverse.IsNull() && !targetConfig.Dataverse.IsUnknown() {
		count++
	}

	if count != 1 {
		resp.Diagnostics.AddAttributeError(
			path.Root("target"),
			"Exactly one target type must be specified",
			fmt.Sprintf(
				"You have configured %d targets; please set exactly one of onelake, adls_gen2, amazon_s3, google_cloud_storage, s3_compatible, external_data_share or dataverse.",
				count,
			),
		)
	}
}

func NewResourceOneLakeShortcut() resource.Resource {
	return &resourceOnelakeShortcut{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceOnelakeShortcut) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceOnelakeShortcut) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema(false).GetResource(ctx)
}

func (r *resourceOnelakeShortcut) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewOneLakeShortcutsClient()
}

func (r *resourceOnelakeShortcut) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan, state resourceOneLakeShortcutModel

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

	var reqCreate requestCreateOnelakeShortcut

	if resp.Diagnostics.Append(reqCreate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	// options := fabcore.OneLakeShortcutsClientCreateShortcutOptions{
	// 	ShortcutConflictPolicy: (*fabcore.ShortcutConflictPolicy)(plan.ShortcutConflictPolicy.ValueStringPointer()),
	// }

	respCreate, err := r.client.CreateShortcut(ctx, plan.WorkspaceID.ValueString(), plan.ItemID.ValueString(), reqCreate.CreateShortcutRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}
	state.ID = types.StringValue(*respCreate.Shortcut.Name + *respCreate.Shortcut.Path)

	if resp.Diagnostics.Append(state.set(ctx, plan.WorkspaceID.ValueString(), plan.ItemID.ValueString(), respCreate.Shortcut)...); resp.Diagnostics.HasError() {
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

func (r *resourceOnelakeShortcut) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceOneLakeShortcutModel

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

func (r *resourceOnelakeShortcut) get(ctx context.Context, model *resourceOneLakeShortcutModel) diag.Diagnostics {
	respGet, err := r.client.GetShortcut(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), model.Path.ValueString(), model.Name.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		return diags
	}

	model.set(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), respGet.Shortcut)

	return nil
}

func (r *resourceOnelakeShortcut) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	var state resourceOneLakeShortcutModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.DeleteShortcut(ctx, state.WorkspaceID.ValueString(), state.ItemID.ValueString(), state.Path.ValueString(), state.Name.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceOnelakeShortcut) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{"action": "start"})

	var plan, state resourceOneLakeShortcutModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	state.Timeouts = plan.Timeouts

	var reqCreate requestCreateOnelakeShortcut

	if resp.Diagnostics.Append(reqCreate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	overwriteOnlyPolicy := fabcore.ShortcutConflictPolicyOverwriteOnly
	options := fabcore.OneLakeShortcutsClientCreateShortcutOptions{
		ShortcutConflictPolicy: &overwriteOnlyPolicy,
	}
	// Call the API to update the resource
	respCreate, err := r.client.CreateShortcut(ctx, plan.WorkspaceID.ValueString(), plan.ItemID.ValueString(), reqCreate.CreateShortcutRequest, &options)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(state.set(ctx, plan.WorkspaceID.ValueString(), plan.ItemID.ValueString(), respCreate.Shortcut)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "UPDATE", map[string]any{"action": "end"})
}
