// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package kqldatabase

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
	fabkqldatabase "github.com/microsoft/fabric-sdk-go/fabric/kqldatabase"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigValidators = (*dataSourceKQLDatabase)(nil)
	_ datasource.DataSourceWithConfigure        = (*dataSourceKQLDatabase)(nil)
)

type dataSourceKQLDatabase struct {
	pConfigData *pconfig.ProviderData
	client      *fabkqldatabase.ItemsClient
}

func NewDataSourceKQLDatabase() datasource.DataSource {
	return &dataSourceKQLDatabase{}
}

func (d *dataSourceKQLDatabase) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ItemTFName
}

func (d *dataSourceKQLDatabase) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	markdownDescription := "Get a Fabric " + ItemName + ".\n\n" +
		"Use this data source to fetch a [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
		ItemDocsSPNSupport

	properties := schema.SingleNestedAttribute{
		MarkdownDescription: "The KQL Database properties.",
		Computed:            true,
		CustomType:          supertypes.NewSingleNestedObjectTypeOf[kqlDatabasePropertiesModel](ctx),
		Attributes: map[string]schema.Attribute{
			"database_type": schema.StringAttribute{
				MarkdownDescription: "The type of the database. Possible values:" + utils.ConvertStringSlicesToString(fabkqldatabase.PossibleKqlDatabaseTypeValues(), true, true) + ".",
				Computed:            true,
			},
			"eventhouse_id": schema.StringAttribute{
				MarkdownDescription: "Parent Eventhouse ID.",
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
			},
			"ingestion_service_uri": schema.StringAttribute{
				MarkdownDescription: "Ingestion service URI.",
				Computed:            true,
				CustomType:          customtypes.URLType{},
			},
			"query_service_uri": schema.StringAttribute{
				MarkdownDescription: "Query service URI.",
				Computed:            true,
				CustomType:          customtypes.URLType{},
			},
		},
	}

	itemConfig := fabricitem.DataSourceFabricItem{
		Type:                ItemType,
		Name:                ItemName,
		TFName:              ItemTFName,
		MarkdownDescription: markdownDescription,
		IsDisplayNameUnique: true,
	}

	resp.Schema = fabricitem.GetDataSourceFabricItemPropertiesSchema(ctx, itemConfig, properties)
}

func (d *dataSourceKQLDatabase) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
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

func (d *dataSourceKQLDatabase) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fabkqldatabase.NewClientFactoryWithClient(*pConfigData.FabricClient).NewItemsClient()
}

func (d *dataSourceKQLDatabase) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceKQLDatabaseModel

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

func (d *dataSourceKQLDatabase) getByID(ctx context.Context, model *dataSourceKQLDatabaseModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY ID", map[string]any{
		"id": model.ID.ValueString(),
	})

	respGet, err := d.client.GetKQLDatabase(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	model.set(respGet.KQLDatabase)

	return model.setProperties(ctx, respGet.KQLDatabase)
}

func (d *dataSourceKQLDatabase) getByDisplayName(ctx context.Context, model *dataSourceKQLDatabaseModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY DISPLAY NAME", map[string]any{
		"display_name": model.DisplayName.ValueString(),
	})

	var diags diag.Diagnostics

	pager := d.client.NewListKQLDatabasesPager(model.WorkspaceID.ValueString(), nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Value {
			if *entity.DisplayName == model.DisplayName.ValueString() {
				model.set(entity)

				return nil
			}
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		fmt.Sprintf("Unable to find %s with display_name: '%s' in the Workspace ID: %s ", ItemName, model.DisplayName.ValueString(), model.WorkspaceID.ValueString()),
	)

	return diags
}
