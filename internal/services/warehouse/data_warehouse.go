// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package warehouse

import (
	"context"
	"fmt"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabwarehouse "github.com/microsoft/fabric-sdk-go/fabric/warehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigValidators = (*dataSourceWarehouse)(nil)
	_ datasource.DataSourceWithConfigure        = (*dataSourceWarehouse)(nil)
)

type dataSourceWarehouse struct {
	pConfigData *pconfig.ProviderData
	client      *fabwarehouse.ItemsClient
}

func NewDataSourceWarehouse() datasource.DataSource {
	return &dataSourceWarehouse{}
}

func (d *dataSourceWarehouse) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ItemTFName
}

func (d *dataSourceWarehouse) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get a Fabric Warehouse.\n\n" +
			"Use this data source to fetch a [Warehouse](https://learn.microsoft.com/fabric/data-warehouse/data-warehousing).\n\n" +
			ItemDocsSPNSupport,
		Attributes: map[string]schema.Attribute{
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The Workspace ID.",
				Required:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The Warehouse ID.",
				Optional:            true,
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The Warehouse display name.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The Warehouse description.",
				Computed:            true,
			},
			"properties": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The Warehouse properties.",
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[warehousePropertiesModel](ctx),
				Attributes: map[string]schema.Attribute{
					"connection_string": schema.StringAttribute{
						MarkdownDescription: "Connection String",
						Computed:            true,
					},
					"created_date": schema.StringAttribute{
						MarkdownDescription: "Created Date",
						Computed:            true,
						CustomType:          timetypes.RFC3339Type{},
					},
					"last_updated_time": schema.StringAttribute{
						MarkdownDescription: "Last Updated Time",
						Computed:            true,
						CustomType:          timetypes.RFC3339Type{},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *dataSourceWarehouse) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("id"),
			path.MatchRoot("display_name"),
		),
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("display_name"),
		),
	}
}

func (d *dataSourceWarehouse) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fabwarehouse.NewClientFactoryWithClient(*pConfigData.FabricClient).NewItemsClient()
}

// Read refreshes the Terraform state with the latest data.
func (d *dataSourceWarehouse) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceWarehouseModel

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

func (d *dataSourceWarehouse) getByID(ctx context.Context, model *dataSourceWarehouseModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting Warehouse by 'id'")

	respGet, err := d.client.GetWarehouse(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	return model.set(ctx, respGet.Warehouse)
}

func (d *dataSourceWarehouse) getByDisplayName(ctx context.Context, model *dataSourceWarehouseModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting Warehouse by 'display_name'")

	var diags diag.Diagnostics

	pager := d.client.NewListWarehousesPager(model.WorkspaceID.ValueString(), nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Value {
			if *entity.DisplayName == model.DisplayName.ValueString() {
				return model.set(ctx, entity)
			}
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		fmt.Sprintf("Unable to find Warehouse with 'display_name': %s in the Workspace ID: %s ", model.DisplayName.ValueString(), model.WorkspaceID.ValueString()),
	)

	return diags
}
