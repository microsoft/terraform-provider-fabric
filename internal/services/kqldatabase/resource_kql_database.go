// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package kqldatabase

import (
	"context"
	"fmt"
	"strings"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	superstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"
	fabkqldatabase "github.com/microsoft/fabric-sdk-go/fabric/kqldatabase"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure   = (*resourceKQLDatabase)(nil)
	_ resource.ResourceWithImportState = (*resourceKQLDatabase)(nil)
)

type resourceKQLDatabase struct {
	pConfigData *pconfig.ProviderData
	client      *fabkqldatabase.ItemsClient
}

func NewResourceKQLDatabase() resource.Resource {
	return &resourceKQLDatabase{}
}

func (r *resourceKQLDatabase) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ItemTFName
}

func (r *resourceKQLDatabase) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	markdownDescription := "This resource manages a Fabric " + ItemName + ".\n\n" +
		"See [" + ItemName + "](" + ItemDocsURL + ") for more information.\n\n" +
		ItemDocsSPNSupport

	properties := schema.SingleNestedAttribute{
		MarkdownDescription: "The " + ItemName + " properties.",
		Computed:            true,
		CustomType:          supertypes.NewSingleNestedObjectTypeOf[kqlDatabasePropertiesModel](ctx),
		Attributes: map[string]schema.Attribute{
			"database_type": schema.StringAttribute{
				MarkdownDescription: "The type of the database. Possible values:" + utils.ConvertStringSlicesToString(fabkqldatabase.PossibleKqlDatabaseTypeValues(), true, true) + ".",
				Computed:            true,
			},
			"eventhouse_id": schema.StringAttribute{
				MarkdownDescription: "Parent Eventhouse ID.",
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"ingestion_service_uri": schema.StringAttribute{
				MarkdownDescription: "Ingestion service URI.",
				Computed:            true,
				CustomType:          customtypes.URLType{},
			},
			"query_service_uri": schema.StringAttribute{
				MarkdownDescription: "Query service URI.",
				Computed:            true,
				CustomType:          customtypes.URLType{},
			},
		},
	}

	configuration := schema.SingleNestedAttribute{
		MarkdownDescription: "The " + ItemName + " creation configuration.\n\n" +
			"Any changes to this configuration will result in recreation of the " + ItemName + ".",
		Required:   true,
		CustomType: supertypes.NewSingleNestedObjectTypeOf[kqlDatabaseConfigurationModel](ctx),
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
		Attributes: map[string]schema.Attribute{
			"database_type": schema.StringAttribute{
				MarkdownDescription: "The type of the KQL database. Accepted values: " + utils.ConvertStringSlicesToString(fabkqldatabase.PossibleKqlDatabaseTypeValues(), true, true) + ".\n\n" +
					"`" + string(fabkqldatabase.TypeReadWrite) + "` Allows read and write operations on the database.\n\n" +
					"`" + string(fabkqldatabase.TypeShortcut) + "` A shortcut is an embedded reference allowing read only operations on a source database. The source can be in the same or different tenants, either in an Azure Data Explorer cluster or a Fabric Eventhouse.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(utils.ConvertEnumsToStringSlices(fabkqldatabase.PossibleKqlDatabaseTypeValues(), false)...),
				},
			},
			"eventhouse_id": schema.StringAttribute{
				MarkdownDescription: "Parent Eventhouse ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"invitation_token": schema.StringAttribute{
				MarkdownDescription: "Invitation token to follow the source database. Only allowed when `database_type` is `" + string(fabkqldatabase.TypeShortcut) + "`.",
				Optional:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRelative().AtParent().AtName("source_cluster_uri"),
						path.MatchRelative().AtParent().AtName("source_database_name"),
					),
					superstringvalidator.NullIfAttributeIsOneOf(
						path.MatchRelative().AtParent().AtName("database_type"),
						[]attr.Value{types.StringValue(string(fabkqldatabase.TypeReadWrite))},
					),
				},
			},
			"source_cluster_uri": schema.StringAttribute{
				MarkdownDescription: "The URI of the source Eventhouse or Azure Data Explorer cluster. Only allowed when `database_type` is `" + string(fabkqldatabase.TypeShortcut) + "`.",
				Optional:            true,
				CustomType:          customtypes.URLType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("invitation_token")),
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("source_database_name")),
					superstringvalidator.NullIfAttributeIsOneOf(
						path.MatchRelative().AtParent().AtName("database_type"),
						[]attr.Value{types.StringValue(string(fabkqldatabase.TypeReadWrite))},
					),
				},
			},
			"source_database_name": schema.StringAttribute{
				MarkdownDescription: "The name of the database to follow in the source Eventhouse or Azure Data Explorer cluster. Only allowed when `database_type` is `" + string(fabkqldatabase.TypeShortcut) + "`.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("invitation_token")),
					superstringvalidator.NullIfAttributeIsOneOf(
						path.MatchRelative().AtParent().AtName("database_type"),
						[]attr.Value{types.StringValue(string(fabkqldatabase.TypeReadWrite))},
					),
				},
			},
		},
	}

	resp.Schema = fabricitem.GetResourceFabricItemPropertiesCreationSchema(ctx, ItemName, markdownDescription, 123, 256, true, properties, configuration)
}

