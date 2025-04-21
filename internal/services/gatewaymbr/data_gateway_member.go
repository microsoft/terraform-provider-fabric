// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package gatewaymbr

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

var (
	_ datasource.DataSourceWithConfigValidators = (*dataSourceGatewayMember)(nil)
	_ datasource.DataSourceWithConfigure        = (*dataSourceGatewayMember)(nil)
)

type dataSourceGatewayMember struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.GatewaysClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceGatewayMember() datasource.DataSource {
	return &dataSourceGatewayMember{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceGatewayMember) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

func (d *dataSourceGatewayMember) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = itemSchema(false).GetDataSource(ctx)
}

func (d *dataSourceGatewayMember) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
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

func (d *dataSourceGatewayMember) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewGatewaysClient()
}

func (d *dataSourceGatewayMember) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceGatewayMemberModel

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

func (d *dataSourceGatewayMember) getByID(ctx context.Context, model *dataSourceGatewayMemberModel) diag.Diagnostics {
	return d.get(ctx, true, model)
}

func (d *dataSourceGatewayMember) getByDisplayName(ctx context.Context, model *dataSourceGatewayMemberModel) diag.Diagnostics {
	return d.get(ctx, false, model)
}

func (d *dataSourceGatewayMember) get(ctx context.Context, byID bool, model *dataSourceGatewayMemberModel) diag.Diagnostics {
	var diags diag.Diagnostics
	var notFound string

	respList, err := d.client.ListGatewayMembers(ctx, model.GatewayID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	for _, entity := range respList.Value {
		switch byID {
		case true:
			if *entity.ID == model.ID.ValueString() {
				return model.set(ctx, model.GatewayID.ValueString(), entity)
			}

			notFound = "Unable to find " + d.TypeInfo.Name + " with 'id': " + model.ID.ValueString()
		default:
			if *entity.DisplayName == model.DisplayName.ValueString() {
				return model.set(ctx, model.GatewayID.ValueString(), entity)
			}

			notFound = "Unable to find " + d.TypeInfo.Name + " with 'display_name': " + model.DisplayName.ValueString()
		}
	}

	diags.AddError(
		common.ErrorReadHeader,
		notFound,
	)

	return diags
}
