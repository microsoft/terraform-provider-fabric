// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedataaccesssecurity

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

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigure        = (*dataSourceOneLakeDataAccessSecurity)(nil)
	_ datasource.DataSourceWithConfigValidators = (*dataSourceOneLakeDataAccessSecurity)(nil)
)

type dataSourceOneLakeDataAccessSecurity struct {
	pConfigData *pconfig.ProviderData
	client      *fabcore.OneLakeDataAccessSecurityClient
	TypeInfo    tftypeinfo.TFTypeInfo
}

func NewDataSourceOneLakeDataAccessSecurity() datasource.DataSource {
	return &dataSourceOneLakeDataAccessSecurity{
		TypeInfo: ItemTypeInfo,
	}
}

func (d *dataSourceOneLakeDataAccessSecurity) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

func (d *dataSourceOneLakeDataAccessSecurity) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = itemSchema(false).GetDataSource(ctx)
}

func (d *dataSourceOneLakeDataAccessSecurity) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("id"),
			path.MatchRoot("role_name"),
		),
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("role_name"),
		),
	}
}

func (d *dataSourceOneLakeDataAccessSecurity) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewOneLakeDataAccessSecurityClient()
}

func (d *dataSourceOneLakeDataAccessSecurity) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data dataSourceOneLakeDataAccessSecurityModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Read(ctx, d.pConfigData.Timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var getDiags diag.Diagnostics
	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		getDiags = d.getByID(ctx, &data)
	} else {
		getDiags = d.getByRoleName(ctx, &data)
	}

	if resp.Diagnostics.Append(getDiags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})
}

func (d *dataSourceOneLakeDataAccessSecurity) getByID(ctx context.Context, model *dataSourceOneLakeDataAccessSecurityModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY ID", map[string]any{
		"id": model.ID.ValueString(),
	})

	respList, err := d.client.ListDataAccessRoles(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	wantID := model.ID.ValueString()
	for _, role := range respList.Value {
		if role.ID != nil && *role.ID == wantID {
			return model.setFromListItem(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), role)
		}
	}

	var diags diag.Diagnostics
	diags.AddError(
		common.ErrorReadHeader,
		fmt.Sprintf("Unable to find "+d.TypeInfo.Name+" with ID: '%s' in Item ID: %s", wantID, model.ItemID.ValueString()),
	)

	return diags
}

func (d *dataSourceOneLakeDataAccessSecurity) getByRoleName(ctx context.Context, model *dataSourceOneLakeDataAccessSecurityModel) diag.Diagnostics {
	tflog.Trace(ctx, "GET BY ROLE NAME", map[string]any{
		"role_name": model.RoleName.ValueString(),
	})

	respList, err := d.client.ListDataAccessRoles(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), nil)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationList, nil); diags.HasError() {
		return diags
	}

	roleName := model.RoleName.ValueString()
	for _, role := range respList.Value {
		if role.Name != nil && *role.Name == roleName {
			return model.setFromListItem(ctx, model.WorkspaceID.ValueString(), model.ItemID.ValueString(), role)
		}
	}

	var diags diag.Diagnostics
	diags.AddError(
		common.ErrorReadHeader,
		fmt.Sprintf("Unable to find "+d.TypeInfo.Name+" with name: '%s' in Item ID: %s", roleName, model.ItemID.ValueString()),
	)

	return diags
}
