// // Copyright (c) Microsoft Corporation
// // SPDX-License-Identifier: MPL-2.0

package onelakeshortcut

// import (
// 	"context"
// 	"fmt"
// 	"reflect"
// 	"time"

// 	"github.com/hashicorp/terraform-plugin-framework/diag"
// 	"github.com/hashicorp/terraform-plugin-framework/path"
// 	"github.com/hashicorp/terraform-plugin-framework/resource"
// 	"github.com/hashicorp/terraform-plugin-log/tflog"
// 	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

// 	"github.com/microsoft/terraform-provider-fabric/internal/common"
// 	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
// 	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
// 	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
// 	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
// )

// // Ensure the implementation satisfies the expected interfaces.
// var (
// 	_ resource.ResourceWithConfigure      = (*resourceOnelakeShortcut)(nil)
// 	_ resource.ResourceWithValidateConfig = (*resourceOnelakeShortcut)(nil)
// 	_ resource.ResourceWithImportState    = (*resourceOnelakeShortcut)(nil)
// )

// type resourceOnelakeShortcut struct {
// 	pConfigData *pconfig.ProviderData
// 	client      *fabcore.OneLakeShortcutsClient
// 	TypeInfo    tftypeinfo.TFTypeInfo
// }

// func NewResourceWorkspace() resource.Resource {
// 	return &resourceOnelakeShortcut{
// 		TypeInfo: ItemTypeInfo,
// 	}
// }

// func (r *resourceOnelakeShortcut) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
// 	resp.TypeName = r.TypeInfo.FullTypeName(false)
// }

// func (r *resourceOnelakeShortcut) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
// 	resp.Schema = itemSchema(false).GetResource(ctx)
// }

// func (r *resourceOnelakeShortcut) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
// 	if req.ProviderData == nil {
// 		return
// 	}

// 	pConfigData, ok := req.ProviderData.(*pconfig.ProviderData)
// 	if !ok {
// 		resp.Diagnostics.AddError(
// 			common.ErrorResourceConfigType,
// 			fmt.Sprintf(common.ErrorFabricClientType, req.ProviderData),
// 		)

// 		return
// 	}

// 	r.pConfigData = pConfigData
// 	r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewOneLakeShortcutsClient()
// }

// func (r *resourceOnelakeShortcut) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
// 	var config resourceOneLakeShortcutModel
// 	var diags diag.Diagnostics

// 	if resp.Diagnostics.Append(req.Config.Get(ctx, &config)...); resp.Diagnostics.HasError() {
// 		return
// 	}

// 	return
// }

// func (r *resourceOnelakeShortcut) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
// 	tflog.Debug(ctx, "CREATE", map[string]any{
// 		"action": "start",
// 	})

// 	var plan, state resourceOneLakeShortcutModel

// 	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
// 		return
// 	}

// 	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
// 	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(ctx, timeout)
// 	defer cancel()

// 	state.Timeouts = plan.Timeouts

// 	var reqCreate requestCreateOnelakeShortcut

// 	reqCreate.set(plan)

// 	respCreate, err := r.client.CreateWorkspace(ctx, reqCreate.CreateWorkspaceRequest, nil)
// 	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
// 		return
// 	}

// 	plan.ID = customtypes.NewUUIDPointerValue(respCreate.ID)
// 	state.ID = plan.ID

// 	if resp.Diagnostics.Append(r.get(ctx, &state)...); resp.Diagnostics.HasError() {
// 		return
// 	}

// 	if resp.Diagnostics.Append(resp.State.Set(ctx, state)...); resp.Diagnostics.HasError() {
// 		return
// 	}

// 	identityPlan := &workspaceIdentityModel{}

// 	if !plan.Identity.IsNull() && !plan.Identity.IsUnknown() {
// 		identityPlan, diags = plan.Identity.Get(ctx)
// 		if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
// 			return
// 		}
// 	}

// 	if !identityPlan.Type.IsNull() {
// 		tflog.Debug(ctx, "PROVISION IDENTITY", map[string]any{
// 			"action": "start",
// 			"id":     plan.ID.ValueString(),
// 		})

// 		_, err := r.client.ProvisionIdentity(ctx, plan.ID.ValueString(), nil)
// 		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
// 			return
// 		}

// 		tflog.Debug(ctx, "PROVISION IDENTITY", map[string]any{
// 			"action": "end",
// 			"id":     plan.ID.ValueString(),
// 		})
// 	}

// 	if resp.Diagnostics.Append(r.get(ctx, &state)...); resp.Diagnostics.HasError() {
// 		return
// 	}

// 	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

// 	tflog.Debug(ctx, "CREATE", map[string]any{
// 		"action": "end",
// 	})

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}
// }

