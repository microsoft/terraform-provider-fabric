// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package workspacempe

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
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigure        = (*dataSourceWorkspaceManagedPrivateEndpoint)(nil)
	_ datasource.DataSourceWithConfigValidators = (*dataSourceWorkspaceManagedPrivateEndpoint)(nil)
)

type dataSourceWorkspaceManagedPrivateEndpoint struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.ManagedPrivateEndpointsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceWorkspaceManagedPrivateEndpoint() datasource.DataSource {
	return &dataSourceWorkspaceManagedPrivateEndpoint{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceWorkspaceManagedPrivateEndpoint) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

func (d *dataSourceWorkspaceManagedPrivateEndpoint) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = itemSchema(false).GetDataSource(ctx)
}

func (d *dataSourceWorkspaceManagedPrivateEndpoint) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

func (d *dataSourceWorkspaceManagedPrivateEndpoint) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(d.TypeInfo.Name, d.TypeInfo.IsPreview, d.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}

	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewManagedPrivateEndpointsClient()
}

func (d *dataSourceWorkspaceManagedPrivateEndpoint) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceWorkspaceManagedPrivateEndpointModel

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
		diags = d.getByName(ctx, &data)
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

func (d *dataSourceWorkspaceManagedPrivateEndpoint) getByID(ctx context.Context, model *dataSourceWorkspaceManagedPrivateEndpointModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY ID", map[string]any{
		"action": "start",
	})

	respGet, err := d.client.GetWorkspaceManagedPrivateEndpoint(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	tflog.Trace(ctx, "GET BY ID", map[string]any{
		"action": "end",
	})

	return model.set(ctx, model.WorkspaceID.ValueString(), respGet.ManagedPrivateEndpoint)
}

func (d *dataSourceWorkspaceManagedPrivateEndpoint) getByName(ctx context.Context, model *dataSourceWorkspaceManagedPrivateEndpointModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY NAME", map[string]any{
		"model": model,
	})

	var diags diag.Diagnostics

	pager := d.client.NewListWorkspaceManagedPrivateEndpointsPager(model.WorkspaceID.ValueString(), nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Value {
			if *entity.Name == model.Name.ValueString() {
				return model.set(ctx, model.WorkspaceID.ValueString(), entity)
			}
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		fmt.Sprintf("Unable to find %s with 'display_name': %s", d.TypeInfo.Name, model.Name.ValueString()),
	)

	return diags
}
