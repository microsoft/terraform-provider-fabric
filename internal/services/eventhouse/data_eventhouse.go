// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package eventhouse

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
	fabeventhouse "github.com/microsoft/fabric-sdk-go/fabric/eventhouse"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigValidators = (*dataSourceEventhouse)(nil)
	_ datasource.DataSourceWithConfigure        = (*dataSourceEventhouse)(nil)
)

type dataSourceEventhouse struct {
	pConfigData *pconfig.ProviderData
	client      *fabeventhouse.ItemsClient
}

func NewDataSourceEventhouse() datasource.DataSource {
	return &dataSourceEventhouse{}
}

func (d *dataSourceEventhouse) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ItemTFName
}

func (d *dataSourceEventhouse) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	markdownDescription := "Get a Fabric " + ItemName + ".\n\n" +
		"Use this data source to fetch a [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
		ItemDocsSPNSupport

	properties := schema.SingleNestedAttribute{
		MarkdownDescription: "The " + ItemName + " properties.",
		Computed:            true,
		CustomType:          supertypes.NewSingleNestedObjectTypeOf[eventhousePropertiesModel](ctx),
		Attributes: map[string]schema.Attribute{
			"ingestion_service_uri": schema.StringAttribute{
				MarkdownDescription: "Ingestion service URI.",
				Computed:            true,
			},
			"query_service_uri": schema.StringAttribute{
				MarkdownDescription: "Query service URI.",
				Computed:            true,
			},
			"database_ids": schema.ListAttribute{
				MarkdownDescription: "The IDs list of KQL Databases.",
				Computed:            true,
				CustomType:          supertypes.NewListTypeOf[string](ctx),
			},
		},
	}

	resp.Schema = fabricitem.GetDataSourceFabricItemPropertiesSchema(ctx, ItemName, markdownDescription, true, properties)
}

func (d *dataSourceEventhouse) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
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

func (d *dataSourceEventhouse) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fabeventhouse.NewClientFactoryWithClient(*pConfigData.FabricClient).NewItemsClient()
}

func (d *dataSourceEventhouse) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceEventhouseModel

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

func (d *dataSourceEventhouse) getByID(ctx context.Context, model *dataSourceEventhouseModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY ID", map[string]any{
		"workspace_id": model.WorkspaceID.ValueString(),
		"id":           model.ID.ValueString(),
	})

	respGet, err := d.client.GetEventhouse(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	model.set(respGet.Eventhouse)

	return model.setProperties(ctx, respGet.Eventhouse)
}

func (d *dataSourceEventhouse) getByDisplayName(ctx context.Context, model *dataSourceEventhouseModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY DISPLAY NAME", map[string]any{
		"workspace_id": model.WorkspaceID.ValueString(),
		"display_name": model.DisplayName.ValueString(),
	})

	var diags diag.Diagnostics

	pager := d.client.NewListEventhousesPager(model.WorkspaceID.ValueString(), nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Value {
			if *entity.DisplayName == model.DisplayName.ValueString() {
				model.set(entity)

				return model.setProperties(ctx, entity)
			}
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		fmt.Sprintf("Unable to find %s with display_name: '%s' in the Workspace ID: %s ", ItemName, model.DisplayName.ValueString(), model.WorkspaceID.ValueString()),
	)

	return diags
}
