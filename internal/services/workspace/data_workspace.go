// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspace

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	_ datasource.DataSourceWithConfigure        = (*dataSourceWorkspace)(nil)
	_ datasource.DataSourceWithConfigValidators = (*dataSourceWorkspace)(nil)
)

// DataSource is the data source for the Fabric Workspace.
type dataSourceWorkspace struct {
	pConfigData    *pconfig.ProviderData
	client         *fabcore.WorkspacesClient
	clientCapacity *fabcore.CapacitiesClient
	TypeInfo       tftypeinfo.TFTypeInfo
}

// NewDataSource creates a new data source for the Fabric Workspace.
func NewDataSourceWorkspace() datasource.DataSource {
	return &dataSourceWorkspace{
		TypeInfo: ItemTypeInfo,
	}
}

// Metadata sets metadata for the Fabric Workspace data source.
func (d *dataSourceWorkspace) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

// Schema sets the schema for the Fabric Workspace data source.
func (d *dataSourceWorkspace) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = itemSchema(false).GetDataSource(ctx)
}

func (d *dataSourceWorkspace) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("id"),
			path.MatchRoot("display_name"),
		),
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("display_name"),
		),
	}
}

func (d *dataSourceWorkspace) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	pConfigData, ok := req.ProviderData.(*pconfig.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			common.ErrorDataSourceConfigType,
			fmt.Sprintf(common.ErrorFabricClientType, req.ProviderData),
		)

		return
	}

	d.pConfigData = pConfigData
	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewWorkspacesClient()
	d.clientCapacity = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewCapacitiesClient()
}

func (d *dataSourceWorkspace) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceWorkspaceModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, d.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if data.ID.ValueString() != "" {
		diags = d.getByID(ctx, &data)
	} else {
		diags = d.getByDisplayName(ctx, &data)
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dataSourceWorkspace) getByID(ctx context.Context, model *dataSourceWorkspaceModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY ID", map[string]any{
		"id": model.ID.ValueString(),
	})

	respGet, err := d.client.GetWorkspace(ctx, model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	if diags := checkWorkspaceType(respGet.WorkspaceInfo); diags.HasError() {
		return diags
	}

	diags := model.set(ctx, respGet.WorkspaceInfo)
	if diags.HasError() {
		return diags
	}

	return validateCapacityState(ctx, d.clientCapacity, model.CapacityID.ValueStringPointer())
}

func (d *dataSourceWorkspace) getByDisplayName(ctx context.Context, model *dataSourceWorkspaceModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY DISPLAY NAME", map[string]any{
		"display_name": model.DisplayName.ValueString(),
	})

	var diags diag.Diagnostics

	pager := d.client.NewListWorkspacesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Value {
			if *entity.DisplayName == model.DisplayName.ValueString() {
				model.ID = customtypes.NewUUIDPointerValue(entity.ID)

				return d.getByID(ctx, model)
			}
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		"Unable to find Workspace with 'display_name': "+model.DisplayName.ValueString(),
	)

	return diags
}
