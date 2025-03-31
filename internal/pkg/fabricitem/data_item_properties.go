// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package fabricitem

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigValidators = (*DataSourceFabricItemProperties[struct{}, struct{}])(nil)
	_ datasource.DataSourceWithConfigure        = (*DataSourceFabricItemProperties[struct{}, struct{}])(nil)
)

type DataSourceFabricItemProperties[Ttfprop, Titemprop any] struct {
	DataSourceFabricItem
	PropertiesAttributes map[string]schema.Attribute
	PropertiesSetter     func(ctx context.Context, from *Titemprop, to *DataSourceFabricItemPropertiesModel[Ttfprop, Titemprop]) diag.Diagnostics
	ItemGetter           func(ctx context.Context, fabricClient fabric.Client, model DataSourceFabricItemPropertiesModel[Ttfprop, Titemprop], fabricItem *FabricItemProperties[Titemprop]) error
	ItemListGetter       func(ctx context.Context, fabricClient fabric.Client, model DataSourceFabricItemPropertiesModel[Ttfprop, Titemprop], errNotFound fabcore.ResponseError, fabricItem *FabricItemProperties[Titemprop]) error
}

func NewDataSourceFabricItemProperties[Ttfprop, Titemprop any](config DataSourceFabricItemProperties[Ttfprop, Titemprop]) datasource.DataSource {
	return &config
}

func (d *DataSourceFabricItemProperties[Ttfprop, Titemprop]) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.TypeInfo.FullTypeName(false)
}

func (d *DataSourceFabricItemProperties[Ttfprop, Titemprop]) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = getDataSourceFabricItemPropertiesSchema(ctx, *d)
}

func (d *DataSourceFabricItemProperties[Ttfprop, Titemprop]) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
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

func (d *DataSourceFabricItemProperties[Ttfprop, Titemprop]) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	if resp.Diagnostics.Append(IsPreviewMode(d.TypeInfo.Name, d.TypeInfo.IsPreview, d.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}
}

func (d *DataSourceFabricItemProperties[Ttfprop, Titemprop]) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data DataSourceFabricItemPropertiesModel[Ttfprop, Titemprop]

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

func (d *DataSourceFabricItemProperties[Ttfprop, Titemprop]) getByID(
	ctx context.Context,
	model *DataSourceFabricItemPropertiesModel[Ttfprop, Titemprop],
) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s by ID: %s", d.TypeInfo.Name, model.ID.ValueString()))

	var fabricItem FabricItemProperties[Titemprop]

	err := d.ItemGetter(ctx, *d.pConfigData.FabricClient, *model, &fabricItem)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	model.set(fabricItem)

	return d.PropertiesSetter(ctx, fabricItem.Properties, model)
}

func (d *DataSourceFabricItemProperties[Ttfprop, Titemprop]) getByDisplayName(
	ctx context.Context,
	model *DataSourceFabricItemPropertiesModel[Ttfprop, Titemprop],
) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s by Display Name: %s", d.TypeInfo.Name, model.DisplayName.ValueString()))

	errNotFoundCode := fabcore.ErrCommon.EntityNotFound.Error()
	errNotFoundMsg := fmt.Sprintf("Unable to find %s with 'display_name': %s in the Workspace ID: %s", d.TypeInfo.Name, model.DisplayName.ValueString(), model.WorkspaceID.ValueString())

	errNotFound := fabcore.ResponseError{
		ErrorCode:  errNotFoundCode,
		StatusCode: http.StatusNotFound,
		ErrorResponse: &fabcore.ErrorResponse{
			ErrorCode: &errNotFoundCode,
			Message:   &errNotFoundMsg,
		},
	}

	var fabricItem FabricItemProperties[Titemprop]

	err := d.ItemListGetter(ctx, *d.pConfigData.FabricClient, *model, errNotFound, &fabricItem)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	model.set(fabricItem)

	return d.PropertiesSetter(ctx, fabricItem.Properties, model)
}
