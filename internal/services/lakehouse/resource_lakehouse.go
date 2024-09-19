// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	"context"
	"fmt"
	"strings"
	"time"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fablakehouse "github.com/microsoft/fabric-sdk-go/fabric/lakehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure   = (*resourceLakehouse)(nil)
	_ resource.ResourceWithImportState = (*resourceLakehouse)(nil)
)

type resourceLakehouse struct {
	pConfigData *pconfig.ProviderData
	client      *fablakehouse.ItemsClient
}

func NewResourceLakehouse() resource.Resource {
	return &resourceLakehouse{}
}

func (r *resourceLakehouse) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ItemTFName
}

func (r *resourceLakehouse) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	markdownDescription := "This resource manages a Fabric " + ItemName + ".\n\n" +
		"See [" + ItemName + "](" + ItemDocsURL + ") for more information.\n\n" +
		ItemDocsSPNSupport

	properties := schema.SingleNestedAttribute{
		MarkdownDescription: "The " + ItemName + " properties.",
		Computed:            true,
		CustomType:          supertypes.NewSingleNestedObjectTypeOf[lakehousePropertiesModel](ctx),
		Attributes: map[string]schema.Attribute{
			"onelake_files_path": schema.StringAttribute{
				MarkdownDescription: "OneLake path to the Lakehouse files directory",
				Computed:            true,
			},
			"onelake_tables_path": schema.StringAttribute{
				MarkdownDescription: "OneLake path to the Lakehouse tables directory.",
				Computed:            true,
			},
			"sql_endpoint_properties": schema.SingleNestedAttribute{
				MarkdownDescription: "An object containing the properties of the SQL endpoint.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[lakehouseSQLEndpointPropertiesModel](ctx),
				Attributes: map[string]schema.Attribute{
					"provisioning_status": schema.StringAttribute{
						MarkdownDescription: "The SQL endpoint provisioning status.",
						Computed:            true,
					},
					"connection_string": schema.StringAttribute{
						MarkdownDescription: "SQL endpoint connection string.",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "SQL endpoint ID.",
						Computed:            true,
						CustomType:          customtypes.UUIDType{},
					},
				},
			},
			"default_schema": schema.StringAttribute{
				MarkdownDescription: "Default schema of the Lakehouse. This property is returned only for schema enabled Lakehouse.",
				Computed:            true,
			},
		},
	}
	configuration := schema.SingleNestedAttribute{
		MarkdownDescription: "The " + ItemName + " creation configuration.\n\n" +
			"Any changes to this configuration will result in recreation of the " + ItemName + ".",
		Optional:   true,
		CustomType: supertypes.NewSingleNestedObjectTypeOf[lakehouseConfigurationModel](ctx),
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
		Attributes: map[string]schema.Attribute{
			"enable_schemas": schema.BoolAttribute{
				MarkdownDescription: "Schema enabled Lakehouse.",
				Required:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}

	resp.Schema = fabricitem.GetResourceFabricItemPropertiesCreationSchema(ctx, ItemName, markdownDescription, 123, 256, true, properties, configuration)
}

func (r *resourceLakehouse) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = fablakehouse.NewClientFactoryWithClient(*pConfigData.FabricClient).NewItemsClient()
}

func (r *resourceLakehouse) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "CREATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
	})

	var plan resourceLakehouseModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqCreate requestCreateLakehouse

	if resp.Diagnostics.Append(reqCreate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	respCreate, err := r.client.CreateLakehouse(ctx, plan.WorkspaceID.ValueString(), reqCreate.CreateLakehouseRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.set(ctx, respCreate.Lakehouse)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	if resp.Diagnostics.Append(r.get(ctx, &plan)...); resp.Diagnostics.HasError() {
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

func (r *resourceLakehouse) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"state": req.State,
	})

	var state resourceLakehouseModel

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

func (r *resourceLakehouse) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "UPDATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
		"state":  req.State,
	})

	var plan resourceLakehouseModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateLakehouse

	reqUpdate.set(plan)

	respUpdate, err := r.client.UpdateLakehouse(ctx, plan.WorkspaceID.ValueString(), plan.ID.ValueString(), reqUpdate.UpdateLakehouseRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(plan.set(ctx, respUpdate.Lakehouse)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

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

func (r *resourceLakehouse) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "DELETE", map[string]any{
		"state": req.State,
	})

	var state resourceLakehouseModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.DeleteLakehouse(ctx, state.WorkspaceID.ValueString(), state.ID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceLakehouse) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "IMPORT", map[string]any{
		"id": req.ID,
	})

	workspaceID, lakehouseID, found := strings.Cut(req.ID, "/")

	if !found {
		resp.Diagnostics.AddError(
			common.ErrorImportIdentifierHeader,
			fmt.Sprintf(common.ErrorImportIdentifierDetails, "WorkspaceID/LakehouseID"),
		)

		return
	}

	uuidWorkspaceID, diags := customtypes.NewUUIDValueMust(workspaceID)
	resp.Diagnostics.Append(diags...)

	uuidID, diags := customtypes.NewUUIDValueMust(lakehouseID)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var configuration supertypes.SingleNestedObjectValueOf[lakehouseConfigurationModel]
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("configuration"), &configuration)...); resp.Diagnostics.HasError() {
		return
	}

	var timeout timeouts.Value
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeout)...); resp.Diagnostics.HasError() {
		return
	}

	state := resourceLakehouseModel{}
	state.ID = uuidID
	state.WorkspaceID = uuidWorkspaceID
	state.Configuration = configuration
	state.Timeouts = timeout

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

func (r *resourceLakehouse) get(ctx context.Context, model *resourceLakehouseModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET", map[string]any{
		"workspace_id": model.WorkspaceID.ValueString(),
		"id":           model.ID.ValueString(),
	})

	for {
		respGet, err := r.client.GetLakehouse(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
			return diags
		}

		if respGet.Properties == nil || respGet.Properties.SQLEndpointProperties == nil {
			tflog.Info(ctx, "Lakehouse provisioning not done, waiting 30 seconds before retrying")
			time.Sleep(30 * time.Second) // lintignore:R018

			continue
		}

		switch *respGet.Properties.SQLEndpointProperties.ProvisioningStatus {
		case fablakehouse.SQLEndpointProvisioningStatusFailed:
			var diags diag.Diagnostics

			diags.AddError(
				"provisioning failed",
				"Lakehouse SQL endpoint provisioning failed")

			return diags

		case fablakehouse.SQLEndpointProvisioningStatusSuccess:
			return model.set(ctx, respGet.Lakehouse)
		default:
			tflog.Info(ctx, "Lakehouse provisioning in progress, waiting 30 seconds before retrying")
			time.Sleep(30 * time.Second) // lintignore:R018
		}
	}
}