func (r *resourceKQLDatabase) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = fabkqldatabase.NewClientFactoryWithClient(*pConfigData.FabricClient).NewItemsClient()
}

func (r *resourceKQLDatabase) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "CREATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "CREATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
	})

	var plan resourceKQLDatabaseModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqCreate requestCreateKQLDatabase

	if resp.Diagnostics.Append(reqCreate.set(ctx, plan)...); resp.Diagnostics.HasError() {
		return
	}

	respCreate, err := r.client.CreateKQLDatabase(ctx, plan.WorkspaceID.ValueString(), reqCreate.CreateKQLDatabaseRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationCreate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	plan.set(respCreate.KQLDatabase)

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

func (r *resourceKQLDatabase) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"state": req.State,
	})

	var state resourceKQLDatabaseModel

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

func (r *resourceKQLDatabase) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "UPDATE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "UPDATE", map[string]any{
		"config": req.Config,
		"plan":   req.Plan,
		"state":  req.State,
	})

	var plan resourceKQLDatabaseModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var reqUpdate requestUpdateKQLDatabase

	reqUpdate.set(plan)

	respUpdate, err := r.client.UpdateKQLDatabase(ctx, plan.WorkspaceID.ValueString(), plan.ID.ValueString(), reqUpdate.UpdateKQLDatabaseRequest, nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationUpdate, nil)...); resp.Diagnostics.HasError() {
		return
	}

	plan.set(respUpdate.KQLDatabase)

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

func (r *resourceKQLDatabase) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "DELETE", map[string]any{
		"state": req.State,
	})

	var state resourceKQLDatabaseModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, r.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.client.DeleteKQLDatabase(ctx, state.WorkspaceID.ValueString(), state.ID.ValueString(), nil)
	if resp.Diagnostics.Append(utils.GetDiagsFromError(ctx, err, utils.OperationDelete, nil)...); resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DELETE", map[string]any{
		"action": "end",
	})
}

func (r *resourceKQLDatabase) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "IMPORT", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "IMPORT", map[string]any{
		"id": req.ID,
	})

	workspaceID, kqlDatabaseID, found := strings.Cut(req.ID, "/")

	if !found {
		resp.Diagnostics.AddError(
			common.ErrorImportIdentifierHeader,
			fmt.Sprintf(common.ErrorImportIdentifierDetails, "WorkspaceID/KQLDatabaseID"),
		)

		return
	}

	uuidWorkspaceID, diags := customtypes.NewUUIDValueMust(workspaceID)
	resp.Diagnostics.Append(diags...)

	uuidID, diags := customtypes.NewUUIDValueMust(kqlDatabaseID)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var configuration supertypes.SingleNestedObjectValueOf[kqlDatabaseConfigurationModel]
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("configuration"), &configuration)...); resp.Diagnostics.HasError() {
		return
	}

	var timeout timeouts.Value
	if resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeout)...); resp.Diagnostics.HasError() {
		return
	}

	state := resourceKQLDatabaseModel{}
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

func (r *resourceKQLDatabase) get(ctx context.Context, model *resourceKQLDatabaseModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET", map[string]any{
		"workspace_id": model.WorkspaceID.ValueString(),
		"id":           model.ID.ValueString(),
	})

	respGet, err := r.client.GetKQLDatabase(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, fabcore.ErrCommon.EntityNotFound); diags.HasError() {
		return diags
	}

	model.set(respGet.KQLDatabase)

	return model.setProperties(ctx, respGet.KQLDatabase)
}