// func (r *resourceOnelakeShortcut) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
// 	tflog.Debug(ctx, "READ", map[string]any{
// 		"action": "start",
// 	})

// 	var state resourceOneLakeShortcutModel

// 	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
// 		return
// 	}

// 	timeout, diags := state.Timeouts.Read(ctx, r.pConfigData.Timeout)
// 	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(ctx, timeout)
// 	defer cancel()

// 	diags = r.get(ctx, &state)
// 	if utils.IsErrNotFound(state.ID.ValueString(), &diags, fabcore.ErrWorkspace.WorkspaceNotFound) {
// 		resp.State.RemoveResource(ctx)

// 		resp.Diagnostics.Append(diags...)

// 		return
// 	}

// 	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
// 		return
// 	}

// 	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

// 	tflog.Debug(ctx, "READ", map[string]any{
// 		"action": "end",
// 	})

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}
// }

// func (r *resourceOnelakeShortcut) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { //nolint:gocognit
// 	tflog.Debug(ctx, "UPDATE", map[string]any{
// 		"action": "start",
// 	})

// 	var plan, state, intermediary resourceOneLakeShortcutModel

// 	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
// 	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	intermediary.ID = plan.ID

// 	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
// 	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
// 		return
// 	}

// 	intermediary.Timeouts = plan.Timeouts

// 	ctx, cancel := context.WithTimeout(ctx, timeout)
// 	defer cancel()

// 	var reqUpdatePlan requestUpdateWorkspace
// 	var reqUpdateState requestUpdateWorkspace

// 	reqUpdatePlan.set(plan)
// 	reqUpdateState.set(state)

// 	if !reflect.DeepEqual(reqUpdatePlan.UpdateWorkspaceRequest, reqUpdateState.UpdateWorkspaceRequest) {
// 		_, err := r.client.UpdateWorkspace(ctx, plan.ID.ValueString(), reqUpdatePlan.UpdateWorkspaceRequest, nil)
// 		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
// 			return
// 		}

// 		if resp.Diagnostics.Append(r.get(ctx, &intermediary)...); resp.Diagnostics.HasError() {
// 			return
// 		}

// 		if resp.Diagnostics.Append(resp.State.Set(ctx, intermediary)...); resp.Diagnostics.HasError() {
// 			return
// 		}
// 	}

// 	if !plan.CapacityID.Equal(state.CapacityID) {
// 		var err error

// 		if plan.CapacityID.IsNull() {
// 			tflog.Debug(ctx, "UNASSIGN CAPACITY", map[string]any{
// 				"action": "start",
// 				"id":     plan.ID.ValueString(),
// 			})

// 			_, err = r.client.UnassignFromCapacity(ctx, plan.ID.ValueString(), nil)

// 			tflog.Debug(ctx, "UNASSIGN CAPACITY", map[string]any{
// 				"action": "end",
// 				"id":     plan.ID.ValueString(),
// 			})
// 		} else {
// 			tflog.Debug(ctx, "ASSIGN CAPACITY", map[string]any{
// 				"action": "start",
// 				"id":     plan.ID.ValueString(),
// 			})

// 			var reqUpdateCapacity assignWorkspaceToCapacityRequest

// 			reqUpdateCapacity.set(plan)

// 			_, err = r.client.AssignToCapacity(ctx, plan.ID.ValueString(), reqUpdateCapacity.AssignWorkspaceToCapacityRequest, nil)

// 			tflog.Debug(ctx, "ASSIGN CAPACITY", map[string]any{
// 				"action": "end",
// 				"id":     plan.ID.ValueString(),
// 			})
// 		}

// 		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
// 			return
// 		}

// 		if resp.Diagnostics.Append(r.get(ctx, &intermediary)...); resp.Diagnostics.HasError() {
// 			return
// 		}

// 		if resp.Diagnostics.Append(resp.State.Set(ctx, intermediary)...); resp.Diagnostics.HasError() {
// 			return
// 		}
// 	}

// 	if !plan.Identity.Equal(state.Identity) { //nolint:nestif
// 		identityPlan := &workspaceIdentityModel{}
// 		identityState := &workspaceIdentityModel{}

// 		if !plan.Identity.IsNull() && !plan.Identity.IsUnknown() {
// 			identityPlan, diags = plan.Identity.Get(ctx)
// 			if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
// 				return
// 			}
// 		}

// 		if !state.Identity.IsNull() && !state.Identity.IsUnknown() {
// 			identityState, diags = state.Identity.Get(ctx)
// 			if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
// 				return
// 			}
// 		}

// 		if !identityPlan.Type.Equal(identityState.Type) {
// 			var err error

