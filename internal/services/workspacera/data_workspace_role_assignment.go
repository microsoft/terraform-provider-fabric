// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package workspacera

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceWorkspaceRoleAssignment)(nil)

type dataSourceWorkspaceRoleAssignment struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.WorkspacesClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceWorkspaceRoleAssignment() datasource.DataSource {
	return &dataSourceWorkspaceRoleAssignment{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceWorkspaceRoleAssignment) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

func (d *dataSourceWorkspaceRoleAssignment) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = itemSchema(false).GetDataSource(ctx)
}

func (d *dataSourceWorkspaceRoleAssignment) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
}

func (d *dataSourceWorkspaceRoleAssignment) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceWorkspaceRoleAssignmentModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, d.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(d.get(ctx, &data)...); resp.Diagnostics.HasError() {
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

func (d *dataSourceWorkspaceRoleAssignment) get(ctx context.Context, model *dataSourceWorkspaceRoleAssignmentModel) diag.Diagnostics {
	respGet, err := d.client.GetWorkspaceRoleAssignment(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	return model.set(ctx, model.WorkspaceID.ValueString(), respGet.WorkspaceRoleAssignment)
}
