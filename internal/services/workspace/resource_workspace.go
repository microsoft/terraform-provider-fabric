// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"
	"fmt"
	"reflect"
	"time"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure      = (*resourceWorkspace)(nil)
	_ resource.ResourceWithValidateConfig = (*resourceWorkspace)(nil)
	_ resource.ResourceWithImportState    = (*resourceWorkspace)(nil)
)

type resourceWorkspace struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.WorkspacesClient
}

func NewResourceWorkspace() resource.Resource {
	return &resourceWorkspace{}
}

func (r *resourceWorkspace) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ItemTFName
}

func (r *resourceWorkspace) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage a Fabric " + ItemName + ".\n\n" +
			"See [" + ItemName + "](" + ItemDocsURL + ") for more information.\n\n" +
			ItemDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The " + ItemName + " ID.",
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The " + ItemName + " display name.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
					stringvalidator.NoneOfCaseInsensitive("Admin monitoring", "My workspace"),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The " + ItemName + " description.",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(4000),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The " + ItemName + " type.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"capacity_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the Fabric Capacity to assign to the Workspace.",
				Optional:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"capacity_assignment_progress": schema.StringAttribute{
				MarkdownDescription: "A Workspace assignment to capacity progress status.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identity": schema.SingleNestedAttribute{
				MarkdownDescription: "A workspace identity.\n\n" +
					"See [Workspace Identity](https://learn.microsoft.com/fabric/security/workspace-identity) for more information.",
				Optional:   true,
				CustomType: supertypes.NewSingleNestedObjectTypeOf[workspaceIdentityModel](ctx),
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "Provision a workspace identity. Accepted values: " + utils.ConvertStringSlicesToString(workspaceIdentityTypes, true, true) + ".",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(workspaceIdentityTypes...),
						},
					},
					"application_id": schema.StringAttribute{
						MarkdownDescription: "The application ID.",
						Computed:            true,
						CustomType:          customtypes.UUIDType{},
					},
					"service_principal_id": schema.StringAttribute{
						MarkdownDescription: "The service principal ID.",
						Computed:            true,
						CustomType:          customtypes.UUIDType{},
					},
				},
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *resourceWorkspace) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewWorkspacesClient()
}

func (r *resourceWorkspace) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config resourceWorkspaceModel
	var diags diag.Diagnostics

	if resp.Diagnostics.Append(req.Config.Get(ctx, &config)...); resp.Diagnostics.HasError() {
		return
	}

	identityConfig := &workspaceIdentityModel{}

	if !config.Identity.IsNull() && !config.Identity.IsUnknown() {
		identityConfig, diags = config.Identity.Get(ctx)
		if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
			return
		}
	}

	if !identityConfig.Type.IsNull() && config.CapacityID.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("capacity_id"),
			common.ErrorAttConfigMissing,
			"Expected 'capacity_id' to be configured if 'identity.enabled' is true.",
		)
	}
}

func (r *resourceWorkspace) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "CREATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
	})

	var plan resourceWorkspaceModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqCreate requestCreateWorkspace

	reqCreate.set(plan)

	respCreate, err := r.client.CreateWorkspace(ctx, reqCreate.CreateWorkspaceRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	plan.ID = customtypes.NewUUIDPointerValue(respCreate.ID)

	identityPlan := &workspaceIdentityModel{}

	if !plan.Identity.IsNull() && !plan.Identity.IsUnknown() {
		identityPlan, diags = plan.Identity.Get(ctx)
		if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
			return
		}
	}

	if !identityPlan.Type.IsNull() {
		tflog.Debug(ctx, "PROVISION IDENTITY", map[string]any{
			"action": "start",
			"id":     plan.ID.ValueString(),
		})

		_, err := r.client.ProvisionIdentity(ctx, plan.ID.ValueString(), nil)
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
			return
		}

		tflog.Debug(ctx, "PROVISION IDENTITY", map[string]any{
			"action": "end",
			"id":     plan.ID.ValueString(),
		})
	}

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

