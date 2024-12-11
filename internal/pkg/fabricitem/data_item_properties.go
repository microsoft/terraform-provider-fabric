// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"
	"fmt"

	azto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigValidators = (*DataSourceFabricItemProperties[struct{}, struct{}])(nil)
	_ datasource.DataSourceWithConfigure        = (*DataSourceFabricItemProperties[struct{}, struct{}])(nil)
)

type DataSourceFabricItemProperties[T any, Tm any] struct {
	pConfigData         *pconfig.ProviderData
	client              *fabcore.ItemsClient
	Type                fabcore.ItemType
	Name                string
	TFName              string
	MarkdownDescription string
	IsDisplayNameUnique bool
	PropertiesSchema    schema.SingleNestedAttribute
	PropertiesSetter    func(ctx context.Context, from *Tm, to *DataSourceFabricItemPropertiesModel[T, Tm]) diag.Diagnostics
	ItemGetter          func(ctx context.Context, fabClient fabric.Client, model DataSourceFabricItemPropertiesModel[T, Tm], fabItem *FabricItem[Tm]) diag.Diagnostics
}

func NewDataSourceFabricItemProperties[T any, Tm any](config DataSourceFabricItemProperties[T, Tm]) datasource.DataSource {
	return &config
}

func (d *DataSourceFabricItemProperties[T, Tm]) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.TFName
}

func (d *DataSourceFabricItemProperties[T, Tm]) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = GetDataSourceFabricItemPropertiesSchema1[T](ctx, *d)
}

func (d *DataSourceFabricItemProperties[T, Tm]) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	if d.IsDisplayNameUnique {
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

	return []datasource.ConfigValidator{}
}

func (d *DataSourceFabricItemProperties[T, Tm]) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = fabcore.NewClientFactoryWithClient(*pConfigData.FabricClient).NewItemsClient()
}

func (d *DataSourceFabricItemProperties[T, Tm]) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})
	tflog.Trace(ctx, "READ", map[string]any{
		"config": req.Config,
	})

	var data DataSourceFabricItemPropertiesModel[T, Tm]

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

func (d *DataSourceFabricItemProperties[T, Tm]) getByID(ctx context.Context, model *DataSourceFabricItemPropertiesModel[T, Tm]) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s by ID: %s", d.Name, model.ID.ValueString()))

	var fabItem FabricItem[Tm]

	if d.ItemGetter != nil {
		diags := d.ItemGetter(ctx, *d.pConfigData.FabricClient, *model, &fabItem)
		if diags.HasError() {
			return diags
		}
	} else {
		respGet, err := d.client.GetItem(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
			return diags
		}

		fabItem.Set(respGet.Item)
	}

	model.set1(fabItem)

	// usee PropertiesSetter to set the properties of the model
	if d.PropertiesSetter != nil {
		diags := d.PropertiesSetter(ctx, fabItem.Properties, model)
		if diags.HasError() {
			return diags
		}
	}

	return nil
}

func (d *DataSourceFabricItemProperties[T, Tm]) getByDisplayName(ctx context.Context, model *DataSourceFabricItemPropertiesModel[T, Tm]) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s by Display Name: %s", d.Name, model.DisplayName.ValueString()))

	var diags diag.Diagnostics

	opts := &fabcore.ItemsClientListItemsOptions{
		Type: azto.Ptr(string(d.Type)),
	}

	pager := d.client.NewListItemsPager(model.WorkspaceID.ValueString(), opts)
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
		fmt.Sprintf("Unable to find %s with 'display_name': %s in the Workspace ID: %s", d.Name, model.DisplayName.ValueString(), model.WorkspaceID.ValueString()),
	)

	return diags
}
