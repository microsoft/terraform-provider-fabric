// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package lakehouse

import (
	"context"
	"fmt"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fablakehouse "github.com/microsoft/fabric-sdk-go/fabric/lakehouse"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigValidators = (*dataSourceLakehouse)(nil)
	_ datasource.DataSourceWithConfigure        = (*dataSourceLakehouse)(nil)
)

type dataSourceLakehouse struct {
	pConfigData *pconfig.ProviderData
	client      *fablakehouse.ItemsClient
}

func NewDataSourceLakehouse() datasource.DataSource {
	return &dataSourceLakehouse{}
}

func (d *dataSourceLakehouse) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ItemTFName
}

func (d *dataSourceLakehouse) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	markdownDescription := "Get a Fabric " + ItemName + ".\n\n" +
		"Use this data source to fetch a [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
		ItemDocsSPNSupport

	properties := schema.SingleNestedAttribute{
		MarkdownDescription: "The " + ItemName + " properties.",
		Computed:            true,
		CustomType:          supertypes.NewSingleNestedObjectTypeOf[lakehousePropertiesModel](ctx),
		Attributes: map[string]schema.Attribute{
			"onelake_files_path": schema.StringAttribute{
				MarkdownDescription: "OneLake path to the Lakehouse files directory",
				Computed:            true,
			},
			"onelake_tables_path": schema.StringAttribute{
				MarkdownDescription: "OneLake path to the Lakehouse tables directory.",
				Computed:            true,
			},
			"sql_endpoint_properties": schema.SingleNestedAttribute{
				MarkdownDescription: "An object containing the properties of the SQL endpoint.",
				Computed:            true,
				CustomType:          supertypes.NewSingleNestedObjectTypeOf[lakehouseSQLEndpointPropertiesModel](ctx),
				Attributes: map[string]schema.Attribute{
					"provisioning_status": schema.StringAttribute{
						MarkdownDescription: "The SQL endpoint provisioning status.",
						Computed:            true,
					},
					"connection_string": schema.StringAttribute{
						MarkdownDescription: "SQL endpoint connection string.",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "SQL endpoint ID.",
						Computed:            true,
						CustomType:          customtypes.UUIDType{},
					},
				},
			},
			"default_schema": schema.StringAttribute{
				MarkdownDescription: "Default schema of the Lakehouse. This property is returned only for schema enabled Lakehouse.",
				Computed:            true,
			},
		},
	}

	resp.Schema = fabricitem.GetDataSourceFabricItemPropertiesSchema(ctx, ItemName, markdownDescription, true, properties)
}

func (d *dataSourceLakehouse) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
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

func (d *dataSourceLakehouse) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fablakehouse.NewClientFactoryWithClient(*pConfigData.FabricClient).NewItemsClient()
}

func (d *dataSourceLakehouse) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceLakehouseModel

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

func (d *dataSourceLakehouse) getByID(ctx context.Context, model *dataSourceLakehouseModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting Lakehouse by 'id'")

	respGet, err := d.client.GetLakehouse(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	return model.set(ctx, respGet.Lakehouse)
}

func (d *dataSourceLakehouse) getByDisplayName(ctx context.Context, model *dataSourceLakehouseModel) diag.Diagnostics {
	tflog.Trace(ctx, "getting Lakehouse by 'display_name'")

	var diags diag.Diagnostics

	pager := d.client.NewListLakehousesPager(model.WorkspaceID.ValueString(), nil)
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
		fmt.Sprintf("Unable to find Lakehouse with 'display_name': %s in the Workspace ID: %s ", model.DisplayName.ValueString(), model.WorkspaceID.ValueString()),
	)

	return diags
}
