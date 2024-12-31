// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package connection

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigValidators = (*dataSourceConnection)(nil)
	_ datasource.DataSourceWithConfigure        = (*dataSourceConnection)(nil)
)

type dataSourceConnection struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.ConnectionsClient
}

func NewDataSourceConnection() datasource.DataSource {
	return &dataSourceConnection{}
}

func (d *dataSourceConnection) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ItemTFName
}

func (d *dataSourceConnection) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attributes := getDataSourceConnectionAttributes(ctx)
	attributes["timeouts"] = timeouts.Attributes(ctx)

	resp.Schema = schema.Schema{
		MarkdownDescription: "Get a Fabric " + ItemName + ".\n\n" +
			"Use this data source to fetch a [" + ItemName + "](" + ItemDocsURL + ").\n\n" +
			ItemDocsSPNSupport,
		Attributes: attributes,
	}
}

func (d *dataSourceConnection) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
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

func (d *dataSourceConnection) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewConnectionsClient()
}

func (d *dataSourceConnection) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data dataSourceConnectionModel

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

func (d *dataSourceConnection) getByID(ctx context.Context, model *dataSourceConnectionModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY ID", map[string]any{
		"id": model.ID.ValueString(),
	})

	respGet, err := d.client.GetConnection(ctx, model.ID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	model.set(respGet.Connection)

	if diags := model.setConnectionDetails(ctx, respGet.ConnectionDetails); diags.HasError() {
		return diags
	}

	if diags := model.setCredentialDetails(ctx, respGet.CredentialDetails); diags.HasError() {
		return diags
	}

	return nil
}

func (d *dataSourceConnection) getByDisplayName(ctx context.Context, model *dataSourceConnectionModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY DISPLAY NAME", map[string]any{
		"display_name": model.DisplayName.ValueString(),
	})

	var diags diag.Diagnostics

	pager := d.client.NewListConnectionsPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Value {
			if *entity.DisplayName == model.DisplayName.ValueString() {
				model.set(entity)

				if diags := model.setConnectionDetails(ctx, entity.ConnectionDetails); diags.HasError() {
					return diags
				}

				if diags := model.setCredentialDetails(ctx, entity.CredentialDetails); diags.HasError() {
					return diags
				}

				return nil
			}
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		fmt.Sprintf("Unable to find %s with display_name: '%s'", ItemName, model.DisplayName.ValueString()),
	)

	return diags
}
