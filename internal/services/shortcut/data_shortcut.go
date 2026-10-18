// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package shortcut

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceShortcut)(nil)

type dataSourceShortcut struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.OneLakeShortcutsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceShortcut() datasource.DataSource {
	return &dataSourceShortcut{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceShortcut) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

func (d *dataSourceShortcut) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = itemSchema(false).GetDataSource(ctx)
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

	if resp.Diagnostics.Append(fabricitem.IsPreviewMode(d.TypeInfo.Name, d.TypeInfo.IsPreview, d.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}

	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewOneLakeShortcutsClient()
}

func (d *dataSourceShortcut) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceShortcutModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, d.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	diags = d.getShortcut(ctx, &data)

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

func (d *dataSourceShortcut) getShortcut(ctx context.Context, model *dataSourceShortcutModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET SHORTCUT", map[string]any{
		"item_Id":      model.ItemID.ValueString(),
		"workspace_id": model.WorkspaceID.ValueString(),
		"path":         model.Path.ValueString(),
		"name":         model.Name.ValueString(),
	})

	respGet, err := d.client.GetShortcut(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), model.Path.ValueString(), model.Name.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	// Normalize mismatch between configuration Path and API path
	if respGet.Path != nil && model.Path.ValueString() != "" {
		apiResponsePath := *respGet.Path
		configPath := model.Path.ValueString()

		if strings.TrimPrefix(configPath, "/") == strings.TrimPrefix(apiResponsePath, "/") {
			respGet.Path = &configPath
		}
	}

	return model.set(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), respGet.Shortcut)
}
