// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var (
	_ datasource.DataSourceWithConfigure        = (*dataSourceTag)(nil)
	_ datasource.DataSourceWithConfigValidators = (*dataSourceTag)(nil)
)

type dataSourceTag struct {
	pConfigData *pconfig.ProviderData
	client      *fabadmin.TagsClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceTag() datasource.DataSource {
	return &dataSourceTag{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceTag) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

func (d *dataSourceTag) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = itemSchema(false).GetDataSource(ctx)
}

func (d *dataSourceTag) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = fabadmin.NewClientFactoryWithClient(*pConfigData.FabricClient).NewTagsClient()
}

func (d *dataSourceTag) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
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

func (d *dataSourceTag) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceTagModel

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

func (d *dataSourceTag) getByID(ctx context.Context, model *dataSourceTagModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY ID", map[string]any{
		"id": model.ID.ValueString(),
	})

	var diags diag.Diagnostics

	pager := d.client.NewListTagsPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Value {
			if *entity.ID == model.ID.ValueString() {
				return model.set(ctx, entity)
			}
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		"Unable to find Tag with 'display_name': "+model.DisplayName.ValueString(),
	)

	return diags
}

func (d *dataSourceTag) getByDisplayName(ctx context.Context, model *dataSourceTagModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY DISPLAY NAME", map[string]any{
		"display_name": model.DisplayName.ValueString(),
	})

	var diags diag.Diagnostics

	pager := d.client.NewListTagsPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
			return diags
		}

		for _, entity := range page.Value {
			if *entity.DisplayName == model.DisplayName.ValueString() {
				model.ID = customtypes.NewUUIDPointerValue(entity.ID)

				return d.getByID(ctx, model)
			}
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		"Unable to find Tag with 'display_name': "+model.DisplayName.ValueString(),
	)

	return diags
}
