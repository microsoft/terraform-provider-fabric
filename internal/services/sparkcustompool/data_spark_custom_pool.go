// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package sparkcustompool

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabspark "github.com/microsoft/fabric-sdk-go/fabric/spark"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var (
	_ datasource.DataSourceWithConfigValidators = (*dataSourceSparkCustomPool)(nil)
	_ datasource.DataSourceWithConfigure        = (*dataSourceSparkCustomPool)(nil)
)

type dataSourceSparkCustomPool struct {
	pConfigData *pconfig.ProviderData
	client      *fabspark.CustomPoolsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceSparkCustomPool() datasource.DataSource {
	return &dataSourceSparkCustomPool{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceSparkCustomPool) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

func (d *dataSourceSparkCustomPool) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = itemSchema().GetDataSource(ctx)
}

func (d *dataSourceSparkCustomPool) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
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

func (d *dataSourceSparkCustomPool) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = fabspark.NewClientFactoryWithClient(*pConfigData.FabricClient).NewCustomPoolsClient()
}

func (d *dataSourceSparkCustomPool) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceSparkCustomPoolModel

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

func (d *dataSourceSparkCustomPool) getByID(ctx context.Context, model *dataSourceSparkCustomPoolModel) diag.Diagnostics {
	respGet, err := d.client.GetWorkspaceCustomPool(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	return model.set(ctx, respGet.CustomPool)
}

func (d *dataSourceSparkCustomPool) getByDisplayName(ctx context.Context, model *dataSourceSparkCustomPoolModel) diag.Diagnostics {
	pager := d.client.NewListWorkspaceCustomPoolsPager(model.WorkspaceID.ValueString(), nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Value {
			if *entity.Name == model.Name.ValueString() {
				return model.set(ctx, entity)
			}
		}
	}

	var diags diag.Diagnostics

	diags.AddError(
		common.ErrorReadHeader,
		fmt.Sprintf("Unable to find "+d.TypeInfo.Name+" with name: '%s' in the Workspace ID: %s ", model.Name.ValueString(), model.WorkspaceID.ValueString()),
	)

	return diags
}
