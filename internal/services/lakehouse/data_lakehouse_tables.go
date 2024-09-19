// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	"context"
	"fmt"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
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

var _ datasource.DataSourceWithConfigure = (*dataSourceLakehouseTables)(nil)

type dataSourceLakehouseTables struct {
	pConfigData *pconfig.ProviderData
	client      *fablakehouse.TablesClient
}

func NewDataSourceLakehouseTables() datasource.DataSource {
	return &dataSourceLakehouseTables{}
}

func (d *dataSourceLakehouseTables) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + LakehouseTablesTFName
}

func (d *dataSourceLakehouseTables) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List a Fabric " + LakehouseTablesName + ".\n\n" +
			"Use this data source to list [" + LakehouseTablesName + "](" + LakehouseTableDocsURL + ").\n\n" +
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
			"values": schema.ListNestedAttribute{
				MarkdownDescription: "The list of Lakehouse Tables.",
				Computed:            true,
				CustomType:          supertypes.NewListNestedObjectTypeOf[lakehouseTableModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The Name of the table.",
							Computed:            true,
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
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceLakehouseTables) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dataSourceLakehouseTables) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceLakehouseTablesModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, d.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if resp.Diagnostics.Append(d.list(ctx, &data)...); resp.Diagnostics.HasError() {
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

func (d *dataSourceLakehouseTables) list(ctx context.Context, model *dataSourceLakehouseTablesModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting Lakehouse Tables")

	respList, err := d.client.ListTables(ctx, model.WorkspaceID.ValueString(), model.LakehouseID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	return model.setValues(ctx, respList)
}
