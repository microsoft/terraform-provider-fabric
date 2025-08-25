package externaldatasharesprovider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceExternalDataSharesProvider)(nil)

type dataSourceExternalDataSharesProvider struct {
	pConfigData *pconfig.ProviderData
	client      *fabadmin.ExternalDataSharesProviderClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceExternalDataSharesProvider() datasource.DataSource {
	return &dataSourceExternalDataSharesProvider{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceExternalDataSharesProvider) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(true)
}

func (d *dataSourceExternalDataSharesProvider) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = itemSchema().GetDataSource(ctx)
}

func (d *dataSourceExternalDataSharesProvider) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = fabadmin.NewClientFactoryWithClient(*pConfigData.FabricClient).NewExternalDataSharesProviderClient()
}

func (d *dataSourceExternalDataSharesProvider) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data baseExternalDataSharesProviderModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(d.list(ctx, &data)...); resp.Diagnostics.HasError() {
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

func (d *dataSourceExternalDataSharesProvider) list(ctx context.Context, model *baseExternalDataSharesProviderModel) diag.Diagnostics {
	respList, err := d.client.ListExternalDataShares(ctx, nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		diags.AddError(
			common.ErrorReadHeader,
			"Unable to find any items.",
		)

		return diags
	}

	return model.set(ctx, respList)
}
