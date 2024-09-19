// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fablakehouse "github.com/microsoft/fabric-sdk-go/fabric/lakehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var _ datasource.DataSourceWithConfigure = (*dataSourceLakehouseTable)(nil)

type dataSourceLakehouseTable struct {
	pConfigData *pconfig.ProviderData
	client      *fablakehouse.TablesClient
}

func NewDataSourceLakehouseTable() datasource.DataSource {
	return &dataSourceLakehouseTable{}
}

func (d *dataSourceLakehouseTable) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + LakehouseTableTFName
}

func (d *dataSourceLakehouseTable) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get a Fabric " + LakehouseTableName + ".\n\n" +
			"Use this data source to get [" + LakehouseTableName + "](" + LakehouseTableDocsURL + ").\n\n" +
			LakehouseTableDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"lakehouse_id": schema.StringAttribute{
				MarkdownDescription: "The Lakehouse ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The Workspace ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The Name of the table.",
				Required:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "The Location of the table.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The Type of the table. Possible values: " + utils.ConvertStringSlicesToString(fablakehouse.PossibleTableTypeValues(), true, true) + ".",
				Computed:            true,
			},
			"format": schema.StringAttribute{
				MarkdownDescription: "The Format of the table.",
				Computed:            true,
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceLakehouseTable) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fablakehouse.NewClientFactoryWithClient(*pConfigData.FabricClient).NewTablesClient()
}

func (d *dataSourceLakehouseTable) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceLakehouseTableModel

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

func (d *dataSourceLakehouseTable) get(ctx context.Context, model *dataSourceLakehouseTableModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting Lakehouse Table")

	var diags diag.Diagnostics

	pager := d.client.NewListTablesPager(model.WorkspaceID.ValueString(), model.LakehouseID.ValueString(), nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Data {
			if *entity.Name == model.Name.ValueString() {
				model.set(entity)

				return nil
			}
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		"Unable to find Table with 'name': "+model.Name.ValueString(),
	)

	return diags
}
