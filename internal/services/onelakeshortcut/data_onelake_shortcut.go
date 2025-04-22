// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package onelakeshortcut

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
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigure        = (*dataSourceShortcut)(nil)
	_ datasource.DataSourceWithConfigValidators = (*dataSourceShortcut)(nil)
)

// DataSource is the data source for the Fabric OnelakeShortcut.
type dataSourceShortcut struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.OneLakeShortcutsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

// NewDataSource creates a new data source for the Fabric OnelakeShortcut.
func NewDataSourceOnelakeShortcut() datasource.DataSource {
	return &dataSourceShortcut{
		TypeInfo: ItemTypeInfo,
	}
}

// Metadata sets metadata for the Fabric OnelakeShortcut data source.
func (d *dataSourceShortcut) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

// Schema sets the schema for the Fabric OnelakeShortcut data source.
func (d *dataSourceShortcut) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = itemSchema(false).GetDataSource(ctx)
}

func (d *dataSourceShortcut) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("path"),
			path.MatchRoot("name"),
			path.MatchRoot("target"),
		),
	}
}

func (d *dataSourceShortcut) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewOneLakeShortcutsClient()
}

func (d *dataSourceShortcut) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceOnelakeShortcutModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, d.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	diags = d.getByDisplayName(ctx, &data)

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

func (d *dataSourceShortcut) getByID(ctx context.Context, model *dataSourceOnelakeShortcutModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY ID", map[string]any{
		// TODO: concatenate name and path
		"id": model.ItemID.ValueString(),
	})

	respGet, err := d.client.GetShortcut(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), model.Path.ValueString(), model.Name.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	return model.set(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), respGet.Shortcut)
}

// TODO rename get functions
func (d *dataSourceShortcut) getByDisplayName(ctx context.Context, model *dataSourceOnelakeShortcutModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY DISPLAY NAME", map[string]any{
		"name": model.Name.ValueString(),
	})

	var diags diag.Diagnostics

	pager := d.client.NewListShortcutsPager(model.WorkspaceID.ValueString(), model.ItemID.ValueString(), nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Value {
			if *entity.Name == model.Name.ValueString() {
				return d.getByID(ctx, model)
			}
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		"Unable to find OnelakeShortcut with 'display_name': "+model.Name.ValueString(),
	)

	return diags
}
