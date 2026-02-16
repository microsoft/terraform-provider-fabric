// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sparkenvsettings

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure        = (*resourceSparkEnvironmentSettings)(nil)
	_ resource.ResourceWithConfigValidators = (*resourceSparkEnvironmentSettings)(nil)
)

type resourceSparkEnvironmentSettings struct {
	pConfigData     *pconfig.ProviderData
	publishedClient *fabenvironment.PublishedClient
	stagingClient   *fabenvironment.StagingClient
	itemsClient     *fabenvironment.ItemsClient
	TypeInfo        tftypeinfo.TFTypeInfo
}

func NewResourceSparkEnvironmentSettings() resource.Resource {
	return &resourceSparkEnvironmentSettings{
		TypeInfo: ItemTypeInfo,
	}
}

func (r *resourceSparkEnvironmentSettings) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeInfo.FullTypeName(false)
}

func (r *resourceSparkEnvironmentSettings) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = itemSchema().GetResource(ctx)
}

func (r *resourceSparkEnvironmentSettings) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.publishedClient = fabenvironment.NewClientFactoryWithClient(*pConfigData.FabricClient).NewPublishedClient()
	r.stagingClient = fabenvironment.NewClientFactoryWithClient(*pConfigData.FabricClient).NewStagingClient()
	r.itemsClient = fabenvironment.NewClientFactoryWithClient(*pConfigData.FabricClient).NewItemsClient()
}

func (r *resourceSparkEnvironmentSettings) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("driver_cores"),
			path.MatchRoot("driver_memory"),
			path.MatchRoot("dynamic_executor_allocation"),
			path.MatchRoot("executor_cores"),
			path.MatchRoot("executor_memory"),
			path.MatchRoot("pool"),
			path.MatchRoot("runtime_version"),
			path.MatchRoot("spark_properties"),
		),
	}
}

func (r *resourceSparkEnvironmentSettings) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})

	var plan resourceSparkEnvironmentSettingsModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqCreate requestUpdateSparkEnvironmentSettings

	if resp.Diagnostics.Append(reqCreate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	{
		respCurrent, err := r.stagingClient.GetSparkCompute(ctx, plan.WorkspaceID.ValueString(), plan.EnvironmentID.ValueString(), false, nil)
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
			return
		}

		if reqCreate.SparkProperties != nil || len(respCurrent.SparkCompute.SparkProperties) > 0 {
			reqCreate.SparkProperties = diffSparkProperties(reqCreate.SparkProperties, respCurrent.SparkCompute.SparkProperties)
		}
	}

	respCreate, err := r.stagingClient.UpdateSparkCompute(
		ctx,
		plan.WorkspaceID.ValueString(),
		plan.EnvironmentID.ValueString(),
		false,
		reqCreate.UpdateEnvironmentSparkComputeRequest,
		nil,
	)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.set(ctx, respCreate.SparkCompute)...); resp.Diagnostics.HasError() {
		return
	}

	plan.ID = plan.EnvironmentID

	if resp.Diagnostics.Append(r.publish(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSparkEnvironmentSettings) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var state resourceSparkEnvironmentSettingsModel

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
	if utils.IsErrNotFound(state.WorkspaceID.ValueString(), &diags, fabcore.ErrCommon.EntityNotFound) {
		resp.State.RemoveResource(ctx)

		resp.Diagnostics.Append(diags...)

		return
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	state.ID = state.EnvironmentID

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSparkEnvironmentSettings) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})

	var plan resourceSparkEnvironmentSettingsModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateSparkEnvironmentSettings

	if resp.Diagnostics.Append(reqUpdate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	{
		respCurrent, err := r.stagingClient.GetSparkCompute(ctx, plan.WorkspaceID.ValueString(), plan.EnvironmentID.ValueString(), false, nil)
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
			return
		}

		if reqUpdate.SparkProperties != nil || len(respCurrent.SparkCompute.SparkProperties) > 0 {
			reqUpdate.SparkProperties = diffSparkProperties(reqUpdate.SparkProperties, respCurrent.SparkCompute.SparkProperties)
		}
	}

	respUpdate, err := r.stagingClient.UpdateSparkCompute(
		ctx,
		plan.WorkspaceID.ValueString(),
		plan.EnvironmentID.ValueString(),
		false,
		reqUpdate.UpdateEnvironmentSparkComputeRequest,
		nil,
	)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.set(ctx, respUpdate.SparkCompute)...); resp.Diagnostics.HasError() {
		return
	}

	plan.ID = plan.EnvironmentID

	if resp.Diagnostics.Append(r.publish(ctx, plan)...); resp.Diagnostics.HasError() {
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

func (r *resourceSparkEnvironmentSettings) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})

	resp.Diagnostics.AddWarning(
		"delete operation not supported",
		fmt.Sprintf(
			"Resource %s does not support deletion. It will be removed from Terraform state, but no action will be taken in the Fabric. All current settings will remain.",
			r.TypeInfo.Name,
		),
	)

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceSparkEnvironmentSettings) get(ctx context.Context, model *resourceSparkEnvironmentSettingsModel) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s for Workspace ID: %s", r.TypeInfo.Name, model.WorkspaceID.ValueString()))

	var respEntity fabenvironment.SparkCompute

	if model.PublicationStatus.ValueString() == SparkEnvironmentPublicationStatusPublished {
		respGet, err := r.publishedClient.GetSparkCompute(ctx, model.WorkspaceID.ValueString(), model.EnvironmentID.ValueString(), false, nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
			return diags
		}

		respEntity = respGet.SparkCompute
	} else {
		respGet, err := r.stagingClient.GetSparkCompute(ctx, model.WorkspaceID.ValueString(), model.EnvironmentID.ValueString(), false, nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
			return diags
		}

		respEntity = respGet.SparkCompute
	}

	return model.set(ctx, respEntity)
}

func (r *resourceSparkEnvironmentSettings) publish(ctx context.Context, model resourceSparkEnvironmentSettingsModel) diag.Diagnostics {
	if model.PublicationStatus.ValueString() == SparkEnvironmentPublicationStatusPublished {
		respPublish, err := r.itemsClient.PublishEnvironment(ctx, model.WorkspaceID.ValueString(), model.EnvironmentID.ValueString(), false, nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil); diags.HasError() {
			return diags
		}

		if respPublish.Properties.PublishDetails != nil && respPublish.Properties.PublishDetails.State != nil {
			state := string(*respPublish.Properties.PublishDetails.State)

			if strings.EqualFold(state, string(fabenvironment.PublishStateFailed)) {
				var diags diag.Diagnostics
				diags.AddError(
					"publishing failed",
					fmt.Sprintf("Environment publishing failed for Environment ID: %s in Workspace ID: %s", model.EnvironmentID.ValueString(), model.WorkspaceID.ValueString()),
				)

				return diags
			}

			if strings.EqualFold(state, string(fabenvironment.PublishStateCancelled)) {
				var diags diag.Diagnostics
				diags.AddError(
					"publishing cancelled",
					fmt.Sprintf("Environment publishing cancelled for Environment ID: %s in Workspace ID: %s", model.EnvironmentID.ValueString(), model.WorkspaceID.ValueString()),
				)

				return diags
			}
		}
	}

	return nil
}