func (r *resourceWorkspace) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"state": req.State,
	})

	var state resourceWorkspaceModel

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
	if utils.IsErrNotFound(state.ID.ValueString(), &diags, fabcore.ErrWorkspace.WorkspaceNotFound) {
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

func (r *resourceWorkspace) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { //nolint:gocognit
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "UPDATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
		"state":  req.State,
	})

	var plan, state resourceWorkspaceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdatePlan requestUpdateWorkspace
	var reqUpdateState requestUpdateWorkspace

	reqUpdatePlan.set(plan)
	reqUpdateState.set(state)

	if !reflect.DeepEqual(reqUpdatePlan.UpdateWorkspaceRequest, reqUpdateState.UpdateWorkspaceRequest) {
		_, err := r.client.UpdateWorkspace(ctx, plan.ID.ValueString(), reqUpdatePlan.UpdateWorkspaceRequest, nil)
		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
			return
		}
	}

	if !plan.CapacityID.Equal(state.CapacityID) {
		var err error

		if plan.CapacityID.IsNull() {
			tflog.Debug(ctx, "UNASSIGN CAPACITY", map[string]any{
				"action": "start",
				"id":     plan.ID.ValueString(),
			})

			_, err = r.client.UnassignFromCapacity(ctx, plan.ID.ValueString(), nil)

			tflog.Debug(ctx, "UNASSIGN CAPACITY", map[string]any{
				"action": "end",
				"id":     plan.ID.ValueString(),
			})
		} else {
			tflog.Debug(ctx, "ASSIGN CAPACITY", map[string]any{
				"action": "start",
				"id":     plan.ID.ValueString(),
			})

			var reqUpdateCapacity assignWorkspaceToCapacityRequest

			reqUpdateCapacity.set(plan)

			_, err = r.client.AssignToCapacity(ctx, plan.ID.ValueString(), reqUpdateCapacity.AssignWorkspaceToCapacityRequest, nil)

			tflog.Debug(ctx, "ASSIGN CAPACITY", map[string]any{
				"action": "end",
				"id":     plan.ID.ValueString(),
			})
		}

		if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
			return
		}
	}

	if !plan.Identity.Equal(state.Identity) { //nolint:nestif
		identityPlan := &workspaceIdentityModel{}
		identityState := &workspaceIdentityModel{}

		if !plan.Identity.IsNull() && !plan.Identity.IsUnknown() {
			identityPlan, diags = plan.Identity.Get(ctx)
			if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
				return
			}
		}

		if !state.Identity.IsNull() && !state.Identity.IsUnknown() {
			identityState, diags = state.Identity.Get(ctx)
			if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
				return
			}
		}

		if !identityPlan.Type.Equal(identityState.Type) {
			var err error

			if !identityPlan.Type.IsNull() && identityState.Type.IsNull() {
				tflog.Debug(ctx, "PROVISION IDENTITY", map[string]any{
					"action": "start",
					"id":     plan.ID.ValueString(),
				})

				_, err = r.client.ProvisionIdentity(ctx, plan.ID.ValueString(), nil)

				tflog.Debug(ctx, "PROVISION IDENTITY", map[string]any{
					"action": "end",
					"id":     plan.ID.ValueString(),
				})
			} else if identityPlan.Type.IsNull() && !identityState.Type.IsNull() {
				tflog.Debug(ctx, "DEPROVISION IDENTITY", map[string]any{
					"action": "start",
					"id":     plan.ID.ValueString(),
				})

				_, err = r.client.DeprovisionIdentity(ctx, plan.ID.ValueString(), nil)

				tflog.Debug(ctx, "DEPROVISION IDENTITY", map[string]any{
					"action": "end",
					"id":     plan.ID.ValueString(),
				})
			}

			if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
				return
			}
		}
	}

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

func (r *resourceWorkspace) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "DELETE", map[string]any{
		"state": req.State,
	})

	var state resourceWorkspaceModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.DeleteWorkspace(ctx, state.ID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceWorkspace) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *resourceWorkspace) get(ctx context.Context, model *resourceWorkspaceModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET", map[string]any{
		"id": model.ID.ValueString(),
	})

	var diags diag.Diagnostics

	for {
		respGet, err := r.client.GetWorkspace(ctx, model.ID.ValueString(), nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrWorkspace.WorkspaceNotFound); diags.HasError() {
			return diags
		}

		if diags := checkWorkspaceType(respGet.WorkspaceInfo); diags.HasError() {
			return diags
		}

		switch *respGet.CapacityAssignmentProgress {
		case fabcore.CapacityAssignmentProgressFailed:
			diags.AddError(
				"capacity assignment operation",
				"Workspace capacity assignment failed",
			)

			return diags

		case fabcore.CapacityAssignmentProgressCompleted:
			model.set(ctx, respGet.WorkspaceInfo)

			return nil
		default:
			tflog.Info(ctx, "Workspace capacity assignment in progress, waiting 30 seconds before retrying")
			time.Sleep(30 * time.Second) // lintignore:R018
		}
	}
}
