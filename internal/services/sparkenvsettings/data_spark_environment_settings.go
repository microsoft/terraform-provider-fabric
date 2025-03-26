// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package sparkenvsettings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabenvironment "github.com/microsoft/fabric-sdk-go/fabric/environment"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceSparkEnvironmentSettings)(nil)

type dataSourceSparkEnvironmentSettings struct {
	pConfigData *pconfig.ProviderData
	client      *fabenvironment.SparkComputeClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceSparkEnvironmentSettings() datasource.DataSource {
	return &dataSourceSparkEnvironmentSettings{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceSparkEnvironmentSettings) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

func (d *dataSourceSparkEnvironmentSettings) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = itemSchema().GetDataSource(ctx)
}

func (d *dataSourceSparkEnvironmentSettings) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = fabenvironment.NewClientFactoryWithClient(*pConfigData.FabricClient).NewSparkComputeClient()
}

func (d *dataSourceSparkEnvironmentSettings) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceSparkEnvironmentSettingsModel

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

	data.ID = data.EnvironmentID

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dataSourceSparkEnvironmentSettings) get(ctx context.Context, model *dataSourceSparkEnvironmentSettingsModel) diag.Diagnostics {
	var respEntity fabenvironment.SparkCompute

	if model.PublicationStatus.ValueString() == SparkEnvironmentPublicationStatusPublished {
		respGet, err := d.client.GetPublishedSettings(ctx, model.WorkspaceID.ValueString(), model.EnvironmentID.ValueString(), nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
			return diags
		}

		respEntity = respGet.SparkCompute
	} else {
		respGet, err := d.client.GetStagingSettings(ctx, model.WorkspaceID.ValueString(), model.EnvironmentID.ValueString(), nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
			return diags
		}

		respEntity = respGet.SparkCompute
	}

	return model.set(ctx, respEntity)
}