// 			if !identityPlan.Type.IsNull() && identityState.Type.IsNull() {
// 				tflog.Debug(ctx, "PROVISION IDENTITY", map[string]any{
// 					"action": "start",
// 					"id":     plan.ID.ValueString(),
// 				})

// 				_, err = r.client.ProvisionIdentity(ctx, plan.ID.ValueString(), nil)

// 				tflog.Debug(ctx, "PROVISION IDENTITY", map[string]any{
// 					"action": "end",
// 					"id":     plan.ID.ValueString(),
// 				})
// 			} else if identityPlan.Type.IsNull() && !identityState.Type.IsNull() {
// 				tflog.Debug(ctx, "DEPROVISION IDENTITY", map[string]any{
// 					"action": "start",
// 					"id":     plan.ID.ValueString(),
// 				})

// 				_, err = r.client.DeprovisionIdentity(ctx, plan.ID.ValueString(), nil)

// 				tflog.Debug(ctx, "DEPROVISION IDENTITY", map[string]any{
// 					"action": "end",
// 					"id":     plan.ID.ValueString(),
// 				})
// 			}

// 			if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
// 				return
// 			}
// 		}
// 	}

// 	if resp.Diagnostics.Append(r.get(ctx, &plan)...); resp.Diagnostics.HasError() {
// 		return
// 	}

// 	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

// 	tflog.Debug(ctx, "UPDATE", map[string]any{
// 		"action": "end",
// 	})

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}
// }

// func (r *resourceOnelakeShortcut) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
// 	tflog.Debug(ctx, "DELETE", map[string]any{
// 		"action": "start",
// 	})

// 	var state resourceOneLakeShortcutModel

// 	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
// 		return
// 	}

// 	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
// 	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(ctx, timeout)
// 	defer cancel()

// 	if !state.Identity.IsNull() && !state.Identity.IsUnknown() {
// 		identityState, diags := state.Identity.Get(ctx)
// 		if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
// 			return
// 		}

// 		if !identityState.ApplicationID.IsNull() && !identityState.ApplicationID.IsUnknown() {
// 			tflog.Debug(ctx, "DEPROVISION IDENTITY", map[string]any{
// 				"action": "start",
// 				"id":     state.ID.ValueString(),
// 			})

// 			_, err := r.client.DeprovisionIdentity(ctx, state.ID.ValueString(), nil)
// 			if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
// 				return
// 			}

// 			tflog.Debug(ctx, "DEPROVISION IDENTITY", map[string]any{
// 				"action": "end",
// 				"id":     state.ID.ValueString(),
// 			})
// 		}
// 	}

// 	_, err := r.client.DeleteWorkspace(ctx, state.ID.ValueString(), nil)
// 	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
// 		return
// 	}

// 	tflog.Debug(ctx, "DELETE", map[string]any{
// 		"action": "end",
// 	})
// }

// func (r *resourceOnelakeShortcut) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	tflog.Debug(ctx, "IMPORT", map[string]any{
// 		"action": "start",
// 	})
// 	tflog.Trace(ctx, "IMPORT", map[string]any{
// 		"id": req.ID,
// 	})

// 	_, diags := customtypes.NewUUIDValueMust(req.ID)
// 	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
// 		return
// 	}

// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

// 	tflog.Debug(ctx, "IMPORT", map[string]any{
// 		"action": "end",
// 	})

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}
// }

// func (r *resourceOnelakeShortcut) get(ctx context.Context, model *resourceOneLakeShortcutModel) diag.Diagnostics {
// 	tflog.Trace(ctx, "GET", map[string]any{
// 		"id": model.ID.ValueString(),
// 	})

// 	var diags diag.Diagnostics

// 	for {
// 		respGet, err := r.client.GetWorkspace(ctx, model.ID.ValueString(), nil)
// 		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrWorkspace.WorkspaceNotFound); diags.HasError() {
// 			return diags
// 		}

// 		if diags := checkWorkspaceType(respGet.WorkspaceInfo); diags.HasError() {
// 			return diags
// 		}

// 		switch *respGet.CapacityAssignmentProgress {
// 		case fabcore.CapacityAssignmentProgressFailed:
// 			diags.AddError(
// 				"capacity assignment operation",
// 				"Workspace capacity assignment failed",
// 			)

// 			return diags

// 		case fabcore.CapacityAssignmentProgressCompleted:
// 			return model.set(ctx, respGet.WorkspaceInfo)
// 		default:
// 			tflog.Info(ctx, "Workspace capacity assignment in progress, waiting 30 seconds before retrying")
// 			time.Sleep(30 * time.Second) // lintignore:R018
// 		}
// 	}
// }
