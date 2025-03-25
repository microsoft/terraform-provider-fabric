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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
	pconfig "github.com/microsoft/terraform-provider-fabric/internal/provider/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSourceWithConfigValidators = (*DataSourceFabricItemDefinitionProperties[struct{}, struct{}])(nil)
	_ datasource.DataSourceWithConfigure        = (*DataSourceFabricItemDefinitionProperties[struct{}, struct{}])(nil)
)

type DataSourceFabricItemDefinitionProperties[Ttfprop, Titemprop any] struct {
	DataSourceFabricItemDefinition
	PropertiesAttributes map[string]schema.Attribute
	PropertiesSetter     func(ctx context.Context, from *Titemprop, to *DataSourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]) diag.Diagnostics
	ItemGetter           func(ctx context.Context, fabricClient fabric.Client, model DataSourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop], fabricItem *FabricItemProperties[Titemprop]) error
	ItemListGetter       func(ctx context.Context, fabricClient fabric.Client, model DataSourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop], errNotFound fabcore.ResponseError, fabricItem *FabricItemProperties[Titemprop]) error
}

func NewDataSourceFabricItemDefinitionProperties[Ttfprop, Titemprop any](config DataSourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) datasource.DataSource {
	return &config
}

func (d *DataSourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + d.TFName
}

func (d *DataSourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = getDataSourceFabricItemDefinitionPropertiesSchema(ctx, *d)
}

func (d *DataSourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
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

func (d *DataSourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
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

	if resp.Diagnostics.Append(IsPreviewMode(d.Name, d.IsPreview, d.pConfigData.Preview)...); resp.Diagnostics.HasError() {
		return
	}
}

func (d *DataSourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "READ", map[string]any{
		"action": "start",
	})

	var data DataSourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]

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

	if data.OutputDefinition.IsNull() || data.OutputDefinition.IsUnknown() {
		data.OutputDefinition = types.BoolValue(false)
	}

	if data.OutputDefinition.ValueBool() {
		if resp.Diagnostics.Append(d.getDefinition(ctx, &data)...); resp.Diagnostics.HasError() {
			return
		}

		tflog.Debug(ctx, "Definition parts content is gzip base64. Use `provider::fabric::content_decode` function to decode content.")

		resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
	}

	tflog.Debug(ctx, "READ", map[string]any{
		"action": "end",
	})

	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *DataSourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) getByID(
	ctx context.Context,
	model *DataSourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop],
) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s by ID: %s", d.Name, model.ID.ValueString()))

	var fabricItem FabricItemProperties[Titemprop]

	err := d.ItemGetter(ctx, *d.pConfigData.FabricClient, *model, &fabricItem)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	model.set(fabricItem)

	return d.PropertiesSetter(ctx, fabricItem.Properties, model)
}

func (d *DataSourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) getByDisplayName(
	ctx context.Context,
	model *DataSourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop],
) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s by Display Name: %s", d.Name, model.DisplayName.ValueString()))

	errNotFoundCode := fabcore.ErrCommon.EntityNotFound.Error()
	errNotFoundMsg := fmt.Sprintf("Unable to find %s with 'display_name': %s in the Workspace ID: %s", d.Name, model.DisplayName.ValueString(), model.WorkspaceID.ValueString())

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

func (d *DataSourceFabricItemDefinitionProperties[Ttfprop, Titemprop]) getDefinition(ctx context.Context, model *DataSourceFabricItemDefinitionPropertiesModel[Ttfprop, Titemprop]) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("getting %s definition (WorkspaceID: %s ItemID: %s)", d.Name, model.WorkspaceID.ValueString(), model.ID.ValueString()))

	respGetOpts := &fabcore.ItemsClientBeginGetItemDefinitionOptions{}

	if !model.Format.IsNull() {
		apiFormat := getDefinitionFormatAPI(d.DefinitionFormats, model.Format.ValueString())

		if apiFormat != "" {
			respGetOpts.Format = &apiFormat
		}
	}

	respGet, err := d.client.GetItemDefinition(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), respGetOpts)
	if diags := utils.GetDiagsFromError(ctx, err, utils.OperationRead, nil); diags.HasError() {
		return diags
	}

	definition, diags := getDataSourceDefinitionModel(ctx, *respGet.Definition)
	if diags.HasError() {
		return diags
	}

	model.setDefinition(definition)

	return nil
}
