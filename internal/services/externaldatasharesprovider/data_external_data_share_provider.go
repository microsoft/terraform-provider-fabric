// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package externaldatasharesprovider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceExternalDataShareProvider)(nil)

type dataSourceExternalDataShareProvider struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.ExternalDataSharesProviderClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceExternalDataShareProvider() datasource.DataSource {
	return &dataSourceExternalDataShareProvider{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceExternalDataShareProvider) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

func (d *dataSourceExternalDataShareProvider) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	s := itemSchema(false).GetDataSource(ctx)

	resp.Schema = schema.Schema{
		MarkdownDescription: s.GetMarkdownDescription(),
		Attributes: map[string]schema.Attribute{
			"external_data_share_id": schema.StringAttribute{
				MarkdownDescription: "The External Data Share ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
		},
	}

	for k, v := range s.Attributes {
		if _, exists := resp.Schema.Attributes[k]; exists {
			continue
		}

		resp.Schema.Attributes[k] = v
	}
}

func (d *dataSourceExternalDataShareProvider) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewExternalDataSharesProviderClient()
}

func (d *dataSourceExternalDataShareProvider) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceExternalDataShareProviderModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(d.getByID(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dataSourceExternalDataShareProvider) getByID(ctx context.Context, model *dataSourceExternalDataShareProviderModel) diag.Diagnostics {
	respList, err := d.client.GetExternalDataShare(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), model.ExternalDataShareID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		diags.AddError(
			common.ErrorReadHeader,
			"Unable to find any items.",
		)

		return diags
	}

	model.set(ctx, &respList.ExternalDataShare)

	return nil
}
